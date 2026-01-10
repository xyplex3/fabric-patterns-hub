package calculator

import (
	"errors"
	"math"
)

// Calculator provides basic arithmetic operations.
type Calculator struct {
	precision int
}

// NewCalculator creates a new Calculator with the specified precision.
func NewCalculator(precision int) *Calculator {
	return &Calculator{precision: precision}
}

// Add returns the sum of two numbers.
func (c *Calculator) Add(a, b float64) float64 {
	return a + b
}

// Subtract returns the difference of two numbers.
func (c *Calculator) Subtract(a, b float64) float64 {
	return a - b
}

// Multiply returns the product of two numbers.
func (c *Calculator) Multiply(a, b float64) float64 {
	return a * b
}

// Divide returns the quotient of two numbers.
// It returns an error if the divisor is zero.
func (c *Calculator) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Power returns a raised to the power of b.
func (c *Calculator) Power(a, b float64) float64 {
	return math.Pow(a, b)
}

// SquareRoot returns the square root of a number.
// It returns an error if the number is negative.
func (c *Calculator) SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, errors.New("cannot compute square root of negative number")
	}
	return math.Sqrt(a), nil
}

// Average returns the average of a slice of numbers.
// It returns an error if the slice is empty.
func Average(numbers []float64) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("cannot compute average of empty slice")
	}

	sum := 0.0
	for _, n := range numbers {
		sum += n
	}
	return sum / float64(len(numbers)), nil
}

// Max returns the maximum value from a slice of numbers.
// It returns an error if the slice is empty.
func Max(numbers []float64) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("cannot find max of empty slice")
	}

	max := numbers[0]
	for _, n := range numbers[1:] {
		if n > max {
			max = n
		}
	}
	return max, nil
}

// Min returns the minimum value from a slice of numbers.
// It returns an error if the slice is empty.
func Min(numbers []float64) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("cannot find min of empty slice")
	}

	min := numbers[0]
	for _, n := range numbers[1:] {
		if n < min {
			min = n
		}
	}
	return min, nil
}

// IsEven reports whether n is even.
func IsEven(n int) bool {
	return n%2 == 0
}

// IsPrime reports whether n is a prime number.
func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}
