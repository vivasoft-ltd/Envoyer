package controller

import (
	"envoyer/config"
	"envoyer/config/consts"
	"envoyer/config/service_name"
	"envoyer/model/entity"
	"envoyer/model/serializers"
	"envoyer/model/service"
	"envoyer/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"net/http"
)

type AuthController struct {
	*BaseController
	authService service.AuthServiceInterface
}

func NewAuthController(base *BaseController) *AuthController {
	return &AuthController{
		BaseController: base,
		authService:    base.container.Get(service_name.AuthService).(service.AuthServiceInterface),
	}
}

func (a AuthController) LogIn(context *gin.Context) {

	var tokenRequest serializers.LoginReq
	if err := context.ShouldBindBodyWith(&tokenRequest, binding.JSON); err != nil {
		a.ReplyError(context, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := a.authService.LoginUser(tokenRequest)
	if err != nil {
		a.ReplyError(context, err.Error(), err.Status)
		return
	}

	a.ReplySuccess(context, response)
	return
}

func (a AuthController) RefreshAccessToken(context *gin.Context) {
	var refreshTokenRequest serializers.RefreshTokenRequest
	if err := context.ShouldBindBodyWith(&refreshTokenRequest, binding.JSON); err != nil {
		a.ReplyError(context, err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := utils.ValidateToken(refreshTokenRequest.RefreshToken, config.Config.RefreshTokenPublicKey)
	if err != nil {
		a.ReplyError(context, err.Error(), http.StatusForbidden)
		return
	}

	id := int(sub.(float64))
	if id == 0 {
		token, tokenErr := a.authService.GetAccessToken(&entity.User{
			Model: gorm.Model{
				ID: 0,
			},
			Role: consts.SuperAdminRole,
		})
		if tokenErr != nil {
			a.ReplyError(context, tokenErr.Error(), tokenErr.Status)
			return
		}
		a.ReplySuccess(context, token)
		return
	}
	user, getErr := a.authService.FindUserById(uint(id))
	if getErr != nil || user.ID == 0 {
		a.ReplyError(context, getErr.Error(), http.StatusForbidden)
		return
	}

	token, tokenErr := a.authService.GetAccessToken(user)
	if tokenErr != nil {
		a.ReplyError(context, tokenErr.Error(), tokenErr.Status)
		return
	}
	a.ReplySuccess(context, token)
	return
}
