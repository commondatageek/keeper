package lib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

func UnixTimeNow() int64 {
	return time.Now().Unix()
}

func SafeWriteFile(filename string, data io.Reader) (int64, error) {
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