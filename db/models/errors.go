package models

import (
	"log"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

type RequestError struct {
	Err 		error
	StatusCode 	int
	Code 		string
}

func (e *RequestError) ErrToMap() gin.H{
	return gin.H{
		"message": e.Err,
		"code": e.Code,
	}
}

func (e *RequestError) Log(){
	red := color.New(color.FgRed).SprintFunc()
	log.Println(red(e.Code) + ":", e.Err)
}

func (e *RequestError) LogAndReturn(c *gin.Context){
	red := color.New(color.FgRed).SprintFunc()
	log.Println(red(e.Code) + ":", e.Err)
	c.JSON(e.StatusCode, e.Code)
}