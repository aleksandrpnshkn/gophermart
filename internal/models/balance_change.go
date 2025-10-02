package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type BalanceChange struct {
	OrderNumber string
	UserID      int64
	Amount      decimal.Decimal
	ProcessedAt time.Time
}
