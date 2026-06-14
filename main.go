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
