package app

import (
	"b3lb/pkg/api"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) ChecksumValidation(c *gin.Context) {
	error := api.DefaultChecksumError()

	checksum_param, exists := c.GetQuery("checksum")
	if !exists {
		c.XML(http.StatusOK, error)
		c.Abort()
		return
	}

	params := c.Request.URL.Query()
	params.Del("checksum")

	checksum := &Checksum{
		Secret: s.Config.BigBlueButton.Secret,
		Action: strings.TrimPrefix(c.FullPath(), "/bigbluebutton/api/"),
		Params: params,
	}

	sha, err := StringToSHA1(checksum.Value())

	if err != nil {
		panic(err)
	}

	if checksum_param != string(sha) {
		c.XML(http.StatusOK, error)
		c.Abort()
		return
	}

	c.Next()
}
