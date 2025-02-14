package db

import (
	"log"
	"log/slog"

	"github.com/gocql/gocql"
	"github.com/tanjed/go-sso/internal/config"
)

const DB_DRIVER_POSTGRES = "postgres"

type DB struct {
	Conn *gocql.Session
}

func Init() *DB{
	var db DB

	cluster := gocql.NewCluster(config.AppConfig.DB_HOST)
	cluster.Port = config.AppConfig.DB_PORT
    cluster.Keyspace = config.AppConfig.DB_NAME     
	
    cluster.Consistency = gocql.Quorum
	// cluster.Timeout = 60 * time.Second  // Increase the timeout
	// cluster.ConnectTimeout = 60 * time.Second


	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.AppConfig.DB_USER,
		Password: config.AppConfig.DB_PASSWORD,
	}

	session, err := cluster.CreateSession()

	if err != nil {
        slog.Error("Error creating Cassandra session:", "error", err)
		log.Fatal("Exiting due to error:", err) 
		return nil
    }
	db.Conn = session
	return &db
}

func (db *DB) Close() {
	db.Conn.Close()
}