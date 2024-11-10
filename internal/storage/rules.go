package storage

import (
	"context"
	"database/sql"
	"errors"
	"inbox451/internal/models"
)

func (r *repository) ListRules(ctx context.Context, limit, offset int) ([]*models.ForwardRule, int, error) {
	var total int
	err := r.queries.CountRules.GetContext(ctx, &total)
	if err != nil {
		return nil, 0, err
	}

	var rules []*models.ForwardRule
	err = r.queries.ListRules.SelectContext(ctx, &rules, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

func (r *repository) ListRulesByInbox(ctx context.Context, inboxID, limit, offset int) ([]*models.ForwardRule, int, error) {
	var total int
	err := r.queries.CountRulesByInbox.GetContext(ctx, &total, inboxID)
	if err != nil {
		return nil, 0, err
	}

	var rules []*models.ForwardRule
	err = r.queries.ListRulesByInbox.SelectContext(ctx, &rules, inboxID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

func (r *repository) GetRule(ctx context.Context, id int) (*models.ForwardRule, error) {
	var rule models.ForwardRule
	err := r.queries.GetRule.GetContext(ctx, &rule, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

func (r *repository) CreateRule(ctx context.Context, rule *models.ForwardRule) error {
	return r.queries.CreateRule.QueryRowContext(ctx, rule.InboxID, rule.Sender, rule.Receiver, rule.Subject).
		Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
}

func (r *repository) UpdateRule(ctx context.Context, rule *models.ForwardRule) error {
	result, err := r.queries.UpdateRule.ExecContext(ctx, rule.Sender, rule.Receiver, rule.Subject, rule.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("rule not found")
	}
	return nil
}

func (r *repository) DeleteRule(ctx context.Context, id int) error {
	result, err := r.queries.DeleteRule.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("rule not found")
	}
	return nil
}
