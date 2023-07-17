package anvil

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/sampiiiii-dev/anvil_server/anvil/config"
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
	c config.Config
	s *zap.Logger
}

func (a *Anvil) Run() {
	makeBanner(a.c.Server.Version, a.c.Server.BannerWebsite)
	a.e.Logger.Fatal(a.e.Start(a.c.Server.Address))
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
	c := config.LoadConfig(s)

	// Echo
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.OFF)

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middlewares.ScribeLogger(s)) // Pass the logger to the middleware
	// Routes
	e.GET("/", hello)

	return &Anvil{e, c, s}
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
