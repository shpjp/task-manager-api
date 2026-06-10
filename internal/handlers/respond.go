package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// init makes validation errors report json tag names (e.g. "due_date")
// instead of Go struct field names.
func init() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

type apiError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"error": apiError{Code: code, Message: message}})
}

func respondData(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}

// respondBindingError converts gin binding/validation failures into a
// consistent 400 payload with per-field messages.
func respondBindingError(c *gin.Context, err error) {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		fields := make(map[string]string, len(verrs))
		for _, fe := range verrs {
			fields[fe.Field()] = fieldMessage(fe)
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": apiError{
			Code:    "VALIDATION_ERROR",
			Message: "One or more fields are invalid",
			Fields:  fields,
		}})
		return
	}
	respondError(c, http.StatusBadRequest, "INVALID_BODY", "Request body is missing or malformed")
}

func fieldMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		if fe.Kind().String() == "string" {
			return fmt.Sprintf("Must be at least %s characters", fe.Param())
		}
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		if fe.Kind().String() == "string" {
			return fmt.Sprintf("Must be at most %s characters", fe.Param())
		}
		return fmt.Sprintf("Must be at most %s", fe.Param())
	case "oneof":
		return "Must be one of: " + strings.Join(strings.Fields(fe.Param()), ", ")
	default:
		return "Invalid value"
	}
}
