package helpers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DecodeReqBody (c *gin.Context, reqBody any) bool {
	if err := c.BindJSON(reqBody); err != nil {
		log.Println("Failed to decode request body:", err)
		c.JSON(http.StatusBadRequest, err)
		return false
	}

	return true;
}