package telegram

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"

	"geeksonator/internal/provider/telegram/mocks"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	type args struct {
		bot BotAPI
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Success",
			args: args{
				bot: nil,
			},
			want: &Service{
				bot: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewService(tt.args.bot)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_GetChatAdministrators(t *testing.T) {
	t.Parallel()

	type args struct {
		chatConfig tgbotapi.ChatConfig
	}
	tests := []struct {
		name    string
		srv     func() *Service
		args    args
		want    []tgbotapi.ChatMember
		wantErr error
	}{
		{
			name: "Success",
			srv: func() *Service {
				bot := mocks.NewBotAPIMock(t)

				bot.EXPECT().
					GetChatAdministrators(
						tgbotapi.ChatAdministratorsConfig{
							ChatConfig: tgbotapi.ChatConfig{},
						},
					).
					Return([]tgbotapi.ChatMember{}, nil)

				return &Service{
					bot: bot,
				}
			},
			args: args{
				chatConfig: tgbotapi.ChatConfig{},
			},
			want:    []tgbotapi.ChatMember{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.srv().GetChatAdministrators(tt.args.chatConfig)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_NewMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		chatID int64
		text   string
	}
	tests := []struct {
		name string
		srv  *Service
		args args
		want tgbotapi.MessageConfig
	}{
		{
			name: "Success",
			srv: &Service{
				bot: nil,
			},
			args: args{
				chatID: 100500,
				text:   "message text",
			},
			want: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: 100500,
				},
				Text: "message text",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.srv.NewMessage(tt.args.chatID, tt.args.text)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_Send(t *testing.T) {
	t.Parallel()

	type args struct {
		c tgbotapi.Chattable
	}
	tests := []struct {
		name    string
		srv     func() *Service
		args    args
		want    tgbotapi.Message
		wantErr error
	}{
		{
			name: "Success",
			srv: func() *Service {
				bot := mocks.NewBotAPIMock(t)

				bot.EXPECT().
					Send(
						tgbotapi.MessageConfig{
							BaseChat: tgbotapi.BaseChat{
								ChatID: 100500,
							},
							Text: "message text",
						},
					).
					Return(tgbotapi.Message{}, nil)

				return &Service{
					bot: bot,
				}
			},
			args: args{
				c: tgbotapi.NewMessage(100500, "message text"),
			},
			want:    tgbotapi.Message{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.srv().Send(tt.args.c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
