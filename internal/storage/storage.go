package storage

import (
	"bufio"
	"os"
)

type Storage struct {
	file os.File
	buf  bufio.Scanner
}

func NewLocalStorage(filepath string) *Storage {
	file, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}

	buf := bufio.NewScanner(file)

	return &Storage{
		file: *file,
		buf:  *buf,
	}
}

func (s *Storage) Close() {
	s.file.Close()
}

func (s *Storage) Read() string {
	s.buf.Scan()
	return s.buf.Text()
}

func (s *Storage) FileIsEmpty() bool {
	stat, err := s.file.Stat()

	if err != nil {
		panic(err)
	}

	return stat.Size() == 0
}
