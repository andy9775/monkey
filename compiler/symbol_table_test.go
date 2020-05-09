package compiler_test

import (
	"testing"

	"github.com/andy9775/monkey/compiler"
)

func TestResolveLocal(t *testing.T) {
	global := compiler.NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := compiler.NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := []compiler.Symbol{
		compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0},
		compiler.Symbol{Name: "b", Scope: compiler.GlobalScope, Index: 1},
		compiler.Symbol{Name: "c", Scope: compiler.LocalScope, Index: 0},
		compiler.Symbol{Name: "d", Scope: compiler.LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := local.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}
		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v",
				sym.Name, sym, result)
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	global := compiler.NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := compiler.NewEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := compiler.NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	tests := []struct {
		table           *compiler.SymbolTable
		expectedSymbols []compiler.Symbol
	}{
		{
			firstLocal,
			[]compiler.Symbol{
				compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0},
				compiler.Symbol{Name: "b", Scope: compiler.GlobalScope, Index: 1},
				compiler.Symbol{Name: "c", Scope: compiler.LocalScope, Index: 0},
				compiler.Symbol{Name: "d", Scope: compiler.LocalScope, Index: 1},
			},
		},
		{
			secondLocal,
			[]compiler.Symbol{
				compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0},
				compiler.Symbol{Name: "b", Scope: compiler.GlobalScope, Index: 1},
				compiler.Symbol{Name: "e", Scope: compiler.LocalScope, Index: 0},
				compiler.Symbol{Name: "f", Scope: compiler.LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v",
					sym.Name, sym, result)
			}
		}

	}
}

func TestDefine(t *testing.T) {
	expected := map[string]compiler.Symbol{
		"a": compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0},
		"b": compiler.Symbol{Name: "b", Scope: compiler.GlobalScope, Index: 1},
		"c": compiler.Symbol{Name: "c", Scope: compiler.LocalScope, Index: 0},
		"d": compiler.Symbol{Name: "d", Scope: compiler.LocalScope, Index: 1},
		"e": compiler.Symbol{Name: "e", Scope: compiler.LocalScope, Index: 0},
		"f": compiler.Symbol{Name: "f", Scope: compiler.LocalScope, Index: 1},
	}
	global := compiler.NewSymbolTable()
	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
	}

	firstLocal := compiler.NewEnclosedSymbolTable(global)
	c := firstLocal.Define("c")
	if c != expected["c"] {
		t.Errorf("expected c=%+v, got=%+v", expected["c"], c)
	}

	d := firstLocal.Define("d")
	if d != expected["d"] {
		t.Errorf("expected d=%+v, got=%+v", expected["d"], d)
	}

	secondLocal := compiler.NewEnclosedSymbolTable(firstLocal)

	e := secondLocal.Define("e")
	if e != expected["e"] {
		t.Errorf("expected e=%+v, got=%+v", expected["e"], e)
	}

	f := secondLocal.Define("f")
	if f != expected["f"] {
		t.Errorf("expected f=%+v, got=%+v", expected["f"], f)
	}
}

func TestResolveGlobal(t *testing.T) {
	global := compiler.NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []compiler.Symbol{
		compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0},
		compiler.Symbol{Name: "b", Scope: compiler.GlobalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}
}
