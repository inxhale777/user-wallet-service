package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user-wallet-service/internal/domain"
	"user-wallet-service/internal/postgres"
	"user-wallet-service/internal/repo/pg_transactions"
	"user-wallet-service/internal/service/pg_locker"
	"user-wallet-service/internal/service/wallet"
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
	rg.POST("/hold/:id", h.Hold)
	rg.POST("/charge/:txID", h.Charge)
}

func errorResponse(ctx *gin.Context, code int, e error) {
	type response struct {
		Error string `json:"error"`
	}

	ctx.AbortWithStatusJSON(code, response{
		Error: e.Error(),
	})
}

func (h *hh) Balance(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	balance, err := h.wallet.Balance(ctx, userID)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}

func (h *hh) Deposit(ctx *gin.Context) {

	type deposit struct {
		Amount int `json:"amount"`
	}

	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	var data deposit
	err = ctx.ShouldBindJSON(&data)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	if data.Amount < 1 {
		errorResponse(ctx, http.StatusBadRequest, domain.ErrInvalidAmount)
		return
	}

	// start tx
	tx, err := h.db.Begin(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// create service, repo and locker wrapped on database TX
	walletServiceTx := wallet.New(pg_transactions.New(tx), pg_locker.New(tx))

	err = walletServiceTx.Deposit(ctx, userID, data.Amount)
	if err != nil {
		_ = tx.Rollback(ctx)
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (h *hh) Hold(ctx *gin.Context) {

	type hold struct {
		Amount int `json:"amount"`
	}

	var data hold
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	if data.Amount < 1 {
		errorResponse(ctx, http.StatusBadRequest, domain.ErrInvalidAmount)
		return
	}

	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	// start tx
	tx, err := h.db.Begin(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	// create service, repo and locker wrapped on database TX
	walletServiceTx := wallet.New(pg_transactions.New(tx), pg_locker.New(tx))
	txID, err := walletServiceTx.Hold(ctx, userID, data.Amount)
	if err != nil {
		_ = tx.Rollback(ctx)
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	type response struct {
		TransactionID int `json:"transaction_id"`
	}

	ctx.JSON(http.StatusOK, response{
		TransactionID: txID,
	})
}

func (h *hh) Charge(ctx *gin.Context) {
	txID, err := strconv.Atoi(ctx.Param("txID"))
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.wallet.Charge(ctx, txID)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
