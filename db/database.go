package db

import (
	"database/sql"

	"discord-rpg-bot/game"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	d := &Database{db: db}
	if err := d.migrate(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS characters (
			user_id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			level INTEGER DEFAULT 1,
			xp INTEGER DEFAULT 0,
			xp_to_next INTEGER DEFAULT 100,
			hp INTEGER DEFAULT 50,
			max_hp INTEGER DEFAULT 50,
			mp INTEGER DEFAULT 20,
			max_mp INTEGER DEFAULT 20,
			attack INTEGER DEFAULT 10,
			defense INTEGER DEFAULT 5,
			speed INTEGER DEFAULT 5,
			magic INTEGER DEFAULT 5,
			gold INTEGER DEFAULT 50,
			current_area TEXT DEFAULT 'village',
			weapon TEXT DEFAULT 'wooden_sword',
			shield TEXT DEFAULT '',
			armor TEXT DEFAULT 'cloth_tunic',
			story_progress INTEGER DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS inventories (
			user_id TEXT NOT NULL,
			item_id TEXT NOT NULL,
			quantity INTEGER DEFAULT 1,
			PRIMARY KEY (user_id, item_id)
		)`,
		`CREATE TABLE IF NOT EXISTS quest_progress (
			user_id TEXT NOT NULL,
			quest_id TEXT NOT NULL,
			status TEXT DEFAULT 'available',
			PRIMARY KEY (user_id, quest_id)
		)`,
		`CREATE TABLE IF NOT EXISTS boss_kills (
			user_id TEXT NOT NULL,
			boss_id TEXT NOT NULL,
			PRIMARY KEY (user_id, boss_id)
		)`,
	}

	for _, q := range queries {
		if _, err := d.db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) CreateCharacter(userID, name string) (*game.Character, error) {
	c := game.NewCharacter(userID, name)
	_, err := d.db.Exec(
		`INSERT INTO characters (user_id, name, level, xp, xp_to_next, hp, max_hp, mp, max_mp, attack, defense, speed, magic, gold, current_area, weapon, shield, armor, story_progress)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.UserID, c.Name, c.Level, c.XP, c.XPToNext, c.HP, c.MaxHP, c.MP, c.MaxMP,
		c.Attack, c.Defense, c.Speed, c.Magic, c.Gold, c.CurrentArea,
		c.Weapon, c.Shield, c.Armor, c.StoryProgress,
	)
	if err != nil {
		return nil, err
	}

	// Starting items
	d.AddItem(userID, "health_potion", 3)

	return c, nil
}

func (d *Database) GetCharacter(userID string) (*game.Character, error) {
	c := &game.Character{}
	err := d.db.QueryRow(
		`SELECT user_id, name, level, xp, xp_to_next, hp, max_hp, mp, max_mp, attack, defense, speed, magic, gold, current_area, weapon, shield, armor, story_progress
		 FROM characters WHERE user_id = ?`, userID,
	).Scan(
		&c.UserID, &c.Name, &c.Level, &c.XP, &c.XPToNext, &c.HP, &c.MaxHP, &c.MP, &c.MaxMP,
		&c.Attack, &c.Defense, &c.Speed, &c.Magic, &c.Gold, &c.CurrentArea,
		&c.Weapon, &c.Shield, &c.Armor, &c.StoryProgress,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (d *Database) SaveCharacter(c *game.Character) error {
	_, err := d.db.Exec(
		`UPDATE characters SET name=?, level=?, xp=?, xp_to_next=?, hp=?, max_hp=?, mp=?, max_mp=?, attack=?, defense=?, speed=?, magic=?, gold=?, current_area=?, weapon=?, shield=?, armor=?, story_progress=?
		 WHERE user_id=?`,
		c.Name, c.Level, c.XP, c.XPToNext, c.HP, c.MaxHP, c.MP, c.MaxMP,
		c.Attack, c.Defense, c.Speed, c.Magic, c.Gold, c.CurrentArea,
		c.Weapon, c.Shield, c.Armor, c.StoryProgress, c.UserID,
	)
	return err
}

func (d *Database) GetInventory(userID string) ([]game.InventorySlot, error) {
	rows, err := d.db.Query(`SELECT item_id, quantity FROM inventories WHERE user_id = ? AND quantity > 0`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inv []game.InventorySlot
	for rows.Next() {
		var slot game.InventorySlot
		if err := rows.Scan(&slot.ItemID, &slot.Quantity); err != nil {
			return nil, err
		}
		inv = append(inv, slot)
	}
	return inv, nil
}

func (d *Database) AddItem(userID, itemID string, qty int) error {
	_, err := d.db.Exec(
		`INSERT INTO inventories (user_id, item_id, quantity) VALUES (?, ?, ?)
		 ON CONFLICT(user_id, item_id) DO UPDATE SET quantity = quantity + ?`,
		userID, itemID, qty, qty,
	)
	return err
}

func (d *Database) RemoveItem(userID, itemID string, qty int) error {
	_, err := d.db.Exec(
		`UPDATE inventories SET quantity = quantity - ? WHERE user_id = ? AND item_id = ?`,
		qty, userID, itemID,
	)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(
		`DELETE FROM inventories WHERE user_id = ? AND item_id = ? AND quantity <= 0`,
		userID, itemID,
	)
	return err
}

func (d *Database) HasItem(userID, itemID string) (bool, error) {
	var count int
	err := d.db.QueryRow(
		`SELECT COALESCE(SUM(quantity), 0) FROM inventories WHERE user_id = ? AND item_id = ?`,
		userID, itemID,
	).Scan(&count)
	return count > 0, err
}

func (d *Database) RecordBossKill(userID, bossID string) error {
	_, err := d.db.Exec(
		`INSERT OR IGNORE INTO boss_kills (user_id, boss_id) VALUES (?, ?)`,
		userID, bossID,
	)
	return err
}

func (d *Database) HasKilledBoss(userID, bossID string) (bool, error) {
	var count int
	err := d.db.QueryRow(
		`SELECT COUNT(*) FROM boss_kills WHERE user_id = ? AND boss_id = ?`,
		userID, bossID,
	).Scan(&count)
	return count > 0, err
}
