package repository

import (
	"context"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/errx"
)

func (r *topicRepository) BulkCreate(ctx context.Context, topics []entity.Topic) error {
	query := `
		INSERT INTO topics (title, count)
		VALUES (:title, :count)
	`

	_, err := r.db.NamedExecContext(
		ctx,
		query,
		topics,
	)
	if err != nil {
		return errx.ErrInternalServer.WithLocation("topicRepository.BulkCreate").WithError(err)
	}

	return nil
}

func (r *topicRepository) GetHotTopics(ctx context.Context) ([]entity.Topic, error) {
	query := `
		SELECT
			id,
			title,
			count,
			created_at
		FROM topics
		WHERE created_at >= CURRENT_DATE - 30 * INTERVAL '1 day'
		GROUP BY title
		ORDER BY count DESC
		LIMIT 5
	`

	var results []entity.Topic
	err := r.db.SelectContext(ctx, &results, query)
	if err != nil {
		return nil, errx.ErrInternalServer.WithLocation("topicRepository.GetHotTopics").WithError(err)
	}

	if results == nil {
		results = []entity.Topic{}
	}

	return results, nil
}
