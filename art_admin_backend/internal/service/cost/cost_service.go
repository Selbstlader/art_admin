package cost

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"
	"art_admin_backend/internal/repository"
)

// CostService 成本估算服务接口
// Cost estimate service interface
type CostService interface {
	// CRUD operations
	CreateCostEstimate(req *request.CreateCostEstimateRequest) (*response.CostEstimateResponse, error)
	UpdateCostEstimate(req *request.UpdateCostEstimateRequest) (*response.CostEstimateResponse, error)
	DeleteCostEstimate(id uint) error
	GetCostEstimateByID(id uint) (*response.CostEstimateResponse, error)
	GetCostEstimateByProjectID(projectID uint) (*response.CostEstimateResponse, error)

	// Calculate operations
	CalculateCost(req *request.CalculateCostRequest) (*response.CostCalculateResponse, error)
	CalculateTotalCost(materialCost, laborCost, equipmentCost, managementCost float64) float64
	CheckBudgetWarning(totalCost, budgetLimit float64) (bool, float64, string)

	// Summary operations
	GetCostSummary(userID uint) (*response.CostSummaryResponse, error)
}

// costService 成本估算服务实现
// Cost estimate service implementation
type costService struct {
	costRepo     repository.CostEstimateRepository
	projectRepo  *repository.DesignerProjectRepository
	materialRepo repository.MaterialRepository
}

// NewCostService 创建成本估算服务实例
// Create cost estimate service instance
func NewCostService() CostService {
	return &costService{
		costRepo:     repository.NewCostEstimateRepository(),
		projectRepo:  repository.NewDesignerProjectRepository(database.DB),
		materialRepo: repository.NewMaterialRepository(),
	}
}

// CreateCostEstimate 创建成本估算
// Create cost estimate
func (s *costService) CreateCostEstimate(req *request.CreateCostEstimateRequest) (*response.CostEstimateResponse, error) {
	// 检查项目是否存在 / Check if project exists
	project, err := s.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, errors.New("项目不存在 / Project not found")
	}

	// 检查是否已存在成本估算 / Check if cost estimate already exists
	existing, _ := s.costRepo.FindByProjectID(req.ProjectID)
	if existing != nil {
		return nil, errors.New("该项目已存在成本估算，请使用更新接口 / Cost estimate already exists for this project, please use update API")
	}

	// 计算各项成本 / Calculate costs
	materialCost := s.calculateMaterialCost(req.Items)
	laborCost := s.calculateLaborCost(project.Area, req.LaborRate)
	equipmentCost := s.calculateEquipmentCost(project.Area, req.EquipRate)

	// 计算基础成本(不含管理费) / Calculate base cost (without management)
	baseCost := materialCost + laborCost + equipmentCost
	managementCost := s.calculateManagementCost(baseCost, req.MgmtRate)

	// 计算总成本 / Calculate total cost
	totalCost := s.CalculateTotalCost(materialCost, laborCost, equipmentCost, managementCost)

	// 序列化明细项 / Serialize items
	itemsJSON, err := json.Marshal(req.Items)
	if err != nil {
		return nil, fmt.Errorf("序列化明细项失败 / Failed to serialize items: %v", err)
	}

	estimate := &model.CostEstimate{
		ProjectID:      req.ProjectID,
		MaterialCost:   materialCost,
		LaborCost:      laborCost,
		EquipmentCost:  equipmentCost,
		ManagementCost: managementCost,
		TotalCost:      totalCost,
		BudgetLimit:    req.BudgetLimit,
		Items:          string(itemsJSON),
	}

	if err := s.costRepo.Create(estimate); err != nil {
		return nil, fmt.Errorf("创建成本估算失败 / Failed to create cost estimate: %v", err)
	}

	return s.convertToResponse(estimate, project.Name), nil
}

// UpdateCostEstimate 更新成本估算
// Update cost estimate
func (s *costService) UpdateCostEstimate(req *request.UpdateCostEstimateRequest) (*response.CostEstimateResponse, error) {
	// 查找成本估算 / Find cost estimate
	estimate, err := s.costRepo.FindByID(req.ID)
	if err != nil {
		return nil, errors.New("成本估算不存在 / Cost estimate not found")
	}

	// 获取项目信息 / Get project info
	project, err := s.projectRepo.GetByID(estimate.ProjectID)
	if err != nil {
		return nil, errors.New("关联项目不存在 / Associated project not found")
	}

	// 计算各项成本 / Calculate costs
	materialCost := s.calculateMaterialCost(req.Items)
	laborCost := s.calculateLaborCost(project.Area, req.LaborRate)
	equipmentCost := s.calculateEquipmentCost(project.Area, req.EquipRate)

	// 计算基础成本(不含管理费) / Calculate base cost (without management)
	baseCost := materialCost + laborCost + equipmentCost
	managementCost := s.calculateManagementCost(baseCost, req.MgmtRate)

	// 计算总成本 / Calculate total cost
	totalCost := s.CalculateTotalCost(materialCost, laborCost, equipmentCost, managementCost)

	// 序列化明细项 / Serialize items
	itemsJSON, err := json.Marshal(req.Items)
	if err != nil {
		return nil, fmt.Errorf("序列化明细项失败 / Failed to serialize items: %v", err)
	}

	// 更新字段 / Update fields
	estimate.MaterialCost = materialCost
	estimate.LaborCost = laborCost
	estimate.EquipmentCost = equipmentCost
	estimate.ManagementCost = managementCost
	estimate.TotalCost = totalCost
	estimate.BudgetLimit = req.BudgetLimit
	estimate.Items = string(itemsJSON)

	if err := s.costRepo.Update(estimate); err != nil {
		return nil, fmt.Errorf("更新成本估算失败 / Failed to update cost estimate: %v", err)
	}

	return s.convertToResponse(estimate, project.Name), nil
}

// DeleteCostEstimate 删除成本估算
// Delete cost estimate
func (s *costService) DeleteCostEstimate(id uint) error {
	// 检查是否存在 / Check if exists
	_, err := s.costRepo.FindByID(id)
	if err != nil {
		return errors.New("成本估算不存在 / Cost estimate not found")
	}

	if err := s.costRepo.Delete(id); err != nil {
		return fmt.Errorf("删除成本估算失败 / Failed to delete cost estimate: %v", err)
	}

	return nil
}

// GetCostEstimateByID 根据ID获取成本估算
// Get cost estimate by ID
func (s *costService) GetCostEstimateByID(id uint) (*response.CostEstimateResponse, error) {
	estimate, err := s.costRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("成本估算不存在 / Cost estimate not found")
	}

	projectName := ""
	if estimate.Project != nil {
		projectName = estimate.Project.Name
	}

	return s.convertToResponse(estimate, projectName), nil
}

// GetCostEstimateByProjectID 根据项目ID获取成本估算
// Get cost estimate by project ID
func (s *costService) GetCostEstimateByProjectID(projectID uint) (*response.CostEstimateResponse, error) {
	estimate, err := s.costRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, errors.New("该项目暂无成本估算 / No cost estimate for this project")
	}

	projectName := ""
	if estimate.Project != nil {
		projectName = estimate.Project.Name
	}

	return s.convertToResponse(estimate, projectName), nil
}

// CalculateCost 计算成本(不保存)
// Calculate cost without saving
func (s *costService) CalculateCost(req *request.CalculateCostRequest) (*response.CostCalculateResponse, error) {
	// 计算各项成本 / Calculate costs
	materialCost := s.calculateMaterialCost(req.Items)
	laborCost := s.calculateLaborCost(req.Area, req.LaborRate)
	equipmentCost := s.calculateEquipmentCost(req.Area, req.EquipRate)

	// 计算基础成本(不含管理费) / Calculate base cost (without management)
	baseCost := materialCost + laborCost + equipmentCost
	managementCost := s.calculateManagementCost(baseCost, req.MgmtRate)

	// 计算总成本 / Calculate total cost
	totalCost := s.CalculateTotalCost(materialCost, laborCost, equipmentCost, managementCost)

	// 检查预算警告 / Check budget warning
	budgetExceeded, exceededAmount, budgetWarning := s.CheckBudgetWarning(totalCost, req.BudgetLimit)

	// 转换明细项 / Convert items
	items := make([]response.CostItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = response.CostItem{
			MaterialID:   item.MaterialID,
			MaterialName: item.MaterialName,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			TotalPrice:   item.TotalPrice,
			Category:     item.Category,
			Unit:         item.Unit,
		}
	}

	return &response.CostCalculateResponse{
		MaterialCost:   materialCost,
		LaborCost:      laborCost,
		EquipmentCost:  equipmentCost,
		ManagementCost: managementCost,
		TotalCost:      totalCost,
		BudgetLimit:    req.BudgetLimit,
		BudgetExceeded: budgetExceeded,
		ExceededAmount: exceededAmount,
		BudgetWarning:  budgetWarning,
		Items:          items,
	}, nil
}

// CalculateTotalCost 计算总成本
// Calculate total cost: materialCost + laborCost + equipmentCost + managementCost
// **Property 14: Cost Report Sum Consistency**
// **Validates: Requirements 10.3**
func (s *costService) CalculateTotalCost(materialCost, laborCost, equipmentCost, managementCost float64) float64 {
	total := materialCost + laborCost + equipmentCost + managementCost
	// 四舍五入到两位小数 / Round to 2 decimal places
	return math.Round(total*100) / 100
}

// CheckBudgetWarning 检查预算警告
// Check budget warning and return warning status
// **Property 15: Budget Warning Trigger**
// **Validates: Requirements 10.4**
func (s *costService) CheckBudgetWarning(totalCost, budgetLimit float64) (bool, float64, string) {
	// 如果未设置预算上限，不触发警告 / If no budget limit set, no warning
	if budgetLimit <= 0 {
		return false, 0, ""
	}

	// 检查是否超预算 / Check if over budget
	if totalCost > budgetLimit {
		exceededAmount := math.Round((totalCost-budgetLimit)*100) / 100
		exceededPercent := math.Round((exceededAmount/budgetLimit)*10000) / 100
		warning := fmt.Sprintf("警告：总成本(%.2f元)已超出预算上限(%.2f元)，超出%.2f元(%.2f%%)",
			totalCost, budgetLimit, exceededAmount, exceededPercent)
		return true, exceededAmount, warning
	}

	return false, 0, ""
}

// GetCostSummary 获取成本汇总
// Get cost summary
func (s *costService) GetCostSummary(userID uint) (*response.CostSummaryResponse, error) {
	summary, err := s.costRepo.GetSummary(userID)
	if err != nil {
		return nil, fmt.Errorf("获取成本汇总失败 / Failed to get cost summary: %v", err)
	}

	return &response.CostSummaryResponse{
		TotalProjects:     int(summary.TotalProjects),
		TotalBudget:       summary.TotalBudget,
		TotalCost:         summary.TotalCost,
		AvgCostPerProject: summary.AvgCostPerProject,
		OverBudgetCount:   int(summary.OverBudgetCount),
	}, nil
}

// calculateMaterialCost 计算材料费
// Calculate material cost from items
func (s *costService) calculateMaterialCost(items []request.CostItem) float64 {
	var total float64
	for _, item := range items {
		if item.Category == "material" || item.Category == "" {
			// 如果已有总价则使用总价，否则计算 / Use total price if available, otherwise calculate
			if item.TotalPrice > 0 {
				total += item.TotalPrice
			} else {
				total += item.Quantity * item.UnitPrice
			}
		}
	}
	return math.Round(total*100) / 100
}

// calculateLaborCost 计算人工费
// Calculate labor cost based on area and rate
func (s *costService) calculateLaborCost(area, laborRate float64) float64 {
	if laborRate <= 0 {
		// 默认人工费率：150元/平方米 / Default labor rate: 150 yuan per square meter
		laborRate = 150
	}
	return math.Round(area*laborRate*100) / 100
}

// calculateEquipmentCost 计算设备费
// Calculate equipment cost based on area and rate
func (s *costService) calculateEquipmentCost(area, equipRate float64) float64 {
	if equipRate <= 0 {
		// 默认设备费率：50元/平方米 / Default equipment rate: 50 yuan per square meter
		equipRate = 50
	}
	return math.Round(area*equipRate*100) / 100
}

// calculateManagementCost 计算管理费
// Calculate management cost as percentage of base cost
func (s *costService) calculateManagementCost(baseCost, mgmtRate float64) float64 {
	if mgmtRate <= 0 {
		// 默认管理费率：8% / Default management rate: 8%
		mgmtRate = 0.08
	}
	if mgmtRate > 1 {
		// 如果传入的是百分比数值(如8)，转换为小数 / Convert percentage to decimal
		mgmtRate = mgmtRate / 100
	}
	return math.Round(baseCost*mgmtRate*100) / 100
}

// convertToResponse 转换为响应对象
// Convert to response object
func (s *costService) convertToResponse(estimate *model.CostEstimate, projectName string) *response.CostEstimateResponse {
	// 解析明细项 / Parse items
	var items []response.CostItem
	if estimate.Items != "" {
		_ = json.Unmarshal([]byte(estimate.Items), &items)
	}

	// 检查预算警告 / Check budget warning
	budgetExceeded, exceededAmount, _ := s.CheckBudgetWarning(estimate.TotalCost, estimate.BudgetLimit)

	return &response.CostEstimateResponse{
		ID:             estimate.ID,
		ProjectID:      estimate.ProjectID,
		ProjectName:    projectName,
		MaterialCost:   estimate.MaterialCost,
		LaborCost:      estimate.LaborCost,
		EquipmentCost:  estimate.EquipmentCost,
		ManagementCost: estimate.ManagementCost,
		TotalCost:      estimate.TotalCost,
		BudgetLimit:    estimate.BudgetLimit,
		BudgetExceeded: budgetExceeded,
		ExceededAmount: exceededAmount,
		Items:          items,
		CreatedAt:      estimate.CreatedAt,
		UpdatedAt:      estimate.UpdatedAt,
	}
}
