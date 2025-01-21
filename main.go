package main

import (
	"github.com/coredns/coredns/coremain"
	_ "github.com/jpizzribeiro/coredns-duckdb-plugin/plugin/duckdblog"
)

func main() {
	coremain.Run()
}
