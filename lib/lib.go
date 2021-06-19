package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

const ioCopyBufferSize = 4_194_304

type ItemList []interface{}

func UnixTimeNow() int64 {
	return time.Now().Unix()
}

func SafeWriteFile(filename string, data io.Reader) (int64, error) {
	// write first to temporary file so as not to corrupt existing file
	var tmpFilename string = fmt.Sprintf("%s.new", filename)

	// create or truncate file
	outFile, createErr := os.Create(tmpFilename)
	if createErr != nil {
		return -1, createErr
	} else {
		defer outFile.Close()
	}

	// write to temporary file
	writer := bufio.NewWriter(outFile)
	writeN, writeErr := io.Copy(writer, data)
	if writeErr != nil {
		return writeN, writeErr
	}
	writer.Flush()

	// atomically replace original file with temporary file
	if renameErr := os.Rename(tmpFilename, filename); renameErr != nil {
		return writeN, renameErr
	}

	return writeN, nil
}

type DatabaseWriter interface {
	Write(ItemList) (n int64, err error)
}

type DatabaseReader interface {
	Read() (ItemList, error)
}

type LocalDatabase struct {
	Path string
}

func NewLocalDatabase(path string) *LocalDatabase {
	return &LocalDatabase{
		Path: path,
	}
}

func (d *LocalDatabase) Write(database ItemList) (n int64, err error) {
	reader, writer := io.Pipe()
	bufReader := bufio.NewReaderSize(reader, ioCopyBufferSize)
	bufWriter := bufio.NewWriterSize(writer, ioCopyBufferSize)
	done := make(chan error)

	// asynchronous write to file
	go func() {
		_, writeErr := SafeWriteFile(d.Path, bufReader)
		done <- writeErr
	}()

	// while we stream each item
	for _, x := range database {
		jsonBytes, marshalError := json.Marshal(x)
		if marshalError != nil {
			return n, marshalError
		}
		_, writeErr := bufWriter.Write(jsonBytes)
		if writeErr != nil {
			return n, writeErr
		}
		_, writeErr = bufWriter.Write([]byte{'\n'}) // newline between each record
		if writeErr != nil {
			return n, writeErr
		}
		n++ // keep track of the number of items we've successfully written
	}
	bufWriter.Flush()
	writer.Close()

	result := <-done
	return n, result
}

func (d *LocalDatabase) Read() error {
	panic("LocalDatabase.Read() not yet implemented")
}
