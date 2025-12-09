package services

import (
	"context"
	"diploma/internal/models"
	"diploma/internal/repo"
	"diploma/pkg/config"
	"diploma/pkg/crypto"
	"diploma/pkg/db"
	"errors"
	"fmt"
	"strings"
	"time"
)

// CardsService управляет созданием, хранением данных платежных карт пользователя.
type CardsService struct {
	cfg   *config.Config
	db    *db.SqlConnection
	users *UsersService
}

// NewCardsService создаёт новый сервис платежных карт.
func NewCardsService(cfg *config.Config, dbConn *db.SqlConnection, users *UsersService) *CardsService {
	return &CardsService{
		cfg:   cfg,
		db:    dbConn,
		users: users,
	}
}

// maskCardPAN уникальный id по номеру карты: первые 4 цифры + 8 * + последние 4 цифры.
func maskCardPAN(number string) (string, error) {
	num := strings.ReplaceAll(number, " ", "")
	num = strings.ReplaceAll(num, "-", "")

	if len(num) < 8 {
		return "", fmt.Errorf("card number too short")
	}
	return num[:4] + "********" + num[len(num)-4:], nil
}

// CreateOrUpdateCard создаёт или обновляет запись карты пользователя.
func (s *CardsService) CreateOrUpdateCard(
	ctx context.Context,
	login string,
	req *models.CardRequest,
) (*models.CardResponse, error) {
	userID, err := s.users.GetUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	cardPAN, err := maskCardPAN(req.Number)
	if err != nil {
		return nil, err
	}

	secret := models.CardSecret{
		Number: req.Number,
		Holder: req.Holder,
		Expire: req.Expire,
		CVV:    req.CVV,
	}
	encData, err := crypto.EncryptJSON(secret, s.cfg.DataEncKey)
	if err != nil {
		return nil, err
	}

	rec := &models.CardRecord{
		UserID:    userID,
		CardPAN:   cardPAN,
		Data:      encData,
		Meta:      req.Meta,
		CreatedAt: time.Time{},
		IsDeleted: false,
	}

	if err := repo.InsertUpdateCard(ctx, s.db, rec); err != nil {
		return nil, err
	}

	saved, err := repo.GetCard(ctx, s.db, userID, cardPAN)
	if err != nil {
		return nil, err
	}
	if saved == nil {
		return nil, errors.New("failed to fetch saved card record")
	}

	var secretOut models.CardSecret
	if err := crypto.DecryptJSON(saved.Data, s.cfg.DataEncKey, &secretOut); err != nil {
		return nil, ErrDecryptionFailed
	}

	return &models.CardResponse{
		CardPAN: saved.CardPAN,
		Number:  secretOut.Number,
		Holder:  secretOut.Holder,
		Expire:  secretOut.Expire,
		CVV:     secretOut.CVV,
		Meta:    saved.Meta,
	}, nil
}

// GetCard возвращает одну запись карты по логину пользователя и идентификатору записи.
func (s *CardsService) GetCard(
	ctx context.Context,
	login, cardPAN string,
) (*models.CardResponse, error) {
	userID, err := s.users.GetUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	rec, err := repo.GetCard(ctx, s.db, userID, cardPAN)
	if err != nil {
		return nil, err
	}
	if rec == nil || rec.IsDeleted {
		return nil, ErrCardNotFound
	}

	var secret models.CardSecret
	if err := crypto.DecryptJSON(rec.Data, s.cfg.DataEncKey, &secret); err != nil {
		return nil, ErrDecryptionFailed
	}

	return &models.CardResponse{
		CardPAN: rec.CardPAN,
		Number:  secret.Number,
		Holder:  secret.Holder,
		Expire:  secret.Expire,
		CVV:     secret.CVV,
		Meta:    rec.Meta,
	}, nil
}

// ListCards возвращает все платежные карты пользователя.
func (s *CardsService) ListCards(
	ctx context.Context,
	login string,
) ([]models.CardResponse, error) {
	userID, err := s.users.GetUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	recs, err := repo.ListCards(ctx, s.db, userID)
	if err != nil {
		return nil, err
	}

	resp := make([]models.CardResponse, 0, len(recs))
	for _, r := range recs {
		if r.IsDeleted {
			continue
		}
		var secret models.CardSecret
		if err := crypto.DecryptJSON(r.Data, s.cfg.DataEncKey, &secret); err != nil {
			return nil, ErrDecryptionFailed
		}
		resp = append(resp, models.CardResponse{
			CardPAN: r.CardPAN,
			Number:  secret.Number,
			Holder:  secret.Holder,
			Expire:  secret.Expire,
			CVV:     secret.CVV,
			Meta:    r.Meta,
		})
	}

	return resp, nil
}
