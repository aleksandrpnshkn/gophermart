package types

type OrderStatus string

const (
	// заказ зарегистрирован, но вознаграждение не рассчитано
	OrderStatusNew OrderStatus = "NEW"

	// вознаграждение за заказ рассчитывается
	OrderStatusProcessing OrderStatus = "PROCESSING"

	// система расчёта вознаграждений отказала в расчёте
	OrderStatusInvalid OrderStatus = "INVALID"

	// данные по заказу проверены и информация о расчёте успешно получена
	OrderStatusProcessed OrderStatus = "PROCESSED"
)
