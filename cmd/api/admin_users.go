package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

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

	profile, _ := h.Services.Auth.GetBusinessProfile(r.Context(), user.ID)
	credits, _ := h.Services.Credit.GetBalance(r.Context(), user.ID)

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status": "success", 
		"user": Envelope{
			"id":                  user.ID,
			"name":                user.Name,
			"email":               user.Email,
			"is_profile_complete": user.IsProfileComplete,
			"is_activated":        user.IsActivated,
			"created_at":          user.CreatedAt,
			"role":                user.Role,
			"credits":             credits,
			"profile":             profile,
		},
	}, nil)
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

func (h *ApiHandler) AdminForcePasswordReset(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID uuid.UUID `json:"user_id"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.Services.Auth.GetUserByID(r.Context(), input.UserID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFoundResponse(w, r)
			return
		}
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.Services.Auth.RequestPasswordReset(r.Context(), user.Email)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "message": "Password reset email dispatched successfully"}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminImpersonateUser(w http.ResponseWriter, r *http.Request) {
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

	roleName := "user"
	if user.Role.Name != "" {
		roleName = user.Role.Name
	}

	accessToken, err := h.JWTManager.GenerateToken(user.ID, user.Email, roleName, user.TokenVersion)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	refreshToken, err := h.Models.RefreshTokens.Create(r.Context(), user.ID, 30*24*time.Hour)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status": "success",
		"tokens": Envelope{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
		"message": "Impersonation session created",
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminExportUsersCSV(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	users, _, err := h.Services.Admin.GetAllUsers(r.Context(), search, 1, 100000)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"users_export.csv\"")

	header := "ID,Name,Email,Role,IsActivated,CreatedAt\n"
	_, _ = w.Write([]byte(header))

	for _, u := range users {
		roleName := "user"
		if u.Role.Name != "" {
			roleName = u.Role.Name
		}
		
		name := strings.ReplaceAll(u.Name, "\"", "\"\"")
		email := strings.ReplaceAll(u.Email, "\"", "\"\"")
		dateStr := u.CreatedAt.Format("2006-01-02 15:04:05")

		row := fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%t\",\"%s\"\n",
			u.ID.String(), name, email, roleName, u.IsActivated, dateStr)
		_, _ = w.Write([]byte(row))
	}
}
