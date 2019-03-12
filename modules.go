package boleto

import (
	"errors"
	"strconv"
	"time"
)

const (
	// min multiplier for module 10
	minModule10 = 1
	// Max multiplier for module 10
	maxModule10 = 2
	// min multiplier for module 11
	minModule11 = 2
	// Max multiplier for module 11
	maxModule11 = 9
)

var (
	// ErrNotFutureDueDate is used when the due date is not in the future
	ErrNotFutureDueDate = errors.New("due date must be in the future")
)

// DateDueFactor use a DateDue type time.Time to return a int,
// with is the quantity of days subsequents from 1997-10-07
func DateDueFactor(dateDue time.Time) (int, error) {
	var dateDueFixed = time.Date(1997, 10, 07, 0, 0, 0, 0, time.UTC)
	dif := dateDue.Sub(dateDueFixed)
	factor := int(dif.Hours() / 24)
	if factor <= 0 {
		return 0, ErrNotFutureDueDate
	}
	return factor, nil
}

// Module10 takes a number and returns his verifier digit (spect an string
// because it may contain left zeros and pad numbers)
// Each digit that makes the Barcode digitable number is multiplied by his multiplier weight,
// the multipliers range from 2 to 1, from left to right
// Multiplication results are summed and divided by ten
func Module10(s string, p int) int {
	// initial multiplier weight, verify if range match
	if p < minModule10 || p > maxModule10 {
		p = maxModule10
	}

	// Create a slice with the numbers
	total := 0
	for _, r := range s {
		c := string(r)
		n, isDot := strconv.Atoi(c)

		// if the multiplier weight is lower then minimal
		if p < minModule10 {
			p = maxModule10
		}

		// if the number could not be found, equals to "."
		if isDot != nil {
			p--
			continue
		}

		// Multiply all numbers using multiplier weight
		m := n * p

		// If the multiplication result is higher then 9,
		// the numbers must be summed between then,
		// For example: m == 18, need to sum 1+8
		if m > 9 {
			// Convert to string and create a range
			multipliers := strconv.Itoa(m)
			numbers := []int{}
			for _, number := range multipliers {
				i, _ := strconv.Atoi(string(number))
				numbers = append(numbers, i)
			}
			// Sum the slice of integers
			m = 0
			for _, number := range numbers {
				m += number
			}
		}

		total += m
		p--

	}

	// End by dividing
	dv := total % 10
	if dv >= 10 {
		dv = 0
	}
	return dv
}

// Module11 takes a number and returns his verifier digit (spect an string
// because it may contain left zeros and pad numbers)
// Each digit that makes up our number is multiplied by his multiplier weight,
// the multipliers range from 9 to 2, from right to left
// Multiplication results are summed and divided by eleven
func Module11(s string) int {
	// Create a slice with the numbers
	numbers := make([]int, len(s))
	for i, r := range s {
		c := string(r)
		n, _ := strconv.Atoi(c)
		numbers[i] = n
	}
	numbersLen := len(numbers)

	// initial multiplier weight
	var p = maxModule11
	if numbersLen > 11 {
		p = minModule11
	}

	// Inverse the numbers creating for loop
	// Multiply all numbers using multiplier weight
	total := 0
	for i := len(numbers) - 1; i >= 0; i-- {
		n := numbers[i]
		total += n * p

		// If the numbers length is higher than 11,
		// we need to inverse the min and max
		if numbersLen > 11 {
			p++
			// if the multiplier weight is higher then max
			if p > maxModule11 {
				p = minModule11
			}
			continue
		}

		p--
		// if the multiplier weight is lower then minimal
		if p < minModule11 {
			p = maxModule11
		}
	}

	// If the numbers length is higher than 11,
	// we need to divide also by 11
	if numbersLen > 11 {
		dv := total % 11
		dv = 11 - dv
		// If the verifier digit is equal 0, 10, 11,
		// need to be always 1
		if dv == 0 || dv == 10 || dv == 11 {
			dv = 1
		}
		return dv
	}

	// End by dividing
	return total % 11
}
