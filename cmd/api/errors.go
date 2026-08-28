package main

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidID            = errors.New("invalid ID parameter")
	ErrMissingID            = errors.New("id parameter is required")
	ErrMissingToken         = errors.New("token parameter is required")
	ErrInvalidToken         = errors.New("invalid authentication token")
	ErrBadJSON              = errors.New("body contains badly-formed JSON")
	ErrEmptyJSON            = errors.New("body must not be empty")
	ErrMultipleJSONValues   = errors.New("body must only contain a single JSON value")
	ErrInvalidCategoryID    = errors.New("invalid category_id format")
	ErrMultipartParseFailed = errors.New("failed to parse multipart form")
	ErrLogoRequired         = errors.New("logo file is required")
	ErrInvalidImageFormat   = errors.New("only JPG, PNG and WEBP images are allowed")
	ErrInvalidClientID      = errors.New("invalid client_id")
	ErrMissingReference     = errors.New("missing reference parameter")
	ErrFailedToVerify       = errors.New("failed to verify paystack transaction")
	ErrTransactionNotSuccessful = errors.New("transaction was not successful")
	ErrMissingPackageID     = errors.New("transaction missing package_id in metadata")
)
	// h.BadRequestResponse(w, r, errors.New("missing reference parameter"))
	// 	return
	// }

	// paystackRes, err := h.Services.Paystack.VerifyTransaction(reference)
	// if err != nil {
	// 	h.ServerErrorResponse(w, r, fmt.Errorf("failed to verify paystack transaction: %w", err))
	// 	return
	// }

	// if paystackRes.Data.Status != "success" {
	// 	h.BadRequestResponse(w, r, errors.New("transaction was not successful"))
	// 	return
	// }

	// pkgIDStr, ok := paystackRes.Data.Metadata["package_id"].(string)
	// if !ok {
	// 	h.ServerErrorResponse(w, r, errors.New("transaction missing package_id in metadata"))

func (h *ApiHandler) LogApiError(r *http.Request, err error) {
	if h.JsonLogger != nil {
		h.JsonLogger.PrintError(err, map[string]string{
			"request_method": r.Method,
			"request_url":    r.URL.String(),
		})
	}
}

func (h *ApiHandler) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"status": "error", "error": message}
	err := h.WriteJSON(w, status, env, nil)
	if err != nil {
		h.LogApiError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ApiHandler) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.LogApiError(r, err)
	message := "the server encountered a problem and could not process your request"
	h.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (h *ApiHandler) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	h.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (h *ApiHandler) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (h *ApiHandler) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	h.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (h *ApiHandler) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	h.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (h *ApiHandler) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	h.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (h *ApiHandler) ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	message := "you do not have permission to access this resource"
	h.ErrorResponse(w, r, http.StatusForbidden, message)
}
