package project_material

import (
	"encoding/json"
	"errors"
	"fmt"

	"art_admin_backend/internal/dto/request"
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/pkg/database"

	"gorm.io/gorm"
)

// ProjectMaterialService 项目材料清单服务
// Project material list service
type ProjectMaterialService struct {
	db *gorm.DB
}

// NewProjectMaterialService 创建项目材料清单服务实例
// Create project material service instance
func NewProjectMaterialService() *ProjectMaterialService {
	return &ProjectMaterialService{
		db: database.DB,
	}
}

// GetProjectMaterials 获取项目材料清单
// Get project material list
func (s *ProjectMaterialService) GetProjectMaterials(projectID uint) (*response.ProjectMaterialListResponse, error) {
	// 检查项目是否存在 / Check if project exists
	var project model.DesignerProject
	if err := s.db.First(&project, projectID).Error; err != nil {
		return nil, errors.New("项目不存在 / Project not found")
	}

	// 获取材料清单 / Get material list
	var materials []model.ProjectMaterial
	if err := s.db.Where("project_id = ?", projectID).Order("created_at ASC").Find(&materials).Error; err != nil {
		return nil, fmt.Errorf("获取材料清单失败 / Failed to get material list: %v", err)
	}

	// 计算材料费合计 / Calculate total material cost
	var materialCost float64
	items := make([]response.ProjectMaterialItem, len(materials))
	for i, m := range materials {
		items[i] = response.ProjectMaterialItem{
			ID:            m.ID,
			ProjectID:     m.ProjectID,
			MaterialID:    m.MaterialID,
			Name:          m.Name,
			Category:      m.Category,
			Specification: m.Specification,
			Unit:          m.Unit,
			UnitPrice:     m.UnitPrice,
			Quantity:      m.Quantity,
			TotalPrice:    m.TotalPrice,
			Brand:         m.Brand,
			Supplier:      m.Supplier,
			Remark:        m.Remark,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
		}
		materialCost += m.TotalPrice
	}

	return &response.ProjectMaterialListResponse{
		ProjectID:    projectID,
		ProjectName:  project.Name,
		Items:        items,
		MaterialCost: materialCost,
		ItemCount:    len(items),
	}, nil
}

// SaveProjectMaterials 保存项目材料清单
// Save project material list (replace all)
func (s *ProjectMaterialService) SaveProjectMaterials(req *request.SaveProjectMaterialRequest) (*response.ProjectMaterialListResponse, error) {
	// 检查项目是否存在 / Check if project exists
	var project model.DesignerProject
	if err := s.db.First(&project, req.ProjectID).Error; err != nil {
		return nil, errors.New("项目不存在 / Project not found")
	}

	// 开启事务 / Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除现有材料清单 / Delete existing materials
	if err := tx.Where("project_id = ?", req.ProjectID).Delete(&model.ProjectMaterial{}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("删除现有材料清单失败 / Failed to delete existing materials: %v", err)
	}

	// 插入新材料清单 / Insert new materials
	var materialCost float64
	items := make([]response.ProjectMaterialItem, len(req.Items))
	for i, item := range req.Items {
		material := model.ProjectMaterial{
			ProjectID:     req.ProjectID,
			MaterialID:    item.MaterialID,
			Name:          item.Name,
			Category:      item.Category,
			Specification: item.Specification,
			Unit:          item.Unit,
			UnitPrice:     item.UnitPrice,
			Quantity:      item.Quantity,
			TotalPrice:    item.UnitPrice * item.Quantity,
			Brand:         item.Brand,
			Supplier:      item.Supplier,
			Remark:        item.Remark,
		}

		if err := tx.Create(&material).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("保存材料失败 / Failed to save material: %v", err)
		}

		items[i] = response.ProjectMaterialItem{
			ID:            material.ID,
			ProjectID:     material.ProjectID,
			MaterialID:    material.MaterialID,
			Name:          material.Name,
			Category:      material.Category,
			Specification: material.Specification,
			Unit:          material.Unit,
			UnitPrice:     material.UnitPrice,
			Quantity:      material.Quantity,
			TotalPrice:    material.TotalPrice,
			Brand:         material.Brand,
			Supplier:      material.Supplier,
			Remark:        material.Remark,
			CreatedAt:     material.CreatedAt,
			UpdatedAt:     material.UpdatedAt,
		}
		materialCost += material.TotalPrice
	}

	// 提交事务 / Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败 / Failed to commit transaction: %v", err)
	}

	return &response.ProjectMaterialListResponse{
		ProjectID:    req.ProjectID,
		ProjectName:  project.Name,
		Items:        items,
		MaterialCost: materialCost,
		ItemCount:    len(items),
	}, nil
}

// AddProjectMaterial 添加单个材料到项目
// Add single material to project
func (s *ProjectMaterialService) AddProjectMaterial(projectID uint, req *request.ProjectMaterialItem) (*response.ProjectMaterialItem, error) {
	// 检查项目是否存在 / Check if project exists
	var project model.DesignerProject
	if err := s.db.First(&project, projectID).Error; err != nil {
		return nil, errors.New("项目不存在 / Project not found")
	}

	material := model.ProjectMaterial{
		ProjectID:     projectID,
		MaterialID:    req.MaterialID,
		Name:          req.Name,
		Category:      req.Category,
		Specification: req.Specification,
		Unit:          req.Unit,
		UnitPrice:     req.UnitPrice,
		Quantity:      req.Quantity,
		TotalPrice:    req.UnitPrice * req.Quantity,
		Brand:         req.Brand,
		Supplier:      req.Supplier,
		Remark:        req.Remark,
	}

	if err := s.db.Create(&material).Error; err != nil {
		return nil, fmt.Errorf("添加材料失败 / Failed to add material: %v", err)
	}

	return &response.ProjectMaterialItem{
		ID:            material.ID,
		ProjectID:     material.ProjectID,
		MaterialID:    material.MaterialID,
		Name:          material.Name,
		Category:      material.Category,
		Specification: material.Specification,
		Unit:          material.Unit,
		UnitPrice:     material.UnitPrice,
		Quantity:      material.Quantity,
		TotalPrice:    material.TotalPrice,
		Brand:         material.Brand,
		Supplier:      material.Supplier,
		Remark:        material.Remark,
		CreatedAt:     material.CreatedAt,
		UpdatedAt:     material.UpdatedAt,
	}, nil
}

// UpdateProjectMaterial 更新项目材料项
// Update project material item
func (s *ProjectMaterialService) UpdateProjectMaterial(projectID, itemID uint, req *request.ProjectMaterialItem) (*response.ProjectMaterialItem, error) {
	var material model.ProjectMaterial
	if err := s.db.Where("id = ? AND project_id = ?", itemID, projectID).First(&material).Error; err != nil {
		return nil, errors.New("材料项不存在 / Material item not found")
	}

	// 更新字段 / Update fields
	material.Name = req.Name
	material.Category = req.Category
	material.Specification = req.Specification
	material.Unit = req.Unit
	material.UnitPrice = req.UnitPrice
	material.Quantity = req.Quantity
	material.TotalPrice = req.UnitPrice * req.Quantity
	material.Brand = req.Brand
	material.Supplier = req.Supplier
	material.Remark = req.Remark

	if err := s.db.Save(&material).Error; err != nil {
		return nil, fmt.Errorf("更新材料失败 / Failed to update material: %v", err)
	}

	return &response.ProjectMaterialItem{
		ID:            material.ID,
		ProjectID:     material.ProjectID,
		MaterialID:    material.MaterialID,
		Name:          material.Name,
		Category:      material.Category,
		Specification: material.Specification,
		Unit:          material.Unit,
		UnitPrice:     material.UnitPrice,
		Quantity:      material.Quantity,
		TotalPrice:    material.TotalPrice,
		Brand:         material.Brand,
		Supplier:      material.Supplier,
		Remark:        material.Remark,
		CreatedAt:     material.CreatedAt,
		UpdatedAt:     material.UpdatedAt,
	}, nil
}

// DeleteProjectMaterial 删除项目材料项
// Delete project material item
func (s *ProjectMaterialService) DeleteProjectMaterial(projectID, itemID uint) error {
	result := s.db.Where("id = ? AND project_id = ?", itemID, projectID).Delete(&model.ProjectMaterial{})
	if result.Error != nil {
		return fmt.Errorf("删除材料失败 / Failed to delete material: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("材料项不存在 / Material item not found")
	}
	return nil
}

// GetProjectCostSummary 获取项目成本汇总
// Get project cost summary
func (s *ProjectMaterialService) GetProjectCostSummary(projectID uint) (*response.ProjectCostSummary, error) {
	// 检查项目是否存在 / Check if project exists
	var project model.DesignerProject
	if err := s.db.First(&project, projectID).Error; err != nil {
		return nil, errors.New("项目不存在 / Project not found")
	}

	// 计算材料费 / Calculate material cost
	var materialCost float64
	s.db.Model(&model.ProjectMaterial{}).Where("project_id = ?", projectID).Select("COALESCE(SUM(total_price), 0)").Scan(&materialCost)

	// 获取成本估算配置 / Get cost estimate config
	var costEstimate model.CostEstimate
	s.db.Where("project_id = ?", projectID).First(&costEstimate)

	// 解析自定义费用项 / Parse custom cost items
	var customCosts []response.CustomCostItem
	var customCostTotal float64
	if costEstimate.Items != "" {
		if err := json.Unmarshal([]byte(costEstimate.Items), &customCosts); err == nil {
			for _, item := range customCosts {
				customCostTotal += item.Amount
			}
		}
	}

	totalCost := materialCost + costEstimate.LaborCost + costEstimate.EquipmentCost + costEstimate.ManagementCost + customCostTotal
	budgetExceeded := costEstimate.BudgetLimit > 0 && totalCost > costEstimate.BudgetLimit
	exceededAmount := float64(0)
	if budgetExceeded {
		exceededAmount = totalCost - costEstimate.BudgetLimit
	}

	return &response.ProjectCostSummary{
		ProjectID:       projectID,
		MaterialCost:    materialCost,
		LaborCost:       costEstimate.LaborCost,
		EquipmentCost:   costEstimate.EquipmentCost,
		ManagementCost:  costEstimate.ManagementCost,
		CustomCosts:     customCosts,
		CustomCostTotal: customCostTotal,
		TotalCost:       totalCost,
		BudgetLimit:     costEstimate.BudgetLimit,
		BudgetExceeded:  budgetExceeded,
		ExceededAmount:  exceededAmount,
	}, nil
}

// SaveProjectCost 保存项目成本配置
// Save project cost config
func (s *ProjectMaterialService) SaveProjectCost(req *request.SaveProjectCostRequest) (*response.ProjectCostSummary, error) {
	// 检查项目是否存在 / Check if project exists
	var project model.DesignerProject
	if err := s.db.First(&project, req.ProjectID).Error; err != nil {
		return nil, errors.New("项目不存在 / Project not found")
	}

	// 计算材料费 / Calculate material cost
	var materialCost float64
	s.db.Model(&model.ProjectMaterial{}).Where("project_id = ?", req.ProjectID).Select("COALESCE(SUM(total_price), 0)").Scan(&materialCost)

	// 计算自定义费用小计 / Calculate custom cost total
	var customCostTotal float64
	var customCostsResp []response.CustomCostItem
	for _, item := range req.CustomCosts {
		customCostTotal += item.Amount
		customCostsResp = append(customCostsResp, response.CustomCostItem{
			Name:   item.Name,
			Amount: item.Amount,
		})
	}

	// 序列化自定义费用项 / Serialize custom costs
	customCostsJSON := ""
	if len(req.CustomCosts) > 0 {
		jsonBytes, err := json.Marshal(req.CustomCosts)
		if err == nil {
			customCostsJSON = string(jsonBytes)
		}
	}

	// 计算总成本 / Calculate total cost
	totalCost := materialCost + req.LaborCost + req.EquipmentCost + req.ManagementCost + customCostTotal

	// 查找或创建成本估算 / Find or create cost estimate
	var costEstimate model.CostEstimate
	result := s.db.Where("project_id = ?", req.ProjectID).First(&costEstimate)

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新记录 / Create new record
		costEstimate = model.CostEstimate{
			ProjectID:      req.ProjectID,
			MaterialCost:   materialCost,
			LaborCost:      req.LaborCost,
			EquipmentCost:  req.EquipmentCost,
			ManagementCost: req.ManagementCost,
			TotalCost:      totalCost,
			BudgetLimit:    req.BudgetLimit,
			Items:          customCostsJSON,
		}
		if err := s.db.Create(&costEstimate).Error; err != nil {
			return nil, fmt.Errorf("创建成本配置失败 / Failed to create cost config: %v", err)
		}
	} else if result.Error != nil {
		return nil, fmt.Errorf("查询成本配置失败 / Failed to query cost config: %v", result.Error)
	} else {
		// 更新现有记录 / Update existing record
		costEstimate.MaterialCost = materialCost
		costEstimate.LaborCost = req.LaborCost
		costEstimate.EquipmentCost = req.EquipmentCost
		costEstimate.ManagementCost = req.ManagementCost
		costEstimate.TotalCost = totalCost
		costEstimate.BudgetLimit = req.BudgetLimit
		costEstimate.Items = customCostsJSON
		if err := s.db.Save(&costEstimate).Error; err != nil {
			return nil, fmt.Errorf("更新成本配置失败 / Failed to update cost config: %v", err)
		}
	}

	budgetExceeded := req.BudgetLimit > 0 && totalCost > req.BudgetLimit
	exceededAmount := float64(0)
	if budgetExceeded {
		exceededAmount = totalCost - req.BudgetLimit
	}

	return &response.ProjectCostSummary{
		ProjectID:       req.ProjectID,
		MaterialCost:    materialCost,
		LaborCost:       req.LaborCost,
		EquipmentCost:   req.EquipmentCost,
		ManagementCost:  req.ManagementCost,
		CustomCosts:     customCostsResp,
		CustomCostTotal: customCostTotal,
		TotalCost:       totalCost,
		BudgetLimit:     req.BudgetLimit,
		BudgetExceeded:  budgetExceeded,
		ExceededAmount:  exceededAmount,
	}, nil
}

// ImportMaterialsToProject 从材料库导入材料到项目
// Import materials from library to project
func (s *ProjectMaterialService) ImportMaterialsToProject(projectID uint, materialIDs []uint, defaultQuantity float64) (*response.ProjectMaterialListResponse, error) {
	// 检查项目是否存在 / Check if project exists
	var project model.DesignerProject
	if err := s.db.First(&project, projectID).Error; err != nil {
		return nil, errors.New("项目不存在 / Project not found")
	}

	// 获取材料库中的材料 / Get materials from library
	var materials []model.Material
	if err := s.db.Where("id IN ?", materialIDs).Find(&materials).Error; err != nil {
		return nil, fmt.Errorf("获取材料失败 / Failed to get materials: %v", err)
	}

	if defaultQuantity <= 0 {
		defaultQuantity = 1
	}

	// 添加到项目材料清单 / Add to project material list
	for _, m := range materials {
		projectMaterial := model.ProjectMaterial{
			ProjectID:     projectID,
			MaterialID:    &m.ID,
			Name:          m.Name,
			Category:      m.Category,
			Specification: m.Specification,
			Unit:          m.Unit,
			UnitPrice:     m.UnitPrice,
			Quantity:      defaultQuantity,
			TotalPrice:    m.UnitPrice * defaultQuantity,
			Brand:         m.Brand,
			Supplier:      m.Supplier,
		}
		s.db.Create(&projectMaterial)
	}

	// 返回更新后的材料清单 / Return updated material list
	return s.GetProjectMaterials(projectID)
}
