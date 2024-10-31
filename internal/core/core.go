package core

import (
	"database/sql"
	"mercury/internal/logger"
	"mercury/internal/models"
	"mercury/internal/storage"
	"os"
)

type Config struct {
	SMTPPort    string
	HTTPPort    string
	DatabaseURL string
	LogLevel    logger.Level
}

type Core struct {
	Config     *Config
	Logger     *logger.Logger
	Repository storage.Repository

	AccountService AccountService
	InboxService   InboxService
	RuleService    RuleService
	MessageService MessageService
}

func NewCore(config *Config, db *sql.DB) *Core {
	// Create logger with stdout only
	logger := logger.New(os.Stdout, config.LogLevel)

	core := &Core{
		Config:     config,
		Logger:     logger,
		Repository: storage.NewRepository(db),
	}

	core.AccountService = NewAccountService(core)
	core.InboxService = NewInboxService(core)
	core.RuleService = NewRuleService(core)
	core.MessageService = NewMessageService(core)

	return core
}

func LoadConfig() *Config {
	return &Config{
		SMTPPort:    ":1025",
		HTTPPort:    ":8080",
		DatabaseURL: "./email.db",
		LogLevel:    logger.INFO, // Default log level
	}
}

func (c *Core) StoreMessage(message *models.Message) error {
	return c.MessageService.Store(message)
}
