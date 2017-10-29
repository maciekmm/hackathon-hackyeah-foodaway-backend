package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/maciekmm/HackYeah/controller"
	"github.com/maciekmm/HackYeah/model"
	"github.com/rs/cors"
)

type Application struct {
	Database *gorm.DB
	Logger   *log.Logger
	router   *mux.Router
}

func (a *Application) Init() error {
	a.Logger.Println("starting HackYeah")

	a.Logger.Println("setting up database connection")
	con, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@database/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB")))
	if err != nil {
		return fmt.Errorf("could not open database connection: %s", err.Error())
	}

	a.Logger.Println("establishing database connection")
	deadline := time.After(10 * time.Second)
out:
	for {
		select {
		case <-deadline:
			return fmt.Errorf("could not establish database connection, last error: %s", err.Error())
		default:
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			err = con.PingContext(ctx)
			if err == nil {
				break out
			} else {
				a.Logger.Printf("pinging database failed: %s\n", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}
	// auto-migrating model
	a.Database, err = gorm.Open("postgres", con)
	a.Database.SetLogger(a.Logger)
	if err != nil {
		return err
	}

	a.Database.AutoMigrate(&model.User{}, &model.Offer{})

	// setup routes
	a.Logger.Println("setting up routes")
	a.router = mux.NewRouter()

	// shitty hack to counter browser CORS clients, it's hackathon all tricks are allowed
	a.router.Methods("OPTIONS").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	// accounts
	accountController := &controller.Accounts{Database: a.Database}
	accountController.Register(a.router.PathPrefix("/accounts/").Subrouter())

	// offers
	offerController := &controller.Offers{Database: a.Database}
	offerController.Register(a.router.PathPrefix("/offers/").Subrouter())
	return nil
}

func (a *Application) Serve() error {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
	})
	server := http.Server{
		Addr:           ":3000",
		Handler:        c.Handler(a.router),
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return server.ListenAndServe()
}
