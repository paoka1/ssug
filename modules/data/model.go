package data

// Mapping 映射模型
// ShortURL 短链
// OriginalURL 原始链接
// ExpirationTime 过期时间（秒）
type Mapping struct {
	ShortURL       string `json:"short_url"`
	OriginalURL    string `json:"original_url"`
	ExpirationTime int64  `json:"expiration_time"`
}
