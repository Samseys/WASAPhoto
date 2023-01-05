package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.Session))
	rt.router.PUT("/users/:UserID/name", rt.wrap(rt.ChangeName))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
