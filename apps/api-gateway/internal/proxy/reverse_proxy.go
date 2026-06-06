package proxy

import (
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewReverseProxy(target string, stripPrefix string) (gin.HandlerFunc, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(c *gin.Context) {
		originalPath := c.Request.URL.Path
		newPath := strings.TrimPrefix(originalPath, stripPrefix)

		if newPath == "" {
			newPath = "/"
		}

		c.Request.URL.Scheme = targetURL.Scheme
		c.Request.URL.Host = targetURL.Host
		c.Request.URL.Path = newPath
		c.Request.Host = targetURL.Host

		proxy.ServeHTTP(c.Writer, c.Request)
	}, nil
}
