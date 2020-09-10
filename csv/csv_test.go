package csv

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCSV(t *testing.T) {
	Convey("Given a csv file and a row parser", t, func() {
		f, err := os.Open(filepath.Join("..", "data", "input.csv"))
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		defer f.Close()

		secondRow := []string{"15.599452", "5.808361", "86.617615", "60.111501", "70.807258"}

		rp := NewRowParser(f)
		So(rp, ShouldNotBeNil)

		Convey("You can get the first row of values", func() {
			vals, err := rp.GetRow(1)
			So(err, ShouldBeNil)
			So(vals, ShouldResemble, []string{"37.454012", "95.071431", "73.199394", "59.865848", "nan"})

			Convey("And then the next, without needing to seek", func() {
				vals, err := rp.GetRow(2)
				So(err, ShouldBeNil)
				So(vals, ShouldResemble, secondRow)
				So(rp.seeks, ShouldEqual, 0)
			})
		})

		Convey("You can get the second row of values, by seeking", func() {
			vals, err := rp.GetRow(2)
			So(err, ShouldBeNil)
			So(vals, ShouldResemble, secondRow)
			So(rp.seeks, ShouldEqual, 1)
		})

		Convey("You can get the last row of values", func() {
			vals, err := rp.GetRow(5)
			So(err, ShouldBeNil)
			So(vals, ShouldResemble, []string{"61.185289", "13.949386", "29.214465", "nan", "45.606998"})
		})

		Convey("You can't get non-existant rows", func() {
			vals, err := rp.GetRow(999)
			So(err, ShouldNotBeNil)
			So(vals, ShouldBeNil)
			So(err, ShouldEqual, io.EOF)
		})
	})
}
