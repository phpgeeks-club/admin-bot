package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotAPI interface for telegram bot.
type BotAPI interface {
	// GetChatAdministrators returns list of administrators.
	GetChatAdministrators(config tgbotapi.ChatAdministratorsConfig) ([]tgbotapi.ChatMember, error)

	// NewMessage creates new message.
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}
