// Package govalidate implements utility types and functions for
// manipulating validation of variables with standarized methods
package govalidate

import "fmt"

// Struct holding a generic value, a set of conditions and a boolean
// indicating if all validations must be performed (true) or only until
// one of them fails (false)
type Validations[T any] struct {
	Value          T
	Conditions     []func(T) error
	ValidateMethod bool
}

// Function to create an empty Validations variable
func CreateValidations[T any]() Validations[T] {
	return Validations[T]{
		Value:          *new(T), // Initialize with the zero value of T
		Conditions:     []func(T) error{},
		ValidateMethod: false, // Default to validating until a condition fails
	}
}

// Function to add another condition to be validated to an existing Validations
// variable
func AddCondition[T any](v *Validations[T], condition func(T) error) {
	v.Conditions = append(v.Conditions, condition)
}

// Function to copy a Validations variable
func CopyValidations[T any](v Validations[T]) Validations[T] {
	copyConditions := append([]func(T) error(nil), v.Conditions...)
	copy := Validations[T]{
		Value:          v.Value,
		Conditions:     copyConditions,
		ValidateMethod: v.ValidateMethod,
	}
	return copy
}

// Function to set the value of a Validations variable
func SetValueValidations[T any](v *Validations[T], value T) {
	v.Value = value
}

// Function to set the method of validation to a Validations variable
// where true means validating all conditions and false validation until one
// of the condition fails
func SetMethodValidations[T any](v *Validations[T], validateMethod bool) {
	v.ValidateMethod = validateMethod
}

// Function to check all the conditions to the value of a Validations variable
// and returning the error (if any)
func Validate[T any](v Validations[T]) error {
	if v.ValidateMethod {
		return validate(v.Value, v.Conditions)
	} else {
		return validateUntilFailure(v.Value, v.Conditions)
	}
}

// Function to validate all conditions, combining all the errors encountered.
func validate[T any](value T, funcs []func(T) error) error {
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
func validateUntilFailure[T any](value T, funcs []func(T) error) error {
	for _, f := range funcs {
		err := f(value)
		if err != nil {
			return err
		}
	}
	return nil
}
