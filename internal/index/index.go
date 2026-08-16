// Package index 构建日志的倒排索引，支持按关键词检索。
package index

import (
	"strings"

	"log-shipper/internal/collect"
)

// Hit 是一条检索命中。
type Hit struct {
	File string
	Line int // 1-based 行号
	Text string
}

// Index 是倒排索引：词 -> 命中列表。
type Index struct {
	term map[string][]Hit
}

// Build 从记录集合构建倒排索引。
func Build(recs []collect.Record) *Index {
	idx := &Index{term: map[string][]Hit{}}
	for _, r := range recs {
		for i, line := range r.Lines {
			for _, tok := range tokenize(line) {
				idx.term[tok] = append(idx.term[tok], Hit{
					File: r.Name, Line: i + 1, Text: line,
				})
			}
		}
	}
	return idx
}

// Search 按关键词（取首个分词）检索，返回命中列表。
func (idx *Index) Search(keyword string) []Hit {
	if idx == nil || idx.term == nil {
		return nil
	}
	toks := tokenize(keyword)
	if len(toks) == 0 {
		return nil
	}
	hits := idx.term[toks[0]]
	if hits == nil {
		return nil
	}
	out := make([]Hit, len(hits))
	copy(out, hits)
	return out
}

// Count 返回关键词命中的条数。
func (idx *Index) Count(keyword string) int {
	return len(idx.Search(keyword))
}

func tokenize(line string) []string {
	line = strings.ToLower(line)
	var toks []string
	var cur strings.Builder
	flush := func() {
		if cur.Len() > 0 {
			toks = append(toks, cur.String())
			cur.Reset()
		}
	}
	for _, r := range line {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			cur.WriteRune(r)
		} else {
			flush()
		}
	}
	flush()
	return toks
}
