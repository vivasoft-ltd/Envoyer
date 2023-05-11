package dic

import (
	"envoyer/config"
	"envoyer/config/service_name"
	"envoyer/controller"
	"envoyer/logger"
	"envoyer/model/db"
	"envoyer/model/repository"
	"envoyer/model/service"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

var Builder *di.Builder
var Container di.Container

func InitContainer() di.Container {
	builder := InitBuilder()
	Container = builder.Build()
	return Container
}

func InitBuilder() *di.Builder {
	Builder, _ = di.NewBuilder()
	RegisterServices(Builder)
	return Builder
}

func RegisterServices(builder *di.Builder) {
	_ = builder.Add(di.Def{
		Name: service_name.LogWriter,
		Build: func(ctn di.Container) (interface{}, error) {
			return logger.NewZapWriter()
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.LoggerService,
		Build: func(ctn di.Container) (interface{}, error) {
			return logger.NewLogger(config.Config.LogLevel, ctn.Get(service_name.LogWriter).(logger.Writer)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.DbService,
		Build: func(ctn di.Container) (interface{}, error) {
			return db.NewDb(), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.BaseRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewBaseRepository(ctn.Get(service_name.DbService).(*gorm.DB), ctn.Get(service_name.LoggerService).(logger.Logger)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.BaseService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewBaseService(ctn, ctn.Get(service_name.LoggerService).(logger.Logger)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.BaseController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewBaseController(ctn, ctn.Get(service_name.LoggerService).(logger.Logger)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.ConsumerService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewRabbitMQSubscriber(ctn.Get(service_name.LoggerService).(logger.Logger))
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.DispatcherService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewDispatcherService(ctn.Get(service_name.ConsumerService).(service.Subscriber), ctn.Get(service_name.LoggerService).(logger.Logger), ctn.Get(service_name.HandlerService).(service.HandlerServiceInterface), ctn), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.SubscriberController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewSubscriberController(ctn.Get(service_name.BaseController).(*controller.BaseController)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.PublisherService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewRabbitMQPublisher(ctn.Get(service_name.LoggerService).(logger.Logger))
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.PublisherController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewPublisherController(ctn.Get(service_name.BaseController).(*controller.BaseController)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.AppRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewAppRepository(ctn.Get(service_name.BaseRepository).(*repository.BaseRepository)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.AppService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewAppService(ctn.Get(service_name.BaseService).(*service.BaseService)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.AppController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewAppController(ctn.Get(service_name.BaseController).(*controller.BaseController), ctn.Get(service_name.SubscriberController).(*controller.SubscriberController)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.UserRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewUserRepository(ctn.Get(service_name.BaseRepository).(*repository.BaseRepository)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.UserService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewUserService(ctn.Get(service_name.BaseService).(*service.BaseService)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.UserController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewUserController(ctn.Get(service_name.BaseController).(*controller.BaseController)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.AuthService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewAuthService(ctn.Get(service_name.BaseService).(*service.BaseService)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.AuthController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewAuthController(ctn.Get(service_name.BaseController).(*controller.BaseController)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.ClientRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewClientRepository(ctn.Get(service_name.BaseRepository).(*repository.BaseRepository)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.ClientService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewClientService(ctn.Get(service_name.BaseService).(*service.BaseService)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.ClientController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewClientController(ctn.Get(service_name.BaseController).(*controller.BaseController)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.EventRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewEventRepository(ctn.Get(service_name.BaseRepository).(*repository.BaseRepository)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.EventService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewEventService(ctn.Get(service_name.BaseService).(*service.BaseService)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.EventController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewEventController(ctn.Get(service_name.BaseController).(*controller.BaseController)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.TemplateRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewTemplateRepository(ctn.Get(service_name.BaseRepository).(*repository.BaseRepository)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.TemplateService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewTemplateService(ctn.Get(service_name.BaseService).(*service.BaseService)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.TemplateController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewTemplateController(ctn.Get(service_name.BaseController).(*controller.BaseController)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.ProviderRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewProviderRepository(ctn.Get(service_name.BaseRepository).(*repository.BaseRepository)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.ProviderService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewProviderService(ctn.Get(service_name.BaseService).(*service.BaseService)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.ProviderController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewProviderController(ctn.Get(service_name.BaseController).(*controller.BaseController)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.HandlerService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewHandlerService(ctn.Get(service_name.BaseService).(*service.BaseService)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.LogRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repository.NewLogRepository(ctn.Get(service_name.BaseRepository).(*repository.BaseRepository)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.LogService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewLogService(ctn.Get(service_name.BaseService).(*service.BaseService)), nil
		},
	})

	_ = builder.Add(di.Def{
		Name: service_name.LogController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewLogController(ctn.Get(service_name.BaseController).(*controller.BaseController)), nil
		},
	})
}
