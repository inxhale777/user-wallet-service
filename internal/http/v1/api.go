package v1

import (
	"github.com/gin-gonic/gin"
	"user-wallet-service/internal/domain"
	"user-wallet-service/internal/http/v1/handlers"
	"user-wallet-service/internal/postgres"
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
