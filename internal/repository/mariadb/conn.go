package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/I-Van-Radkov/messenger/internal/config"
)

func MustConnect(cfg *config.DBConfig) *sql.DB {
	connectString := getConnectString(cfg)

	db, err := sql.Open("mysql", connectString)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}

func getConnectString(cfg *config.DBConfig) string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}
