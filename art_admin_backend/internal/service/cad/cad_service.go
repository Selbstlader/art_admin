package cad

import (
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// CadService CAD文件服务
// CAD file service for handling CAD file operations
type CadService struct {
	repo          *repository.CadFileRepository
	uploadPath    string
	baseURL       string
	converterType string // 转换器类型: oda, libredwg / Converter type
	odaPath       string // ODA File Converter路径 / ODA path
	libredwgPath  string // LibreDWG dwg2dxf路径 / LibreDWG path
}

// NewCadService 创建CAD文件服务
// Create new CAD file service
func NewCadService(repo *repository.CadFileRepository, uploadPath, baseURL, converterType, odaPath, libredwgPath string) *CadService {
	// 确保上传目录存在 / Ensure upload directory exists
	cadPath := filepath.Join(uploadPath, "cad")
	os.MkdirAll(cadPath, 0755)
	os.MkdirAll(filepath.Join(cadPath, "original"), 0755)
	os.MkdirAll(filepath.Join(cadPath, "parsed"), 0755)

	return &CadService{
		repo:          repo,
		uploadPath:    uploadPath,
		baseURL:       baseURL,
		converterType: converterType,
		odaPath:       odaPath,
		libredwgPath:  libredwgPath,
	}
}

// Upload 上传CAD文件
// Upload CAD file
func (s *CadService) Upload(file *multipart.FileHeader, projectID uint) (*response.CadFileUploadResponse, error) {
	// 验证文件格式 / Validate file format
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".dwg" && ext != ".dxf" {
		return nil, fmt.Errorf("不支持的文件格式，请上传DWG或DXF文件")
	}

	// 打开上传的文件 / Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 生成存储文件名 / Generate storage file name
	storageName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dateDir := time.Now().Format("2006/01/02")

	// 创建目录 / Create directory
	fullDir := filepath.Join(s.uploadPath, "cad", "original", dateDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 完整文件路径 / Full file path
	relativePath := filepath.Join("cad", "original", dateDir, storageName)
	fullPath := filepath.Join(s.uploadPath, relativePath)

	// 创建目标文件 / Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容 / Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 确定文件格式 / Determine file format
	fileFormat := "dxf"
	if ext == ".dwg" {
		fileFormat = "dwg"
	}

	// 创建数据库记录 / Create database record
	cadFile := &model.CadFile{
		ProjectID:    projectID,
		FileName:     file.Filename,
		OriginalPath: relativePath,
		FileFormat:   fileFormat,
		ParseStatus:  "pending",
		FileSize:     file.Size,
		Layers:       "[]", // 初始化为空JSON数组 / Initialize as empty JSON array
	}

	if err := s.repo.Create(cadFile); err != nil {
		os.Remove(fullPath)
		return nil, fmt.Errorf("保存文件记录失败: %w", err)
	}

	return &response.CadFileUploadResponse{
		ID:         cadFile.ID,
		FileName:   cadFile.FileName,
		FileFormat: cadFile.FileFormat,
		FileSize:   cadFile.FileSize,
		ProjectID:  cadFile.ProjectID,
		Status:     cadFile.ParseStatus,
	}, nil
}

// Parse 解析CAD文件
// Parse CAD file
func (s *CadService) Parse(cadFileID uint) (*response.CadParseResult, error) {
	// 获取CAD文件记录 / Get CAD file record
	cadFile, err := s.repo.GetByID(cadFileID)
	if err != nil {
		return nil, fmt.Errorf("CAD文件不存在: %w", err)
	}

	// 更新状态为处理中 / Update status to processing
	if err := s.repo.UpdateParseStatus(cadFileID, "processing", ""); err != nil {
		return nil, fmt.Errorf("更新状态失败: %w", err)
	}

	// 获取原始文件路径 / Get original file path
	originalPath := filepath.Join(s.uploadPath, cadFile.OriginalPath)

	// 如果是DWG文件，需要先转换为DXF / Convert DWG to DXF if needed
	dxfPath := originalPath
	if cadFile.FileFormat == "dwg" {
		convertedPath, err := s.convertDwgToDxf(originalPath)
		if err != nil {
			s.repo.UpdateParseStatus(cadFileID, "failed", fmt.Sprintf("DWG转换失败: %v", err))
			return nil, fmt.Errorf("DWG转换失败: %w", err)
		}
		dxfPath = convertedPath
	}

	// 解析DXF文件 / Parse DXF file
	parseResult, err := s.parseDxfFile(dxfPath)
	if err != nil {
		s.repo.UpdateParseStatus(cadFileID, "failed", fmt.Sprintf("DXF解析失败: %v", err))
		return nil, fmt.Errorf("DXF解析失败: %w", err)
	}

	// 保存解析结果到JSON文件 / Save parse result to JSON file
	parsedPath, err := s.saveParsedData(cadFileID, parseResult)
	if err != nil {
		s.repo.UpdateParseStatus(cadFileID, "failed", fmt.Sprintf("保存解析结果失败: %v", err))
		return nil, fmt.Errorf("保存解析结果失败: %w", err)
	}

	// 序列化图层列表 / Serialize layer list
	layerNames := make([]string, len(parseResult.Layers))
	for i, layer := range parseResult.Layers {
		layerNames[i] = layer.Name
	}
	layersJSON, _ := json.Marshal(layerNames)

	// 更新数据库记录 / Update database record
	if err := s.repo.UpdateParseResult(cadFileID, parsedPath, len(parseResult.Layers), string(layersJSON), parseResult.Has3D); err != nil {
		return nil, fmt.Errorf("更新解析结果失败: %w", err)
	}

	parseResult.CadFileID = cadFileID
	parseResult.FileName = cadFile.FileName
	parseResult.FileFormat = cadFile.FileFormat
	parseResult.ParsedData = parsedPath

	return parseResult, nil
}

// convertDwgToDxf 将DWG文件转换为DXF
// Convert DWG file to DXF using configured converter
func (s *CadService) convertDwgToDxf(dwgPath string) (string, error) {
	// 根据配置选择转换器 / Select converter based on config
	switch s.converterType {
	case "oda":
		return s.convertWithODA(dwgPath)
	case "libredwg":
		return s.convertWithLibreDWG(dwgPath)
	default:
		// 尝试自动检测可用的转换器 / Try to auto-detect available converter
		if s.odaPath != "" {
			if _, err := os.Stat(s.odaPath); err == nil {
				return s.convertWithODA(dwgPath)
			}
		}
		if s.libredwgPath != "" {
			if _, err := os.Stat(s.libredwgPath); err == nil {
				return s.convertWithLibreDWG(dwgPath)
			}
		}
		return "", fmt.Errorf("DWG格式需要配置转换器才能解析。请在config.yaml中配置dwgconverter，支持ODA File Converter或LibreDWG，或直接上传DXF格式文件")
	}
}

// convertWithODA 使用ODA File Converter转换
// Convert using ODA File Converter
func (s *CadService) convertWithODA(dwgPath string) (string, error) {
	if s.odaPath == "" {
		return "", fmt.Errorf("ODA File Converter路径未配置")
	}
	if _, err := os.Stat(s.odaPath); os.IsNotExist(err) {
		return "", fmt.Errorf("ODA File Converter不存在: %s", s.odaPath)
	}

	inputDir := filepath.Dir(dwgPath)
	outputDir := filepath.Join(s.uploadPath, "cad", "converted", time.Now().Format("2006/01/02"))
	os.MkdirAll(outputDir, 0755)

	// ODA命令格式: ODAFileConverter <input_folder> <output_folder> <output_version> <output_type> <recurse> <audit>
	cmd := exec.Command(s.odaPath, inputDir, outputDir, "ACAD2018", "DXF", "0", "1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ODA转换失败: %v, 输出: %s", err, string(output))
	}

	baseName := strings.TrimSuffix(filepath.Base(dwgPath), filepath.Ext(dwgPath))
	dxfPath := filepath.Join(outputDir, baseName+".dxf")

	if _, err := os.Stat(dxfPath); os.IsNotExist(err) {
		return "", fmt.Errorf("转换后的DXF文件不存在: %s", dxfPath)
	}
	return dxfPath, nil
}

// convertWithLibreDWG 使用LibreDWG转换
// Convert using LibreDWG dwg2dxf
func (s *CadService) convertWithLibreDWG(dwgPath string) (string, error) {
	if s.libredwgPath == "" {
		return "", fmt.Errorf("LibreDWG路径未配置")
	}
	if _, err := os.Stat(s.libredwgPath); os.IsNotExist(err) {
		return "", fmt.Errorf("LibreDWG dwg2dxf不存在: %s", s.libredwgPath)
	}

	outputDir := filepath.Join(s.uploadPath, "cad", "converted", time.Now().Format("2006/01/02"))
	os.MkdirAll(outputDir, 0755)

	baseName := strings.TrimSuffix(filepath.Base(dwgPath), filepath.Ext(dwgPath))
	dxfPath := filepath.Join(outputDir, baseName+".dxf")

	// LibreDWG dwg2dxf命令格式: dwg2dxf -o output.dxf input.dwg
	cmd := exec.Command(s.libredwgPath, "-o", dxfPath, dwgPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("LibreDWG转换失败: %v, 输出: %s", err, string(output))
	}

	if _, err := os.Stat(dxfPath); os.IsNotExist(err) {
		return "", fmt.Errorf("转换后的DXF文件不存在: %s", dxfPath)
	}
	return dxfPath, nil
}

// parseDxfFile 解析DXF文件
// Parse DXF file and extract layer and entity information
func (s *CadService) parseDxfFile(dxfPath string) (*response.CadParseResult, error) {
	// 读取DXF文件内容 / Read DXF file content
	content, err := os.ReadFile(dxfPath)
	if err != nil {
		return nil, fmt.Errorf("读取DXF文件失败: %w", err)
	}

	// 解析DXF内容 / Parse DXF content
	parser := NewDxfParser(string(content))
	return parser.Parse()
}

// saveParsedData 保存解析数据到JSON文件
// Save parsed data to JSON file
func (s *CadService) saveParsedData(cadFileID uint, result *response.CadParseResult) (string, error) {
	// 创建解析数据目录 / Create parsed data directory
	dateDir := time.Now().Format("2006/01/02")
	parsedDir := filepath.Join(s.uploadPath, "cad", "parsed", dateDir)
	os.MkdirAll(parsedDir, 0755)

	// 生成JSON文件名 / Generate JSON file name
	jsonFileName := fmt.Sprintf("%d_%s.json", cadFileID, uuid.New().String())
	relativePath := filepath.Join("cad", "parsed", dateDir, jsonFileName)
	fullPath := filepath.Join(s.uploadPath, relativePath)

	// 序列化并保存 / Serialize and save
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化解析结果失败: %w", err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("写入解析结果文件失败: %w", err)
	}

	return relativePath, nil
}

// GetByID 根据ID获取CAD文件
// Get CAD file by ID
func (s *CadService) GetByID(id uint) (*response.CadFileResponse, error) {
	cadFile, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.toCadFileResponse(cadFile), nil
}

// GetParseResult 获取解析结果
// Get parse result
func (s *CadService) GetParseResult(id uint) (*response.CadParseResult, error) {
	cadFile, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("CAD文件不存在: %w", err)
	}

	if cadFile.ParseStatus != "completed" {
		return nil, fmt.Errorf("CAD文件尚未解析完成，当前状态: %s", cadFile.ParseStatus)
	}

	if cadFile.ParsedPath == "" {
		return nil, fmt.Errorf("解析数据不存在")
	}

	// 读取解析结果文件 / Read parse result file
	fullPath := filepath.Join(s.uploadPath, cadFile.ParsedPath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("读取解析结果失败: %w", err)
	}

	var result response.CadParseResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("解析结果格式错误: %w", err)
	}

	return &result, nil
}

// GetLayers 获取图层列表
// Get layer list
func (s *CadService) GetLayers(id uint) (*response.CadLayerListResponse, error) {
	cadFile, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("CAD文件不存在: %w", err)
	}

	// 解析图层JSON / Parse layers JSON
	var layerNames []string
	if cadFile.Layers != "" {
		if err := json.Unmarshal([]byte(cadFile.Layers), &layerNames); err != nil {
			return nil, fmt.Errorf("图层数据格式错误: %w", err)
		}
	}

	// 如果有解析结果文件，从中获取详细图层信息 / Get detailed layer info from parse result if available
	var layers []response.CadLayerInfo
	if cadFile.ParsedPath != "" && cadFile.ParseStatus == "completed" {
		fullPath := filepath.Join(s.uploadPath, cadFile.ParsedPath)
		data, err := os.ReadFile(fullPath)
		if err == nil {
			var result response.CadParseResult
			if json.Unmarshal(data, &result) == nil {
				layers = result.Layers
			}
		}
	}

	// 如果没有详细信息，使用简单的图层名称列表 / Use simple layer names if no detailed info
	if layers == nil {
		layers = make([]response.CadLayerInfo, len(layerNames))
		for i, name := range layerNames {
			layers[i] = response.CadLayerInfo{
				Name:    name,
				Visible: true,
			}
		}
	}

	return &response.CadLayerListResponse{
		CadFileID:  cadFile.ID,
		FileName:   cadFile.FileName,
		LayerCount: cadFile.LayerCount,
		Layers:     layers,
	}, nil
}

// List 获取CAD文件列表
// List CAD files
func (s *CadService) List(projectID uint, page, pageSize int, parseStatus, fileFormat string) (*response.CadFileListResponse, error) {
	cadFiles, total, err := s.repo.List(projectID, page, pageSize, parseStatus, fileFormat)
	if err != nil {
		return nil, err
	}

	records := make([]response.CadFileResponse, len(cadFiles))
	for i, cf := range cadFiles {
		records[i] = *s.toCadFileResponse(&cf)
	}

	return &response.CadFileListResponse{
		Records: records,
		Current: page,
		Size:    pageSize,
		Total:   total,
	}, nil
}

// Delete 删除CAD文件
// Delete CAD file
func (s *CadService) Delete(id uint) error {
	cadFile, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除原始文件 / Delete original file
	if cadFile.OriginalPath != "" {
		os.Remove(filepath.Join(s.uploadPath, cadFile.OriginalPath))
	}

	// 删除解析结果文件 / Delete parsed result file
	if cadFile.ParsedPath != "" {
		os.Remove(filepath.Join(s.uploadPath, cadFile.ParsedPath))
	}

	return s.repo.Delete(id)
}

// BatchDelete 批量删除CAD文件
// Batch delete CAD files
func (s *CadService) BatchDelete(ids []uint) error {
	for _, id := range ids {
		s.Delete(id)
	}
	return nil
}

// toCadFileResponse 转换为响应结构
// Convert to response structure
func (s *CadService) toCadFileResponse(cadFile *model.CadFile) *response.CadFileResponse {
	var layers []string
	if cadFile.Layers != "" {
		json.Unmarshal([]byte(cadFile.Layers), &layers)
	}

	return &response.CadFileResponse{
		ID:           cadFile.ID,
		ProjectID:    cadFile.ProjectID,
		FileName:     cadFile.FileName,
		OriginalPath: cadFile.OriginalPath,
		ParsedPath:   cadFile.ParsedPath,
		FileFormat:   cadFile.FileFormat,
		ParseStatus:  cadFile.ParseStatus,
		LayerCount:   cadFile.LayerCount,
		Layers:       layers,
		Has3D:        cadFile.Has3D,
		FileSize:     cadFile.FileSize,
		ErrorMessage: cadFile.ErrorMessage,
		CreatedAt:    cadFile.CreatedAt,
		UpdatedAt:    cadFile.UpdatedAt,
	}
}
