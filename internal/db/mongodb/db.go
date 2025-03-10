package mongodb

import (
	"context"
	"log"
	"log/slog"
	"math/rand"
	"strconv"
	"time"

	"github.com/tanjed/go-sso/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const DB_DRIVER_POSTGRES = "postgres"

type DB struct {
	id int
	Conn *mongo.Client
}


func NewDB(c *config.Config) *DB{
	client, err := mongo.Connect(options.Client().ApplyURI(getConnectionString(*c)))

	if err != nil {
		slog.Error("Error creating mongo session:", "error", err)
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err:= client.Ping(ctx, nil); err != nil {
		slog.Error("Error creating mongo session:", "error", err)
		log.Fatal(err)
	}
	return &DB{
		id: rand.Intn(101),
		Conn: client,
	}
}

func (db DB)Close() {
	if db.Conn != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()
		db.Conn.Disconnect(ctx)
	}
}

func getConnectionString(c config.Config) string {
	return "mongodb://"+c.DB_USER+":"+c.DB_PASSWORD+"@"+c.DB_HOST+":"+strconv.Itoa(c.DB_PORT)
}