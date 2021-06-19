package lib

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"
)

func UnixTimeNow() int64 {
	return time.Now().Unix()
}

func SafeWriteFile(filename string, data []byte) error {
	var tmpFilename string = fmt.Sprintf("%s.new", filename)

	// create or truncate file
	outFile, createErr := os.Create(tmpFilename)
	if createErr != nil {
		return createErr
	} else {
		defer outFile.Close()
	}

	// write to temporary file
	writer := bufio.NewWriter(outFile)
	writeN, writeErr := writer.Write(data)
	if writeErr != nil {
		return writeErr
	}
	if writeN != len(data) {
		return errors.New(fmt.Sprintf("%d bytes written of %d total bytes", writeN, len(data)))
	}
	writer.Flush()

	// atomically replace original file with temporary file
	if renameErr := os.Rename(tmpFilename, filename); renameErr != nil {
		return renameErr
	}

	return nil
}
