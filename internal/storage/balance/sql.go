package balance

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shopspring/decimal"
)

type SQLStorage struct {
	pgxpool *pgxpool.Pool
}

const (
	ChangeBalanceQuery = `
        INSERT INTO balance_logs (id, order_number, user_id, amount, processed_at) 
        VALUES (DEFAULT, @order_number, @user_id, @amount, DEFAULT)
    `
)

var (
	ErrNotEnoughFunds = errors.New("not enough funds on user balance")
)

func (s *SQLStorage) Ping(ctx context.Context) error {
	return s.pgxpool.Ping(ctx)
}

func (s *SQLStorage) Withdraw(ctx context.Context, withdraw models.BalanceChange) error {
	tx, err := s.pgxpool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var balance decimal.Decimal

	row := tx.QueryRow(ctx, `
        SELECT COALESCE(SUM(amount), 0) 
        FROM (
            SELECT amount FROM balance_logs 
            WHERE user_id = $1
            FOR UPDATE
        )
    `, withdraw.UserID)
	err = row.Scan(&balance)
	if err != nil {
		return err
	}

	if balance.LessThan(withdraw.Amount) {
		return ErrNotEnoughFunds
	}

	_, err = tx.Exec(ctx, ChangeBalanceQuery, pgx.NamedArgs{
		"order_number": withdraw.OrderNumber,
		"user_id":      withdraw.UserID,
		"amount":       withdraw.Amount,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *SQLStorage) GetWithdrawals(
	ctx context.Context,
	user models.User,
) ([]models.BalanceChange, error) {
	withdrawals := []models.BalanceChange{}

	rows, err := s.pgxpool.Query(ctx, `
        SELECT order_number, user_id, amount, processed_at 
        FROM balance_logs 
        WHERE user_id = $1 AND amount < 0
        ORDER BY processed_at DESC
    `, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var balanceChange models.BalanceChange

		err = rows.Scan(
			&balanceChange.OrderNumber,
			&balanceChange.UserID,
			&balanceChange.Amount,
			&balanceChange.ProcessedAt,
		)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, balanceChange)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (s *SQLStorage) GetBalance(
	ctx context.Context,
	user models.User,
) (models.Balance, error) {
	var balance models.Balance

	row := s.pgxpool.QueryRow(ctx, `
        SELECT 
            COALESCE(SUM(amount), 0) AS current,
            COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0) AS withdrawn
        FROM balance_logs 
        WHERE user_id = $1
    `, user.ID)
	err := row.Scan(&balance.Current, &balance.Withdrawn)

	return balance, err
}

func (s *SQLStorage) Close() error {
	s.pgxpool.Close()
	return nil
}

func NewSQLStorage(ctx context.Context, databaseDSN string) (*SQLStorage, error) {
	pool, err := pgxpool.New(ctx, databaseDSN)
	if err != nil {
		return nil, err
	}

	storage := SQLStorage{
		pgxpool: pool,
	}

	err = storage.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &storage, nil
}
