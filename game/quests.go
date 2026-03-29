package game

type QuestStatus string

const (
	QuestAvailable QuestStatus = "available"
	QuestActive    QuestStatus = "active"
	QuestCompleted QuestStatus = "completed"
)

type Quest struct {
	ID           string
	Name         string
	Description  string
	Area         string
	StoryIndex   int
	XPReward     int
	GoldReward   int
	ItemReward   string
	RequiresBoss string
}

var Quests = map[string]*Quest{
	"awakening": {
		ID:          "awakening",
		Name:        "The Awakening",
		Description: "The village Elder senses a great darkness approaching. Speak with him to begin your journey.",
		Area:        "village",
		StoryIndex:  1,
		XPReward:    50,
		GoldReward:  30,
	},
	"forest_shadows": {
		ID:           "forest_shadows",
		Name:         "Shadows in the Forest",
		Description:  "Dark creatures have infested the Whispering Forest. Defeat the Great Deku Blight to restore peace.",
		Area:         "forest",
		StoryIndex:   2,
		XPReward:     200,
		GoldReward:   150,
		ItemReward:   "iron_sword",
		RequiresBoss: "forest_boss",
	},
	"crystal_heart": {
		ID:           "crystal_heart",
		Name:         "Heart of the Mountain",
		Description:  "A powerful crystal lies deep within the caves. Defeat the Crystal Golem and retrieve the Crystal Heart.",
		Area:         "caves",
		StoryIndex:   3,
		XPReward:     400,
		GoldReward:   300,
		ItemReward:   "chain_mail",
		RequiresBoss: "cave_boss",
	},
	"frozen_summit": {
		ID:           "frozen_summit",
		Name:         "The Frozen Summit",
		Description:  "A Frost Drake terrorizes the mountain peaks. Slay the beast and claim the Frost Gem.",
		Area:         "peaks",
		StoryIndex:   4,
		XPReward:     600,
		GoldReward:   500,
		ItemReward:   "mirror_shield",
		RequiresBoss: "peaks_boss",
	},
	"final_shadow": {
		ID:           "final_shadow",
		Name:         "The Final Shadow",
		Description:  "Storm the Shadow Citadel and defeat Ganondrath, the Dark Lord, to save Eldoria!",
		Area:         "citadel",
		StoryIndex:   5,
		XPReward:     1000,
		GoldReward:   1000,
		ItemReward:   "master_sword",
		RequiresBoss: "dark_lord",
	},
}

var QuestOrder = []string{"awakening", "forest_shadows", "crystal_heart", "frozen_summit", "final_shadow"}

func GetCurrentQuest(storyProgress int) *Quest {
	for _, qID := range QuestOrder {
		q := Quests[qID]
		if q.StoryIndex == storyProgress+1 {
			return q
		}
	}
	return nil
}
