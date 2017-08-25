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

type Value struct {
	Name			string
	Connection		*Block
}

type Input struct {
	Type			int
	Name			string
	Text			string
	Value			*Value
	Statement		*Statement
}

type Block struct {
	Type			    string
	Root                bool

	PreviousConnection	*Block
	NextConnection		*Block
	OutputConnection	*Block

	Inputs			    []*Input
}

type Statement struct {
	Name			string
	Connection		*Block
}

func NewWorkspace(xml io.Reader) Workspace {
	if xml == nil {
		panic("Xml cannot be nil")
	}

	return Workspace{
		ExportedXml:	xml,
		Root:		[]*Block{},
	}
}

func (w *Workspace) Parse() error {
	dec := xml.NewDecoder(w.ExportedXml)

	for {
		t, err := dec.Token()
		if t == nil {
			break
		}else if err != nil {
			return err
		}

		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local ==  "block" {
				block := NewBlock()
				block.Root = true
				block.Type = getAttr(t, "type")
				w.Root = append(w.Root, &block)

				if err := w.parseStatement(dec, &block); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (w *Workspace) parseStatement(dec *xml.Decoder, ptr *Block) error{
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
					ptr.Type = getAttr(t, "type")
				}
				case "next": {
					b := NewBlock()

					b.PreviousConnection = ptr
					ptr.NextConnection = &b
					ptr = ptr.NextConnection
				}
				case "field": {
					var field struct {
						Name	string	`xml:"name,attr"`
						Value	string	`xml:",chardata"`
					}

					if err := dec.DecodeElement(&field, &t); err != nil {
						return err
					}

					ptr.Inputs = append(ptr.Inputs, &Input{
						Type: InputField,
						Name: field.Name,
						Text: field.Value,
					})
				}
				case "statement": {
					block := NewBlock()

					name := getAttr(t, "name")
					input := &Input {
						Type: InputStatement,
						Name: name,
						Statement: &Statement {
							Name: name,
							Connection: &block,
						},
					}

					if err := w.parseStatement(dec, &block); err != nil {
						return err
					}

					ptr.Inputs = append(ptr.Inputs, input)
				}
				case "value": {
					block := NewBlock()

					name := getAttr(t, "name")
					input := &Input {
						Type: InputValue,
						Name: name,
						Value: &Value {
							Name: name,
							Connection: &block,
						},
					}

					if err := w.parseStatement(dec, &block); err != nil {
						return err
					}
					ptr.Inputs = append(ptr.Inputs, input)
				}
			}
		case xml.EndElement:
			switch t.Name.Local {
				case "next": {
					ptr = ptr.PreviousConnection

					if ptr.Root {
						return nil
					}
				}
				case "statement", "value": {
					return nil
				}
			}
		}
	}
	return nil
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
		ptr := root

		printStatement(ptr, 0)
	}
}

func printStatement(block *Block, depth int){
	ptr := block

	for ptr != nil {
		fmt.Println(strings.Repeat("\t", depth) + "<block type=\"" + ptr.Type + "\">")

		for _, input := range ptr.Inputs {
			if input.Type == InputField {
				fmt.Println(strings.Repeat("\t", depth + 1) + "<field name=\"" + input.Name + "\">" + input.Text + "</field>")
			}else if input.Type == InputStatement {
				fmt.Println(strings.Repeat("\t", depth + 1) + "<statement name=\"" + input.Statement.Name + "\">")
				printStatement(input.Statement.Connection, depth + 2)
			}else if input.Type == InputValue {
				fmt.Println(strings.Repeat("\t", depth + 1) + "<value name=\"" + input.Value.Name + "\">")
				printStatement(input.Value.Connection, depth + 2)
			}
		}

		ptr = ptr.NextConnection
		depth++
	}
}
