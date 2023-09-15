package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"user-balance-service/config"
	"user-balance-service/pkg/domain"
	"user-balance-service/pkg/postgres"
	"user-balance-service/pkg/repo"
	"user-balance-service/pkg/service"
)

type SetupRequest struct {
	CFG      *config.Config
	Postgres *postgres.Postgres
	Wallet   domain.WalletService
}

func errorResponse(c *gin.Context, code int, e error) {
	type response struct {
		Error string `json:"error"`
	}

	c.AbortWithStatusJSON(code, response{
		Error: e.Error(),
	})
}

func Run(req *SetupRequest) http.Handler {
	r := gin.Default()

	r.GET("/balance/:id", func(c *gin.Context) {
		balance, err := req.Wallet.Balance(c.Request.Context(), c.Param("id"))
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"balance": balance,
		})
	})

	r.POST("/deposit/:id", func(c *gin.Context) {

		type deposit struct {
			Amount int `json:"amount"`
		}

		var data deposit
		userID := c.Param("id")
		ctx := c.Request.Context()

		err := c.ShouldBindJSON(&data)
		if err != nil {
			errorResponse(c, http.StatusBadRequest, err)
			return
		}

		if data.Amount < 1 {
			errorResponse(c, http.StatusBadRequest, errors.New("amount must be greater or equal 1"))
			return
		}

		// start tx
		tx, err := req.Postgres.Begin(ctx)
		if err != nil {
			_ = c.Error(err)
			return
		}

		// create service & repo wrapped on database TX
		walletServiceTx := service.NewWallet(repo.NewTransactionPGRepo(tx))
		lockerTX := service.NewLocker(tx)

		err = lockerTX.Lock(ctx, userID)
		if err != nil {
			_ = tx.Rollback(ctx)
			errorResponse(c, http.StatusInternalServerError, err)
			return
		}

		err = walletServiceTx.Deposit(ctx, userID, data.Amount)
		if err != nil {
			_ = tx.Rollback(ctx)
			errorResponse(c, http.StatusInternalServerError, err)
			return
		}

		err = tx.Commit(ctx)
		if err != nil {
			_ = tx.Rollback(ctx)
			errorResponse(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	})

	return r
}
