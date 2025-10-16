package services

import (
	"context"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type AppUni struct {
	uni *ut.UniversalTranslator
}

const (
	defaultTrans = "en"
)

func (u *AppUni) ResolveUserTrans(ctx context.Context) ut.Translator {
	trans, _ := u.uni.GetTranslator(defaultTrans) // можно будет брать из контекста
	return trans
}

func (u *AppUni) RegisterValidationTranslations(validate *validator.Validate) {
	defaultTrans, _ := u.uni.GetTranslator(defaultTrans)
	en_translations.RegisterDefaultTranslations(validate, defaultTrans)
}

func NewAppUni() *AppUni {
	en := en.New()
	uni := ut.New(en, en)

	return &AppUni{
		uni: uni,
	}
}
