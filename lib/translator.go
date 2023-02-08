package lib

import (
	"fmt"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
)

type Translator = ut.Translator

var translator Translator

func InitTranslator(locale string) (ut.Translator, error) {
	if translator != nil {
		return translator, nil
	}
	zhTranslator := zh.New()
	enTranslator := en.New()
	utTranslator := ut.New(enTranslator, zhTranslator, enTranslator)
	trans, ok := utTranslator.GetTranslator(locale)
	if !ok {
		return nil, fmt.Errorf("failed to get translator of %s", locale)
	}
	translator = trans
	return trans, nil
}
