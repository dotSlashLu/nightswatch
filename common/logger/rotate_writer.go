package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotateWriter struct {
	lock      sync.Mutex
	filename  string
	fp        *os.File
	wrote     int
	sizeLimit int
}

func NewRotateWriter(filename string, sizeLimit int) (*RotateWriter, error) {
	w := &RotateWriter{filename: filename, sizeLimit: sizeLimit}
	logFileStat, err := os.Stat(filename)
	// Log containing folder must be created beforehand
	dir := filepath.Dir(filename)
	_, err = os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		panic(fmt.Sprintf("Log directory %s does not exist", dir))
	}
	// It's ok if the folder exists but the file doesn't
	if err != nil && !os.IsNotExist(err) {
		panic(fmt.Sprintf("Can't stat log file %s: %s\n",
			filename, err.Error()))
	}
	if logFileStat.Size() >= int64(sizeLimit) {
		if err := w.Rotate(); err != nil {
			return nil, err
		}
		return w, nil
	}
	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	w.fp, err = os.OpenFile(filename, flags, 0644)
	if err != nil {
		panic("can't open log file " + filename)
	}
	return w, nil
}

// TODO: buffered IO
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	wrote, err := w.fp.Write(output)
	if err != nil {
		return 0, err
	}
	w.wrote += wrote
	if w.wrote >= w.sizeLimit {
		err = w.Rotate()
		if err != nil {
			return 0, err
		}
	}
	return wrote, err
}

func (w *RotateWriter) Rotate() (err error) {
	// Close existing file if open
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}

	// Rename dest file if it already exists
	_, err = os.Stat(w.filename)
	if err == nil {
		t := time.Now().Format(time.RFC3339)
		err = os.Rename(w.filename, w.filename+"."+t)
		if err != nil {
			return
		}
	}

	// Create a file.
	w.fp, err = os.Create(w.filename)
	return
}
