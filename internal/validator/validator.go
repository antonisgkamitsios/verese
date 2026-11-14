package validator

import (
	"net/mail"
	"slices"
	"strings"
)

type Validator struct {
	Errors         Error    `db:"-"`
	NonFieldErrors []string `db:"-"`
}

type Error map[string]string

func New() *Validator {
	return &Validator{Errors: make(Error)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if v.Errors == nil {
		v.Errors = make(Error)
	}

	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func FieldsEqual(value1, value2 string) bool {
	return value1 == value2
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

func ValidMail(value string) bool {
	_, err := mail.ParseAddress(value)
	return err == nil
}
