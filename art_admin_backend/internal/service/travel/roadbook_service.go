package travel

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"errors"
	"time"
)

// RoadbookService 路书服务
type RoadbookService struct {
	repo *repository.RoadbookRepository
}

// NewRoadbookService 创建路书服务实例
func NewRoadbookService() *RoadbookService {
	return &RoadbookService{
		repo: repository.NewRoadbookRepository(),
	}
}

// GetRoadbookListRequest 获取路书列表请求
type GetRoadbookListRequest struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
	UserID     int64  `form:"-"`
	Keyword    string `form:"keyword"`
	TagID      int64  `form:"tagId"`
	Visibility *int   `form:"visibility"`
	Status     *int   `form:"status"`
}

// CreateRoadbookRequest 创建路书请求
type CreateRoadbookRequest struct {
	UserID      int64                   `json:"-"`
	Title       string                  `json:"title" binding:"required,max=100"`
	Description string                  `json:"description"`
	CoverURL    string                  `json:"coverUrl"`
	StartDate   *time.Time              `json:"startDate"`
	EndDate     *time.Time              `json:"endDate"`
	Visibility  int                     `json:"visibility"`
	TravelMode  string                  `json:"travelMode"`
	TotalBudget float64                 `json:"totalBudget"`
	Waypoints   []CreateWaypointRequest `json:"waypoints"`
	TagIDs      []int64                 `json:"tagIds"`
}

// CreateWaypointRequest 创建途经点请求
type CreateWaypointRequest struct {
	Name         string   `json:"name" binding:"required,max=100"`
	Address      string   `json:"address"`
	Longitude    float64  `json:"longitude" binding:"required"`
	Latitude     float64  `json:"latitude" binding:"required"`
	PoiID        string   `json:"poiId"`
	PoiType      string   `json:"poiType"`
	DayIndex     int      `json:"dayIndex"`
	SortOrder    int      `json:"sortOrder"`
	StayDuration int      `json:"stayDuration"`
	Budget       float64  `json:"budget"`
	Notes        string   `json:"notes"`
	Images       []string `json:"images"`
	WaypointType int      `json:"waypointType"`
}

// UpdateRoadbookRequest 更新路书请求
type UpdateRoadbookRequest struct {
	ID          int64      `json:"-"`
	UserID      int64      `json:"-"`
	Title       string     `json:"title" binding:"max=100"`
	Description string     `json:"description"`
	CoverURL    string     `json:"coverUrl"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Visibility  *int       `json:"visibility"`
	TravelMode  string     `json:"travelMode"`
	TotalBudget *float64   `json:"totalBudget"`
	Status      *int       `json:"status"`
	TagIDs      []int64    `json:"tagIds"`
}

// GetRoadbookList 获取路书列表
func (s *RoadbookService) GetRoadbookList(req *GetRoadbookListRequest) ([]model.Roadbook, int64, error) {
	return s.repo.FindWithPagination(req.Page, req.PageSize, req.UserID, req.Keyword, req.TagID, req.Visibility, req.Status)
}

// GetRoadbookDetail 获取路书详情
func (s *RoadbookService) GetRoadbookDetail(id int64, userID int64) (*model.Roadbook, error) {
	roadbook, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 检查访问权限
	if roadbook.Visibility == model.VisibilityPrivate && roadbook.UserID != userID {
		return nil, errors.New("forbidden")
	}

	// 增加浏览量
	_ = s.repo.IncrementViewCount(id)

	return roadbook, nil
}

// CreateRoadbook 创建路书
func (s *RoadbookService) CreateRoadbook(req *CreateRoadbookRequest) (*model.Roadbook, error) {
	roadbook := &model.Roadbook{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		CoverURL:    req.CoverURL,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Visibility:  req.Visibility,
		TravelMode:  req.TravelMode,
		TotalBudget: req.TotalBudget,
		Status:      model.RoadbookStatusDraft,
	}

	// 转换途经点
	waypoints := make([]model.RoadbookWaypoint, len(req.Waypoints))
	for i, wp := range req.Waypoints {
		waypoints[i] = model.RoadbookWaypoint{
			Name:         wp.Name,
			Address:      wp.Address,
			Longitude:    wp.Longitude,
			Latitude:     wp.Latitude,
			PoiID:        wp.PoiID,
			PoiType:      wp.PoiType,
			DayIndex:     wp.DayIndex,
			SortOrder:    wp.SortOrder,
			StayDuration: wp.StayDuration,
			Budget:       wp.Budget,
			Notes:        wp.Notes,
			Images:       wp.Images,
			WaypointType: wp.WaypointType,
		}
	}

	err := s.repo.CreateWithWaypoints(roadbook, waypoints, req.TagIDs)
	if err != nil {
		return nil, err
	}

	return roadbook, nil
}

// UpdateRoadbook 更新路书
func (s *RoadbookService) UpdateRoadbook(req *UpdateRoadbookRequest) error {
	roadbook, err := s.repo.FindByID(req.ID)
	if err != nil {
		return err
	}

	// 检查所有权
	if roadbook.UserID != req.UserID {
		return errors.New("forbidden")
	}

	// 更新字段
	if req.Title != "" {
		roadbook.Title = req.Title
	}
	if req.Description != "" {
		roadbook.Description = req.Description
	}
	if req.CoverURL != "" {
		roadbook.CoverURL = req.CoverURL
	}
	if req.StartDate != nil {
		roadbook.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		roadbook.EndDate = req.EndDate
	}
	if req.Visibility != nil {
		roadbook.Visibility = *req.Visibility
	}
	if req.TravelMode != "" {
		roadbook.TravelMode = req.TravelMode
	}
	if req.TotalBudget != nil {
		roadbook.TotalBudget = *req.TotalBudget
	}
	if req.Status != nil {
		roadbook.Status = *req.Status
	}

	return s.repo.Update(roadbook, req.TagIDs)
}

// DeleteRoadbook 删除路书
func (s *RoadbookService) DeleteRoadbook(id int64, userID int64) error {
	roadbook, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 检查所有权
	if roadbook.UserID != userID {
		return errors.New("forbidden")
	}

	return s.repo.Delete(id)
}
