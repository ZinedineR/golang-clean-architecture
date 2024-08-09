package http

import (
	_ "boiler-plate-clean/internal/delivery/http/response"
	"boiler-plate-clean/internal/entity"
	service "boiler-plate-clean/internal/services"
	"github.com/gin-gonic/gin"
)

type CategoryTransactionHTTPHandler struct {
	Handler
	CategoryTransactionService service.CategoryTransactionService
}

func NewCategoryTransactionHTTPHandler(example service.CategoryTransactionService) *CategoryTransactionHTTPHandler {
	return &CategoryTransactionHTTPHandler{
		CategoryTransactionService: example,
	}
}

func (h CategoryTransactionHTTPHandler) Create(ctx *gin.Context) {
	request := entity.CategoryTransaction{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	if errException := h.CategoryTransactionService.Create(ctx, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

func (h CategoryTransactionHTTPHandler) FindOne(ctx *gin.Context) {
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	result, errException := h.CategoryTransactionService.Detail(ctx, idParam)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

func (h CategoryTransactionHTTPHandler) Find(ctx *gin.Context) {
	_, order, filter, err := h.ParsePaginationParams(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	result, errException := h.CategoryTransactionService.Find(ctx, order, filter)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

func (h CategoryTransactionHTTPHandler) Delete(ctx *gin.Context) {
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	errException := h.CategoryTransactionService.Delete(ctx, idParam)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.SuccessJSON(ctx)
}

func (h CategoryTransactionHTTPHandler) Update(ctx *gin.Context) {
	// Get Info
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	request := entity.CategoryTransaction{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	if errException := h.CategoryTransactionService.Update(ctx, idParam, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}
