package app

import (
	"b3lb/pkg/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) healthCheck(c *gin.Context) {
	status := &api.HealtCheck{
		Return_code: "SUCCESS",
		Version:     "2.0",
	}

	c.XML(http.StatusOK, status)
}

func (s *Server) getMeetings(c *gin.Context) {
	c.String(http.StatusOK, c.FullPath())
}
