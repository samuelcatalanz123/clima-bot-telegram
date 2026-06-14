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
