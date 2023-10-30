# smtp2webhook

SMTP to Webhook Relay.

This rendition will forward all mail to the specified webhook.

Build yourself with the Dockerfile, or pull registry.white.fm/smtp2webhook:latest

To build, 
```docker build -t smtp2webhook:2.2 .  ```

---
docker-compose

```yml
version: "3"
services:
  smtp2webhook:
    restart: always
    image: registry.white.fm/smtp2webhook:latest
    ports:
      - "25:25"
    environment:
      - WEBHOOK_URL=https://EXAMPLE.COM
```


