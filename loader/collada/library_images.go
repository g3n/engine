// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collada

import (
	"encoding/xml"
	"fmt"
	"io"
)

//
// LibraryImages
//
type LibraryImages struct {
	Id    string
	Name  string
	Asset *Asset
	Image []*Image
}

// Dump prints out information about the LibraryImages
func (li *LibraryImages) Dump(out io.Writer, indent int) {

	if li == nil {
		return
	}
	fmt.Fprintf(out, "%sLibraryImages id:%s name:%s\n", sIndent(indent), li.Id, li.Name)
	for _, img := range li.Image {
		img.Dump(out, indent+step)
	}
}

//
// Image
//
type Image struct {
	Id          string
	Name        string
	Format      string
	Height      uint
	Width       uint
	Depth       uint
	ImageSource interface{}
}

// Dump prints out information about the Image
func (img *Image) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sImage id:%s name:%s\n", sIndent(indent), img.Id, img.Name)
	ind := indent + step
	switch is := img.ImageSource.(type) {
	case InitFrom:
		is.Dump(out, ind)
	}
}

//
// InitFrom
//
type InitFrom struct {
	Uri string
}

// Dump prints out information about the InitFrom
func (initf *InitFrom) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sInitFrom:%s\n", sIndent(indent), initf.Uri)
}

func (d *Decoder) decLibraryImages(start xml.StartElement, dom *Collada) error {

	li := new(LibraryImages)
	dom.LibraryImages = li
	li.Id = findAttrib(start, "id").Value
	li.Name = findAttrib(start, "name").Value

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "image" {
			err := d.decImage(child, li)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decImage(start xml.StartElement, li *LibraryImages) error {

	img := new(Image)
	img.Id = findAttrib(start, "id").Value
	img.Name = findAttrib(start, "name").Value
	li.Image = append(li.Image, img)

	for {
		child, data, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "init_from" {
			err := d.decImageSource(child, data, img)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decImageSource(start xml.StartElement, cdata []byte, img *Image) error {

	img.ImageSource = InitFrom{string(cdata)}
	return nil
}
