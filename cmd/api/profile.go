package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"raven.go.invoice-builder/internal/currency"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/validator"
)

func (h *ApiHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	profile, err := h.Services.Auth.GetBusinessProfile(r.Context(), user.ID)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":     "success",
		"profile":    profile,
		"currencies": currency.SupportedCurrencies(),
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		Name               string `json:"name"`
		Role               string `json:"role"`
		CompanyName        string `json:"company_name"`
		Address            string `json:"address"`
		TaxID              string `json:"tax_id"`
		RegistrationNumber string `json:"registration_number"`
		RegistrationDate   string `json:"registration_date"`
		BusinessType       string `json:"business_type"`
		RegisteredAddress  string `json:"registered_address"`
		DefaultCurrency    string `json:"default_currency"`
		LogoURL            string `json:"logo_url"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.CheckField(validator.NotBlank(input.CompanyName), "company_name", "Company Name is required")

	if !v.Valid() {
		h.FailedValidationResponse(w, r, v.FieldErrors)
		return
	}

	profile := &models.BusinessProfile{
		UserID:             user.ID,
		Role:               strings.TrimSpace(input.Role),
		CompanyName:        strings.TrimSpace(input.CompanyName),
		Address:            strings.TrimSpace(input.Address),
		TaxID:              strings.TrimSpace(input.TaxID),
		RegistrationNumber: strings.TrimSpace(input.RegistrationNumber),
		RegistrationDate:   strings.TrimSpace(input.RegistrationDate),
		BusinessType:       strings.TrimSpace(input.BusinessType),
		RegisteredAddress:  strings.TrimSpace(input.RegisteredAddress),
		DefaultCurrency:    strings.TrimSpace(input.DefaultCurrency),
		LogoURL:            strings.TrimSpace(input.LogoURL),
	}

	existingProfile, _ := h.Services.Auth.GetBusinessProfile(r.Context(), user.ID)
	if existingProfile != nil {
		profile.ID = existingProfile.ID
	}

	// Update user name if changed
	if strings.TrimSpace(input.Name) != "" && input.Name != user.Name {
		user.Name = strings.TrimSpace(input.Name)
	}
	user.IsProfileComplete = true

	// Persist the updates
	err = h.Services.Auth.UpdateBusinessProfile(r.Context(), profile)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.Services.Auth.UpdateUser(r.Context(), user)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":  "success",
		"message": "Profile updated successfully",
		"profile": profile,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) UploadProfileLogo(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		h.BadRequestResponse(w, r, ErrMultipartParseFailed)
		return
	}

	file, header, err := r.FormFile("logo")
	if err != nil {
		h.BadRequestResponse(w, r, ErrLogoRequired)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		h.BadRequestResponse(w, r, ErrInvalidImageFormat)
		return
	}

	uploadDir := "./ui/asset/img/uploads"
	_ = os.MkdirAll(uploadDir, 0755)

	filename := fmt.Sprintf("logo_%s_%s%s", user.ID.String(), os.Getenv("APP_PORT"), ext)
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	logoURL := fmt.Sprintf("/asset/img/uploads/%s", filename)

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":   "success",
		"logo_url": logoURL,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
