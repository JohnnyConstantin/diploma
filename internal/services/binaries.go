package services

import (
	"context"
	"diploma/internal/models"
	"diploma/internal/repo"
	"diploma/pkg/config"
	"diploma/pkg/crypto"
	"diploma/pkg/db"
	"errors"
	"time"
)

// BinariesService управляет созданием, хранением бинарных данных.
type BinariesService struct {
	Cfg   *config.Config
	DB    *db.SqlConnection
	Users *UsersService
}

// NewBinariesService создаёт новый сервис бинарных данных.
func NewBinariesService(cfg *config.Config, dbConn *db.SqlConnection, users *UsersService) *BinariesService {
	return &BinariesService{
		Cfg:   cfg,
		DB:    dbConn,
		Users: users,
	}
}

// CreateOrUpdateBinary создаёт или обновляет запись с бинарными данными.
func (s *BinariesService) CreateOrUpdateBinary(
	ctx context.Context,
	login string,
	req *models.BinaryRequest,
) (*models.BinaryResponse, error) {
	userID, err := s.Users.GetUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	secret := models.BinarySecret{Data: req.Data}
	encData, err := crypto.EncryptJSON(secret, s.Cfg.DataEncKey)
	if err != nil {
		return nil, err
	}

	rec := &models.BinaryRecord{
		UserID:    userID,
		ID:        req.ID,
		Data:      encData,
		Meta:      req.Meta,
		CreatedAt: time.Time{},
		IsDeleted: false,
	}

	if err := repo.UpsertBinary(ctx, s.DB, rec); err != nil {
		return nil, err
	}

	saved, err := repo.GetBinary(ctx, s.DB, userID, req.ID)
	if err != nil {
		return nil, err
	}
	if saved == nil {
		return nil, errors.New("failed to fetch saved binary record")
	}

	// Как и в Passwords/Text CreateOrUpdate — отдаём зашифрованное значение.
	return &models.BinaryResponse{
		ID:        saved.ID,
		Data:      saved.Data,
		Meta:      saved.Meta,
		CreatedAt: saved.CreatedAt,
		UpdatedAt: saved.UpdatedAt,
		IsDeleted: saved.IsDeleted,
	}, nil
}

// GetBinary возвращает одну запись бинарных данных по логину пользователя и идентификатору записи.
func (s *BinariesService) GetBinary(
	ctx context.Context,
	login, id string,
) (*models.BinaryResponse, error) {
	userID, err := s.Users.GetUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	rec, err := repo.GetBinary(ctx, s.DB, userID, id)
	if err != nil {
		return nil, err
	}
	if rec == nil || rec.IsDeleted {
		return nil, ErrBinaryNotFound
	}

	var secret models.BinarySecret
	if err := crypto.DecryptJSON(rec.Data, s.Cfg.DataEncKey, &secret); err != nil {
		return nil, ErrDecryptionFailed
	}

	return &models.BinaryResponse{
		ID:        rec.ID,
		Data:      secret.Data, // уже расшифрованное содержимое
		Meta:      rec.Meta,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
		IsDeleted: rec.IsDeleted,
	}, nil
}

// ListBinaries возвращает все бинарные данные пользователя.
func (s *BinariesService) ListBinaries(
	ctx context.Context,
	login string,
) ([]models.BinaryResponse, error) {
	userID, err := s.Users.GetUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	recs, err := repo.ListBinaries(ctx, s.DB, userID)
	if err != nil {
		return nil, err
	}

	resp := make([]models.BinaryResponse, 0, len(recs))
	for _, r := range recs {
		if r.IsDeleted {
			continue
		}
		var secret models.BinarySecret
		if err := crypto.DecryptJSON(r.Data, s.Cfg.DataEncKey, &secret); err != nil {
			return nil, ErrDecryptionFailed
		}
		resp = append(resp, models.BinaryResponse{
			ID:        r.ID,
			Data:      secret.Data,
			Meta:      r.Meta,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			IsDeleted: r.IsDeleted,
		})
	}

	return resp, nil
}
