package v1

import (
	"github.com/gin-gonic/gin"
	"user-balance-service/internal/domain"
	"user-balance-service/internal/http/v1/handlers"
	"user-balance-service/internal/postgres"
)

type SetupRequest struct {
	DB     postgres.DB
	Wallet domain.WalletService
}

func Run(req *SetupRequest) *gin.Engine {
	r := gin.Default()

	rg := r.Group("/wallet")
	handlers.SetupWallet(rg, req.Wallet, req.DB)

	return r
}
