package commands

import (
	"github.com/disgoorg/disgo/discord"
)

var (
	Commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "say",
			Description: "says what you say",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "message",
					Description: "What to say",
					Required:    true,
				},
				discord.ApplicationCommandOptionBool{
					Name:        "ephemeral",
					Description: "If the response should only be visible to you",
					Required:    true,
				},
			},
		},
	}
)
