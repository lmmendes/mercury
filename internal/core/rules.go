package core

import (
	"log"
	"mercury/internal/models"
	"mercury/internal/storage"
)

type RuleService interface {
	Create(rule *models.Rule) error
	Get(id int) (*models.Rule, error)
	Update(rule *models.Rule) error
	Delete(id int) error
	GetByInboxID(inboxID int) ([]*models.Rule, error)
}

type ruleService struct {
	repo   storage.Repository
	logger *log.Logger
}

func NewRuleService(core *Core) RuleService {
	return &ruleService{
		repo:   core.Repository,
		logger: core.Logger,
	}
}

func (s *ruleService) Create(rule *models.Rule) error {
	return s.repo.CreateRule(rule)
}

func (s *ruleService) Get(id int) (*models.Rule, error) {
	return s.repo.GetRule(id)
}

func (s *ruleService) Update(rule *models.Rule) error {
	return s.repo.UpdateRule(rule)
}

func (s *ruleService) Delete(id int) error {
	return s.repo.DeleteRule(id)
}

func (s *ruleService) GetByInboxID(inboxID int) ([]*models.Rule, error) {
	return s.repo.ListRulesByInbox(inboxID)
}
