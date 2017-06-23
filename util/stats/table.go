package stats

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
)

// StatsTable is a gui.Table panel with statistics
type StatsTable struct {
	*gui.Table          // embedded table panel
	fields     []*field // array of fields to show
}

type field struct {
	id  string
	row int
}

// NewStatsTable creates and returns a pointer to a new statistics table panel
func NewStatsTable(width, height float32, gs *gls.GLS) *StatsTable {

	// Creates table panel
	st := new(StatsTable)
	t, err := gui.NewTable(width, height, []gui.TableColumn{
		{Id: "f", Header: "Stat", Width: 50, Minwidth: 32, Align: gui.AlignRight, Format: "%s", Resize: true, Expand: 2},
		{Id: "v", Header: "Value", Width: 50, Minwidth: 32, Align: gui.AlignRight, Format: "%d", Resize: false, Expand: 1},
	})
	if err != nil {
		panic(err)
	}
	st.Table = t
	st.ShowHeader(false)

	// Add rows
	st.addRow("shaders", "Shaders:")
	st.addRow("vaos", "Vaos:")
	st.addRow("buffers", "Buffers:")
	st.addRow("textures", "Textures:")
	st.addRow("unisets", "Uniforms/frame:")
	st.addRow("drawcalls", "Draw calls/frame:")
	st.addRow("cgocalls", "CGO calls/frame:")
	return st
}

// Update updates the table values from the specified stats table
func (st *StatsTable) Update(s *Stats) {

	for i := 0; i < len(st.fields); i++ {
		f := st.fields[i]
		switch f.id {
		case "shaders":
			st.Table.SetCell(f.row, "v", s.Glstats.Shaders)
		case "vaos":
			st.Table.SetCell(f.row, "v", s.Glstats.Vaos)
		case "buffers":
			st.Table.SetCell(f.row, "v", s.Glstats.Buffers)
		case "textures":
			st.Table.SetCell(f.row, "v", s.Glstats.Textures)
		case "unisets":
			st.Table.SetCell(f.row, "v", s.Unisets)
		case "drawcalls":
			st.Table.SetCell(f.row, "v", s.Drawcalls)
		case "cgocalls":
			st.Table.SetCell(f.row, "v", s.Cgocalls)
		}
	}
}

func (st *StatsTable) addRow(id, label string) {

	f := new(field)
	f.id = id
	f.row = st.Table.RowCount()
	st.Table.AddRow(map[string]interface{}{"f": label, "v": 0})
	st.fields = append(st.fields, f)
}
