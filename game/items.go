package game

type ItemType string

const (
	ItemTypeWeapon     ItemType = "weapon"
	ItemTypeShield     ItemType = "shield"
	ItemTypeArmor      ItemType = "armor"
	ItemTypeConsumable ItemType = "consumable"
	ItemTypeQuest      ItemType = "quest"
)

type Item struct {
	ID          string
	Name        string
	Description string
	Type        ItemType
	Attack      int
	Defense     int
	Magic       int
	HPRestore   int
	MPRestore   int
	BuyPrice    int
	SellPrice   int
	LevelReq    int
}

var Items = map[string]*Item{
	// --- Weapons ---
	"wooden_sword": {
		ID: "wooden_sword", Name: "Wooden Sword", Description: "A basic training sword.",
		Type: ItemTypeWeapon, Attack: 5, BuyPrice: 0, SellPrice: 10,
	},
	"iron_sword": {
		ID: "iron_sword", Name: "Iron Sword", Description: "A sturdy iron blade.",
		Type: ItemTypeWeapon, Attack: 12, BuyPrice: 150, SellPrice: 75, LevelReq: 3,
	},
	"steel_blade": {
		ID: "steel_blade", Name: "Steel Blade", Description: "A finely forged steel sword.",
		Type: ItemTypeWeapon, Attack: 20, BuyPrice: 400, SellPrice: 200, LevelReq: 6,
	},
	"master_sword": {
		ID: "master_sword", Name: "Blade of Eldoria", Description: "The legendary blade that seals the darkness.",
		Type: ItemTypeWeapon, Attack: 35, Magic: 10, BuyPrice: 0, SellPrice: 0, LevelReq: 10,
	},
	"fire_rod": {
		ID: "fire_rod", Name: "Fire Rod", Description: "A magical rod that channels fire.",
		Type: ItemTypeWeapon, Attack: 8, Magic: 20, BuyPrice: 300, SellPrice: 150, LevelReq: 5,
	},

	// --- Shields ---
	"wooden_shield": {
		ID: "wooden_shield", Name: "Wooden Shield", Description: "A simple wooden shield.",
		Type: ItemTypeShield, Defense: 3, BuyPrice: 50, SellPrice: 25,
	},
	"iron_shield": {
		ID: "iron_shield", Name: "Iron Shield", Description: "A reliable iron shield.",
		Type: ItemTypeShield, Defense: 8, BuyPrice: 200, SellPrice: 100, LevelReq: 4,
	},
	"mirror_shield": {
		ID: "mirror_shield", Name: "Mirror Shield", Description: "A mystical shield that reflects magic.",
		Type: ItemTypeShield, Defense: 15, Magic: 5, BuyPrice: 500, SellPrice: 250, LevelReq: 8,
	},

	// --- Armor ---
	"cloth_tunic": {
		ID: "cloth_tunic", Name: "Cloth Tunic", Description: "Basic cloth clothing.",
		Type: ItemTypeArmor, Defense: 2, BuyPrice: 0, SellPrice: 5,
	},
	"leather_armor": {
		ID: "leather_armor", Name: "Leather Armor", Description: "Sturdy leather protection.",
		Type: ItemTypeArmor, Defense: 6, BuyPrice: 120, SellPrice: 60, LevelReq: 2,
	},
	"chain_mail": {
		ID: "chain_mail", Name: "Chain Mail", Description: "Interlocking metal rings for solid defense.",
		Type: ItemTypeArmor, Defense: 12, BuyPrice: 350, SellPrice: 175, LevelReq: 5,
	},
	"hero_tunic": {
		ID: "hero_tunic", Name: "Hero's Tunic", Description: "The legendary armor of the hero of Eldoria.",
		Type: ItemTypeArmor, Defense: 20, Magic: 5, BuyPrice: 0, SellPrice: 0, LevelReq: 10,
	},

	// --- Consumables ---
	"health_potion": {
		ID: "health_potion", Name: "Health Potion", Description: "Restores 30 HP.",
		Type: ItemTypeConsumable, HPRestore: 30, BuyPrice: 25, SellPrice: 12,
	},
	"mana_potion": {
		ID: "mana_potion", Name: "Mana Potion", Description: "Restores 15 MP.",
		Type: ItemTypeConsumable, MPRestore: 15, BuyPrice: 30, SellPrice: 15,
	},
	"elixir": {
		ID: "elixir", Name: "Elixir", Description: "Fully restores HP and MP.",
		Type: ItemTypeConsumable, HPRestore: 999, MPRestore: 999, BuyPrice: 200, SellPrice: 100,
	},

	// --- Quest Items ---
	"forest_key": {
		ID: "forest_key", Name: "Forest Temple Key", Description: "Opens the Forest Temple.",
		Type: ItemTypeQuest,
	},
	"crystal_heart": {
		ID: "crystal_heart", Name: "Crystal Heart", Description: "A glowing crystal from deep within the caves.",
		Type: ItemTypeQuest,
	},
	"frost_gem": {
		ID: "frost_gem", Name: "Frost Gem", Description: "A gem of eternal ice from the Frozen Peaks.",
		Type: ItemTypeQuest,
	},
	"shadow_key": {
		ID: "shadow_key", Name: "Shadow Key", Description: "Opens the gates of the Shadow Citadel.",
		Type: ItemTypeQuest,
	},
}

var ShopItems = map[string][]string{
	"village": {"health_potion", "mana_potion", "wooden_shield", "leather_armor"},
	"forest":  {"health_potion", "mana_potion", "iron_sword", "iron_shield"},
	"caves":   {"health_potion", "mana_potion", "elixir", "chain_mail", "fire_rod"},
	"peaks":   {"health_potion", "mana_potion", "elixir", "steel_blade", "mirror_shield"},
	"citadel": {"health_potion", "mana_potion", "elixir"},
}
