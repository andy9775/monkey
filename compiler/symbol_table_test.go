package compiler_test

import (
	"testing"

	"github.com/andy9775/monkey/compiler"
)

func TestShadowingFunctionName(t *testing.T) {
	/*
		ensure the following works:

		let fb = fn() {
			let fb = 1;
			fb;
		};
	*/
	global := compiler.NewSymbolTable()
	global.DefineFunctionName("a")
	global.Define("a")
	expected := compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0}

	result, ok := global.Resolve(expected.Name)
	if !ok {
		t.Fatalf("function name %s not resolvable", expected.Name)
	}
	if result != expected {
		t.Errorf("expected %s to resolve to %+v, got=%+v",
			expected.Name, expected, result)
	}
}

func TestDefineAndResolveFunctionName(t *testing.T) {
	global := compiler.NewSymbolTable()
	global.DefineFunctionName("a")
	expected := compiler.Symbol{Name: "a", Scope: compiler.FunctionScope, Index: 0}

	result, ok := global.Resolve(expected.Name)
	if !ok {
		t.Fatalf("function name %s not resolvable", expected.Name)
	}

	if result != expected {
		t.Errorf("expected %s to resolve to %+v, got=%+v",
			expected.Name, expected, result)
	}
}

func TestResolveFree(t *testing.T) {
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
		table               *compiler.SymbolTable
		expectedSymbols     []compiler.Symbol
		expectedFreeSymbols []compiler.Symbol
	}{
		{
			firstLocal,
			[]compiler.Symbol{
				compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0},
				compiler.Symbol{Name: "b", Scope: compiler.GlobalScope, Index: 1},
				compiler.Symbol{Name: "c", Scope: compiler.LocalScope, Index: 0},
				compiler.Symbol{Name: "d", Scope: compiler.LocalScope, Index: 1},
			},
			[]compiler.Symbol{},
		},
		{
			secondLocal,
			[]compiler.Symbol{
				compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0},
				compiler.Symbol{Name: "b", Scope: compiler.GlobalScope, Index: 1},
				compiler.Symbol{Name: "c", Scope: compiler.FreeScope, Index: 0},
				compiler.Symbol{Name: "d", Scope: compiler.FreeScope, Index: 1},
				compiler.Symbol{Name: "e", Scope: compiler.LocalScope, Index: 0},
				compiler.Symbol{Name: "f", Scope: compiler.LocalScope, Index: 1},
			},
			[]compiler.Symbol{
				compiler.Symbol{Name: "c", Scope: compiler.LocalScope, Index: 0},
				compiler.Symbol{Name: "d", Scope: compiler.LocalScope, Index: 1},
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

		if len(tt.table.FreeSymbols) != len(tt.expectedFreeSymbols) {
			t.Errorf("wrong number of free symbols. got=%d, want=%d",
				len(tt.table.FreeSymbols), len(tt.expectedFreeSymbols))
			continue
		}

		for i, sym := range tt.expectedFreeSymbols {
			result := tt.table.FreeSymbols[i]
			if result != sym {
				t.Errorf("wrong free symbol. got=%+v, want=%+v",
					result, sym)
			}
		}
	}
}

func TestResolveUnresolvableFree(t *testing.T) {
	global := compiler.NewSymbolTable()
	global.Define("a")

	firstLocal := compiler.NewEnclosedSymbolTable(global)
	firstLocal.Define("c")

	secondLocal := compiler.NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	expected := []compiler.Symbol{
		compiler.Symbol{Name: "a", Scope: compiler.GlobalScope, Index: 0},
		compiler.Symbol{Name: "c", Scope: compiler.FreeScope, Index: 0},
		compiler.Symbol{Name: "e", Scope: compiler.LocalScope, Index: 0},
		compiler.Symbol{Name: "f", Scope: compiler.LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := secondLocal.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}
		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v",
				sym.Name, sym, result)
		}
	}

	expectedUnresolvable := []string{
		"b",
		"d",
	}

	for _, name := range expectedUnresolvable {
		_, ok := secondLocal.Resolve(name)
		if ok {
			t.Errorf("name %s resolved, but was expected not to", name)
		}
	}
}

func TestDefineResolveBuiltins(t *testing.T) {
	global := compiler.NewSymbolTable()
	firstLocal := compiler.NewEnclosedSymbolTable(global)
	secondLocal := compiler.NewEnclosedSymbolTable(firstLocal)

	expected := []compiler.Symbol{
		compiler.Symbol{Name: "a", Scope: compiler.BuiltinScope, Index: 0},
		compiler.Symbol{Name: "c", Scope: compiler.BuiltinScope, Index: 1},
		compiler.Symbol{Name: "e", Scope: compiler.BuiltinScope, Index: 2},
		compiler.Symbol{Name: "f", Scope: compiler.BuiltinScope, Index: 3},
	}
	for i, v := range expected {
		global.DefineBuiltin(i, v.Name)

	}

	for _, table := range []*compiler.SymbolTable{global, firstLocal, secondLocal} {
		for _, sym := range expected {
			result, ok := table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
			}
		}
	}
}

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
