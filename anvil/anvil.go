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
	s *echo.Echo
	c config.Config
}

func (a *Anvil) Run() {
	makeBanner(a.c.Server.Version, a.c.Server.BannerWebsite)
	a.s.Logger.Fatal(a.s.Start(a.c.Server.Address))
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
	// Scribe
	scribe, _ := zap.NewProduction()
	defer scribe.Sync()

	// Load configuration
	c := config.LoadConfig(scribe)

	// Echo
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.OFF)
	e.Use(ZapLogger(scribe)) // Pass the logger to the middleware

	// Middleware
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	return &Anvil{s: e, c: c}
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
