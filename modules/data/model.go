package data

// Key 原始链接
// Value 短链
// ExpirationTime 过期时间（秒）
type mapping struct {
	Key            string
	Value          string
	ExpirationTime int64
}
