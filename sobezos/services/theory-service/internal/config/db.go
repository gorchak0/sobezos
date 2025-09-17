package config

import (
	"database/sql"
	"sobezos/services/theory-service/config"

	_ "github.com/lib/pq"
)

func OpenDB() (*sql.DB, error) {
	cfg := config.LoadEnvConfig()
	connStr := "host=" + cfg.DBHost +
		" port=" + cfg.DBPort +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" sslmode=" + cfg.DBSSLMode
	return sql.Open("postgres", connStr)
}
