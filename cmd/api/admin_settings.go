package main

import (
	"net/http"
)

func (h *ApiHandler) AdminGetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.Services.Admin.GetSystemSettings(r.Context())
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "settings": settings}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminUpdateSettings(w http.ResponseWriter, r *http.Request) {
	var input struct {
		DefaultSignupBonus  string `json:"default_signup_credits"`
		MaintenanceMode     string `json:"maintenance_mode"`
		SupportContactEmail string `json:"support_email"`
		LegalTerms          string `json:"legal_terms"`
		LegalPrivacy        string `json:"legal_privacy"`
		LegalRefund         string `json:"legal_refund"`
		LegalSecurity       string `json:"legal_security"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	err = h.Services.Admin.UpdateSystemSettings(
		r.Context(),
		input.DefaultSignupBonus,
		input.MaintenanceMode,
		input.SupportContactEmail,
		input.LegalTerms,
		input.LegalPrivacy,
		input.LegalRefund,
		input.LegalSecurity,
	)

	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	_ = h.WriteJSON(w, http.StatusOK, Envelope{"status": "success", "message": "settings updated"}, nil)
}

func (h *ApiHandler) GetPublicSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.Services.Admin.GetSystemSettings(r.Context())
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	_ = h.WriteJSON(w, http.StatusOK, Envelope{
		"status": "success",
		"settings": map[string]any{
			"legal_terms":    settings["LegalTerms"],
			"legal_privacy":  settings["LegalPrivacy"],
			"legal_refund":   settings["LegalRefund"],
			"legal_security": settings["LegalSecurity"],
			"support_email":  settings["SupportContactEmail"],
		},
	}, nil)
}
