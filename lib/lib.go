package lib

import (
	"bufio"
	"encoding/json"
	"errors"
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
		var record Record
		record.ItemType = xArch.ItemType()
			ItemType: x.ItemType
		}
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

func (d *LocalDatabase) Read() (ItemList, error) {
	items := make(ItemList, 0)

	db, openErr := os.Open(d.Path)
	if openErr != nil {
		return nil, openErr
	}
	defer db.Close()

	scanner := bufio.NewScanner(db)
	for scanner.Scan() {
		if scanErr := scanner.Err(); scanErr != nil {
			return nil, scanErr
		}
		i, parseErr := ParseJsonItem(scanner.Text())
		if parseErr != nil {
			return nil, parseErr
		}
		items = append(items, i)
	}
	return items, nil
}

type Record struct {
	ItemType string      `json:"item_type"`
	ItemData interface{} `json:"item_data"`
}

type WebSiteItem struct {
	ItemType    string `json:"item_type"`
	DateCreated int64  `json:"date_created"`
	DateRead    int64  `json:"date_read"`
	Url         string `json:"url"`
	Read        bool   `json:"read"`
}

type WebSiteRecord struct {
	ItemType string      `json:"item_type"`
	ItemData WebSiteItem `json:"item_data"`
}

func NewWebSite(url string) *WebSiteItem {
	return &WebSiteItem{
		ItemType:    "web_site",
		DateCreated: UnixTimeNow(),
		DateRead:    -1,
		Url:         url,
		Read:        false,
	}
}

func (w *WebSiteItem) MarkRead() {
	w.Read = true
	w.DateRead = UnixTimeNow()
}

// inspired by https://eagain.net/articles/go-dynamic-json/
func ParseJsonItem(itemJson string) (interface{}, error) {
	var itemData json.RawMessage
	var env Record = Record{
		ItemData: &itemData,
	}
	if unmarshalErr := json.Unmarshal([]byte(itemJson), &env); unmarshalErr != nil {
		return nil, unmarshalErr
	}
	switch env.ItemType {
	case "web_site":
		var item WebSiteItem
		if unmarshalErr := json.Unmarshal(itemData, &item); unmarshalErr != nil {
			return nil, unmarshalErr
		} else {
			return item, nil
		}
	default:
		err := errors.New(fmt.Sprintf("Attempted to parse unknown item type: %s", env.ItemType))
		return nil, err
	}
}

