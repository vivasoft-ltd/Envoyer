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

type AppController struct {
	*BaseController
	appService service.AppServiceInterface
	*SubscriberController
}

func NewAppController(base *BaseController, subscriptionController *SubscriberController) *AppController {
	return &AppController{
		BaseController:       base,
		SubscriberController: subscriptionController,
		appService:           base.container.Get(service_name.AppService).(service.AppServiceInterface),
	}
}

func (c AppController) CreateApp(context *gin.Context) {
	var payload serializers.CreateAppReq
	bindErr := context.ShouldBindJSON(&payload)
	if bindErr != nil {
		c.ReplyError(context, bindErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.appService.Create(payload)
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	if resp.Active {
		c.StartSubscriber(resp.Name, "v2")
	}
	c.ReplySuccess(context, resp)
	return
}

func (c AppController) GetAllApp(context *gin.Context) {
	resp, err := c.appService.GetAll()
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c AppController) GetApp(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	if context.GetString("user_access") != consts.SuperAdminRole {
		appid := context.GetUint("app_id")
		if appid != uint(id) {
			c.ReplyError(context, utils.Trans("userDoesNotHavePermissionToAccess", nil), http.StatusUnauthorized)
			return
		}
	}

	resp, err := c.appService.Get(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c AppController) UpdateApp(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	var payload serializers.UpdateAppReq
	bindErr := context.ShouldBindJSON(&payload)
	if bindErr != nil {
		c.ReplyError(context, bindErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.appService.Update(payload, uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}

	if resp.Active {
		c.StartSubscriber(resp.Name, "v2")
	} else {
		c.StopSubscriber(resp.Name)
	}

	c.ReplySuccess(context, resp)
	return
}

func (c AppController) DeleteApp(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.appService.Get(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}

	err = c.appService.Delete(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}

	c.StopSubscriber(resp.Name)
	c.ReplySuccess(context, nil)
	return
}

func (c AppController) GetAppByAppKey(context *gin.Context) {
	appKey := context.Params.ByName("id")
	if len(appKey) <= 0 {
		c.ReplyError(context, "app id not found", http.StatusBadRequest)
		return
	}

	resp, err := c.appService.GetByAppKey(appKey)
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}
