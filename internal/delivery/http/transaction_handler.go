package http

import (
	_ "boiler-plate-clean/internal/delivery/http/response"
	"boiler-plate-clean/internal/entity"
	service "boiler-plate-clean/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TransactionHTTPHandler struct {
	Handler
	TransactionService service.TransactionService
	FileService        service.FileService
}

func NewTransactionHTTPHandler(example service.TransactionService, file service.FileService) *TransactionHTTPHandler {
	return &TransactionHTTPHandler{
		TransactionService: example,
		FileService:        file,
	}
}

func (h TransactionHTTPHandler) Create(ctx *gin.Context) {
	request := entity.Transaction{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	if errException := h.TransactionService.Create(ctx, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

func (h TransactionHTTPHandler) FindOne(ctx *gin.Context) {
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	result, errException := h.TransactionService.Detail(ctx, idParam)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

func (h TransactionHTTPHandler) Update(ctx *gin.Context) {
	// Get Info
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	request := entity.Transaction{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	if errException := h.TransactionService.Update(ctx, idParam, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

func (h TransactionHTTPHandler) Credit(ctx *gin.Context) {
	walletId, err := h.QueryInt(ctx, "wallet_id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	categoryId, err := h.QueryInt(ctx, "category_id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	amount, err := h.QueryFloat64(ctx, "amount")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	if errException := h.TransactionService.Credit(ctx, walletId, categoryId, amount); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.SuccessJSON(ctx)
}

func (h TransactionHTTPHandler) Transfer(ctx *gin.Context) {
	senderId, err := h.QueryInt(ctx, "sender_id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	receiverId, err := h.QueryInt(ctx, "receiver_id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	amount, err := h.QueryFloat64(ctx, "amount")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	if errException := h.TransactionService.Transfer(ctx, senderId, receiverId, amount); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.SuccessJSON(ctx)
}

func (h TransactionHTTPHandler) Delete(ctx *gin.Context) {
	idParam, err := h.ParamInt(ctx, "id")
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
	}
	errException := h.TransactionService.Delete(ctx, idParam)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.SuccessJSON(ctx)
}

func (h TransactionHTTPHandler) RetrieveDocument(c *gin.Context) {

	document, exc := h.FileService.Download(c)
	if exc != nil {
		h.ExceptionJSON(c, exc)
		return
	}

	mimeType := http.DetectContentType(document)
	c.Header("Content-Description", "File Transfer")
	c.Data(http.StatusOK, mimeType, document)
}
