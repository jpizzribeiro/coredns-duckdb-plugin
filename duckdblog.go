package duckdb

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/marcboeker/go-duckdb"
	"github.com/miekg/dns"
	"net"
	"time"
)

type DuckDBLogger struct {
	Next   plugin.Handler
	DBPath string
	conn   *duckdb.Conn
}

// ServeDNS processa a consulta DNS e registra o log no DuckDB.
func (d DuckDBLogger) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	// Obter informações da consulta
	queryName := state.Name()
	queryType := state.QType()
	clientIP, _, _ := net.SplitHostPort(state.IP())
	timestamp := time.Now()

	// Inserir log no DuckDB
	query := `INSERT INTO dns_logs (timestamp, client_ip, query_name, query_type) VALUES (?, ?, ?, ?)`
	_, err := d.conn.Exec(query, timestamp, clientIP, queryName, queryType)
	if err != nil {
		fmt.Printf("Erro ao registrar log: %v\n", err)
	}

	// Passa para o próximo plugin na cadeia
	return plugin.NextOrFailure(d.Name(), d.Next, ctx, w, r)
}

// Name retorna o nome do plugin.
func (d DuckDBLogger) Name() string { return "duckdblog" }
