package main

import (
	"net/http"
)

func (h *ApiHandler) AdminGetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	totalUsers, activeUsers, unverifiedUsers, totalRevenue, totalPayments, totalPurchasedCredits, totalUsedCredits, recentUsers, recentLogs, err := h.Services.Admin.GetDashboardStats(ctx)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status": "success",
		"stats": Envelope{
			"total_revenue":          totalRevenue,
			"total_payments":         totalPayments,
			"total_users":            totalUsers,
			"active_users":           activeUsers,
			"unverified_users":       unverifiedUsers,
			"total_purchased_credits": totalPurchasedCredits,
			"total_used_credits":      totalUsedCredits,
		},
		"recent_logs":  recentLogs,
		"recent_users": recentUsers,
	}, nil)

	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) AdminGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page, limit := h.ParsePagination(r)
	search := r.URL.Query().Get("search")

	users, totalCount, err := h.Services.Admin.GetAllUsers(ctx, search, page, limit)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status": "success",
		"users":  users,
		"meta": Envelope{
			"total_count": totalCount,
			"page":        page,
			"limit":       limit,
		},
	}, nil)

	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
