package object

import (
	"context"
	"os"
	"strings"

	"github.com/go-basic/uuid"
	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
)

func useMem() bool {
	env := os.Environ()
	for _, v := range env {
		if strings.Contains(strings.ToLower(v), "usemem") {
			return true
		}
	}
	return false
}
func NewEnvironment() *Environment {
	if useMem() {
		return NewMemoryEnvironment()
	}
	cli, err := rawkv.NewClient(context.TODO(), []string{"127.0.0.1:2379"}, config.DefaultConfig().Security)
	if err != nil {
		panic(err)
	}
	clean := func() {
		cli.Close()
	}
	ns := uuid.New()
	return &Environment{
		Namespace: ns,
		IsMem:     false,
		cli:       cli, outer: nil, clean: clean}
}

type Environment struct {
	Namespace string
	IsMem     bool
	store     map[string]Object
	cli       *rawkv.Client
	outer     *Environment
	clean     func()
}

func (e *Environment) Get(name string) (Object, bool) {
	if e.IsMem {
		return e.MemGet(name)
	}
	k := []byte(e.Namespace + "_" + name)
	val, err := e.cli.Get(context.TODO(), k)
	ok := true
	if err != nil {
		ok = false
	}
	obj := FromBytesArray(val)
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}
func (e *Environment) Set(name string, val Object) Object {
	if e.IsMem {
		return e.MemSet(name, val)
	}
	k := []byte(e.Namespace + "_" + name)
	v := val.ToBytesArray()
	err := e.cli.Put(context.TODO(), k, v)
	if err != nil {
		panic(err)
	}
	return val
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
