package interpolate

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInterpolateMiddleRow(t *testing.T) {
	firstRow := []string{"2.23", "3.34", "4.45"}
	thirdRow := []string{"5.56", "6.67", "7.78"}

	Convey("Interpolate.MiddleRow doesn't change a fully specified row", t, func() {
		noNanRow := []string{"8.89", "9.99", "10"}
		result := MiddleRow([][]string{firstRow, noNanRow, thirdRow}, 2)
		So(result, ShouldResemble, noNanRow)
	})

	Convey("It interpoaltes a nan on the left", t, func() {
		nanLeftRow := []string{"nan", "9.99", "10"}
		result := MiddleRow([][]string{firstRow, nanLeftRow, thirdRow}, 2)
		So(result, ShouldResemble, []string{"5.93", "9.99", "10"})
	})

	Convey("It interpoaltes a nan on the right", t, func() {
		nanRightRow := []string{"8.89", "9.99", "nan"}
		result := MiddleRow([][]string{firstRow, nanRightRow, thirdRow}, 2)
		So(result, ShouldResemble, []string{"8.89", "9.99", "7.41"})
	})

	Convey("It interpoaltes a nan in the middle", t, func() {
		nanMidRow := []string{"8.89", "nan", "10"}
		result := MiddleRow([][]string{firstRow, nanMidRow, thirdRow}, 2)
		So(result, ShouldResemble, []string{"8.89", "7.22", "10"})
	})

	Convey("It handles adjacent nans in a row", t, func() {
		nanAdjacentRow := []string{"8.89", "nan", "nan"}
		result := MiddleRow([][]string{firstRow, nanAdjacentRow, thirdRow}, 2)
		So(result, ShouldResemble, []string{"8.89", "6.30", "6.12"})
	})

	Convey("It handles first and third rows being nil", t, func() {
		nanAdjacentRow := []string{"8.89", "nan", "nan"}
		result := MiddleRow([][]string{nil, nanAdjacentRow, nil}, 2)
		So(result, ShouldResemble, []string{"8.89", "8.89", nan})
	})

	Convey("It handles nonsense data", t, func() {
		result := MiddleRow([][]string{{"2.23", "lghlh&3", "4.45"}, {"8.89", "asdgh!sdf", "10"}, thirdRow}, 2)
		So(result, ShouldResemble, []string{"8.89", "8.52", "10"})
	})
}

func TestInterpolateCSV(t *testing.T) {
	input := filepath.Join("..", "data", "input.csv")
	expected := [][]string{
		{"37.454012", "95.071431", "73.199394", "59.865848", "65.336553"},
		{"15.599452", "5.808361", "86.617615", "60.111501", "70.807258"},
		{"2.058449", "96.990985", "64.329538", "21.233911", "18.182497"},
		{"31.222654", "30.424224", "52.475643", "43.194502", "29.122914"},
		{"61.185289", "13.949386", "29.214465", "39.338655", "45.606998"},
	}

	Convey("Given a csv file of 2d matrix data with nans, and a CSVInterpolator(", t, func() {
		f, err := os.Open(input)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		defer f.Close()

		c, err := NewCSVInterpolator(f, 6)
		So(err, ShouldBeNil)

		Convey("You can interpolate the nans", func() {
			var i int
			for {
				row, err := c.NextRow()
				if err != nil {
					break
				}
				So(row, ShouldResemble, expected[i])
				i++
			}
			So(i, ShouldEqual, 5)
		})
	})
}
