package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/internal/validator"
)

func (h *ApiHandler) GetFinanceTransactions(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	txnType := r.URL.Query().Get("type")
	categoryID := r.URL.Query().Get("category_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	search := r.URL.Query().Get("search")

	transactions, err := h.Services.Finance.GetAllTransactions(r.Context(), user.ID, txnType, categoryID, startDate, endDate, search)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":       "success",
		"transactions": transactions,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) GetFinanceSummary(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	totalIncome, totalExpenses, netProfit, err := h.Services.Finance.GetFinancialStats(r.Context(), user.ID, startDate, endDate)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status": "success",
		"summary": Envelope{
			"total_income":   totalIncome,
			"total_expenses": totalExpenses,
			"net_profit":     netProfit,
		},
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) CreateFinanceTransaction(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		Title           string  `json:"title"`
		Amount          float64 `json:"amount"`
		Type            string  `json:"type"`
		CategoryID      string  `json:"category_id"`
		Currency        string  `json:"currency"`
		PayeeOrPayer    string  `json:"payee_or_payer"`
		TransactionDate string  `json:"transaction_date"`
		Description     string  `json:"description"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.CheckField(validator.NotBlank(input.Title), "title", "Title is required")
	v.CheckField(input.Amount > 0, "amount", "Amount must be greater than 0")
	v.CheckField(input.Type == "income" || input.Type == "expense", "type", "Type must be 'income' or 'expense'")
	v.CheckField(validator.NotBlank(input.CategoryID), "category_id", "Category ID is required")

	if !v.Valid() {
		h.FailedValidationResponse(w, r, v.FieldErrors)
		return
	}

	categoryUUID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		h.BadRequestResponse(w, r, ErrInvalidCategoryID)
		return
	}

	txnDate := time.Now()
	if input.TransactionDate != "" {
		if t, parseErr := time.Parse("2006-01-02", input.TransactionDate); parseErr == nil {
			txnDate = t
		}
	}

	txn := &models.FinancialTransaction{
		ID:              uuid.New(),
		UserID:          user.ID,
		CategoryID:      categoryUUID,
		Type:            input.Type,
		Amount:          input.Amount,
		Currency:        strings.TrimSpace(input.Currency),
		Title:           strings.TrimSpace(input.Title),
		Description:     strings.TrimSpace(input.Description),
		TransactionDate: txnDate,
		PayeeOrPayer:    strings.TrimSpace(input.PayeeOrPayer),
		Status:          "completed",
	}

	err = h.Services.Finance.CreateTransaction(r.Context(), txn)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusCreated, Envelope{
		"status":      "success",
		"transaction": txn,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) GetFinanceCategories(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	categories, err := h.Services.Finance.GetCategories(r.Context(), user.ID)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":     "success",
		"categories": categories,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) CreateFinanceCategory(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Icon string `json:"icon"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.CheckField(validator.NotBlank(input.Name), "name", "Category name is required")
	v.CheckField(input.Type == "income" || input.Type == "expense", "type", "Type must be 'income' or 'expense'")

	if !v.Valid() {
		h.FailedValidationResponse(w, r, v.FieldErrors)
		return
	}

	cat, err := h.Services.Finance.CreateCategory(r.Context(), user.ID, input.Name, input.Type, "primary", input.Icon)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusCreated, Envelope{
		"status":   "success",
		"category": cat,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) DeleteFinanceTransaction(w http.ResponseWriter, r *http.Request) {
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

	txnID, err := uuid.Parse(input.ID)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	err = h.Services.Finance.DeleteTransaction(r.Context(), txnID, user.ID)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":  "success",
		"message": "Transaction deleted",
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) ExportFinanceTransactionsCSV(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	txnType := r.URL.Query().Get("type")
	categoryID := r.URL.Query().Get("category")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	search := r.URL.Query().Get("search")

	txns, err := h.Services.Finance.GetAllTransactions(r.Context(), user.ID, txnType, categoryID, startDate, endDate, search)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"finance_export.csv\"")

	header := "Date,Type,Category,Title,Amount,Currency,Payee/Payer,Status\n"
	_, _ = w.Write([]byte(header))

	for _, t := range txns {
		catName := "Uncategorized"
		if t.Category != nil {
			catName = t.Category.Name
		}
		
		dateStr := t.TransactionDate.Format("2006-01-02")
		
		title := strings.ReplaceAll(t.Title, "\"", "\"\"")
		payee := strings.ReplaceAll(t.PayeeOrPayer, "\"", "\"\"")
		catName = strings.ReplaceAll(catName, "\"", "\"\"")
		
		row := fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%.2f\",\"%s\",\"%s\",\"%s\"\n",
			dateStr, t.Type, catName, title, t.Amount, t.Currency, payee, t.Status)
		_, _ = w.Write([]byte(row))
	}
}

func (h *ApiHandler) UpdateFinanceTransaction(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		ID              string  `json:"id"`
		Title           string  `json:"title"`
		Amount          float64 `json:"amount"`
		Type            string  `json:"type"`
		CategoryID      string  `json:"category_id"`
		Currency        string  `json:"currency"`
		PayeeOrPayer    string  `json:"payee_or_payer"`
		TransactionDate string  `json:"transaction_date"`
		Description     string  `json:"description"`
	}

	err := h.ReadJSON(w, r, &input)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.CheckField(validator.NotBlank(input.ID), "id", "Transaction ID is required")
	v.CheckField(validator.NotBlank(input.Title), "title", "Title is required")
	v.CheckField(input.Amount > 0, "amount", "Amount must be greater than 0")
	v.CheckField(input.Type == "income" || input.Type == "expense", "type", "Type must be 'income' or 'expense'")
	v.CheckField(validator.NotBlank(input.CategoryID), "category_id", "Category ID is required")

	if !v.Valid() {
		h.FailedValidationResponse(w, r, v.FieldErrors)
		return
	}

	txnUUID, err := uuid.Parse(input.ID)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	categoryUUID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		h.BadRequestResponse(w, r, ErrInvalidCategoryID)
		return
	}

	txn, err := h.Services.Finance.GetTransactionByID(r.Context(), txnUUID, user.ID)
	if err != nil {
		if err == models.ErrNoRecord {
			h.NotFoundResponse(w, r)
			return
		}
		h.ServerErrorResponse(w, r, err)
		return
	}

	txnDate := txn.TransactionDate
	if input.TransactionDate != "" {
		if t, parseErr := time.Parse("2006-01-02", input.TransactionDate); parseErr == nil {
			txnDate = t
		}
	}

	txn.Title = strings.TrimSpace(input.Title)
	txn.Amount = input.Amount
	txn.Type = input.Type
	txn.CategoryID = categoryUUID
	if strings.TrimSpace(input.Currency) != "" {
		txn.Currency = strings.TrimSpace(input.Currency)
	}
	txn.PayeeOrPayer = strings.TrimSpace(input.PayeeOrPayer)
	txn.TransactionDate = txnDate
	txn.Description = strings.TrimSpace(input.Description)

	err = h.Services.Finance.UpdateTransaction(r.Context(), txn)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status":      "success",
		"transaction": txn,
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ApiHandler) GetFinanceAnalytics(w http.ResponseWriter, r *http.Request) {
	user := h.ContextGetUser(r)
	if user == nil {
		h.AuthenticationRequiredResponse(w, r)
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	totalIncome, totalExpenses, netProfit, err := h.Services.Finance.GetFinancialStats(r.Context(), user.ID, startDate, endDate)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	profitMarginPct := 0.0
	if totalIncome > 0 {
		profitMarginPct = (netProfit / totalIncome) * 100.0
	}

	now := time.Now()
	type MonthData struct {
		MonthKey string  `json:"month_key"`
		Month    string  `json:"month"`
		Income   float64 `json:"income"`
		Expenses float64 `json:"expenses"`
		Net      float64 `json:"net"`
	}

	monthsMap := make(map[string]*MonthData)
	var monthlyTrend []*MonthData

	for i := 5; i >= 0; i-- {
		mDate := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		key := mDate.Format("2006-01")
		item := &MonthData{
			MonthKey: key,
			Month:    mDate.Format("Jan 2006"),
			Income:   0,
			Expenses: 0,
			Net:      0,
		}
		monthlyTrend = append(monthlyTrend, item)
		monthsMap[key] = item
	}

	sixMonthsAgo := time.Date(now.Year(), now.Month()-5, 1, 0, 0, 0, 0, time.UTC)
	var trendTxns []*models.FinancialTransaction
	err = h.Models.DB.WithContext(r.Context()).
		Where("user_id = ? AND transaction_date >= ?", user.ID, sixMonthsAgo).
		Find(&trendTxns).Error
	if err == nil {
		for _, t := range trendTxns {
			k := t.TransactionDate.Format("2006-01")
			if bucket, exists := monthsMap[k]; exists {
				if t.Type == "income" {
					bucket.Income += t.Amount
				} else if t.Type == "expense" {
					bucket.Expenses += t.Amount
				}
			}
		}
		for _, bucket := range monthlyTrend {
			bucket.Net = bucket.Income - bucket.Expenses
		}
	}

	type CatDist struct {
		CategoryID string  `json:"category_id"`
		Name       string  `json:"name"`
		Type       string  `json:"type"`
		Color      string  `json:"color"`
		Icon       string  `json:"icon"`
		Amount     float64 `json:"amount"`
		Percentage float64 `json:"percentage"`
	}

	expTxns, _ := h.Services.Finance.GetAllTransactions(r.Context(), user.ID, "expense", "", startDate, endDate, "")

	catTotals := make(map[string]*CatDist)
	var totalExpAmount float64
	for _, t := range expTxns {
		cid := "uncategorized"
		cname := "Uncategorized"
		ccolor := "slate"
		cicon := "more_horiz"
		if t.Category != nil {
			cid = t.Category.ID.String()
			cname = t.Category.Name
			ccolor = t.Category.Color
			cicon = t.Category.Icon
		}
		if _, exists := catTotals[cid]; !exists {
			catTotals[cid] = &CatDist{
				CategoryID: cid,
				Name:       cname,
				Type:       t.Type,
				Color:      ccolor,
				Icon:       cicon,
				Amount:     0,
				Percentage: 0,
			}
		}
		catTotals[cid].Amount += t.Amount
		totalExpAmount += t.Amount
	}

	var categoryDistribution []*CatDist
	for _, cd := range catTotals {
		if totalExpAmount > 0 {
			cd.Percentage = (cd.Amount / totalExpAmount) * 100.0
		}
		categoryDistribution = append(categoryDistribution, cd)
	}

	var unpaidInvoices []*models.Invoice
	_ = h.Models.DB.WithContext(r.Context()).
		Where("user_id = ? AND is_paid = false", user.ID).
		Find(&unpaidInvoices).Error

	var unpaidTotal float64
	var unpaidCount int64
	var overdueTotal float64
	var overdueCount int64
	for _, inv := range unpaidInvoices {
		unpaidTotal += inv.Total
		unpaidCount++
		if inv.DueDate.Before(now) {
			overdueTotal += inv.Total
			overdueCount++
		}
	}

	err = h.WriteJSON(w, http.StatusOK, Envelope{
		"status": "success",
		"analytics": Envelope{
			"summary": Envelope{
				"total_income":      totalIncome,
				"total_expenses":    totalExpenses,
				"net_profit":        netProfit,
				"profit_margin_pct": profitMarginPct,
			},
			"monthly_trend":         monthlyTrend,
			"category_distribution": categoryDistribution,
			"receivables": Envelope{
				"unpaid_total":  unpaidTotal,
				"unpaid_count":  unpaidCount,
				"overdue_total": overdueTotal,
				"overdue_count": overdueCount,
			},
		},
	}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

