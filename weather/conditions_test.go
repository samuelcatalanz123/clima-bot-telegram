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
