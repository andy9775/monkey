package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	/*
		In the case of a function, e.store is the function's environment and
		e.outer is the scoped environment. This allows for:
		`
		let i = 10;
		let add = fn(i) { i + 10; };
		add(i);
		add(20);
		`
		Where the i in the outer scope is preserved. Therefore, in a function call,
		the current scope is this, and the outer scope (enclosing scope) is outer
		therefore we check this before outer.
	*/
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
