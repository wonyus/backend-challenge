package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Validator struct {
	emailRegex *regexp.Regexp
}

func New() *Validator {
	return &Validator{
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`),
	}
}

func (v *Validator) Validate(s interface{}) error {
	return v.validateStruct(reflect.ValueOf(s))
}

func (v *Validator) validateStruct(val reflect.Value) error {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		if err := v.validateField(field, fieldType.Name, tag); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateField(field reflect.Value, fieldName, tag string) error {
	rules := strings.Split(tag, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)

		switch {
		case rule == "required":
			if v.isEmpty(field) {
				return errors.New(fieldName + " is required")
			}
		case rule == "email":
			if field.Kind() == reflect.String && field.String() != "" {
				if !v.emailRegex.MatchString(strings.ToLower(field.String())) {
					return errors.New(fieldName + " must be a valid email")
				}
			}
		case strings.HasPrefix(rule, "min="):
			minStr := strings.TrimPrefix(rule, "min=")
			min, err := strconv.Atoi(minStr)
			if err != nil {
				continue
			}
			if field.Kind() == reflect.String && len(field.String()) < min {
				return errors.New(fieldName + " must be at least " + minStr + " characters")
			}
		case rule == "omitempty":
			if v.isEmpty(field) {
				return nil // Skip other validations if field is empty
			}
		}
	}

	return nil
}

func (v *Validator) isEmpty(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.String:
		return field.String() == ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return field.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return field.IsNil()
	default:
		return false
	}
}
