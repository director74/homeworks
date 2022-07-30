package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrEmptyOptions    = errors.New("empty tag options")
	ErrNotFilledOption = errors.New("tag option not filled properly")
	ErrNotStruct       = errors.New("incoming variable is not structure")
	ErrWrongLen        = errors.New("length not equal required")
	ErrLower           = errors.New("lower than")
	ErrBigger          = errors.New("bigger than")
	ErrNotExistIn      = errors.New("value not exist in")
	ErrRegexpNotMatch  = errors.New("string dont match regexp")
)

type validationOption struct {
	Name  string
	Value string
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError
type validationOptions []validationOption

func (v ValidationErrors) Error() string {
	var sb strings.Builder

	for _, ver := range v {
		sb.WriteString(ver.Field)
		sb.WriteString(": ")
		sb.WriteString(ver.Err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

type validator struct {
	options   validationOptions
	errors    ValidationErrors
	field     reflect.Value
	fieldName string
}

func (v *validator) readOptions(tagValue string) error {
	if tagValue == "" {
		return ErrEmptyOptions
	}

	actions := strings.Split(tagValue, "|")
	for _, action := range actions {
		parts := strings.Split(action, ":")
		if len(parts) == 2 {
			v.options = append(v.options, validationOption{strings.Title(parts[0]), parts[1]})
		} else {
			return ErrNotFilledOption
		}
	}

	return nil
}

func (v *validator) setFieldInfo(fieldName string, field reflect.Value) {
	v.field = field
	v.fieldName = fieldName
}

func (v *validator) appendError(err error) {
	v.errors = append(v.errors, ValidationError{Field: v.fieldName, Err: err})
}

func (v *validator) getValidationErrors() ValidationErrors {
	return v.errors
}

func (v *validator) LenString(lenVal int64) {
	val := v.field.String()
	if int64(len(val)) != lenVal {
		v.appendError(fmt.Errorf("%v %w %v", val, ErrWrongLen, lenVal))
	}
}

func (v *validator) LenSliceString(lenVal int64) {
	for i := 0; i < v.field.Len(); i++ {
		val := v.field.Index(i).String()
		if int64(len(val)) != lenVal {
			v.appendError(fmt.Errorf("%v %w %v", val, ErrWrongLen, lenVal))
		}
	}
}

func (v *validator) MinInt(minVal int64) {
	if v.field.Int() < minVal {
		v.appendError(fmt.Errorf("%v %w %v", v.field.Int(), ErrLower, minVal))
	}
}

func (v *validator) MinIntSlice(minVal int64) {
	for i := 0; i < v.field.Len(); i++ {
		val := v.field.Index(i).Int()
		if val < minVal {
			v.appendError(fmt.Errorf("%v %w %v", val, ErrLower, minVal))
		}
	}
}

func (v *validator) MaxInt(maxVal int64) {
	if v.field.Int() > maxVal {
		v.appendError(fmt.Errorf("%v %w %v", v.field.Int(), ErrBigger, maxVal))
	}
}

func (v *validator) MaxIntSlice(maxVal int64) {
	for i := 0; i < v.field.Len(); i++ {
		val := v.field.Index(i).Int()
		if val > maxVal {
			v.appendError(fmt.Errorf("%v %w %v", val, ErrBigger, maxVal))
		}
	}
}

func (v *validator) inStringCompare(val string, incomingSet []string) {
	var founded bool

	for _, compareItem := range incomingSet {
		if val == compareItem {
			founded = true
		}
	}

	if !founded {
		v.appendError(fmt.Errorf("%v %w %v", val, ErrNotExistIn, incomingSet))
	}
}

func (v *validator) inIntCompare(val int64, incomingSet []string) {
	var founded bool

	for _, compareItem := range incomingSet {
		converted, _ := strconv.ParseInt(compareItem, 10, 64)
		if val == converted {
			founded = true
		}
	}

	if !founded {
		v.appendError(fmt.Errorf("%v %w %v", val, ErrNotExistIn, incomingSet))
	}
}

func (v *validator) InString(incomingSet string) {
	v.inStringCompare(v.field.String(), strings.Split(incomingSet, ","))
}

func (v *validator) InSliceString(incomingSet string) {
	iSet := strings.Split(incomingSet, ",")

	for i := 0; i < v.field.Len(); i++ {
		v.inStringCompare(v.field.Index(i).String(), iSet)
	}
}

func (v *validator) InInt(incomingSet string) {
	v.inIntCompare(v.field.Int(), strings.Split(incomingSet, ","))
}

func (v *validator) InSliceInt(incomingSet string) {
	iSet := strings.Split(incomingSet, ",")

	for i := 0; i < v.field.Len(); i++ {
		v.inIntCompare(v.field.Index(i).Int(), iSet)
	}
}

func (v *validator) RegexpString(expression string) error {
	val := v.field.String()
	r, err := regexp.Compile(expression)
	if err != nil {
		return err
	}
	if !r.MatchString(val) {
		v.appendError(fmt.Errorf("%v %w %v", val, ErrRegexpNotMatch, expression))
	}

	return nil
}

func (v *validator) RegexpSliceString(expression string) error {
	r, err := regexp.Compile(expression)
	if err != nil {
		return err
	}
	for i := 0; i < v.field.Len(); i++ {
		val := v.field.Index(i).String()
		founded := r.MatchString(val)
		if !founded {
			v.appendError(fmt.Errorf("%v %w %v", val, ErrRegexpNotMatch, expression))
		}
	}
	return nil
}

func (v *validator) check() error {
	var in []reflect.Value
	var methodName string

	if v.field.Kind() == reflect.Slice {
		elemKind := v.field.Type().Elem().Kind().String()
		methodName = strings.Title(v.field.Kind().String()) + strings.Title(elemKind)
	} else {
		methodName = strings.Title(v.field.Kind().String())
	}

	for _, option := range v.options {
		method := reflect.ValueOf(v).MethodByName(option.Name + methodName)
		mtType := method.Type()

		if mtType.In(0).Kind() == reflect.Int64 {
			optValue, _ := strconv.ParseInt(option.Value, 10, 64)
			in = []reflect.Value{reflect.ValueOf(optValue)}
		} else {
			in = []reflect.Value{reflect.ValueOf(option.Value)}
		}

		result := method.Call(in)
		if result != nil {
			if errCallResult, ok := result[0].Interface().(error); ok && errCallResult != nil {
				return errCallResult
			}
		}
	}

	return nil
}

func Validate(v interface{}) error {
	var allFieldErrors ValidationErrors

	reflectedStruct := reflect.ValueOf(v)

	if reflectedStruct.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	reflectedStructType := reflectedStruct.Type()
	for i := 0; i < reflectedStruct.NumField(); i++ {
		structField := reflectedStruct.Field(i)
		structFieldType := reflectedStructType.Field(i)

		// Проверять только публичные свойства
		if structFieldType.Name == strings.Title(structFieldType.Name) {
			strOptions, needValidate := structFieldType.Tag.Lookup("validate")
			if needValidate {
				validator := validator{}
				errParse := validator.readOptions(strOptions)
				if errParse == nil {
					validator.setFieldInfo(structFieldType.Name, structField)
					errCheckProcess := validator.check()
					if errCheckProcess != nil {
						return errCheckProcess
					}
					fieldErrors := validator.getValidationErrors()
					allFieldErrors = append(allFieldErrors, fieldErrors...)
				} else {
					return errParse
				}
			}
		}
	}

	return allFieldErrors
}
