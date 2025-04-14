package geeksonator

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"geeksonator/internal/observer"
	"geeksonator/internal/provider/telegram"
	cacher "geeksonator/pkg/cache"
)

const (
	cacheMaxSize = 1
	cacheTTL     = 24 * time.Hour
)

// Start starts the application.
func Start() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	defer stop()

	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("LoadConfig: %v", err)
	}

	logger, err := newLogger(cfg.DebugMode)
	if err != nil {
		return fmt.Errorf("newLogger: %v", err)
	}
	defer logger.Sync() //nolint:errcheck

	var tgBotToken string
	if cfg.DebugMode {
		tgBotToken = cfg.DebugTgBotToken
	} else {
		tgBotToken = cfg.TgBotToken
	}
	botAPI, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		return fmt.Errorf("tgbotapi.NewBotAPI: %v", err)
	}
	logger.Info("Authorized on account",
		zap.String("account", botAPI.Self.UserName),
	)

	telegramService := telegram.NewService(botAPI)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = cfg.TgTimeoutSeconds // long polling

	updatesChan := botAPI.GetUpdatesChan(updateConfig)

	cache, err := cacher.NewCacher[string, []tgbotapi.ChatMember](
		cacheMaxSize,
		cacheTTL,
		cacher.WithDebug[string, []tgbotapi.ChatMember](logger),
	)
	if err != nil {
		return fmt.Errorf("cacher.NewCacher: %v", err)
	}

	var observerManager *observer.Manager
	if cfg.DebugMode {
		observerManager = observer.NewManager(
			telegramService,
			updatesChan,
			cache,
			observer.WithDebug(logger),
			observer.WithSkipAdminCheck(),
		)
	} else {
		observerManager = observer.NewManager(
			telegramService,
			updatesChan,
			cache,
			observer.WithDebug(logger),
		)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err = observerManager.Run(ctx)
		if err != nil {
			logger.Error("Observer manager Run error",
				zap.Error(err),
			)

			return
		}
		logger.Info("Observer manager gracefully stopped")
	}()

	wg.Wait()

	logger.Info("Application stopped")

	return nil
}

// newLogger creates new logger.
func newLogger(debugMode bool) (*zap.Logger, error) {
	if debugMode {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("zap.NewDevelopment: %v", err)
		}

		logger.Debug("Debug mode running")

		return logger, nil
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("zap.NewProduction: %v", err)
	}

	return logger, nil
}
