package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gookit/config/v2"
	"github.com/rs/zerolog/log"
	"net"
	"time"
)

func main() {
	loadConfig()
	setupLogger()
	app := createFiberApp()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.
			Status(fiber.StatusOK).
			SendString("Hello, world!")
	})

	defer func() {
		if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
			log.Error().Err(err).Msgf("Error while shutting down HTTP server")
		} else {
			log.Info().Msgf("HTTP server has stopped")
		}
	}()

	runInBackground(func() {
		listener, err := createListener()
		if err != nil {
			log.Error().Err(err).Msgf("Error while starting network listener")
			return
		}

		log.Info().Msgf("HTTP server has started")

		if err = app.Listener(listener); err != nil {
			log.Error().Err(err).Msgf("Error while starting HTTP server")
			return
		}
	})
}

func createListener() (net.Listener, error) {
	tlsCert := config.String("server.tls.cert")
	tlsKey := config.String("server.tls.key")

	switch {
	case tlsCert != "" && tlsKey != "":
		cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
		if err != nil {
			return nil, err
		}

		return tls.Listen(
			config.String("server.network", "tcp4"),
			config.String("server.address", "0.0.0.0:8080"),
			&tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		)
	default:
		return net.Listen(
			config.String("server.network", "tcp4"),
			config.String("server.address", "0.0.0.0:8080"),
		)
	}
}

func createFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:            errorHandler,
		ReadTimeout:             5 * time.Second,
		WriteTimeout:            10 * time.Second,
		IdleTimeout:             2 * time.Minute,
		Concurrency:             256 * 1024,
		BodyLimit:               4 * 1024 * 1024,
		ReadBufferSize:          4096,
		WriteBufferSize:         4096,
		DisableStartupMessage:   true,
		EnablePrintRoutes:       false,
		EnableIPValidation:      false,
		EnableTrustedProxyCheck: true,
		ProxyHeader:             "X-Forwarded-For",
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		TrustedProxies: []string{
			"127.0.0.0/8",
			"10.0.0.0/8",
			"172.16.0.0/12",
			"192.168.0.0/16",
		},
	})

	app.Use(
		recover.New(recover.Config{
			StackTraceHandler: panicHandler,
		}),
		securityHeadersHandler,
	)

	return app
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
	}

	c.Status(code)
	return nil
}

func panicHandler(_ *fiber.Ctx, r any) {
	log.Error().
		Stack().
		Err(fmt.Errorf("%v", r)).
		Msg("Panic while handling HTTP request")
}

func securityHeadersHandler(c *fiber.Ctx) error {
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-XSS-Protection", "0")

	if c.Protocol() == "https" {
		c.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	}

	return c.Next()
}
