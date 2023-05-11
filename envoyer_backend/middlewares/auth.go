package middlewares

import (
	"envoyer/config"
	"envoyer/config/consts"
	"envoyer/config/service_name"
	"envoyer/dic"
	"envoyer/model/service"
	"envoyer/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		sub, err := utils.ValidateToken(accessToken, config.Config.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		ctx.Set("user_id", int64(sub.(float64)))

		if int64(sub.(float64)) == 0 {
			ctx.Set("user_access", consts.SuperAdminRole)
			ctx.Next()
		} else {
			authService := dic.Container.Get(service_name.AuthService).(service.AuthServiceInterface)
			user, getErr := authService.FindUserById(uint(sub.(float64)))
			if getErr != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": getErr.Error()})
				return
			}

			ctx.Set("app_id", user.AppId)
			ctx.Set("user_access", user.Role)
			ctx.Next()
		}
	}
}
