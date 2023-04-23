package main

import (
	"context"
	"flag"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"shortenLink"
	"shortenLink/pkg/handler"
	"shortenLink/pkg/repository"
	"shortenLink/pkg/repository/postgres"
	"shortenLink/pkg/service"
	"syscall"
	"time"
)

func main() {
	var imFlag, dbFlag bool
	flag.BoolVar(&dbFlag, "db", false, "Run with the postgres")
	flag.BoolVar(&imFlag, "im", false, "Run with the in-memory")
	flag.Parse()

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	if dbFlag {
		startApp("db")
	} else if imFlag {
		startApp("im")
	} else {
		log.Fatal("This flag does not exist")
	}

}

func startApp(flag string) {
	var repos *repository.Repository
	if flag == "db" {
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
		repos = repository.NewDB(db)
		if err != nil {
			log.Fatalf("failed to initialize db: %s", err.Error())
		}
	} else {
		repos = repository.NewIm()
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
		date := time.Now().Format("2006-01-02")
		err := services.Delete(date)
		if err != nil {
			log.Println(err)
		}
	})
	s.StartAsync()

	log.Print("Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
