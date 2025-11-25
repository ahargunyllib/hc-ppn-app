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

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, phone_number, label, assigned_to, notes, created_at, updated_at)
		VALUES (:id, :phone_number, :label, :assigned_to, :notes, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(
		ctx,
		query,
		user,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErrors := []pg.PgError{
				{
					Code:           pg.UniqueViolation,
					ConstraintName: "users_phone_number_key",
					Err: errx.ErrUserPhoneExists.WithDetails(map[string]any{
						"phone_number": user.PhoneNumber,
					}).WithLocation("userRepository.Create"),
				},
			}

			if customPgErr := pg.HandlePgError(pgErr, pgErrors); customPgErr != nil {
				return customPgErr
			}
		}

		return errx.ErrInternalServer.WithLocation("userRepository.Create").WithError(err)
	}

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, phone_number, label, assigned_to, notes, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user entity.User
	err := r.db.GetContext(ctx, &user, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrUserNotFound.WithDetails(map[string]any{
				"id": id,
			}).WithLocation("userRepository.FindByID")
		}

		return nil, errx.ErrInternalServer.WithLocation("userRepository.FindByID").WithError(err)
	}

	return &user, nil
}

func (r *userRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error) {
	query := `
		SELECT id, phone_number, label, assigned_to, notes, created_at, updated_at
		FROM users
		WHERE phone_number = $1
	`

	var user entity.User
	err := r.db.GetContext(ctx, &user, query, phoneNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrUserNotFound.WithDetails(map[string]any{
				"phone_number": phoneNumber,
			}).WithLocation("userRepository.FindByPhoneNumber")
		}

		return nil, errx.ErrInternalServer.WithLocation("userRepository.FindByPhoneNumber").WithError(err)
	}

	return &user, nil
}

func (r *userRepository) List(ctx context.Context, filter entity.GetUsersFilter) ([]entity.User, int64, error) {
	offset := min(max(filter.Offset, 0), 10000)
	limit := min(max(filter.Limit, 10), 100)

	var qb strings.Builder
	var args []any

	qb.WriteString(`
		SELECT id, phone_number, label, assigned_to, notes, created_at, updated_at
		FROM users
		WHERE 1=1
	`)

	if filter.Search != "" {
		qb.WriteString(fmt.Sprintf(" AND (phone_number ILIKE $%d OR label ILIKE $%d)", len(args)+1, len(args)+1))
		args = append(args, "%"+filter.Search+"%")
	}

	if filter.AssignedTo != nil && *filter.AssignedTo != "" {
		qb.WriteString(fmt.Sprintf(" AND assigned_to = $%d", len(args)+1))
		args = append(args, *filter.AssignedTo)
	}

	var total int64
	err := r.db.GetContext(ctx, &total, qb.String(), args...)
	if err != nil {
		return nil, 0, errx.ErrInternalServer.WithLocation("userRepository.List.Count").WithError(err)
	}

	qb.WriteString(" ORDER BY created_at DESC")
	qb.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2))

	args = append(args, limit, offset)

	var users []entity.User
	err = r.db.SelectContext(ctx, &users, qb.String(), args...)
	if err != nil {
		return nil, 0, errx.ErrInternalServer.WithLocation("userRepository.List.Select").WithError(err)
	}

	if users == nil {
		users = []entity.User{}
	}

	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET phone_number = :phone_number, label = :label, assigned_to = :assigned_to, notes = :notes, updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(
		ctx,
		query,
		user,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return errx.ErrUserPhoneExists.WithError(err)
		}

		return errx.ErrInternalServer.WithLocation("userRepository.Update").WithError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errx.ErrInternalServer.WithLocation("userRepository.Update.RowsAffected").WithError(err)
	}

	if rowsAffected == 0 {
		return errx.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errx.ErrInternalServer.WithLocation("userRepository.Delete").WithError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errx.ErrInternalServer.WithLocation("userRepository.Delete.RowsAffected").WithError(err)
	}

	if rowsAffected == 0 {
		return errx.ErrUserNotFound
	}

	return nil
}
