# Plan de Implementación — Bot de Telegram del Clima (Go)

> **Para trabajadores agénticos:** SUB-SKILL REQUERIDA: usar
> superpowers:executing-plans (o subagent-driven-development) para implementar
> tarea por tarea. Los pasos usan casillas (`- [ ]`).

**Goal:** Un bot de Telegram en Go que, al recibir el nombre de una ciudad,
responde el clima actual y un pronóstico corto.

**Architecture:** Paquete `weather` (consulta Open-Meteo, adaptado de weather-go),
paquete `bot` (decide la respuesta de texto), y `main.go` (conecta con Telegram
por long polling con la librería go-telegram-bot-api/v5).

**Tech Stack:** Go 1.26, go-telegram-bot-api/v5, API Open-Meteo (sin clave),
pruebas con httptest.

---

## Estructura de archivos

```
clima-bot-telegram/
├── main.go                      arranque: token, conexión Telegram, bucle de mensajes
├── bot/bot.go                   Responder(): decide la respuesta según el texto
├── weather/types.go             tipos Place, Clima, Dia
├── weather/client.go            Geocodificar + Pronostico (Open-Meteo)
├── weather/client_test.go       pruebas con httptest
├── weather/conditions.go        Describir() y EmojiDe() para códigos del clima
├── weather/conditions_test.go   pruebas de descripciones/emojis
├── go.mod
├── .gitignore
└── README.md
```

---

### Task 1: Inicializar el proyecto e instalar la librería de Telegram

- [ ] **Step 1: Inicializar el módulo**

Run:
```bash
cd ~/Repos/clima-bot-telegram
go mod init clima-bot-telegram
```
Expected: crea `go.mod` con `module clima-bot-telegram`.

- [ ] **Step 2: Instalar la librería de bots de Telegram**

Run:
```bash
go get github.com/go-telegram-bot-api/telegram-bot-api/v5
```
Expected: se añade a `go.mod`.

- [ ] **Step 3: Crear .gitignore**

`.gitignore`:
```
clima-bot-telegram
*.log
.DS_Store
.env
```

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum .gitignore
git commit -m "chore: inicializar proyecto e instalar librería de Telegram"
```

---

### Task 2: Tipos y condiciones del clima (con pruebas, TDD)

**Files:**
- Create: `weather/types.go`
- Create: `weather/conditions.go`
- Create: `weather/conditions_test.go`

- [ ] **Step 1: Crear los tipos**

`weather/types.go`:
```go
// Package weather consulta el clima desde la API gratuita Open-Meteo.
package weather

// Place es una ubicación encontrada por su nombre.
type Place struct {
	Name    string
	Country string
	Lat     float64
	Lon     float64
}

// Clima es el clima actual más el pronóstico de los próximos días.
type Clima struct {
	TempC       float64
	Code        int
	Descripcion string
	Emoji       string
	Dias        []Dia
}

// Dia es el pronóstico de un día.
type Dia struct {
	Fecha string
	MaxC  float64
	MinC  float64
	Code  int
	Emoji string
}
```

- [ ] **Step 2: Escribir las pruebas de condiciones (fallarán)**

`weather/conditions_test.go`:
```go
package weather

import "testing"

func TestDescribir(t *testing.T) {
	if got := Describir(0); got != "Despejado" {
		t.Errorf("Describir(0) = %q; se esperaba Despejado", got)
	}
	if got := Describir(99999); got != "Desconocido" {
		t.Errorf("Describir(99999) = %q; se esperaba Desconocido", got)
	}
}

func TestEmojiDe(t *testing.T) {
	if got := EmojiDe(0); got != "☀️" {
		t.Errorf("EmojiDe(0) = %q; se esperaba ☀️", got)
	}
	if got := EmojiDe(99999); got != "🌡️" {
		t.Errorf("EmojiDe(99999) = %q; se esperaba 🌡️ (por defecto)", got)
	}
}
```

- [ ] **Step 3: Ejecutar para ver que fallan**

Run: `go test ./weather/`
Expected: FAIL (Describir y EmojiDe no existen todavía).

- [ ] **Step 4: Implementar las condiciones**

`weather/conditions.go`:
```go
package weather

// descripciones traduce los códigos WMO de Open-Meteo a texto en español.
var descripciones = map[int]string{
	0: "Despejado", 1: "Mayormente despejado", 2: "Parcialmente nublado", 3: "Nublado",
	45: "Niebla", 48: "Niebla con escarcha",
	51: "Llovizna ligera", 53: "Llovizna moderada", 55: "Llovizna densa",
	61: "Lluvia ligera", 63: "Lluvia moderada", 65: "Lluvia fuerte",
	71: "Nieve ligera", 73: "Nieve moderada", 75: "Nieve fuerte",
	80: "Chubascos ligeros", 81: "Chubascos moderados", 82: "Chubascos violentos",
	95: "Tormenta", 96: "Tormenta con granizo ligero", 99: "Tormenta con granizo fuerte",
}

// emojis asocia cada código WMO con un emoji del clima.
var emojis = map[int]string{
	0: "☀️", 1: "🌤️", 2: "⛅", 3: "☁️",
	45: "🌫️", 48: "🌫️",
	51: "🌦️", 53: "🌦️", 55: "🌧️",
	61: "🌧️", 63: "🌧️", 65: "🌧️",
	71: "🌨️", 73: "🌨️", 75: "❄️",
	80: "🌦️", 81: "🌧️", 82: "⛈️",
	95: "⛈️", 96: "⛈️", 99: "⛈️",
}

// Describir devuelve el texto en español de un código del clima.
func Describir(code int) string {
	if d, ok := descripciones[code]; ok {
		return d
	}
	return "Desconocido"
}

// EmojiDe devuelve el emoji de un código del clima (🌡️ si no se reconoce).
func EmojiDe(code int) string {
	if e, ok := emojis[code]; ok {
		return e
	}
	return "🌡️"
}
```

- [ ] **Step 5: Ejecutar para ver que pasan**

Run: `go test ./weather/`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add weather/types.go weather/conditions.go weather/conditions_test.go
git commit -m "feat: tipos del clima y traducción de códigos (texto + emoji) con pruebas"
```

---

### Task 3: Cliente de Open-Meteo (con pruebas, TDD)

**Files:**
- Create: `weather/client_test.go`
- Create: `weather/client.go`

- [ ] **Step 1: Escribir las pruebas (fallarán)**

`weather/client_test.go`:
```go
package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// clienteDePrueba crea un Client que apunta a un servidor de mentira (httptest),
// así las pruebas no llaman a internet de verdad.
func clienteDePrueba(srv *httptest.Server) *Client {
	return &Client{
		geoBaseURL:      srv.URL,
		forecastBaseURL: srv.URL,
		http:            srv.Client(),
	}
}

func TestGeocodificarEncontrado(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"results":[{"name":"Madrid","country":"España","latitude":40.4,"longitude":-3.7}]}`))
	}))
	defer srv.Close()

	c := clienteDePrueba(srv)
	p, err := c.Geocodificar(context.Background(), "Madrid")
	if err != nil {
		t.Fatalf("Geocodificar falló: %v", err)
	}
	if p.Name != "Madrid" || p.Country != "España" {
		t.Errorf("lugar = %+v; se esperaba Madrid, España", p)
	}
}

func TestGeocodificarNoEncontrado(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"results":[]}`))
	}))
	defer srv.Close()

	c := clienteDePrueba(srv)
	if _, err := c.Geocodificar(context.Background(), "Xyzabc"); err == nil {
		t.Error("se esperaba error para una ciudad inexistente")
	}
}

func TestPronosticoParsea(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"current":{"temperature_2m":22.0,"weather_code":2},"daily":{"time":["2026-06-14","2026-06-15","2026-06-16","2026-06-17"],"weather_code":[2,0,61,3],"temperature_2m_max":[25,26,19,21],"temperature_2m_min":[14,15,12,13]}}`))
	}))
	defer srv.Close()

	c := clienteDePrueba(srv)
	clima, err := c.Pronostico(context.Background(), 40.4, -3.7)
	if err != nil {
		t.Fatalf("Pronostico falló: %v", err)
	}
	if clima.TempC != 22.0 {
		t.Errorf("TempC = %v; se esperaba 22.0", clima.TempC)
	}
	if clima.Descripcion != "Parcialmente nublado" {
		t.Errorf("Descripcion = %q; se esperaba Parcialmente nublado", clima.Descripcion)
	}
	if len(clima.Dias) != 3 {
		t.Errorf("se esperaban 3 días de pronóstico; hubo %d", len(clima.Dias))
	}
}
```

- [ ] **Step 2: Ejecutar para ver que fallan**

Run: `go test ./weather/`
Expected: FAIL (Client, Geocodificar, Pronostico no existen).

- [ ] **Step 3: Implementar el cliente**

`weather/client.go`:
```go
package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client habla con las APIs de geocodificación y pronóstico de Open-Meteo.
type Client struct {
	geoBaseURL      string
	forecastBaseURL string
	http            *http.Client
}

// NuevoCliente devuelve un Client apuntando a los endpoints reales de Open-Meteo.
func NuevoCliente() *Client {
	return &Client{
		geoBaseURL:      "https://geocoding-api.open-meteo.com/v1/search",
		forecastBaseURL: "https://api.open-meteo.com/v1/forecast",
		http:            &http.Client{Timeout: 10 * time.Second},
	}
}

// Geocodificar convierte el nombre de una ciudad en una ubicación (lat/lon).
func (c *Client) Geocodificar(ctx context.Context, ciudad string) (Place, error) {
	q := url.Values{}
	q.Set("name", ciudad)
	q.Set("count", "1")
	q.Set("language", "es")

	var body struct {
		Results []struct {
			Name      string  `json:"name"`
			Country   string  `json:"country"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}
	if err := c.getJSON(ctx, c.geoBaseURL+"?"+q.Encode(), &body); err != nil {
		return Place{}, err
	}
	if len(body.Results) == 0 {
		return Place{}, fmt.Errorf("no se encontró la ciudad %q", ciudad)
	}
	r := body.Results[0]
	return Place{Name: r.Name, Country: r.Country, Lat: r.Latitude, Lon: r.Longitude}, nil
}

// Pronostico devuelve el clima actual y hasta 3 días de pronóstico.
func (c *Client) Pronostico(ctx context.Context, lat, lon float64) (Clima, error) {
	q := url.Values{}
	q.Set("latitude", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Set("longitude", strconv.FormatFloat(lon, 'f', -1, 64))
	q.Set("current", "temperature_2m,weather_code")
	q.Set("daily", "weather_code,temperature_2m_max,temperature_2m_min")
	q.Set("forecast_days", "4")
	q.Set("timezone", "auto")

	var body struct {
		Current struct {
			Temperature float64 `json:"temperature_2m"`
			WeatherCode int     `json:"weather_code"`
		} `json:"current"`
		Daily struct {
			Time    []string  `json:"time"`
			Code    []int     `json:"weather_code"`
			TempMax []float64 `json:"temperature_2m_max"`
			TempMin []float64 `json:"temperature_2m_min"`
		} `json:"daily"`
	}
	if err := c.getJSON(ctx, c.forecastBaseURL+"?"+q.Encode(), &body); err != nil {
		return Clima{}, err
	}

	clima := Clima{
		TempC:       body.Current.Temperature,
		Code:        body.Current.WeatherCode,
		Descripcion: Describir(body.Current.WeatherCode),
		Emoji:       EmojiDe(body.Current.WeatherCode),
	}
	// Tomamos hasta 3 días, saltando hoy (índice 0).
	for i := 1; i < len(body.Daily.Time) && i <= 3; i++ {
		clima.Dias = append(clima.Dias, Dia{
			Fecha: body.Daily.Time[i],
			MaxC:  body.Daily.TempMax[i],
			MinC:  body.Daily.TempMin[i],
			Code:  body.Daily.Code[i],
			Emoji: EmojiDe(body.Daily.Code[i]),
		})
	}
	return clima, nil
}

// getJSON hace un GET y decodifica la respuesta JSON en dst.
func (c *Client) getJSON(ctx context.Context, endpoint string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("error de red: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("respuesta inesperada (código %d)", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		return fmt.Errorf("respuesta no válida: %w", err)
	}
	return nil
}
```

- [ ] **Step 4: Ejecutar para ver que pasan**

Run: `go test ./weather/ -v`
Expected: PASS (todas las pruebas de weather).

- [ ] **Step 5: Commit**

```bash
git add weather/client.go weather/client_test.go
git commit -m "feat: cliente de Open-Meteo (geocodificar + pronóstico) con pruebas"
```

---

### Task 4: La lógica del bot (qué responde)

**Files:**
- Create: `bot/bot.go`

- [ ] **Step 1: Crear la lógica del bot**

`bot/bot.go`:
```go
package bot

import (
	"context"
	"fmt"
	"strings"

	"clima-bot-telegram/weather"
)

const ayuda = "Soy un bot del clima 🌤️\n\n" +
	"Escríbeme el nombre de una ciudad (por ejemplo: Madrid) y te diré el clima " +
	"de ahora y el pronóstico de los próximos días.\n\n" +
	"Comandos:\n/start - empezar\n/help - esta ayuda"

// Responder decide el texto de respuesta según lo que escribió el usuario.
func Responder(ctx context.Context, c *weather.Client, texto string) string {
	texto = strings.TrimSpace(texto)

	switch {
	case texto == "":
		return "Mándame el nombre de una ciudad 🌤️"
	case strings.HasPrefix(texto, "/start"):
		return "¡Hola! 👋 Soy tu bot del clima. Mándame el nombre de una ciudad y te digo cómo está el tiempo. 🌤️"
	case strings.HasPrefix(texto, "/help"):
		return ayuda
	}

	lugar, err := c.Geocodificar(ctx, texto)
	if err != nil {
		return "No encontré esa ciudad 🤔, revisa cómo se escribe."
	}
	clima, err := c.Pronostico(ctx, lugar.Lat, lugar.Lon)
	if err != nil {
		return "Hubo un problema al consultar el clima, intenta de nuevo 🙏"
	}
	return formatear(lugar, clima)
}

// formatear arma el mensaje bonito con el clima.
func formatear(lugar weather.Place, clima weather.Clima) string {
	var b strings.Builder
	fmt.Fprintf(&b, "📍 %s, %s\n", lugar.Name, lugar.Country)
	fmt.Fprintf(&b, "%s Ahora: %.0f°C, %s\n", clima.Emoji, clima.TempC, clima.Descripcion)
	if len(clima.Dias) > 0 {
		b.WriteString("\nPróximos días:\n")
		for _, d := range clima.Dias {
			fmt.Fprintf(&b, "  %s %s  %.0f°C / %.0f°C\n", d.Fecha, d.Emoji, d.MaxC, d.MinC)
		}
	}
	return b.String()
}
```

- [ ] **Step 2: Verificar que compila**

Run: `go build ./...`
Expected: sin errores.

- [ ] **Step 3: Commit**

```bash
git add bot/bot.go
git commit -m "feat: lógica del bot (start, help y respuesta del clima)"
```

---

### Task 5: main.go (conexión con Telegram)

**Files:**
- Create: `main.go`

- [ ] **Step 1: Crear main.go**

`main.go`:
```go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"clima-bot-telegram/bot"
	"clima-bot-telegram/weather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("Falta TELEGRAM_TOKEN. Ejemplo: TELEGRAM_TOKEN=\"123:abc\" go run .")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("No se pudo iniciar el bot (¿token inválido o sin internet?): ", err)
	}
	log.Printf("✅ Bot @%s en marcha. Escríbele en Telegram.", api.Self.UserName)

	cliente := weather.NuevoCliente()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		respuesta := bot.Responder(ctx, cliente, update.Message.Text)
		cancel()

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, respuesta)
		if _, err := api.Send(msg); err != nil {
			log.Println("error al enviar respuesta:", err)
		}
	}
}
```

- [ ] **Step 2: Verificar que compila y pasan las pruebas**

Run: `go build ./... && go test ./...`
Expected: compila y los tests PASS.

- [ ] **Step 3: Commit**

```bash
git add main.go
git commit -m "feat: conexión con Telegram (long polling) y bucle de mensajes"
```

---

### Task 6: README

**Files:**
- Create: `README.md`

- [ ] **Step 1: Escribir el README**

`README.md` con: descripción; **cómo crear el bot en @BotFather** (abrir Telegram,
buscar @BotFather, `/newbot`, elegir nombre y usuario, copiar el token); cómo
correrlo (`TELEGRAM_TOKEN="tu_token" go run .`); qué hace (comandos + escribir
ciudad); ejemplo de respuesta; estructura de carpetas; cómo correr las pruebas
(`go test ./...`); nota de que el token nunca se sube a GitHub; stack (Go,
go-telegram-bot-api, Open-Meteo).

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: README con instrucciones de BotFather"
```

---

### Task 7: Crear el bot en BotFather y probar en vivo (manual, guiado)

- [ ] **Step 1: Crear el bot en Telegram (lo hace el usuario)**

En Telegram (teléfono o web): buscar **@BotFather** → `/newbot` → elegir un
nombre y un usuario que termine en `bot` → copiar el **token** que devuelve.

- [ ] **Step 2: Arrancar el bot con el token**

Run:
```bash
cd ~/Repos/clima-bot-telegram
TELEGRAM_TOKEN="PEGAR_TOKEN_AQUI" go run .
```
Expected: imprime `✅ Bot @<nombre> en marcha.`

- [ ] **Step 3: Probar en Telegram**

Buscar el bot por su usuario, enviar `/start`, luego `Madrid`, luego una ciudad
inexistente como `Xyzabc`.
Expected: responde bienvenida; clima + pronóstico; y mensaje amable de "no
encontré esa ciudad".

- [ ] **Step 4: Detener el bot**

`Ctrl + C` en la terminal donde corre.

---

## Self-Review (cobertura del spec)

- Comandos /start y /help → Task 4 ✅
- Ciudad → clima actual + pronóstico → Tasks 3 (Pronostico), 4 (formatear) ✅
- Ciudad inexistente / API falla → Task 4 (mensajes amables) ✅
- Token por variable de entorno, no en el código → Tasks 1 (.gitignore), 5 ✅
- Pruebas (httptest + condiciones) → Tasks 2, 3 ✅
- README con BotFather → Task 6 ✅
- Probar en vivo → Task 7 ✅
