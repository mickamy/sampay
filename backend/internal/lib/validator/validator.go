package validator

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	lib "github.com/go-playground/validator/v10"

	"github.com/mickamy/sampay/internal/lib/logger"
)

var validator = lib.New()

type ValidationErrorMessages = map[string][]string

func init() {
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Struct(ctx context.Context, s interface{}) ValidationErrorMessages {
	if err := validator.StructCtx(ctx, s); err != nil {
		errs := lib.ValidationErrors{}
		if errors.As(err, &errs) {
			return mapValidationErrors(ctx, errs)
		}
		panic(fmt.Errorf("unknown validation error: %w", err))
	}
	return nil
}

func mapValidationErrors(ctx context.Context, errs lib.ValidationErrors) ValidationErrorMessages {
	messages := map[string][]string{}
	for _, err := range errs {
		var message string
		switch err.ActualTag() {
		case "required":
			message = "is required."
		case "min":
			message = fmt.Sprintf("is too short. (min=%s)", err.Param())
		case "max":
			message = fmt.Sprintf("is too long. (max=%s)", err.Param())
		case "email":
			message = "is not a valid email."
		case "url":
			message = "is not a valid URL."
		default:
			logger.Warn(ctx, "failing back to default error message", "ActualTag", err.ActualTag(), "Tag", err.Tag(), "Param", err.Param())
			message = "is invalid."
		}

		field := err.Field()
		if val, ok := messages[field]; ok {
			messages[field] = append(val, message)
		} else {
			messages[field] = []string{message}
		}
	}
	return messages
}
