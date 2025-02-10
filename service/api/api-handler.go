package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	rt.router.PUT("/users/:id/username", rt.wrap(rt.setMyUsername))
	rt.router.GET("/users/:id/username", rt.wrap(rt.getUsername))

	rt.router.PUT("/users/:id/photo", rt.wrap(rt.setMyPhoto))
	rt.router.GET("/users/:id/photo", rt.wrap(rt.getPhoto))

	rt.router.GET("/chats", rt.wrap(rt.getMyConversations))

	rt.router.GET("/chats/:chatId", rt.wrap(rt.getConversation))
	rt.router.POST("/chats/:chatId", rt.wrap(rt.sendMessage))

	rt.router.POST("/chats/:chatId/messages/:messageId", rt.wrap(rt.forwardMessage))
	rt.router.GET("/chats/:chatId/messages/:messageId", rt.wrap(rt.getMessage))
	rt.router.DELETE("/chats/:chatId/messages/:messageId", rt.wrap(rt.deleteMessage))

	rt.router.GET("/chats/:chatId/messages/:messageId/photo", rt.wrap(rt.getMessagePhoto))

	rt.router.POST("/chats/:chatId/messages/:messageId/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/chats/:chatId/messages/:messageId/comments", rt.wrap(rt.uncommentMessage))

	rt.router.PUT("/chats/:chatId/chatName", rt.wrap(rt.setGroupName))
	rt.router.GET("/chats/:chatId/chatName", rt.wrap(rt.getGroupName))

	rt.router.PUT("/chats/:chatId/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.GET("/chats/:chatId/photo", rt.wrap(rt.getGroupPhoto))

	rt.router.PUT("/chats/:chatId/members", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/chats/:chatId/members", rt.wrap(rt.leaveGroup))

	// Added
	rt.router.PUT("/newchat", rt.wrap(rt.newChat))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
