package middlewares

import (
	"envoyer/config/consts"
	"envoyer/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SuperAdminPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		permission := ctx.GetString("user_access")
		if permission != consts.SuperAdminRole {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": utils.Trans("userDoesNotHavePermissionToAccess", nil)})
			return
		}
		ctx.Next()
	}
}

func AppAdminPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		permission := ctx.GetString("user_access")
		switch permission {
		case
			consts.SuperAdminRole,
			consts.DevRole,
			consts.AdminRole:
			ctx.Next()
			return
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": utils.Trans("userDoesNotHavePermissionToAccess", nil)})
		return
	}
}

func DevPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		permission := ctx.GetString("user_access")
		switch permission {
		case
			consts.SuperAdminRole,
			consts.DevRole:
			ctx.Next()
			return
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": utils.Trans("userDoesNotHavePermissionToAccess", nil)})
		return
	}
}
