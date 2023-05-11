package controller

import (
	"envoyer/config/consts"
	"envoyer/config/service_name"
	"envoyer/model/serializers"
	"envoyer/model/service"
	"envoyer/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type EventController struct {
	*BaseController
	eventService service.EventServiceInterface
}

func NewEventController(base *BaseController) *EventController {
	return &EventController{
		BaseController: base,
		eventService:   base.container.Get(service_name.EventService).(service.EventServiceInterface),
	}
}

func (c EventController) CreateEvent(context *gin.Context) {
	var payload serializers.CreateEventReq
	bindErr := context.ShouldBindJSON(&payload)
	if bindErr != nil {
		c.ReplyError(context, bindErr.Error(), http.StatusBadRequest)
		return
	}
	for _, v := range payload.Variables {
		if len(v) < 2 {
			c.ReplyError(context, "Variables must be minimum 2 characters", http.StatusBadRequest)
			return
		}
	}

	resp, err := c.eventService.Create(payload)
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c EventController) GetByAppId(context *gin.Context) {
	appId, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	if context.GetString("user_access") != consts.SuperAdminRole {
		appid := context.GetUint("app_id")
		if appid != uint(appId) {
			c.ReplyError(context, utils.Trans("userDoesNotHavePermissionToAccess", nil), http.StatusUnauthorized)
			return
		}
	}

	resp, err := c.eventService.GetByAppId(uint(appId))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c EventController) GetEvent(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.eventService.Get(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c EventController) UpdateEvent(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	var payload serializers.UpdateEventReq
	bindErr := context.ShouldBindJSON(&payload)
	if bindErr != nil {
		c.ReplyError(context, bindErr.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range payload.Variables {
		if len(v) < 2 {
			c.ReplyError(context, "Variables must be minimum 2 characters", http.StatusBadRequest)
			return
		}
	}

	resp, err := c.eventService.Update(payload, uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c EventController) DeleteEvent(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	err := c.eventService.Delete(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, nil)
	return
}
