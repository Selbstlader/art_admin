package repository

import (
	"art_admin_backend/internal/model"
	"time"

	"gorm.io/gorm"
)

// ChatRoomRepository 聊天室仓储接口
type ChatRoomRepository interface {
	Create(room *model.ChatRoom) error
	Update(room *model.ChatRoom) error
	Delete(id uint) error
	FindByID(id uint) (*model.ChatRoom, error)
	FindAll(page, pageSize int, keyword, roomType string, isActive *bool) ([]*model.ChatRoom, int64, error)
	FindActiveRooms() ([]*model.ChatRoom, error)
}

type chatRoomRepository struct {
	db *gorm.DB
}

// NewChatRoomRepository 创建聊天室仓储实例
func NewChatRoomRepository(db *gorm.DB) ChatRoomRepository {
	return &chatRoomRepository{db: db}
}

func (r *chatRoomRepository) Create(room *model.ChatRoom) error {
	return r.db.Create(room).Error
}

func (r *chatRoomRepository) Update(room *model.ChatRoom) error {
	return r.db.Save(room).Error
}

func (r *chatRoomRepository) Delete(id uint) error {
	return r.db.Delete(&model.ChatRoom{}, id).Error
}

func (r *chatRoomRepository) FindByID(id uint) (*model.ChatRoom, error) {
	var room model.ChatRoom
	err := r.db.First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *chatRoomRepository) FindAll(page, pageSize int, keyword, roomType string, isActive *bool) ([]*model.ChatRoom, int64, error) {
	var rooms []*model.ChatRoom
	var total int64

	query := r.db.Model(&model.ChatRoom{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if roomType != "" {
		query = query.Where("type = ?", roomType)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rooms).Error
	if err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *chatRoomRepository) FindActiveRooms() ([]*model.ChatRoom, error) {
	var rooms []*model.ChatRoom
	err := r.db.Where("is_active = ?", true).Order("created_at DESC").Find(&rooms).Error
	return rooms, err
}

// ChatMessageRepository 聊天消息仓储接口
type ChatMessageRepository interface {
	Create(message *model.ChatMessage) error
	Update(message *model.ChatMessage) error
	Delete(id uint) error
	FindByID(id uint) (*model.ChatMessage, error)
	FindByRoomID(roomID uint, page, pageSize int, beforeID, afterID *uint) ([]*model.ChatMessage, int64, error)
	FindAllByRoomID(roomID uint) ([]*model.ChatMessage, error)
	RecallMessage(id uint) error
	GetLastMessageTime(roomID uint) (*time.Time, error)
}

type chatMessageRepository struct {
	db *gorm.DB
}

// NewChatMessageRepository 创建聊天消息仓储实例
func NewChatMessageRepository(db *gorm.DB) ChatMessageRepository {
	return &chatMessageRepository{db: db}
}

func (r *chatMessageRepository) Create(message *model.ChatMessage) error {
	return r.db.Create(message).Error
}

func (r *chatMessageRepository) Update(message *model.ChatMessage) error {
	return r.db.Save(message).Error
}

func (r *chatMessageRepository) Delete(id uint) error {
	return r.db.Delete(&model.ChatMessage{}, id).Error
}

func (r *chatMessageRepository) FindByID(id uint) (*model.ChatMessage, error) {
	var message model.ChatMessage
	err := r.db.Preload("ReplyTo").First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *chatMessageRepository) FindByRoomID(roomID uint, page, pageSize int, beforeID, afterID *uint) ([]*model.ChatMessage, int64, error) {
	var messages []*model.ChatMessage
	var total int64

	query := r.db.Model(&model.ChatMessage{}).Where("room_id = ?", roomID)

	// 如果指定了 beforeID，获取该消息之前的消息（用于加载历史消息）
	if beforeID != nil {
		query = query.Where("id < ?", *beforeID)
	}

	// 如果指定了 afterID，获取该消息之后的消息（用于加载新消息）
	if afterID != nil {
		query = query.Where("id > ?", *afterID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，预加载回复消息
	offset := (page - 1) * pageSize
	err := query.Preload("ReplyTo").Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (r *chatMessageRepository) FindAllByRoomID(roomID uint) ([]*model.ChatMessage, error) {
	var messages []*model.ChatMessage

	err := r.db.Model(&model.ChatMessage{}).
		Where("room_id = ?", roomID).
		Preload("ReplyTo").
		Order("created_at ASC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *chatMessageRepository) RecallMessage(id uint) error {
	return r.db.Model(&model.ChatMessage{}).Where("id = ?", id).Update("is_recalled", true).Error
}

func (r *chatMessageRepository) GetLastMessageTime(roomID uint) (*time.Time, error) {
	var message model.ChatMessage
	err := r.db.Where("room_id = ?", roomID).Order("created_at DESC").First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &message.CreatedAt, nil
}

// ChatRoomMemberRepository 聊天室成员仓储接口
type ChatRoomMemberRepository interface {
	Create(member *model.ChatRoomMember) error
	Update(member *model.ChatRoomMember) error
	Delete(id uint) error
	DeleteByRoomAndUser(roomID, userID uint) error
	FindByID(id uint) (*model.ChatRoomMember, error)
	FindByRoomAndUser(roomID, userID uint) (*model.ChatRoomMember, error)
	FindByRoomID(roomID uint) ([]*model.ChatRoomMember, error)
	FindByUserID(userID uint) ([]*model.ChatRoomMember, error)
	CountByRoomID(roomID uint) (int64, error)
	UpdateLastReadTime(roomID, userID uint, readTime time.Time) error
}

type chatRoomMemberRepository struct {
	db *gorm.DB
}

// NewChatRoomMemberRepository 创建聊天室成员仓储实例
func NewChatRoomMemberRepository(db *gorm.DB) ChatRoomMemberRepository {
	return &chatRoomMemberRepository{db: db}
}

func (r *chatRoomMemberRepository) Create(member *model.ChatRoomMember) error {
	return r.db.Create(member).Error
}

func (r *chatRoomMemberRepository) Update(member *model.ChatRoomMember) error {
	return r.db.Save(member).Error
}

func (r *chatRoomMemberRepository) Delete(id uint) error {
	return r.db.Delete(&model.ChatRoomMember{}, id).Error
}

func (r *chatRoomMemberRepository) DeleteByRoomAndUser(roomID, userID uint) error {
	return r.db.Where("room_id = ? AND user_id = ?", roomID, userID).Delete(&model.ChatRoomMember{}).Error
}

func (r *chatRoomMemberRepository) FindByID(id uint) (*model.ChatRoomMember, error) {
	var member model.ChatRoomMember
	err := r.db.First(&member, id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *chatRoomMemberRepository) FindByRoomAndUser(roomID, userID uint) (*model.ChatRoomMember, error) {
	var member model.ChatRoomMember
	err := r.db.Where("room_id = ? AND user_id = ?", roomID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *chatRoomMemberRepository) FindByRoomID(roomID uint) ([]*model.ChatRoomMember, error) {
	var members []*model.ChatRoomMember
	err := r.db.Where("room_id = ?", roomID).Order("joined_at ASC").Find(&members).Error
	return members, err
}

func (r *chatRoomMemberRepository) FindByUserID(userID uint) ([]*model.ChatRoomMember, error) {
	var members []*model.ChatRoomMember
	err := r.db.Preload("Room").Where("user_id = ?", userID).Order("joined_at DESC").Find(&members).Error
	return members, err
}

func (r *chatRoomMemberRepository) CountByRoomID(roomID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.ChatRoomMember{}).Where("room_id = ?", roomID).Count(&count).Error
	return count, err
}

func (r *chatRoomMemberRepository) UpdateLastReadTime(roomID, userID uint, readTime time.Time) error {
	return r.db.Model(&model.ChatRoomMember{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("last_read_at", readTime).Error
}
