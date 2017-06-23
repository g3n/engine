package stats

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
)

type Stats struct {
	*gui.Table          // embedded table
	fields     []*field // array of fields
}

type field struct {
	id  string
	row int
}

func NewStats(width, height float32) *Stats {

	s := new(Stats)
	t, err := gui.NewTable(width, height, []gui.TableColumn{
		{Id: "f", Header: "Stat", Width: 100, Minwidth: 32, Align: gui.AlignRight, Format: "%s", Resize: false, Expand: 0},
		{Id: "v", Header: "Value", Width: 100, Minwidth: 32, Align: gui.AlignRight, Format: "%d", Resize: false, Expand: 0},
	})
	if err != nil {
		panic(err)
	}
	s.Table = t
	s.addRow("shaders", "Shaders:")
	s.addRow("vaos", "Vaos:")
	s.addRow("vbos", "Vbos:")
	s.addRow("textures", "Textures:")
	s.addRow("ccalls", "CGO calls/frame:")
	return s
}

func (s *Stats) Update(gs *gls.GLS) {

	var stats gls.Stats
	gs.Stats(&stats)
	for i := 0; i < len(s.fields); i++ {
		f := s.fields[i]
		switch f.id {
		case "shaders":
			s.Table.SetCell(f.row, "v", stats.Shaders)
		case "vaos":
			s.Table.SetCell(f.row, "v", stats.Vaos)
		case "vbos":
			s.Table.SetCell(f.row, "v", stats.Vbos)
		case "textures":
			s.Table.SetCell(f.row, "v", stats.Textures)
		}
	}
}

func (s *Stats) addRow(id, label string) {

	f := new(field)
	f.id = id
	f.row = s.Table.RowCount()
	s.Table.AddRow(map[string]interface{}{"f": label, "v": 0})
	s.fields = append(s.fields, f)
}

//type field struct {
//	id     string
//	label  *gui.Label
//	value  *gui.Label
//	align  gui.Align
//	format string
//}
//
//func NewStats() *Stats {
//
//	s := new(Stats)
//	s.Panel.Initialize(0, 0)
//	s.addField("shaders", "Shaders", "%d", gui.AlignRight)
//	s.addField("vaos", "VAOs", "%d", gui.AlignRight)
//	s.addField("vbos", "VBOs", "%d", gui.AlignRight)
//	s.addField("textures", "Textures", "%d", gui.AlignRight)
//	s.recalc()
//	return s
//}
//
//func (s *Stats) Update(gs *gls.GLS) {
//
//	var stats gls.Stats
//	gs.Stats(&stats)
//	for i := 0; i < len(s.fields); i++ {
//		f := s.fields[i]
//		switch f.id {
//		case "shaders":
//			f.value.SetText(fmt.Sprintf(f.format, stats.Shaders))
//		case "vaos":
//			f.value.SetText(fmt.Sprintf(f.format, stats.Vaos))
//		case "vbos":
//			f.value.SetText(fmt.Sprintf(f.format, stats.Vbos))
//		case "textures":
//			f.value.SetText(fmt.Sprintf(f.format, stats.Textures))
//		}
//	}
//}
//
//func (s *Stats) addField(id, label string, format string, align gui.Align) {
//
//	f := new(field)
//	f.id = id
//	f.label = gui.NewLabel(label)
//	f.value = gui.NewLabel("")
//	f.align = align
//}
//
//func (s *Stats) recalc() {
//
//	maxLabelWidth := 0
//	for i := 0; i < len(s.fields); i ++ {
//
//
//
//
//
//
//
//
//
//
//	}
//}
