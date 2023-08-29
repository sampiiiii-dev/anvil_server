package anvil

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sampiiiii-dev/anvil_server/anvil/resources"
	"github.com/sampiiiii-dev/anvil_server/anvil/workers"
	"github.com/sampiiiii-dev/anvil_server/anvil/workers/jobs"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/sampiiiii-dev/anvil_server/anvil/config"
	"github.com/sampiiiii-dev/anvil_server/anvil/db"
	"github.com/sampiiiii-dev/anvil_server/anvil/logs"
	"github.com/sampiiiii-dev/anvil_server/anvil/middlewares"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// Anvil is an instance of the server
type Anvil struct {
	e *echo.Echo
	c *config.Config
	s *zap.Logger
	d *gorm.DB
	r *redis.Client
}

var isShuttingDown int32 // 0 means not shutting down, 1 means shutting down

func (a *Anvil) Run() {
	makeBanner(a.c.Server.Version, a.c.Server.BannerWebsite)

	// Thread-safe flag for shutdown. Use this flag in your request-handling code.
	atomic.StoreInt32(&isShuttingDown, 0)

	go func() {
		if err := a.e.Start(a.c.Server.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.s.Fatal("Shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	// Setting the shutdown flag to 1 to indicate that the server is in the process of shutting down.
	atomic.StoreInt32(&isShuttingDown, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.s.Info("Cleaning up...")
	if err := db.CloseDB(); err != nil {
		a.s.Error("Failed to close database:", zap.Error(err))
	}
	a.s.Info("Database closed.")

	if err := a.r.Close(); err != nil {
		a.s.Fatal(err.Error())
	}
	a.s.Info("Redis client closed.")

	if err := a.e.Shutdown(ctx); err != nil {
		a.s.Error("Shutdown failed:", zap.Error(err))
	} else {
		a.s.Info("Shutdown successful.")
	}
}

func makeBanner(version string, website string) {
	// http://patorjk.com/software/taag/#p=display&v=0&f=Epic&t=iForge%0A
	banner := `
_________ _______  _______  _______  _______  _______ 
\__   __/(  ____ \(  ___  )(  ____ )(  ____ \(  ____ \
   ) (   | (    \/| (   ) || (    )|| (    \/| (    \/
   | |   | (__    | |   | || (____)|| |      | (__    
   | |   |  __)   | |   | ||     __)| | ____ |  __)   
   | |   | (      | |   | || (\ (   | | \_  )| (      
___) (___| )      | (___) || ) \ \__| (___) || (____/\
\_______/|/       (_______)|/   \__/(_______)(_______/
														  
%s
anvil, an iForge product
%s
____________________________________O/_______
                                    O\

`

	// Create color functions
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Print the banner with color
	fmt.Printf(banner, cyan(version), yellow(website))
}

// NewAnvil creates a new Anvil instance with injected dependencies
func NewAnvil(e *echo.Echo, c *config.Config, s *zap.Logger, d *gorm.DB, r *redis.Client) *Anvil {
	return &Anvil{
		e: e,
		c: c,
		s: s,
		d: d,
		r: r,
	}
}

func Forge() *Anvil {
	s := logs.HireScribe()

	// Load configuration
	c := config.GetConfigInstance(s)

	// Echo & Echo Config
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.OFF)

	d := db.GetDBInstance()
	r := db.InitializeRedisClient(c)

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middlewares.ScribeLogger(s)) // Pass the logger to the middleware
	// Routes
	e.GET("/", hello)
	e.POST("/test_email", test_email)
	resources.RegisterUserRoutes(e, d)

	return NewAnvil(e, c, s, d, r)
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func test_email(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" { // Assume isValidEmail validates the email
		return c.String(http.StatusBadRequest, "Invalid email")
	}

	// Create a new job
	job := &jobs.EmailJob{
		Email: email,
	}

	redisJobQueue := workers.NewRedisJobQueue(db.InitializeRedisClient(nil), db.InitializeRedisClient(nil).Context())

	// Add the job to the queue
	if err := redisJobQueue.Enqueue(job, "email"); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to enqueue job")
	}

	return c.String(http.StatusOK, "Email sent!")
}
