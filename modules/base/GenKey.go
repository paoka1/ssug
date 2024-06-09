package base

import (
	"fmt"
	"ssug/internal/base"
	"ssug/internal/utils"
	"ssug/modules/data"
)

var initLen int

// SetInitLen 设置最短短链长度
func SetInitLen(len int) {
	if len <= 0 {
		initLen = 3
		utils.Logger.Info("短链最短长度不合法，将使用默认值")
	} else {
		initLen = len
	}
}

// GenValue 根据原始链接生成短链
func GenValue(originalURL string) string {
	timesTotal := 1
	times := 1
	nowLen := initLen
	initKey := originalURL
	md5 := utils.CalculateMD5(originalURL)
	ok := data.Redirect.HasShortURL(md5[:nowLen])
	for ok {
		originalURL += "z"
		if times > 3 {
			nowLen++
			times = 0
			originalURL = initKey
		}
		md5 = utils.CalculateMD5(originalURL)
		ok = data.Redirect.HasShortURL(md5[:nowLen])
		times++
		timesTotal++
	}
	if base.Debug {
		utils.Logger.Debug(fmt.Sprintf("生成%s，短链：%s，长度：%d，消耗次数：%d", originalURL, md5[:nowLen], nowLen, timesTotal))
	}
	return md5[:nowLen]
}
