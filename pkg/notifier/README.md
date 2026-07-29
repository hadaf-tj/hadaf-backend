```md
# Telegram Bot Setup

This project uses a Telegram bot for sending backend alerts.

## 1. Create a Telegram Bot

1. Open Telegram
2. Find **@BotFather**
3. Send:

```

/newbot

```

4. Follow the instructions.
5. Copy the generated bot token.

Example:

```

1234567890:AAxxxxxxxxxxxxxxxxxxxx

```

---

## 2. Connect Bot to Chat

Open your new bot in Telegram and send any message to it:

```

Hello

```

Then open this URL in your browser:

```

[https://api.telegram.org/bot](https://api.telegram.org/bot)<YOUR_BOT_TOKEN>/getUpdates

```

Replace:

```

<YOUR_BOT_TOKEN>

````

with your actual token.

Find:

```json
"chat": {
    "id": 123456789
}
````

The value of `id` is your chat ID.

Example:

```
TELEGRAM_ALERT_CHAT_ID=123456789
```

---

## 3. Configure Environment Variables

Add these variables to your `.env` file:

```env
TELEGRAM_ALERT_TOKEN=your_bot_token
TELEGRAM_ALERT_CHAT_ID=your_chat_id
TELEGRAM_BASE_URL=https://api.telegram.org
```

Example:

```env
TELEGRAM_ALERT_TOKEN=1234567890:AAxxxxxxxxxxxxxxxx
TELEGRAM_ALERT_CHAT_ID=123456789
TELEGRAM_BASE_URL=https://api.telegram.org
```

---

## 4. Restart Application

After changing `.env`, restart the application:

```bash
docker compose down
docker compose up --build
```

The backend will now be able to send Telegram notifications.

---

## Notes

* Keep your bot token private.
* Do not commit `.env` files.
* If the bot does not send messages, make sure:

  * you started a chat with the bot
  * the token is correct
  * the chat ID belongs to the correct Telegram account/group

```
```
