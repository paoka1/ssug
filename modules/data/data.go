package data

import (
	"errors"
	"fmt"
	"ssug/internal/base"
	"ssug/internal/utils"
	"sync"
	"time"
)

var (
	Redirect redirect
)

// redirectCache 原始链接 -> 短链
// timeExpirationCache 短链 -> 过期时间
// ttl 短链存活时长，单位秒
type redirect struct {
	l                   *sync.Mutex
	accessKey           string
	redirectCache       map[string]string
	timeExpirationCache map[string]int64
	ttl                 int64
	db                  database
}

// Init 创建数据结构，初始化数据库
func (r *redirect) Init(key string, ttl int64) {
	r.accessKey = key
	r.ttl = ttl
	r.redirectCache = make(map[string]string)
	r.timeExpirationCache = make(map[string]int64)
	r.l = &sync.Mutex{}
	r.db = getDatabase()
	d, err := r.db.open()
	if err != nil {
		utils.Logger.Fatal(err)
	}
	r.db.db = d
	utils.Logger.Info("成功加载数据库")
	if base.Debug {
		utils.Logger.Info("数据库路径：" + r.db.path)
	}
}

// GetKey 获取 accessKey
func (r *redirect) GetKey() string {
	return r.accessKey
}

func (r *redirect) Close() {
	r.db.close()
}

// AddMapping 添加新的映射
// 参数：key 原始链接，value 短链
// 返回值：err，原始链接，短链
func (r *redirect) AddMapping(key string, value string) (error, string, string) {
	r.l.Lock()
	defer r.l.Unlock()
	v, ok := r.redirectCache[key]
	if ok {
		return errors.New("短链映射已存在"), key, v
	}
	// t 短链过期时间
	t := time.Now().Unix() + r.ttl
	err := r.db.addMapping(t, key, value)
	if err != nil {
		return err, key, v
	}
	r.timeExpirationCache[value] = t
	r.redirectCache[key] = value
	return nil, key, value
}

// RemoveRCacheMapping 去除缓存 redirectCache 映射
// 参数 key：原始链接
// 返回值：操作是否成功，原始链接，短链
func (r *redirect) RemoveRCacheMapping(key string) (bool, string, string) {
	r.l.Lock()
	defer r.l.Unlock()
	v, ok := r.redirectCache[key]
	if ok {
		delete(r.redirectCache, key)
		return true, key, v
	} else {
		return false, key, ""
	}
}

// RemoveTCacheMapping 去除缓存 timeExpirationCache 映射
// 参数 key：短链
// 返回值：操作是否成功，短链，短链过期时间
func (r *redirect) RemoveTCacheMapping(key string) (bool, string, int64) {
	r.l.Lock()
	defer r.l.Unlock()
	v, ok := r.timeExpirationCache[key]
	if ok {
		delete(r.timeExpirationCache, key)
		return true, key, v
	} else {
		return false, key, 0
	}
}

// RemovingDBMapping 去除数据库里过期的的映射
// 参数：time 应为操作时的时间戳
// 返回值：删除的映射（原始链接 -> 短链）
func (r *redirect) RemovingDBMapping(time int64) map[string]string {
	data, _ := r.db.getRemove(time)
	_ = r.db.autoRemove(time)
	ret := make(map[string]string)
	for _, m := range data {
		ret[m.Key] = m.Value
	}
	return ret
}

// GetMappingKey 获取短链对应的原始链接
// 参数：短链
// 返回值：err，原始链接，短链
func (r *redirect) GetMappingKey(value string) (error, string, string) {
	r.l.Lock()
	defer r.l.Unlock()
	for k, v := range r.redirectCache {
		if v == value {
			if base.Debug {
				utils.Logger.Info(fmt.Sprintf("从缓存查找到%s -> %s", k, v))
			}
			return nil, k, v
		}
	}
	err, m := r.db.getMappingByV(value)
	if err != nil {
		return err, "", value
	} else {
		r.timeExpirationCache[m.Value] = m.ExpirationTime
		r.redirectCache[m.Key] = m.Value
		if base.Debug {
			utils.Logger.Info(fmt.Sprintf("从数据库查找到%s -> %s", m.Key, m.Value))
		}
		return nil, m.Key, value
	}
}

// GetMappingValue 获取原始链接对应的短链
// 参数：原始链接
// 返回值：err，原始链接，短链
func (r *redirect) GetMappingValue(key string) (error, string, string) {
	r.l.Lock()
	defer r.l.Unlock()
	v, ok := r.redirectCache[key]
	if ok {
		return nil, key, v
	}
	err, m := r.db.getMappingByK(key)
	if err != nil {
		return err, key, ""
	}
	return nil, key, m.Value
}

// HasKey 是否存在原始链接
// 参数：原始链接
// 返回值：bool
func (r *redirect) HasKey(key string) bool {
	r.l.Lock()
	defer r.l.Unlock()
	_, ok := r.redirectCache[key]
	if ok {
		return true
	}
	return r.db.hasKey(key)
}

// HasValue 是否存在短链
// 参数：短链
// 返回值：bool
func (r *redirect) HasValue(value string) bool {
	r.l.Lock()
	defer r.l.Unlock()
	for _, v := range r.redirectCache {
		if v == value {
			return true
		}
	}
	ok := r.db.hasValue(value)
	return ok
}

// GetCacheMappingKV 获取缓存中所有 KV
// 返回值：原始链接 -> 短链
func (r *redirect) GetCacheMappingKV() map[string]string {
	r.l.Lock()
	defer r.l.Unlock()
	kv := make(map[string]string)
	for k, v := range r.redirectCache {
		kv[k] = v
	}
	return kv
}

// GetCacheTimeMapping 获取缓存中短链过期的时间
// 参数：短链
// 返回值：过期时间
func (r *redirect) GetCacheTimeMapping(key string) int64 {
	r.l.Lock()
	defer r.l.Unlock()
	v, _ := r.timeExpirationCache[key]
	return v
}
