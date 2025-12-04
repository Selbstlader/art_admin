package travel

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"errors"
)

// WaypointService 途经点服务
type WaypointService struct {
	waypointRepo *repository.WaypointRepository
	roadbookRepo *repository.RoadbookRepository
}

// NewWaypointService 创建途经点服务实例
func NewWaypointService() *WaypointService {
	return &WaypointService{
		waypointRepo: repository.NewWaypointRepository(),
		roadbookRepo: repository.NewRoadbookRepository(),
	}
}

// NewWaypointServiceWithRepos 创建带自定义仓库的途经点服务实例（用于测试）
func NewWaypointServiceWithRepos(waypointRepo *repository.WaypointRepository, roadbookRepo *repository.RoadbookRepository) *WaypointService {
	return &WaypointService{
		waypointRepo: waypointRepo,
		roadbookRepo: roadbookRepo,
	}
}

// AddWaypointRequest 添加途经点请求
type AddWaypointRequest struct {
	RoadbookID   int64    `json:"-"`
	UserID       int64    `json:"-"`
	Name         string   `json:"name" binding:"required,max=100"`
	Address      string   `json:"address"`
	Longitude    float64  `json:"longitude" binding:"required"`
	Latitude     float64  `json:"latitude" binding:"required"`
	PoiID        string   `json:"poiId"`
	PoiType      string   `json:"poiType"`
	DayIndex     int      `json:"dayIndex"`
	SortOrder    *int     `json:"sortOrder"`
	StayDuration int      `json:"stayDuration"`
	Budget       float64  `json:"budget"`
	Notes        string   `json:"notes"`
	Images       []string `json:"images"`
	WaypointType int      `json:"waypointType"`
}

// UpdateWaypointRequest 更新途经点请求
type UpdateWaypointRequest struct {
	ID           int64    `json:"-"`
	RoadbookID   int64    `json:"-"`
	UserID       int64    `json:"-"`
	Name         string   `json:"name" binding:"max=100"`
	Address      string   `json:"address"`
	Longitude    *float64 `json:"longitude"`
	Latitude     *float64 `json:"latitude"`
	PoiID        string   `json:"poiId"`
	PoiType      string   `json:"poiType"`
	DayIndex     *int     `json:"dayIndex"`
	SortOrder    *int     `json:"sortOrder"`
	StayDuration *int     `json:"stayDuration"`
	Budget       *float64 `json:"budget"`
	Notes        string   `json:"notes"`
	Images       []string `json:"images"`
	WaypointType *int     `json:"waypointType"`
}

// BatchUpdateWaypointsRequest 批量更新途经点请求
type BatchUpdateWaypointsRequest struct {
	RoadbookID int64                   `json:"-"`
	UserID     int64                   `json:"-"`
	Waypoints  []UpdateWaypointRequest `json:"waypoints" binding:"required"`
}

// UpdateSortOrderRequest 更新排序请求
type UpdateSortOrderRequest struct {
	RoadbookID int64           `json:"-"`
	UserID     int64           `json:"-"`
	Orders     []SortOrderItem `json:"orders" binding:"required"`
}

// SortOrderItem 排序项
type SortOrderItem struct {
	ID        int64 `json:"id" binding:"required"`
	SortOrder int   `json:"sortOrder"`
}

// WaypointsByDay 按天分组的途经点
type WaypointsByDay struct {
	DayIndex  int                      `json:"dayIndex"`
	Waypoints []model.RoadbookWaypoint `json:"waypoints"`
}

// AddWaypoint 添加途经点
func (s *WaypointService) AddWaypoint(req *AddWaypointRequest) (*model.RoadbookWaypoint, error) {
	// 验证路书存在且用户有权限
	roadbook, err := s.roadbookRepo.GetByID(req.RoadbookID)
	if err != nil {
		return nil, errors.New("roadbook not found")
	}
	if roadbook.UserID != req.UserID {
		return nil, errors.New("forbidden")
	}

	// 如果没有指定排序，获取当前最大排序值+1
	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	} else {
		maxOrder, err := s.waypointRepo.GetMaxSortOrder(req.RoadbookID, req.DayIndex)
		if err == nil {
			sortOrder = maxOrder + 1
		}
	}

	// 设置默认值
	dayIndex := req.DayIndex
	if dayIndex <= 0 {
		dayIndex = 1
	}
	stayDuration := req.StayDuration
	if stayDuration <= 0 {
		stayDuration = 60 // 默认60分钟
	}
	waypointType := req.WaypointType
	if waypointType <= 0 {
		waypointType = model.WaypointTypeWaypoint
	}

	waypoint := &model.RoadbookWaypoint{
		RoadbookID:   req.RoadbookID,
		Name:         req.Name,
		Address:      req.Address,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		PoiID:        req.PoiID,
		PoiType:      req.PoiType,
		DayIndex:     dayIndex,
		SortOrder:    sortOrder,
		StayDuration: stayDuration,
		Budget:       req.Budget,
		Notes:        req.Notes,
		Images:       req.Images,
		WaypointType: waypointType,
	}

	if err := s.waypointRepo.Create(waypoint); err != nil {
		return nil, err
	}

	return waypoint, nil
}

// UpdateWaypoint 更新途经点
func (s *WaypointService) UpdateWaypoint(req *UpdateWaypointRequest) error {
	// 验证路书存在且用户有权限
	roadbook, err := s.roadbookRepo.GetByID(req.RoadbookID)
	if err != nil {
		return errors.New("roadbook not found")
	}
	if roadbook.UserID != req.UserID {
		return errors.New("forbidden")
	}

	// 获取途经点
	waypoint, err := s.waypointRepo.GetByID(req.ID)
	if err != nil {
		return errors.New("waypoint not found")
	}
	if waypoint.RoadbookID != req.RoadbookID {
		return errors.New("waypoint does not belong to this roadbook")
	}

	// 更新字段
	if req.Name != "" {
		waypoint.Name = req.Name
	}
	if req.Address != "" {
		waypoint.Address = req.Address
	}
	if req.Longitude != nil {
		waypoint.Longitude = *req.Longitude
	}
	if req.Latitude != nil {
		waypoint.Latitude = *req.Latitude
	}
	if req.PoiID != "" {
		waypoint.PoiID = req.PoiID
	}
	if req.PoiType != "" {
		waypoint.PoiType = req.PoiType
	}
	if req.DayIndex != nil {
		waypoint.DayIndex = *req.DayIndex
	}
	if req.SortOrder != nil {
		waypoint.SortOrder = *req.SortOrder
	}
	if req.StayDuration != nil {
		waypoint.StayDuration = *req.StayDuration
	}
	if req.Budget != nil {
		waypoint.Budget = *req.Budget
	}
	if req.Notes != "" {
		waypoint.Notes = req.Notes
	}
	if req.Images != nil {
		waypoint.Images = req.Images
	}
	if req.WaypointType != nil {
		waypoint.WaypointType = *req.WaypointType
	}

	return s.waypointRepo.Update(waypoint)
}

// DeleteWaypoint 删除途经点
func (s *WaypointService) DeleteWaypoint(roadbookID, waypointID, userID int64) error {
	// 验证路书存在且用户有权限
	roadbook, err := s.roadbookRepo.GetByID(roadbookID)
	if err != nil {
		return errors.New("roadbook not found")
	}
	if roadbook.UserID != userID {
		return errors.New("forbidden")
	}

	// 验证途经点存在且属于该路书
	waypoint, err := s.waypointRepo.GetByID(waypointID)
	if err != nil {
		return errors.New("waypoint not found")
	}
	if waypoint.RoadbookID != roadbookID {
		return errors.New("waypoint does not belong to this roadbook")
	}

	return s.waypointRepo.Delete(waypointID)
}

// GetWaypointsByRoadbook 获取路书的所有途经点
func (s *WaypointService) GetWaypointsByRoadbook(roadbookID, userID int64) ([]model.RoadbookWaypoint, error) {
	// 验证路书存在
	roadbook, err := s.roadbookRepo.GetByID(roadbookID)
	if err != nil {
		return nil, errors.New("roadbook not found")
	}

	// 检查访问权限（私有路书只有所有者可以查看）
	if roadbook.Visibility == model.VisibilityPrivate && roadbook.UserID != userID {
		return nil, errors.New("forbidden")
	}

	return s.waypointRepo.ListByRoadbookID(roadbookID)
}

// GetWaypointsGroupedByDay 获取路书的途经点并按天分组
func (s *WaypointService) GetWaypointsGroupedByDay(roadbookID, userID int64) ([]WaypointsByDay, error) {
	// 验证路书存在
	roadbook, err := s.roadbookRepo.GetByID(roadbookID)
	if err != nil {
		return nil, errors.New("roadbook not found")
	}

	// 检查访问权限
	if roadbook.Visibility == model.VisibilityPrivate && roadbook.UserID != userID {
		return nil, errors.New("forbidden")
	}

	// 获取按天分组的途经点
	grouped, err := s.waypointRepo.ListByRoadbookIDGroupedByDay(roadbookID)
	if err != nil {
		return nil, err
	}

	// 转换为有序的结果
	result := make([]WaypointsByDay, 0, len(grouped))
	for dayIndex, waypoints := range grouped {
		result = append(result, WaypointsByDay{
			DayIndex:  dayIndex,
			Waypoints: waypoints,
		})
	}

	// 按天数排序
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].DayIndex > result[j].DayIndex {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

// BatchUpdateWaypoints 批量更新途经点
func (s *WaypointService) BatchUpdateWaypoints(req *BatchUpdateWaypointsRequest) error {
	// 验证路书存在且用户有权限
	roadbook, err := s.roadbookRepo.GetByID(req.RoadbookID)
	if err != nil {
		return errors.New("roadbook not found")
	}
	if roadbook.UserID != req.UserID {
		return errors.New("forbidden")
	}

	// 更新每个途经点
	for _, wpReq := range req.Waypoints {
		wpReq.RoadbookID = req.RoadbookID
		wpReq.UserID = req.UserID
		if err := s.UpdateWaypoint(&wpReq); err != nil {
			return err
		}
	}

	return nil
}

// UpdateSortOrder 更新途经点排序
func (s *WaypointService) UpdateSortOrder(req *UpdateSortOrderRequest) error {
	// 验证路书存在且用户有权限
	roadbook, err := s.roadbookRepo.GetByID(req.RoadbookID)
	if err != nil {
		return errors.New("roadbook not found")
	}
	if roadbook.UserID != req.UserID {
		return errors.New("forbidden")
	}

	// 构建更新映射
	updates := make(map[int64]int)
	for _, item := range req.Orders {
		updates[item.ID] = item.SortOrder
	}

	return s.waypointRepo.BatchUpdateSortOrder(updates)
}
