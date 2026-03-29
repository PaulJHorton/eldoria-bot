package game

import "fmt"

type Character struct {
	UserID        string
	Name          string
	Level         int
	XP            int
	XPToNext      int
	HP            int
	MaxHP         int
	MP            int
	MaxMP         int
	Attack        int
	Defense       int
	Speed         int
	Magic         int
	Gold          int
	CurrentArea   string
	Weapon        string
	Shield        string
	Armor         string
	StoryProgress int
}

type InventorySlot struct {
	ItemID   string
	Quantity int
}

func NewCharacter(userID, name string) *Character {
	return &Character{
		UserID:        userID,
		Name:          name,
		Level:         1,
		XP:            0,
		XPToNext:      100,
		HP:            50,
		MaxHP:         50,
		MP:            20,
		MaxMP:         20,
		Attack:        10,
		Defense:       5,
		Speed:         5,
		Magic:         5,
		Gold:          50,
		CurrentArea:   "village",
		Weapon:        "wooden_sword",
		Shield:        "",
		Armor:         "cloth_tunic",
		StoryProgress: 0,
	}
}

func (c *Character) GainXP(amount int) (leveled bool) {
	c.XP += amount
	for c.XP >= c.XPToNext {
		c.XP -= c.XPToNext
		c.LevelUp()
		leveled = true
	}
	return
}

func (c *Character) LevelUp() {
	c.Level++
	c.XPToNext = c.Level * 100
	c.MaxHP += 10
	c.MaxMP += 5
	c.Attack += 3
	c.Defense += 2
	c.Speed += 2
	c.Magic += 2
	c.HP = c.MaxHP
	c.MP = c.MaxMP
}

func (c *Character) TotalAttack() int {
	total := c.Attack
	if item, ok := Items[c.Weapon]; ok {
		total += item.Attack
	}
	return total
}

func (c *Character) TotalDefense() int {
	total := c.Defense
	if item, ok := Items[c.Shield]; ok {
		total += item.Defense
	}
	if item, ok := Items[c.Armor]; ok {
		total += item.Defense
	}
	return total
}

func (c *Character) TotalMagic() int {
	total := c.Magic
	if item, ok := Items[c.Weapon]; ok {
		total += item.Magic
	}
	return total
}

func (c *Character) TakeDamage(amount int) {
	c.HP -= amount
	if c.HP < 0 {
		c.HP = 0
	}
}

func (c *Character) Heal(amount int) {
	c.HP += amount
	if c.HP > c.MaxHP {
		c.HP = c.MaxHP
	}
}

func (c *Character) RestoreMP(amount int) {
	c.MP += amount
	if c.MP > c.MaxMP {
		c.MP = c.MaxMP
	}
}

func (c *Character) IsAlive() bool {
	return c.HP > 0
}

func (c *Character) StatsString() string {
	return fmt.Sprintf("Level %d | HP: %d/%d | MP: %d/%d\nATK: %d | DEF: %d | SPD: %d | MAG: %d\nGold: %d | XP: %d/%d",
		c.Level, c.HP, c.MaxHP, c.MP, c.MaxMP,
		c.TotalAttack(), c.TotalDefense(), c.Speed, c.TotalMagic(),
		c.Gold, c.XP, c.XPToNext)
}
