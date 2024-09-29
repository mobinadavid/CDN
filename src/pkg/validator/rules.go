package validator

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"ticket/src/pkg/i18h"
)

func RegisterRules(val *validator.Validate, trans *ut.UniversalTranslator) {
	// Map of rule names to their corresponding validation functions
	ruleToFunc := map[string]validator.Func{
		"is-uuid": isValidUuid,
	}

	for ruleName, ruleFunc := range ruleToFunc {

		// Register the validation.
		_ = val.RegisterValidation(ruleName, ruleFunc)

		// Register validation messages as well.
		for _, lang := range i18n.Locales {
			translator, _ := trans.GetTranslator(lang)
			_ = val.RegisterTranslation(ruleName, translator, func(ut ut.Translator) error {
				return ut.Add(ruleName, i18n.Localize(lang, fmt.Sprintf("invalid-%s", ruleName)), false)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(fe.Tag(), fe.Field())
				return t
			})
		}
	}
}

// isValidUuid Custom validator function to validate UUID format
func isValidUuid(fl validator.FieldLevel) bool {
	_, err := uuid.Parse(fl.Field().String())
	return err == nil
}
