package main

import (
	"flag"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"shortenLink"
	"shortenLink/pkg/handler"
	"shortenLink/pkg/repository/inMemory"
	"shortenLink/pkg/repository/postgres"
	"shortenLink/pkg/service"
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
	var repos service.Repository
	if PostgresFlag {
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

	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error running http serever: %s", err.Error())
	}
	//go func() {
	//	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
	//		log.Fatalf("error occured while running http server: %s", err.Error())
	//	}
	//}()

	//log.Print("Started")

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
