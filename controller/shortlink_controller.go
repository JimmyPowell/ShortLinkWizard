package controller

import (
	"net/http"
	"net/url"
	"shortlink/service"

	"github.com/gin-gonic/gin"
)

type ShortLinkController struct {
	service *service.ShortLinkService
}

func NewShortLinkController() *ShortLinkController {
	return &ShortLinkController{
		service: service.NewShortLinkService(),
	}
}

type CreateShortLinkRequest struct {
	URL string `json:"url" binding:"required"` // 移除无效的 url 验证
}

type CreateShortLinkResponse struct {
	ShortURL string `json:"short_url"`
}

func (c *ShortLinkController) CreateShortLink(ctx *gin.Context) {
	var req CreateShortLinkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "请求体缺失或格式错误: " + err.Error(), // 输出具体错误
		})
		return
	}

	// 手动验证 URL 格式
	if !isValidURL(req.URL) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "URL 格式无效",
		})
		return
	}

	shortLink, err := c.service.CreateShortLink(req.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建短链接失败",
		})
		return
	}

	// 拼接完整的短链接
	fullShortURL := "http://localhost:8081/" + shortLink.Code

	ctx.JSON(http.StatusCreated, CreateShortLinkResponse{
		ShortURL: fullShortURL,
	})
}

func (c *ShortLinkController) Redirect(ctx *gin.Context) {
	code := ctx.Param("code") // 获取短码

	// 根据短码查找原始 URL
	originalURL, err := c.service.GetOriginalURL(code)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "短链接未找到",
		})
		return
	}

	// 返回 302 重定向到原始 URL
	ctx.Redirect(http.StatusFound, originalURL)
}

// 自定义 URL 验证函数
func isValidURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}
