// Package collect 读取目录中的日志文件，整理为可归档的记录。
package collect

import (
	"os"
	"path/filepath"
	"strings"
)

// Record 是一条日志记录（一个源文件 + 其按行切分的内容）。
type Record struct {
	Name  string
	Lines []string
}

// Collect 读取 srcDir 下所有 .log 文件，返回记录切片。
// 目录不可读时返回错误；目录下没有 .log 文件则返回空切片（不视为错误）。
func Collect(srcDir string) ([]Record, error) {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}
	var recs []Record
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(e.Name()), ".log") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(srcDir, e.Name()))
		if err != nil {
			return nil, err
		}
		recs = append(recs, Record{
			Name:  e.Name(),
			Lines: strings.Split(string(data), "\n"),
		})
	}
	return recs, nil
}
