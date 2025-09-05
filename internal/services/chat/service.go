package chat

import (
	"context"
	"database/sql"

	"github.com/I-Van-Radkov/messenger/internal/models"
	"github.com/I-Van-Radkov/messenger/internal/repository"
)

type MessageProvider interface {
	GetMessagesByDialogID(ctx context.Context, dialogID int64, limit, offset int) ([]*models.Message, error)
	GetLastMessagesByDialogID(ctx context.Context, dialogIDs []int64) (map[int64]*models.Message, error)
}

type ChatService struct {
	chatRepo       repository.ChatRepository
	messageService MessageProvider
}

func NewChatService(chatRepo repository.ChatRepository, messageService MessageProvider) *ChatService {
	return &ChatService{
		chatRepo:       chatRepo,
		messageService: messageService,
	}
}

func (s *ChatService) GetChats(ctx context.Context, userID int64, limit, offset int) ([]*models.Chat, map[int64]*models.Message, error) {
	chats, err := s.chatRepo.GetChatsByUserID(ctx, userID, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Chat{}, nil, nil
		}
		return nil, nil, err
	}

	if len(chats) == 0 {
		return chats, map[int64]*models.Message{}, nil
	}

	dialogIDs := make([]int64, len(chats))
	for i, chat := range chats {
		dialogIDs[i] = chat.ID
	}

	lastMsgInChats, err := s.messageService.GetLastMessagesByDialogID(ctx, dialogIDs)
	if err != nil {
		return nil, nil, err
	}

	return chats, lastMsgInChats, nil
}

func (s *ChatService) GetUserChat(ctx context.Context, dialogID int64, limit, offset int) ([]*models.Message, error) {
	messages, err := s.messageService.GetMessagesByDialogID(ctx, dialogID, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Message{}, nil
		}
		return nil, err
	}

	return messages, err
}
