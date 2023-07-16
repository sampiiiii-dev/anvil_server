package anvil

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/sampiiiii-dev/anvil_server/anvil/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

// Anvil is an instance of the server
type Anvil struct {
	Echo *echo.Echo
}

func (a *Anvil) Run(address string) {
	makeBanner()
	a.Echo.Logger.Fatal(a.Echo.Start(address))
}

func makeBanner() {
	version := `0.0.1`
	website := `https://iforgesheffield.org`
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
	// Scribe
	scribe, _ := zap.NewProduction()
	defer scribe.Sync()

	_, c_err := config.LoadConfig(".")
	if c_err != nil {
		scribe.Warn("unable to load configuration", zap.Error(c_err))
	}

	// Echo
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.OFF)
	e.Use(ZapLogger(scribe)) // Pass the logger to the middleware

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	return &Anvil{Echo: e}
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// Zap Middleware
func ZapLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			logger.Info("incoming request",
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.Int("status", res.Status),
			)
			return next(c)
		}
	}
}
