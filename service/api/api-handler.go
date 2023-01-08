package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.Session))
	rt.router.GET("/userids/:Username", rt.wrap(rt.GetUserID))
	rt.router.GET("/users/:UserID/profile", rt.wrap(rt.GetUserProfile))
	rt.router.PUT("/users/:UserID/name", rt.wrap(rt.ChangeName))
	rt.router.PUT("/users/:UserID/followed/:OtherUserID", rt.wrap(rt.FollowUser))
	rt.router.DELETE("/users/:UserID/followed/:OtherUserID", rt.wrap(rt.UnfollowUser))
	rt.router.PUT("/users/:UserID/banned/:OtherUserID", rt.wrap(rt.BanUser))
	rt.router.DELETE("/users/:UserID/banned/:OtherUserID", rt.wrap(rt.UnbanUser))
	rt.router.POST("/photos", rt.wrap(rt.uploadPhoto))
	rt.router.PUT("/photos/:PhotoID/likes/:UserID", rt.wrap(rt.LikePhoto))
	rt.router.POST("/photos/:PhotoID/comments", rt.wrap(rt.CommentPhoto))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
