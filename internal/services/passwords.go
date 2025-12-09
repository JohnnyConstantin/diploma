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

// PasswordsService управляет созданием, хранением и выдачей паролей пользователя.
type PasswordsService struct {
	cfg   *config.Config
	db    *db.SqlConnection
	users *UsersService
}

// NewPasswordsService создаёт новый сервис паролей.
func NewPasswordsService(cfg *config.Config, dbConn *db.SqlConnection, users *UsersService) *PasswordsService {
	return &PasswordsService{
		cfg:   cfg,
		db:    dbConn,
		users: users,
	}
}

// CreateOrUpdatePassword создаёт или обновляет запись пароля для указанного логина.
func (s *PasswordsService) CreateOrUpdatePassword(
	ctx context.Context,
	login string,
	req *models.PasswordRequest,
) (*models.PasswordResponse, error) {
	userID, err := s.users.GetUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	secret := models.PasswordSecret{Password: req.Password}
	encPassword, err := crypto.EncryptJSON(secret, s.cfg.DataEncKey)
	if err != nil {
		return nil, err
	}

	rec := &models.PasswordRecord{
		UserID:    userID,
		ID:        req.ID,
		Login:     req.Login,
		Password:  encPassword,
		Meta:      req.Meta,
		CreatedAt: time.Time{},
		IsDeleted: false,
	}

	if err := repo.UpsertPassword(ctx, s.db, rec); err != nil {
		return nil, err
	}

	saved, err := repo.GetPassword(ctx, s.db, userID, req.ID)
	if err != nil {
		return nil, err
	}
	if saved == nil {
		return nil, errors.New("failed to fetch saved record")
	}

	return &models.PasswordResponse{
		ID:        saved.ID,
		Login:     saved.Login,
		Password:  saved.Password,
		Meta:      saved.Meta,
		CreatedAt: saved.CreatedAt,
		UpdatedAt: saved.UpdatedAt,
		IsDeleted: saved.IsDeleted,
	}, nil
}

// GetPassword возвращает одну запись пароля по логину пользователя и идентификатору записи.
func (s *PasswordsService) GetPassword(
	ctx context.Context,
	login, id string,
) (*models.PasswordResponse, error) {
	userID, err := s.users.GetUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	rec, err := repo.GetPassword(ctx, s.db, userID, id)
	if err != nil {
		return nil, err
	}
	if rec == nil || rec.IsDeleted {
		return nil, ErrPasswordNotFound
	}

	var secret models.PasswordSecret
	if err := crypto.DecryptJSON(rec.Password, s.cfg.DataEncKey, &secret); err != nil {
		return nil, ErrDecryptionFailed
	}

	return &models.PasswordResponse{
		ID:        rec.ID,
		Login:     rec.Login,
		Password:  secret.Password,
		Meta:      rec.Meta,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
		IsDeleted: rec.IsDeleted,
	}, nil
}

// ListPasswords возвращает все пароли пользователя.
func (s *PasswordsService) ListPasswords(
	ctx context.Context,
	login string,
) ([]models.PasswordResponse, error) {
	userID, err := s.users.GetUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	recs, err := repo.ListPasswords(ctx, s.db, userID)
	if err != nil {
		return nil, err
	}

	resp := make([]models.PasswordResponse, 0, len(recs))
	for _, r := range recs {
		var secret models.PasswordSecret
		if err := crypto.DecryptJSON(r.Password, s.cfg.DataEncKey, &secret); err != nil {
			return nil, ErrDecryptionFailed
		}
		resp = append(resp, models.PasswordResponse{
			ID:        r.ID,
			Login:     r.Login,
			Password:  secret.Password,
			Meta:      r.Meta,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			IsDeleted: r.IsDeleted,
		})
	}

	return resp, nil
}
