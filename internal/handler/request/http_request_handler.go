package request

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/tanjed/go-sso/internal/customerror"
)

type HttpRequest interface{
	Validated(d io.Reader) error
}


func GetValidated[T HttpRequest](d io.Reader, s T) (error) {
	err := json.NewDecoder(d).Decode(&s); 
	
	if err != nil {
		slog.Error("Unable to process the request", "error", err)
	}

	validate := validator.New()
	err = validate.Struct(s)

	validationErrors := make(map[string][]string)
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range errors {
				validationErrors[fieldError.Field()] = append(
					validationErrors[fieldError.Field()], 
					fmt.Sprintf("'%s' failed validation: %s", fieldError.Field(), fieldError.Tag())) 
			}
			
		} else {
			slog.Error("Unable to process the request", "error", err)
			validationErrors["default"] = append(validationErrors["default"], "unknown validation error")
		}

		return &customerror.ValidationError{
			ErrMessage: "validation error occured",
			ErrCode: http.StatusUnprocessableEntity,
			ErrBag: validationErrors,
		}
    }

	return nil
}
