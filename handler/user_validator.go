package handler

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
)

func (s *Server) ValidateCreateUser(request interface{}) (errors generated.ErrorValidationResponse) {
	v := reflect.ValueOf(request)
	validations := make(map[string][]string)

	for i := 0; i < v.NumField(); i++ {
		var min, max int
		tag := v.Type().Field(i).Tag.Get("json") 
		val := v.Field(i).Interface()

		const (
			FullNameTag    = "full_name"
			PasswordTag    = "password"
			PhoneNumberTag = "phone_number"
		)

		switch tag {
		case FullNameTag:
			min, max = 3, 60
		case PasswordTag:
			min, max = 6, 64
			password := v.Field(i).String()
			specialRegexp := regexp.MustCompile(`[!@#$%^&*()_+\[\]{};':"\|,.<>?]`)
			capitialRegexp := regexp.MustCompile(`[A-Z]`)
			numericalRegexp := regexp.MustCompile(`[0-9]`)
			if !specialRegexp.MatchString(password) {
				validations[tag] = append(validations[tag], "at least 1 special character needed")
			}
			if !capitialRegexp.MatchString(password) {
				validations[tag] = append(validations[tag], "at least 1 capital character needed")
			}
			if !numericalRegexp.MatchString(password) {
				validations[tag] = append(validations[tag], "at least 1 number needed")
			}
		case PhoneNumberTag:
			min, max = 10, 13
			s := strings.Split(v.Field(i).String(), "+62")
			if s[0] != "" {
				validations[tag] = append(validations[tag], "must be indonesian(+62) format")
			}
			if _, err := strconv.Atoi(s[1]);err != nil {
        validations[tag] = append(validations[tag], "must be number")
    	}
		}

		if val == "" {
			validations[tag] = append([]string{}, fmt.Sprintf("%s is required", tag))
		} else if (max != 0 && min != 0) && v.Field(i).Len() < min || v.Field(i).Len() > max { 
			validations[tag] = append(validations[tag], fmt.Sprintf("must be more than %d and less than %d characters long", min, max))
		}
	}

	if len(validations) != 0 {
		for field, validation := range validations {
			errors.Messages = append(errors.Messages, fmt.Sprintf("%s : %s", field, strings.Join(validation, ", ")))
		}
		return errors
	}
	return
}

func (s *Server) ValidateLoginUser(request interface{}) (errors generated.ErrorValidationResponse) {
	v := reflect.ValueOf(request)
	validations := make(map[string][]string)

	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get("json") 
		val := v.Field(i).Interface()

		if val == "" {
			validations[tag] = append([]string{}, fmt.Sprintf("%s is required", tag))
		}
	}

	if len(validations) != 0 {
		for field, validation := range validations {
			errors.Messages = append(errors.Messages, fmt.Sprintf("%s : %s", field, strings.Join(validation, ", ")))
		}
		return errors
	}
	return
}

func (s *Server) ValidateUpdateUser(request interface{}) (errors generated.ErrorValidationResponse) {
	v := reflect.ValueOf(request)
	validations := make(map[string][]string)

	for i := 0; i < v.NumField(); i++ {
		var min, max int
		tag := v.Type().Field(i).Tag.Get("json") 
		val := v.Field(i).Interface()

		const (
			FullNameTag    = "full_name"
			PhoneNumberTag = "phone_number"
		)

		switch tag {
		case FullNameTag:
			min, max = 3, 60
		case PhoneNumberTag:
			min, max = 10, 13
			if val != "" {
				s := strings.Split(v.Field(i).String(), "+62")
				if s[0] != "" {
					validations[tag] = append(validations[tag], "must be indonesian(+62) format")
				}
				if _, err := strconv.Atoi(s[1]);err != nil {
					validations[tag] = append(validations[tag], "must be number")
				}
			}
		}

		if val != "" && (max != 0 && min != 0) && v.Field(i).Len() < min || v.Field(i).Len() > max { 
			validations[tag] = append(validations[tag], fmt.Sprintf("must be more than %d and less than %d characters long", min, max))
		}
	}

	if len(validations) != 0 {
		for field, validation := range validations {
			errors.Messages = append(errors.Messages, fmt.Sprintf("%s : %s", field, strings.Join(validation, ", ")))
		}
		return errors
	}
	return
}
