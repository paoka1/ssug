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

// GenValue 生成短链
// 参数：原始链接
// 返回值：短链
func GenValue(key string) string {
	timesTotal := 1
	times := 1
	nowLen := initLen
	initKey := key
	md5 := utils.CalculateMD5(key)
	ok := data.Redirect.HasValue(md5[:nowLen])
	for ok {
		key += "z"
		if times > 3 {
			nowLen++
			times = 0
			key = initKey
		}
		md5 = utils.CalculateMD5(key)
		ok = data.Redirect.HasValue(md5[:nowLen])
		times++
		timesTotal++
	}
	if base.Debug {
		utils.Logger.Info(fmt.Sprintf("生成%s短链%s，长度%d，消耗次数%d", key, md5[:nowLen], nowLen, timesTotal))
	}
	return md5[:nowLen]
}
