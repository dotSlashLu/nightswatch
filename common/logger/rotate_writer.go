package logger

// https://stackoverflow.com/a/28797984/798042
import (
    "os"
    "sync"
    "time"
)

type RotateWriter struct {
    lock        sync.Mutex
    filename    string // should be set to the actual filename
    fp          *os.File
    wrote       int
    sizeLimit   int
}

// Make a new RotateWriter. Return nil if error occurs during setup.
func NewRotateWriter(filename string, sizeLimit int) (*RotateWriter, error) {
    w := &RotateWriter{filename: filename, sizeLimit: sizeLimit}
    err := w.Rotate()
    if err != nil {
        return nil, err
    }
    return w, nil
}

// Write satisfies the io.Writer interface.
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
        err = os.Rename(w.filename, w.filename + "." + t)
        if err != nil {
            return
        }
    }

    // Create a file.
    w.fp, err = os.Create(w.filename)
    return
}
