package validator

import (
	"fmt"
	"motico-api/internal/rest/response"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateRequest(r *http.Request, v interface{}) error {
	if err := validate.Struct(v); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()
			validationErrors[field] = fmt.Sprintf("Field '%s' failed validation: %s", field, tag)
		}
		return fmt.Errorf("validation failed: %v", validationErrors)
	}
	return nil
}

func HandleValidationError(w http.ResponseWriter, err error) {
	response.Error(w, http.StatusBadRequest, "Validation error", err.Error())
}
