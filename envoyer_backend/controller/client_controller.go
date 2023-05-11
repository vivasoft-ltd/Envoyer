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

type ClientController struct {
	*BaseController
	clientService service.ClientServiceInterface
}

func NewClientController(base *BaseController) *ClientController {
	return &ClientController{
		BaseController: base,
		clientService:  base.container.Get(service_name.ClientService).(service.ClientServiceInterface),
	}
}

func (c ClientController) CreateClient(context *gin.Context) {
	var payload serializers.CreateClientReq
	bindErr := context.ShouldBindJSON(&payload)
	if bindErr != nil {
		c.ReplyError(context, bindErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.clientService.Create(payload)
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c ClientController) GetByAppId(context *gin.Context) {
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

	resp, err := c.clientService.GetByAppId(uint(appId))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c ClientController) GetClient(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.clientService.Get(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c ClientController) UpdateClient(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	var payload serializers.UpdateClientReq
	bindErr := context.ShouldBindJSON(&payload)
	if bindErr != nil {
		c.ReplyError(context, bindErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.clientService.Update(payload, uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c ClientController) DeleteClient(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	err := c.clientService.Delete(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, nil)
	return
}

func (c ClientController) GetByClientId(context *gin.Context) {
	clientKey := context.Params.ByName("id")
	if len(clientKey) <= 0 {
		c.ReplyError(context, "client id not found", http.StatusBadRequest)
		return
	}

	resp, err := c.clientService.GetByClientKey(clientKey)
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}
