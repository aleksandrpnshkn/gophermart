package models

import (
	"time"

	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"github.com/shopspring/decimal"
)

type Order struct {
	OrderNumber string
	UserID      int64
	Status      types.OrderStatus
	Accrual     decimal.Decimal
	UploadedAt  time.Time
}
