// main.go
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Malix-Learning/stateless-discord-bot-Zlatovo/src/commands"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
)

var (
	token     string = os.Getenv("disgo_token")
	publicKey string = os.Getenv("disgo_public_key")
)

func main() {
	slog.Info(token)
	slog.Info("starting example...")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	// use custom ed25519 verify implementation
	httpserver.Verify = func(publicKey httpserver.PublicKey, message, sig []byte) bool {
		return ed25519.Verify(publicKey, message, sig)
	}

	client, err := disgo.New(token,
		bot.WithHTTPServerConfigOpts(publicKey,
			httpserver.WithURL("/interactions/callback"),
			httpserver.WithAddress(":80"),
		),
		bot.WithEventListenerFunc(commandListener),
	)
	if err != nil {
		panic("error while building disgo instance: " + err.Error())
	}

	defer client.Close(context.TODO())

	if _, err = client.Rest().SetGlobalCommands(client.ApplicationID(), commands.Commands); err != nil {
		panic("error while registering commands: " + err.Error())
	}

	if err = client.OpenHTTPServer(); err != nil {
		panic("error while starting http server: " + err.Error())
	}

	slog.Info("example is now running. Press CTRL-C to exit.")
	keepAlive()
}

func keepAlive() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func commandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	if data.CommandName() == "say" {
		if err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent(data.String("message")).
			SetEphemeral(data.Bool("ephemeral")).
			Build(),
		); err != nil {
			event.Client().Logger().Error("error on sending response: ", err)
		}
	}
}
