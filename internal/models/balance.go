package models

import (
	"github.com/shopspring/decimal"
)

type Balance struct {
	Current   decimal.Decimal
	Withdrawn decimal.Decimal
}
