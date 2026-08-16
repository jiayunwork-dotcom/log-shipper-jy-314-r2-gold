// Package bundle 把日志记录打包为 tar.gz 归档，并可重新读回。
package bundle

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"strings"

	"log-shipper/internal/collect"
)

// Write 将记录压缩并写入一个 tar.gz 归档文件。
func Write(path string, recs []collect.Record) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	for _, r := range recs {
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		if _, err := zw.Write([]byte(strings.Join(r.Lines, "\n"))); err != nil {
			return err
		}
		if err := zw.Close(); err != nil {
			return err
		}
		hdr := &tar.Header{
			Name: r.Name + ".gz",
			Mode: 0644,
			Size: int64(buf.Len()),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tw.Write(buf.Bytes()); err != nil {
			return err
		}
	}
	if err := tw.Close(); err != nil {
		return err
	}
	return gw.Close()
}

// Read 从 tar.gz 归档读回记录。
func Read(path string) ([]collect.Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	var recs []collect.Record
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if !strings.HasSuffix(hdr.Name, ".gz") {
			continue
		}
		zr, err := gzip.NewReader(tr)
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(zr)
		zr.Close()
		if err != nil {
			return nil, err
		}
		name := strings.TrimSuffix(hdr.Name, ".gz")
		recs = append(recs, collect.Record{
			Name:  name,
			Lines: strings.Split(string(data), "\n"),
		})
	}
	return recs, nil
}
