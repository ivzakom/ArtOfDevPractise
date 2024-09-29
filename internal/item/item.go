// internal/item/item.go
package item

import "fmt"

type item struct {
	ID   int
	Name string
	Type string
}

// Item - интерфейс для всех предметов
type Item interface {
	GetName() string
	GetID() int
	Use()
	GetType() string
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

func (i *item) GetName() string {
	return i.Name
}

func (i *item) GetID() int {
	return i.ID
}

func (i *item) GetType() string {
	return i.Type
}

type ErrorItemDoesNotExist struct {
	Message string
	Err     error
}

func (e *ErrorItemDoesNotExist) Error() string {
	return e.Message
}

func (e *ErrorItemDoesNotExist) Unwrap() error {
	return e.Err
}

// Типы предметов для фабричной функции
const (
	WeaponType = "weapon"
	ArmorType  = "armor"
)

// ItemConfig - структура для передачи параметров при создании предмета
type ItemConfig struct {
	ID     int
	Name   string
	Damage int // для оружия
	Defend int // для брони
	Hp     int // для брони
	Type   string
}

// NewItem - фабричная функция для создания предметов
func NewItem(itemType string, config ItemConfig) (Item, error) {
	switch itemType {
	case WeaponType:
		return &weapon{
			item: item{
				config.ID,
				config.Name,
				itemType,
			},
			Damage: config.Damage,
		}, nil
	case ArmorType:
		return &armor{
			item: item{
				config.ID,
				config.Name,
				itemType,
			},
			Defend: config.Defend,
			Hp:     config.Hp,
		}, nil
	default:
		return nil, fmt.Errorf("unknown item type: %s", itemType)
	}
}

func GetItemConfig(item Item) ItemConfig {

	сonfig := ItemConfig{}

	switch i := item.(type) {
	case *armor:
		сonfig.ID = i.GetID()
		сonfig.Name = i.GetName()
		сonfig.Type = i.GetType()
		сonfig.Hp = i.Hp
		сonfig.Defend = i.Defend
	case *weapon:
		сonfig.ID = i.GetID()
		сonfig.Name = i.GetName()
		сonfig.Type = i.GetType()
		сonfig.Damage = i.Damage
	}

	return сonfig

}
