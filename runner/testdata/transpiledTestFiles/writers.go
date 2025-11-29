package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// format constants
const (
	__dgi_FormatCSV    = "csv"
	__dgi_FormatJSON   = "json"
	__dgi_FormatXML    = "xml"
	__dgi_FormatStdout = "stdout"
)

type __dgi_Record interface {
	ToCSV() []string
	CSVHeaders() []string
	ToJSON() string
	ToXML() string
}

type __dgi_RecordGenerator func(i int) __dgi_Record

type __dgi_OutputWriter func(name string, records []__dgi_Record, outPath string) error

func __dgi_resolveOutputFilePath(outPath, name, ext string) (string, error) {
	if outPath == "" {
		return fmt.Sprintf("%s.%s", name, ext), nil
	}
	if fi, err := os.Stat(outPath); err == nil {
		if fi.IsDir() {
			return filepath.Join(outPath, fmt.Sprintf("%s.%s", name, ext)), nil
		}
		return outPath, nil
	} else if errors.Is(err, os.ErrNotExist) {
		if strings.EqualFold(filepath.Ext(outPath), "."+ext) {
			return outPath, nil
		}
		return filepath.Join(outPath, fmt.Sprintf("%s.%s", name, ext)), nil
	} else {
		return "", fmt.Errorf("error accessing output path %s: %v", outPath, err)
	}
}

func __dgi_getOutputFile(outPath, name, ext string) (*os.File, error) {
	filePath, err := __dgi_resolveOutputFilePath(outPath, name, ext)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("error creating output directory for %s: %v", name, err)
	}
	outputFile, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("error creating output file for %s: %v", name, err)
	}
	return outputFile, nil
}

func __dgi_writeCSV(name string, records []__dgi_Record, outPath string) error {
	csvFile, err := __dgi_getOutputFile(outPath, name, __dgi_FormatCSV)
	if err != nil {
		return fmt.Errorf("error creating CSV file for %s: %v", name, err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	headersWritten := false
	for _, record := range records {
		row := record.ToCSV()
		if !headersWritten {
			if err := writer.Write(record.CSVHeaders()); err != nil {
				return fmt.Errorf("error writing CSV headers for %s: %w", name, err)
			}
			headersWritten = true
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing CSV row for %s: %w", name, err)
		}
	}
	slog.Info(fmt.Sprintf("generated CSV file %s with %d records", csvFile.Name(), len(records)))
	return nil
}

func __dgi_writeJSON(name string, records []__dgi_Record, outPath string) error {
	jsonFile, jsonErr := __dgi_getOutputFile(outPath, name, __dgi_FormatJSON)
	if jsonErr != nil {
		return fmt.Errorf("error creating JSON file for %s: %v", name, jsonErr)
	}
	defer jsonFile.Close()
	jsonWriter := bufio.NewWriter(jsonFile)
	defer jsonWriter.Flush()
	for _, record := range records {
		jsonRow := record.ToJSON()
		fmt.Fprintln(jsonWriter, jsonRow)
	}
	slog.Info(fmt.Sprintf("generated JSON file %s with %d records", jsonFile.Name(), len(records)))
	return nil
}

func __dgi_writeXML(name string, records []__dgi_Record, outPath string) error {
	xmlFile, xmlErr := __dgi_getOutputFile(outPath, name, __dgi_FormatXML)
	if xmlErr != nil {
		return fmt.Errorf("error creating XML file for %s: %v", name, xmlErr)
	}
	defer xmlFile.Close()
	xmlWriter := bufio.NewWriter(xmlFile)
	defer xmlWriter.Flush()
	for _, record := range records {
		xmlRow := record.ToXML()
		fmt.Fprintln(xmlWriter, xmlRow)
	}
	slog.Info(fmt.Sprintf("generated XML file %s with %d records", xmlFile.Name(), len(records)))
	return nil
}

func __dgi_writeStdout(name string, records []__dgi_Record, outPath string) error {
	outWriter := bufio.NewWriter(os.Stdout)
	defer outWriter.Flush()
	for _, record := range records {
		val := interface{}(record)
		rv := reflect.ValueOf(record)
		if rv.Kind() == reflect.Ptr && !rv.IsNil() {
			val = rv.Elem().Interface()
		}
		fmt.Fprintf(outWriter, "%+v%+v\n", name, val)
	}
	return nil
}
