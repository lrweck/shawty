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
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := red.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("choose REDIS")
		return repo
	case "mongo":

		mongoURL := os.Getenv("MONGO_URL")
		mongoDB := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		fmt.Println("choose MONGO")
		repo, err := mdb.NewMongoRepo(mongoURL, mongoDB, mongoTimeout)
		fmt.Printf("Mongo URL: %s \n", mongoURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}

	fmt.Println("choose NOTHING")

	return nil

}
