package bundle

import (
	"path/filepath"
	"testing"

	"log-shipper/internal/collect"
)

func TestBundle_RoundTrip(t *testing.T) {
	recs := []collect.Record{
		{Name: "app.log", Lines: []string{"line one", "line two"}},
		{Name: "sys.log", Lines: []string{"only line"}},
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "logs.tar.gz")
	if err := Write(path, recs); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := Read(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(got) != len(recs) {
		t.Fatalf("expected %d records, got %d", len(recs), len(got))
	}
	if got[0].Name != "app.log" || len(got[0].Lines) != 2 {
		t.Fatalf("first record mismatch: %+v", got[0])
	}
	if got[1].Lines[0] != "only line" {
		t.Fatalf("second record content mismatch: %+v", got[1])
	}
}

func TestBundle_ReadMissing(t *testing.T) {
	if _, err := Read("/no/such/bundle.tar.gz"); err == nil {
		t.Fatalf("expected error reading missing bundle")
	}
}
