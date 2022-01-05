package validators

import (
	validator "github.com/go-playground/validator/v10"
)

type (
	OrderValidatorInterface interface {
		Validate(payload interface{}) error
	}

	OrderValidator struct {
		Engine *validator.Validate
	}

	OrderValidationErrorInterface interface {
		Error() string
	}

	OrderValidationError struct {
		Message string                   `json:"message"`
		Details []map[string]interface{} `json:"details"`
	}
)

func (err *OrderValidationError) Error() string {
	return err.Message
}

func NewValidator() OrderValidatorInterface {
	engine := validator.New()

	return OrderValidator{
		Engine: engine,
	}
}

func (v OrderValidator) Validate(payload interface{}) error {
	err := v.Engine.Struct(payload)

	if err != nil {
		validationError := &OrderValidationError{
			Message: "Invalid input for order creation",
		}

		for _, err := range err.(validator.ValidationErrors) {
			validationError.Details = append(validationError.Details, map[string]interface{}{
				"field":  err.Field(),
				"value":  err.Value(),
				"reason": err.Tag(),
			})

		}

		return validationError
	}

	return nil
}
