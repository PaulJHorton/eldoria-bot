package bot

import (
	"fmt"
	"strings"

	"discord-rpg-bot/game"

	"github.com/bwmarrin/discordgo"
)

func welcomeEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "\u2694\uFE0F Legends of Eldoria",
		Description: "Welcome, brave adventurer! A great darkness threatens the land of Eldoria. Only a true hero can save us.\n\nWill you answer the call?",
		Color:       0x00AE86,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "How to Play", Value: "Create your character and explore the world! Fight monsters, collect loot, and complete quests to save Eldoria.", Inline: false},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: "Click below to begin your adventure!"},
	}
}

func welcomeComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Begin Adventure",
					Style:    discordgo.SuccessButton,
					CustomID: "create_character",
					Emoji:    &discordgo.ComponentEmoji{Name: "\u2694\uFE0F"},
				},
			},
		},
	}
}

func villageEmbed(c *game.Character) *discordgo.MessageEmbed {
	area := game.Areas["village"]
	currentQuest := game.GetCurrentQuest(c.StoryProgress)
	questInfo := "All quests complete! You are the Hero of Eldoria!"
	if currentQuest != nil {
		questInfo = fmt.Sprintf("**%s**: %s", currentQuest.Name, currentQuest.Description)
	}

	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s %s", area.Emoji, area.Name),
		Description: area.Description,
		Color:       0x2ECC71,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Hero", Value: fmt.Sprintf("**%s** \u2014 Level %d", c.Name, c.Level), Inline: true},
			{Name: "HP", Value: fmt.Sprintf("%s %d/%d", makeBar(c.HP, c.MaxHP, 10), c.HP, c.MaxHP), Inline: true},
			{Name: "Gold", Value: fmt.Sprintf("%d", c.Gold), Inline: true},
			{Name: "Current Quest", Value: questInfo, Inline: false},
		},
	}
}

func villageComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "Explore", Style: discordgo.PrimaryButton, CustomID: "explore", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F5FA\uFE0F"}},
				discordgo.Button{Label: "Shop", Style: discordgo.PrimaryButton, CustomID: "shop", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F6D2"}},
				discordgo.Button{Label: "Inventory", Style: discordgo.SecondaryButton, CustomID: "inventory", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F392"}},
				discordgo.Button{Label: "Quests", Style: discordgo.SecondaryButton, CustomID: "quests", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F4DC"}},
				discordgo.Button{Label: "Stats", Style: discordgo.SecondaryButton, CustomID: "stats", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F4CA"}},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "Rest at Inn (10 Gold)", Style: discordgo.SuccessButton, CustomID: "rest", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F6CF\uFE0F"}},
				discordgo.Button{Label: "Travel", Style: discordgo.PrimaryButton, CustomID: "travel", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F9ED"}},
			},
		},
	}
}

func combatEmbed(cs *game.CombatSession) *discordgo.MessageEmbed {
	logText := strings.Join(cs.Log, "\n")
	if logText == "" {
		logText = "..."
	}

	hpBar := makeBar(cs.Player.HP, cs.Player.MaxHP, 10)
	mpBar := makeBar(cs.Player.MP, cs.Player.MaxMP, 10)
	enemyHPBar := makeBar(cs.Enemy.HP, cs.Enemy.MaxHP, 10)

	color := 0xE74C3C
	if cs.Finished && cs.Victory {
		color = 0x2ECC71
	} else if cs.Finished && !cs.Victory {
		color = 0x95A5A6
	}

	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("\u2694\uFE0F Battle: %s vs %s", cs.Player.Name, cs.Enemy.Name),
		Color: color,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("%s (Lv.%d)", cs.Player.Name, cs.Player.Level),
				Value:  fmt.Sprintf("HP: %s %d/%d\nMP: %s %d/%d", hpBar, cs.Player.HP, cs.Player.MaxHP, mpBar, cs.Player.MP, cs.Player.MaxMP),
				Inline: true,
			},
			{
				Name:   cs.Enemy.Name,
				Value:  fmt.Sprintf("HP: %s %d/%d", enemyHPBar, cs.Enemy.HP, cs.Enemy.MaxHP),
				Inline: true,
			},
			{
				Name:   "Combat Log",
				Value:  logText,
				Inline: false,
			},
		},
	}
}

func combatComponents(cs *game.CombatSession) []discordgo.MessageComponent {
	if cs.Finished {
		return []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{Label: "Continue", Style: discordgo.SuccessButton, CustomID: "continue_after_combat", Emoji: &discordgo.ComponentEmoji{Name: "\u25B6\uFE0F"}},
				},
			},
		}
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "Attack", Style: discordgo.DangerButton, CustomID: "combat_attack", Emoji: &discordgo.ComponentEmoji{Name: "\u2694\uFE0F"}},
				discordgo.Button{Label: "Defend", Style: discordgo.PrimaryButton, CustomID: "combat_defend", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F6E1\uFE0F"}},
				discordgo.Button{Label: "Magic", Style: discordgo.PrimaryButton, CustomID: "combat_magic", Emoji: &discordgo.ComponentEmoji{Name: "\u2728"}},
				discordgo.Button{Label: "Item", Style: discordgo.SecondaryButton, CustomID: "combat_item", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F9EA"}},
				discordgo.Button{Label: "Flee", Style: discordgo.SecondaryButton, CustomID: "combat_flee", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F3C3"}},
			},
		},
	}
}

func statsEmbed(c *game.Character) *discordgo.MessageEmbed {
	weaponName := "None"
	if item, ok := game.Items[c.Weapon]; ok {
		weaponName = item.Name
	}
	shieldName := "None"
	if item, ok := game.Items[c.Shield]; ok {
		shieldName = item.Name
	}
	armorName := "None"
	if item, ok := game.Items[c.Armor]; ok {
		armorName = item.Name
	}

	area := game.Areas[c.CurrentArea]
	areaName := c.CurrentArea
	if area != nil {
		areaName = fmt.Sprintf("%s %s", area.Emoji, area.Name)
	}

	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("\U0001F4CA %s \u2014 Level %d", c.Name, c.Level),
		Color: 0x3498DB,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "HP", Value: fmt.Sprintf("%s %d/%d", makeBar(c.HP, c.MaxHP, 10), c.HP, c.MaxHP), Inline: false},
			{Name: "MP", Value: fmt.Sprintf("%s %d/%d", makeBar(c.MP, c.MaxMP, 10), c.MP, c.MaxMP), Inline: false},
			{Name: "Attack", Value: fmt.Sprintf("%d (+%d)", c.Attack, c.TotalAttack()-c.Attack), Inline: true},
			{Name: "Defense", Value: fmt.Sprintf("%d (+%d)", c.Defense, c.TotalDefense()-c.Defense), Inline: true},
			{Name: "Speed", Value: fmt.Sprintf("%d", c.Speed), Inline: true},
			{Name: "Magic", Value: fmt.Sprintf("%d (+%d)", c.Magic, c.TotalMagic()-c.Magic), Inline: true},
			{Name: "Gold", Value: fmt.Sprintf("%d", c.Gold), Inline: true},
			{Name: "XP", Value: fmt.Sprintf("%d / %d", c.XP, c.XPToNext), Inline: true},
			{Name: "Equipment", Value: fmt.Sprintf("\u2694\uFE0F %s\n\U0001F6E1\uFE0F %s\n\U0001F9E5 %s", weaponName, shieldName, armorName), Inline: false},
			{Name: "Location", Value: areaName, Inline: true},
			{Name: "Story Progress", Value: fmt.Sprintf("%d/5", c.StoryProgress), Inline: true},
		},
	}
}

func backButton() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "Back", Style: discordgo.SecondaryButton, CustomID: "back_village", Emoji: &discordgo.ComponentEmoji{Name: "\u25C0\uFE0F"}},
			},
		},
	}
}

func inventoryEmbed(c *game.Character, inv []game.InventorySlot) *discordgo.MessageEmbed {
	var items strings.Builder
	if len(inv) == 0 {
		items.WriteString("Your inventory is empty.")
	} else {
		for _, slot := range inv {
			if item, ok := game.Items[slot.ItemID]; ok {
				items.WriteString(fmt.Sprintf("**%s** x%d \u2014 %s\n", item.Name, slot.Quantity, item.Description))
			}
		}
	}

	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("\U0001F392 %s's Inventory", c.Name),
		Description: items.String(),
		Color:       0xF39C12,
	}
}

func inventoryComponents(inv []game.InventorySlot) []discordgo.MessageComponent {
	components := []discordgo.MessageComponent{}

	var options []discordgo.SelectMenuOption
	for _, slot := range inv {
		item, ok := game.Items[slot.ItemID]
		if !ok {
			continue
		}
		label := fmt.Sprintf("%s (x%d)", item.Name, slot.Quantity)
		var desc string
		switch item.Type {
		case game.ItemTypeConsumable:
			desc = "Use this item"
		case game.ItemTypeWeapon:
			desc = "Equip as weapon"
		case game.ItemTypeShield:
			desc = "Equip as shield"
		case game.ItemTypeArmor:
			desc = "Equip as armor"
		default:
			desc = "Quest item"
		}
		options = append(options, discordgo.SelectMenuOption{
			Label:       label,
			Description: desc,
			Value:       slot.ItemID,
		})
	}

	if len(options) > 0 {
		if len(options) > 25 {
			options = options[:25]
		}
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "inv_select",
					Placeholder: "Select an item to use/equip",
					Options:     options,
				},
			},
		})
	}

	components = append(components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "Back", Style: discordgo.SecondaryButton, CustomID: "back_village", Emoji: &discordgo.ComponentEmoji{Name: "\u25C0\uFE0F"}},
		},
	})

	return components
}

func shopEmbed(c *game.Character, areaID string) *discordgo.MessageEmbed {
	shopItems, ok := game.ShopItems[areaID]
	if !ok {
		shopItems = game.ShopItems["village"]
	}

	var items strings.Builder
	for _, itemID := range shopItems {
		if item, ok := game.Items[itemID]; ok {
			items.WriteString(fmt.Sprintf("**%s** \u2014 %d Gold \u2014 %s\n", item.Name, item.BuyPrice, item.Description))
		}
	}

	return &discordgo.MessageEmbed{
		Title:       "\U0001F6D2 Shop",
		Description: fmt.Sprintf("Your Gold: **%d**\n\n%s", c.Gold, items.String()),
		Color:       0xE67E22,
	}
}

func shopComponents(areaID string) []discordgo.MessageComponent {
	shopItems, ok := game.ShopItems[areaID]
	if !ok {
		shopItems = game.ShopItems["village"]
	}

	var options []discordgo.SelectMenuOption
	for _, itemID := range shopItems {
		if item, ok := game.Items[itemID]; ok {
			options = append(options, discordgo.SelectMenuOption{
				Label:       fmt.Sprintf("%s \u2014 %d Gold", item.Name, item.BuyPrice),
				Description: item.Description,
				Value:       itemID,
			})
		}
	}

	components := []discordgo.MessageComponent{}
	if len(options) > 0 {
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "shop_buy",
					Placeholder: "Select an item to buy",
					Options:     options,
				},
			},
		})
	}

	components = append(components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "Back", Style: discordgo.SecondaryButton, CustomID: "back_village", Emoji: &discordgo.ComponentEmoji{Name: "\u25C0\uFE0F"}},
		},
	})

	return components
}

func questsEmbed(c *game.Character) *discordgo.MessageEmbed {
	var desc strings.Builder

	for i, qID := range game.QuestOrder {
		q := game.Quests[qID]
		status := "\U0001F512"
		if i < c.StoryProgress {
			status = "\u2705"
		} else if i == c.StoryProgress {
			status = "\U0001F4CC"
		}
		desc.WriteString(fmt.Sprintf("%s **%s**\n%s\n\n", status, q.Name, q.Description))
	}

	return &discordgo.MessageEmbed{
		Title:       "\U0001F4DC Quest Log",
		Description: desc.String(),
		Color:       0x9B59B6,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Story Progress: %d/5", c.StoryProgress)},
	}
}

func questsComponents(c *game.Character) []discordgo.MessageComponent {
	components := []discordgo.MessageComponent{}

	currentQuest := game.GetCurrentQuest(c.StoryProgress)
	if currentQuest != nil && currentQuest.ID == "awakening" {
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "Speak to the Elder", Style: discordgo.SuccessButton, CustomID: "quest_complete_awakening", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F5E3\uFE0F"}},
			},
		})
	}

	components = append(components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "Back", Style: discordgo.SecondaryButton, CustomID: "back_village", Emoji: &discordgo.ComponentEmoji{Name: "\u25C0\uFE0F"}},
		},
	})

	return components
}

func travelEmbed(c *game.Character) *discordgo.MessageEmbed {
	var desc strings.Builder
	for _, areaID := range game.AreaOrder {
		area := game.Areas[areaID]
		accessible := game.CanAccessArea(area, c.StoryProgress)
		status := "\u2705"
		if !accessible {
			status = "\U0001F512"
		}
		if areaID == c.CurrentArea {
			status = "\U0001F4CD"
		}
		desc.WriteString(fmt.Sprintf("%s %s %s", status, area.Emoji, area.Name))
		if !accessible {
			desc.WriteString(" *(locked)*")
		}
		if area.MinLevel > 0 {
			desc.WriteString(fmt.Sprintf(" \u2014 Recommended Lv.%d+", area.MinLevel))
		}
		desc.WriteString("\n")
	}

	currentArea := game.Areas[c.CurrentArea]
	footer := c.CurrentArea
	if currentArea != nil {
		footer = currentArea.Name
	}

	return &discordgo.MessageEmbed{
		Title:       "\U0001F9ED Travel",
		Description: desc.String(),
		Color:       0x1ABC9C,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Current Location: %s", footer)},
	}
}

func travelComponents(c *game.Character) []discordgo.MessageComponent {
	var options []discordgo.SelectMenuOption
	for _, areaID := range game.AreaOrder {
		area := game.Areas[areaID]
		if !game.CanAccessArea(area, c.StoryProgress) {
			continue
		}
		label := area.Name
		if areaID == c.CurrentArea {
			label += " (current)"
		}
		desc := area.Description
		if len(desc) > 100 {
			desc = desc[:97] + "..."
		}
		options = append(options, discordgo.SelectMenuOption{
			Label:       label,
			Description: desc,
			Value:       areaID,
		})
	}

	components := []discordgo.MessageComponent{}
	if len(options) > 0 {
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "travel_select",
					Placeholder: "Choose destination",
					Options:     options,
				},
			},
		})
	}

	components = append(components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "Back", Style: discordgo.SecondaryButton, CustomID: "back_village", Emoji: &discordgo.ComponentEmoji{Name: "\u25C0\uFE0F"}},
		},
	})

	return components
}

func areaEmbed(c *game.Character, area *game.Area) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s %s", area.Emoji, area.Name),
		Description: fmt.Sprintf("%s\n\n**%s** \u2014 Level %d | HP: %s %d/%d", area.Description, c.Name, c.Level, makeBar(c.HP, c.MaxHP, 10), c.HP, c.MaxHP),
		Color:       0xE74C3C,
	}
}

func areaComponents(area *game.Area) []discordgo.MessageComponent {
	buttons := []discordgo.MessageComponent{
		discordgo.Button{Label: "Hunt Monsters", Style: discordgo.DangerButton, CustomID: "hunt", Emoji: &discordgo.ComponentEmoji{Name: "\u2694\uFE0F"}},
	}
	if area.BossID != "" {
		buttons = append(buttons, discordgo.Button{Label: "Challenge Boss", Style: discordgo.DangerButton, CustomID: "fight_boss", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F480"}})
	}
	if area.HasShop {
		buttons = append(buttons, discordgo.Button{Label: "Shop", Style: discordgo.PrimaryButton, CustomID: "shop", Emoji: &discordgo.ComponentEmoji{Name: "\U0001F6D2"}})
	}
	buttons = append(buttons, discordgo.Button{Label: "Back to Village", Style: discordgo.SecondaryButton, CustomID: "back_village", Emoji: &discordgo.ComponentEmoji{Name: "\u25C0\uFE0F"}})

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: buttons},
	}
}

func makeBar(current, max, length int) string {
	if max <= 0 {
		max = 1
	}
	filled := (current * length) / max
	if filled < 0 {
		filled = 0
	}
	if filled > length {
		filled = length
	}
	empty := length - filled
	return "[" + strings.Repeat("\u2588", filled) + strings.Repeat("\u2591", empty) + "]"
}
