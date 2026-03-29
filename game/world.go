package game

import "math/rand"

type Area struct {
	ID          string
	Name        string
	Description string
	Emoji       string
	Enemies     []string
	BossID      string
	MinLevel    int
	UnlockStory int
	HasShop     bool
}

var Areas = map[string]*Area{
	"village": {
		ID:          "village",
		Name:        "Eldoria Village",
		Description: "A peaceful village nestled in a green valley. The villagers go about their daily lives, unaware of the growing darkness.",
		Emoji:       "\U0001F3E1",
		HasShop:     true,
	},
	"forest": {
		ID:          "forest",
		Name:        "Whispering Forest",
		Description: "An ancient forest shrouded in mist. Strange sounds echo between the twisted trees.",
		Emoji:       "\U0001F332",
		Enemies:     []string{"slime", "wolf", "bokoblin"},
		BossID:      "forest_boss",
		MinLevel:    1,
		UnlockStory: 1,
		HasShop:     true,
	},
	"caves": {
		ID:          "caves",
		Name:        "Crystal Caves",
		Description: "Deep underground caverns illuminated by glowing crystals. Danger lurks in every shadow.",
		Emoji:       "\U0001F48E",
		Enemies:     []string{"bat", "skeleton", "cave_troll"},
		BossID:      "cave_boss",
		MinLevel:    4,
		UnlockStory: 2,
		HasShop:     true,
	},
	"peaks": {
		ID:          "peaks",
		Name:        "Frozen Peaks",
		Description: "Treacherous mountain peaks blanketed in eternal snow. The air is thin and biting cold.",
		Emoji:       "\U0001F3D4\uFE0F",
		Enemies:     []string{"ice_wolf", "golem"},
		BossID:      "peaks_boss",
		MinLevel:    7,
		UnlockStory: 3,
		HasShop:     true,
	},
	"citadel": {
		ID:          "citadel",
		Name:        "Shadow Citadel",
		Description: "A dark fortress radiating evil energy. The final stronghold of the Dark Lord.",
		Emoji:       "\U0001F3F0",
		Enemies:     []string{"dark_knight", "shadow_mage"},
		BossID:      "dark_lord",
		MinLevel:    10,
		UnlockStory: 4,
		HasShop:     true,
	},
}

var AreaOrder = []string{"village", "forest", "caves", "peaks", "citadel"}

func GetRandomEnemy(areaID string) *Enemy {
	area, ok := Areas[areaID]
	if !ok || len(area.Enemies) == 0 {
		return nil
	}
	templateID := area.Enemies[rand.Intn(len(area.Enemies))]
	return SpawnEnemy(templateID)
}

func GetBoss(areaID string) *Enemy {
	area, ok := Areas[areaID]
	if !ok || area.BossID == "" {
		return nil
	}
	return SpawnEnemy(area.BossID)
}

func CanAccessArea(area *Area, storyProgress int) bool {
	return storyProgress >= area.UnlockStory
}
