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

type ProviderController struct {
	*BaseController
	providerService service.ProviderServiceInterface
}

func NewProviderController(base *BaseController) *ProviderController {
	return &ProviderController{
		BaseController:  base,
		providerService: base.container.Get(service_name.ProviderService).(service.ProviderServiceInterface),
	}
}

func (c ProviderController) CreateProvider(context *gin.Context) {
	var payload serializers.CreateProviderReq
	bindErr := context.ShouldBindJSON(&payload)
	if bindErr != nil {
		c.ReplyError(context, bindErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.providerService.Create(payload)
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c ProviderController) GetProvider(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.providerService.Get(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c ProviderController) UpdateProvider(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	var payload serializers.UpdateProviderReq
	bindErr := context.ShouldBindJSON(&payload)
	if bindErr != nil {
		c.ReplyError(context, bindErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.providerService.Update(payload, uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c ProviderController) DeleteProvider(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	err := c.providerService.Delete(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, nil)
	return
}

func (c ProviderController) GetProviderByAppId(context *gin.Context) {
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

	resp, err := c.providerService.GetByAppId(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c ProviderController) GetProvidersByAppIdAndType(context *gin.Context) {
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

	Type := context.Params.ByName("type")
	if len(Type) <= 0 {
		c.ReplyError(context, "type is not found", http.StatusBadRequest)
		return
	}

	resp, err := c.providerService.GetByAppIdAndType(uint(id), Type)
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c ProviderController) GetProviderByAppIdAndTypeWithTopPriority(context *gin.Context) {
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

	Type := context.Params.ByName("type")
	if len(Type) <= 0 {
		c.ReplyError(context, "type is not found", http.StatusBadRequest)
		return
	}

	resp, err := c.providerService.GetByAppIdAndTypeWithTopPriority(uint(id), Type)
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c ProviderController) UpdateProviderPriority(context *gin.Context) {
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

	Type := context.Params.ByName("type")
	if len(Type) <= 0 {
		c.ReplyError(context, "type is not found", http.StatusBadRequest)
		return
	}

	var payload serializers.UpdatePriorityReq
	bindErr := context.ShouldBindJSON(&payload)
	if bindErr != nil {
		c.ReplyError(context, bindErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.providerService.UpdatePriority(uint(id), Type, payload)
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}

	c.ReplySuccess(context, resp)
	return
}
