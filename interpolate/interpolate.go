// interpolate package interpolates missing values in 2d matrices.
package interpolate

import (
	"io"
	"strconv"

	"github.com/sb10/interpolate/csv"
)

const nan = "nan"
const rowsNeededToInterpolate = 3

// MiddleRow takes a slice of 3 string slices (each representing a row in a
// 2d matrix, with each element position respresenting a column), and
// interpolates "nan" values as the average of non-diagonal adjacent values,
// which must be string representations of numbers, or "nan".
//
// There must be 3 rows. All rows must be equal length.
//
// It returns a copy of the middle row with any "nan" values interpolated and
// rounded to decimalPlaces decimal places.
func MiddleRow(rows [][]string, decimalPlaces int) []string {
	middleRow := rows[1]
	result := make([]string, len(middleRow))

	for i, val := range middleRow {
		if isNaN(val) {
			result[i] = averageOf(
				decimalPlaces,
				valueAt(rows[0], i),
				valueAt(rows[2], i),
				valueLeftOf(middleRow, i),
				valueRightOf(middleRow, i),
			)
		} else {
			result[i] = val
		}
	}

	return result
}

// isNaN returns true if the given string can't be interpreted as a number.
func isNaN(val string) bool {
	if val == nan {
		return true
	}

	_, err := strconv.ParseFloat(val, 64)

	return err != nil
}

// averageOf returns the average of the given floats as a string, rounded to
// decimalPlaces decimal places. You can supply nil floats and they will be
// ignored. If you only supply nils, nan will be returned.
func averageOf(decimalPlaces int, floats ...*float64) string {
	var total, num float64

	for _, f := range floats {
		if f == nil {
			continue
		}

		num++

		total += *f
	}

	if num == 0 {
		return nan
	}

	return strconv.FormatFloat(total/num, 'f', decimalPlaces, 64)
}

// valueAt returns the stringToFloat() of the value in the row at the given
// index position. If row is nil, returns nil.
func valueAt(row []string, position int) *float64 {
	if row == nil {
		return nil
	}

	return stringToFloat(row[position])
}

// stringToFloat converts the given string to a float. If the string is "nan",
// or otherwise not interpretable as a float, returns nil.
func stringToFloat(num string) *float64 {
	if num == nan {
		return nil
	}

	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return nil
	}

	return &f
}

// valueLeftOf returns the stringToFloat() of the value to the left of the given
// index in the row. If i is 0, returns nil.
func valueLeftOf(row []string, position int) *float64 {
	if position <= 0 {
		return nil
	}

	return stringToFloat(row[position-1])
}

// valueRightOf returns the stringToFloat() of the value to the right of the
// given index in the row. If i is the last element, returns nil.
func valueRightOf(row []string, position int) *float64 {
	if position >= len(row)-1 {
		return nil
	}

	return stringToFloat(row[position+1])
}

// CSVInterpolator interpolates nan values in csv data using MiddleRow().
type CSVInterpolator struct {
	rp            *csv.CachedRowParser
	row           int64
	decimalPlaces int
}

func NewCSVInterpolator(data io.ReadSeeker, decimalPlaces int) (*CSVInterpolator, error) {
	rp, err := csv.NewCachedRowParser(data, rowsNeededToInterpolate)
	if err != nil {
		return nil, err
	}

	return &CSVInterpolator{
		rp:            rp,
		decimalPlaces: decimalPlaces,
	}, nil
}

// NextRow returns the next row of CSV data as a slice of strings, with any
// nan values in the input data replaced with interpreted values.
func (c *CSVInterpolator) NextRow() ([]string, error) {
	c.row++
	rowsToGet := c.numRows()

	rows, err := c.rp.GetRows(c.rowToStartWith(), rowsToGet)
	if err != nil {
		return nil, err
	}

	if rowsToGet != rowsNeededToInterpolate {
		rows = [][]string{nil, rows[0], rows[1]}
	}

	if rows[1] == nil {
		return nil, io.EOF
	}

	return MiddleRow(rows, c.decimalPlaces), nil
}

// numRows returns 3, except at start.
func (c *CSVInterpolator) numRows() int64 {
	if c.row == 1 {
		return rowsNeededToInterpolate - 1
	}

	return rowsNeededToInterpolate
}

// rowToStartWith returns c.row-1, except at start.
func (c *CSVInterpolator) rowToStartWith() int64 {
	if c.row == 1 {
		return 1
	}

	return c.row - 1
}
