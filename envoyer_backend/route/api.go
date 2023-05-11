package route

import (
	"envoyer/config"
	"envoyer/config/service_name"
	"envoyer/controller"
	"envoyer/dic"
	log "envoyer/logger"
	"envoyer/middlewares"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"net/http"
)

func Setup(_ *di.Builder) *gin.Engine {
	logger := dic.Container.Get(service_name.LoggerService).(log.Logger)
	logger.Info("logger is initialized", log.Extra("logLevel", logger.MinLevel()))

	gin.SetMode(config.Config.GinMode)

	r := gin.New()
	r.Use(gin.Recovery())

	// Set localization
	r.Use(ginI18n.Localize(ginI18n.WithBundle(&ginI18n.BundleCfg{
		RootPath:         "./localize",
		AcceptLanguage:   []language.Tag{language.Bengali, language.English},
		DefaultLanguage:  language.English,
		UnmarshalFunc:    yaml.Unmarshal,
		FormatBundleFile: "yaml",
	})))

	subscriberController := dic.Container.Get(service_name.SubscriberController).(*controller.SubscriberController)
	//start consuming and handling messages
	//go subscriberController.StartDefaultSubscriber()
	subscriberController.StartAllSubscribers()

	publisherController := dic.Container.Get(service_name.PublisherController).(*controller.PublisherController)
	appController := dic.Container.Get(service_name.AppController).(*controller.AppController)
	clientController := dic.Container.Get(service_name.ClientController).(*controller.ClientController)
	eventController := dic.Container.Get(service_name.EventController).(*controller.EventController)
	templateController := dic.Container.Get(service_name.TemplateController).(*controller.TemplateController)
	userController := dic.Container.Get(service_name.UserController).(*controller.UserController)
	authController := dic.Container.Get(service_name.AuthController).(*controller.AuthController)
	providerController := dic.Container.Get(service_name.ProviderController).(*controller.ProviderController)
	logController := dic.Container.Get(service_name.LogController).(*controller.LogController)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	api := r.Group("/api")
	{
		v2 := api.Group("/v2")
		{
			v2.POST("/publish/:type", publisherController.PublishInQueueV2)

			auth := v2.Group("/auth")
			{
				auth.POST("/login", authController.LogIn)
				auth.POST("/refresh", authController.RefreshAccessToken)
			}

			app := v2.Group("/app").Use(middlewares.Auth())
			{
				app.Use(middlewares.AppAdminPermission()).GET("/:id", appController.GetApp)         //all
				app.Use(middlewares.SuperAdminPermission()).GET("", appController.GetAllApp)        //super
				app.Use(middlewares.SuperAdminPermission()).POST("", appController.CreateApp)       //super
				app.Use(middlewares.SuperAdminPermission()).PUT("/:id", appController.UpdateApp)    //super
				app.Use(middlewares.SuperAdminPermission()).DELETE("/:id", appController.DeleteApp) //super
			}

			client := v2.Group("/client").Use(middlewares.Auth()).Use(middlewares.DevPermission()) //super, dev
			{
				client.POST("", clientController.CreateClient)
				client.GET("/:id", clientController.GetClient)
				client.PUT("/:id", clientController.UpdateClient)
				client.DELETE("/:id", clientController.DeleteClient)
				client.GET("/app_id/:id", clientController.GetByAppId)
			}

			event := v2.Group("/event").Use(middlewares.Auth())
			{
				event.Use(middlewares.AppAdminPermission()).GET("/:id", eventController.GetEvent)          //all
				event.Use(middlewares.AppAdminPermission()).GET("/app_id/:id", eventController.GetByAppId) //all
				event.Use(middlewares.DevPermission()).POST("", eventController.CreateEvent)               //super, dev
				event.Use(middlewares.DevPermission()).PUT("/:id", eventController.UpdateEvent)            //super, dev
				event.Use(middlewares.DevPermission()).DELETE("/:id", eventController.DeleteEvent)         //super, dev
			}

			template := v2.Group("/template").Use(middlewares.Auth()).Use(middlewares.AppAdminPermission()) //all
			{
				template.POST("", templateController.CreateTemplate)
				template.GET("/:id", templateController.GetTemplate)
				template.PUT("/:id", templateController.UpdateTemplate)
				template.DELETE("/:id", templateController.DeleteTemplate)
				template.GET("/event_id/:id", templateController.GetByEventId)
			}

			user := v2.Group("/user").Use(middlewares.Auth())
			{
				user.Use(middlewares.AppAdminPermission()).GET("/:id", userController.GetUser)         //all
				user.Use(middlewares.DevPermission()).GET("/app_id/:id", userController.GetByAppId)    //super, dev
				user.Use(middlewares.SuperAdminPermission()).POST("", userController.CreateUser)       //super
				user.Use(middlewares.SuperAdminPermission()).PUT("/:id", userController.UpdateUser)    //super
				user.Use(middlewares.SuperAdminPermission()).DELETE("/:id", userController.DeleteUser) //super
			}

			provider := v2.Group("/provider").Use(middlewares.Auth()).Use(middlewares.AppAdminPermission()) //all
			{
				provider.POST("", providerController.CreateProvider)
				provider.GET("/:id", providerController.GetProvider)
				provider.PUT("/:id", providerController.UpdateProvider)
				provider.DELETE("/:id", providerController.DeleteProvider)
				provider.GET("/app_id/:id", providerController.GetProviderByAppId)
				provider.GET("/app_id/:id/type/:type", providerController.GetProvidersByAppIdAndType)
				provider.PUT("/app_id/:id/type/:type", providerController.UpdateProviderPriority)
				provider.GET("/app_id/:id/type/:type/top", providerController.GetProviderByAppIdAndTypeWithTopPriority)
			}

			errorLog := v2.Group("/log").Use(middlewares.Auth()).Use(middlewares.DevPermission()) //super, dev
			{
				errorLog.DELETE("/:id", logController.DeleteLog)
				errorLog.GET("/app_id/:id", logController.GetByAppId)
			}
		}
	}

	return r
}
