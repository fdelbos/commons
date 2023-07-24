package validation

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	once  sync.Once
	v     *validator.Validate
	trans ut.Translator
)

func Start() {
	once.Do(func() {

		v = validator.New()

		en := en.New()
		uni := ut.New(en, en)

		// this is usually know or extracted from http 'Accept-Language' header
		// also see uni.FindTranslator(...)
		trans, _ = uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(v, trans)

		// extract json name from tag
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

			if name == "-" {
				return ""
			}

			return name
		})
	})
}

func Validator() *validator.Validate {
	Start()

	return v
}

func Translator() ut.Translator {
	Start()
	return trans
}
