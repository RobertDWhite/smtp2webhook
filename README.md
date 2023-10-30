# smtp2webhook

SMTP to Webhook Relay.

```yml
version: "3"
services:
  smtp2webhook:
    restart: always
    image: registry.white.fm/smtp2webhook:latest
    ports:
      - "25:25"
    environment:
      - WEBHOOK_URL=https://webhooks.white.fm/post/UNUmcoRPEj4mCmmTUEozFVbwtj2Fv6xb4qCdbfV9ZhBrZx8MrQk8jqTpHiq2Jyjs
```

Will forward all mail to the specified Webhook.
