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
		in    interface{}
		is    error
		equal error
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
			equal: ValidationErrors(nil),
		},
		{
			in: App{
				Version: "49210",
			},
			equal: ValidationErrors(nil),
		},
		{
			in: Response{
				Code: 200,
				Body: "",
			},
			equal: ValidationErrors(nil),
		},
		{
			in: []string{},
			is: ErrNotStruct,
		},
		{
			in: User{
				ID:     "389582290458082",
				Name:   "Tod",
				Age:    49,
				Email:  "demo@localhost.com",
				Role:   UserRole("stuff"),
				Phones: []string{"89999999999", "1"},
				meta:   nil,
			},
			is: ErrWrongLen,
		},
		{
			in: Response{
				Code: 205,
				Body: "",
			},
			is: ErrNotExistIn,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.is != nil {
				ve, ok := err.(ValidationErrors)
				if ok {
					require.ErrorIs(t, ve[0].Err, tt.is)
				} else {
					require.ErrorIs(t, err, tt.is)
				}
			}

			if tt.equal != nil {
				require.Equal(t, tt.equal, err)
			}
		})
	}
}
