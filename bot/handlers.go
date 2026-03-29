package bot

import (
	"fmt"
	"strings"

	"discord-rpg-bot/db"
	"discord-rpg-bot/game"

	"github.com/bwmarrin/discordgo"
)

type Handlers struct {
	DB            *db.Database
	CombatManager *game.CombatManager
}

func (h *Handlers) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		h.handleCommand(s, i)
	case discordgo.InteractionMessageComponent:
		h.handleComponent(s, i)
	}
}

func getUserID(i *discordgo.InteractionCreate) string {
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID
	}
	if i.User != nil {
		return i.User.ID
	}
	return ""
}

func getUsername(i *discordgo.InteractionCreate) string {
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.Username
	}
	if i.User != nil {
		return i.User.Username
	}
	return "Hero"
}

// --- Slash Command Handlers ---

func (h *Handlers) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := getUserID(i)

	switch i.ApplicationCommandData().Name {
	case "adventure":
		h.cmdAdventure(s, i, userID)
	case "stats":
		h.cmdStats(s, i, userID)
	case "inventory":
		h.cmdInventory(s, i, userID)
	case "quests":
		h.cmdQuests(s, i, userID)
	}
}

func (h *Handlers) cmdAdventure(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, err := h.DB.GetCharacter(userID)
	if err != nil {
		respondMsg(s, i, "An error occurred. Please try again.")
		return
	}

	if char == nil {
		respondEmbed(s, i, welcomeEmbed(), welcomeComponents())
		return
	}

	if char.CurrentArea != "village" {
		if area, ok := game.Areas[char.CurrentArea]; ok {
			respondEmbed(s, i, areaEmbed(char, area), areaComponents(area))
			return
		}
	}
	respondEmbed(s, i, villageEmbed(char), villageComponents())
}

func (h *Handlers) cmdStats(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		respondMsg(s, i, "You haven't started your adventure yet! Use `/adventure` to begin.")
		return
	}
	respondEmbed(s, i, statsEmbed(char), backButton())
}

func (h *Handlers) cmdInventory(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		respondMsg(s, i, "You haven't started your adventure yet! Use `/adventure` to begin.")
		return
	}
	inv, _ := h.DB.GetInventory(userID)
	respondEmbed(s, i, inventoryEmbed(char, inv), inventoryComponents(inv))
}

func (h *Handlers) cmdQuests(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		respondMsg(s, i, "You haven't started your adventure yet! Use `/adventure` to begin.")
		return
	}
	respondEmbed(s, i, questsEmbed(char), questsComponents(char))
}

// --- Component (Button/Menu) Handlers ---

func (h *Handlers) handleComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := getUserID(i)
	customID := i.MessageComponentData().CustomID
	values := i.MessageComponentData().Values

	switch {
	case customID == "create_character":
		h.createCharacter(s, i, userID)
	case customID == "back_village":
		h.backToVillage(s, i, userID)
	case customID == "explore":
		h.explore(s, i, userID)
	case customID == "hunt":
		h.hunt(s, i, userID)
	case customID == "fight_boss":
		h.fightBoss(s, i, userID)
	case customID == "shop":
		h.openShop(s, i, userID)
	case customID == "shop_buy" && len(values) > 0:
		h.buyItem(s, i, userID, values[0])
	case customID == "inventory":
		h.showInventory(s, i, userID)
	case customID == "inv_select" && len(values) > 0:
		h.useOrEquipItem(s, i, userID, values[0])
	case customID == "quests":
		h.showQuests(s, i, userID)
	case customID == "quest_complete_awakening":
		h.completeAwakening(s, i, userID)
	case customID == "stats":
		h.showStats(s, i, userID)
	case customID == "rest":
		h.rest(s, i, userID)
	case customID == "travel":
		h.showTravel(s, i, userID)
	case customID == "travel_select" && len(values) > 0:
		h.travelTo(s, i, userID, values[0])
	case strings.HasPrefix(customID, "combat_"):
		h.handleCombat(s, i, userID, customID)
	case customID == "continue_after_combat":
		h.backToVillage(s, i, userID)
	}
}

func (h *Handlers) createCharacter(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	username := getUsername(i)

	char, err := h.DB.CreateCharacter(userID, username)
	if err != nil {
		updateMsg(s, i, "Failed to create character. You may already have one! Try `/adventure`.")
		return
	}

	updateEmbed(s, i, villageEmbed(char), villageComponents())
}

func (h *Handlers) backToVillage(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}

	h.CombatManager.EndSession(userID)
	char.CurrentArea = "village"
	h.DB.SaveCharacter(char)

	updateEmbed(s, i, villageEmbed(char), villageComponents())
}

func (h *Handlers) explore(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}

	if char.CurrentArea == "village" {
		h.showTravel(s, i, userID)
		return
	}

	area := game.Areas[char.CurrentArea]
	if area == nil {
		return
	}

	updateEmbed(s, i, areaEmbed(char, area), areaComponents(area))
}

func (h *Handlers) hunt(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}

	enemy := game.GetRandomEnemy(char.CurrentArea)
	if enemy == nil {
		updateMsg(s, i, "No enemies found here.")
		return
	}

	cs := h.CombatManager.StartCombat(char, enemy)
	updateEmbed(s, i, combatEmbed(cs), combatComponents(cs))
}

func (h *Handlers) fightBoss(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}

	area := game.Areas[char.CurrentArea]
	if area == nil || area.BossID == "" {
		return
	}

	killed, _ := h.DB.HasKilledBoss(userID, area.BossID)
	if killed {
		updateMsg(s, i, "You have already defeated this area's boss!")
		return
	}

	boss := game.GetBoss(char.CurrentArea)
	if boss == nil {
		return
	}

	cs := h.CombatManager.StartCombat(char, boss)
	updateEmbed(s, i, combatEmbed(cs), combatComponents(cs))
}

func (h *Handlers) handleCombat(s *discordgo.Session, i *discordgo.InteractionCreate, userID, action string) {
	cs := h.CombatManager.GetSession(userID)
	if cs == nil {
		updateMsg(s, i, "No active combat session. Use `/adventure` to continue.")
		return
	}

	switch action {
	case "combat_attack":
		cs.PlayerAttack()
	case "combat_defend":
		cs.PlayerDefend()
	case "combat_magic":
		cs.PlayerMagic()
	case "combat_item":
		inv, _ := h.DB.GetInventory(userID)
		used := false
		for _, slot := range inv {
			item := game.Items[slot.ItemID]
			if item != nil && item.Type == game.ItemTypeConsumable {
				if cs.PlayerUseItem(slot.ItemID) {
					h.DB.RemoveItem(userID, slot.ItemID, 1)
					used = true
					break
				}
			}
		}
		if !used {
			cs.Log = []string{"No consumable items in your inventory!"}
		}
	case "combat_flee":
		cs.PlayerFlee()
	}

	h.DB.SaveCharacter(cs.Player)

	if cs.Finished {
		if cs.Victory {
			for _, itemID := range cs.LootDropped {
				h.DB.AddItem(userID, itemID, 1)
			}
			if cs.Enemy.IsBoss {
				h.DB.RecordBossKill(userID, cs.Enemy.ID)
				h.checkQuestCompletion(cs.Player, cs.Enemy.ID)
				h.DB.SaveCharacter(cs.Player)
			}
		}
		h.CombatManager.EndSession(userID)
	}

	updateEmbed(s, i, combatEmbed(cs), combatComponents(cs))
}

func (h *Handlers) checkQuestCompletion(char *game.Character, bossID string) {
	currentQuest := game.GetCurrentQuest(char.StoryProgress)
	if currentQuest == nil {
		return
	}
	if currentQuest.RequiresBoss == bossID {
		char.StoryProgress = currentQuest.StoryIndex
		char.GainXP(currentQuest.XPReward)
		char.Gold += currentQuest.GoldReward
		if currentQuest.ItemReward != "" {
			h.DB.AddItem(char.UserID, currentQuest.ItemReward, 1)
		}
	}
}

func (h *Handlers) openShop(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}
	updateEmbed(s, i, shopEmbed(char, char.CurrentArea), shopComponents(char.CurrentArea))
}

func (h *Handlers) buyItem(s *discordgo.Session, i *discordgo.InteractionCreate, userID, itemID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}

	item, ok := game.Items[itemID]
	if !ok {
		return
	}

	if char.Gold < item.BuyPrice {
		updateMsg(s, i, fmt.Sprintf("Not enough gold! You need **%d Gold** but only have **%d**.", item.BuyPrice, char.Gold))
		return
	}

	char.Gold -= item.BuyPrice
	h.DB.SaveCharacter(char)
	h.DB.AddItem(userID, itemID, 1)

	updateEmbed(s, i, shopEmbed(char, char.CurrentArea), shopComponents(char.CurrentArea))
}

func (h *Handlers) showInventory(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}
	inv, _ := h.DB.GetInventory(userID)
	updateEmbed(s, i, inventoryEmbed(char, inv), inventoryComponents(inv))
}

func (h *Handlers) useOrEquipItem(s *discordgo.Session, i *discordgo.InteractionCreate, userID, itemID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}

	item, ok := game.Items[itemID]
	if !ok {
		return
	}

	switch item.Type {
	case game.ItemTypeWeapon:
		char.Weapon = itemID
		h.DB.SaveCharacter(char)
	case game.ItemTypeShield:
		char.Shield = itemID
		h.DB.SaveCharacter(char)
	case game.ItemTypeArmor:
		char.Armor = itemID
		h.DB.SaveCharacter(char)
	case game.ItemTypeConsumable:
		if item.HPRestore > 0 {
			char.Heal(item.HPRestore)
		}
		if item.MPRestore > 0 {
			char.RestoreMP(item.MPRestore)
		}
		h.DB.SaveCharacter(char)
		h.DB.RemoveItem(userID, itemID, 1)
	}

	inv, _ := h.DB.GetInventory(userID)
	updateEmbed(s, i, inventoryEmbed(char, inv), inventoryComponents(inv))
}

func (h *Handlers) showQuests(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}
	updateEmbed(s, i, questsEmbed(char), questsComponents(char))
}

func (h *Handlers) completeAwakening(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil || char.StoryProgress != 0 {
		return
	}

	quest := game.Quests["awakening"]
	char.StoryProgress = quest.StoryIndex
	char.GainXP(quest.XPReward)
	char.Gold += quest.GoldReward
	h.DB.SaveCharacter(char)

	embed := &discordgo.MessageEmbed{
		Title:       "\U0001F4DC Quest Complete: The Awakening",
		Description: "The Elder speaks:\n\n*\"A great darkness stirs in the Shadow Citadel. You must journey through the lands, grow stronger, and stop it before all is lost. Take this gold and begin your journey through the Whispering Forest...\"*",
		Color:       0x2ECC71,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Rewards", Value: fmt.Sprintf("+%d XP | +%d Gold\n\n**The Whispering Forest is now unlocked!**", quest.XPReward, quest.GoldReward)},
		},
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "Continue", Style: discordgo.SuccessButton, CustomID: "back_village", Emoji: &discordgo.ComponentEmoji{Name: "\u25B6\uFE0F"}},
			},
		},
	}

	updateEmbed(s, i, embed, components)
}

func (h *Handlers) showStats(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}
	updateEmbed(s, i, statsEmbed(char), backButton())
}

func (h *Handlers) rest(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}

	if char.Gold < 10 {
		updateMsg(s, i, "You don't have enough gold to rest! (Need 10 Gold)")
		return
	}

	char.Gold -= 10
	char.HP = char.MaxHP
	char.MP = char.MaxMP
	h.DB.SaveCharacter(char)

	updateEmbed(s, i, villageEmbed(char), villageComponents())
}

func (h *Handlers) showTravel(s *discordgo.Session, i *discordgo.InteractionCreate, userID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}
	updateEmbed(s, i, travelEmbed(char), travelComponents(char))
}

func (h *Handlers) travelTo(s *discordgo.Session, i *discordgo.InteractionCreate, userID, areaID string) {
	char, _ := h.DB.GetCharacter(userID)
	if char == nil {
		return
	}

	area, ok := game.Areas[areaID]
	if !ok || !game.CanAccessArea(area, char.StoryProgress) {
		return
	}

	char.CurrentArea = areaID
	h.DB.SaveCharacter(char)

	if areaID == "village" {
		updateEmbed(s, i, villageEmbed(char), villageComponents())
	} else {
		updateEmbed(s, i, areaEmbed(char, area), areaComponents(area))
	}
}

// --- Response Helpers ---

func respondMsg(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func respondEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, components []discordgo.MessageComponent) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
}

func updateMsg(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func updateEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, components []discordgo.MessageComponent) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}
