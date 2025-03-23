package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

func (c *Config) Validate() error {
	validate := validator.New()

	validate.RegisterValidation("file_strategy", validateFileStrategy)
	validate.RegisterValidation("ignore_strategy", validateIgnoreStrategy)
	validate.RegisterValidation("input_file", validateInputFile)
	validate.RegisterValidation("file_name_template", validateFileNameTemplate)

	if err := validate.Struct(c); err != nil {
		return wrapValidationErrors(err)
	}

	if c.Generation.Count < 1 {
		return fmt.Errorf("count must be at least 1")
	}

	return nil
}

func validateFileStrategy(fl validator.FieldLevel) bool {
	fmt.Println("validateFileStrategy")
	value, ok := fl.Field().Interface().(OutputStrategy)
	fmt.Println(value)
	if !ok {
		return false
	}
	return value >= FilePerStruct && value <= SingleFile
}

func validateIgnoreStrategy(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(FieldIgnoreStrategy)
	if !ok {
		return false
	}
	return value >= IgnoreUntagged && value <= IncludeAll
}

func validateInputFile(fl validator.FieldLevel) bool {
	filePath := fl.Field().String()
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func validateFileNameTemplate(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	return strings.Contains(value, "{struct}")
}

func wrapValidationErrors(err error) error {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	var errMsgs []string
	for _, e := range err.(validator.ValidationErrors) {
		switch e.Tag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", e.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be at least %s", e.Field(), e.Param()))
		case "oneof":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be one of %s", e.Field(), e.Param()))
		case "file_strategy":
			errMsgs = append(errMsgs, fmt.Sprintf("invalid file strategy in field %s", e.Field()))
		case "ignore_strategy":
			errMsgs = append(errMsgs, fmt.Sprintf("invalid ignore strategy in field %s", e.Field()))
		case "input_file":
			errMsgs = append(errMsgs, fmt.Sprintf("invalid input file path: %s", e.Value()))
		case "file_name_template":
			errMsgs = append(errMsgs, "file name template should contain {struct} placeholder")
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("validation failed for field %s", e.Field()))
		}
	}

	return fmt.Errorf("configuration errors:\n  • %s", strings.Join(errMsgs, "\n  • "))
}
