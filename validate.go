// Package govalidate implements utility types and functions for
// manipulating validation of variables with standarized methods
package govalidate

import "fmt"

// Struct holding a generic value, a set of conditions and a boolean
// indicating if all validations must be performed (true) or only until
// one of them fails (false)
type Validations struct {
	Value          interface{}
	Conditions     []func(interface{}) error
	ValidateMethod bool
}

// Function to create an empty Validations variable
func CreateValidations[T any]() Validations {
	return Validations{
		Value:          *new(interface{}), // Initialize with the zero value of T
		Conditions:     []func(interface{}) error{},
		ValidateMethod: false, // Default to validating until a condition fails
	}
}

// Function to add another condition to be validated to an existing Validations
// variable
func AddCondition(v *Validations, condition func(interface{}) error) {
	v.Conditions = append(v.Conditions, condition)
}

// Function to copy a Validations variable
func CopyValidations(v Validations) Validations {
	copyConditions := append([]func(interface{}) error(nil), v.Conditions...)
	copy := Validations{
		Value:          v.Value,
		Conditions:     copyConditions,
		ValidateMethod: v.ValidateMethod,
	}
	return copy
}

// Function to set the value of a Validations variable
func SetValueValidations(v *Validations, value interface{}) {
	v.Value = value
}

// Function to set the method of validation to a Validations variable
// where true means validating all conditions and false validation until one
// of the condition fails
func SetMethodValidations(v *Validations, validateMethod bool) {
	v.ValidateMethod = validateMethod
}

// Function to check all the conditions to the value of a Validations variable
// and returning the error (if any)
func Validate(v Validations) error {
	if v.ValidateMethod {
		return validate(v.Value, v.Conditions)
	} else {
		return validateUntilFailure(v.Value, v.Conditions)
	}
}

// Function to validate all conditions, combining all the errors encountered.
func validate(value interface{}, funcs []func(interface{}) error) error {
	var combinedErr error
	for _, f := range funcs {
		err := f(value)
		if err != nil {
			if combinedErr == nil {
				combinedErr = err
			} else {
				combinedErr = fmt.Errorf("%v; %v", combinedErr, err)
			}
		}
	}
	return combinedErr
}

// Function to validate until a condition fails, returning the error of that
// condition
func validateUntilFailure(value interface{}, funcs []func(interface{}) error) error {
	for _, f := range funcs {
		err := f(value)
		if err != nil {
			return err
		}
	}
	return nil
}
