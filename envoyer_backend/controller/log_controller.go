package controller

import (
	"envoyer/config/service_name"
	"envoyer/model/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LogController struct {
	*BaseController
	logService service.LogServiceInterface
}

func NewLogController(base *BaseController) *LogController {
	return &LogController{
		BaseController: base,
		logService:     base.container.Get(service_name.LogService).(service.LogServiceInterface),
	}
}

func (c LogController) DeleteLog(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	err := c.logService.Delete(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, nil)
	return
}

func (c LogController) GetByAppId(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.logService.GetByAppId(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}
