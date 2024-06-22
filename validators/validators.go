// Package validators creates Validations type variables with the
// most common conditions for different types that are used frequently
// reducing code duplication
package validators

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Intisma/govalidate"
)

// Function isMoney checks if a given value checks the conditions to be
// considered as a Money type as stated by PostgresSQL documentation
func isMoney(v interface{}) error {

	// Unify to string all possible value types
	var moneyValue string
	switch val := v.(type) {
	case string:
		moneyValue = val
	case int, int64, float32, float64:
		moneyValue = fmt.Sprintf("%v", val)
	default:
		return errors.New("value type is not supported")
	}

	// Check if it's a valid float
	_, err := strconv.ParseFloat(moneyValue, 64)
	if err != nil {
		return errors.New("value is not a valid float")
	}

	// Convert string to float64
	moneyFloat, err := strconv.ParseFloat(moneyValue, 64)
	if err != nil {
		return errors.New("value is not a valid float")
	}

	// Check if value is inside PostgreSQL money type range
	if moneyFloat < -92233720368547758.08 || moneyFloat > 92233720368547758.07 {
		return fmt.Errorf("value %f is out of the valid range for money type", moneyFloat)
	}

	// Check if the value has at most two decimal places
	parts := strings.Split(moneyValue, ".")
	if len(parts) == 2 && len(parts[1]) > 2 {
		return errors.New("value has more than two decimal places")
	}

	// Return nil if all conditions passed
	return nil

}

// Validations to check if a value is a Money type
var MoneyValidation = govalidate.Validations[interface{}]{
	Conditions: []func(interface{}) error{
		func(v interface{}) error { return isMoney(v) },
	}, ValidateMethod: false,
}
