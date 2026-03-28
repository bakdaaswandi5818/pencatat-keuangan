package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bakdaaswandi5818/pencatat-keuangan/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

// TransactionHandler wires HTTP routes to the service layer.
type TransactionHandler struct {
	svc service.TransactionService
}

// NewTransactionHandler creates a new TransactionHandler.
func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

// Register attaches routes to the provided Echo group.
func (h *TransactionHandler) Register(e *echo.Echo) {
	e.GET("/health", h.HealthCheck)
	e.GET("/transactions", h.ListTransactions)
	e.POST("/transactions", h.CreateTransaction)
	e.GET("/transactions/:id", h.GetTransaction)
	e.DELETE("/transactions/:id", h.DeleteTransaction)
	e.GET("/summary", h.GetSummary)
}

// HealthCheck godoc
// @Summary Health check
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *TransactionHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// ListTransactions godoc
// @Summary List transactions with pagination and filtering
// @Produce json
// @Param type     query string false "Filter by type (income|expense)"
// @Param category query string false "Filter by category"
// @Param date_from query string false "Filter start date (YYYY-MM-DD)"
// @Param date_to   query string false "Filter end date (YYYY-MM-DD)"
// @Param limit  query int false "Page size (default 10, max 100)"
// @Param offset query int false "Page offset (default 0)"
// @Success 200
// @Router /transactions [get]
func (h *TransactionHandler) ListTransactions(c echo.Context) error {
	input := service.ListTransactionsInput{}

	input.Type = c.QueryParam("type")
	input.Category = c.QueryParam("category")

	if ds := c.QueryParam("date_from"); ds != "" {
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid date_from format, use YYYY-MM-DD")
		}
		input.DateFrom = &t
	}
	if ds := c.QueryParam("date_to"); ds != "" {
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid date_to format, use YYYY-MM-DD")
		}
		input.DateTo = &t
	}

	if ls := c.QueryParam("limit"); ls != "" {
		l, err := strconv.Atoi(ls)
		if err != nil || l < 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid limit")
		}
		input.Limit = l
	}
	if os := c.QueryParam("offset"); os != "" {
		o, err := strconv.Atoi(os)
		if err != nil || o < 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid offset")
		}
		input.Offset = o
	}

	out, err := h.svc.List(input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, out)
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Accept  json
// @Produce json
// @Param   body body service.CreateTransactionInput true "Transaction payload"
// @Success 201
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	var input service.CreateTransactionInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(input); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	tx, err := h.svc.Create(input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, tx)
}

// GetTransaction godoc
// @Summary Get a single transaction by ID
// @Produce json
// @Param id path string true "Transaction UUID"
// @Success 200
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid uuid")
	}
	tx, err := h.svc.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "transaction not found")
	}
	return c.JSON(http.StatusOK, tx)
}

// DeleteTransaction godoc
// @Summary Soft-delete a transaction by ID
// @Produce json
// @Param id path string true "Transaction UUID"
// @Success 204
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid uuid")
	}
	if err := h.svc.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// GetSummary godoc
// @Summary Get financial summary (total income, expense, balance)
// @Produce json
// @Success 200
// @Router /summary [get]
func (h *TransactionHandler) GetSummary(c echo.Context) error {
	s, err := h.svc.GetSummary()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, s)
}
