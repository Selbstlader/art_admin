package mood_analytics

import (
	"art_admin_backend/internal/dto/response"
	"art_admin_backend/internal/repository"
	"errors"
)

// MeditationFavoriteService 冥想收藏服务
type MeditationFavoriteService struct {
	favoriteRepo   *repository.MeditationFavoriteRepository
	playRecordRepo *repository.MeditationPlayRecordRepository
	contentRepo    *repository.MeditationContentRepository
}

// NewMeditationFavoriteService 创建冥想收藏服务
func NewMeditationFavoriteService() *MeditationFavoriteService {
	return &MeditationFavoriteService{
		favoriteRepo:   repository.NewMeditationFavoriteRepository(),
		playRecordRepo: repository.NewMeditationPlayRecordRepository(),
		contentRepo:    repository.NewMeditationContentRepository(),
	}
}

// AddFavorite 添加收藏
func (s *MeditationFavoriteService) AddFavorite(userID, contentID int64) error {
	// 验证内容是否存在
	_, err := s.contentRepo.FindByID(contentID)
	if err != nil {
		return errors.New("冥想内容不存在")
	}

	return s.favoriteRepo.AddFavorite(userID, contentID)
}

// RemoveFavorite 取消收藏
func (s *MeditationFavoriteService) RemoveFavorite(userID, contentID int64) error {
	return s.favoriteRepo.RemoveFavorite(userID, contentID)
}

// IsFavorite 检查是否已收藏
func (s *MeditationFavoriteService) IsFavorite(userID, contentID int64) (bool, error) {
	return s.favoriteRepo.IsFavorite(userID, contentID)
}

// GetUserFavorites 获取用户收藏列表
func (s *MeditationFavoriteService) GetUserFavorites(userID int64, page, pageSize int) ([]response.MeditationContentItem, int64, error) {
	favorites, total, err := s.favoriteRepo.GetUserFavorites(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	items := make([]response.MeditationContentItem, 0, len(favorites))
	for _, fav := range favorites {
		if fav.Content != nil {
			items = append(items, response.MeditationContentItem{
				ID:              fav.Content.ID,
				Title:           fav.Content.Title,
				Category:        fav.Content.Category,
				DifficultyLevel: fav.Content.DifficultyLevel,
				Duration:        fav.Content.Duration,
				AudioURL:        fav.Content.AudioURL,
				CoverImage:      fav.Content.CoverImage,
				ViewCount:       int64(fav.Content.ViewCount),
				LikeCount:       int64(fav.Content.LikeCount),
			})
		}
	}

	return items, total, nil
}

// UpdatePlayRecord 更新播放记录
func (s *MeditationFavoriteService) UpdatePlayRecord(userID, contentID int64, duration int, progress float64, completed bool) error {
	// 验证内容是否存在
	_, err := s.contentRepo.FindByID(contentID)
	if err != nil {
		return errors.New("冥想内容不存在")
	}

	return s.playRecordRepo.CreateOrUpdatePlayRecord(userID, contentID, duration, progress, completed)
}

// GetPlayRecord 获取播放记录
func (s *MeditationFavoriteService) GetPlayRecord(userID, contentID int64) (*response.PlayRecordItem, error) {
	record, err := s.playRecordRepo.GetPlayRecord(userID, contentID)
	if err != nil {
		return nil, err
	}

	return &response.PlayRecordItem{
		ID:           record.ID,
		ContentID:    record.ContentID,
		Duration:     record.Duration,
		Progress:     record.Progress,
		Completed:    record.Completed,
		LastPlayedAt: record.LastPlayedAt,
		PlayCount:    record.PlayCount,
	}, nil
}

// GetUserPlayHistory 获取用户播放历史
func (s *MeditationFavoriteService) GetUserPlayHistory(userID int64, page, pageSize int) ([]response.PlayHistoryItem, int64, error) {
	records, total, err := s.playRecordRepo.GetUserPlayHistory(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	items := make([]response.PlayHistoryItem, 0, len(records))
	for _, record := range records {
		item := response.PlayHistoryItem{
			ID:           record.ID,
			ContentID:    record.ContentID,
			Duration:     record.Duration,
			Progress:     record.Progress,
			Completed:    record.Completed,
			LastPlayedAt: record.LastPlayedAt,
			PlayCount:    record.PlayCount,
		}
		if record.Content != nil {
			item.ContentTitle = record.Content.Title
			item.ContentImage = record.Content.CoverImage
			item.ContentDuration = record.Content.Duration
		}
		items = append(items, item)
	}

	return items, total, nil
}

// GetRecentlyPlayed 获取最近播放
func (s *MeditationFavoriteService) GetRecentlyPlayed(userID int64, limit int) ([]response.PlayHistoryItem, error) {
	records, err := s.playRecordRepo.GetRecentlyPlayed(userID, limit)
	if err != nil {
		return nil, err
	}

	items := make([]response.PlayHistoryItem, 0, len(records))
	for _, record := range records {
		item := response.PlayHistoryItem{
			ID:           record.ID,
			ContentID:    record.ContentID,
			Duration:     record.Duration,
			Progress:     record.Progress,
			Completed:    record.Completed,
			LastPlayedAt: record.LastPlayedAt,
			PlayCount:    record.PlayCount,
		}
		if record.Content != nil {
			item.ContentTitle = record.Content.Title
			item.ContentImage = record.Content.CoverImage
			item.ContentDuration = record.Content.Duration
		}
		items = append(items, item)
	}

	return items, nil
}

// GetContinuePlaying 获取继续播放列表
func (s *MeditationFavoriteService) GetContinuePlaying(userID int64, limit int) ([]response.PlayHistoryItem, error) {
	records, err := s.playRecordRepo.GetContinuePlaying(userID, limit)
	if err != nil {
		return nil, err
	}

	items := make([]response.PlayHistoryItem, 0, len(records))
	for _, record := range records {
		item := response.PlayHistoryItem{
			ID:           record.ID,
			ContentID:    record.ContentID,
			Duration:     record.Duration,
			Progress:     record.Progress,
			Completed:    record.Completed,
			LastPlayedAt: record.LastPlayedAt,
			PlayCount:    record.PlayCount,
		}
		if record.Content != nil {
			item.ContentTitle = record.Content.Title
			item.ContentImage = record.Content.CoverImage
			item.ContentDuration = record.Content.Duration
		}
		items = append(items, item)
	}

	return items, nil
}

// GetPlayStats 获取播放统计
func (s *MeditationFavoriteService) GetPlayStats(userID int64) (map[string]interface{}, error) {
	stats, err := s.playRecordRepo.GetPlayStats(userID)
	if err != nil {
		return nil, err
	}

	// 添加收藏数量
	favoriteCount, _ := s.favoriteRepo.GetFavoriteCount(userID)
	stats["favorite_count"] = favoriteCount

	return stats, nil
}
