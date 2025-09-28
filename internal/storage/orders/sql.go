package orders

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type SQLStorage struct {
	pgxpool *pgxpool.Pool
}

func (s *SQLStorage) Ping(ctx context.Context) error {
	return s.pgxpool.Ping(ctx)
}

func (s *SQLStorage) GetByNumber(
	ctx context.Context,
	orderNumber string,
) (models.Order, error) {
	var order models.Order

	row := s.pgxpool.QueryRow(ctx, `
        SELECT number, user_id, status FROM orders 
        WHERE number = $1
    `, orderNumber)
	err := row.Scan(&order.OrderNumber, &order.UserID, &order.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, ErrOrderNotFound
		}
		return models.Order{}, err
	}

	return order, nil
}

func (s *SQLStorage) Create(ctx context.Context, order models.Order) error {
	_, err := s.pgxpool.Exec(ctx, `
        INSERT INTO orders (number, user_id, status, accrual, uploaded_at) 
        VALUES (@number, @user_id, @status, @accrual, @uploaded_at)
    `, pgx.NamedArgs{
		"number":      order.OrderNumber,
		"user_id":     order.UserID,
		"status":      order.Status,
		"accrual":     order.Accrual,
		"uploaded_at": order.UploadedAt,
	})

	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code != pgerrcode.UniqueViolation {
		return err
	}

	existedOrder, err := s.GetByNumber(ctx, order.OrderNumber)
	if err != nil {
		return err
	}

	if existedOrder.UserID != order.UserID {
		return ErrOrderAlreadyCreatedByAnotherUser
	}

	return ErrOrderAlreadyCreated
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
