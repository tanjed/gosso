package db

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"sync"

	"github.com/gocql/gocql"
	"github.com/tanjed/go-sso/internal/config"
)

const DB_DRIVER_POSTGRES = "postgres"

type DB struct {
	Id int
	Conn *gocql.Session
}

var dbInstance *DB
var dbInstanceOnce sync.Once

func InitDB() *DB{

	dbInstanceOnce.Do(func ()  {
		cluster := gocql.NewCluster(config.AppConfig.DB_HOST)
		cluster.Port = config.AppConfig.DB_PORT
		cluster.Keyspace = config.AppConfig.DB_NAME     
		cluster.NumConns = 500
		cluster.Consistency = gocql.Quorum
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: config.AppConfig.DB_USER,
			Password: config.AppConfig.DB_PASSWORD,
		}

		session, err := cluster.CreateSession()

		if err != nil {
			slog.Error("Error creating Cassandra session:", "error", err)
		}else {
			dbInstance = &DB{
				Id: rand.Intn(101),
				Conn: session,
			}
		}
			
	})

	if dbInstance.Conn == nil {
		log.Fatal("Unable to establish DB connection") 
	}
	fmt.Println("DB CONNECTION ID:", dbInstance.Id)
	return dbInstance
}

func CloseDB() {
	if dbInstance != nil && dbInstance.Conn != nil {
		dbInstance.Conn.Close()
	}
}