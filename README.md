# Geeksonator Bot for PanteleevGroup chats

## Install

```
export GEEKSONATOR_TELEGRAM_BOT_TOKEN='bot_token_here'
sudo chmod 755 /opt/geeksonator
sudo cp geeksonator.service /lib/systemd/user
sudo systemctl enable /lib/systemd/user/geeksonator.service
sudo service geeksonator start
```

Also, the bot must disable Privacy mode (in BotFather) before being included in groups (otherwise it will not have access to messages to do reply)
