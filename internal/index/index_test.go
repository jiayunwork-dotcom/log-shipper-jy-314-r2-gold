package index

import (
	"testing"

	"log-shipper/internal/collect"
)

func sample() []collect.Record {
	return []collect.Record{
		{Name: "app.log", Lines: []string{"info: started", "error: boom occurred", "done"}},
		{Name: "sys.log", Lines: []string{"warning: low disk", "error: timeout"}},
	}
}

func TestBuildIndex_SearchKeyword(t *testing.T) {
	idx := Build(sample())
	hits := idx.Search("error")
	if len(hits) != 2 {
		t.Fatalf("expected 2 hits for 'error', got %d", len(hits))
	}
	for _, h := range hits {
		if h.Line < 1 {
			t.Fatalf("line must be 1-based, got %d", h.Line)
		}
	}
}

func TestIndex_SearchMissing(t *testing.T) {
	idx := Build(sample())
	if got := idx.Search("zzz-nope"); got != nil {
		t.Fatalf("missing keyword should return nil, got %v", got)
	}
	if idx.Count("zzz-nope") != 0 {
		t.Fatalf("count for missing should be 0")
	}
}

func TestIndex_Count(t *testing.T) {
	idx := Build(sample())
	if c := idx.Count("error"); c != 2 {
		t.Fatalf("expected count 2 for 'error', got %d", c)
	}
}

func TestSearchHitsIndependent(t *testing.T) {
	idx := Build(sample())
	hits := idx.Search("error")
	if len(hits) != 2 {
		t.Fatalf("expected 2 hits, got %d", len(hits))
	}
	hits[0].Text = "mutated"
	hits[0].File = "x"
	again := idx.Search("error")
	if len(again) != 2 {
		t.Fatalf("second Search expected 2 hits, got %d", len(again))
	}
	if again[0].Text == "mutated" || again[0].File == "x" {
		t.Fatalf("mutating first Search result changed later Search: %+v", again[0])
	}
	if again[0].Text != "error: boom occurred" || again[0].File != "app.log" {
		t.Fatalf("second Search = %+v", again[0])
	}
}

func TestBuild_LineNumbersOneBased(t *testing.T) {
	idx := Build(sample())
	started := idx.Search("started")
	if len(started) != 1 || started[0].Line != 1 {
		t.Fatalf("started should be line 1, got %+v", started)
	}
	done := idx.Search("done")
	if len(done) != 1 || done[0].Line != 3 {
		t.Fatalf("done should be line 3, got %+v", done)
	}
}
