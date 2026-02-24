package router

import (
	"community-elderly-care-platform/internal/config"
	"community-elderly-care-platform/internal/handler"
	"community-elderly-care-platform/internal/notify"
	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"
	"community-elderly-care-platform/internal/storage"
	"community-elderly-care-platform/pkg/database"
	"community-elderly-care-platform/pkg/jwt"
	"community-elderly-care-platform/pkg/redis"
)

// Deps 依赖注入容器
type Deps struct {
	// Config
	Config *config.Config

	// Infrastructure
	DB         *database.DB
	Redis      *redis.Client
	JWTManager *jwt.Manager

	// Repositories
	UserRepo           *repository.UserRepository
	UserIdentityRepo   *repository.UserIdentityRepository
	CustomerRepo       *repository.CustomerRepository
	StationRepo        *repository.StationRepository
	ZoneRepo           *repository.ZoneRepository
	RequestRepo        *repository.RequestRepository
	TaskRepo           *repository.TaskRepository
	NotificationRepo   *repository.NotificationRepository
	NewsRepo           *repository.NewsRepository
	MenuRepo           *repository.MenuRepository
	RoleRepo           *repository.RoleRepository
	PermissionRepo     *repository.PermissionRepository
	RolePermissionRepo *repository.RolePermissionRepository
	ReportRepo         *repository.ReportRepository

	// Services
	ElderlyService      *service.ElderlyService
	AuthService         *service.AuthService
	SMSService          *service.SMSService
	BannerService       *service.BannerService
	GeocodeService      *service.GeocodeService
	GeofenceService     *service.GeofenceService
	StationService      *service.StationService
	ZoneService         *service.ZoneService
	RequestService      *service.RequestService
	TaskService         *service.TaskService
	NotificationService *service.NotificationService
	UserService         *service.UserService
	UserIdentityService *service.UserIdentityService
	NewsService         *service.NewsService
	MenuService         *service.MenuService
	PermissionService   *service.PermissionService
	RoleService         *service.RoleService
	BlacklistService    *service.TokenBlacklistService
	StorageService      *service.StorageService
	StatisticsService   *service.StatisticsService
	ReportService       *service.ReportService

	// Handlers
	CustomerHandler     *handler.CustomerHandler
	BannerHandler       *handler.BannerHandler
	BAuthHandler        *handler.BAuthHandler
	CAuthHandler        *handler.CAuthHandler
	CProfileHandler     *handler.CProfileHandler
	StationHandler      *handler.StationHandler
	ZoneHandler         *handler.ZoneHandler
	RequestHandler      *handler.RequestHandler
	TaskHandler         *handler.TaskHandler
	NotificationHandler *handler.NotificationHandler
	UploadHandler       *handler.UploadHandler
	UserHandler         *handler.UserHandler
	LogoutHandler       *handler.LogoutHandler
	PermissionHandler   *handler.PermissionHandler
	RoleHandler         *handler.RoleHandler
	NewsHandler         *handler.NewsHandler
	MenuHandler         *handler.MenuHandler
	StatisticsHandler   *handler.StatisticsHandler
	ReportHandler       *handler.ReportHandler
}

// NewDeps 创建并初始化依赖
func NewDeps(db *database.DB, rdb *redis.Client, cfg *config.Config) (*Deps, error) {
	d := &Deps{
		Config: cfg,
		DB:     db,
		Redis:  rdb,
	}

	// 初始化 Repositories
	bannerRepo := repository.NewBannerRepository(db.DB)
	d.UserRepo = repository.NewUserRepository(db.DB)
	d.UserIdentityRepo = repository.NewUserIdentityRepository(db.DB)
	d.CustomerRepo = repository.NewCustomerRepository(db.DB)
	d.StationRepo = repository.NewStationRepository(db.DB)
	d.ZoneRepo = repository.NewZoneRepository(db.DB)
	d.RequestRepo = repository.NewRequestRepository(db.DB)
	d.TaskRepo = repository.NewTaskRepository(db.DB)
	d.NotificationRepo = repository.NewNotificationRepository(db.DB)
	d.NewsRepo = repository.NewNewsRepository(db.DB)
	d.MenuRepo = repository.NewMenuRepository(db.DB)
	d.RoleRepo = repository.NewRoleRepository(db.DB)
	d.PermissionRepo = repository.NewPermissionRepository(db.DB)
	d.RolePermissionRepo = repository.NewRolePermissionRepository(db.DB)
	d.ReportRepo = repository.NewReportRepository(db.DB)

	// 初始化 JWT Manager
	d.JWTManager = jwt.NewManager(cfg.JWT.Secret, cfg.JWT.ExpiresIn, cfg.JWT.RefreshExpiresIn)

	// 初始化基础 Services
	d.SMSService = service.NewSMSService(rdb, cfg.Server.Mode)
	d.GeocodeService = service.NewGeocodeService(cfg.Amap.Key)
	d.BlacklistService = service.NewTokenBlacklistService(rdb)

	// 初始化权限服务（替代 Casbin）
	d.PermissionService = service.NewPermissionService(
		db.DB,
		d.PermissionRepo,
		d.RoleRepo,
		d.RolePermissionRepo,
		d.BlacklistService,
	)
	d.RoleService = service.NewRoleService(d.RoleRepo)

	// 初始化 GeofenceService
	d.GeofenceService = service.NewGeofenceService(d.ZoneRepo)
	if err := d.GeofenceService.Reload(); err != nil {
		return nil, err
	}

	// 初始化 AuthService
	d.AuthService = service.NewAuthService(d.UserRepo, d.UserIdentityRepo, d.CustomerRepo, d.JWTManager, d.SMSService, db.DB)
	d.AuthService.SetGeofenceService(d.GeofenceService)
	d.AuthService.SetStationRepo(d.StationRepo)

	// 初始化邮件发送器
	var mailSender notify.MailSender
	if cfg.Mail.SMTPHost != "" && cfg.Mail.SMTPPort != 0 {
		sender, err := notify.NewSMTPSender(cfg.Mail.SMTPHost, cfg.Mail.SMTPPort, cfg.Mail.SMTPUser, cfg.Mail.SMTPPass, cfg.Mail.SMTPTLS)
		if err != nil {
			return nil, err
		}
		mailSender = sender
	}

	// 初始化存储
	storageProvider, err := storage.NewProvider(cfg.Storage)
	if err != nil {
		return nil, err
	}

	// 初始化业务 Services
	d.StationService = service.NewStationService(d.StationRepo)
	d.StationService.SetCustomerRepo(d.CustomerRepo)
	d.StationService.SetGeofenceService(d.GeofenceService)
	d.StationService.SetGeocodeService(d.GeocodeService)

	d.ZoneService = service.NewZoneService(d.ZoneRepo, d.GeofenceService)
	d.NotificationService = service.NewNotificationService(d.NotificationRepo, d.UserRepo, mailSender)
	d.RequestService = service.NewRequestService(db.DB, d.RequestRepo, d.TaskRepo, d.StationRepo, d.GeofenceService, d.GeocodeService, d.NotificationService, d.UserIdentityRepo)
	d.TaskService = service.NewTaskService(db.DB, d.TaskRepo, d.RequestRepo, d.NotificationService)
	d.StorageService = service.NewStorageService(storageProvider)
	d.UserService = service.NewUserService(d.UserRepo, d.UserIdentityRepo)
	d.UserIdentityService = service.NewUserIdentityService(d.UserIdentityRepo, d.BlacklistService, db.DB)
	d.NewsService = service.NewNewsService(d.NewsRepo)
	d.MenuService = service.NewMenuService(d.MenuRepo)
	d.BannerService = service.NewBannerService(bannerRepo)
	d.StatisticsService = service.NewStatisticsService(d.TaskRepo, d.RequestRepo, d.UserRepo)
	d.ReportService = service.NewReportService(d.ReportRepo, d.RequestRepo, d.TaskRepo, d.StationRepo, d.UserRepo, cfg.Storage.Local.BasePath)
	d.ElderlyService = service.NewElderlyService(db.DB, d.UserRepo, d.UserIdentityRepo, d.CustomerRepo, d.RequestRepo, d.TaskRepo)

	// 初始化 Handlers
	d.CustomerHandler = handler.NewCustomerHandler(d.ElderlyService)
	d.BannerHandler = handler.NewBannerHandler(d.BannerService)
	d.BAuthHandler = handler.NewBAuthHandler(d.AuthService, d.UserRepo, d.UserIdentityRepo, d.PermissionService)
	d.CAuthHandler = handler.NewCAuthHandler(d.AuthService, d.UserRepo, d.CustomerRepo, d.SMSService)
	d.CProfileHandler = handler.NewCProfileHandler(d.CustomerRepo, d.GeocodeService)
	d.StationHandler = handler.NewStationHandler(d.StationService)
	d.ZoneHandler = handler.NewZoneHandler(d.ZoneService)
	d.RequestHandler = handler.NewRequestHandler(d.RequestService)
	d.TaskHandler = handler.NewTaskHandler(d.TaskService)
	d.NotificationHandler = handler.NewNotificationHandler(d.NotificationService)
	d.UploadHandler = handler.NewUploadHandler(d.StorageService)
	d.UserHandler = handler.NewUserHandler(d.UserService, cfg.JWT.Secret, cfg.Security.IDCardEncryptKey)
	d.LogoutHandler = handler.NewLogoutHandler(d.BlacklistService)
	d.PermissionHandler = handler.NewPermissionHandler(d.PermissionService)
	d.RoleHandler = handler.NewRoleHandler(d.PermissionService, d.UserIdentityService)
	d.NewsHandler = handler.NewNewsHandler(d.NewsService)
	d.MenuHandler = handler.NewMenuHandler(d.MenuService, d.UserRepo, d.UserIdentityRepo, d.PermissionService)
	d.StatisticsHandler = handler.NewStatisticsHandler(d.StatisticsService)
	d.ReportHandler = handler.NewReportHandler(d.ReportService)

	return d, nil
}
