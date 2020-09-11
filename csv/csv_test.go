package csv

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRowParser(t *testing.T) {
	input := filepath.Join("..", "data", "input.csv")
	firstRow := []string{"37.454012", "95.071431", "73.199394", "59.865848", "nan"}
	secondRow := []string{"15.599452", "5.808361", "86.617615", "60.111501", "70.807258"}
	thirdRow := []string{"2.058449", "96.990985", "nan", "21.233911", "18.182497"}
	fourthRow := []string{"nan", "30.424224", "52.475643", "43.194502", "29.122914"}
	lastRow := []string{"61.185289", "13.949386", "29.214465", "nan", "45.606998"}

	Convey("Given a csv file", t, func() {
		f, err := os.Open(input)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		defer f.Close()

		Convey("And a RowParser", func() {
			rp := NewRowParser(f)
			So(rp, ShouldNotBeNil)

			Convey("You can get the first row of values", func() {
				vals, err := rp.GetRow(1)
				So(err, ShouldBeNil)
				So(vals, ShouldResemble, firstRow)

				Convey("And then the next, without needing to seek", func() {
					vals, err = rp.GetRow(2)
					So(err, ShouldBeNil)
					So(vals, ShouldResemble, secondRow)
					So(rp.seeks, ShouldEqual, 0)

					Convey("And then the previous, with a seek", func() {
						vals, err = rp.GetRow(1)
						So(err, ShouldBeNil)
						So(rp.seeks, ShouldEqual, 1)
						So(vals, ShouldResemble, firstRow)
					})
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
				So(vals, ShouldResemble, lastRow)
			})

			Convey("You can't get a non-existent row", func() {
				vals, err := rp.GetRow(999)
				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, io.EOF)
				So(vals, ShouldBeNil)
			})

			Convey("You can get the first 3 rows of values", func() {
				vals, err := rp.GetRows(1, 3)
				So(err, ShouldBeNil)
				So(vals, ShouldResemble, [][]string{firstRow, secondRow, thirdRow})

				Convey("And then the 3 rows starting from row 2, then row 3", func() {
					vals, err = rp.GetRows(2, 3)
					So(err, ShouldBeNil)
					So(vals, ShouldResemble, [][]string{secondRow, thirdRow, fourthRow})

					vals, err = rp.GetRows(3, 3)
					So(err, ShouldBeNil)
					So(vals, ShouldResemble, [][]string{thirdRow, fourthRow, lastRow})

					Convey("And then subsequent iterations start returning nil rows past the end of the file", func() {
						vals, err = rp.GetRows(4, 3)
						So(err, ShouldBeNil)
						So(vals, ShouldResemble, [][]string{fourthRow, lastRow, nil})

						vals, err = rp.GetRows(5, 3)
						So(err, ShouldBeNil)
						So(vals, ShouldResemble, [][]string{lastRow, nil, nil})
					})
				})
			})

			Convey("You can get the 3 rows of values starting from row 2", func() {
				vals, err := rp.GetRows(2, 3)
				So(err, ShouldBeNil)
				So(vals, ShouldResemble, [][]string{secondRow, thirdRow, fourthRow})
			})

			Convey("You can't get rows that don't exist at all", func() {
				vals, err := rp.GetRows(999, 3)
				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, io.EOF)
				So(vals, ShouldResemble, [][]string{nil, nil, nil})
			})
		})

		Convey("And a CachedRowParser", func() {
			rp, err := NewCachedRowParser(f, 3)
			So(err, ShouldBeNil)
			So(rp, ShouldNotBeNil)

			Convey("You can get the first 3 rows of values, which get cached", func() {
				vals, err := rp.GetRows(1, 3)
				So(err, ShouldBeNil)
				So(vals, ShouldResemble, [][]string{firstRow, secondRow, thirdRow})
				So(rp.cache.Contains(int64(1)), ShouldBeTrue)
				So(rp.cache.Contains(int64(2)), ShouldBeTrue)
				So(rp.cache.Contains(int64(3)), ShouldBeTrue)
				So(rp.cache.Contains(int64(4)), ShouldBeFalse)
				So(rp.reads, ShouldEqual, 3)

				Convey("And then the same again without seeking or reading new data", func() {
					vals, err := rp.GetRows(1, 3)
					So(err, ShouldBeNil)
					So(vals, ShouldResemble, [][]string{firstRow, secondRow, thirdRow})
					So(rp.seeks, ShouldEqual, 0)
					So(rp.reads, ShouldEqual, 3)
				})

				Convey("And then the next 3 rows starting at row 3, which gets 2 from the cache and does 1 read", func() {
					vals, err := rp.GetRows(2, 3)
					So(err, ShouldBeNil)
					So(vals, ShouldResemble, [][]string{secondRow, thirdRow, fourthRow})
					So(rp.seeks, ShouldEqual, 0)
					So(rp.reads, ShouldEqual, 4)
					So(rp.cache.Contains(int64(1)), ShouldBeFalse)
					So(rp.cache.Contains(int64(2)), ShouldBeTrue)
					So(rp.cache.Contains(int64(3)), ShouldBeTrue)
					So(rp.cache.Contains(int64(4)), ShouldBeTrue)
					So(rp.cache.Contains(int64(5)), ShouldBeFalse)
				})
			})

			Convey("You can't get rows that don't exist at all", func() {
				vals, err := rp.GetRows(999, 3)
				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, io.EOF)
				So(vals, ShouldResemble, [][]string{nil, nil, nil})
			})
		})

		Convey("You can't make a CachedRowParser with a non-sensicle cache size", func() {
			rp, err := NewCachedRowParser(f, -1)
			So(err, ShouldNotBeNil)
			So(rp, ShouldBeNil)
		})
	})
}

func BenchmarkRowParserGetRow(b *testing.B) {
	benchmarkRowParserGetRow(newRowParser(b), b)
}

func benchmarkRowParserGetRow(rp csvRowParser, b *testing.B) {
	for n := 0; n < b.N; n++ {
		var err error
		for i := int64(1); err == nil; i++ {
			_, err = rp.GetRow(i)
		}
	}
}

func newRowParser(b *testing.B) *RowParser {
	return NewRowParser(openInput(b))
}

func openInput(b *testing.B) *os.File {
	f, err := os.Open(filepath.Join("..", "data", "input.csv"))
	if err != nil {
		b.Fatal(err)
	}

	return f
}

func BenchmarkCachedRowParserGetRow(b *testing.B) {
	benchmarkRowParserGetRow(newCachedRowParser(b), b)
}

func newCachedRowParser(b *testing.B) *CachedRowParser {
	crp, err := NewCachedRowParser(openInput(b), 3)
	if err != nil {
		b.Fatal(err)
	}

	return crp
}

func BenchmarkRowParserGetRows(b *testing.B) {
	benchmarkRowParserGetRows(newRowParser(b), b)
}

func benchmarkRowParserGetRows(rp csvRowParser, b *testing.B) {
	for n := 0; n < b.N; n++ {
		var err error
		for i := int64(1); err == nil; i++ {
			_, err = rp.GetRows(i, 3)
		}
	}
}

func BenchmarkCachedRowParserGetRows(b *testing.B) {
	benchmarkRowParserGetRows(newCachedRowParser(b), b)
}
