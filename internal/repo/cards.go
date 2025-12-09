package repo

import (
	"context"
	"diploma/internal/models"
	"diploma/pkg/db"
	"fmt"
	"time"
)

// InsertUpdateCard создаёт или обновляет запись карты.
func InsertUpdateCard(ctx context.Context, conn *db.SqlConnection, rec *models.CardRecord) error {
	if rec == nil {
		return fmt.Errorf("nil CardRecord")
	}

	const q = `
		INSERT INTO cards (
			user_id, card_pan, data, meta, created_at, updated_at, is_deleted
		) VALUES (
			$1, $2, $3, $4, COALESCE($5, now()), now(), $6
		)
		ON CONFLICT (user_id, card_pan) DO UPDATE SET
			data       = EXCLUDED.data,
			meta       = EXCLUDED.meta,
			updated_at = now(),
			is_deleted = EXCLUDED.is_deleted
	`

	created := rec.CreatedAt
	if created.IsZero() {
		_, err := db.ExecuteDBExec(ctx, conn.PgSql, conn.Timeout, q,
			rec.UserID,
			rec.CardPAN,
			rec.Data,
			rec.Meta,
			nil,
			rec.IsDeleted,
		)
		return err
	}

	_, err := db.ExecuteDBExec(ctx, conn.PgSql, conn.Timeout, q,
		rec.UserID,
		rec.CardPAN,
		rec.Data,
		rec.Meta,
		rec.CreatedAt,
		rec.IsDeleted,
	)
	return err
}

// GetCard возвращает одну запись по user_id + card_pan.
func GetCard(ctx context.Context, conn *db.SqlConnection, userID int64, cardPAN string) (*models.CardRecord, error) {
	const q = `
		SELECT user_id, card_pan, data, meta, created_at, updated_at, is_deleted
		FROM cards
		WHERE user_id = $1 AND card_pan = $2
		LIMIT 1
	`

	row, err := db.ExecuteDBQuery(ctx, conn.PgSql, conn.Timeout, q, userID, cardPAN)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}

	rec, err := scanCardRow(row)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// ListCards возвращает все записи пользователя.
func ListCards(ctx context.Context, conn *db.SqlConnection, userID int64) ([]*models.CardRecord, error) {
	const q = `
		SELECT user_id, card_pan, data, meta, created_at, updated_at, is_deleted
		FROM cards
		WHERE user_id = $1
		ORDER BY updated_at ASC
	`

	rows, err := db.ExecuteDBQueryAll(ctx, conn.PgSql, conn.Timeout, q, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*models.CardRecord, 0, len(rows))
	for _, r := range rows {
		rec, err := scanCardRow(r)
		if err != nil {
			return nil, err
		}
		result = append(result, rec)
	}
	return result, nil
}

// scanCardRow — преобразует []any в CardRecord.
func scanCardRow(row []any) (*models.CardRecord, error) {
	if len(row) != 7 {
		return nil, fmt.Errorf("unexpected column count: %d", len(row))
	}

	var (
		userIDRaw    = row[0]
		cardPANRaw   = row[1]
		dataRaw      = row[2]
		metaRaw      = row[3]
		createdRaw   = row[4]
		updatedRaw   = row[5]
		isDeletedRaw = row[6]
	)

	rec := &models.CardRecord{}

	// user_id (bigint)
	switch v := userIDRaw.(type) {
	case int64:
		rec.UserID = v
	default:
		return nil, fmt.Errorf("bad type for user_id: %T", userIDRaw)
	}

	// card_pan (text)
	if s, ok := cardPANRaw.(string); ok {
		rec.CardPAN = s
	} else {
		return nil, fmt.Errorf("bad type for card_pan: %T", cardPANRaw)
	}

	// data (text)
	if s, ok := dataRaw.(string); ok {
		rec.Data = s
	} else {
		return nil, fmt.Errorf("bad type for data: %T", dataRaw)
	}

	// meta (nullable text)
	if metaRaw == nil {
		rec.Meta = nil
	} else if s, ok := metaRaw.(string); ok {
		rec.Meta = &s
	} else {
		return nil, fmt.Errorf("bad type for meta: %T", metaRaw)
	}

	// created_at
	if t, ok := createdRaw.(time.Time); ok {
		rec.CreatedAt = t
	} else {
		return nil, fmt.Errorf("bad type for created_at: %T", createdRaw)
	}

	// updated_at
	if t, ok := updatedRaw.(time.Time); ok {
		rec.UpdatedAt = t
	} else {
		return nil, fmt.Errorf("bad type for updated_at: %T", updatedRaw)
	}

	// is_deleted
	if b, ok := isDeletedRaw.(bool); ok {
		rec.IsDeleted = b
	} else {
		return nil, fmt.Errorf("bad type for is_deleted: %T", isDeletedRaw)
	}

	return rec, nil
}
