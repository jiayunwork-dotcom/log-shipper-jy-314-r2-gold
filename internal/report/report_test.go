package report

import (
	"bytes"
	"testing"

	"log-shipper/internal/index"
)

func TestRender_NoMatches(t *testing.T) {
	var buf bytes.Buffer
	Render(nil, &buf)
	if got := buf.String(); got != "no matches\n" {
		t.Fatalf("nil hits: got %q", got)
	}
	buf.Reset()
	Render([]index.Hit{}, &buf)
	if got := buf.String(); got != "no matches\n" {
		t.Fatalf("empty hits: got %q", got)
	}
}
