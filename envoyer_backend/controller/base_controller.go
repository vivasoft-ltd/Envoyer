package controller

import (
	"envoyer/config/consts"
	"envoyer/logger"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
	"net/http"
)

type BaseController struct {
	logger    logger.Logger
	container di.Container
}

func NewBaseController(container di.Container, logger logger.Logger) *BaseController {
	return &BaseController{logger: logger, container: container}
}

func (c BaseController) ReplySuccess(context *gin.Context, data interface{}) {
	c.response(context, gin.H{"data": data, "status": consts.StatusOk}, http.StatusOK)
}

func (c BaseController) SuccessResponse(context *gin.Context, data interface{}, data2 gin.H) {
	d := gin.H{"data": data, "status": consts.StatusOk}
	if len(data2) > 0 {
		for k, v := range data2 {
			d[k] = v
		}
	}
	c.response(context, d, http.StatusOK)
}

func (c BaseController) ReplyError(context *gin.Context, message string, code int) {
	c.response(context, gin.H{"message": message, "status": consts.StatusError}, code)
}

func (c BaseController) response(context *gin.Context, obj interface{}, code int) {
	switch context.GetHeader("Accept") {
	case "application/xml":
		context.XML(code, obj)
	default:
		context.JSON(code, obj)
	}
}
