package resolver

// this is a test file. should be removed later.
import (
	"github.com/KrabbaTek/krabbascript/internal/parser"
	"testing"
)

func TestTypeString(t *testing.T) {
	i32 := &parser.Node{Type: parser.NodeI32Type}
	str := &parser.Node{Type: parser.NodeStrType}
	tests := []struct {
		name string
		node *parser.Node
		want string
	}{
		{"builtin", i32, "I32"},
		{"struct name", &parser.Node{Type: parser.NodeLit, Lexeme: "Point"}, "Point"},
		{"array of I32", &parser.Node{Type: parser.NodeArray, Left: i32}, "Arr<I32>"},
		{"array of Str", &parser.Node{Type: parser.NodeArray, Left: str}, "Arr<Str>"},
		{"pointer", &parser.Node{Type: parser.NodePointer, Left: i32}, "*I32"},
		{"nested", &parser.Node{Type: parser.NodeArray,
			Left: &parser.Node{Type: parser.NodePointer, Left: i32}}, "Arr<*I32>"},
		{"nil", nil, "?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := typeString(tt.node); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
