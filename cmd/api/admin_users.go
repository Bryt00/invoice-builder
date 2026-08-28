package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"raven.go.invoice-builder/internal/models"
)


func (h *ApiHandler) AdminCreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		RoleName    string `json:"role"`
		IsActivated bool   `json:"is_activated"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.Services.Admin.CreateUser(r.Context(), input.Name, input.Email, input.Password, input.RoleName, input.IsActivated)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			h.FailedValidationResponse(w, r, map[string]string{"email": "a user with this email address already exists"})
			return
		}
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusCreated, Envelope{"status": "success", "user": user}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminGetUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	user, err := h.Services.Auth.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFoundResponse(w, r)
			return
		}
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "user": user}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		RoleName    string    `json:"role"`
		IsActivated bool      `json:"is_activated"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.Services.Admin.UpdateUser(r.Context(), input.ID, input.Name, input.Email, input.RoleName, input.IsActivated)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFoundResponse(w, r)
			return
		}
		if errors.Is(err, models.ErrDuplicateEmail) {
			h.FailedValidationResponse(w, r, map[string]string{"email": "a user with this email address already exists"})
			return
		}
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "user": user}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	err = h.Services.Admin.DeleteUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFoundResponse(w, r)
			return
		}
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "message": "user deleted"}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminGrantCredits(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID uuid.UUID `json:"user_id"`
		Amount int       `json:"amount"`
		Reason string    `json:"reason"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	err = h.Services.Credit.AdminGrantCredits(r.Context(), input.UserID, input.Amount, input.Reason)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFoundResponse(w, r)
			return
		}
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "message": "credits updated successfully"}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
