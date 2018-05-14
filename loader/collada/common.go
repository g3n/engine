// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collada

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//
// Source
//
type Source struct {
	Id              string      // Source id
	Name            string      // Source name
	ArrayElement    interface{} // Array element (FloatArray|others)
	TechniqueCommon struct {
		Accessor Accessor
	}
}

// Dump prints out information about the Source
func (s *Source) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sSource id:%s name:%s\n", sIndent(indent), s.Id, s.Name)
	ind := indent + step
	switch at := s.ArrayElement.(type) {
	case *FloatArray:
		at.Dump(out, ind)
	case *NameArray:
		at.Dump(out, ind)
	}
	fmt.Fprintf(out, "%sTechniqueCommon\n", sIndent(ind))
	s.TechniqueCommon.Accessor.Dump(out, ind+3)
}

//
// NameArray
//
type NameArray struct {
	Id    string
	Name  string
	Count int
	Data  []string
}

// Dump prints out information about the NameArray
func (na *NameArray) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sNameArray id:%s count:%d\n", sIndent(indent), na.Id, na.Count)
	ind := indent + step
	fmt.Fprintf(out, "%sData(%d):%s\n", sIndent(ind), len(na.Data), na.Data)
}

//
// FloatArray
//
type FloatArray struct {
	Id    string
	Count int
	Data  []float32
}

// Dump prints out information about the FloatArray
func (fa *FloatArray) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sFloatArray id:%s count:%d\n", sIndent(indent), fa.Id, fa.Count)
	ind := indent + step
	fmt.Fprintf(out, "%sData(%d):%s\n", sIndent(ind), len(fa.Data), f32sToString(fa.Data, 20))
}

//
// Accessor
//
type Accessor struct {
	Source string
	Count  int
	Stride int
	Params []Param
}

// Dump prints out information about the Accessor
func (ac *Accessor) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sAccessor source:%s count:%d stride:%d\n",
		sIndent(indent), ac.Source, ac.Count, ac.Stride)
	ind := indent + step
	for _, p := range ac.Params {
		p.Dump(out, ind)
	}
}

//
// Param for <bind_material> and <accessor>
//
type Param struct {
	Name string
	Type string
}

// Dump prints out information about the Param
func (p *Param) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sParam name:%s type:%s\n", sIndent(indent), p.Name, p.Type)
}

// decSource decodes the source from the specified mesh
func (d *Decoder) decSource(start xml.StartElement) (*Source, error) {

	// Create source and adds it to the mesh
	source := new(Source)
	source.Id = findAttrib(start, "id").Value
	source.Name = findAttrib(start, "name").Value

	// Decodes source children
	for {
		// Get next child
		child, data, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return source, err
		}
		if child.Name.Local == "float_array" {
			err = d.decFloatArray(child, data, source)
			if err != nil {
				return nil, err
			}
			continue
		}
		if child.Name.Local == "Name_array" {
			err = d.decNameArray(child, data, source)
			if err != nil {
				return nil, err
			}
			continue
		}
		// Decodes technique_common which should contain an Acessor
		if child.Name.Local == "technique_common" {
			err = d.decSourceTechniqueCommon(child, source)
			if err != nil {
				return nil, err
			}
			continue
		}
	}
}

// decSource decodes the float array from the specified source
func (d *Decoder) decFloatArray(start xml.StartElement, data []byte, source *Source) error {

	// Create float array and associates it to the parent source
	farray := &FloatArray{}
	farray.Id = findAttrib(start, "id").Value
	farray.Count, _ = strconv.Atoi(findAttrib(start, "count").Value)
	source.ArrayElement = farray

	// Allocates memory for array
	farray.Data = make([]float32, farray.Count, farray.Count)

	// Reads the numbers from the data
	err := decFloat32Sequence(data, farray.Data)
	if err != nil {
		return err
	}
	return nil
}

func (d *Decoder) decNameArray(start xml.StartElement, data []byte, source *Source) error {

	narray := new(NameArray)
	narray.Id = findAttrib(start, "id").Value
	narray.Count, _ = strconv.Atoi(findAttrib(start, "count").Value)
	source.ArrayElement = narray

	// Allocates memory for array
	narray.Data = make([]string, narray.Count, narray.Count)

	// Reads the strings from the data
	err := decStringSequence(data, narray.Data)
	if err != nil {
		return err
	}
	return nil
}

func (d *Decoder) decSourceTechniqueCommon(start xml.StartElement, source *Source) error {

	// Decodes source technique common children
	for {
		// Get next child
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		//log.Debug("decSourceTechniqueCommon(%s): %s", start.Name.Local, child.Name.Local)
		if child.Name.Local == "accessor" {
			err = d.decAcessor(child, source)
			if err != nil {
				return err
			}
			continue
		}
	}
}

// decAcessore decodes the acessor from the specified source
func (d *Decoder) decAcessor(start xml.StartElement, source *Source) error {

	// Sets accessor fields
	source.TechniqueCommon.Accessor.Source = findAttrib(start, "source").Value
	source.TechniqueCommon.Accessor.Count, _ = strconv.Atoi(findAttrib(start, "count").Value)
	source.TechniqueCommon.Accessor.Stride, _ = strconv.Atoi(findAttrib(start, "stride").Value)

	// Decodes accessor children
	for {
		// Get next child
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		// param
		if child.Name.Local == "param" {
			err = d.decParam(child, &source.TechniqueCommon.Accessor)
			if err != nil {
				return err
			}
		}
	}
}

func (d *Decoder) decParam(start xml.StartElement, accessor *Accessor) error {

	p := Param{}
	p.Name = findAttrib(start, "name").Value
	p.Type = findAttrib(start, "type").Value
	accessor.Params = append(accessor.Params, p)
	return nil
}

func (d *Decoder) decNextChild(parent xml.StartElement) (xml.StartElement, []byte, error) {

	for {
		var tok interface{}
		var err error
		// Reads next token
		if d.lastToken == nil {
			tok, err = d.xmldec.Token()
			if err != nil {
				return xml.StartElement{}, nil, err
			}
		} else {
			tok = d.lastToken
			d.lastToken = nil
		}
		// Checks if it is the end element of this parent
		el, ok := tok.(xml.EndElement)
		if ok {
			if el.Name.Local == parent.Name.Local {
				return xml.StartElement{}, nil, nil
			}
			continue
		}
		// Checks if it is a start element
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		// Get this start element optional char data (should be next token)
		tok, err = d.xmldec.Token()
		if err != nil {
			return xml.StartElement{}, nil, err
		}
		// If token read is CharData, return the start element and its CharData
		cdata, ok := tok.(xml.CharData)
		if ok {
			return start, cdata, nil
		}
		// Token read was not CharData and was not processed
		// Save it into "lastToken" to be processed at the next call
		d.lastToken = tok
		return start, nil, nil
	}
}

func findAttrib(s xml.StartElement, name string) xml.Attr {

	for _, attr := range s.Attr {
		if attr.Name.Local == name {
			return attr
		}
	}
	return xml.Attr{}
}

const tokenSep string = " \r\n\t"

type bytesReader struct {
	pos    int
	source []byte
}

func (br *bytesReader) Init(source []byte) {

	br.pos = 0
	br.source = source
}

func (br *bytesReader) TokenNext() []byte {

	// Skip leading separators
	for br.pos < len(br.source) {
		if bytes.IndexByte([]byte(tokenSep), br.source[br.pos]) < 0 {
			break
		}
		br.pos++
	}
	if br.pos >= len(br.source) {
		return nil
	}

	// Advance till the end of the token
	start := br.pos
	for br.pos < len(br.source) {
		if bytes.IndexByte([]byte(tokenSep), br.source[br.pos]) >= 0 {
			break
		}
		br.pos++
	}
	res := br.source[start:br.pos]
	if len(res) == 0 {
		return nil
	}
	return res
}

const step = 3

func sIndent(indent int) string {

	return strings.Repeat(" ", indent)
}

// decFloat32Sequence receives a byte slice with float numbers separated
// by spaces and a preallocated destination slice.
// It reads numbers from the source byte slice, converts them to float32 and
// stores in the destination array.
func decFloat32Sequence(cdata []byte, dest []float32) error {

	var br bytesReader
	br.Init(cdata)
	idx := 0
	for {
		tok := br.TokenNext()
		if tok == nil {
			break
		}
		if idx >= len(dest) {
			return fmt.Errorf("To much float array data")
		}
		v, err := strconv.ParseFloat(string(tok), 32)
		if err != nil {
			return err
		}
		dest[idx] = float32(v)
		idx++
	}
	if idx < len(dest)-1 {
		return fmt.Errorf("Expected %d floats, got %d", len(dest), idx)
	}
	return nil
}

// decStringSequence receives a byte slice with strings separated
// by spaces and a preallocated destination slice.
// It reads strings from the source byte slice and
// stores in the destination array.
func decStringSequence(cdata []byte, dest []string) error {

	var br bytesReader
	br.Init(cdata)
	idx := 0
	for {
		tok := br.TokenNext()
		if tok == nil {
			break
		}
		if idx >= len(dest) {
			return fmt.Errorf("To many string array data")
		}
		dest[idx] = string(tok)
		idx++
	}
	if idx < len(dest)-1 {
		return fmt.Errorf("Expected %d strings, got %d", len(dest), idx)
	}
	return nil
}

func f32sToString(a []float32, max int) string {

	parts := []string{"["}
	if len(a) > max {
		for i := 0; i < max/2; i++ {
			parts = append(parts, strconv.FormatFloat(float64(a[i]), 'f', -1, 32))
		}
		parts = append(parts, " ... ")
		for i := len(a) - max/2; i < len(a); i++ {
			parts = append(parts, strconv.FormatFloat(float64(a[i]), 'f', -1, 32))
		}
	} else {
		for i := 0; i < len(a); i++ {
			parts = append(parts, strconv.FormatFloat(float64(a[i]), 'f', -1, 32))
		}
	}
	parts = append(parts, "]")
	return strings.Join(parts, " ")
}

func intsToString(a []int, max int) string {

	parts := []string{"["}
	if len(a) > max {
		for i := 0; i < max/2; i++ {
			parts = append(parts, strconv.FormatInt(int64(a[i]), 10))
		}
		parts = append(parts, " ... ")
		for i := len(a) - max/2; i < len(a); i++ {
			parts = append(parts, strconv.FormatInt(int64(a[i]), 10))
		}
	} else {
		for i := 0; i < len(a); i++ {
			parts = append(parts, strconv.FormatInt(int64(a[i]), 10))
		}
	}
	parts = append(parts, "]")
	return strings.Join(parts, " ")
}
