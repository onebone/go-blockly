package blockly

import (
	"testing"
	"os"
)

var simple = []*Block{
	&Block {
		Type:		"b1",
		Inputs:		[]*Input{},
		NextConnection: &Block {
			Type:		"b2",
			Inputs:		[]*Input{},
		},
	},
}
func TestSimple(t *testing.T){
	f, err := os.Open("test/simple.xml")
	if err != nil {
		t.Error("Error opening file")
	}

	defer f.Close()

	w := NewWorkspace(f)
	if err := w.Parse(); err != nil {
		t.Error("Error during parse", err)
	}

	/*if !reflect.DeepEqual(w.Root, simple) {
		t.Error("Unexpected output")
	}*/
	printWorkspace(w)
}

func TestStatement(t *testing.T){
	f, err := os.Open("test/statement.xml")
	if err != nil {
		t.Error("Error opening file")
	}

	defer f.Close()

	w := NewWorkspace(f)
	if err := w.Parse(); err != nil {
		t.Error("Error during parse", err)
	}

	/*if !reflect.DeepEqual(w.Root, simple) {
		t.Error("Unexpected output")
	}*/
	printWorkspace(w)
}

func TestValue(t *testing.T){
	f, err := os.Open("test/value.xml")
	if err != nil {
		t.Error("Error opening file")
	}

	defer f.Close()

	w := NewWorkspace(f)
	if err := w.Parse(); err != nil {
		t.Error("Error during parse", err)
	}

	/*if !reflect.DeepEqual(w.Root, simple) {
		t.Error("Unexpected output")
	}*/
	printWorkspace(w)
}