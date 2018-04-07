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
// Library Animations
//
type LibraryAnimations struct {
	Id        string
	Name      string
	Asset     *Asset
	Animation []*Animation
}

// Dump prints out information about the LibraryAnimations
func (la *LibraryAnimations) Dump(out io.Writer, indent int) {

	if la == nil {
		return
	}
	fmt.Fprintf(out, "%sLibraryAnimations id:%s name:%s\n", sIndent(indent), la.Id, la.Name)
	for _, an := range la.Animation {
		an.Dump(out, indent+step)
	}
}

//
// Animation
//
type Animation struct {
	Id        string
	Name      string
	Animation []*Animation
	Source    []*Source
	Sampler   []*Sampler
	Channel   []*Channel
}

// Dump prints out information about the Animation
func (an *Animation) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sAnimation id:%s name:%s\n", sIndent(indent), an.Id, an.Name)
	ind := indent + step
	for _, source := range an.Source {
		source.Dump(out, ind)
	}
	for _, sampler := range an.Sampler {
		sampler.Dump(out, ind)
	}
	for _, channel := range an.Channel {
		channel.Dump(out, ind)
	}
	for _, child := range an.Animation {
		child.Dump(out, ind)
	}
}

//
// Sampler
//
type Sampler struct {
	Id    string
	Input []Input // One or more

}

// Dump prints out information about the Sampler
func (sp *Sampler) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sSampler id:%s\n", sIndent(indent), sp.Id)
	ind := indent + step
	for _, inp := range sp.Input {
		inp.Dump(out, ind)
	}
}

//
// Channel
//
type Channel struct {
	Source string
	Target string
}

// Dump prints out information about the Channel
func (ch *Channel) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sChannel source:%s target:%s\n", sIndent(indent), ch.Source, ch.Target)
}

func (d *Decoder) decLibraryAnimations(start xml.StartElement, dom *Collada) error {

	la := new(LibraryAnimations)
	dom.LibraryAnimations = la
	la.Id = findAttrib(start, "id").Value
	la.Name = findAttrib(start, "name").Value

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "animation" {
			err := d.decAnimation(child, la)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decAnimation(start xml.StartElement, la *LibraryAnimations) error {

	anim := new(Animation)
	la.Animation = append(la.Animation, anim)
	anim.Id = findAttrib(start, "id").Value
	anim.Name = findAttrib(start, "name").Value

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "source" {
			source, err := d.decSource(child)
			if err != nil {
				return err
			}
			anim.Source = append(anim.Source, source)
			continue
		}
		if child.Name.Local == "sampler" {
			err = d.decSampler(child, anim)
			if err != nil {
				return nil
			}
			continue
		}
		if child.Name.Local == "channel" {
			err = d.decChannel(child, anim)
			if err != nil {
				return nil
			}
			continue
		}
	}
}

func (d *Decoder) decSampler(start xml.StartElement, anim *Animation) error {

	sp := new(Sampler)
	anim.Sampler = append(anim.Sampler, sp)
	sp.Id = findAttrib(start, "id").Value

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "input" {
			inp, err := d.decInput(child)
			if err != nil {
				return err
			}
			sp.Input = append(sp.Input, inp)
			continue
		}
	}
}

func (d *Decoder) decChannel(start xml.StartElement, anim *Animation) error {

	ch := new(Channel)
	ch.Source = findAttrib(start, "source").Value
	ch.Target = findAttrib(start, "target").Value
	anim.Channel = append(anim.Channel, ch)
	return nil
}
