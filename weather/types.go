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
