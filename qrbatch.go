package main

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type CSVParseResult struct {
	Headers []string          `json:"headers"`
	Rows    []map[string]string `json:"rows"`
	Error   string            `json:"error,omitempty"`
}

type BatchResult struct {
	Index   int    `json:"index"`
	Image   string `json:"image"`
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

func (a *App) PickCSVFile() (string, error) {
	filePath, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select CSV File",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "CSV Files", Pattern: "*.csv"},
			{DisplayName: "Text Files", Pattern: "*.txt"},
			{DisplayName: "All Files", Pattern: "*"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("dialog failed: %v", err)
	}
	if filePath == "" {
		return "", nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read failed: %v", err)
	}

	return string(data), nil
}

func (a *App) ParseCSV(csvData string) CSVParseResult {
	reader := csv.NewReader(strings.NewReader(csvData))
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		return CSVParseResult{Error: fmt.Sprintf("failed to read headers: %v", err)}
	}

	var rows []map[string]string
	lineNum := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		lineNum++
		if err != nil {
			return CSVParseResult{Error: fmt.Sprintf("error at line %d: %v", lineNum, err)}
		}

		row := make(map[string]string)
		for i, value := range record {
			if i < len(headers) {
				row[headers[i]] = value
			}
		}
		rows = append(rows, row)
	}

	if rows == nil {
		rows = []map[string]string{}
	}

	return CSVParseResult{Headers: headers, Rows: rows}
}

func (a *App) GenerateBatch(req QRRequest, csvData string) []BatchResult {
	parsed := a.ParseCSV(csvData)
	if parsed.Error != "" {
		return []BatchResult{{Error: parsed.Error}}
	}

	if len(parsed.Rows) == 0 {
		return []BatchResult{{Error: "CSV has no data rows"}}
	}

	var results []BatchResult

	for i, row := range parsed.Rows {
		rowReq := QRRequest{
			Type:    req.Type,
			Fields:  row,
			ECL:     req.ECL,
			Size:    req.Size,
			FgColor: req.FgColor,
			BgColor: req.BgColor,
		}

		content, err := formatContent(req.Type, row)
		if err != nil {
			results = append(results, BatchResult{
				Index: i,
				Error: fmt.Sprintf("row %d: %v", i+1, err),
			})
			continue
		}

		ecl := parseECL(req.ECL)
		size := req.Size
		if size <= 0 {
			size = 256
		}

		qr, err := qrcode.New(content, ecl)
		if err != nil {
			results = append(results, BatchResult{
				Index: i,
				Error: fmt.Sprintf("row %d: %v", i+1, err),
			})
			continue
		}

		if req.FgColor != "" || req.BgColor != "" {
			qr.ForegroundColor = parseHexColor(req.FgColor, qr.ForegroundColor)
			qr.BackgroundColor = parseHexColor(req.BgColor, qr.BackgroundColor)
		}

		png, err := qr.PNG(size)
		if err != nil {
			results = append(results, BatchResult{
				Index: i,
				Error: fmt.Sprintf("row %d: %v", i+1, err),
			})
			continue
		}

		if err := a.AddToHistory(rowReq, content); err != nil {
			println("history error:", err.Error())
		}

		results = append(results, BatchResult{
			Index:   i,
			Image:   base64.StdEncoding.EncodeToString(png),
			Content: content,
		})
	}

	return results
}
