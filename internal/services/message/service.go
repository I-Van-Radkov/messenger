package message

import (
	"context"

	"github.com/I-Van-Radkov/messenger/internal/models"
	"github.com/I-Van-Radkov/messenger/internal/repository"
)

type MessageService struct {
	messageRepo repository.MessageRepository
}

func NewMessageService(messageRepo repository.MessageRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

func (s *MessageService) Create(ctx context.Context, message *models.Message) (int64, error) {
	return s.messageRepo.Create(ctx, message)
}

func (s *MessageService) UpdateStatus(ctx context.Context, msgId int64, status string) {
	s.messageRepo.UpdateStatus(ctx, msgId, status)
}

func (s *MessageService) GetMessagesByDialogID(ctx context.Context, dialogID int64, limit, offset int) ([]*models.Message, error) {
	return s.messageRepo.GetMessagesByDialogID(ctx, dialogID, limit, offset)
}
func (s *MessageService) GetLastMessagesByDialogID(ctx context.Context, dialogIDs []int64) (map[int64]*models.Message, error) {
	return s.messageRepo.GetLastMessagesByDialogID(ctx, dialogIDs)
}
