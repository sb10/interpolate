# interpolate

interpolate takes a csv file representing a 2d matrix of number strings, and
outputs the same matrix, but with any "nan" values interpolated as the average
of non-diagonal neighbours.

## Installation

Requires go v1.14 or later. Follow https://golang.org/doc/install, then:

git clone https://github.com/sb10/interpolate.git
cd interpolate
make install

## Usage

interpolate your.csv > interpolated.csv

See the example input and output csv files in the data subdirectory.

## Developers
To develop this code base, you should use TDD. To aid this, the test suite is
written using GoConvey.

To install goconvey:
```
cd ~/somewhere_else
git checkout https://github.com/smartystreets/goconvey.git
go build
mv goconvey $GOPATH/bin/
```

To use goconvey:
```
cd ~/your_clone_of_this_repository
goconvey &
```
This will pop up a browser window which will aid in the red-green-refactor
cycle.

To run the tests on the command line:
`go test ./...` or `make test`

To run the benchmarks:
`go test -run Bench -bench=. ./...` or `make bench`

Before committing any code, you should make sure you haven't introduced any
linting errors. First install the linters:
`curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.31.0`

Then:
`make lint`
