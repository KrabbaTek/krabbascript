package parser

import (
	"testing"

	"github.com/KrabbaTek/krabbascript/internal/lexer"
)

func TestListInitializer(t *testing.T) {
	source := `var messages: Arr<Arr<Str>> = Arr<Arr<Str>>{ { "Krabbas", "are", "awesome!" }, { "I", "love", "KrabbaScript!" } };`
	l := lexer.NewLexerFromStr(source, "test")
	p := NewParser(l.Scan(), "test")

	ast, err := p.Parse()
	if len(err) != 0 {
		t.Fatalf("parser reported %d errors", len(err))
	}

	block, ok := ast.ExtraInfo.(*NodeInfoBlock)
	if !ok || len(block.Elements) != 1 {
		t.Fatalf("expected one top-level declaration, got %#v", ast.ExtraInfo)
	}

	declaration := block.Elements[0]
	if declaration.Type != NodeVariableDef {
		t.Fatalf("expected variable definition, got %s", declaration.Type)
	}

	if declaration.Left.Type != NodeArray ||
		declaration.Left.Left.Type != NodeArray ||
		declaration.Left.Left.Left.Type != NodeStrType {
		t.Fatal("declared type is not Arr<Arr<Str>>")
	}

	initializer := declaration.Right
	if initializer.Type != NodeListInit ||
		initializer.Left.Type != NodeArray ||
		initializer.Left.Left.Type != NodeArray ||
		initializer.Left.Left.Left.Type != NodeStrType {
		t.Fatal("expected an Arr<Arr<Str>> list initializer")
	}

	rows, ok := initializer.ExtraInfo.(*NodeInfoBlock)
	if !ok || len(rows.Elements) != 2 {
		t.Fatalf("expected two nested array initializers, got %#v", initializer.ExtraInfo)
	}

	want := [][]string{
		{"Krabbas", "are", "awesome!"},
		{"I", "love", "KrabbaScript!"},
	}
	for rowIndex, row := range rows.Elements {
		elements, ok := row.ExtraInfo.(*NodeInfoBlock)
		if !ok || len(elements.Elements) != len(want[rowIndex]) {
			t.Fatalf("row %d has unexpected elements: %#v", rowIndex, row.ExtraInfo)
		}

		for elementIndex, element := range elements.Elements {
			if element.Type != NodeStrLit || element.Lexeme != want[rowIndex][elementIndex] {
				t.Errorf("row %d element %d = (%s, %q), want string %q",
					rowIndex, elementIndex, element.Type, element.Lexeme, want[rowIndex][elementIndex])
			}
		}
	}
}
