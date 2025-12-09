package repo

import (
	"context"
	"diploma/internal/models"
	"diploma/pkg/db"
	"fmt"
	"time"
)

// InsertUpdatePassword создаёт или обновляет запись пароля пользователя.
func InsertUpdatePassword(ctx context.Context, conn *db.SqlConnection, rec *models.PasswordRecord) error {
	if rec == nil {
		return fmt.Errorf("nil PasswordRecord")
	}

	// ON CONFLICT по составному ключу (user_id, id)
	const q = `
		INSERT INTO passwords (
			user_id, id, login, password, meta, created_at, updated_at, is_deleted
		) VALUES (
			$1, $2, $3, $4, $5, COALESCE($6, now()), now(), $7
		)
		ON CONFLICT (user_id, id) DO UPDATE SET
			login      = EXCLUDED.login,
			password   = EXCLUDED.password,
			meta       = EXCLUDED.meta,
			updated_at = now(),
			is_deleted = EXCLUDED.is_deleted
	`

	created := rec.CreatedAt
	if created.IsZero() {
		// по умолчанию created_date = now()
		returnWithCreatedNil := (func() error {
			_, err := db.ExecuteDBExec(ctx, conn.PgSql, conn.Timeout, q,
				rec.UserID,
				rec.ID,
				rec.Login,
				rec.Password,
				rec.Meta,
				nil,
				rec.IsDeleted,
			)
			return err
		})()
		return returnWithCreatedNil
	}

	_, err := db.ExecuteDBExec(ctx, conn.PgSql, conn.Timeout, q,
		rec.UserID,
		rec.ID,
		rec.Login,
		rec.Password,
		rec.Meta,
		rec.CreatedAt,
		rec.IsDeleted,
	)
	return err
}

// GetPassword возвращает одну запись пароля по user_id и идентификатору записи.
func GetPassword(ctx context.Context, conn *db.SqlConnection, userID int64, id string) (*models.PasswordRecord, error) {
	const q = `
		SELECT user_id, id, login, password, meta, created_at, updated_at, is_deleted
		FROM passwords
		WHERE user_id = $1 AND id = $2
		LIMIT 1
	`

	row, err := db.ExecuteDBQuery(ctx, conn.PgSql, conn.Timeout, q, userID, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}

	rec, err := scanPasswordRow(row)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// ListPasswords возвращает все записи паролей пользователя, отсортированные по времени обновления.
func ListPasswords(ctx context.Context, conn *db.SqlConnection, userID int64) ([]*models.PasswordRecord, error) {
	const q = `
		SELECT user_id, id, login, password, meta, created_at, updated_at, is_deleted
		FROM passwords
		WHERE user_id = $1
		ORDER BY updated_at ASC
	`

	rows, err := db.ExecuteDBQueryAll(ctx, conn.PgSql, conn.Timeout, q, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*models.PasswordRecord, 0, len(rows))
	for _, r := range rows {
		rec, err := scanPasswordRow(r)
		if err != nil {
			return nil, err
		}
		result = append(result, rec)
	}
	return result, nil
}

// scanPasswordRow конвертер из []any в PasswordRecord.
func scanPasswordRow(row []any) (*models.PasswordRecord, error) {
	if len(row) != 8 {
		return nil, fmt.Errorf("unexpected column count: %d", len(row))
	}

	var (
		userIDRaw    = row[0]
		idRaw        = row[1]
		loginRaw     = row[2]
		passwordRaw  = row[3]
		metaRaw      = row[4]
		createdRaw   = row[5]
		updatedRaw   = row[6]
		isDeletedRaw = row[7]
	)

	rec := &models.PasswordRecord{}

	// user_id (bigint)
	switch v := userIDRaw.(type) {
	case int64:
		rec.UserID = v
	default:
		return nil, fmt.Errorf("bad type for user_id: %T", userIDRaw)
	}

	// id (text)
	if s, ok := idRaw.(string); ok {
		rec.ID = s
	} else {
		return nil, fmt.Errorf("bad type for id: %T", idRaw)
	}

	// login (text)
	if s, ok := loginRaw.(string); ok {
		rec.Login = s
	} else {
		return nil, fmt.Errorf("bad type for login: %T", loginRaw)
	}

	// password (text)
	if s, ok := passwordRaw.(string); ok {
		rec.Password = s
	} else {
		return nil, fmt.Errorf("bad type for password: %T", passwordRaw)
	}

	// meta (nullable text)
	if metaRaw == nil {
		rec.Meta = nil
	} else if s, ok := metaRaw.(string); ok {
		rec.Meta = &s
	} else {
		return nil, fmt.Errorf("bad type for meta: %T", metaRaw)
	}

	// created_at (timestamptz)
	if t, ok := createdRaw.(time.Time); ok {
		rec.CreatedAt = t
	} else {
		return nil, fmt.Errorf("bad type for created_at: %T", createdRaw)
	}

	// updated_at (timestamptz)
	if t, ok := updatedRaw.(time.Time); ok {
		rec.UpdatedAt = t
	} else {
		return nil, fmt.Errorf("bad type for updated_at: %T", updatedRaw)
	}

	// is_deleted (bool)
	if b, ok := isDeletedRaw.(bool); ok {
		rec.IsDeleted = b
	} else {
		return nil, fmt.Errorf("bad type for is_deleted: %T", isDeletedRaw)
	}

	return rec, nil
}
