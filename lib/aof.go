package lib

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

type AppendOnlyFile struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.RWMutex
}

func NewAppendOnlyFile(path string) (*AppendOnlyFile, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666) //o666 is for read-write permission for files
	if err != nil {
		return nil, err
	}

	aof := AppendOnlyFile{
		file: f,
		rd:   bufio.NewReader(f),
	}

	go syncFileEverySecond(&aof)

	return &aof, nil
}

func syncFileEverySecond(a *AppendOnlyFile) {
	for {
		a.mu.Lock()
		// fmt.Println("Lock acquired on file")

		a.file.Sync()
		// fmt.Println("File synchronization starts")

		a.mu.Unlock()
		// fmt.Println("File synchronization ends & Lock released")

		time.Sleep(time.Second)
	}
}

func (a *AppendOnlyFile) Close() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.file.Close()
}

func (a *AppendOnlyFile) Write(v Value) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, err := a.file.Write(v.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (a *AppendOnlyFile) Read(fn func(v Value)) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.file.Seek(0, io.SeekStart)

	reader := NewResp(a.file)

	for {
		value, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				break
			}
			return err
		}
		fn(value)
	}

	return nil
}
