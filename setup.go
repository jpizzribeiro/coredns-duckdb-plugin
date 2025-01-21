package duckdb

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin"
	"github.com/marcboeker/go-duckdb"
)

// Setup configura o plugin DuckDBLogger
func setup(c *caddy.Controller) error {
	var dbPath string

	// Lê as configurações do Corefile
	for c.Next() {
		if !c.Args(&dbPath) {
			return plugin.Error("duckdblog", c.ArgErr())
		}
	}

	// Conecta ao DuckDB
	conn, err := duckdb.Open(dbPath)
	if err != nil {
		return plugin.Error("duckdblog", err)
	}

	// Cria a tabela de logs, se não existir
	createTable := `CREATE TABLE IF NOT EXISTS dns_logs (
		timestamp TIMESTAMP,
		client_ip TEXT,
		query_name TEXT,
		query_type INTEGER
	)`
	_, err = conn.Exec(createTable)
	if err != nil {
		return plugin.Error("duckdblog", err)
	}

	// Registra o plugin
	plugin.Register("duckdblog", func(next plugin.Handler) plugin.Handler {
		return DuckDBLogger{Next: next, DBPath: dbPath, conn: conn}
	})

	return nil
}
