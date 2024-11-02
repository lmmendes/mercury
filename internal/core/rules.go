package core

import (
	"mercury/internal/logger"
	"mercury/internal/models"
	"mercury/internal/storage"
)

type RuleService interface {
	Create(rule *models.Rule) error
	Get(id int) (*models.Rule, error)
	Update(rule *models.Rule) error
	Delete(id int) error
	ListByInbox(inboxID, limit, offset int) (*models.PaginatedResponse, error)
}

type ruleService struct {
	repo   storage.Repository
	logger *logger.Logger
}

func NewRuleService(core *Core) RuleService {
	return &ruleService{
		repo:   core.Repository,
		logger: core.Logger,
	}
}

func (s *ruleService) Create(rule *models.Rule) error {
	if err := s.repo.CreateRule(rule); err != nil {
		s.logger.Error("Failed to create rule: %v", err)
		return err
	}
	s.logger.Info("Created rule %d for inbox %d", rule.ID, rule.InboxID)
	return nil
}

func (s *ruleService) Get(id int) (*models.Rule, error) {
	rule, err := s.repo.GetRule(id)
	if err != nil {
		s.logger.Error("Failed to get rule %d: %v", id, err)
		return nil, err
	}
	s.logger.Debug("Retrieved rule: %d", id)
	return rule, nil
}

func (s *ruleService) Update(rule *models.Rule) error {
	if err := s.repo.UpdateRule(rule); err != nil {
		s.logger.Error("Failed to update rule %d: %v", rule.ID, err)
		return err
	}
	s.logger.Info("Updated rule: %d", rule.ID)
	return nil
}

func (s *ruleService) Delete(id int) error {
	if err := s.repo.DeleteRule(id); err != nil {
		s.logger.Error("Failed to delete rule %d: %v", id, err)
		return err
	}
	s.logger.Info("Deleted rule: %d", id)
	return nil
}

func (s *ruleService) ListByInbox(inboxID, limit, offset int) (*models.PaginatedResponse, error) {
	rules, total, err := s.repo.ListRulesByInbox(inboxID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list rules for inbox %d: %v", inboxID, err)
		return nil, err
	}

	response := &models.PaginatedResponse{
		Data: rules,
	}
	response.Pagination.Total = total
	response.Pagination.Limit = limit
	response.Pagination.Offset = offset

	s.logger.Debug("Retrieved %d rules for inbox %d (total: %d)", len(rules), inboxID, total)
	return response, nil
}
