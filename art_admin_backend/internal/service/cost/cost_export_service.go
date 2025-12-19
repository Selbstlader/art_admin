package cost

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/repository"

	"github.com/xuri/excelize/v2"
)

// CostExportService 成本报告导出服务接口
// Cost report export service interface
type CostExportService interface {
	ExportCostReport(req *request.ExportCostReportRequest) (*response.CostReportExportResponse, error)
	ExportToExcel(projectID uint) ([]byte, string, error)
	ExportToPDF(projectID uint) ([]byte, string, error)
}

// costExportService 成本报告导出服务实现
// Cost report export service implementation
type costExportService struct {
	costRepo    repository.CostEstimateRepository
	projectRepo *repository.DesignerProjectRepository
}

// NewCostExportService 创建成本报告导出服务实例
// Create cost report export service instance
func NewCostExportService() CostExportService {
	return &costExportService{
		costRepo:    repository.NewCostEstimateRepository(),
		projectRepo: repository.NewDesignerProjectRepository(database.DB),
	}
}

// ExportCostReport 导出成本报告
// Export cost report
func (s *costExportService) ExportCostReport(req *request.ExportCostReportRequest) (*response.CostReportExportResponse, error) {
	var data []byte
	var fileName string
	var err error

	switch req.Format {
	case "excel", "xlsx":
		data, fileName, err = s.ExportToExcel(req.ProjectID)
	case "pdf":
		data, fileName, err = s.ExportToPDF(req.ProjectID)
	default:
		return nil, fmt.Errorf("不支持的导出格式: %s / Unsupported export format: %s", req.Format, req.Format)
	}

	if err != nil {
		return nil, err
	}

	// 保存文件到临时目录 / Save file to temp directory
	tempDir := "temp/exports"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("创建导出目录失败 / Failed to create export directory: %v", err)
	}

	filePath := filepath.Join(tempDir, fileName)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("保存导出文件失败 / Failed to save export file: %v", err)
	}

	return &response.CostReportExportResponse{
		FileURL:  "/api/file/download/" + fileName,
		FileName: fileName,
		FileSize: int64(len(data)),
	}, nil
}

// ExportToExcel 导出为Excel格式
// Export to Excel format
func (s *costExportService) ExportToExcel(projectID uint) ([]byte, string, error) {
	// 获取成本估算数据 / Get cost estimate data
	estimate, err := s.costRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, "", fmt.Errorf("获取成本估算失败 / Failed to get cost estimate: %v", err)
	}

	// 获取项目信息 / Get project info
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, "", fmt.Errorf("获取项目信息失败 / Failed to get project info: %v", err)
	}

	// 创建Excel文件 / Create Excel file
	f := excelize.NewFile()
	defer f.Close()

	// 设置工作表名称 / Set sheet name
	sheetName := "成本报告"
	f.SetSheetName("Sheet1", sheetName)

	// 设置列宽 / Set column width
	f.SetColWidth(sheetName, "A", "A", 20)
	f.SetColWidth(sheetName, "B", "B", 25)
	f.SetColWidth(sheetName, "C", "C", 15)
	f.SetColWidth(sheetName, "D", "D", 15)
	f.SetColWidth(sheetName, "E", "E", 15)

	// 创建标题样式 / Create title style
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 16,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	// 创建表头样式 / Create header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E0E0E0"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	// 创建数据样式 / Create data style
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
	})

	// 创建金额样式 / Create amount style
	amountStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		NumFmt: 4, // #,##0.00
	})

	// 写入标题 / Write title
	f.MergeCell(sheetName, "A1", "E1")
	f.SetCellValue(sheetName, "A1", "项目成本预算报告")
	f.SetCellStyle(sheetName, "A1", "E1", titleStyle)
	f.SetRowHeight(sheetName, 1, 30)

	// 写入项目信息 / Write project info
	f.SetCellValue(sheetName, "A3", "项目名称:")
	f.SetCellValue(sheetName, "B3", project.Name)
	f.SetCellValue(sheetName, "A4", "项目面积:")
	f.SetCellValue(sheetName, "B4", fmt.Sprintf("%.2f 平方米", project.Area))
	f.SetCellValue(sheetName, "A5", "设计风格:")
	f.SetCellValue(sheetName, "B5", project.Style)
	f.SetCellValue(sheetName, "A6", "报告日期:")
	f.SetCellValue(sheetName, "B6", time.Now().Format("2006-01-02"))

	// 写入成本汇总表头 / Write cost summary header
	row := 8
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "成本类别")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "金额(元)")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "占比")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), headerStyle)

	// 计算各项占比 / Calculate percentages
	totalCost := estimate.TotalCost
	materialPercent := 0.0
	laborPercent := 0.0
	equipPercent := 0.0
	mgmtPercent := 0.0
	if totalCost > 0 {
		materialPercent = estimate.MaterialCost / totalCost * 100
		laborPercent = estimate.LaborCost / totalCost * 100
		equipPercent = estimate.EquipmentCost / totalCost * 100
		mgmtPercent = estimate.ManagementCost / totalCost * 100
	}

	// 写入成本数据 / Write cost data
	costData := []struct {
		Category string
		Amount   float64
		Percent  float64
	}{
		{"材料费", estimate.MaterialCost, materialPercent},
		{"人工费", estimate.LaborCost, laborPercent},
		{"设备费", estimate.EquipmentCost, equipPercent},
		{"管理费", estimate.ManagementCost, mgmtPercent},
		{"总计", estimate.TotalCost, 100},
	}

	for i, item := range costData {
		r := row + 1 + i
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", r), item.Category)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", r), item.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", r), fmt.Sprintf("%.2f%%", item.Percent))
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", r), fmt.Sprintf("A%d", r), dataStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", r), fmt.Sprintf("B%d", r), amountStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("C%d", r), fmt.Sprintf("C%d", r), dataStyle)
	}

	// 写入预算信息 / Write budget info
	budgetRow := row + 7
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", budgetRow), "预算上限:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", budgetRow), estimate.BudgetLimit)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", budgetRow), fmt.Sprintf("B%d", budgetRow), amountStyle)

	// 检查是否超预算 / Check if over budget
	if estimate.BudgetLimit > 0 && estimate.TotalCost > estimate.BudgetLimit {
		exceeded := estimate.TotalCost - estimate.BudgetLimit
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", budgetRow+1), "超出预算:")
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", budgetRow+1), exceeded)
		// 设置红色警告样式 / Set red warning style
		warningStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Color: "#FF0000",
				Bold:  true,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "right",
			},
			NumFmt: 4,
		})
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", budgetRow+1), fmt.Sprintf("B%d", budgetRow+1), warningStyle)
	}

	// 生成文件名 / Generate file name
	fileName := fmt.Sprintf("成本报告_%s_%s.xlsx", project.Name, time.Now().Format("20060102150405"))

	// 写入缓冲区 / Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", fmt.Errorf("生成Excel文件失败 / Failed to generate Excel file: %v", err)
	}

	return buf.Bytes(), fileName, nil
}

// ExportToPDF 导出为PDF格式
// Export to PDF format (simplified implementation using HTML to PDF conversion)
func (s *costExportService) ExportToPDF(projectID uint) ([]byte, string, error) {
	// 获取成本估算数据 / Get cost estimate data
	estimate, err := s.costRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, "", fmt.Errorf("获取成本估算失败 / Failed to get cost estimate: %v", err)
	}

	// 获取项目信息 / Get project info
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, "", fmt.Errorf("获取项目信息失败 / Failed to get project info: %v", err)
	}

	// 计算各项占比 / Calculate percentages
	totalCost := estimate.TotalCost
	materialPercent := 0.0
	laborPercent := 0.0
	equipPercent := 0.0
	mgmtPercent := 0.0
	if totalCost > 0 {
		materialPercent = estimate.MaterialCost / totalCost * 100
		laborPercent = estimate.LaborCost / totalCost * 100
		equipPercent = estimate.EquipmentCost / totalCost * 100
		mgmtPercent = estimate.ManagementCost / totalCost * 100
	}

	// 生成HTML内容 / Generate HTML content
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>项目成本预算报告</title>
    <style>
        body { font-family: "Microsoft YaHei", Arial, sans-serif; padding: 20px; }
        h1 { text-align: center; color: #333; }
        .info { margin: 20px 0; }
        .info p { margin: 5px 0; }
        table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 10px; text-align: left; }
        th { background-color: #f5f5f5; }
        .amount { text-align: right; }
        .warning { color: red; font-weight: bold; }
        .total { font-weight: bold; background-color: #f0f0f0; }
    </style>
</head>
<body>
    <h1>项目成本预算报告</h1>
    <div class="info">
        <p><strong>项目名称:</strong> %s</p>
        <p><strong>项目面积:</strong> %.2f 平方米</p>
        <p><strong>设计风格:</strong> %s</p>
        <p><strong>报告日期:</strong> %s</p>
    </div>
    <table>
        <tr>
            <th>成本类别</th>
            <th>金额(元)</th>
            <th>占比</th>
        </tr>
        <tr>
            <td>材料费</td>
            <td class="amount">%.2f</td>
            <td>%.2f%%</td>
        </tr>
        <tr>
            <td>人工费</td>
            <td class="amount">%.2f</td>
            <td>%.2f%%</td>
        </tr>
        <tr>
            <td>设备费</td>
            <td class="amount">%.2f</td>
            <td>%.2f%%</td>
        </tr>
        <tr>
            <td>管理费</td>
            <td class="amount">%.2f</td>
            <td>%.2f%%</td>
        </tr>
        <tr class="total">
            <td>总计</td>
            <td class="amount">%.2f</td>
            <td>100%%</td>
        </tr>
    </table>
    <div class="info">
        <p><strong>预算上限:</strong> %.2f 元</p>
        %s
    </div>
</body>
</html>`,
		project.Name,
		project.Area,
		project.Style,
		time.Now().Format("2006-01-02"),
		estimate.MaterialCost, materialPercent,
		estimate.LaborCost, laborPercent,
		estimate.EquipmentCost, equipPercent,
		estimate.ManagementCost, mgmtPercent,
		estimate.TotalCost,
		estimate.BudgetLimit,
		func() string {
			if estimate.BudgetLimit > 0 && estimate.TotalCost > estimate.BudgetLimit {
				exceeded := estimate.TotalCost - estimate.BudgetLimit
				return fmt.Sprintf(`<p class="warning">⚠️ 超出预算: %.2f 元</p>`, exceeded)
			}
			return ""
		}(),
	)

	// 生成文件名 / Generate file name
	fileName := fmt.Sprintf("成本报告_%s_%s.html", project.Name, time.Now().Format("20060102150405"))

	// 注意：完整的PDF生成需要额外的库（如wkhtmltopdf或chromedp）
	// 这里暂时返回HTML格式，前端可以使用浏览器打印功能生成PDF
	// Note: Full PDF generation requires additional libraries (like wkhtmltopdf or chromedp)
	// Here we return HTML format, frontend can use browser print function to generate PDF

	return []byte(html), fileName, nil
}
