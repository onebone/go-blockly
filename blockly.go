package blockly

import (
	"io"
	"encoding/xml"
	"fmt"
	"strings"
)

const (
	InputDummy = iota
	InputField
	InputValue
	InputStatement
)

type Workspace struct {
	ExportedXml		io.Reader

	// TODO: Variables
	Root			[]*Block
}

type Input struct {
	Type			int
	Name			string
	Value			string
	Statement		*Statement
}

type Block struct {
	Type			string

	PreviousConnection	*Block
	NextConnection		*Block
	OutputConnection	*Block

	Inputs			[]*Input
}

type Statement struct {
	Name			string
	Connection		*Block
}

func NewWorkspace(xml io.Reader) Workspace {
	if xml == nil {
		panic("Xml cannot nil")
	}

	return Workspace{
		ExportedXml:	xml,
		Root:		[]*Block{},
	}
}

func (w *Workspace) Parse() error {
	dec := xml.NewDecoder(w.ExportedXml)

	block := NewBlock()
	var ptr *Block = &block
	w.Root = append(w.Root, ptr) // Defining initial root block

	for {
		t, err := dec.Token()
		if t == nil {
			break
		}else if err != nil {
			return err
		}

		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "block": {
				blType := getAttr(t, "type")

				ptr.Type = blType
			}
			case "next": {
				b := NewBlock()

				b.PreviousConnection = ptr
				ptr.NextConnection = &b
				ptr = ptr.NextConnection
			}
			case "field": {
				var field struct {
					Name	string `xml:"name,attr"`
					Value	string	`xml:",chardata"`
				}

				if err := dec.DecodeElement(&field, &t); err != nil {
					return err
				}

				ptr.Inputs = append(ptr.Inputs, &Input{
					Type: InputField,
					Name: field.Name,
					Value: field.Value,
				})
			}
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "next": {
				ptr = ptr.PreviousConnection
			}
			}
		}
	}

	return nil
}

func (w *Workspace) parseStatement(dec *xml.Decoder){
	// TODO
}

func NewBlock() Block {
	return Block {
		Inputs: []*Input{},
	}
}

func getAttr(t xml.StartElement, key string) string{
	for _, a := range t.Attr {
		if a.Name.Local == key {
			return a.Value
		}
	}

	return ""
}

func printWorkspace(w Workspace){
	for _, root := range w.Root {
		depth := 0

		ptr := root

		for ptr != nil {
			fmt.Println(strings.Repeat("\t", depth) + "<block type=\"" + ptr.Type + "\">")

			for _, input := range ptr.Inputs {
				fmt.Println(strings.Repeat("\t", depth + 1) + "<field name=\"" + input.Name + "\">" + input.Value + "</field>")
			}

			ptr = ptr.NextConnection
			depth++
		}
	}
}
