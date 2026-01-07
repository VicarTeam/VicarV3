package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"vicar-backend/auth"
	"vicar-backend/cache"
	"vicar-backend/controllers"
	"vicar-backend/db"
	"vicar-backend/log"
	"vicar-backend/sync"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go start()

	<-ctx.Done()
}

func start() {
	if err := db.Initialize(); err != nil {
		panic(err)
	}

	if err := auth.InitializeWebauthn(); err != nil {
		log.Error(log.Auth, "❌", "Webauthn initialization failed: %v", err)
		return
	}

	if err := cache.Initialize(); err != nil {
		log.Error(log.Cache, "❌", "Cache initialization failed: %v", err)
		return
	}

	app := fiber.New(fiber.Config{
		ErrorHandler:          handleErrors,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_ALLOW_ORIGINS"),
		AllowMethods:     "GET,POST,PATCH,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, User-Agent, X-Device-Name, X-Csrf-Token",
		AllowCredentials: true,
	}))

	app.Use(helmet.New())
	app.Use(logger.New())

	v1 := app.Group("/v1")

	controllers.Register(v1)
	sync.Register(v1)

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}

func handleErrors(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	codeMessage := "Internal Server Error"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code

		log.Error(log.Server, "❌", "Handling request failed: %v", e.Message)
	} else if db.IsRecordNotFound(err) {
		code = fiber.StatusNotFound
		codeMessage = "Entity not found"
	} else {
		log.Error(log.Server, "❌", "Handling request failed: %v", err)
	}

	if code == fiber.StatusNotFound {
		codeMessage = "Not Found"
	} else if code != fiber.StatusInternalServerError {
		return ctx.SendStatus(code)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"error": codeMessage,
	})
}
