package http

import (
	_ "boiler-plate-clean/internal/delivery/http/response"
	"boiler-plate-clean/internal/entity"
	service "boiler-plate-clean/internal/services"
	"github.com/gin-gonic/gin"
)

type WalletHTTPHandler struct {
	Handler
	WalletService service.WalletService
}

func NewWalletHTTPHandler(example service.WalletService) *WalletHTTPHandler {
	return &WalletHTTPHandler{
		WalletService: example,
	}
}

func (h WalletHTTPHandler) Create(ctx *gin.Context) {
	request := entity.Wallet{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	if errException := h.WalletService.Create(ctx, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

func (h WalletHTTPHandler) FindOne(ctx *gin.Context) {
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}

	fromDate, toDate, err := h.ParseDateParam(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	result, errException := h.WalletService.DetailWalletTransaction(ctx, idParam, fromDate, toDate)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

func (h WalletHTTPHandler) Last10(ctx *gin.Context) {
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}

	result, errException := h.WalletService.Last10(ctx, idParam)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

func (h WalletHTTPHandler) RecapCategory(ctx *gin.Context) {
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}

	category, err := h.QueryInt(ctx, "category_id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}

	fromDate, toDate, err := h.ParseDateParam(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	result, errException := h.WalletService.RecapCategory(ctx, idParam, category, fromDate, toDate)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

func (h WalletHTTPHandler) Update(ctx *gin.Context) {
	// Get Info
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	request := entity.Wallet{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	if errException := h.WalletService.Update(ctx, idParam, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

func (h WalletHTTPHandler) Delete(ctx *gin.Context) {
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	errException := h.WalletService.Delete(ctx, idParam)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.SuccessJSON(ctx)
}
