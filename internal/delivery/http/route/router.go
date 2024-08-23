package route

import (
	"boiler-plate-clean/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

type Router struct {
	App            *gin.Engine
	ExampleHandler *http.UserHTTPHandler
}

func (h *Router) Setup() {
	api := h.App.Group("/api/v2")
	{

		//Example Routes
		campaignApi := api.Group("/user")
		//campaignApi.Use(h.RequestMiddleware.RequestHeader)
		{
			campaignApi.POST("", h.ExampleHandler.Create)
			campaignApi.GET("", h.ExampleHandler.Find)
			campaignApi.GET("/:id", h.ExampleHandler.FindOne)
		}
	}
}
