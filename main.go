package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Bot struct {
	bot   *telego.Bot
	bh    *th.BotHandler
	token string
}

type SymplaResponse struct {
	Data []Event `json:"data"`
}

type Location struct {
	Country      string  `json:"country"`
	Address      string  `json:"address"`
	AddressAlt   string  `json:"address_alt"`
	City         string  `json:"city"`
	AddressNum   string  `json:"address_num"`
	Name         string  `json:"name"`
	Longitude    float64 `json:"lon"`
	State        string  `json:"state"`
	Neighborhood string  `json:"neighborhood"`
	ZipCode      string  `json:"zip_code"`
	Latitude     float64 `json:"lat"`
}

type Images struct {
	Original string `json:"original"`
	XS       string `json:"xs"`
	LG       string `json:"lg"`
}

type StartDateFormats struct {
	Pt string `json:"pt"`
	En string `json:"en"`
	Es string `json:"es"`
}
type EndDateFormats struct {
	Pt string `json:"pt"`
	En string `json:"en"`
	Es string `json:"es"`
}
type Event struct {
	Name             string           `json:"name"`
	Images           Images           `json:"images"`
	Location         Location         `json:"location"`
	StartDateFormats StartDateFormats `json:"start_date_formats"`
	EndDateFormats   EndDateFormats   `json:"end_date_formats"`
	URL              string           `json:"url"`
}

func NewBot(token string) (*Bot, error) {
	bot, err := telego.NewBot(token, telego.WithDefaultDebugLogger())
	if err != nil {
		return nil, err
	}

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		return nil, err
	}

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot:   bot,
		bh:    bh,
		token: token,
	}, nil
}

func (b *Bot) Start() {
	defer b.bh.Stop()
	defer b.bot.StopLongPolling()

	b.registerCommands()

	b.bh.Start()
}

func (b *Bot) registerCommands() {
	b.registerBotCommand()
	b.registerEventCommands()
}

func (b *Bot) registerBotCommand() {
	b.bh.Handle(func(bot *telego.Bot, update telego.Update) {
		infoMessage := `
👋 Bem-vindo ao Bot da Comunidade Devs Norte! 🚀

Este bot está aqui para ajudá-lo a encontrar os eventos mais recentes e emocionantes hospedados no Sympla pela nossa comunidade.

Para consultar os eventos disponíveis, basta digitar /disponiveis. E se estiver interessado nos eventos que já passaram, digite /encerrados.

Fique à vontade para explorar e participar dos eventos que mais lhe interessarem!😊
`

		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			infoMessage,
		))
	}, th.CommandEqual("start"))
}

func (b *Bot) registerEventCommands() {
	b.registerAvailableEventsCommand()
	b.registerClosedEventsCommand()
}

func (b *Bot) registerAvailableEventsCommand() {
	b.bh.Handle(func(bot *telego.Bot, update telego.Update) {
		events, err := fetchSymplaEvents("future")
		if err != nil {
			fmt.Println("Erro ao buscar eventos:", err)
			return
		}
		message := formatEventsMessage(events)
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			message,
		))
	}, th.CommandEqual("disponiveis"))
}

func (b *Bot) registerClosedEventsCommand() {
	b.bh.Handle(func(bot *telego.Bot, update telego.Update) {
		events, err := fetchSymplaEvents("past")
		if err != nil {
			fmt.Println("Erro ao buscar eventos:", err)
			return
		}
		message := formatEventsMessage(events)
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			message,
		))
	}, th.CommandEqual("encerrados"))
}

// Função para buscar eventos do Sympla
func fetchSymplaEvents(eventType string) ([]Event, error) {
	organizerIDs := []int{3125215, 5478152}

	service := "/v4/search"
	if eventType == "past" {
		service = "/v4/events/past"
	}
	requestBody := fmt.Sprintf(`{
		"service": "%s",
		"params": {
			"only": "name,images,location,start_date_formats,end_date_formats,url",
			"organizer_id": %s,
			"sort": "date",
			"order_by": "desc",
			"limit": "6",
			"page": 1
		},
		"ignoreLocation": true
	}`, service, intArrayToString(organizerIDs))

	fmt.Println("Request body:", requestBody)

	resp, err := http.Post("https://www.sympla.com.br/api/v1/search", "application/json", strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("Status code:", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("Response body:", string(body))

	var symplaResp SymplaResponse
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&symplaResp); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	return symplaResp.Data, nil
}

func intArrayToString(arr []int) string {
	strArr := make([]string, len(arr))
	for i, num := range arr {
		strArr[i] = fmt.Sprint(num)
	}
	return "[" + strings.Join(strArr, ",") + "]"
}

func formatEventsMessage(events []Event) string {
	message := "#BOT Devs Norte 🤖\n\n\n"
	if events == nil || len(events) == 0 {
		message += "Ops... Nem um evento disponivel no momento, mas não fique triste logo estaremos fazendo mais eventos! 🥺\n\n\n"
	} else {
		message += "🎉 Eventos: 🎉\n\n\n"
		for _, event := range events {
			message += fmt.Sprintf("- %s\n  Local: %s\n  Data: %s\n  URL: %s\n \n\n\n", event.Name, event.Location.City, event.StartDateFormats.Pt, event.URL)
			message += "----------------------------------------\n\n\n"
		}
	}
	return message
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar o arquivo .env:", err)
		os.Exit(1)
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		fmt.Println("Token do bot do Telegram não fornecido")
		os.Exit(1)
	}

	bot, err := NewBot(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bot.Start()
}
