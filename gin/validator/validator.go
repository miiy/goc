package validator

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

type Options struct {
	Locale string
}

// Custom format and translation
type ValidationErrorsTranslations validator.ValidationErrorsTranslations

var (
	trans ut.Translator
)

func NewValidator(o *Options) (*validator.Validate, error) {
	var locale = localeToLang(o.Locale)

	var ok bool

	// validator
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errors.New("validator error")
	}

	// RegisterTagNameFunc registers a function to get alternate names for StructFields.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// translator
	enT := en.New()
	zhT := zh.New()
	// fallback(备用), supported ...
	uni := ut.New(enT, enT, zhT)
	// the specified translator for the given locale
	trans, ok = uni.GetTranslator(locale)
	if !ok {
		return nil, fmt.Errorf("uni.GetTranslator(%s) failed", locale)
	}

	var err error
	switch locale {
	case "en":
		err = enTranslations.RegisterDefaultTranslations(v, trans)
	case "zh":
		err = zhTranslations.RegisterDefaultTranslations(v, trans)
	default:
		err = enTranslations.RegisterDefaultTranslations(v, trans)
	}

	if err != nil {
		return nil, err
	}

	if err := translateOverride(trans, v); err != nil {
		return nil, err
	}

	return v, nil
}

func translateOverride(trans ut.Translator, v *validator.Validate) error {
	//  is_exists
	//  Example:
	// 	err = validate.RegisterValidation("is_exists", func(fl validator.FieldLevel) bool {
	//		user := find(fl.Field().String())
	//		return user == nil
	//	})
	err := v.RegisterTranslation("is_exists", trans, func(ut ut.Translator) error {
		var txt = "{0} already is exists"
		if trans.Locale() == "zh" {
			txt = "{0}已存在"
		}
		return ut.Add("is_exists", txt, true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("is_exists", fe.Field())

		return t
	})
	return err
}

// TODO local to language
func localeToLang(locale string) string {
	if locale == "zh-CN" {
		return "zh"
	}
	return "en"
}

// Custom format and translation
func VErrorsTranslations(e error) ValidationErrorsTranslations {
	if e == nil {
		return nil
	}

	var validationErrorsTranslations = make(ValidationErrorsTranslations)
	validationErrors := e.(validator.ValidationErrors)
	//fmt.Println(validationErrors.Translate(trans))
	for _, e := range validationErrors {
		field := e.Field()
		validationErrorsTranslations[field] = e.Translate(trans)
	}
	return validationErrorsTranslations
}
