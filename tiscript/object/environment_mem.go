package object

func NewMemoryEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s,
		IsMem:     true,
		Namespace: "mem",
		outer:     nil}
}

func (e *Environment) MemGet(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}
func (e *Environment) MemSet(name string, val Object) Object {
	e.store[name] = val
	return val
}

func NewEnclosedMemoryEnvironment(outer *Environment) *Environment {
	env := NewMemoryEnvironment()
	env.outer = outer
	return env
}
