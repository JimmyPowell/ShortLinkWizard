package service

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"shortlink/config"
	"shortlink/model"
	"strconv"
	"strings"
	"time"
)

const (
	base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	codeLength  = 6 // 固定长度为 6 位
)

type ShortLinkService struct{}

func NewShortLinkService() *ShortLinkService {
	return &ShortLinkService{}
}

// 创建短链接
func (s *ShortLinkService) CreateShortLink(originalURL string) (*model.ShortLink, error) {
	// 验证 URL 格式
	if !isValidURL(originalURL) {
		return nil, errors.New("invalid URL format")
	}

	// 创建记录并获取自增ID
	shortLink := &model.ShortLink{
		OriginalURL: originalURL,
	}

	if err := config.DB.Create(shortLink).Error; err != nil {
		return nil, fmt.Errorf("failed to create short link record: %w", err)
	}

	// 生成混淆的短码
	code, err := s.generateObfuscatedCode(shortLink.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate short code: %w", err)
	}

	// 更新短码
	shortLink.Code = code
	if err := config.DB.Save(shortLink).Error; err != nil {
		return nil, fmt.Errorf("failed to save short code: %w", err)
	}

	return shortLink, nil
}

// 增加访问计数
func (s *ShortLinkService) IncrementViewCount(code string) error {
	var shortLink model.ShortLink
	if err := config.DB.Where("code = ?", code).First(&shortLink).Error; err != nil {
		return fmt.Errorf("failed to find short link by code: %w", err)
	}
	return s._incrementViewCount(&shortLink)
}

// 获取原始URL并增加访问计数
func (s *ShortLinkService) GetOriginalURL(code string) (string, error) {
	var shortLink model.ShortLink
	if err := config.DB.Where("code = ?", code).First(&shortLink).Error; err != nil {
		return "", fmt.Errorf("failed to find short link by code: %w", err)
	}

	// 更新访问计数
	if err := s._incrementViewCount(&shortLink); err != nil {
		return "", fmt.Errorf("failed to increment view count: %w", err)
	}

	return shortLink.OriginalURL, nil
}

// 私有方法：增加访问计数
func (s *ShortLinkService) _incrementViewCount(shortLink *model.ShortLink) error {
	shortLink.ViewCount++
	return config.DB.Save(shortLink).Error
}

// 生成混淆的短码
func (s *ShortLinkService) generateObfuscatedCode(id uint) (string, error) {
	if id == 0 {
		return "", errors.New("invalid ID for short code generation")
	}

	// 1. 将ID转换为固定长度的字符串
	paddedID := s.padID(id)

	// 2. 反转字符串
	reversedID := s.reverseString(paddedID)

	// 3. 转换为Base62编码
	base62Code := s.encodeBase62(reversedID)
	if base62Code == "" {
		return "", errors.New("failed to encode Base62")
	}

	// 4. 填充到固定长度
	base62Code = s.padBase62Code(base62Code)

	// 5. 混淆处理
	obfuscatedCode := s.obfuscateCode(base62Code)

	return obfuscatedCode, nil
}

// 将ID转换为固定长度的字符串
func (s *ShortLinkService) padID(id uint) string {
	// 将ID转换为固定长度的字符串，前面补0
	return fmt.Sprintf("%06d", id)
}

// 反转字符串
func (s *ShortLinkService) reverseString(input string) string {
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// 转换为Base62编码
func (s *ShortLinkService) encodeBase62(input string) string {
	num, err := strconv.Atoi(input)
	if err != nil {
		return ""
	}

	var result strings.Builder
	for num > 0 {
		remainder := num % 62
		result.WriteByte(base62Chars[remainder])
		num = num / 62
	}

	// Base62编码结果需要反转
	return s.reverseString(result.String())
}

// 填充Base62编码到固定长度
func (s *ShortLinkService) padBase62Code(code string) string {
	// 如果长度不足6位，前面补 'a'（或其他字符）
	for len(code) < codeLength {
		code = "a" + code
	}
	return code
}

// 混淆短码
func (s *ShortLinkService) obfuscateCode(code string) string {
	runes := []rune(code)
	seed := time.Now().UnixNano() % 100 // 使用时间戳作为随机种子

	for i := len(runes) - 1; i > 0; i-- {
		// 动态计算乘数
		multiplier := 31 + int(seed)
		j := int(math.Mod(float64(i*multiplier), float64(i+1)))

		// 交换字符
		runes[i], runes[j] = runes[j], runes[i]

		// 确保字符在Base62范围内
		if i%2 == 0 {
			newCharIndex := (strings.IndexRune(base62Chars, runes[i]) + 1) % len(base62Chars)
			runes[i] = rune(base62Chars[newCharIndex])
		}
	}

	return string(runes)
}

// 验证 URL 格式是否合法
func isValidURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}
