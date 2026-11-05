package resolver

// this is a test file.
// it will be deleted later.
import (
	"github.com/KrabbaTek/krabbascript/internal/lexer"
	"github.com/KrabbaTek/krabbascript/internal/parser"
	"testing"
)

// resolveSource lexes, parses and resolves src, and returns every symbol key
// found in the returned scopes along with the resolver errors.
func resolveSource(t *testing.T, src string) (map[string]bool, []*ResolverError) {
	t.Helper()

	l := lexer.NewLexerFromStr(src, "test")
	toks, lerrs := l.Scan()
	if len(lerrs) > 0 {
		t.Fatalf("lexer reported %d errors", len(lerrs))
	}

	p := parser.NewParser(toks, "test")
	ast, perrs := p.Parse()
	if len(perrs) > 0 {
		t.Fatalf("parser reported %d errors", len(perrs))
	}

	r := NewResolver("test")
	table, errs := r.Walk(ast)

	keys := make(map[string]bool)
	for _, scope := range table {
		for k := range scope {
			keys[k] = true
		}
	}

	return keys, errs
}

func TestMethodKeys(t *testing.T) {
	const body = " {\n    var a = 1;\n}\n"

	tests := []struct {
		name     string
		src      string
		wantKeys []string
		noKeys   []string
		wantErrs []string
	}{
		{
			name:     "same name on different types",
			src:      "method foo() of I32" + body + "method foo() of Str" + body,
			wantKeys: []string{"I32.foo", "Str.foo"},
			noKeys:   []string{"foo"},
		},
		{
			name:     "same name on different array types",
			src:      "method f() of Arr<I32>" + body + "method f() of Arr<Str>" + body,
			wantKeys: []string{"Arr<I32>.f", "Arr<Str>.f"},
		},
		{
			name:     "same name on the same type",
			src:      "method foo() of I32" + body + "method foo() of I32" + body,
			wantKeys: []string{"I32.foo"},
			wantErrs: []string{"test:4:1: redeclaration of method foo on I32 at test:1:1"},
		},
		{
			name:     "variable and method share a name",
			src:      "var foo = 1;\nmethod foo() of I32" + body,
			wantKeys: []string{"foo", "I32.foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, errs := resolveSource(t, tt.src)

			for _, k := range tt.wantKeys {
				if !keys[k] {
					t.Errorf("expected symbol %q, got %v", k, keys)
				}
			}
			for _, k := range tt.noKeys {
				if keys[k] {
					t.Errorf("did not expect symbol %q, got %v", k, keys)
				}
			}

			if len(errs) != len(tt.wantErrs) {
				t.Fatalf("got %d errors %v, want %d", len(errs), errs, len(tt.wantErrs))
			}
			for i, w := range tt.wantErrs {
				if got := errs[i].What.Error(); got != w {
					t.Errorf("[%d] got %q, want %q", i, got, w)
				}
			}
		})
	}
}
