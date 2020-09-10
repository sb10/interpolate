package csv

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/stoicperlman/fls"
)

// RowParser can parse particular rows from a csv file.
type RowParser struct {
	fls     *fls.File
	parser  *csv.Reader
	lastRow int64
	seeks   int64
}

// NewRowParser creates a new RowParser given an io.ReadSeeker, such as an
// opened file.
func NewRowParser(data io.ReadSeeker) *RowParser {
	return &RowParser{
		fls:    fls.LineFile(data.(*os.File)),
		parser: csv.NewReader(data),
	}
}

// GetRow returns a slice of column values for the given row. Returns an io.EOF
// error if the given row doesn't exist.
func (r *RowParser) GetRow(i int64) ([]string, error) {
	last := r.lastRow
	r.lastRow = i

	if i != last+1 {
		_, err := r.fls.SeekLine(i-1, io.SeekStart)
		r.seeks++
		if err != nil {
			return nil, err
		}
	}

	return r.parser.Read()
}
