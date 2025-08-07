package observer

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotProvider interface for telegram bot.
type BotProvider interface {
	// GetChatAdministrators returns list of administrators.
	GetChatAdministrators(chatConfig tgbotapi.ChatConfig) ([]tgbotapi.ChatMember, error)

	// NewMessage creates new message.
	NewMessage(chatID int64, text string) tgbotapi.MessageConfig

	// Send sends message.
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

// Cache interface for cache.
type Cache interface {
	// Get looks up a key's value from the cache.
	Get(key int64) (value []tgbotapi.ChatMember, ok bool)

	// Set adds a value to the cache.
	Set(key int64, value []tgbotapi.ChatMember) error
}
