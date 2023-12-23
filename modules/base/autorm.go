package base

import (
	"fmt"
	"ssug/internal/base"
	"ssug/internal/utils"
	"ssug/modules/data"
	"time"
)

// AutoRemove 每隔一段时间移除过期的短链
func AutoRemove() {
	for {
		time.Sleep(10 * time.Second)
		remove()
	}
}

// RemoveExp 移除数据库中过期的短链
func RemoveExp() {
	utils.Logger.Info("尝试移除数据库中已过期的短链...")
	remove()
	utils.Logger.Info("移除过期数据完成！")
}

func remove() {
	tn := time.Now().Unix()
	for ou, su := range data.Redirect.GetCacheMappingKV() {
		t := data.Redirect.GetCacheTimeMapping(su)
		if tn >= t && ou != "" {
			_, _ = data.Redirect.RemoveRCacheMapping(ou)
			_, _ = data.Redirect.RemoveTCacheMapping(su)
			if base.Debug {
				utils.Logger.Info(fmt.Sprintf("移除缓存映射%s -> %s，存活时间结束", su, ou))
			}
		}
	}
	dataRm := data.Redirect.RemovingDBMapping(tn)
	for _, m := range dataRm {
		if base.Debug {
			utils.Logger.Info(fmt.Sprintf("移除数据库映射%s -> %s，存活时间结束", m.ShortURL, m.OriginalURL))
		} else {
			utils.Logger.Info(fmt.Sprintf("移除映射%s -> %s，存活时间结束", m.ShortURL, m.OriginalURL))
		}
	}
}
