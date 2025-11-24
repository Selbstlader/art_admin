package service

import (
	"art_admin_backend/internal/dto"
	"art_admin_backend/internal/model"
	"art_admin_backend/internal/repository"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ChatService 聊天服务接口
type ChatService interface {
	// 聊天室管理
	CreateRoom(req *dto.CreateChatRoomRequest, userID uint) (*dto.ChatRoomResponse, error)
	UpdateRoom(req *dto.UpdateChatRoomRequest) error
	DeleteRoom(id uint) error
	GetRoomByID(id uint) (*dto.ChatRoomResponse, error)
	GetRoomList(req *dto.ChatRoomListRequest) ([]*dto.ChatRoomResponse, int64, error)

	// 成员管理
	JoinRoom(roomID, userID uint, username string) error
	LeaveRoom(roomID, userID uint) error
	GetRoomMembers(roomID uint) ([]*dto.ChatRoomMemberResponse, error)
	GetUserRooms(userID uint) ([]*dto.ChatRoomResponse, error)

	// 消息管理
	SendMessage(req *dto.SendMessageRequest, userID uint, username string) (*dto.ChatMessageResponse, error)
	GetMessages(req *dto.MessageListRequest) ([]*dto.ChatMessageResponse, int64, error)
	GetAllMessages(roomID uint) ([]*dto.ChatMessageResponse, error)
	RecallMessage(messageID, userID uint) error
	UpdateLastReadTime(roomID, userID uint) error
}

type chatService struct {
	roomRepo   repository.ChatRoomRepository
	msgRepo    repository.ChatMessageRepository
	memberRepo repository.ChatRoomMemberRepository
}

// NewChatService 创建聊天服务实例
func NewChatService(
	roomRepo repository.ChatRoomRepository,
	msgRepo repository.ChatMessageRepository,
	memberRepo repository.ChatRoomMemberRepository,
) ChatService {
	return &chatService{
		roomRepo:   roomRepo,
		msgRepo:    msgRepo,
		memberRepo: memberRepo,
	}
}

func (s *chatService) CreateRoom(req *dto.CreateChatRoomRequest, userID uint) (*dto.ChatRoomResponse, error) {
	room := &model.ChatRoom{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		MaxMembers:  req.MaxMembers,
		IsActive:    true,
		CreatedBy:   userID,
	}

	if err := s.roomRepo.Create(room); err != nil {
		return nil, err
	}

	// 创建者自动加入聊天室
	member := &model.ChatRoomMember{
		RoomID:   room.ID,
		UserID:   userID,
		Role:     "owner",
		JoinedAt: time.Now(),
	}
	if err := s.memberRepo.Create(member); err != nil {
		return nil, err
	}

	return s.GetRoomByID(room.ID)
}

func (s *chatService) UpdateRoom(req *dto.UpdateChatRoomRequest) error {
	room, err := s.roomRepo.FindByID(req.ID)
	if err != nil {
		return err
	}

	room.Name = req.Name
	room.Description = req.Description
	room.MaxMembers = req.MaxMembers
	if req.IsActive != nil {
		room.IsActive = *req.IsActive
	}

	return s.roomRepo.Update(room)
}

func (s *chatService) DeleteRoom(id uint) error {
	return s.roomRepo.Delete(id)
}

func (s *chatService) GetRoomByID(id uint) (*dto.ChatRoomResponse, error) {
	room, err := s.roomRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 获取成员数量
	memberCount, _ := s.memberRepo.CountByRoomID(id)

	// 获取最后一条消息时间
	lastMessageAt, _ := s.msgRepo.GetLastMessageTime(id)

	return &dto.ChatRoomResponse{
		ID:            room.ID,
		Name:          room.Name,
		Description:   room.Description,
		Type:          room.Type,
		MaxMembers:    room.MaxMembers,
		IsActive:      room.IsActive,
		CreatedBy:     room.CreatedBy,
		MemberCount:   int(memberCount),
		OnlineCount:   0, // 将由 WebSocket 管理器提供
		LastMessageAt: lastMessageAt,
		CreatedAt:     room.CreatedAt,
		UpdatedAt:     room.UpdatedAt,
	}, nil
}

func (s *chatService) GetRoomList(req *dto.ChatRoomListRequest) ([]*dto.ChatRoomResponse, int64, error) {
	rooms, total, err := s.roomRepo.FindAll(req.Page, req.PageSize, req.Keyword, req.Type, req.IsActive)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.ChatRoomResponse
	for _, room := range rooms {
		memberCount, _ := s.memberRepo.CountByRoomID(room.ID)
		lastMessageAt, _ := s.msgRepo.GetLastMessageTime(room.ID)

		responses = append(responses, &dto.ChatRoomResponse{
			ID:            room.ID,
			Name:          room.Name,
			Description:   room.Description,
			Type:          room.Type,
			MaxMembers:    room.MaxMembers,
			IsActive:      room.IsActive,
			CreatedBy:     room.CreatedBy,
			MemberCount:   int(memberCount),
			OnlineCount:   0,
			LastMessageAt: lastMessageAt,
			CreatedAt:     room.CreatedAt,
			UpdatedAt:     room.UpdatedAt,
		})
	}

	return responses, total, nil
}

func (s *chatService) JoinRoom(roomID, userID uint, username string) error {
	// 检查聊天室是否存在
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return err
	}

	if !room.IsActive {
		return errors.New("聊天室已禁用")
	}

	// 检查是否已经是成员
	_, err = s.memberRepo.FindByRoomAndUser(roomID, userID)
	if err == nil {
		// 已经是成员，直接返回成功
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// 检查人数限制
	if room.MaxMembers > 0 {
		count, err := s.memberRepo.CountByRoomID(roomID)
		if err != nil {
			return err
		}
		if count >= int64(room.MaxMembers) {
			return errors.New("聊天室已满")
		}
	}

	// 加入聊天室
	member := &model.ChatRoomMember{
		RoomID:   roomID,
		UserID:   userID,
		Username: username,
		Role:     "member",
		JoinedAt: time.Now(),
	}

	return s.memberRepo.Create(member)
}

func (s *chatService) LeaveRoom(roomID, userID uint) error {
	return s.memberRepo.DeleteByRoomAndUser(roomID, userID)
}

func (s *chatService) GetRoomMembers(roomID uint) ([]*dto.ChatRoomMemberResponse, error) {
	members, err := s.memberRepo.FindByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	var responses []*dto.ChatRoomMemberResponse
	for _, member := range members {
		responses = append(responses, &dto.ChatRoomMemberResponse{
			ID:         member.ID,
			RoomID:     member.RoomID,
			UserID:     member.UserID,
			Username:   member.Username,
			Role:       member.Role,
			IsMuted:    member.IsMuted,
			IsOnline:   false, // 将由 WebSocket 管理器提供
			LastReadAt: member.LastReadAt,
			JoinedAt:   member.JoinedAt,
		})
	}

	return responses, nil
}

func (s *chatService) GetUserRooms(userID uint) ([]*dto.ChatRoomResponse, error) {
	members, err := s.memberRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []*dto.ChatRoomResponse
	for _, member := range members {
		if member.Room != nil {
			memberCount, _ := s.memberRepo.CountByRoomID(member.Room.ID)
			lastMessageAt, _ := s.msgRepo.GetLastMessageTime(member.Room.ID)

			responses = append(responses, &dto.ChatRoomResponse{
				ID:            member.Room.ID,
				Name:          member.Room.Name,
				Description:   member.Room.Description,
				Type:          member.Room.Type,
				MaxMembers:    member.Room.MaxMembers,
				IsActive:      member.Room.IsActive,
				CreatedBy:     member.Room.CreatedBy,
				MemberCount:   int(memberCount),
				OnlineCount:   0,
				LastMessageAt: lastMessageAt,
				CreatedAt:     member.Room.CreatedAt,
				UpdatedAt:     member.Room.UpdatedAt,
			})
		}
	}

	return responses, nil
}

func (s *chatService) SendMessage(req *dto.SendMessageRequest, userID uint, username string) (*dto.ChatMessageResponse, error) {
	// 检查是否是聊天室成员
	member, err := s.memberRepo.FindByRoomAndUser(req.RoomID, userID)
	if err != nil {
		return nil, errors.New("不是聊天室成员")
	}

	// 检查是否被禁言
	if member.IsMuted {
		return nil, errors.New("已被禁言")
	}

	// 创建消息
	message := &model.ChatMessage{
		RoomID:      req.RoomID,
		UserID:      userID,
		Username:    username,
		Content:     req.Content,
		MessageType: req.MessageType,
		ReplyToID:   req.ReplyToID,
	}

	if err := s.msgRepo.Create(message); err != nil {
		return nil, err
	}

	// 如果有回复消息，加载回复消息信息
	if req.ReplyToID != nil {
		replyMsg, _ := s.msgRepo.FindByID(*req.ReplyToID)
		if replyMsg != nil {
			message.ReplyTo = replyMsg
		}
	}

	return s.convertMessageToResponse(message), nil
}

func (s *chatService) GetMessages(req *dto.MessageListRequest) ([]*dto.ChatMessageResponse, int64, error) {
	messages, total, err := s.msgRepo.FindByRoomID(req.RoomID, req.Page, req.PageSize, req.BeforeID, req.AfterID)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.ChatMessageResponse
	for _, msg := range messages {
		responses = append(responses, s.convertMessageToResponse(msg))
	}

	return responses, total, nil
}

func (s *chatService) GetAllMessages(roomID uint) ([]*dto.ChatMessageResponse, error) {
	messages, err := s.msgRepo.FindAllByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	var responses []*dto.ChatMessageResponse
	for _, msg := range messages {
		responses = append(responses, s.convertMessageToResponse(msg))
	}

	return responses, nil
}

func (s *chatService) RecallMessage(messageID, userID uint) error {
	message, err := s.msgRepo.FindByID(messageID)
	if err != nil {
		return err
	}

	// 只能撤回自己的消息
	if message.UserID != userID {
		return errors.New("只能撤回自己的消息")
	}

	// 检查消息发送时间（例如：只能撤回2分钟内的消息）
	if time.Since(message.CreatedAt) > 2*time.Minute {
		return errors.New("消息发送时间超过2分钟，无法撤回")
	}

	return s.msgRepo.RecallMessage(messageID)
}

func (s *chatService) UpdateLastReadTime(roomID, userID uint) error {
	return s.memberRepo.UpdateLastReadTime(roomID, userID, time.Now())
}

func (s *chatService) convertMessageToResponse(msg *model.ChatMessage) *dto.ChatMessageResponse {
	resp := &dto.ChatMessageResponse{
		ID:          msg.ID,
		RoomID:      msg.RoomID,
		UserID:      msg.UserID,
		Username:    msg.Username,
		Content:     msg.Content,
		MessageType: msg.MessageType,
		ReplyToID:   msg.ReplyToID,
		IsRecalled:  msg.IsRecalled,
		CreatedAt:   msg.CreatedAt,
	}

	if msg.ReplyTo != nil {
		resp.ReplyTo = &dto.ChatMessageResponse{
			ID:          msg.ReplyTo.ID,
			RoomID:      msg.ReplyTo.RoomID,
			UserID:      msg.ReplyTo.UserID,
			Username:    msg.ReplyTo.Username,
			Content:     msg.ReplyTo.Content,
			MessageType: msg.ReplyTo.MessageType,
			IsRecalled:  msg.ReplyTo.IsRecalled,
			CreatedAt:   msg.ReplyTo.CreatedAt,
		}
	}

	return resp
}
