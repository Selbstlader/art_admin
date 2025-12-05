package travel

import (
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"errors"

	"gorm.io/gorm"
)

// CommentService 评论服务
type CommentService struct {
	commentRepo  *repository.CommentRepository
	roadbookRepo *repository.RoadbookRepository
}

// NewCommentService 创建评论服务实例
func NewCommentService() *CommentService {
	return &CommentService{
		commentRepo:  repository.NewCommentRepository(),
		roadbookRepo: repository.NewRoadbookRepository(),
	}
}

// NewCommentServiceWithDB 创建带自定义DB的评论服务实例（用于测试）
func NewCommentServiceWithDB(db *gorm.DB) *CommentService {
	return &CommentService{
		commentRepo:  repository.NewCommentRepositoryWithDB(db),
		roadbookRepo: repository.NewRoadbookRepositoryWithDB(db),
	}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	RoadbookID int64  `json:"roadbookId" binding:"required"`
	UserID     int64  `json:"-"`
	ParentID   int64  `json:"parentId"`
	Content    string `json:"content" binding:"required"`
}

// GetCommentListRequest 获取评论列表请求
type GetCommentListRequest struct {
	RoadbookID int64 `form:"roadbookId" binding:"required"`
	Page       int   `form:"page"`
	PageSize   int   `form:"pageSize"`
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(req *CreateCommentRequest) (*model.RoadbookComment, error) {
	// 验证评论内容
	if !model.ValidateCommentContent(req.Content) {
		return nil, errors.New("invalid_content")
	}

	// 验证路书是否存在
	exists, err := s.roadbookRepo.ExistsByID(req.RoadbookID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("roadbook_not_found")
	}

	// 如果是回复评论，验证父评论是否存在且属于同一路书
	if req.ParentID > 0 {
		exists, err := s.commentRepo.ExistsByIDAndRoadbookID(req.ParentID, req.RoadbookID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("parent_comment_not_found")
		}
	}

	comment := &model.RoadbookComment{
		RoadbookID: req.RoadbookID,
		UserID:     req.UserID,
		ParentID:   req.ParentID,
		Content:    req.Content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// 更新路书评论计数
	_ = s.roadbookRepo.IncrementCommentCount(req.RoadbookID, 1)

	return comment, nil
}

// GetCommentList 获取评论列表（树形结构）
func (s *CommentService) GetCommentList(req *GetCommentListRequest) ([]model.RoadbookComment, int64, error) {
	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	return s.commentRepo.GetCommentTree(req.RoadbookID, req.Page, req.PageSize)
}

// GetFlatCommentList 获取评论列表（扁平结构，按时间倒序）
func (s *CommentService) GetFlatCommentList(roadbookID int64, page, pageSize int) ([]model.RoadbookComment, int64, error) {
	return s.commentRepo.ListByRoadbookID(roadbookID, page, pageSize)
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(commentID, userID int64) error {
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return errors.New("comment_not_found")
	}

	// 检查所有权
	if comment.UserID != userID {
		return errors.New("forbidden")
	}

	if err := s.commentRepo.SoftDelete(commentID); err != nil {
		return err
	}

	// 更新路书评论计数
	_ = s.roadbookRepo.IncrementCommentCount(comment.RoadbookID, -1)

	return nil
}

// LikeComment 点赞评论
func (s *CommentService) LikeComment(commentID int64) error {
	exists, err := s.commentRepo.ExistsByID(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment_not_found")
	}

	return s.commentRepo.IncrementLikeCount(commentID, 1)
}

// UnlikeComment 取消点赞评论
func (s *CommentService) UnlikeComment(commentID int64) error {
	exists, err := s.commentRepo.ExistsByID(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment_not_found")
	}

	return s.commentRepo.IncrementLikeCount(commentID, -1)
}

// GetCommentByID 根据ID获取评论
func (s *CommentService) GetCommentByID(id int64) (*model.RoadbookComment, error) {
	return s.commentRepo.GetByID(id)
}

// ValidateParentComment 验证父评论是否有效（存在且属于同一路书）
func (s *CommentService) ValidateParentComment(parentID, roadbookID int64) (bool, error) {
	if parentID == 0 {
		return true, nil
	}
	return s.commentRepo.ExistsByIDAndRoadbookID(parentID, roadbookID)
}
