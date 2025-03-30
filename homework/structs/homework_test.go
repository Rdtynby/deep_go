package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.nameLevelHouseGunFamilyType[0:42], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthRespectStrengthExperience = setBitsToInt(person.manaHealthRespectStrengthExperience, uint32(mana), 0, 10)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthRespectStrengthExperience = setBitsToInt(person.manaHealthRespectStrengthExperience, uint32(health), 10, 10)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthRespectStrengthExperience = setBitsToInt(person.manaHealthRespectStrengthExperience, uint32(respect), 20, 4)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthRespectStrengthExperience = setBitsToInt(person.manaHealthRespectStrengthExperience, uint32(strength), 24, 4)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthRespectStrengthExperience = setBitsToInt(person.manaHealthRespectStrengthExperience, uint32(experience), 28, 4)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.nameLevelHouseGunFamilyType[42] = uint8(level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.nameLevelHouseGunFamilyType[43] = setBitsToInt(person.nameLevelHouseGunFamilyType[43], 1, 0, 1)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.nameLevelHouseGunFamilyType[43] = setBitsToInt(person.nameLevelHouseGunFamilyType[43], 1, 1, 1)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.nameLevelHouseGunFamilyType[43] = setBitsToInt(person.nameLevelHouseGunFamilyType[43], 1, 2, 1)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.nameLevelHouseGunFamilyType[43] = setBitsToInt(person.nameLevelHouseGunFamilyType[43], uint8(personType), 3, 2)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type Number interface {
	uint8 | uint32
}

type GamePerson struct {
	x                                   int32
	y                                   int32
	z                                   int32
	gold                                int32
	manaHealthRespectStrengthExperience uint32
	nameLevelHouseGunFamilyType         [44]uint8
}

func getBitsFromInt[T Number](value T, start int, length int) T {
	mask := T(1<<length - 1)
	return value >> start & mask
}

func setBitsToInt[T Number](value T, bits T, start int, length int) T {
	mask := T(1<<length-1) << start
	return value & ^mask | bits<<start
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}

	for _, option := range options {
		option(&person)
	}

	return person
}

func (p *GamePerson) Name() string {
	return string(p.nameLevelHouseGunFamilyType[0:42])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(getBitsFromInt(p.manaHealthRespectStrengthExperience, 0, 10))
}

func (p *GamePerson) Health() int {
	return int(getBitsFromInt(p.manaHealthRespectStrengthExperience, 10, 10))
}

func (p *GamePerson) Respect() int {
	return int(getBitsFromInt(p.manaHealthRespectStrengthExperience, 20, 4))
}

func (p *GamePerson) Strength() int {
	return int(getBitsFromInt(p.manaHealthRespectStrengthExperience, 24, 4))
}

func (p *GamePerson) Experience() int {
	return int(getBitsFromInt(p.manaHealthRespectStrengthExperience, 28, 4))
}

func (p *GamePerson) Level() int {
	return int(p.nameLevelHouseGunFamilyType[42])
}

func (p *GamePerson) HasHouse() bool {
	return getBitsFromInt(p.nameLevelHouseGunFamilyType[43], 0, 1) > 0
}

func (p *GamePerson) HasGun() bool {
	return getBitsFromInt(p.nameLevelHouseGunFamilyType[43], 1, 1) > 0
}

func (p *GamePerson) HasFamily() bool {
	return getBitsFromInt(p.nameLevelHouseGunFamilyType[43], 2, 1) > 0
}

func (p *GamePerson) Type() int {
	return int(getBitsFromInt(p.nameLevelHouseGunFamilyType[43], 3, 2))
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32 - 1, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 999
	const health = 1000
	const respect = 7
	const strength = 8
	const experience = 9
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
