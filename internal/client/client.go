package client

import "time"

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type PasswordRequest struct {
	ID       string  `json:"id"`
	Login    string  `json:"login"`
	Password string  `json:"password"`
	Meta     *string `json:"meta,omitempty"`
}

type PasswordResponse struct {
	ID        string    `json:"id"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	Meta      *string   `json:"meta,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}

type TextRequest struct {
	ID   string  `json:"id"`
	Text string  `json:"text"`
	Meta *string `json:"meta,omitempty"`
}

type TextResponse struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Meta      *string   `json:"meta,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}

type CardRequest struct {
	Number string  `json:"number"`
	Holder string  `json:"holder"`
	Expire string  `json:"expire"`
	CVV    string  `json:"cvv"`
	Meta   *string `json:"meta,omitempty"`
}

type CardResponse struct {
	CardPAN string  `json:"card_pan"`
	Number  string  `json:"number"`
	Holder  string  `json:"holder"`
	Expire  string  `json:"expire"`
	CVV     string  `json:"cvv"`
	Meta    *string `json:"meta,omitempty"`
}

type BinaryRequest struct {
	ID   string  `json:"id"`
	Data string  `json:"data"`
	Meta *string `json:"meta,omitempty"`
}

type BinaryResponse struct {
	ID        string    `json:"id"`
	Data      string    `json:"data"`
	Meta      *string   `json:"meta,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}
