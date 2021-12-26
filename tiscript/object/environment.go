package object

import (
	"context"
	"os"
	"strings"

	"github.com/go-basic/uuid"
	"github.com/pingcap/log"
	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
	"go.uber.org/zap"
)

func init() {
	beQuiet()
}
func beQuiet() {
	logger := zap.NewNop()
	defer logger.Sync()

	log.ReplaceGlobals(logger, nil)
	log.Info("you should not see this line.")
}
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
	s := make(map[string]Object)
	ns := "tiscript-" + uuid.New()
	return &Environment{
		Namespace: ns,
		store:     s,
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
	ok := false
	if err != nil {
		ok = false
	}
	var obj Object
	if len(val) > 0 {
		if val[0] == byte(0x00) { // 对于不方便序列化的 放内存
			return e.MemGet(name)
		}
		obj = FromBytesArray(val)
		ok = true
	}
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
	if len(v) == 1 && v[0] == 0x00 {
		// 对于不方便序列化的 放内存
		e.MemSet(name, val)
	}
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
