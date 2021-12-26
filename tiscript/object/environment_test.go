package object

import (
	"context"
	"fmt"
	"testing"

	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
)

func TestGoTikvClient(t *testing.T) {
	cli, err := rawkv.NewClient(context.TODO(), []string{"127.0.0.1:2379"}, config.DefaultConfig().Security)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	t.Logf("cluster ID: %d\n", cli.ClusterID())

	key := []byte("Company")
	val := []byte("PingCAP")

	// put key into tikv
	err = cli.Put(context.TODO(), key, val)
	if err != nil {
		panic(err)
	}
	t.Logf("Successfully put %s:%s to tikv\n", key, val)

	// get key from tikv
	val, err = cli.Get(context.TODO(), key)
	if err != nil {
		panic(err)
	}
	t.Logf("found val: %s for key: %s\n", val, key)

	// delete key from tikv
	err = cli.Delete(context.TODO(), key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("key: %s deleted\n", key)

	// get key again from tikv
	val, err = cli.Get(context.TODO(), key)
	if err != nil {
		panic(err)
	}
	t.Logf("found val: %s for key: %s\n", val, key)
}
