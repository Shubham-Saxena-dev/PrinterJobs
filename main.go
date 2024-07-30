package main

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"karlstorz/internal/routes"
	"karlstorz/pkg/controller"
	"karlstorz/pkg/repositories"
	"karlstorz/pkg/service"
	"sync"
)

const (
	Database   = "goMongo"
	Collection = "printers"
	MongoDbUrl = "mongodb://mongodb:27017/"
)

var (
	collection *mongo.Collection
	ctx        = context.TODO()
	once       = sync.Once{}
)

func main() {
	log.Info("Hi, Welcome !")

	initDatabase()

	repo := repositories.NewPrinterRepositories(collection, ctx)
	serv := service.NewPrinterService(repo)
	contr := controller.NewPrinterController(serv)

	r := gin.Default()

	routes.RegisterHandlers(r, contr).RegisterRoutest()

	err := r.Run()
	if err != nil {
		panic(err)
	}
}

func initDatabase() {
	once.Do(func() {
		log.Info("Connecting to datastore")
		clientOptions := options.Client().ApplyURI(MongoDbUrl)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}
		collection = client.Database(Database).Collection(Collection)
	})
}
