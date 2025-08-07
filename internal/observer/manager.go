package observer

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Manager is manager for observer.
type Manager struct {
	bot            BotProvider
	chanUpdates    tgbotapi.UpdatesChannel
	cache          Cache
	logger         *zap.Logger
	skipAdminCheck bool
}

// NewManager creates new manager.
func NewManager(bot BotProvider, chanUpdates tgbotapi.UpdatesChannel, cache Cache, opts ...ManagerOption) *Manager {
	m := &Manager{
		bot:         bot,
		chanUpdates: chanUpdates,
		cache:       cache,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// ManagerOption is functional option.
type ManagerOption func(m *Manager)

// WithDebug enables debug logging.
func WithDebug(logger *zap.Logger) ManagerOption {
	return func(m *Manager) {
		m.logger = logger.Named("observer_manager")
	}
}

// WithSkipAdminCheck skips admin check.
func WithSkipAdminCheck() ManagerOption {
	return func(m *Manager) {
		m.skipAdminCheck = true
	}
}

// Run runs manager.
func (m *Manager) Run(ctx context.Context) error {
	for update := range m.chanUpdates {
		select {
		case <-ctx.Done():
			m.log("Gracefully stopped")

			return nil
		default:
			err := m.processingUpdate(update)
			if err != nil {
				return fmt.Errorf("m.processingUpdate: %v", err)
			}
		}
	}

	return nil
}

// processingUpdate processes update.
func (m *Manager) processingUpdate(update tgbotapi.Update) error {
	msgText, err := m.processingMessage(update.Message)
	if err != nil {
		return fmt.Errorf("m.processingMessage: %v", err)
	}
	if msgText == "" {
		return nil
	}

	if err := m.sendMessage(update.Message, msgText); err != nil {
		return fmt.Errorf("m.sendMessage: %v", err)
	}

	return nil
}

// processingMessage processes message.
func (m *Manager) processingMessage(message *tgbotapi.Message) (string, error) {
	if message == nil {
		return "", nil
	}
	m.log("Received message",
		zap.String("message", message.Text),
	)

	msgText := getMessageText(message.Text)
	if msgText == "" {
		return "", nil
	}
	m.log("Output message",
		zap.String("msgText", msgText),
	)

	if m.skipAdminCheck {
		return msgText, nil
	}

	admins, err := m.getAdmins(message.Chat.ChatConfig())
	if err != nil {
		return "", fmt.Errorf("m.getAdmins: %v", err)
	}

	if authorIsAdmin(admins, message.From.ID) {
		return msgText, nil
	}

	return "", nil
}

// sendMessage sends message.
func (m *Manager) sendMessage(updateMsg *tgbotapi.Message, message string) error { //nolint:gocognit // it's ok
	msg := m.bot.NewMessage(updateMsg.Chat.ID, message)
	msg.ParseMode = "html"
	msg.DisableWebPagePreview = true

	if updateMsg.ReplyToMessage != nil { //nolint:nestif // it's ok
		msg.ReplyToMessageID = updateMsg.ReplyToMessage.MessageID

		if updateMsg.ReplyToMessage.From != nil {
			var fromUser string
			if updateMsg.ReplyToMessage.From.UserName != "" {
				fromUser = "@" + updateMsg.ReplyToMessage.From.UserName
			} else {
				uName := updateMsg.ReplyToMessage.From.FirstName
				if updateMsg.ReplyToMessage.From.LastName != "" {
					uName += " " + updateMsg.ReplyToMessage.From.LastName
				}

				fromUser = fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`,
					updateMsg.ReplyToMessage.From.ID,
					uName,
				)
			}

			msg.Text = fromUser + " " + msg.Text
		}
	}

	_, err := m.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("m.bot.Send(%s): %v", msg.Text, err)
	}

	return nil
}

// log debug message.
func (m *Manager) log(msg string, fields ...zapcore.Field) {
	if m.logger != nil {
		m.logger.Debug(msg, fields...)
	}
}

// getAdmins returns admins.
func (m *Manager) getAdmins(chatCfg tgbotapi.ChatConfig) ([]tgbotapi.ChatMember, error) {
	adminsFromCache, ok := m.cache.Get(chatCfg.ChatID)
	if ok {
		return adminsFromCache, nil
	}

	admins, err := m.bot.GetChatAdministrators(chatCfg)
	if err != nil {
		return nil, fmt.Errorf("m.bot.GetChatAdministrators: %v", err)
	}

	if err := m.cache.Set(chatCfg.ChatID, admins); err != nil {
		m.log("Set admins in cache",
			zap.Error(err),
		)
	}

	return admins, nil
}

// authorIsAdmin returns true if author is admin.
func authorIsAdmin(admins []tgbotapi.ChatMember, userID int64) bool {
	for _, admin := range admins {
		if admin.User != nil && admin.User.ID == userID {
			return true
		}
	}

	return false
}

// getMessageText returns message text.
func getMessageText(text string) string {
	switch text {
	case "/help", "/хелп":
		return `БОТ РАБОТАЕТ ТОЛЬКО У АДМИНОВ.

Команды можно писать обычным сообщением и ответом на сообщение.

Список доступных команд:
[<code>/help</code>, <code>/хелп</code>] Список доступных команд бота
[<code>/php</code>, <code>/пхп</code>] @phpGeeks - Best PHP chat
[<code>/jun</code>, <code>/джун</code>] @phpGeeksJunior - Группа для новичков. Не стесняйтесь задавать вопросы по php.
[<code>/go</code>, <code>/го</code>] @golangGeeks - Приветствуем всех в нашем гетеросексуальном чате гоферов!
[<code>/db</code>, <code>/дб</code>] @dbGeeks - Чат про базы данных, их устройство и приемы работы с ними.
[<code>/lara</code>, <code>/лара</code>] @laravel_pro - Официальный чат для всех Laravel программистов.
[<code>/js</code>, <code>/жс</code>] @jsChat - Чат посвященный программированию на языке JavaScript.
[<code>/hr</code>, <code>/хр</code>] @jobGeeks - Топ вакансии (250 000+ р/мес).
[<code>/fl</code>, <code>/фл</code>] @freelanceGeeks - IT фриланс, ищем исполнителей и заказчиков, делимся опытом и проблемами связанными с фрилансом.
[<code>/job</code>, <code>/раб</code>] Объединяет сразу две команды: <code>/hr</code> и <code>/fl</code>.
[<code>/code</code>, <code>/код</code>] Код в нашем чате <a href="https://t.me/phpGeeks/1318040">ложут</a> на pastebin.org, gist.github.com или любой аналогичный ресурс (с)der_Igel
[<code>/nometa</code>, <code>/номета</code>] nometa.xyz
[<code>/wtf</code>, <code>/втф</code>] А причём тут пхп?`
	case "/php", "/пхп":
		return "@phpGeeks - Best PHP chat"
	case "/jun", "/джун":
		return "@phpGeeksJunior - Группа для новичков. Не стесняйтесь задавать вопросы по php."
	case "/go", "/го":
		return "@golangGeeks - Приветствуем всех в нашем гетеросексуальном чате гоферов!"
	case "/db", "/бд":
		return "@dbGeeks - Чат про базы данных, их устройство и приемы работы с ними."
	case "/lara", "/лара":
		return "@laravel_pro - Официальный чат для всех Laravel программистов."
	case "/js", "/жс":
		return "@jsChat - Чат посвященный программированию на языке JavaScript."
	case "/hr", "/хр":
		return "@jobGeeks - Топ вакансии (250 000+ р/мес)."
	case "/fl", "/фл":
		return "@freelanceGeeks - IT фриланс, ищем исполнителей и заказчиков, делимся опытом и проблемами связанными с фрилансом."
	case "/job", "/раб":
		return `@jobGeeks - Топ вакансии (250 000+ р/мес).
@freelanceGeeks - IT фриланс, ищем исполнителей и заказчиков, делимся опытом и проблемами связанными с фрилансом.`
	case "/code", "/код":
		return "Код в нашем чате <a href=\"https://t.me/phpGeeks/1318040\">ложут</a> на pastebin.org, gist.github.com или любой аналогичный ресурс (с)der_Igel"
	case "/nometa", "/номета":
		return "nometa.xyz"
	case "/wtf", "/втф":
		return "А причём тут пхп?"
	}

	return ""
}
