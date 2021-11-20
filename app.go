package main

import (
	"fmt"
	"net/http"
	"time"

	ginzerolog "github.com/easonlin404/gin-zerolog"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/v-denis-e/gojwt/conf"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// App представляет главный контекст приложения
type App struct {
	config conf.AppConf // Конфигурация приложения
	db     *gorm.DB     // Источник данных с БД
	router *gin.Engine  // Роутер приложения
}

// NewApp создает новый экземпляр приложения
// с конфигурацией из переменных среды
func NewApp() *App {
	config, err := conf.FromEnv()
	if err != nil {
		log.Fatal().Err(err)
	}

	return &App{config: config}
}

// Init инициализирует главный контекст приложения
func (a *App) Init() {
	log.Info().Msg("Initializing application context ...")

	a.initDb()
	a.initRouter()
}

// initDb инициализирует подключение к базе данных
func (a *App) initDb() {
	log.Info().Msg("Initializing database connectivity ...")

	// Использование zerolog вместо стандартного логера
	logger := logger.New(
		&log.Logger,
		logger.Config{
			SlowThreshold:             time.Second, // Предел до секунды
			Colorful:                  false,       // Не использовать цвета при выводе логов
			IgnoreRecordNotFoundError: true,        // Не логировать ErrRecordNotFound
			LogLevel:                  logger.Info, // Писать только на уровне Info
		},
	)

	dsn := fmt.Sprintf(dsnTemplate, a.config.DB.Host, a.config.DB.Port, a.config.DB.User, a.config.DB.Password, a.config.DB.Name)

	var err error
	a.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Использовать наименование таблиц в единственном числе
		},
		Logger: logger,
	})
	fatal(err)

	log.Info().Msg("Database inited.")
}

// initRouter инициализирует движок gin
func (a *App) initRouter() {
	log.Info().Msg("Initializing application router ...")

	a.router = gin.New()

	// Gin должен использовать zerolog для логирования
	a.router.Use(ginzerolog.Logger())

	// Восстановление после panic по умолчанию
	a.router.Use(gin.Recovery())

	// Добавление параметров безопсности
	a.router.Use(secure.New(secure.Config{
		AllowedHosts:          []string{"localhost", "localhost:9000"}, // CORS allowings
		SSLRedirect:           false,                                   // use http
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true, // J Frame restriction
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,                 // cross site script block
		ContentSecurityPolicy: "default-src 'self'", // This must be better content security policy
		IsDevelopment:         true,                 // This is development environment (for prod disable this)
		IENoOpen:              true,
	}))

	a.initRoutes()
}

// initRoutes инциализирует все маршруты приложения
func (a *App) initRoutes() {
	a.router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}

// Run запускает приложение по указанному адресу
func (a *App) Run(addr string) {
	fatal(a.router.Run(addr))
}

func fatal(err error) {
	if err != nil {
		log.Fatal().Err(err)
	}
}

// Шаблон формата подключения к postgres
var dsnTemplate = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
