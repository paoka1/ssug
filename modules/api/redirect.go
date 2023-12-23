package api

import (
	"github.com/gin-gonic/gin"
	"ssug/internal/utils"
	"ssug/modules/data"
	"ssug/modules/handlers"
)

func AddMapping(c *gin.Context) {
	accKey := c.DefaultPostForm("key", "")
	if accKey != data.Redirect.GetKey() {
		c.JSON(401, utils.ResultFail(401, "操作未授权"))
		return
	}

	originalURL := c.DefaultPostForm("url", "")
	v, err := handlers.AddMappingHandler(originalURL)
	if err != nil {
		c.JSON(403, utils.ResultFailWD(403, err.Error(), v))
	} else {
		c.JSON(200, utils.ResultSuccess(v))
	}
}

func GetMapping(c *gin.Context) {
	shortURL := c.Param("short")
	url, err := handlers.GetMappingHandler(shortURL)
	if err != nil {
		c.JSON(403, utils.ResultFail(403, err.Error()))
	} else {
		c.Redirect(302, url)
	}
}

func Happy(c *gin.Context) {
	c.JSON(200, utils.ResultSuccess("Service is running :)"))
}
