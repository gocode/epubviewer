package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
	"html/template"
	"io"
)

type epub struct {
	fs vfs.FileSystem
}

func New(name string) (*epub, error) {
	rc, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}
	//defer rc.Close()

	return &epub{fs: zipfs.New(rc, name)}, nil
}

func (e *epub) WriteToc(w io.Writer) error {

	buf, err := vfs.ReadFile(e.fs, "/toc.ncx")
	if err != nil {
		return err
	}

	v := ncx{}

	err = xml.Unmarshal(buf, &v)
	if err != nil {
		return err
	}

	t, err := template.ParseFiles("ncx.template")
	if err != nil {
		return err
	}

	err = t.Execute(w, v)
	if err != nil {
		return err
	}

	return nil
}

func (e *epub) WriteSpine(w io.Writer) error {
	buf, err := vfs.ReadFile(e.fs, "/content.opf")
	if err != nil {
		return err
	}

	v := OPF{}

	err = xml.Unmarshal(buf, &v)
	if err != nil {
		return err
	}

	for _, t := range tocFromSpine(v.Spine, v.Manifest) {
		fmt.Fprintln(w, "<a href="+t.Href+">"+t.Text+"</a><br/>")
	}

	return nil
}

func (e *epub) WriteFile(w io.Writer, path string) error {
	buf, err := vfs.ReadFile(e.fs, path)
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
