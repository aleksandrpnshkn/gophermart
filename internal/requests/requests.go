package requests

type (
	Login struct {
		Login    string `json:"login" validate:"required,alphanum,min=3,max=30"`
		Password string `json:"password" validate:"required,alphanum,min=6,max=50"`
	}

	Register struct {
		Login    string `json:"login" validate:"required,alphanum,min=3,max=30"`
		Password string `json:"password" validate:"required,alphanum,min=6,max=50"`
	}

	Withdraw struct {
		OrderNumber string  `json:"order" validate:"required,numeric,min=3,max=100,luhn"`
		Amount      float64 `json:"sum" validate:"required,number,min=1"`
	}
)
