package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	h "github.com/lrweck/shawty/api"
	mdb "github.com/lrweck/shawty/repository/mongodb"
	pg "github.com/lrweck/shawty/repository/postgresql"
	red "github.com/lrweck/shawty/repository/redis"
	"github.com/lrweck/shawty/shortener"
)

func main() {

	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	hand := h.NewHandler(service)

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/:code", hand.Get)
	app.Post("/", hand.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port", httpPort())
		errs <- app.Listen(httpPort())
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated: %s", <-errs)

}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func chooseRepo() shortener.RedirectRepository {

	fmt.Println("choose db")
	switch os.Getenv("URL_DB") {
	case "pg":
		pgURL := os.Getenv("PG_URL")
		pgTimeout, _ := strconv.Atoi(os.Getenv("PG_TIMEOUT"))
		repo, err := pg.NewPGRepo(pgURL, pgTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo

	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := red.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("choose REDIS")
		return repo
	case "mongo":

		mongoURL := "mongodb+srv://luis:4ulR36RwsoP9SH96@cluster0.qzt86.mongodb.net/urlshortener?retryWrites=true&w=majority" //os.Getenv("MONGO_URL")
		mongoDB := "redirects"                                                                                                //os.Getenv("MONGO_DB")
		mongoTimeout := 30                                                                                                    //, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		fmt.Println("choose MONGO")
		repo, err := mdb.NewMongoRepo(mongoURL, mongoDB, mongoTimeout)
		// fmt.Printf("Mongo URL: %s \n", mongoURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}

	fmt.Println("choose NOTHING")

	// return nil

}
