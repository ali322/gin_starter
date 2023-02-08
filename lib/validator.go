package lib

import (
	"log"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

func TranslateValidatorErrors(err validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)
	for f, err := range err.Translate(translator) {
		stripedFieldName := f[strings.Index(f, ".")+1:]
		errs[stripedFieldName] = err
	}
	return errs
}

// func registerValidatorTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
// 	return func(trans Translator) error {
// 		if err := translator.Add(tag, msg, false); err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// }

// func translateFieldError(trans Translator, fieldError validator.FieldError) string {
// 	msg, err := trans.T(fieldError.Tag(), fieldError.Field())
// 	if err != nil {
// 		return ""
// 	}
// 	return msg
// }

// func checkDate(fl validator.FieldLevel) bool {
// 	date, err := time.Parse("2006-01-02", fl.Field().String())
// 	if err != nil {
// 		return false
// 	}
// 	if date.Before(time.Now()) {
// 		return false
// 	}
// 	return true
// }

func RegisterValidatorTranslations(locale string) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		var err error
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, translator)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, translator)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, translator)
		}
		if err != nil {
			log.Fatal(err)
		}
		// if err := v.RegisterValidation("checkDate", checkDate); err != nil {
		// 	log.Fatal(err)
		// }
		// if err := v.RegisterTranslation(
		// 	"checkDate",
		// 	translator,
		// 	registerValidatorTranslator("checkDate", "{0}必须要晚于当前日期"),
		// 	translateFieldError,
		// ); err != nil {
		// 	log.Fatal(err)
		// }
	}
}
