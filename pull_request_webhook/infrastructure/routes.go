package infrastructure

import "github.com/gin-gonic/gin"

func Routes(router *gin.Engine) {

	// https://38a1-189-150-46-82.ngrok-free.app/pull_request/process

	routes := router.Group("pull_request")

	{
		routes.POST("/process", HandlePullRequestEvent)
	}

}
