package route

import "github.com/labstack/echo/v4"

type RouteConfig struct {
	App *echo.Echo
}

func (c *RouteConfig) Setup() {
	//Setup Route
}

func (c *RouteConfig) SetupOrderRoute() {
	// Setup endpoint
}
