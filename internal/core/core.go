package core

import (
	"database/sql"
	"log"
	"mercury/internal/models"
	"mercury/internal/storage"
)

type Config struct {
	SMTPPort    string
	HTTPPort    string
	DatabaseURL string
}

type Core struct {
	Config     *Config
	Logger     *log.Logger
	Repository storage.Repository

	AccountService AccountService
	InboxService   InboxService
	RuleService    RuleService
	MessageService MessageService
}

func NewCore(config *Config, db *sql.DB, logger *log.Logger) *Core {
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
	}
}

func (c *Core) StoreMessage(message *models.Message) error {
	return c.MessageService.Store(message)
}
