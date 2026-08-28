package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"raven.go.invoice-builder/internal/models"
)

func (h *ApiHandler) AdminGetPackages(w http.ResponseWriter, r *http.Request) {
	includeInactive := true
	packages, err := h.Services.Admin.GetCreditPackages(r.Context(), includeInactive)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "packages": packages}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminCreatePackage(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string  `json:"name"`
		Slug           string  `json:"slug"`
		Description    string  `json:"description"`
		CreditsGranted int     `json:"credits_granted"`
		Price          float64 `json:"price"` // frontend sends 50.00, backend needs 5000
		Currency       string  `json:"currency"`
		BadgeTag       string  `json:"badge_tag"`
		IsActive       bool    `json:"is_active"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	pkg := &models.CreditPackage{
		Name:           input.Name,
		Slug:           input.Slug,
		Description:    input.Description,
		CreditsGranted: input.CreditsGranted,
		Price:          int64(input.Price * 100),
		Currency:       input.Currency,
		BadgeTag:       input.BadgeTag,
		IsActive:       input.IsActive,
	}

	err = h.Services.Admin.CreateCreditPackage(r.Context(), pkg)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusCreated, Envelope{"status": "success", "package": pkg}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminUpdatePackage(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID             uuid.UUID `json:"id"`
		Name           string    `json:"name"`
		Slug           string    `json:"slug"`
		Description    string    `json:"description"`
		CreditsGranted int       `json:"credits_granted"`
		Price          float64   `json:"price"`
		Currency       string    `json:"currency"`
		BadgeTag       string    `json:"badge_tag"`
		IsActive       bool      `json:"is_active"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	pkg := &models.CreditPackage{
		ID:             input.ID,
		Name:           input.Name,
		Slug:           input.Slug,
		Description:    input.Description,
		CreditsGranted: input.CreditsGranted,
		Price:          int64(input.Price * 100),
		Currency:       input.Currency,
		BadgeTag:       input.BadgeTag,
		IsActive:       input.IsActive,
	}

	err = h.Services.Admin.UpdateCreditPackage(r.Context(), pkg)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFoundResponse(w, r)
			return
		}
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "package": pkg}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminDeletePackage(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	err = h.Services.Admin.DeleteCreditPackage(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFoundResponse(w, r)
			return
		}
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "message": "package deleted"}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
