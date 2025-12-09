package models

import "time"

type CardSecret struct {
	Number string `json:"number"`
	Holder string `json:"holder"`
	Expire string `json:"expire"`
	CVV    string `json:"cvv"`
}

// CardRecord модель для таблицы cards.
type CardRecord struct {
	UserID    int64
	CardPAN   string
	Data      string
	Meta      *string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
}

type BinarySecret struct {
	// хранится в виде base64-строки, уже внутри зашифрованного JSON.
	Data string `json:"data"`
}

// BinaryRecord модель для таблицы binaries.
type BinaryRecord struct {
	UserID    int64
	ID        string
	Data      string
	Meta      *string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
}

// PasswordSecret представляет собой структуру, которая шифруется и хранит пароль.
type PasswordSecret struct {
	Password string `json:"password"`
}

// PasswordRecord описывает строку таблицы passwords в базе данных.
type PasswordRecord struct {
	UserID    int64     // FK на logins.id
	ID        string    // ID записи (то, что задаёт клиент)
	Login     string    // логин для этого ресурса
	Password  string    // сам пароль (шифрованный)
	Meta      *string   // произвольное описание / мета
	CreatedAt time.Time // когда запись создана
	UpdatedAt time.Time // когда запись обновлена
	IsDeleted bool      // признак, удален или нет
}

type TextSecret struct {
	Text string `json:"text"`
}

// TextRecord модель для таблицы texts.
type TextRecord struct {
	UserID    int64
	ID        string
	Text      string
	Meta      *string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
}
