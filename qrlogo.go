package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	qrcode "github.com/skip2/go-qrcode"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) PickLogoImage() (string, error) {
	filePath, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select Logo Image",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Images", Pattern: "*.png;*.jpg;*.jpeg;*.gif;*.bmp"},
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

	return base64.StdEncoding.EncodeToString(data), nil
}

func (a *App) GenerateQRWithLogo(req QRRequest, logoBase64 string) QRResponse {
	content, err := formatContent(req.Type, req.Fields)
	if err != nil {
		return QRResponse{Error: err.Error()}
	}

	ecl := parseECL(req.ECL)
	size := req.Size
	if size <= 0 {
		size = 256
	}

	qr, err := qrcode.New(content, ecl)
	if err != nil {
		return QRResponse{Error: fmt.Sprintf("QR generation failed: %v", err)}
	}

	if req.FgColor != "" || req.BgColor != "" {
		qr.ForegroundColor = parseHexColor(req.FgColor, color.Black)
		qr.BackgroundColor = parseHexColor(req.BgColor, color.White)
	}

	qrImg := qr.Image(size)

	logoData, err := base64.StdEncoding.DecodeString(logoBase64)
	if err != nil {
		return QRResponse{Error: fmt.Sprintf("logo decode failed: %v", err)}
	}

	logoImg, _, err := image.Decode(bytes.NewReader(logoData))
	if err != nil {
		return QRResponse{Error: fmt.Sprintf("logo image decode failed: %v", err)}
	}

	logoBounds := logoImg.Bounds()
	logoW := logoBounds.Dx()
	logoH := logoBounds.Dy()

	maxLogoSize := size / 4
	if logoW > maxLogoSize || logoH > maxLogoSize {
		scale := float64(maxLogoSize) / float64(max(logoW, logoH))
		newW := int(float64(logoW) * scale)
		newH := int(float64(logoH) * scale)
		logoImg = resizeImage(logoImg, newW, newH)
		logoBounds = logoImg.Bounds()
		logoW = logoBounds.Dx()
		logoH = logoBounds.Dy()
	}

	offsetX := (size - logoW) / 2
	offsetY := (size - logoH) / 2

	result := image.NewRGBA(qrImg.Bounds())
	draw.Draw(result, result.Bounds(), qrImg, image.Point{}, draw.Src)
	draw.Draw(result, image.Rect(offsetX, offsetY, offsetX+logoW, offsetY+logoH), logoImg, image.Point{}, draw.Over)

	var buf bytes.Buffer
	if err := png.Encode(&buf, result); err != nil {
		return QRResponse{Error: fmt.Sprintf("PNG encode failed: %v", err)}
	}

	resp := QRResponse{
		Image:   base64.StdEncoding.EncodeToString(buf.Bytes()),
		Content: content,
	}

	if resp.Content != "" {
		if err := a.AddToHistory(req, resp.Content); err != nil {
			println("history error:", err.Error())
		}
	}

	return resp
}

func resizeImage(src image.Image, newW, newH int) image.Image {
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			srcX := x * srcW / newW
			srcY := y * srcH / newH
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}

	return dst
}
