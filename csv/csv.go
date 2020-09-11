// csv package provides a way of parsing a csv file by row.
package csv

import (
	"encoding/csv"
	"io"
	"os"

	lru "github.com/hashicorp/golang-lru"
	"github.com/stoicperlman/fls"
)

// rowGetter describes the GetRow() method of a RowParser, allowing it to be
// "overridden".
type rowGetter func(i int64) ([]string, error)

// csvRowParser is implemented by RowParser and CachedRowParser.
type csvRowParser interface {
	GetRow(i int64) ([]string, error)
	GetRows(firstRow int64, numberOfRows int64) ([][]string, error)
}

// RowParser can parse particular rows from a csv file.
type RowParser struct {
	data         io.ReadSeeker
	fls          *fls.File
	parser       *csv.Reader
	lastRow      int64
	seeks        uint64
	getRowMethod rowGetter
}

// NewRowParser creates a new RowParser given an io.ReadSeeker, such as an
// opened file.
func NewRowParser(data io.ReadSeeker) *RowParser {
	rp := &RowParser{
		data:   data,
		fls:    fls.LineFile(data.(*os.File)),
		parser: csv.NewReader(data),
	}
	rp.getRowMethod = rp.GetRow

	return rp
}

// GetRow returns a slice of column values for the given row. Returns an io.EOF
// error if the given row doesn't exist.
func (r *RowParser) GetRow(i int64) ([]string, error) {
	last := r.lastRow
	r.lastRow = i

	if i != last+1 {
		_, err := r.fls.SeekLine(i-1, io.SeekStart)
		if err != nil {
			return nil, err
		}

		r.parser = csv.NewReader(r.data)
		r.seeks++
	}

	return r.parser.Read()
}

// GetRows returns a slice of column values for each of the requested rows.
// Returns an io.EOF error if firstRow doesn't exist, but not if not enough
// rows exist to satisfy numberOfRows.
func (r *RowParser) GetRows(firstRow int64, numberOfRows int64) ([][]string, error) {
	rows := make([][]string, numberOfRows)

	for i := int64(0); i < numberOfRows; i++ {
		row, err := r.getRowMethod(firstRow + i)
		if err != nil {
			if i != 0 {
				err = nil
			}

			return rows, err
		}

		rows[i] = row
	}

	return rows, nil
}

// CachedRowParser is a RowParser that caches the last rows retrieved.
type CachedRowParser struct {
	*RowParser
	cache *lru.Cache
	reads uint64
}

// NewCachedRowParser creates a new CachedRowParser, which is a RowParser that
// will cache the previous cacheSize rows retrieved with GetRow().
// Can return an error if there is a problem creating the cache.
func NewCachedRowParser(data io.ReadSeeker, cacheSize int) (*CachedRowParser, error) {
	cache, err := lru.New(cacheSize)
	if err != nil {
		return nil, err
	}

	crp := &CachedRowParser{
		cache: cache,
	}
	rp := NewRowParser(data)
	rp.getRowMethod = crp.GetRow
	crp.RowParser = rp

	return crp, nil
}

// GetRow is like RowParser.GetRow(), but the result is cached.
func (r *CachedRowParser) GetRow(i int64) ([]string, error) {
	cachedRow, ok := r.cache.Get(i)
	if ok {
		return cachedRow.([]string), nil
	}

	r.reads++

	row, err := r.RowParser.GetRow(i)
	if err != nil {
		return nil, err
	}

	r.cache.Add(i, row)

	return row, err
}
