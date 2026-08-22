package database

import (
	"context"
	"time"

	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

func (d Database) addEntityEvent(ctx context.Context, id uuid.UUID, eventType types.EntityEventType, entityType types.EntityType, userID uuid.UUID) error {
	eid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO entity_events (
			id, item_id, user_id, event_type, entity_type, event_time
		) VALUES (?, ?, ?, ?, ?, ?);
	`

	_, err = d.ext.ExecContext(
		ctx,
		query,
		eid,
		id,
		userID,
		eventType,
		entityType,
		time.Now(),
	)

	return err
}

func (d Database) ListEntityEvents(ctx context.Context, since time.Time) (events []types.EntityEvent, err error) {
	query := `
        SELECT *
        FROM entity_events
        WHERE
          event_time > ?
        ORDER BY event_time DESC
        ;
    `
	err = sqlx.SelectContext(ctx, d.ext, &events, query, since)
	if err != nil {
		return nil, err
	}

	return
}

func (d Database) ListEntityEventsByType(ctx context.Context, eventType types.EntityEventType, entityType types.EntityType, since time.Time) (ids []uuid.UUID, err error) {
	query := `
        SELECT item_id
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

func (d Database) ListEntityEventsByEntityType(ctx context.Context, entityType types.EntityType, since time.Time, userID *uuid.UUID) (ids []types.EntityEvent, err error) {
	if userID == nil {
		query := `
        SELECT *
        FROM entity_events
        WHERE
          entity_type = ?
          AND
          event_time > ?
        ORDER BY event_time
        ;
    `
		err = sqlx.SelectContext(ctx, d.ext, &ids, query, entityType, since)
		if err != nil {
			return nil, err
		}

		return
	}

	query := `
        SELECT *
        FROM entity_events
        WHERE
          entity_type = ?
          AND
          event_time > ?
          AND
          user_id = ?
        ORDER BY event_time
        ;
    `
	err = sqlx.SelectContext(ctx, d.ext, &ids, query, entityType, since, *userID)
	if err != nil {
		return nil, err
	}

	return
}
