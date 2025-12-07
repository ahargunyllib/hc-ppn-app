package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/errx"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/helpers/pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *feedbackRepository) Create(ctx context.Context, feedback *entity.Feedback) error {
	query := `
		INSERT INTO feedbacks (id, user_id, rating, comment, created_at)
		VALUES (:id, :user_id, :rating, :comment, :created_at)
	`

	_, err := r.db.NamedExecContext(
		ctx,
		query,
		feedback,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErrors := []pg.PgError{
				{
					Code:           pg.ForeignKey,
					ConstraintName: "fk_feedback_user",
					Err: errx.ErrUserNotFound.WithDetails(map[string]any{
						"user_id": feedback.UserID,
					}).WithLocation("feedbackRepository.Create"),
				},
			}

			if customPgErr := pg.HandlePgError(pgErr, pgErrors); customPgErr != nil {
				return customPgErr
			}
		}

		return errx.ErrInternalServer.WithLocation("feedbackRepository.Create").WithError(err)
	}

	return nil
}

func (r *feedbackRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Feedback, error) {
	query := `
		SELECT
			feedbacks.id,
			feedbacks.user_id,
			feedbacks.rating,
			feedbacks.comment,
			feedbacks.created_at,

			users.id AS "user.id",
			users.phone_number AS "user.phone_number",
			users.name AS "user.name",
			users.job_title AS "user.job_title",
			users.gender AS "user.gender",
			users.date_of_birth AS "user.date_of_birth",
			users.created_at AS "user.created_at",
			users.updated_at AS "user.updated_at"
		FROM feedbacks
		LEFT JOIN users ON feedbacks.user_id = users.id
		WHERE feedbacks.id = $1
	`

	var feedback entity.Feedback
	err := r.db.GetContext(ctx, &feedback, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrFeedbackNotFound.WithDetails(map[string]any{
				"id": id,
			}).WithLocation("feedbackRepository.FindByID")
		}

		return nil, errx.ErrInternalServer.WithLocation("feedbackRepository.FindByID").WithError(err)
	}

	return &feedback, nil
}

func (r *feedbackRepository) List(ctx context.Context, filter *entity.GetFeedbacksFilter) ([]entity.Feedback, int64, error) {
	offset := min(max(filter.Offset, 0), 10000)
	limit := min(max(filter.Limit, 10), 100)

	var qb strings.Builder
	var whereClauses strings.Builder
	var args []any

	qb.WriteString(`
		SELECT
			feedbacks.id,
			feedbacks.user_id,
			feedbacks.rating,
			feedbacks.comment,
			feedbacks.created_at,

			users.id AS "user.id",
			users.phone_number AS "user.phone_number",
			users.name AS "user.name",
			users.job_title AS "user.job_title",
			users.gender AS "user.gender",
			users.date_of_birth AS "user.date_of_birth",
			users.created_at AS "user.created_at",
			users.updated_at AS "user.updated_at"
		FROM feedbacks
		LEFT JOIN users ON feedbacks.user_id = users.id
	`)

	if filter.UserID != nil {
		whereClauses.WriteString(fmt.Sprintf(" AND user_id = $%d", len(args)+1))
		args = append(args, *filter.UserID)
	}

	if filter.MinRating != nil {
		whereClauses.WriteString(fmt.Sprintf(" AND rating >= $%d", len(args)+1))
		args = append(args, *filter.MinRating)
	}

	if filter.MaxRating != nil {
		whereClauses.WriteString(fmt.Sprintf(" AND rating <= $%d", len(args)+1))
		args = append(args, *filter.MaxRating)
	}

	var total int64
	err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM feedbacks WHERE 1=1"+whereClauses.String(), args...)
	if err != nil {
		return nil, 0, errx.ErrInternalServer.WithLocation("feedbackRepository.List.Count").WithError(err)
	}

	if whereClauses.Len() > 0 {
		qb.WriteString(" WHERE 1=1")
		qb.WriteString(whereClauses.String())
	}
	qb.WriteString(" ORDER BY created_at DESC")
	qb.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2))

	args = append(args, limit, offset)

	var feedbacks []entity.Feedback
	err = r.db.SelectContext(ctx, &feedbacks, qb.String(), args...)
	if err != nil {
		return nil, 0, errx.ErrInternalServer.WithLocation("feedbackRepository.List.Select").WithError(err)
	}

	if feedbacks == nil {
		feedbacks = []entity.Feedback{}
	}

	return feedbacks, total, nil
}

func (r *feedbackRepository) GetMetrics(ctx context.Context) (float64, error) {
	query := `
		SELECT
			COALESCE(
				(SUM(rating) * 100.0) / NULLIF(COUNT(*) * 5, 0),
				0
			) AS satisfaction_score
		FROM feedbacks
	`

	var satisfactionScore float64
	err := r.db.GetContext(ctx, &satisfactionScore, query)
	if err != nil {
		return 0, errx.ErrInternalServer.WithLocation("feedbackRepository.GetMetrics").WithError(err)
	}

	return satisfactionScore, nil
}

func (r *feedbackRepository) GetSatisfactionTrend(ctx context.Context) ([]entity.SatisfactionTrendRow, error) {
	query := `
		WITH date_series AS (
			SELECT generate_series(
				CURRENT_DATE - 30 * INTERVAL '1 day',
				CURRENT_DATE,
				INTERVAL '1 day'
			)::date AS date
		)
		SELECT
			ds.date,
			COALESCE(AVG(f.rating), 0) AS avg_satisfaction
		FROM date_series ds
		LEFT JOIN feedbacks f ON DATE(f.created_at) = ds.date
		GROUP BY ds.date
		ORDER BY ds.date ASC
	`

	var results []entity.SatisfactionTrendRow
	err := r.db.SelectContext(ctx, &results, query)
	if err != nil {
		return nil, errx.ErrInternalServer.WithLocation("feedbackRepository.GetSatisfactionTrend").WithError(err)
	}

	if results == nil {
		results = []entity.SatisfactionTrendRow{}
	}

	return results, nil
}
