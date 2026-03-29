package bot

import "github.com/bwmarrin/discordgo"

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "adventure",
		Description: "Begin or continue your adventure in the Legends of Eldoria!",
	},
	{
		Name:        "stats",
		Description: "View your character's stats and equipment.",
	},
	{
		Name:        "inventory",
		Description: "View your inventory.",
	},
	{
		Name:        "quests",
		Description: "View your quests and story progress.",
	},
}
