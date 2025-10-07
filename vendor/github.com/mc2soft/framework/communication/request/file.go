package request

import "io"

// File is a structure for file sending or receiving.
type File struct {
	Reader   io.Reader
	Filename string
}
