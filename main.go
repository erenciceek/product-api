package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"product-api/internal/controller"
	"product-api/internal/repository"
	"product-api/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB: ", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("product_db")

	productRepo := repository.NewMongoProductRepository(db)
	productService := service.NewProductService(productRepo)
	productController := controller.NewProductController(productService)

	e := echo.New()

	e.Use(LoggerMiddleware(logger))

	productController.RegisterRoutes(e)

	go func() {
		if err := e.Start(":8080"); err != nil {
			logger.Fatal("Failed to start server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal("Failed to shutdown server: ", err)
	}
}

func LoggerMiddleware(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			fields := logrus.Fields{
				"remote_ip":  c.RealIP(),
				"host":       req.Host,
				"method":     req.Method,
				"uri":        req.RequestURI,
				"status":     res.Status,
				"latency":    time.Since(start).String(),
				"user_agent": req.UserAgent(),
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			fields["request_id"] = id

			n := res.Status
			switch {
			case n >= 500:
				logger.WithFields(fields).Error("Server error")
			case n >= 400:
				logger.WithFields(fields).Warn("Client error")
			case n >= 300:
				logger.WithFields(fields).Info("Redirect")
			default:
				logger.WithFields(fields).Info("Success")
			}

			return nil
		}
	}
}
