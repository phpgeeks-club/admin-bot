package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Service is service for telegram bot.
type Service struct {
	bot BotAPI
}

// NewService creates new service for telegram bot.
func NewService(bot BotAPI) *Service {
	return &Service{
		bot: bot,
	}
}

// GetChatAdministrators returns list of administrators.
func (s *Service) GetChatAdministrators(chatConfig tgbotapi.ChatConfig) ([]tgbotapi.ChatMember, error) {
	admins, err := s.bot.GetChatAdministrators(
		tgbotapi.ChatAdministratorsConfig{
			ChatConfig: chatConfig,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("s.bot.GetChatAdministrators: %v", err)
	}

	return admins, nil
}

// NewMessage creates new message.
func (*Service) NewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, text)
}

// Send sends message.
func (s *Service) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	msg, err := s.bot.Send(c)
	if err != nil {
		return tgbotapi.Message{}, fmt.Errorf("s.bot.Send: %v", err)
	}

	return msg, nil
}
