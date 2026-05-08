package main

import (
	"bytes"
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/jung-kurt/gofpdf"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) SaveSheetPDF(req QRRequest, cols, rows int) (string, error) {
	content, err := formatContent(req.Type, req.Fields)
	if err != nil {
		return "", err
	}

	ecl := parseECL(req.ECL)
	qr, err := qrcode.New(content, ecl)
	if err != nil {
		return "", fmt.Errorf("QR generation failed: %v", err)
	}

	if req.FgColor != "" || req.BgColor != "" {
		qr.ForegroundColor = parseHexColor(req.FgColor, qr.ForegroundColor)
		qr.BackgroundColor = parseHexColor(req.BgColor, qr.BackgroundColor)
	}

	qrSize := 200
	qrPNG, err := qr.PNG(qrSize)
	if err != nil {
		return "", fmt.Errorf("PNG encoding failed: %v", err)
	}

	filePath, err := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		DefaultFilename: "qr-sheet.pdf",
		Title:           "Save QR Sheet as PDF",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "PDF Document", Pattern: "*.pdf"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("dialog failed: %v", err)
	}
	if filePath == "" {
		return "", nil
	}

	imgReader := bytes.NewReader(qrPNG)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pageW, pageH := 210.0, 297.0
	margin := 10.0
	cellW := (pageW - 2*margin) / float64(cols)
	cellH := (pageH - 2*margin) / float64(rows)

	qrMM := cellW
	if cellH < cellW {
		qrMM = cellH
	}
	qrMM -= 4

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			x := margin + float64(c)*cellW + (cellW-qrMM)/2
			y := margin + float64(r)*cellH + (cellH-qrMM)/2

			imgReader.Seek(0, 0)
			option := gofpdf.ImageOptions{
				ImageType: "PNG",
				ReadDpi:   false,
			}
			pdf.RegisterImageOptionsReader(fmt.Sprintf("qr_%d_%d", r, c), option, imgReader)
			pdf.ImageOptions(
				fmt.Sprintf("qr_%d_%d", r, c),
				x, y, qrMM, qrMM,
				false, option, 0, "",
			)
		}
	}

	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return "", fmt.Errorf("pdf write failed: %v", err)
	}

	return filePath, nil
}
