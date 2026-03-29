package game

import "math/rand"

type Enemy struct {
	ID         string
	Name       string
	HP         int
	MaxHP      int
	Attack     int
	Defense    int
	Speed      int
	Magic      int
	XPReward   int
	GoldReward int
	LootTable  []LootDrop
	IsBoss     bool
}

type LootDrop struct {
	ItemID string
	Chance float64
}

type EnemyTemplate struct {
	ID         string
	Name       string
	HP         int
	Attack     int
	Defense    int
	Speed      int
	Magic      int
	XPReward   int
	GoldReward int
	LootTable  []LootDrop
	IsBoss     bool
}

var EnemyTemplates = map[string]*EnemyTemplate{
	// --- Whispering Forest ---
	"slime": {
		ID: "slime", Name: "Green Slime", HP: 20, Attack: 6, Defense: 2, Speed: 3,
		XPReward: 15, GoldReward: 8,
		LootTable: []LootDrop{{ItemID: "health_potion", Chance: 0.3}},
	},
	"wolf": {
		ID: "wolf", Name: "Shadow Wolf", HP: 35, Attack: 12, Defense: 4, Speed: 8,
		XPReward: 25, GoldReward: 15,
		LootTable: []LootDrop{{ItemID: "health_potion", Chance: 0.2}},
	},
	"bokoblin": {
		ID: "bokoblin", Name: "Forest Bokoblin", HP: 45, Attack: 14, Defense: 6, Speed: 5,
		XPReward: 35, GoldReward: 25,
		LootTable: []LootDrop{
			{ItemID: "iron_sword", Chance: 0.1},
			{ItemID: "health_potion", Chance: 0.3},
		},
	},
	"forest_boss": {
		ID: "forest_boss", Name: "Great Deku Blight", HP: 120, Attack: 18, Defense: 8, Speed: 6, Magic: 10,
		XPReward: 150, GoldReward: 100, IsBoss: true,
		LootTable: []LootDrop{{ItemID: "forest_key", Chance: 1.0}},
	},

	// --- Crystal Caves ---
	"bat": {
		ID: "bat", Name: "Cave Bat", HP: 25, Attack: 10, Defense: 3, Speed: 10,
		XPReward: 20, GoldReward: 12,
		LootTable: []LootDrop{{ItemID: "mana_potion", Chance: 0.25}},
	},
	"skeleton": {
		ID: "skeleton", Name: "Skeleton Warrior", HP: 50, Attack: 16, Defense: 10, Speed: 4,
		XPReward: 40, GoldReward: 30,
		LootTable: []LootDrop{{ItemID: "iron_shield", Chance: 0.1}},
	},
	"cave_troll": {
		ID: "cave_troll", Name: "Crystal Troll", HP: 80, Attack: 20, Defense: 12, Speed: 3,
		XPReward: 60, GoldReward: 45,
		LootTable: []LootDrop{{ItemID: "elixir", Chance: 0.15}},
	},
	"cave_boss": {
		ID: "cave_boss", Name: "Crystal Golem", HP: 200, Attack: 25, Defense: 15, Speed: 4, Magic: 5,
		XPReward: 300, GoldReward: 200, IsBoss: true,
		LootTable: []LootDrop{{ItemID: "crystal_heart", Chance: 1.0}},
	},

	// --- Frozen Peaks ---
	"ice_wolf": {
		ID: "ice_wolf", Name: "Frost Wolf", HP: 60, Attack: 22, Defense: 10, Speed: 9,
		XPReward: 55, GoldReward: 40,
		LootTable: []LootDrop{{ItemID: "mana_potion", Chance: 0.3}},
	},
	"golem": {
		ID: "golem", Name: "Ice Golem", HP: 100, Attack: 28, Defense: 18, Speed: 2,
		XPReward: 80, GoldReward: 60,
		LootTable: []LootDrop{{ItemID: "elixir", Chance: 0.2}},
	},
	"peaks_boss": {
		ID: "peaks_boss", Name: "Frost Drake", HP: 300, Attack: 32, Defense: 20, Speed: 7, Magic: 15,
		XPReward: 500, GoldReward: 350, IsBoss: true,
		LootTable: []LootDrop{{ItemID: "frost_gem", Chance: 1.0}},
	},

	// --- Shadow Citadel ---
	"dark_knight": {
		ID: "dark_knight", Name: "Dark Knight", HP: 90, Attack: 30, Defense: 16, Speed: 6,
		XPReward: 100, GoldReward: 75,
		LootTable: []LootDrop{
			{ItemID: "steel_blade", Chance: 0.1},
			{ItemID: "elixir", Chance: 0.25},
		},
	},
	"shadow_mage": {
		ID: "shadow_mage", Name: "Shadow Mage", HP: 70, Attack: 20, Defense: 8, Speed: 8, Magic: 25,
		XPReward: 90, GoldReward: 70,
		LootTable: []LootDrop{{ItemID: "mana_potion", Chance: 0.4}},
	},
	"dark_lord": {
		ID: "dark_lord", Name: "Ganondrath, the Dark Lord", HP: 500, Attack: 40, Defense: 25, Speed: 10, Magic: 30,
		XPReward: 1000, GoldReward: 500, IsBoss: true,
		LootTable: []LootDrop{{ItemID: "shadow_key", Chance: 1.0}},
	},
}

func SpawnEnemy(templateID string) *Enemy {
	t, ok := EnemyTemplates[templateID]
	if !ok {
		return nil
	}
	return &Enemy{
		ID:         t.ID,
		Name:       t.Name,
		HP:         t.HP,
		MaxHP:      t.HP,
		Attack:     t.Attack,
		Defense:    t.Defense,
		Speed:      t.Speed,
		Magic:      t.Magic,
		XPReward:   t.XPReward,
		GoldReward: t.GoldReward,
		LootTable:  t.LootTable,
		IsBoss:     t.IsBoss,
	}
}

func (e *Enemy) TakeDamage(amount int) {
	e.HP -= amount
	if e.HP < 0 {
		e.HP = 0
	}
}

func (e *Enemy) IsAlive() bool {
	return e.HP > 0
}

func (e *Enemy) RollLoot() []string {
	var loot []string
	for _, drop := range e.LootTable {
		if rand.Float64() < drop.Chance {
			loot = append(loot, drop.ItemID)
		}
	}
	return loot
}
