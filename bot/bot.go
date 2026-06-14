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
