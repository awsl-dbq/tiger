package main

import (
	"context"
	"fmt"

	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
)

func main() {
	cli, err := rawkv.NewClient(context.TODO(), []string{"127.0.0.1:2379"}, config.DefaultConfig().Security)
	if err != nil {
		panic(err)
	}
	ks, vs, err := cli.Scan(context.TODO(), nil, nil, 8000)
	if err != nil {
		panic(err)
	}
	for idx, k := range ks {
		fmt.Println(string(k), " -> ", string(vs[idx]))
	}
}
