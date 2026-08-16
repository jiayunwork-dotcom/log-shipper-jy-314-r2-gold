package collect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollect_OK(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "app.log"), []byte("start\nerror: boom\nend"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sys.log"), []byte("ok\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	recs, err := Collect(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 2 {
		t.Fatalf("expected 2 records, got %d", len(recs))
	}
}

func TestCollect_MissingDir(t *testing.T) {
	if _, err := Collect("/no/such/dir/xyz"); err == nil {
		t.Fatalf("expected error for missing directory")
	}
}

func TestCollect_MissingDirErrors(t *testing.T) {
	recs, err := Collect(filepath.Join(t.TempDir(), "no-such-logs"))
	if err == nil {
		t.Fatalf("missing directory should return an error, got recs=%v", recs)
	}
}

func TestCollect_UppercaseExt(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "APP.LOG"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	recs, err := Collect(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 1 || recs[0].Name != "APP.LOG" {
		t.Fatalf("expected APP.LOG, got %+v", recs)
	}
}
