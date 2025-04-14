package observer

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"geeksonator/internal/observer/mocks"
)

const (
	laraCmd = "/lara"
	laraTxt = "@laravel_pro - Официальный чат для всех Laravel программистов."
)

func TestNewManager(t *testing.T) {
	t.Parallel()

	type args struct {
		bot         BotProvider
		chanUpdates tgbotapi.UpdatesChannel
		cache       Cache
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		{
			name: "Success",
			args: args{
				bot:         nil,
				chanUpdates: nil,
				cache:       nil,
			},
			want: &Manager{
				bot:         nil,
				chanUpdates: nil,
				cache:       nil,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewManager(tt.args.bot, tt.args.chanUpdates, tt.args.cache)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewManagerWithDebug(t *testing.T) {
	t.Parallel()

	m := NewManager(nil, nil, nil, WithDebug(zap.NewNop()))
	assert.NotEqual(t, nil, m.logger)
}

func TestNewManagerWithSkipAdminCheck(t *testing.T) {
	t.Parallel()

	m := NewManager(nil, nil, nil, WithSkipAdminCheck())
	assert.True(t, m.skipAdminCheck)
}

func TestManager_processingUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		update tgbotapi.Update
	}
	tests := []struct {
		name string
		man  func() *Manager
		args args
		want error
	}{
		{
			name: "Success",
			man: func() *Manager {
				botProvider := mocks.NewBotProviderMock(t)

				botProvider.EXPECT().
					NewMessage(int64(300600), "nometa.xyz").
					Return(tgbotapi.MessageConfig{
						BaseChat: tgbotapi.BaseChat{
							ChatID: 300600,
						},
						Text: "nometa.xyz",
					})

				botProvider.EXPECT().
					Send(
						tgbotapi.MessageConfig{
							BaseChat: tgbotapi.BaseChat{
								ChatID: 300600,
							},
							Text:                  "nometa.xyz",
							ParseMode:             "html",
							DisableWebPagePreview: true,
						},
					).
					Return(tgbotapi.Message{}, nil)

				cache := mocks.NewCacheMock(t)

				cache.EXPECT().
					Get(cacheKey).
					Return(
						[]tgbotapi.ChatMember{
							{
								User: &tgbotapi.User{
									ID: 100500,
								},
							},
						},
						true,
					)

				return &Manager{
					bot:   botProvider,
					cache: cache,
				}
			},
			args: args{
				update: tgbotapi.Update{
					Message: &tgbotapi.Message{
						Chat: &tgbotapi.Chat{
							ID: 300600,
						},
						From: &tgbotapi.User{
							ID: 100500,
						},
						Text: "/nometa",
					},
				},
			},
			want: nil,
		},
		{
			name: "Message is nil",
			man: func() *Manager {
				return &Manager{}
			},
			args: args{
				update: tgbotapi.Update{
					Message: nil,
				},
			},
			want: nil,
		},
		{
			name: "Message is empty",
			man: func() *Manager {
				return &Manager{}
			},
			args: args{
				update: tgbotapi.Update{
					Message: &tgbotapi.Message{
						Text: "",
					},
				},
			},
			want: nil,
		},
		{
			name: "Message is unknown command",
			man: func() *Manager {
				return &Manager{}
			},
			args: args{
				update: tgbotapi.Update{
					Message: &tgbotapi.Message{
						Text: "/unknown",
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.man().processingUpdate(tt.args.update)
			assert.Equal(t, tt.want, err)
		})
	}
}

func TestManager_processingMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		message *tgbotapi.Message
	}
	tests := []struct {
		name    string
		man     func() *Manager
		args    args
		wantMsg string
		wantErr error
	}{
		{
			name: "Message is nil",
			man: func() *Manager {
				return &Manager{}
			},
			args: args{
				message: nil,
			},
			wantMsg: "",
			wantErr: nil,
		},
		{
			name: "Message is empty",
			man: func() *Manager {
				return &Manager{}
			},
			args: args{
				message: &tgbotapi.Message{
					Text: "",
				},
			},
			wantMsg: "",
			wantErr: nil,
		},
		{
			name: "Message is unknown command",
			man: func() *Manager {
				return &Manager{}
			},
			args: args{
				message: &tgbotapi.Message{
					Text: "/unknown",
				},
			},
			wantMsg: "",
			wantErr: nil,
		},
		{
			name: "Skip admin check",
			man: func() *Manager {
				return &Manager{
					skipAdminCheck: true,
				}
			},
			args: args{
				message: &tgbotapi.Message{
					Text: laraCmd,
				},
			},
			wantMsg: laraTxt,
			wantErr: nil,
		},
		{
			name: "getAdmins",
			man: func() *Manager {
				cache := mocks.NewCacheMock(t)

				cache.EXPECT().
					Get(cacheKey).
					Return([]tgbotapi.ChatMember{
						{
							User: &tgbotapi.User{
								ID: 100500,
							},
						},
					}, true)

				return &Manager{
					cache: cache,
				}
			},
			args: args{
				message: &tgbotapi.Message{
					Chat: &tgbotapi.Chat{},
					From: &tgbotapi.User{
						ID: 100500,
					},
					Text: laraCmd,
				},
			},
			wantMsg: laraTxt,
			wantErr: nil,
		},
		{
			name: "Author is admin",
			man: func() *Manager {
				cache := mocks.NewCacheMock(t)

				cache.EXPECT().
					Get(cacheKey).
					Return(
						[]tgbotapi.ChatMember{
							{
								User: &tgbotapi.User{
									ID: 100500,
								},
							},
						},
						true,
					)

				return &Manager{
					cache: cache,
				}
			},
			args: args{
				message: &tgbotapi.Message{
					Chat: &tgbotapi.Chat{},
					From: &tgbotapi.User{
						ID: 100500,
					},
					Text: laraCmd,
				},
			},
			wantMsg: laraTxt,
			wantErr: nil,
		},
		{
			name: "Author is not admin",
			man: func() *Manager {
				cache := mocks.NewCacheMock(t)

				cache.EXPECT().
					Get(cacheKey).
					Return(
						[]tgbotapi.ChatMember{
							{
								User: &tgbotapi.User{
									ID: 100500,
								},
							},
						},
						true,
					)

				return &Manager{
					cache: cache,
				}
			},
			args: args{
				message: &tgbotapi.Message{
					Chat: &tgbotapi.Chat{},
					From: &tgbotapi.User{
						ID: 100501,
					},
					Text: laraCmd,
				},
			},
			wantMsg: "",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.man().processingMessage(tt.args.message)
			assert.Equal(t, tt.wantMsg, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestManager_sendMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		updateMsg *tgbotapi.Message
		message   string
	}
	tests := []struct {
		name    string
		man     func() *Manager
		args    args
		wantErr error
	}{
		{
			name: "Success without reply to message",
			man: func() *Manager {
				botProvider := mocks.NewBotProviderMock(t)

				botProvider.EXPECT().
					NewMessage(int64(100500), "message text").
					Return(tgbotapi.MessageConfig{
						BaseChat: tgbotapi.BaseChat{
							ChatID: 100500,
						},
						Text: "message text",
					})

				botProvider.EXPECT().
					Send(
						tgbotapi.MessageConfig{
							BaseChat: tgbotapi.BaseChat{
								ChatID: 100500,
							},
							Text:                  "message text",
							ParseMode:             "html",
							DisableWebPagePreview: true,
						},
					).
					Return(tgbotapi.Message{}, nil)

				return &Manager{
					bot:         botProvider,
					chanUpdates: make(<-chan tgbotapi.Update, 100),
				}
			},
			args: args{
				updateMsg: &tgbotapi.Message{
					Chat: &tgbotapi.Chat{
						ID: 100500,
					},
				},
				message: "message text",
			},
			wantErr: nil,
		},
		{
			name: "Success with reply to message and user name",
			man: func() *Manager {
				botProvider := mocks.NewBotProviderMock(t)

				botProvider.EXPECT().
					NewMessage(int64(100500), "message text").
					Return(tgbotapi.MessageConfig{
						BaseChat: tgbotapi.BaseChat{
							ChatID: 100500,
						},
						Text: "message text",
					})

				botProvider.EXPECT().
					Send(
						tgbotapi.MessageConfig{
							BaseChat: tgbotapi.BaseChat{
								ChatID:           100500,
								ReplyToMessageID: 100,
							},
							Text:                  "@username message text",
							ParseMode:             "html",
							DisableWebPagePreview: true,
						},
					).
					Return(tgbotapi.Message{}, nil)

				return &Manager{
					bot:         botProvider,
					chanUpdates: make(<-chan tgbotapi.Update, 100),
				}
			},
			args: args{
				updateMsg: &tgbotapi.Message{
					Chat: &tgbotapi.Chat{
						ID: 100500,
					},
					ReplyToMessage: &tgbotapi.Message{
						MessageID: 100,
						From: &tgbotapi.User{
							UserName: "username",
						},
						Text: "reply to message text",
					},
				},
				message: "message text",
			},
			wantErr: nil,
		},
		{
			name: "Success with reply to message and user id",
			man: func() *Manager {
				botProvider := mocks.NewBotProviderMock(t)

				botProvider.EXPECT().
					NewMessage(int64(100500), "message text").
					Return(tgbotapi.MessageConfig{
						BaseChat: tgbotapi.BaseChat{
							ChatID: 100500,
						},
						Text: "message text",
					})

				botProvider.EXPECT().
					Send(
						tgbotapi.MessageConfig{
							BaseChat: tgbotapi.BaseChat{
								ChatID:           100500,
								ReplyToMessageID: 100,
							},
							Text:                  "<a href=\"tg://user?id=300600\">first last</a> message text",
							ParseMode:             "html",
							DisableWebPagePreview: true,
						},
					).
					Return(tgbotapi.Message{}, nil)

				return &Manager{
					bot:         botProvider,
					chanUpdates: make(<-chan tgbotapi.Update, 100),
				}
			},
			args: args{
				updateMsg: &tgbotapi.Message{
					Chat: &tgbotapi.Chat{
						ID: 100500,
					},
					ReplyToMessage: &tgbotapi.Message{
						MessageID: 100,
						From: &tgbotapi.User{
							ID:        300600,
							FirstName: "first",
							LastName:  "last",
						},
						Text: "reply to message text",
					},
				},
				message: "message text",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.man().sendMessage(tt.args.updateMsg, tt.args.message)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestManager_getAdmins(t *testing.T) {
	t.Parallel()

	type args struct {
		chatCfg tgbotapi.ChatConfig
	}
	tests := []struct {
		name       string
		man        func() *Manager
		args       args
		wantResult []tgbotapi.ChatMember
		wantErr    error
	}{
		{
			name: "Success from cache",
			man: func() *Manager {
				cache := mocks.NewCacheMock(t)

				cache.EXPECT().
					Get(cacheKey).
					Return([]tgbotapi.ChatMember{
						{
							User: &tgbotapi.User{
								ID: 100500,
							},
						},
					}, true)

				return &Manager{
					cache: cache,
				}
			},
			args: args{
				chatCfg: tgbotapi.ChatConfig{},
			},
			wantResult: []tgbotapi.ChatMember{
				{
					User: &tgbotapi.User{
						ID: 100500,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "Success from GetChatAdministrators",
			man: func() *Manager {
				botProvider := mocks.NewBotProviderMock(t)

				botProvider.EXPECT().
					GetChatAdministrators(tgbotapi.ChatConfig{}).
					Return(
						[]tgbotapi.ChatMember{
							{
								User: &tgbotapi.User{
									ID: 100500,
								},
							},
						},
						nil,
					)

				cache := mocks.NewCacheMock(t)

				cache.EXPECT().
					Get(cacheKey).
					Return(nil, false)

				cache.EXPECT().
					Set(cacheKey, []tgbotapi.ChatMember{
						{
							User: &tgbotapi.User{
								ID: 100500,
							},
						},
					}).
					Return(nil)

				return &Manager{
					bot:   botProvider,
					cache: cache,
				}
			},
			args: args{
				chatCfg: tgbotapi.ChatConfig{},
			},
			wantResult: []tgbotapi.ChatMember{
				{
					User: &tgbotapi.User{
						ID: 100500,
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.man().getAdmins(tt.args.chatCfg)
			assert.Equal(t, tt.wantResult, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_authorIsAdmin(t *testing.T) {
	t.Parallel()

	type args struct {
		admins []tgbotapi.ChatMember
		userID int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Author is admin",
			args: args{
				admins: []tgbotapi.ChatMember{
					{
						User: &tgbotapi.User{
							ID: 100500,
						},
					},
				},
				userID: 100500,
			},
			want: true,
		},
		{
			name: "Author is not admin",
			args: args{
				admins: []tgbotapi.ChatMember{
					{
						User: &tgbotapi.User{
							ID: 100500,
						},
					},
				},
				userID: 100501,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := authorIsAdmin(tt.args.admins, tt.args.userID)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getMessageText(t *testing.T) {
	t.Parallel()

	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success",
			args: args{
				text: laraCmd,
			},
			want: laraTxt,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := getMessageText(tt.args.text)
			assert.Equal(t, tt.want, got)
		})
	}
}
