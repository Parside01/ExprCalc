package controllers

import "github.com/labstack/echo/v4"

type Controller interface {
	GetGroup() string
	GetHandlers() []ControllerHandler
	GetMiddleware() []echo.MiddlewareFunc
}

type ControllerHandler interface {
	GetMethod() string
	GetPath() string
	GetHandler() func(echo.Context) error
}

type Handler struct {
	Method  string
	Path    string
	Handler func(echo.Context) error
}

func (c *Handler) GetMethod() string {
	return c.Method
}

func (c *Handler) GetPath() string {
	return c.Path
}

func (c *Handler) GetHandler() func(echo.Context) error {
	return c.Handler
}
