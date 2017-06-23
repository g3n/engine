package stats

import (
	"fmt"
	"runtime"
	"time"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
)

type StatsTable struct {
	*gui.Table          // embedded table
	fields     []*field // array of fields
	prev       gls.Stats
	frames     int
	cgocalls   int64
	last       time.Time
}

type field struct {
	id  string
	row int
}

func NewStatsTable(width, height float32) *StatsTable {

	s := new(StatsTable)
	t, err := gui.NewTable(width, height, []gui.TableColumn{
		{Id: "f", Header: "Stat", Width: 50, Minwidth: 32, Align: gui.AlignRight, Format: "%s", Resize: true, Expand: 1.5},
		{Id: "v", Header: "Value", Width: 50, Minwidth: 32, Align: gui.AlignRight, Format: "%d", Resize: false, Expand: 1},
	})
	if err != nil {
		panic(err)
	}
	s.Table = t
	s.ShowHeader(false)
	s.addRow("shaders", "Shaders:")
	s.addRow("vaos", "Vaos:")
	s.addRow("vbos", "Vbos:")
	s.addRow("textures", "Textures:")
	s.addRow("unisets", "Uniforms/frame:")
	s.addRow("cgocalls", "CGO calls/frame:")
	s.last = time.Now()
	s.cgocalls = runtime.NumCgoCall()
	return s
}

func (s *StatsTable) Update(gs *gls.GLS, d time.Duration) {

	now := time.Now()
	s.frames++
	if s.last.Add(d).After(now) {
		return
	}
	s.update(gs)
	s.last = now
	s.frames = 0
}

func (s *StatsTable) update(gs *gls.GLS) {

	var stats gls.Stats
	gs.Stats(&stats)
	for i := 0; i < len(s.fields); i++ {
		f := s.fields[i]
		switch f.id {
		case "shaders":
			if stats.Shaders != s.prev.Shaders {
				s.Table.SetCell(f.row, "v", stats.Shaders)
				fmt.Println("update")
			}
		case "vaos":
			if stats.Vaos != s.prev.Vaos {
				s.Table.SetCell(f.row, "v", stats.Vaos)
				fmt.Println("update")
			}
		case "vbos":
			if stats.Vbos != s.prev.Vbos {
				s.Table.SetCell(f.row, "v", stats.Vbos)
				fmt.Println("update")
			}
		case "textures":
			if stats.Textures != s.prev.Textures {
				s.Table.SetCell(f.row, "v", stats.Textures)
				fmt.Println("update")
			}
		case "cgocalls":
			current := runtime.NumCgoCall()
			calls := current - s.cgocalls
			s.Table.SetCell(f.row, "v", int(float64(calls)/float64(s.frames)))
			s.cgocalls = current
		}
	}
	s.prev = stats
}

func (s *StatsTable) addRow(id, label string) {

	f := new(field)
	f.id = id
	f.row = s.Table.RowCount()
	s.Table.AddRow(map[string]interface{}{"f": label, "v": 0})
	s.fields = append(s.fields, f)
}
