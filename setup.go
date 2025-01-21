package duckdb

import (
	"context"
	"database/sql"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin"
	_ "github.com/marcboeker/go-duckdb"
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
	db, err := sql.Open("duckdb", dbPath)
	if err != nil {
		return plugin.Error("duckdblog", err)
	}
	defer db.Close()

	conn, err := db.Conn(context.Background())
	defer conn.Close()

	// Cria a tabela de logs, se não existir
	createTable := `CREATE TABLE IF NOT EXISTS dns_logs (
		timestamp TIMESTAMP,
		client_ip TEXT,
		query_name TEXT,
		query_type INTEGER
	)`
	_, err = conn.QueryContext(context.Background(), createTable)
	if err != nil {
		return plugin.Error("duckdblog", err)
	}

	// Registra o plugin
	plugin.Register("duckdblog", func(next plugin.Handler) plugin.Handler {
		return DuckDBLogger{Next: next, DBPath: dbPath, conn: conn}
	})

	return nil
}

func init() { plugin.Register("duckdblog", setup) }
