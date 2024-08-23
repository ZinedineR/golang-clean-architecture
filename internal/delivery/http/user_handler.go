package http

import (
	_ "boiler-plate-clean/internal/delivery/http/response"
	"boiler-plate-clean/internal/gateway/messaging"
	"boiler-plate-clean/internal/model"
	service "boiler-plate-clean/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHTTPHandler struct {
	Handler
	UserService  service.UserService
	UserProducer messaging.UserProducer
}

func NewUserHTTPHandler(
	example service.UserService, userWrite messaging.UserProducer,
) *UserHTTPHandler {
	return &UserHTTPHandler{
		UserService:  example,
		UserProducer: userWrite,
	}
}

func (h UserHTTPHandler) Create(ctx *gin.Context) {
	request := model.UserMessage{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	if err := h.UserProducer.Send(ctx, &request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	h.DataJSON(ctx, request)
}

func (h UserHTTPHandler) Find(ctx *gin.Context) {
	var req model.ListReq
	var err error
	req.Page, req.Order, req.Filter, err = h.ParsePaginationParams(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	result, errException := h.UserService.Find(ctx, &req)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

func (h UserHTTPHandler) FindOne(ctx *gin.Context) {
	idParam := ctx.Param("id")
	result, errException := h.UserService.FindById(ctx, idParam)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}
