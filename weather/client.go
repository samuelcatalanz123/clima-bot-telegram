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
