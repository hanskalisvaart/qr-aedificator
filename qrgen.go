package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/url"
	"os"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type QRRequest struct {
	Type   string            `json:"type"`
	Fields map[string]string `json:"fields"`
}

type QRResponse struct {
	Image   string `json:"image"`
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

func (a *App) GenerateQR(req QRRequest) QRResponse {
	content, err := formatContent(req.Type, req.Fields)
	if err != nil {
		return QRResponse{Error: err.Error()}
	}

	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return QRResponse{Error: fmt.Sprintf("QR generation failed: %v", err)}
	}

	png, err := qr.PNG(256)
	if err != nil {
		return QRResponse{Error: fmt.Sprintf("PNG encoding failed: %v", err)}
	}

	return QRResponse{
		Image:   base64.StdEncoding.EncodeToString(png),
		Content: content,
	}
}

func formatContent(qrType string, fields map[string]string) (string, error) {
	switch qrType {
	case "text":
		if v, ok := fields["text"]; ok {
			return v, nil
		}
		return "", fmt.Errorf("missing text field")

	case "url":
		return fields["url"], nil

	case "email":
		to := fields["to"]
		subject := url.QueryEscape(fields["subject"])
		body := url.QueryEscape(fields["body"])
		s := "mailto:" + to
		if subject != "" || body != "" {
			s += "?"
			parts := []string{}
			if subject != "" {
				parts = append(parts, "subject="+subject)
			}
			if body != "" {
				parts = append(parts, "body="+body)
			}
			s += strings.Join(parts, "&")
		}
		return s, nil

	case "phone":
		return "tel:" + fields["phone"], nil

	case "sms":
		phone := fields["phone"]
		message := url.QueryEscape(fields["message"])
		s := "smsto:" + phone
		if message != "" {
			s += ":" + message
		}
		return s, nil

	case "wifi":
		ssid := fields["ssid"]
		password := fields["password"]
		encryption := fields["encryption"]
		hidden := fields["hidden"]

		s := "WIFI:"
		if encryption != "" {
			s += "T:" + encryption + ";"
		}
		if ssid != "" {
			s += "S:" + ssid + ";"
		}
		if password != "" {
			s += "P:" + password + ";"
		}
		if hidden == "true" {
			s += "H:true;"
		}
		s += ";"
		return s, nil

	case "vcard":
		name := fields["name"]
		phone := fields["phone"]
		email := fields["email"]
		org := fields["org"]
		title := fields["title"]
		addr := fields["address"]

		var lines []string
		lines = append(lines, "BEGIN:VCARD")
		lines = append(lines, "VERSION:3.0")
		if name != "" {
			parts := strings.SplitN(name, " ", 2)
			given := parts[0]
			family := ""
			if len(parts) > 1 {
				family = parts[1]
			}
			lines = append(lines, fmt.Sprintf("N:%s;%s;;;", family, given))
			lines = append(lines, fmt.Sprintf("FN:%s", name))
		}
		if phone != "" {
			lines = append(lines, "TEL:"+phone)
		}
		if email != "" {
			lines = append(lines, "EMAIL:"+email)
		}
		if org != "" {
			lines = append(lines, "ORG:"+org)
		}
		if title != "" {
			lines = append(lines, "TITLE:"+title)
		}
		if addr != "" {
			lines = append(lines, "ADR:"+addr)
		}
		lines = append(lines, "END:VCARD")
		return strings.Join(lines, "\n"), nil

	case "geo":
		lat := fields["latitude"]
		lon := fields["longitude"]
		return fmt.Sprintf("geo:%s,%s", lat, lon), nil

	case "calendar":
		title := fields["title"]
		start := fields["start"]
		end := fields["end"]
		location := fields["location"]
		description := fields["description"]

		var lines []string
		lines = append(lines, "BEGIN:VEVENT")
		if title != "" {
			lines = append(lines, "SUMMARY:"+title)
		}
		if start != "" {
			lines = append(lines, "DTSTART:"+start)
		}
		if end != "" {
			lines = append(lines, "DTEND:"+end)
		}
		if location != "" {
			lines = append(lines, "LOCATION:"+location)
		}
		if description != "" {
			lines = append(lines, "DESCRIPTION:"+description)
		}
		lines = append(lines, "END:VEVENT")
		return strings.Join(lines, "\n"), nil

	case "bitcoin":
		addr := fields["address"]
		amount := fields["amount"]
		label := url.QueryEscape(fields["label"])

		s := "bitcoin:" + addr
		if amount != "" || label != "" {
			s += "?"
			parts := []string{}
			if amount != "" {
				parts = append(parts, "amount="+amount)
			}
			if label != "" {
				parts = append(parts, "label="+label)
			}
			s += strings.Join(parts, "&")
		}
		return s, nil

	default:
		return "", fmt.Errorf("unknown QR type: %s", qrType)
	}
}

func (a *App) SaveQRPNG(base64Data string) (string, error) {
	filePath, err := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		DefaultFilename: "qrcode.png",
		Title:           "Save QR Code",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "PNG Image", Pattern: "*.png"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("dialog failed: %v", err)
	}
	if filePath == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("decode failed: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("write failed: %v", err)
	}

	return filePath, nil
}

func (a *App) SaveSheetPNG(req QRRequest, cols, rows int) (string, error) {
	content, err := formatContent(req.Type, req.Fields)
	if err != nil {
		return "", err
	}

	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("QR generation failed: %v", err)
	}

	qrPNG, err := qr.PNG(200)
	if err != nil {
		return "", fmt.Errorf("PNG encoding failed: %v", err)
	}

	qrImg, _, err := image.Decode(bytes.NewReader(qrPNG))
	if err != nil {
		return "", fmt.Errorf("image decode failed: %v", err)
	}

	qrSize := qrImg.Bounds().Dx()
	padding := 20
	margin := 30
	totalW := cols*qrSize + (cols-1)*padding + 2*margin
	totalH := rows*qrSize + (rows-1)*padding + 2*margin

	sheet := image.NewRGBA(image.Rect(0, 0, totalW, totalH))
	draw.Draw(sheet, sheet.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			x := margin + c*(qrSize+padding)
			y := margin + r*(qrSize+padding)
			draw.Draw(sheet, image.Rect(x, y, x+qrSize, y+qrSize), qrImg, image.Point{}, draw.Src)
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, sheet); err != nil {
		return "", fmt.Errorf("PNG encode failed: %v", err)
	}

	filePath, err := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		DefaultFilename: "qr-sheet.png",
		Title:           "Save QR Sheet",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "PNG Image", Pattern: "*.png"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("dialog failed: %v", err)
	}
	if filePath == "" {
		return "", nil
	}

	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("write failed: %v", err)
	}

	return filePath, nil
}

func (a *App) CopyContentToClipboard(text string) error {
	return wailsRuntime.ClipboardSetText(a.ctx, text)
}
