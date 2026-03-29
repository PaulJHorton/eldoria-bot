# Eldoria Bot

A Discord RPG bot set in the world of Eldoria — fight monsters, level up, and adventure with friends.

## Features

- **Character Creation** — Name your hero and begin your journey in Eldoria Village
- **Turn-Based Combat** — Battle monsters with attacks, magic, and items
- **5 Unique Areas** — Explore from the peaceful village to the Shadow Citadel
- **Boss Fights** — Defeat powerful bosses to advance the story
- **Equipment & Items** — Buy weapons, armor, shields, and potions from shops
- **Leveling System** — Gain XP, level up, and grow stronger
- **Quests & Story** — Progress through an unfolding storyline as you explore

## Slash Commands

| Command | Description |
|---------|-------------|
| `/adventure` | Begin or continue your adventure |
| `/stats` | View your character's stats and equipment |
| `/inventory` | View your inventory |
| `/quests` | View your quests and story progress |

## The World of Eldoria

| Area | Level | Description |
|------|-------|-------------|
| Eldoria Village | — | A peaceful village nestled in a green valley |
| Whispering Forest | 1+ | An ancient forest shrouded in mist |
| Crystal Caves | 4+ | Underground caverns illuminated by glowing crystals |
| Frozen Peaks | 7+ | Treacherous mountain peaks blanketed in eternal snow |
| Shadow Citadel | 10+ | The final stronghold of the Dark Lord |

## Tech Stack

- **Language:** Go
- **Discord Library:** [discordgo](https://github.com/bwmarrin/discordgo)
- **Database:** SQLite via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)

## Setup

1. Clone the repo
2. Create a `.env` file with your bot token:
   ```
   DISCORD_TOKEN=your_token_here
   APP_ID=your_app_id_here
   ```
3. Run the bot:
   ```bash
   go run main.go
   ```
