package game

import (
	"fmt"
	"math/rand"
	"sync"
)

type CombatSession struct {
	Player          *Character
	Enemy           *Enemy
	PlayerDefending bool
	Log             []string
	Finished        bool
	Victory         bool
	LootDropped     []string
	XPGained        int
	GoldGained      int
	LeveledUp       bool
}

type CombatManager struct {
	mu       sync.RWMutex
	Sessions map[string]*CombatSession
}

func NewCombatManager() *CombatManager {
	return &CombatManager{
		Sessions: make(map[string]*CombatSession),
	}
}

func (cm *CombatManager) StartCombat(player *Character, enemy *Enemy) *CombatSession {
	session := &CombatSession{
		Player: player,
		Enemy:  enemy,
	}

	session.Log = append(session.Log, fmt.Sprintf("A wild **%s** appears!", enemy.Name))

	if player.Speed < enemy.Speed {
		session.Log = append(session.Log, fmt.Sprintf("The %s is faster and strikes first!", enemy.Name))
		session.enemyTurn()
	}

	cm.mu.Lock()
	cm.Sessions[player.UserID] = session
	cm.mu.Unlock()

	return session
}

func (cm *CombatManager) GetSession(userID string) *CombatSession {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.Sessions[userID]
}

func (cm *CombatManager) EndSession(userID string) {
	cm.mu.Lock()
	delete(cm.Sessions, userID)
	cm.mu.Unlock()
}

func (cs *CombatSession) PlayerAttack() {
	cs.Log = nil
	cs.PlayerDefending = false

	damage := calculateDamage(cs.Player.TotalAttack(), cs.Enemy.Defense)
	cs.Enemy.TakeDamage(damage)
	cs.Log = append(cs.Log, fmt.Sprintf("You strike the %s for **%d** damage!", cs.Enemy.Name, damage))

	if !cs.Enemy.IsAlive() {
		cs.victory()
		return
	}

	cs.enemyTurn()
}

func (cs *CombatSession) PlayerDefend() {
	cs.Log = nil
	cs.PlayerDefending = true
	cs.Log = append(cs.Log, "You raise your guard, bracing for impact!")

	cs.enemyTurn()
}

func (cs *CombatSession) PlayerMagic() {
	cs.Log = nil
	cs.PlayerDefending = false

	if cs.Player.MP < 5 {
		cs.Log = append(cs.Log, "Not enough MP! You need at least 5 MP to cast a spell.")
		return
	}

	cs.Player.MP -= 5
	damage := calculateDamage(cs.Player.TotalMagic()*2, cs.Enemy.Defense/2)
	cs.Enemy.TakeDamage(damage)
	cs.Log = append(cs.Log, fmt.Sprintf("You cast a spell on the %s for **%d** magic damage! (-5 MP)", cs.Enemy.Name, damage))

	if !cs.Enemy.IsAlive() {
		cs.victory()
		return
	}

	cs.enemyTurn()
}

func (cs *CombatSession) PlayerUseItem(itemID string) bool {
	cs.Log = nil
	cs.PlayerDefending = false

	item, ok := Items[itemID]
	if !ok || item.Type != ItemTypeConsumable {
		return false
	}

	if item.HPRestore > 0 {
		cs.Player.Heal(item.HPRestore)
		cs.Log = append(cs.Log, fmt.Sprintf("You use **%s** and restore HP! (HP: %d/%d)", item.Name, cs.Player.HP, cs.Player.MaxHP))
	}
	if item.MPRestore > 0 {
		cs.Player.RestoreMP(item.MPRestore)
		cs.Log = append(cs.Log, fmt.Sprintf("You use **%s** and restore MP! (MP: %d/%d)", item.Name, cs.Player.MP, cs.Player.MaxMP))
	}

	cs.enemyTurn()
	return true
}

func (cs *CombatSession) PlayerFlee() bool {
	cs.Log = nil
	cs.PlayerDefending = false

	if cs.Enemy.IsBoss {
		cs.Log = append(cs.Log, "You can't flee from a boss battle!")
		cs.enemyTurn()
		return false
	}

	fleeChance := float64(cs.Player.Speed)/float64(cs.Player.Speed+cs.Enemy.Speed) + 0.2
	if rand.Float64() < fleeChance {
		cs.Log = append(cs.Log, "You successfully fled from battle!")
		cs.Finished = true
		return true
	}

	cs.Log = append(cs.Log, "You failed to escape!")
	cs.enemyTurn()
	return false
}

func (cs *CombatSession) enemyTurn() {
	if cs.Finished {
		return
	}

	damage := calculateDamage(cs.Enemy.Attack, cs.Player.TotalDefense())
	if cs.PlayerDefending {
		damage = damage / 2
	}
	if damage < 1 {
		damage = 1
	}

	cs.Player.TakeDamage(damage)

	if cs.PlayerDefending {
		cs.Log = append(cs.Log, fmt.Sprintf("The %s attacks for **%d** damage! (reduced by guard)", cs.Enemy.Name, damage))
	} else {
		cs.Log = append(cs.Log, fmt.Sprintf("The %s attacks for **%d** damage!", cs.Enemy.Name, damage))
	}

	if !cs.Player.IsAlive() {
		cs.defeat()
	}
}

func (cs *CombatSession) victory() {
	cs.Finished = true
	cs.Victory = true
	cs.XPGained = cs.Enemy.XPReward
	cs.GoldGained = cs.Enemy.GoldReward
	cs.LootDropped = cs.Enemy.RollLoot()

	cs.Player.Gold += cs.GoldGained
	cs.LeveledUp = cs.Player.GainXP(cs.XPGained)

	cs.Log = append(cs.Log, fmt.Sprintf("\n**Victory!** You defeated the %s!", cs.Enemy.Name))
	cs.Log = append(cs.Log, fmt.Sprintf("Gained **%d XP** and **%d Gold**!", cs.XPGained, cs.GoldGained))

	if cs.LeveledUp {
		cs.Log = append(cs.Log, fmt.Sprintf("**LEVEL UP!** You are now level %d!", cs.Player.Level))
	}

	for _, itemID := range cs.LootDropped {
		if item, ok := Items[itemID]; ok {
			cs.Log = append(cs.Log, fmt.Sprintf("Obtained: **%s**!", item.Name))
		}
	}
}

func (cs *CombatSession) defeat() {
	cs.Finished = true
	cs.Victory = false
	cs.Log = append(cs.Log, "\n**Defeated!** You have fallen in battle...")
	cs.Log = append(cs.Log, "You wake up back in the village with half your gold gone.")

	cs.Player.Gold = cs.Player.Gold / 2
	cs.Player.HP = cs.Player.MaxHP / 2
	cs.Player.MP = cs.Player.MaxMP / 2
	cs.Player.CurrentArea = "village"
}

func calculateDamage(attack, defense int) int {
	base := attack - defense/2
	if base < 1 {
		base = 1
	}
	variance := rand.Intn(base/3+1) - base/6
	damage := base + variance
	if damage < 1 {
		damage = 1
	}
	return damage
}
