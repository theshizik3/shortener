package httpapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/theshizik3/shortener/internal/shortener"
)

type Controller struct {
	s *shortener.Service
}

func NewController(s *shortener.Service) *Controller {
	return &Controller{s: s}
}

func (ctrl *Controller) Redirect(c *gin.Context) {
	code := c.Param("code")
	err := ctrl.s.ValidateCode(code)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "error": err.Error()})
		return
	}
	url, err := ctrl.s.GetUrl(context.Background(), code)
	if err != nil {
		c.JSON(404, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.Redirect(302, url)
}

func (ctrl *Controller) DeleteLink(c *gin.Context) {
	code := c.Param("code")
	err := ctrl.s.ValidateCode(code)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "error": err.Error()})
		return
	}
	err = ctrl.s.DeleteUrl(context.Background(), code)
	if err != nil {
		c.JSON(404, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.Status(204)
}

func (ctrl *Controller) GetClicks(c *gin.Context) {
	code := c.Param("code")
	err := ctrl.s.ValidateCode(code)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "error": err.Error()})
		return
	}
	clicks, err := ctrl.s.GetClicks(context.Background(), code)
	if err != nil {
		c.JSON(404, gin.H{"success": false, "error": err})
		return
	}
	c.JSON(200, gin.H{"success": true, "clicks": clicks})
}

func (ctrl *Controller) MakeLink(c *gin.Context) {
	u, err := url.ParseRequestURI(c.PostForm("url"))
	if err != nil {
		c.JSON(400, gin.H{"success": false, "error": "incorrect url"})
		return
	}
	if u.Scheme == "" || u.Host == "" {
		c.JSON(400, gin.H{"success": false, "error": "incorrect url"})
		return
	}
	code, err := ctrl.s.SetUrl(context.Background(), c.PostForm("url"))
	if err != nil {
		c.JSON(500, gin.H{"success": false, "error": "internal error"})
		fmt.Println(err)
		return
	}
	c.JSON(201, gin.H{"success": true, "code": code})
}
