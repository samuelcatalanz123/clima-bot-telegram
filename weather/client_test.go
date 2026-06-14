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
