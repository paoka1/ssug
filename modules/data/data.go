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
		utils.Logger.Debug("数据库路径：" + r.db.path)
	}
}

// GetKey 获取 accessKey
func (r *redirect) GetKey() string {
	return r.accessKey
}

// AddMapping 添加新的映射
func (r *redirect) AddMapping(originalURL string, shortURL string) (Mapping, error) {
	r.l.Lock()
	defer r.l.Unlock()
	if base.Debug {
		defer r.PrintRedirect()
	}
	su, ok := r.redirectCache[originalURL]
	if ok {
		t, _ := r.timeExpirationCache[su]
		return Mapping{su, originalURL, t}, errors.New("短链映射已存在")
	}
	// t 短链过期时间
	t := time.Now().Unix() + r.ttl
	err := r.db.addMapping(Mapping{shortURL, originalURL, t})
	if err != nil {
		m, _ := r.db.getMappingByS(shortURL)
		return m, err
	}
	r.timeExpirationCache[shortURL] = t
	r.redirectCache[originalURL] = shortURL
	return Mapping{shortURL, originalURL, t}, nil
}

// RemoveRCacheMapping 去除缓存 redirectCache 映射，返回删除的映射的短链
func (r *redirect) RemoveRCacheMapping(originalURL string) (string, bool) {
	r.l.Lock()
	defer r.l.Unlock()
	su, ok := r.redirectCache[originalURL]
	if ok {
		delete(r.redirectCache, originalURL)
		return su, true
	} else {
		return "", false
	}
}

// RemoveTCacheMapping 去除缓存 timeExpirationCache 映射，返回删除的映射的过期时间
func (r *redirect) RemoveTCacheMapping(shortURL string) (int64, bool) {
	r.l.Lock()
	defer r.l.Unlock()
	t, ok := r.timeExpirationCache[shortURL]
	if ok {
		delete(r.timeExpirationCache, shortURL)
		return t, true
	} else {
		return 0, false
	}
}

// RemovingDBMapping 去除数据库里过期的的映射
func (r *redirect) RemovingDBMapping(time int64) []Mapping {
	r.l.Lock()
	defer r.l.Unlock()
	data, _ := r.db.getRemove(time)
	err := r.db.autoRemove(time)
	if err != nil {
		utils.Logger.Warning(err)
	}
	return data
}

// GetMappingFO 通过原始链接获取映射
func (r *redirect) GetMappingFO(originalURL string) (Mapping, error) {
	r.l.Lock()
	defer r.l.Unlock()
	su, ok := r.redirectCache[originalURL]
	if ok {
		t, _ := r.timeExpirationCache[su]
		return Mapping{su, originalURL, t}, nil
	}
	m, err := r.db.getMappingByO(originalURL)
	return m, err
}

// GetMappingO 获取短链对应的原始链接
func (r *redirect) GetMappingO(shortURL string) (string, error) {
	r.l.Lock()
	defer r.l.Unlock()
	if base.Debug {
		defer r.PrintRedirect()
	}
	for ou, su := range r.redirectCache {
		if su == shortURL {
			if base.Debug {
				utils.Logger.Debug(fmt.Sprintf("从缓存查找到%s -> %s", su, ou))
			}
			return ou, nil
		}
	}
	m, err := r.db.getMappingByS(shortURL)
	if err != nil {
		return "", err
	} else {
		// 添加缓存
		r.timeExpirationCache[m.ShortURL] = m.ExpirationTime
		r.redirectCache[m.OriginalURL] = m.ShortURL
		if base.Debug {
			utils.Logger.Debug(fmt.Sprintf("从数据库查找到%s -> %s", m.ShortURL, m.OriginalURL))
		}
		return m.OriginalURL, nil
	}
}

// GetMappingS 获取原始链接对应的短链
func (r *redirect) GetMappingS(originalURL string) (string, error) {
	r.l.Lock()
	defer r.l.Unlock()
	su, ok := r.redirectCache[originalURL]
	if ok {
		return su, nil
	}
	m, err := r.db.getMappingByO(originalURL)
	if err != nil {
		return "", err
	}
	return m.ShortURL, nil
}

// HasOriginalURL 是否存在原始链接
func (r *redirect) HasOriginalURL(originalURL string) bool {
	r.l.Lock()
	defer r.l.Unlock()
	_, ok := r.redirectCache[originalURL]
	if ok {
		return true
	}
	return r.db.hasOriginalURL(originalURL)
}

// HasShortURL 是否存在短链
func (r *redirect) HasShortURL(shortURL string) bool {
	r.l.Lock()
	defer r.l.Unlock()
	for _, su := range r.redirectCache {
		if su == shortURL {
			return true
		}
	}
	return r.db.hasShortURL(shortURL)
}

// GetCacheMappingKV 获取缓存中所有 KV，返回值：原始链接 -> 短链
func (r *redirect) GetCacheMappingKV() map[string]string {
	r.l.Lock()
	defer r.l.Unlock()
	kv := make(map[string]string)
	for ou, su := range r.redirectCache {
		kv[ou] = su
	}
	return kv
}

// GetCacheTimeMapping 获取缓存中短链过期的时间
func (r *redirect) GetCacheTimeMapping(shortURL string) int64 {
	r.l.Lock()
	defer r.l.Unlock()
	v, _ := r.timeExpirationCache[shortURL]
	return v
}

// PrintRedirect 打印 redirect 当前的状态
func (r *redirect) PrintRedirect() {
	utils.Logger.Debug("短链缓存：", r.redirectCache)
	utils.Logger.Debug("到期时间缓存：", r.timeExpirationCache)
}
