# Bot do Telegram para Eventos da Comunidade Devs Norte

Este é um bot do Telegram desenvolvido em Go que permite aos usuários buscar eventos disponíveis e encerrados da comunidade Devs Norte usando a API do Sympla.

## Funcionalidades

- `/disponiveis`: Comando para listar os eventos disponíveis.
- `/encerrados`: Comando para listar os eventos encerrados.

## Pré-requisitos

Antes de executar este bot, você precisa ter:
- Um token de bot do Telegram. Você pode obter um conversando com o [BotFather](https://t.me/BotFather).

## Configuração

1. Clone este repositório:

```bash
git clone https://github.com/ecsistem/bot-telegram-devs-norte.git
```
2. Acesse o diretório do projeto:
```bash
cd bot-telegram-devs-norte
```
3. Crie um arquivo `.env` na raiz do projeto e adicione o token do bot do Telegram:
```bash
TELEGRAM_BOT_TOKEN=SeuTokenDoBotTelegram
```

## Instalação

1. Instale as dependências do projeto:
```bash
go mod tidy
```

## Uso

1. Execute o bot:
```bash
go run main.go
```

2. No Telegram, encontre o seu bot e inicie uma conversa.
3. Use os comandos `/disponiveis` e `/encerrados` para listar os eventos disponíveis e encerrados, respectivamente.

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir uma issue ou enviar um pull request.