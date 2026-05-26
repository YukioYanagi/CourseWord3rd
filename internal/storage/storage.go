package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Format string

const (
	FormatJSON Format = "json"
	FormatXML  Format = "xml"
	FormatSOAP Format = "soap"
)

type Record struct {
	ID        string    `json:"id"`
	Format    Format    `json:"format"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	mu      sync.Mutex
	root    string
	indexMu sync.RWMutex
	records []Record
}

func New(root string) (*Store, error) {
	received := filepath.Join(root, "received")
	if err := os.MkdirAll(received, 0o750); err != nil {
		return nil, err
	}
	s := &Store{root: root, records: nil}
	if err := s.loadIndex(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

func (s *Store) indexPath() string {
	return filepath.Join(s.root, "index.json")
}

func (s *Store) loadIndex() error {
	b, err := os.ReadFile(s.indexPath())
	if err != nil {
		return err
	}
	var recs []Record
	if err := json.Unmarshal(b, &recs); err != nil {
		return err
	}
	s.indexMu.Lock()
	s.records = recs
	s.indexMu.Unlock()
	return nil
}

func (s *Store) persistIndex() error {
	s.indexMu.RLock()
	data, err := json.MarshalIndent(s.records, "", "  ")
	s.indexMu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(s.indexPath(), data, 0o600)
}

func (s *Store) SavePayload(format Format, originalName string, body []byte) (Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	safe := strings.Map(func(r rune) rune {
		if r < 32 || strings.ContainsRune(`<>:"/\|?*`, r) {
			return '_'
		}
		return r
	}, originalName)
	if safe == "" || safe == "_" {
		safe = "payload"
	}
	ext := string(format)
	if ext == "soap" {
		ext = "xml"
	}
	filename := fmt.Sprintf("%s_%s.%s", id, strings.TrimSuffix(safe, filepath.Ext(safe)), ext)
	path := filepath.Join(s.root, "received", filename)
	if err := os.WriteFile(path, body, 0o600); err != nil {
		return Record{}, err
	}
	st, _ := os.Stat(path)
	rec := Record{
		ID:        id,
		Format:    format,
		Filename:  filename,
		Size:      st.Size(),
		CreatedAt: time.Now().UTC(),
	}
	s.indexMu.Lock()
	s.records = append([]Record{rec}, s.records...)
	s.indexMu.Unlock()
	if err := s.persistIndex(); err != nil {
		return Record{}, err
	}
	return rec, nil
}

func (s *Store) List() []Record {
	s.indexMu.RLock()
	defer s.indexMu.RUnlock()
	out := make([]Record, len(s.records))
	copy(out, s.records)
	return out
}

func (s *Store) GetPath(id string) (string, Record, error) {
	s.indexMu.RLock()
	defer s.indexMu.RUnlock()
	for _, r := range s.records {
		if r.ID == id {
			p := filepath.Join(s.root, "received", r.Filename)
			if _, err := os.Stat(p); err != nil {
				return "", Record{}, err
			}
			return p, r, nil
		}
	}
	return "", Record{}, os.ErrNotExist
}

// AssertPathInsideReceived returns nil if absPath is a regular file under this store's received directory.
func (s *Store) AssertPathInsideReceived(absPath string) error {
	base, err := filepath.Abs(filepath.Join(s.root, "received"))
	if err != nil {
		return err
	}
	target, err := filepath.Abs(absPath)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return errors.New("path outside received directory")
	}
	return nil
}
