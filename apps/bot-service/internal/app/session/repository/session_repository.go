package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/errx"
	"github.com/google/uuid"
)

func (r *sessionRepository) Create(ctx context.Context, session *entity.ConversationSession) error {
	query := `
		INSERT INTO conversation_sessions (id, phone_number, status, last_message_at, feedback_prompt_sent_at, created_at, updated_at)
		VALUES (:id, :phone_number, :status, :last_message_at, :feedback_prompt_sent_at, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(
		ctx,
		query,
		session,
	)

	if err != nil {
		return errx.ErrInternalServer.WithLocation("sessionRepository.Create").WithError(err)
	}

	return nil
}

func (r *sessionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.ConversationSession, error) {
	query := `
		SELECT id, phone_number, status, last_message_at, feedback_prompt_sent_at, created_at, updated_at
		FROM conversation_sessions
		WHERE id = $1
	`

	var session entity.ConversationSession
	err := r.db.GetContext(ctx, &session, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrSessionNotFound.WithDetails(map[string]any{
				"id": id,
			}).WithLocation("sessionRepository.FindByID")
		}

		return nil, errx.ErrInternalServer.WithLocation("sessionRepository.FindByID").WithError(err)
	}

	return &session, nil
}

func (r *sessionRepository) FindActiveByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.ConversationSession, error) {
	query := `
		SELECT id, phone_number, status, last_message_at, feedback_prompt_sent_at, created_at, updated_at
		FROM conversation_sessions
		WHERE phone_number = $1 AND status IN ('active', 'waiting_feedback')
		ORDER BY created_at DESC
		LIMIT 1
	`

	var session entity.ConversationSession
	err := r.db.GetContext(ctx, &session, query, phoneNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errx.ErrInternalServer.WithLocation("sessionRepository.FindActiveByPhoneNumber").WithError(err)
	}

	return &session, nil
}

func (r *sessionRepository) List(ctx context.Context, filter *entity.GetSessionsFilter) ([]entity.ConversationSession, int64, error) {
	offset := min(max(filter.Offset, 0), 10000)
	limit := min(max(filter.Limit, 10), 100)

	var qb strings.Builder
	var whereClauses strings.Builder
	var args []any

	qb.WriteString(`
		SELECT id, phone_number, status, last_message_at, feedback_prompt_sent_at, created_at, updated_at
		FROM conversation_sessions
	`)

	if filter.PhoneNumber != nil && *filter.PhoneNumber != "" {
		whereClauses.WriteString(fmt.Sprintf(" AND phone_number = $%d", len(args)+1))
		args = append(args, *filter.PhoneNumber)
	}

	if filter.Status != nil {
		whereClauses.WriteString(fmt.Sprintf(" AND status = $%d", len(args)+1))
		args = append(args, *filter.Status)
	}

	var total int64
	err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM conversation_sessions WHERE 1=1"+whereClauses.String(), args...)
	if err != nil {
		return nil, 0, errx.ErrInternalServer.WithLocation("sessionRepository.List.Count").WithError(err)
	}

	if whereClauses.Len() > 0 {
		qb.WriteString(" WHERE 1=1")
		qb.WriteString(whereClauses.String())
	}
	qb.WriteString(" ORDER BY created_at DESC")
	qb.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2))

	args = append(args, limit, offset)

	var sessions []entity.ConversationSession
	err = r.db.SelectContext(ctx, &sessions, qb.String(), args...)
	if err != nil {
		return nil, 0, errx.ErrInternalServer.WithLocation("sessionRepository.List.Select").WithError(err)
	}

	if sessions == nil {
		sessions = []entity.ConversationSession{}
	}

	return sessions, total, nil
}

func (r *sessionRepository) Update(ctx context.Context, session *entity.ConversationSession) error {
	query := `
		UPDATE conversation_sessions
		SET status = :status, last_message_at = :last_message_at, feedback_prompt_sent_at = :feedback_prompt_sent_at, updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(
		ctx,
		query,
		session,
	)

	if err != nil {
		return errx.ErrInternalServer.WithLocation("sessionRepository.Update").WithError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errx.ErrInternalServer.WithLocation("sessionRepository.Update.RowsAffected").WithError(err)
	}

	if rowsAffected == 0 {
		return errx.ErrSessionNotFound.WithDetails(map[string]any{
			"id": session.ID,
		}).WithLocation("sessionRepository.Update")
	}

	return nil
}

func (r *sessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM conversation_sessions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errx.ErrInternalServer.WithLocation("sessionRepository.Delete").WithError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errx.ErrInternalServer.WithLocation("sessionRepository.Delete.RowsAffected").WithError(err)
	}

	if rowsAffected == 0 {
		return errx.ErrSessionNotFound.WithDetails(map[string]any{
			"id": id,
		}).WithLocation("sessionRepository.Delete")
	}

	return nil
}

func (r *sessionRepository) FindExpiredSessions(ctx context.Context, inactiveMinutes int) ([]entity.ConversationSession, error) {
	query := `
		SELECT id, phone_number, status, last_message_at, feedback_prompt_sent_at, created_at, updated_at
		FROM conversation_sessions
		WHERE status = 'active'
			AND last_message_at < NOW() - INTERVAL '1 minute' * $1
	`

	var sessions []entity.ConversationSession
	err := r.db.SelectContext(ctx, &sessions, query, inactiveMinutes)
	if err != nil {
		return nil, errx.ErrInternalServer.WithLocation("sessionRepository.FindExpiredSessions").WithError(err)
	}

	if sessions == nil {
		sessions = []entity.ConversationSession{}
	}

	return sessions, nil
}
