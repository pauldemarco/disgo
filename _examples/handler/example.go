package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetEnv("disgo_guild_id")

	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "ping",
			Description: "Replies with pong",
		},
		discord.SlashCommandCreate{
			Name:        "test",
			Description: "Replies with test",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommandGroup{
					Name:        "idk",
					Description: "Group command",
					Options: []discord.ApplicationCommandOptionSubCommand{
						{
							Name:        "sub",
							Description: "Sub command",
						},
					},
				},
				discord.ApplicationCommandOptionSubCommandGroup{
					Name:        "idk2",
					Description: "Group2 command",
					Options: []discord.ApplicationCommandOptionSubCommand{
						{
							Name:        "sub",
							Description: "Sub command",
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "sub2",
					Description: "Sub2 command",
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "ping2",
			Description: "Replies with pong2",
		},
	}
)

func main() {
	log.SetLevel(log.LevelInfo)
	log.Info("starting example...")
	log.Infof("disgo version: %s", disgo.Version)

	r := handler.New()
	r.Use(middleware.Logger)
	r.Group(func(r handler.Router) {
		r.Use(middleware.Print("group1"))
		r.Route("/test", func(r handler.Router) {
			r.Command("/sub2", handleContent("/test/sub2"))
			r.Route("/{group}", func(r handler.Router) {
				r.Command("/sub", handleVariableContent)
			})
		})
	})
	r.Group(func(r handler.Router) {
		r.Use(middleware.Print("group2"))
		r.Command("/ping", handlePing)
		r.Command("/ping2", handleContent("pong2"))
		r.Component("button1/{data}", handleComponent)
	})
	r.NotFound(handleNotFound)

	client, err := disgo.New(token,
		bot.WithDefaultGateway(),
		bot.WithEventListeners(r),
	)
	if err != nil {
		log.Fatal("error while building bot: ", err)
	}

	if err = handler.SyncCommands(client, commands, []snowflake.ID{guildID}); err != nil {
		log.Fatal("error while syncing commands: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func handleContent(content string) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		return event.CreateMessage(discord.MessageCreate{Content: content})
	}
}

func handleVariableContent(event *handler.CommandEvent) error {
	group := event.Variables["group"]
	return event.CreateMessage(discord.MessageCreate{Content: "group: " + group})
}

func handlePing(event *handler.CommandEvent) error {
	return event.CreateMessage(discord.MessageCreate{
		Content: "pong",
		Components: []discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.NewPrimaryButton("button1", "button1/testData"),
			},
		},
	})
}

func handleComponent(event *handler.ComponentEvent) error {
	data := event.Variables["data"]
	return event.CreateMessage(discord.MessageCreate{Content: "component: " + data})
}

func handleNotFound(event *events.InteractionCreate) error {
	return event.Respond(discord.InteractionResponseTypeCreateMessage, discord.MessageCreate{Content: "not found"})
}
