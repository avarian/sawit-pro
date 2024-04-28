package handler

import (
	"reflect"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
)

func TestServer_ValidateCreateUser(t *testing.T) {
	s := &Server{} // assuming Server struct is defined somewhere
	tests := []struct {
		name string
		data generated.RegisterJSONBody
		want generated.ErrorValidationResponse
	}{
		{
			name: "Valid Data",
			data: generated.RegisterJSONBody{
				FullName:    "John Doe",
				Password:    "P@ssw0rd",
				PhoneNumber: "+628123456789",
			},
			want: generated.ErrorValidationResponse{},
		},
		{
			name: "Empty Data",
			data: generated.RegisterJSONBody{},
			want: generated.ErrorValidationResponse{
				Messages: []string{"full_name : full_name is required", "password : password is required", "phone_number : phone_number is required"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.ValidateCreateUser(tt.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.ValidateCreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_ValidateLoginUser(t *testing.T) {
	s := &Server{}
	tests := []struct {
		name string
		data generated.LoginJSONBody
		want generated.ErrorValidationResponse
	}{
		{
			name: "Valid Login User",
			data: generated.LoginJSONBody{
				PhoneNumber: "+628536374829",
				Password:    "Password1.",
			},
			want: generated.ErrorValidationResponse{},
		},
		{
			name: "Empty Fields",
			data: generated.LoginJSONBody{},
			want: generated.ErrorValidationResponse{
				Messages: []string{"password : password is required", "phone_number : phone_number is required"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.ValidateLoginUser(tt.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.ValidateLoginUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_ValidateUpdateUser(t *testing.T) {
	s := &Server{}
	tests := []struct {
		name string
		data generated.User
		want generated.ErrorValidationResponse
	}{
		{
			name: "Valid Update User",
			data: generated.User{},
			want: generated.ErrorValidationResponse{},
		},
		{
			name: "Invalid Phone Number Format",
			data: generated.User{
				PhoneNumber: "1234567890",
			},
			want: generated.ErrorValidationResponse{
				Messages: []string{"phone_number : must be indonesian(+62) format"},
			},
		},
		{
			name: "Phone Number Not a Number",
			data: generated.User{
				PhoneNumber: "+62abcd343565",
			},
			want: generated.ErrorValidationResponse{
				Messages: []string{"phone_number : must be number"},
			},
		},
		{
			name: "Invalid FullName Length",
			data: generated.User{
				FullName: "Jo",
			},
			want: generated.ErrorValidationResponse{
				Messages: []string{"full_name : must be more than 3 and less than 60 characters long"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.ValidateUpdateUser(tt.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.ValidateUpdateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
