package http

import (
	_ "boiler-plate-clean/internal/delivery/http/response"
	"boiler-plate-clean/internal/entity"
	service "boiler-plate-clean/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHTTPHandler struct {
	Handler
	UserService service.UserService
}

func NewUserHTTPHandler(example service.UserService) *UserHTTPHandler {
	return &UserHTTPHandler{
		UserService: example,
	}
}

func (h UserHTTPHandler) Create(ctx *gin.Context) {
	request := entity.Users{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	if errException := h.UserService.Create(ctx, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

func (h UserHTTPHandler) FindOne(ctx *gin.Context) {
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	result, errException := h.UserService.Detail(ctx, idParam)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

func (h UserHTTPHandler) Cashflow(ctx *gin.Context) {
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	fromDate, toDate, err := h.ParseDateParam(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	result, errException := h.UserService.Cashflow(ctx, idParam, fromDate, toDate)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

func (h UserHTTPHandler) Update(ctx *gin.Context) {
	// Get Info
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	request := entity.Users{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	if errException := h.UserService.Update(ctx, idParam, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

func (h UserHTTPHandler) Delete(ctx *gin.Context) {
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	errException := h.UserService.Delete(ctx, idParam)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.SuccessJSON(ctx)
}
