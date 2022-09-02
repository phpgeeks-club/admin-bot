package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("GEEKSONATOR_TELEGRAM_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 // long polling

	updates := bot.GetUpdatesChan(u)
	commands := [22]string{"/help", "/хелп", "/php", "/пхп", "/jun", "/джун", "/go", "/го", "/db", "/дб", "/lara", "/лара", "/js", "/жс", "/hr", "/хр", "/fl", "/фл", "/code", "/код", "/nometa", "/номета"}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		authorIsAdmin := false
		message := ""

		admins, _ := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{ChatConfig: update.Message.Chat.ChatConfig()})
		if err != nil {
			log.Printf("GetChatAdministrators error: %v", err)

			continue
		}

		for _, admin := range admins {
			if admin.User.ID == update.Message.From.ID {
				authorIsAdmin = true
			}
		}

		if !authorIsAdmin {
			continue
		}

		for _, command := range commands {
			if update.Message.Text == command || strings.Contains(update.Message.Text, " "+command+" ") {
				switch command {
				case "/help", "/хелп":
					message = `БОТ РАБОТАЕТ ТОЛЬКО У АДМИНОВ.

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
[<code>/code</code>, <code>/код</code>] Код в нашем чате ложут на pastebin.org, gist.github.com или любой аналогичный ресурс (с)der_Igel
[<code>/nometa</code>, <code>/номета</code>] nometa.xyz`
				case "/php", "/пхп":
					message = "@phpGeeks - Best PHP chat"
				case "/jun", "/джун":
					message = "@phpGeeksJunior - Группа для новичков. Не стесняйтесь задавать вопросы по php."
				case "/go", "/го":
					message = "@golangGeeks - Приветствуем всех в нашем гетеросексуальном чате гоферов!"
				case "/db", "/бд":
					message = "@dbGeeks - Чат про базы данных, их устройство и приемы работы с ними."
				case "/lara", "/лара":
					message = "@laravel_pro - Официальный чат для всех Laravel программистов."
				case "/js", "/жс":
					message = "@jsChat - Чат посвященный программированию на языке JavaScript."
				case "/hr", "/хр":
					message = "@jobGeeks - Топ вакансии (250 000+ р/мес)."
				case "/fl", "/фл":
					message = "@freelanceGeeks - IT фриланс, ищем исполнителей и заказчиков, делимся опытом и проблемами связанными с фрилансом."
				case "/code", "/код":
					message = "Код в нашем чате ложут на pastebin.org, gist.github.com или любой аналогичный ресурс (с)der_Igel"
				case "/nometa", "/номета":
					message = "nometa.xyz"
				}
			}
		}

		if message == "" {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
		msg.ParseMode = "html"
		msg.DisableWebPagePreview = true

		if update.Message.ReplyToMessage != nil {
			msg.ReplyToMessageID = update.Message.ReplyToMessage.MessageID
		}

		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("[%s] Send message error: %v", msg.Text, err)
		}
	}
}
