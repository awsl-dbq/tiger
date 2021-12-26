package object

import (
	"context"

	"github.com/go-basic/uuid"
	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
)

func NewEnvironment() *Environment {
	cli, err := rawkv.NewClient(context.TODO(), []string{"127.0.0.1:2379"}, config.DefaultConfig().Security)
	if err != nil {
		panic(err)
	}
	clean := func() {
		cli.Close()
	}
	ns := uuid.New()
	return &Environment{
		namespace: ns,
		cli:       cli, outer: nil, clean: clean}
}

type Environment struct {
	namespace string
	cli       *rawkv.Client
	outer     *Environment
	clean     func()
}

func (e *Environment) Get(name string) (Object, bool) {
	k := []byte(e.namespace + "_" + name)
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
	k := []byte(e.namespace + "_" + name)
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
