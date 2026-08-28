package main

import (
	"math"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/validator"
)

func (h *ApiHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	page, limit := h.ParsePagination(r)

	clients, total, err := h.Services.Client.GetAllByUserID(r.Context(), user.ID, page, limit)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":  "success",
		"clients": clients,
		"meta": Envelope{
			"current_page":  page,
			"page_size":     limit,
			"total_records": total,
			"total_pages":   totalPages,
		},
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Company string `json:"company"`
		Address string `json:"address"`
		TaxID   string `json:"tax_id"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(input.Email)

	v := validator.New()
	v.CheckField(validator.NotBlank(name), "name", "Client Name is required")
	v.CheckField(validator.NotBlank(email), "email", "Client Email is required")
	v.CheckField(validator.Matches(email, validator.EmailRegexp), "email", "Must be a valid email address")

	if !v.Valid() {
		h.FailedValidationResponse(w, r, v.FieldErrors)
		return
	}

	client := &models.Client{
		ID:      uuid.New(),
		UserID:  user.ID,
		Name:    name,
		Email:   email,
		Phone:   strings.TrimSpace(input.Phone),
		Company: strings.TrimSpace(input.Company),
		Address: strings.TrimSpace(input.Address),
		TaxID:   strings.TrimSpace(input.TaxID),
	}

	err = h.Services.Client.CreateClient(r.Context(), client)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusCreated, Envelope{
		"status": "success",
		"client": client,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Company string `json:"company"`
		Address string `json:"address"`
		TaxID   string `json:"tax_id"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	clientID, err := uuid.Parse(input.ID)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	client, err := h.Services.Client.GetByID(r.Context(), clientID, user.ID)
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	if input.Name != "" {
		client.Name = strings.TrimSpace(input.Name)
	}
	if input.Email != "" {
		client.Email = strings.TrimSpace(input.Email)
	}
	client.Phone = strings.TrimSpace(input.Phone)
	client.Company = strings.TrimSpace(input.Company)
	client.Address = strings.TrimSpace(input.Address)
	client.TaxID = strings.TrimSpace(input.TaxID)

	v := validator.New()
	v.CheckField(validator.NotBlank(client.Name), "name", "Client Name is required")
	v.CheckField(validator.NotBlank(client.Email), "email", "Client Email is required")
	v.CheckField(validator.Matches(client.Email, validator.EmailRegexp), "email", "Must be a valid email address")

	if !v.Valid() {
		h.FailedValidationResponse(w, r, v.FieldErrors)
		return
	}

	err = h.Services.Client.UpdateClient(r.Context(), client)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status": "success",
		"client": client,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		ID string `json:"id"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	clientID, err := uuid.Parse(input.ID)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	err = h.Services.Client.DeleteClient(r.Context(), clientID, user.ID)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":  "success",
		"message": "Client deleted successfully",
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
