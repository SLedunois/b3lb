package app

import (
	"b3lb/pkg/api"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// HealthCheck handler returns an health check response
func (s *Server) HealthCheck(c *gin.Context) {
	c.XML(http.StatusOK, api.CreateHealthCheck())
}

// GetMeetings handler returns the getMeetings API. See https://docs.bigbluebutton.org/dev/api.html#getmeetings.
func (s *Server) GetMeetings(c *gin.Context) {
	c.String(http.StatusOK, c.FullPath())
}

// Create handler find a server and create a meeting on balanced server.
func (s *Server) Create(c *gin.Context) {
	ctx := getAPIContext(c)
	instances, err := s.InstanceManager.List()
	if err != nil {
		log.Error("Manager failed to retrieve instance list", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(instances) == 0 {
		log.Error("InstanceManager does not retrieve any instances. Please check you add at least one Bigbluebutton instance")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	target, err := s.Balancer.Process(instances)
	if err != nil {
		log.Error("Balancer failed to process current request", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	instance, err := s.InstanceManager.Get(target)
	if err != nil {
		log.Error("Manager failed to retrieve target instance for current request", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	apiResponse := instance.Create(ctx.Params)

	if apiResponse == nil {
		log.Error("An error occurred while creating remote session")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	addErr := s.SessionManager.Add(apiResponse.MeetingID, instance.URL)
	if addErr != nil {
		log.Error("SessionManager failed to add new session", addErr)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.XML(http.StatusOK, apiResponse)
}
