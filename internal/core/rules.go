package core

import (
	"context"
	"mercury/internal/logger"
	"mercury/internal/models"
	"mercury/internal/storage"
)

type RuleService interface {
	Create(ctx context.Context, rule *models.Rule) error
	Get(ctx context.Context, id int) (*models.Rule, error)
	Update(ctx context.Context, rule *models.Rule) error
	Delete(ctx context.Context, id int) error
	ListByInbox(ctx context.Context, inboxID, limit, offset int) (*models.PaginatedResponse, error)
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

func (s *ruleService) Create(ctx context.Context, rule *models.Rule) error {
	if err := s.repo.CreateRule(ctx, rule); err != nil {
		s.logger.Error("Failed to create rule: %v", err)
		return err
	}
	s.logger.Info("Created rule %d for inbox %d", rule.ID, rule.InboxID)
	return nil
}

func (s *ruleService) Get(ctx context.Context, id int) (*models.Rule, error) {
	rule, err := s.repo.GetRule(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get rule %d: %v", id, err)
		return nil, err
	}
	s.logger.Debug("Retrieved rule: %d", id)
	return rule, nil
}

func (s *ruleService) Update(ctx context.Context, rule *models.Rule) error {
	if err := s.repo.UpdateRule(ctx, rule); err != nil {
		s.logger.Error("Failed to update rule %d: %v", rule.ID, err)
		return err
	}
	s.logger.Info("Updated rule: %d", rule.ID)
	return nil
}

func (s *ruleService) Delete(ctx context.Context, id int) error {
	if err := s.repo.DeleteRule(ctx, id); err != nil {
		s.logger.Error("Failed to delete rule %d: %v", id, err)
		return err
	}
	s.logger.Info("Deleted rule: %d", id)
	return nil
}

func (s *ruleService) ListByInbox(ctx context.Context, inboxID, limit, offset int) (*models.PaginatedResponse, error) {
	rules, total, err := s.repo.ListRulesByInbox(ctx, inboxID, limit, offset)
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
