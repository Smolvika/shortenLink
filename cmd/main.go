package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"shortenLink"
	"shortenLink/pkg/handler"
	"shortenLink/pkg/repository/inMemory"
	"shortenLink/pkg/repository/postgres"
	"shortenLink/pkg/service"
	"syscall"
	"time"
)

func main() {
	var PostgresFlag bool
	flag.BoolVar(&PostgresFlag, "db", false, "Run with the postgres database")
	flag.Parse()

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	startApp(PostgresFlag)

}

func startApp(postgresFlag bool) {
	var repos service.Repository
	if postgresFlag {
		db, err := postgres.NewPostgresDB(postgres.Config{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			Username: viper.GetString("db.username"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
			Password: os.Getenv("DB_PASSWORD"),
		})
		if err != nil {
			log.Fatalf("error initializing db: %s", err.Error())
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("Failed to close database:%v\n", err)
			}
		}()
		repos = postgres.New(db)
		if err != nil {
			log.Fatalf("failed to initialize db: %s", err.Error())
		}
	} else {
		const mapLen = 100
		repos = inMemory.New(mapLen)
	}

	services := service.New(repos)
	handlers := handler.New(services)
	srv := new(shortenLink.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().Do(func() {
		year, month, day := time.Now().Date()
		err := services.Delete(fmt.Sprintf("%d-%02d-%d", year, month, day))
		if err != nil {
			log.Println("error delete")
		}

	})
	s.StartAsync()

	log.Print("notebookApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print(" Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
