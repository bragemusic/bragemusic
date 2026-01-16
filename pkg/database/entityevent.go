package database

import (
	"context"
	"time"

	"github.com/bragemusic/core/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) addEntityEvent(ctx context.Context, id uuid.UUID, eventType types.EntityEventType, entityType types.EntityType) error {
	const query = `
		INSERT INTO entity_events (
			id, event_type, entity_type, event_time
		) VALUES (?, ?, ?, ?);
	`

	_, err := d.ext.ExecContext(
		ctx,
		query,
		id,
		eventType,
		entityType,
		time.Now(),
	)

	return err
}

func (d Database) ListEntityEvents(ctx context.Context, eventType types.EntityEventType, entityType types.EntityType, since time.Time) (ids []uuid.UUID, err error) {
	query := `
        SELECT id
        FROM entity_events
        WHERE
          event_type = ?
          AND
          entity_type = ?
          AND
          event_time > ?
        ;
    `
	err = sqlx.SelectContext(ctx, d.ext, &ids, query, eventType, entityType, since, since)
	if err != nil {
		return nil, err
	}

	return
}
