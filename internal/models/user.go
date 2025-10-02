package models

import "github.com/aleksandrpnshkn/gophermart/internal/types"

type User struct {
	ID    int64
	Login string
	Hash  types.PasswordHash
}
