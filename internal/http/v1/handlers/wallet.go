package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"user-balance-service/internal/domain"
	"user-balance-service/internal/postgres"
	"user-balance-service/internal/repo/pg_transactions"
	"user-balance-service/internal/service/pg_locker"
	"user-balance-service/internal/service/wallet"
)

type hh struct {
	wallet domain.WalletService
	db     postgres.DB
}

func SetupWallet(rg *gin.RouterGroup, w domain.WalletService, db postgres.DB) {
	h := &hh{
		wallet: w,
		db:     db,
	}

	rg.GET("/balance/:id", h.Balance)
	rg.POST("/deposit/:id", h.Deposit)
}

func errorResponse(c *gin.Context, code int, e error) {
	type response struct {
		Error string `json:"error"`
	}

	c.AbortWithStatusJSON(code, response{
		Error: e.Error(),
	})
}

func (h *hh) Balance(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	balance, err := h.wallet.Balance(c.Request.Context(), userID)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}

func (h *hh) Deposit(c *gin.Context) {

	type deposit struct {
		Amount int `json:"amount"`
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	var data deposit
	err = c.ShouldBindJSON(&data)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if data.Amount < 1 {
		errorResponse(c, http.StatusBadRequest, errors.New("amount must be greater or equal 1"))
		return
	}

	ctx := c.Request.Context()

	// start tx
	tx, err := h.db.Begin(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// create service, repo and locker wrapped on database TX
	walletServiceTx := wallet.New(pg_transactions.New(tx), pg_locker.New(tx))

	err = walletServiceTx.Deposit(ctx, userID, data.Amount)
	if err != nil {
		_ = tx.Rollback(ctx)
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
