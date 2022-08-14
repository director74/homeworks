package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
		in interface{}
		is error
		as error
	}{
		{
			in: User{
				ID:     "389582290458082048242804112310583152",
				Name:   "Tod",
				Age:    49,
				Email:  "demo@localhost.com",
				Role:   UserRole("stuff"),
				Phones: []string{"89999999999", "89999992999"},
				meta:   nil,
			},
			as: ValidationErrors{},
		},
		{
			in: App{
				Version: "49210",
			},
			as: ValidationErrors{},
		},
		{
			in: Response{
				Code: 200,
				Body: "",
			},
			as: ValidationErrors{},
		},
		{
			in: []string{},
			is: ErrNotStruct,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.is != nil {
				require.ErrorIs(t, err, tt.is)
			} else {
				require.ErrorAs(t, err, &tt.as)
			}
		})
	}
}
