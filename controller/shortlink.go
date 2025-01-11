package controller

import (
	"net/http"

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
	URL string `json:"url" binding:"required,url"`
}

type CreateShortLinkResponse struct {
	ShortURL string `json:"short_url"`
}

func (c *ShortLinkController) CreateShortLink(ctx *gin.Context) {
	var req CreateShortLinkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	shortLink, err := c.service.CreateShortLink(req.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create short link",
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
			"error": "Short link not found",
		})
		return
	}

	// 返回 302 重定向到原始 URL
	ctx.Redirect(http.StatusFound, originalURL)
}
