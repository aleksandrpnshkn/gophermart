package models

import (
	"time"

	"github.com/aleksandrpnshkn/gophermart/internal/types"
)

type Order struct {
	OrderNumber string
	UserID      int64
	Status      types.OrderStatus
	Accrual     int64
	UploadedAt  time.Time
}
