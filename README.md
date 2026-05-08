# QR Aedificator

Multi-purpose desktop QR code generator built with [Wails v2](https://wails.io/) (Go backend + HTML/CSS/JS frontend).

## Features

- **10 QR content types**: Text, URL, WiFi, vCard, Email, SMS, Phone, Geo Location, Calendar, Bitcoin
- **Error correction**: L / M / Q / H
- **Custom sizes**: 128px to 1024px
- **Colors**: Customizable foreground and background colors
- **Logo overlay**: Embed a logo image in the center of the QR code
- **QR Decoding**: Decode QR codes from image files or drag-and-drop
- **SVG Export**: Save QR code as vector SVG
- **PDF Export**: Save QR sheet as PDF document
- **Batch CSV Import**: Generate multiple QR codes from a CSV file
- **History**: Automatic history with search, load, regenerate, and delete
- **Copy to clipboard**: Copy encoded content

## Requirements

- Go 1.18+
- Node.js 16+
- Wails v2 CLI
- Linux: webkit2gtk-4.1 (Arch: `webkit2gtk-4.1`), `gcc`, `pkg-config`

## Build

```bash
# Install frontend dependencies
cd frontend && npm install && cd ..

# Build with Wails
wails build -tags webkit2_41 -o qr-aedificator

# Binary is in build/bin/
```

## Project Structure

```
├── app.go              — App struct, startup/shutdown
├── main.go             — Wails entry point
├── qrgen.go            — QR generation, formatting, SVG/PNG export
├── qrdecode.go         — QR decoding with gozxing
├── qrlogo.go           — Logo image picker and overlay
├── qrbatch.go          — CSV parsing and batch generation
├── qrpdf.go            — PDF sheet export
├── history.go          — SQLite history (modernc.org/sqlite)
├── frontend/
│   ├── index.html      — Main UI
│   ├── src/
│   │   ├── style.css   — Dark theme styles
│   │   └── templates.js — Form builders for each QR type
│   └── wailsjs/        — Auto-generated Wails bindings
└── build/              — Build artifacts
```

## Dependencies

- [github.com/wailsapp/wails/v2](https://github.com/wailsapp/wails/v2) — Desktop framework
- [github.com/skip2/go-qrcode](https://github.com/skip2/go-qrcode) — QR code generation
- [github.com/makiuchi-d/gozxing](https://github.com/makiuchi-d/gozxing) — QR decoding
- [modernc.org/sqlite](https://modernc.org/sqlite) — Pure Go SQLite for history
- [github.com/jung-kurt/gofpdf](https://github.com/jung-kurt/gofpdf) — PDF export
