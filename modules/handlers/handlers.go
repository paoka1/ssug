package handlers

import (
	"errors"
	"fmt"
	"ssug/internal/utils"
	"ssug/modules/base"
	"ssug/modules/data"
)

// AddMappingHandler 处理添加短链请求
// 参数：原始链接
// 返回值：短链，err
func AddMappingHandler(key string) (string, error) {
	if !utils.IsURL(key) {
		return "", errors.New("添加映射失败，URL非法")
	}
	ok := data.Redirect.HasKey(key)
	if ok {
		_, _, v := data.Redirect.GetMappingValue(key)
		utils.Logger.Info(fmt.Sprintf("添加映射%s -> %s失败，映射已存在", key, v))
		return v, errors.New("添加失败，映射已存在")
	}
	value := base.GenValue(key)
	err, k, v := data.Redirect.AddMapping(key, value)
	if err == nil {
		utils.Logger.Info(fmt.Sprintf("成功添加映射%s -> %s", k, v))
		return v, nil
	} else {
		utils.Logger.Warning(err)
		return v, err
	}
}

// GetMappingHandler 处理原始链接获取请求
// 参数：短链
// 返回值：原始链接，err
func GetMappingHandler(value string) (string, error) {
	if !utils.IsLegalValue(value) {
		return "", errors.New("查询映射失败，value非法")
	}
	err, k, _ := data.Redirect.GetMappingKey(value)
	if err == nil {
		return k, nil
	} else {
		return "", errors.New("查询失败，映射不存在")
	}
}
