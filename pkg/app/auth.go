package app

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/SLedunois/b3lb/v2/pkg/api"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func processParameters(query string) string {
	reg := regexp.MustCompile(`&?checksum(=[^&]*&?)`)
	return reg.ReplaceAllString(query, "")
}

// ChecksumValidation handler validate all requests checksum and returns an error if the checksum is not int the request or if the checksum is invalid
func (s *Server) ChecksumValidation(c *gin.Context) {
	error := api.DefaultChecksumError()

	checksumParam, exists := c.GetQuery("checksum")
	if !exists {
		c.XML(http.StatusOK, error)
		c.Abort()
		return
	}

	secret := s.Config.BigBlueButton.Secret
	tenant, err := s.TenantManager.GetTenant(c.Request.Host)
	if err != nil {
		log.Error("Tenant manager can't retrieve tenant: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if tenantSecret, exists := tenant.Spec["secret"]; exists {
		secret = tenantSecret
	}

	checksum := &api.Checksum{
		Secret: secret,
		Action: strings.TrimPrefix(c.FullPath(), "/bigbluebutton/api/"),
		Params: processParameters(c.Request.URL.RawQuery),
	}

	sha, err := checksum.Process()

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if checksumParam != string(sha) {
		c.XML(http.StatusOK, error)
		c.Abort()
		return
	}

	setAPIContext(c, checksum)

	c.Next()
}

func setAPIContext(c *gin.Context, checksum *api.Checksum) {
	c.Set("api_ctx", checksum)
}

func getAPIContext(c *gin.Context) *api.Checksum {
	return c.MustGet("api_ctx").(*api.Checksum)
}
