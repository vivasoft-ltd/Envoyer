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

type UserController struct {
	*BaseController
	userService service.UserServiceInterface
}

func NewUserController(base *BaseController) *UserController {
	return &UserController{
		BaseController: base,
		userService:    base.container.Get(service_name.UserService).(service.UserServiceInterface),
	}
}

func (c UserController) CreateUser(context *gin.Context) {
	var payload serializers.CreateUserReq
	bindErr := context.ShouldBindJSON(&payload)
	if bindErr != nil {
		c.ReplyError(context, bindErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.userService.Create(payload)
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c UserController) GetByAppId(context *gin.Context) {
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

	resp, err := c.userService.GetByAppId(uint(appId))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c UserController) GetUser(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.userService.Get(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c UserController) UpdateUser(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	var payload serializers.UpdateUserReq
	bindErr := context.ShouldBindJSON(&payload)
	if bindErr != nil {
		c.ReplyError(context, bindErr.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.userService.Update(payload, uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, resp)
	return
}

func (c UserController) DeleteUser(context *gin.Context) {
	id, convErr := strconv.Atoi(context.Params.ByName("id"))
	if convErr != nil {
		c.ReplyError(context, convErr.Error(), http.StatusBadRequest)
		return
	}

	err := c.userService.Delete(uint(id))
	if err != nil {
		c.ReplyError(context, err.Error(), err.Status)
		return
	}
	c.ReplySuccess(context, nil)
	return
}
