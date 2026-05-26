package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreSaveAndList(t *testing.T) {
	dir := t.TempDir()
	st, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	rec, err := st.SavePayload(FormatJSON, "test.json", []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if rec.ID == "" || rec.Filename == "" {
		t.Fatalf("record: %+v", rec)
	}
	list := st.List()
	if len(list) != 1 || list[0].ID != rec.ID {
		t.Fatalf("list: %+v", list)
	}
	p, got, err := st.GetPath(rec.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != rec.ID {
		t.Fatal("mismatch")
	}
	b, err := os.ReadFile(p)
	if err != nil || string(b) != "{}" {
		t.Fatalf("file: %s err=%v", b, err)
	}
}

func TestGetPathMissing(t *testing.T) {
	dir := t.TempDir()
	st, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = st.GetPath("nope")
	if err != os.ErrNotExist {
		t.Fatalf("got %v want ErrNotExist", err)
	}
}

func TestSavePayloadSanitizeName(t *testing.T) {
	dir := t.TempDir()
	st, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.SavePayload(FormatJSON, `bad<>name.json`, []byte(`[]`))
	if err != nil {
		t.Fatal(err)
	}
	files, _ := os.ReadDir(filepath.Join(dir, "received"))
	if len(files) != 1 {
		t.Fatalf("files: %d", len(files))
	}
}
