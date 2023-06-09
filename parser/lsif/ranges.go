package lsif

import (
	"errors"
	"io"
	"strconv"

	"github.com/goccy/go-json"
)

const (
	definitions = "definitions"
	references  = "references"
)

type Ranges struct {
	DefRefs    map[Id]Item
	References *References
	Hovers     *Hovers
	Cache      Cache

	NextMap           map[Id]Id
	TextReferenceMap  map[Id]Id
	TextDefinitionMap map[Id]Id
	RawEdgeMap        map[Id][]RawItem
}

type Next struct {
	Id    Id     `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
	OutV  Id     `json:"outV"`
	InV   Id     `json:"inV"`
}

type TextReference struct {
	Id    Id     `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
	OutV  Id     `json:"outV"`
	InV   Id     `json:"inV"`
}

type RawRange struct {
	Id   Id    `json:"id"`
	Data Range `json:"start"`
	End  Range `json:"end"`
}

type Range struct {
	Line      int32 `json:"line"`
	Character int32 `json:"character"`
	Length    int32 `json:"length"`
	RefId     Id
}

type RawItem struct {
	Property string `json:"property"`
	RefId    Id     `json:"outV"`
	RangeIds []Id   `json:"inVs"`
	DocId    Id     `json:"document"`
}

type Item struct {
	Line    int32
	DocId   Id
	RangeId Id
}

type SerializedRange struct {
	StartLine      int32                 `json:"start_line"`
	StartChar      int32                 `json:"start_char"`
	DefinitionPath string                `json:"definition_path,omitempty"`
	Hover          json.RawMessage       `json:"hover"`
	References     []SerializedReference `json:"references,omitempty"`
}

func NewRanges() (*Ranges, error) {
	hovers, err := NewHovers()
	if err != nil {
		return nil, err
	}

	references, err := NewReferences()
	if err != nil {
		return nil, err
	}

	cache, err := newCache("ranges", Range{})
	if err != nil {
		return nil, err
	}

	return &Ranges{
		DefRefs:           make(map[Id]Item),
		References:        references,
		Hovers:            hovers,
		Cache:             cache,
		NextMap:           make(map[Id]Id),
		TextReferenceMap:  make(map[Id]Id),
		TextDefinitionMap: make(map[Id]Id),
		RawEdgeMap:        map[Id][]RawItem{},
	}, nil
}

func (r *Ranges) Read(label string, line []byte) error {
	switch label {
	case "range":
		if err := r.addRange(line); err != nil {
			return err
		}
	case "item":
		if err := r.addItem(line); err != nil {
			return err
		}
	default:
		switch label {
		case "next":
			var rawNext Next
			if err := json.Unmarshal(line, &rawNext); err != nil {
				return err
			}
			r.NextMap[rawNext.OutV] = rawNext.InV
		case "textDocument/references":
			var textReference TextReference
			if err := json.Unmarshal(line, &textReference); err != nil {
				return err
			}
			r.TextReferenceMap[textReference.OutV] = textReference.InV
		case "textDocument/definition":
			var textReference TextReference
			if err := json.Unmarshal(line, &textReference); err != nil {
				return err
			}
			r.TextDefinitionMap[textReference.OutV] = textReference.InV
		}

		// currently we do not need hover
		// return r.Hovers.Read(label, line)
	}

	return nil
}

func (r *Ranges) Serialize(f io.Writer, rangeIds []Id, docs map[Id]string) error {
	encoder := json.NewEncoder(f)
	n := len(rangeIds)

	if _, err := f.Write([]byte("[")); err != nil {
		return err
	}

	for i, rangeId := range rangeIds {
		entry, err := r.getRange(rangeId)
		if err != nil {
			continue
		}

		serializedRange := SerializedRange{
			StartLine:      entry.Line,
			StartChar:      entry.Character,
			DefinitionPath: r.definitionPathFor(docs, entry.RefId),
			Hover:          r.Hovers.For(entry.RefId),
			References:     r.References.For(docs, entry.RefId),
		}
		if err := encoder.Encode(serializedRange); err != nil {
			return err
		}
		if i+1 < n {
			if _, err := f.Write([]byte(",")); err != nil {
				return err
			}
		}
	}

	if _, err := f.Write([]byte("]")); err != nil {
		return err
	}

	return nil
}

func (r *Ranges) Close() error {
	for _, err := range []error{
		r.Cache.Close(),
		r.References.Close(),
		r.Hovers.Close(),
	} {
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Ranges) definitionPathFor(docs map[Id]string, refId Id) string {
	defRef, ok := r.DefRefs[refId]
	if !ok {
		return ""
	}

	defPath := docs[defRef.DocId] + "#L" + strconv.Itoa(int(defRef.Line))

	return defPath
}

func (r *Ranges) addRange(line []byte) error {
	var rg RawRange
	if err := json.Unmarshal(line, &rg); err != nil {
		return err
	}

	length := rg.End.Character - rg.Data.Character
	if length > 0 {
		rg.Data.Length = length
	}

	return r.Cache.SetEntry(rg.Id, &rg.Data)
}

func (r *Ranges) addItem(line []byte) error {
	var rawItem RawItem
	if err := json.Unmarshal(line, &rawItem); err != nil {
		return err
	}

	// store these edges whatever
	if l, ok := r.RawEdgeMap[rawItem.RefId]; ok {
		l = append(l, rawItem)
	} else {
		r.RawEdgeMap[rawItem.RefId] = []RawItem{rawItem}
	}

	if rawItem.Property != definitions && rawItem.Property != references {
		return nil
	}

	if len(rawItem.RangeIds) == 0 {
		return errors.New("no range IDs")
	}

	var references []Item

	for _, rangeId := range rawItem.RangeIds {
		rg, err := r.getRange(rangeId)
		if err != nil {
			return err
		}

		rg.RefId = rawItem.RefId

		if err := r.Cache.SetEntry(rangeId, rg); err != nil {
			return err
		}

		item := Item{
			Line:    rg.Line + 1,
			DocId:   rawItem.DocId,
			RangeId: rangeId,
		}

		if rawItem.Property == definitions {
			r.DefRefs[rawItem.RefId] = item
		} else {
			references = append(references, item)
		}
	}

	if err := r.References.Store(rawItem.RefId, references); err != nil {
		return err
	}

	return nil
}

func (r *Ranges) getRange(rangeId Id) (*Range, error) {
	var rg Range
	if err := r.Cache.Entry(rangeId, &rg); err != nil {
		return nil, err
	}

	return &rg, nil
}
