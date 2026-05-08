package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type HistoryEntry struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Fields    string `json:"fields"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

func (a *App) initDB() error {
	db, err := sql.Open("sqlite", "qr-history.db")
	if err != nil {
		return fmt.Errorf("open db failed: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		qr_type TEXT NOT NULL,
		fields TEXT NOT NULL DEFAULT '{}',
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		db.Close()
		return fmt.Errorf("create table failed: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_history_created ON history(created_at DESC)`)
	if err != nil {
		db.Close()
		return fmt.Errorf("create index failed: %v", err)
	}

	a.db = db
	return nil
}

func (a *App) AddToHistory(req QRRequest, content string) error {
	if a.db == nil {
		return nil
	}

	fieldsJSON, err := json.Marshal(req.Fields)
	if err != nil {
		return fmt.Errorf("json marshal failed: %v", err)
	}

	_, err = a.db.Exec(
		"INSERT INTO history (qr_type, fields, content, created_at) VALUES (?, ?, ?, ?)",
		req.Type, string(fieldsJSON), content, time.Now().UTC().Format(time.RFC3339),
	)
	return err
}

func (a *App) GetHistory(search string, limit, offset int) ([]HistoryEntry, error) {
	if a.db == nil {
		return []HistoryEntry{}, nil
	}

	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var rows *sql.Rows
	var err error

	if search != "" {
		rows, err = a.db.Query(
			"SELECT id, qr_type, fields, content, created_at FROM history WHERE content LIKE ? OR qr_type LIKE ? ORDER BY created_at DESC LIMIT ? OFFSET ?",
			"%"+search+"%", "%"+search+"%", limit, offset,
		)
	} else {
		rows, err = a.db.Query(
			"SELECT id, qr_type, fields, content, created_at FROM history ORDER BY created_at DESC LIMIT ? OFFSET ?",
			limit, offset,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var entries []HistoryEntry
	for rows.Next() {
		var e HistoryEntry
		var createdAt time.Time
		if err := rows.Scan(&e.ID, &e.Type, &e.Fields, &e.Content, &createdAt); err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		e.CreatedAt = createdAt.Format(time.RFC3339)
		entries = append(entries, e)
	}

	if entries == nil {
		entries = []HistoryEntry{}
	}

	return entries, nil
}

func (a *App) ClearHistory() error {
	if a.db == nil {
		return nil
	}
	_, err := a.db.Exec("DELETE FROM history")
	return err
}

func (a *App) DeleteHistoryEntry(id int64) error {
	if a.db == nil {
		return nil
	}
	_, err := a.db.Exec("DELETE FROM history WHERE id = ?", id)
	return err
}
