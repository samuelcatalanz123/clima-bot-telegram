# Diseño — Bot de Telegram del Clima (Go)

**Autor:** Samuel Catalán · **Fecha:** 2026-06-14

## Objetivo

Un **bot de Telegram** hecho en **Go** que da el clima de cualquier ciudad.
La gente lo usa desde su teléfono. Reutiliza la lógica de Open-Meteo del
proyecto `~/Repos/weather-go`. Para portafolio / aprendizaje.

## Alcance (nivel: Esencial + pronóstico)

- Comandos `/start` y `/help`.
- El usuario manda el nombre de una ciudad → el bot responde el clima **actual**
  (temperatura + estado con emoji) y un **pronóstico corto** (próximos días).
- Mensaje amable si la ciudad no existe o si la API falla.

Fuera de alcance (YAGNI): base de datos, ciudad favorita por usuario, mapas,
alertas. (Se pueden añadir después.)

## Arquitectura

```
clima-bot-telegram/
├── main.go                 lee TELEGRAM_TOKEN, arranca el bot, escucha mensajes
├── bot/bot.go              decide la respuesta según el mensaje/comando
├── weather/client.go       Geocode (buscar ciudad) + Forecast (clima) — Open-Meteo
├── weather/conditions.go   traduce el código del clima a texto + emoji
└── README.md
```

Flujo: `Telegram → main (polling) → bot (decide) → weather (consulta) → bot (responde)`.

Conexión por **long polling** (la librería pregunta a Telegram por mensajes
nuevos). No requiere servidor con IP/dominio público.

## Componentes

- **weather/client.go**: portado de `weather-go`. `Geocode(ciudad)` devuelve
  lat/lon y nombre; `Forecast(lat, lon)` devuelve clima actual + diario. Usa la
  API gratuita Open-Meteo (sin clave).
- **weather/conditions.go**: función `Describe(code int) (texto, emoji)` para los
  códigos WMO de Open-Meteo (0 = despejado ☀️, 61 = lluvia 🌧️, etc.).
- **bot/bot.go**: recibe un mensaje; si es `/start` o `/help`, responde el texto
  de ayuda; si no, lo trata como nombre de ciudad y arma la respuesta del clima.
- **main.go**: lee `TELEGRAM_TOKEN`, crea el bot con la librería
  `go-telegram-bot-api/v5`, y entra al bucle de mensajes.

## Comandos y respuestas

| Entrada | Respuesta |
|---------|-----------|
| `/start` | Bienvenida + cómo usarlo |
| `/help` | Explicación de uso |
| nombre de ciudad | Clima actual + pronóstico corto |
| ciudad inexistente | "No encontré esa ciudad 🤔, revisa cómo se escribe" |

Ejemplo de respuesta:
```
📍 Madrid, España
🌤️ Ahora: 22°C, parcialmente nublado

Próximos días:
  Mañana: ☀️ 25°C / 14°C
  Pasado: 🌧️ 19°C / 12°C
```

## Configuración (el token)

- El token de @BotFather se lee de la variable de entorno `TELEGRAM_TOKEN`.
- NUNCA se escribe en el código ni se sube a GitHub (un `.gitignore` evita
  subir archivos con secretos).
- Si falta al arrancar → el programa termina con un mensaje claro.

## Manejo de errores

- Open-Meteo falla → el bot responde "Hubo un problema, intenta de nuevo".
- Ciudad no encontrada → mensaje amable.
- `TELEGRAM_TOKEN` ausente → `log.Fatal` con mensaje claro.
- Un error procesando un mensaje no debe tumbar el bot (sigue con el siguiente).

## Pruebas

- `weather` con `httptest`: simular respuestas de Open-Meteo y verificar que
  `Geocode` y `Forecast` las interpretan bien (como en weather-go).
- `conditions`: verificar que los códigos se traducen al texto/emoji correctos.

## Criterios de éxito

1. `go build ./...` compila y `go test ./...` pasa.
2. Con un `TELEGRAM_TOKEN` válido, el bot responde en Telegram a `/start`,
   `/help` y a un nombre de ciudad con el clima.
3. Ciudad inexistente y caída de la API se manejan con mensajes amables.
4. README con cómo crear el bot en @BotFather y cómo correrlo.
