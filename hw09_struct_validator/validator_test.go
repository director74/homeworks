package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "389582290458082048242804112310583152",
				Name:   "Tod",
				Age:    49,
				Email:  "demo@localhost.com",
				Role:   UserRole("stuff"),
				Phones: []string{"89999999999", "89999992999"},
				meta:   nil,
			},
			nil,
		},
		{
			App{
				Version: "49210",
			},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "",
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// Place your code here.
			_ = tt
		})
	}
}
