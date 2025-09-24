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
)
