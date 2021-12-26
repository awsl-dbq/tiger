package mem

import (
	obj "github.com/awsl-dbq/tiger/tiscript/object"
)

func NewEnvironment() *Environment {
	s := make(map[string]obj.Object)
	return &Environment{store: s, outer: nil}
}

type Environment struct {
	store map[string]obj.Object
	outer *Environment
}

func (e *Environment) Get(name string) (obj.Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}
func (e *Environment) Set(name string, val obj.Object) obj.Object {
	e.store[name] = val
	return val
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
