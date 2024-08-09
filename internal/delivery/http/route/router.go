package route

import (
	"boiler-plate-clean/config"
	"boiler-plate-clean/internal/delivery/grpc"
	"boiler-plate-clean/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

type Router struct {
	App                *gin.Engine
	UserHandler        *http.UserHTTPHandler
	WalletHandler      *http.WalletHTTPHandler
	CategoryHandler    *http.CategoryTransactionHTTPHandler
	TransactionHandler *http.TransactionHTTPHandler
	GRPCHandler        *grpc.BaseGRPCHandler
	Config             *config.Config
}

func (h *Router) Setup() {
	api := h.App.Group("api/" + h.Config.AppEnvConfig.AppVersion)
	{

		//Example Routes
		userApi := api.Group("/user")
		{
			userApi.POST("", h.UserHandler.Create)
			userApi.GET("/:id", h.UserHandler.FindOne)
			userApi.GET("/cashflow/:id", h.UserHandler.Cashflow)
			userApi.PUT("/:id", h.UserHandler.Update)
			userApi.DELETE("/:id", h.UserHandler.Delete)
		}

		walletApi := api.Group("/wallet")
		{
			walletApi.POST("", h.WalletHandler.Create)
			walletApi.GET("/:id", h.WalletHandler.FindOne)
			walletApi.PUT("/:id", h.WalletHandler.Update)
			walletApi.GET("/latest/:id", h.WalletHandler.Last10)
			walletApi.GET("/category/:id", h.WalletHandler.RecapCategory)
			walletApi.DELETE("/:id", h.WalletHandler.Delete)
		}

		categoryApi := api.Group("/category")
		{
			categoryApi.POST("", h.CategoryHandler.Create)
			categoryApi.GET("/:id", h.CategoryHandler.FindOne)
			categoryApi.PUT("/:id", h.CategoryHandler.Update)
			categoryApi.GET("", h.CategoryHandler.Find)
			categoryApi.DELETE("/:id", h.CategoryHandler.Delete)
		}

		transactionApi := api.Group("/transaction")
		{
			transactionApi.POST("", h.TransactionHandler.Create)
			transactionApi.GET("/:id", h.TransactionHandler.FindOne)
			transactionApi.PUT("/:id", h.TransactionHandler.Update)
			transactionApi.PUT("/credit", h.TransactionHandler.Credit)
			transactionApi.PUT("/transfer", h.TransactionHandler.Transfer)
			transactionApi.DELETE("/id", h.TransactionHandler.Delete)
		}
	}
}
