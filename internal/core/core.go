package core

import (
	"context"
	"fmt"
	"mercury/internal/config"
	"mercury/internal/logger"
	"mercury/internal/models"
	"mercury/internal/storage"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Core struct {
	Config     *config.Config
	Logger     *logger.Logger
	Repository storage.Repository

	ProjectService ProjectService
	InboxService   InboxService
	RuleService    RuleService
	MessageService MessageService
}

func NewCore(cfg *config.Config, db *sqlx.DB) (*Core, error) {
	baseLogger := logger.New(os.Stdout, cfg.Logging.Level)

	repo, err := storage.NewRepository(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	core := &Core{
		Config:     cfg,
		Logger:     baseLogger,
		Repository: repo,
	}

	core.ProjectService = NewProjectService(core)
	core.InboxService = NewInboxService(core)
	core.RuleService = NewRuleService(core)
	core.MessageService = NewMessageService(core)

	return core, nil
}

func (c *Core) StoreMessage(message *models.Message) error {
	ctx := context.Background()
	return c.MessageService.Store(ctx, message)
}
