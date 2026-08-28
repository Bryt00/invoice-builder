package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/paystack"
)

func (h *ApiHandler) GetCreditsBalance(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	balance, totalPurchased, totalUsed, err := h.Services.Credit.GetCreditStats(r.Context(), user.ID)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status": "success",
		"stats": Envelope{
			"balance":         balance,
			"total_purchased": totalPurchased,
			"total_used":      totalUsed,
		},
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) GetCreditsHistory(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	history, err := h.Services.Credit.GetCreditHistory(r.Context(), user.ID)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":  "success",
		"history": history,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) GetCreditsPackages(w http.ResponseWriter, r *http.Request) {
	packages, err := h.Services.Credit.GetAvailablePackages(r.Context())
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":   "success",
		"packages": packages,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) InitializeCreditsTopup(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		PackageID string `json:"package_id"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	pkgID, err := uuid.Parse(input.PackageID)
	if err != nil {
		h.BadRequestResponse(w, r, models.ErrInvalidPackageID)
		return
	}

	pkg, err := h.Models.CreditPackages.GetByID(r.Context(), pkgID)
	if err != nil {
		h.NotFoundResponse(w, r)
		return
	}

	// Generate a unique reference
	reference := fmt.Sprintf("TOPUP_%s_%d", user.ID.String()[:8], time.Now().Unix())

	initReq := &paystack.InitializeRequest{
		Email:       user.Email,
		Amount:      int(pkg.Price), // Already in sub-units? Yes, according to model "Price in smallest subunit"
		Reference:   reference,
		CallbackURL: fmt.Sprintf("%s/user/dashboard", h.FrontendURL),
		Metadata: map[string]any{
			"user_id":    user.ID,
			"package_id": pkg.ID,
			"type":       "credit_topup",
		},
	}

	paystackRes, err := h.Services.Paystack.InitializeTransaction(initReq)
	if err != nil {
		h.ServerErrorResponse(w, r, fmt.Errorf("failed to initialize paystack: %w", err))
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":            "success",
		"authorization_url": paystackRes.Data.AuthorizationURL,
		"reference":         paystackRes.Data.Reference,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) VerifyCreditsTopup(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	reference := r.URL.Query().Get("reference")
	if reference == "" {
		h.BadRequestResponse(w, r, ErrMissingReference)
		return
	}

	paystackRes, err := h.Services.Paystack.VerifyTransaction(reference)
	if err != nil {
		h.ServerErrorResponse(w, r, fmt.Errorf("failed to verify paystack transaction: %w", err))
		return
	}

	if paystackRes.Data.Status != "success" {
		h.BadRequestResponse(w, r, ErrTransactionNotSuccessful)
		return
	}

	pkgIDStr, ok := paystackRes.Data.Metadata["package_id"].(string)
	if !ok {
		h.ServerErrorResponse(w, r, ErrMissingPackageID)
		return
	}

	pkgID, err := uuid.Parse(pkgIDStr)
	if err != nil {
		h.ServerErrorResponse(w, r, fmt.Errorf("invalid package_id in metadata: %w", err))
		return
	}

	err = h.Services.Credit.FulfillTopup(r.Context(), user.ID, pkgID, reference)
	if err != nil {
		if errors.Is(err, models.ErrTransactionAlreadyProcessed) {
			// Idempotent success
			err = h.WriteJSON(w, http.StatusOK, Envelope{
				"status":  "success",
				"message": "Payment verified successfully",
			}, nil)
			if err != nil {
				h.ServerErrorResponse(w, r, err)
			}
			return
		}
		h.ServerErrorResponse(w, r, fmt.Errorf("failed to fulfill topup: %w", err))
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":  "success",
		"message": "Payment verified successfully",
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
