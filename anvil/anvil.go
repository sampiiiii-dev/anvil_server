package anvil

import (
	"context"
	"errors"
	"fmt"
	"github.com/sampiiiii-dev/anvil_server/anvil/resources"
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
	e  *echo.Echo
	c  *config.Config
	s  *zap.Logger
	db *db.Database
	ks *db.Keystore
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
	if err := a.ks.Shutdown(); err != nil {
		a.s.Error("Failed to close Redis client:", zap.Error(err))
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

func Forge() *Anvil {
	s := logs.HireScribe()

	// Load configuration
	c := config.GetConfigInstance(s)

	// Echo & Echo Config
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.OFF)

	// Keystore (Redis)
	ctx := context.Background()
	ks := db.NewKeystore(ctx)

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middlewares.ScribeLogger(s)) // Pass the logger to the middleware
	// Routes
	e.GET("/", hello)
	resources.RegisterUserRoutes(e, db.GetDBInstance())

	return &Anvil{e: e, c: c, s: s, ks: ks}
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
