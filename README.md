# Geeksonator Bot for PanteleevGroup chats

[![Go Report Card](https://goreportcard.com/badge/github.com/phpgeeks-club/admin-bot?style=flat-square)](https://goreportcard.com/report/github.com/phpgeeks-club/admin-bot)
[![Audit](https://github.com/phpgeeks-club/admin-bot/actions/workflows/audit.yml/badge.svg?branch=master)](https://github.com/phpgeeks-club/admin-bot/actions/workflows/audit.yml)
![License](https://img.shields.io/github/license/phpgeeks-club/admin-bot.svg)

## Install

```
curl -O https://raw.githubusercontent.com/phpgeeks-club/admin-bot/master/install.sh
chmod +x install.sh
./install.sh
rm ./install.sh
```

Also, the bot must disable Privacy mode (in BotFather) before being included in groups (otherwise it will not have access to messages to do reply)

#### Defaults

-   `GEEKSONATOR_TELEGRAM_BOT_TOKEN` = `""`
-   `GEEKSONATOR_TELEGRAM_TIMEOUT_SECONDS` = `15`
-   `GEEKSONATOR_DEBUG_MODE` = `false`
-   `GEEKSONATOR_DEBUG_TELEGRAM_BOT_TOKEN` = `""`

## Run in debug mode

1. In file `~/.geeksonator` set variables:
    ```
    GEEKSONATOR_DEBUG_MODE="true"
    GEEKSONATOR_DEBUG_TELEGRAM_BOT_TOKEN="debug_bot_token_here"
    ```
2. Run:
    ```
    docker run -d --env-file ~/.geeksonator --name geeksonator.app ghcr.io/phpgeeks-club/geeksonator:latest
    ```

## Run in development mode

1. Create file `.env` from `.env.example`
2. Set variables
    ```
    GEEKSONATOR_DEBUG_MODE="true"
    GEEKSONATOR_DEBUG_TELEGRAM_BOT_TOKEN="debug_bot_token_here"
    ```
3. Build image
    ```
    make docker_build
    ```
4. Run image
    ```
    make docker_run
    ```

## Comments

-   The use of `fmt.Errorf("error: %v", err)` instead of `fmt.Errorf("error: %w", err)` is due to the fact that this error is not unwrapped anywhere above.
