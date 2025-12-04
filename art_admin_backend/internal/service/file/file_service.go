package file

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// FileService 文件服务
type FileService struct {
	repo       *repository.FileRepository
	uploadPath string
	baseURL    string
}

// NewFileService 创建文件服务
func NewFileService(repo *repository.FileRepository, uploadPath, baseURL string) *FileService {
	// 确保上传目录存在
	os.MkdirAll(uploadPath, 0755)
	return &FileService{
		repo:       repo,
		uploadPath: uploadPath,
		baseURL:    baseURL,
	}
}

// Upload 上传文件
func (s *FileService) Upload(file *multipart.FileHeader, userID uint, category string) (*model.File, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 获取文件信息
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if category == "" {
		category = s.detectCategory(ext)
	}

	// 生成存储文件名
	storageName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 按日期创建子目录
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(s.uploadPath, category, dateDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 完整文件路径
	filePath := filepath.Join(category, dateDir, storageName)
	fullPath := filepath.Join(s.uploadPath, filePath)

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 创建文件记录
	fileRecord := &model.File{
		FileName:     file.Filename,
		StorageName:  storageName,
		FilePath:     filePath,
		FileSize:     file.Size,
		FileType:     file.Header.Get("Content-Type"),
		FileExt:      ext,
		Category:     category,
		UploadUserID: userID,
		URL:          fmt.Sprintf("%s/uploads/%s", s.baseURL, filePath),
	}

	if err := s.repo.Create(fileRecord); err != nil {
		// 删除已上传的文件
		os.Remove(fullPath)
		return nil, fmt.Errorf("保存文件记录失败: %w", err)
	}

	return fileRecord, nil
}

// GetByID 根据ID获取文件
func (s *FileService) GetByID(id uint) (*model.File, error) {
	return s.repo.GetByID(id)
}

// GetFilePath 获取文件完整路径
func (s *FileService) GetFilePath(id uint) (string, error) {
	file, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}
	return filepath.Join(s.uploadPath, file.FilePath), nil
}

// List 获取文件列表
func (s *FileService) List(page, pageSize int, category string, userID uint) ([]model.File, int64, error) {
	return s.repo.List(page, pageSize, category, userID)
}

// Delete 删除文件
func (s *FileService) Delete(id uint) error {
	file, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除物理文件
	fullPath := filepath.Join(s.uploadPath, file.FilePath)
	os.Remove(fullPath)

	// 删除数据库记录
	return s.repo.Delete(id)
}

// BatchDelete 批量删除
func (s *FileService) BatchDelete(ids []uint) error {
	for _, id := range ids {
		if err := s.Delete(id); err != nil {
			continue // 忽略单个文件删除失败
		}
	}
	return nil
}

// detectCategory 根据扩展名检测文件分类
func (s *FileService) detectCategory(ext string) string {
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg", ".ico"}
	videoExts := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv", ".webm"}
	audioExts := []string{".mp3", ".wav", ".ogg", ".flac", ".aac", ".m4a"}
	docExts := []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".md"}

	for _, e := range imageExts {
		if ext == e {
			return "image"
		}
	}
	for _, e := range videoExts {
		if ext == e {
			return "video"
		}
	}
	for _, e := range audioExts {
		if ext == e {
			return "audio"
		}
	}
	for _, e := range docExts {
		if ext == e {
			return "document"
		}
	}
	return "other"
}
