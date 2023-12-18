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

	url := c.DefaultPostForm("url", "")
	v, err := handlers.AddMappingHandler(url)
	if err != nil {
		c.JSON(403, utils.ResultFailD(403, err.Error(), v))
		return
	} else {
		c.JSON(200, utils.ResultSuccess(v))
	}
}

func GetMapping(c *gin.Context) {
	value := c.Param("value")
	url, err := handlers.GetMappingHandler(value)
	if err != nil {
		c.JSON(403, utils.ResultFail(403, err.Error()))
		return
	} else {
		c.Redirect(302, url)
	}
}

func Happy(c *gin.Context) {
	c.JSON(200, utils.ResultSuccess("Service is running :)"))
}
