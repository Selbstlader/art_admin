package api

import (
	"net/http"
	"strconv"

	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/response"
	"art_admin_backend/internal/service/xhs"

	"github.com/gin-gonic/gin"
)

// XHSAPI 小红书总结API
type XHSAPI struct {
	service *xhs.SummaryService
}

// NewXHSAPI 创建小红书API
func NewXHSAPI() *XHSAPI {
	return &XHSAPI{
		service: xhs.NewSummaryService(),
	}
}

// ========== 请求/响应结构 ==========

// SummarizeRequest 总结请求
type SummarizeRequest struct {
	URL            string   `json:"url"`             // 小红书链接（可选，兼容旧字段）
	OriginalURL    string   `json:"original_url"`    // 小红书链接（前端新字段）
	NoteTitle      string   `json:"note_title"`      // 笔记标题
	Content        string   `json:"content"`         // 笔记文字内容
	ImageURLs      []string `json:"image_urls"`      // 图片URL列表
	Style          string   `json:"style"`           // 总结风格: concise, detailed, casual
	MaxLength      int      `json:"max_length"`      // 最大字数
	EnableOCR      bool     `json:"enable_ocr"`      // 是否开启图片OCR
	EnableAnalysis bool     `json:"enable_analysis"` // 是否开启专业分析
	SaveHistory    bool     `json:"save_history"`    // 是否保存到历史记录
}

// OCRRequest OCR识别请求
type OCRRequest struct {
	ImageURL  string `json:"image_url"`  // 图片URL
	ImageData string `json:"image_data"` // 图片Base64数据（二选一）
}

// FavoriteRequest 收藏请求
type FavoriteRequest struct {
	IsFavorite bool  `json:"is_favorite"`
	FolderID   int64 `json:"folder_id"`
}

// FolderRequest 收藏夹请求
type FolderRequest struct {
	Name      string `json:"name" binding:"required"`
	SortOrder int    `json:"sort_order"`
}

// ========== API接口 ==========

// Summarize 生成总结
// @Summary 生成小红书笔记总结
// @Tags 小红书总结
// @Accept json
// @Produce json
// @Param request body SummarizeRequest true "总结请求"
// @Success 200 {object} response.Response
// @Router /api/xhs/summarize [post]
func (a *XHSAPI) Summarize(c *gin.Context) {
	var req SummarizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 合并URL字段（兼容两种传参方式）
	noteURL := req.OriginalURL
	if noteURL == "" {
		noteURL = req.URL
	}

	// 验证：至少需要提供链接、内容或图片之一
	if noteURL == "" && req.Content == "" && len(req.ImageURLs) == 0 {
		response.Error(c, http.StatusBadRequest, "请提供小红书链接、笔记内容或图片")
		return
	}

	// 获取用户ID
	userID := getUserID(c)

	// 构建服务请求
	svcReq := &xhs.SummaryRequest{
		URL:            noteURL,
		Content:        req.Content,
		ImageURLs:      req.ImageURLs,
		Style:          req.Style,
		MaxLength:      req.MaxLength,
		EnableOCR:      req.EnableOCR,
		EnableAnalysis: req.EnableAnalysis,
	}

	// 根据是否保存历史决定调用方法
	if req.SaveHistory && userID > 0 {
		summary, err := a.service.SummarizeAndSave(c.Request.Context(), userID, svcReq, req.NoteTitle)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "生成总结失败: "+err.Error())
			return
		}
		response.Success(c, summary)
	} else {
		result, err := a.service.Summarize(c.Request.Context(), svcReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "生成总结失败: "+err.Error())
			return
		}
		response.Success(c, result)
	}
}

// OCR 图片文字识别
// @Summary OCR图片文字识别
// @Tags 小红书总结
// @Accept json
// @Produce json
// @Param request body OCRRequest true "OCR请求"
// @Success 200 {object} response.Response
// @Router /api/xhs/ocr [post]
func (a *XHSAPI) OCR(c *gin.Context) {
	var req OCRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if req.ImageURL == "" && req.ImageData == "" {
		response.Error(c, http.StatusBadRequest, "请提供图片URL或Base64数据")
		return
	}

	var text string
	var err error

	if req.ImageURL != "" {
		text, err = a.service.ExtractTextFromImage(c.Request.Context(), req.ImageURL)
	} else {
		// Base64解码并识别
		imageBase64 := req.ImageData
		// 移除可能的data:image前缀 (如 data:image/png;base64,xxxxx)
		if len(imageBase64) > 11 && imageBase64[:11] == "data:image/" {
			for i := 0; i < len(imageBase64); i++ {
				if imageBase64[i] == ',' {
					imageBase64 = imageBase64[i+1:]
					break
				}
			}
		}
		text, err = a.service.ExtractTextFromImageBase64(c.Request.Context(), imageBase64)
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "OCR识别失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"text": text})
}

// GetHistory 获取总结历史
// @Summary 获取总结历史列表
// @Tags 小红书总结
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/xhs/history [get]
func (a *XHSAPI) GetHistory(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	summaries, total, err := a.service.GetUserSummaries(userID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取历史记录失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      summaries,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetSummaryDetail 获取总结详情
// @Summary 获取总结详情
// @Tags 小红书总结
// @Produce json
// @Param id path int true "总结ID"
// @Success 200 {object} response.Response
// @Router /api/xhs/summary/{id} [get]
func (a *XHSAPI) GetSummaryDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	summary, err := a.service.GetSummaryByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "总结不存在")
		return
	}

	response.Success(c, summary)
}

// DeleteSummary 删除总结
// @Summary 删除总结
// @Tags 小红书总结
// @Param id path int true "总结ID"
// @Success 200 {object} response.Response
// @Router /api/xhs/summary/{id} [delete]
func (a *XHSAPI) DeleteSummary(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := a.service.DeleteSummary(id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateFavorite 更新收藏状态
// @Summary 更新收藏状态
// @Tags 小红书总结
// @Accept json
// @Produce json
// @Param id path int true "总结ID"
// @Param request body FavoriteRequest true "收藏请求"
// @Success 200 {object} response.Response
// @Router /api/xhs/summary/{id}/favorite [put]
func (a *XHSAPI) UpdateFavorite(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	var req FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	if err := a.service.UpdateFavorite(id, userID, req.IsFavorite, req.FolderID); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新收藏状态失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetFavorites 获取收藏列表
// @Summary 获取收藏列表
// @Tags 小红书总结
// @Produce json
// @Param folder_id query int false "收藏夹ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/xhs/favorites [get]
func (a *XHSAPI) GetFavorites(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	folderID, _ := strconv.ParseInt(c.Query("folder_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	summaries, total, err := a.service.GetFavorites(userID, folderID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取收藏列表失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      summaries,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// SearchSummaries 搜索总结
// @Summary 搜索总结
// @Tags 小红书总结
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/xhs/search [get]
func (a *XHSAPI) SearchSummaries(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	summaries, total, err := a.service.SearchSummaries(userID, keyword, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "搜索失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      summaries,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== 收藏夹接口 ==========

// CreateFolder 创建收藏夹
// @Summary 创建收藏夹
// @Tags 小红书总结
// @Accept json
// @Produce json
// @Param request body FolderRequest true "收藏夹请求"
// @Success 200 {object} response.Response
// @Router /api/xhs/folder [post]
func (a *XHSAPI) CreateFolder(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	var req FolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	folder, err := a.service.CreateFolder(userID, req.Name)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建收藏夹失败: "+err.Error())
		return
	}

	response.Success(c, folder)
}

// GetFolders 获取收藏夹列表
// @Summary 获取收藏夹列表
// @Tags 小红书总结
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/xhs/folders [get]
func (a *XHSAPI) GetFolders(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	folders, err := a.service.GetUserFolders(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取收藏夹失败: "+err.Error())
		return
	}

	response.Success(c, folders)
}

// UpdateFolder 更新收藏夹
// @Summary 更新收藏夹
// @Tags 小红书总结
// @Accept json
// @Produce json
// @Param id path int true "收藏夹ID"
// @Param request body FolderRequest true "收藏夹请求"
// @Success 200 {object} response.Response
// @Router /api/xhs/folder/{id} [put]
func (a *XHSAPI) UpdateFolder(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	var req FolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	folder := &model.XHSFolder{
		ID:        id,
		UserID:    userID,
		Name:      req.Name,
		SortOrder: req.SortOrder,
	}

	if err := a.service.UpdateFolder(folder); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新收藏夹失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteFolder 删除收藏夹
// @Summary 删除收藏夹
// @Tags 小红书总结
// @Param id path int true "收藏夹ID"
// @Success 200 {object} response.Response
// @Router /api/xhs/folder/{id} [delete]
func (a *XHSAPI) DeleteFolder(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := a.service.DeleteFolder(id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除收藏夹失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// ========== 辅助函数 ==========

// getUserID 从上下文获取用户ID
func getUserID(c *gin.Context) int64 {
	// 优先检查 userId (JWT中间件设置的key)
	if id, exists := c.Get("userId"); exists {
		if userID, ok := id.(int64); ok {
			return userID
		}
		if userID, ok := id.(float64); ok {
			return int64(userID)
		}
		if userID, ok := id.(int); ok {
			return int64(userID)
		}
	}
	// 兼容 user_id
	if id, exists := c.Get("user_id"); exists {
		if userID, ok := id.(int64); ok {
			return userID
		}
		if userID, ok := id.(float64); ok {
			return int64(userID)
		}
	}
	return 0
}
