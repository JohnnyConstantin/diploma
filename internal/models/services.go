package models

import (
	"time"
)

// BinaryRequest конструкция входящего запроса
type BinaryRequest struct {
	ID   string  `json:"id" binding:"required"`   // ID записи (задаётся клиентом)
	Data string  `json:"data" binding:"required"` // произвольные бинарные данные (например, base64)
	Meta *string `json:"meta"`                    // произвольная мета
}

// BinaryResponse конструкция исходящего ответа
type BinaryResponse struct {
	ID        string    `json:"id"`
	Data      string    `json:"data"`
	Meta      *string   `json:"meta,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}

// CardRequest конструкция входящего запроса
type CardRequest struct {
	Number string  `json:"number" binding:"required"` // полный номер карты
	Holder string  `json:"holder" binding:"required"` // имя владельца
	Expire string  `json:"expire" binding:"required"` // срок действия (например, 12/30)
	CVV    string  `json:"cvv" binding:"required"`    // CVV
	Meta   *string `json:"meta"`                      // произвольная мета
}

// CardResponse конструкция исходящего ответа
type CardResponse struct {
	CardPAN string  `json:"card_pan"` // маскированный PAN, например 4600********5363
	Number  string  `json:"number"`   // полный номер (для владельца, уже расшифрованный)
	Holder  string  `json:"holder"`
	Expire  string  `json:"expire"`
	CVV     string  `json:"cvv"`
	Meta    *string `json:"meta,omitempty"`
}

// PasswordRequest конструкция входящего запроса
type PasswordRequest struct {
	ID       string  `json:"id" binding:"required"`
	Login    string  `json:"login" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Meta     *string `json:"meta"`
}

// PasswordResponse конструкция исходящего ответа
type PasswordResponse struct {
	ID        string    `json:"id"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	Meta      *string   `json:"meta,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}

// TextRequest конструкция входящего запроса
type TextRequest struct {
	ID   string  `json:"id" binding:"required"`   // заголовок/ID записи
	Text string  `json:"text" binding:"required"` // сам текст
	Meta *string `json:"meta"`                    // произвольная мета
}

// TextResponse конструкция исходящего ответа
type TextResponse struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Meta      *string   `json:"meta,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}
