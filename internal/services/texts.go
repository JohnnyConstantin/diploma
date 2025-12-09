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

// TextsService управляет созданием, хранением текста.
type TextsService struct {
	cfg   *config.Config
	db    *db.SqlConnection
	users *UsersService
}

// NewTextsService создаёт новый сервис текстов.
func NewTextsService(cfg *config.Config, dbConn *db.SqlConnection, users *UsersService) *TextsService {
	return &TextsService{
		cfg:   cfg,
		db:    dbConn,
		users: users,
	}
}

// CreateOrUpdateText создаёт или обновляет запись текста для указанного логина.
func (s *TextsService) CreateOrUpdateText(
	ctx context.Context,
	login string,
	req *models.TextRequest,
) (*models.TextResponse, error) {
	userID, err := s.users.GetUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	secret := models.TextSecret{Text: req.Text}
	encText, err := crypto.EncryptJSON(secret, s.cfg.DataEncKey)
	if err != nil {
		return nil, err
	}

	rec := &models.TextRecord{
		UserID:    userID,
		ID:        req.ID,
		Text:      encText,
		Meta:      req.Meta,
		CreatedAt: time.Time{},
		IsDeleted: false,
	}

	if err := repo.UpsertText(ctx, s.db, rec); err != nil {
		return nil, err
	}

	saved, err := repo.GetText(ctx, s.db, userID, req.ID)
	if err != nil {
		return nil, err
	}
	if saved == nil {
		return nil, errors.New("failed to fetch saved text record")
	}

	// Зашифрованное значение
	return &models.TextResponse{
		ID:        saved.ID,
		Text:      saved.Text,
		Meta:      saved.Meta,
		CreatedAt: saved.CreatedAt,
		UpdatedAt: saved.UpdatedAt,
		IsDeleted: saved.IsDeleted,
	}, nil
}

// GetText возвращает одну запись текста по логину пользователя и идентификатору записи.
func (s *TextsService) GetText(
	ctx context.Context,
	login, id string,
) (*models.TextResponse, error) {
	userID, err := s.users.GetUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	rec, err := repo.GetText(ctx, s.db, userID, id)
	if err != nil {
		return nil, err
	}
	if rec == nil || rec.IsDeleted {
		return nil, ErrTextNotFound
	}

	var secret models.TextSecret
	if err := crypto.DecryptJSON(rec.Text, s.cfg.DataEncKey, &secret); err != nil {
		return nil, ErrDecryptionFailed
	}

	return &models.TextResponse{
		ID:        rec.ID,
		Text:      secret.Text,
		Meta:      rec.Meta,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
		IsDeleted: rec.IsDeleted,
	}, nil
}

// ListTexts возвращает все текста пользователя.
func (s *TextsService) ListTexts(
	ctx context.Context,
	login string,
) ([]models.TextResponse, error) {
	userID, err := s.users.GetUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	recs, err := repo.ListTexts(ctx, s.db, userID)
	if err != nil {
		return nil, err
	}

	resp := make([]models.TextResponse, 0, len(recs))
	for _, r := range recs {
		var secret models.TextSecret
		if err := crypto.DecryptJSON(r.Text, s.cfg.DataEncKey, &secret); err != nil {
			return nil, ErrDecryptionFailed
		}
		resp = append(resp, models.TextResponse{
			ID:        r.ID,
			Text:      secret.Text,
			Meta:      r.Meta,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			IsDeleted: r.IsDeleted,
		})
	}

	return resp, nil
}
