# Bot de Telegram del Clima (Go)

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Telegram](https://img.shields.io/badge/Telegram-Bot-26A5E4?logo=telegram&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)

Un **bot de Telegram** hecho en **Go**: le escribes el nombre de una ciudad y te
responde el **clima actual** y un **pronóstico** de los próximos días. Usa la API
gratuita **Open-Meteo** (sin clave) y maneja concurrencia con la librería oficial
de bots de Telegram.

## Qué hace

| Le escribes... | Responde... |
|----------------|-------------|
| `/start` | Saludo de bienvenida |
| `/help` | Cómo usarlo |
| `Madrid` (una ciudad) | Clima actual + pronóstico de los próximos días |
| una ciudad que no existe | Mensaje amable para que revises cómo se escribe |

Ejemplo:
```
📍 Madrid, España
⛅ Ahora: 22°C, parcialmente nublado

Próximos días:
  2026-06-15 ☀️  26°C / 15°C
  2026-06-16 🌧️  19°C / 12°C
  2026-06-17 ☁️  21°C / 13°C
```

## Crear tu bot en Telegram (una sola vez)

1. En Telegram, busca **@BotFather** y ábrelo.
2. Envía `/newbot`.
3. Elige un **nombre** (lo que verá la gente) y un **usuario** que termine en `bot`
   (por ejemplo `climasamuel_bot`).
4. BotFather te dará un **token** (algo como `123456:ABC-DEF...`). **Cópialo y
   guárdalo en secreto** — quien lo tenga controla tu bot.

## Cómo correrlo

```bash
TELEGRAM_TOKEN="pega_aquí_tu_token" go run .
```

Verás `✅ Bot @tu_bot en marcha`. Búscalo en Telegram por su usuario y escríbele.

> 🔒 El token se lee de la variable de entorno `TELEGRAM_TOKEN` y **nunca** se
> guarda en el código (por eso no se sube a GitHub).

## Pruebas

```bash
go test ./...
```

Las pruebas del paquete `weather` usan un servidor de mentira (`httptest`), así
que no llaman a internet de verdad. Cubren: buscar una ciudad, ciudad no
encontrada, interpretar el pronóstico, y la traducción de códigos a texto/emoji.

## Estructura

```
main.go              arranque: token, conexión con Telegram y bucle de mensajes
bot/bot.go           decide la respuesta según el mensaje (/start, /help, ciudad)
weather/client.go    busca la ciudad y consulta el clima (Open-Meteo)
weather/conditions.go traduce los códigos del clima a texto + emoji
weather/types.go     tipos Place, Clima, Dia
```

## Stack

Go · [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) ·
API [Open-Meteo](https://open-meteo.com) (gratis, sin clave).
