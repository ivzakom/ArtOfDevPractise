package item

import "encoding/json"

type armor struct {
	item
	Defend int
	Hp     int
}

func (a *armor) Use() {
	a.Hp--
}

func (a *armor) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name   string `json:"name"`
		ID     int    `json:"id"`
		Hp     int    `json:"hp"`
		Type   string `json:"type"`
		Defend int    `json:"defend"`
	}{
		Name:   a.GetName(),
		ID:     a.GetID(),
		Hp:     a.Hp,
		Type:   a.GetType(),
		Defend: a.Defend,
	})
}

func (a *armor) UnmarshalJSON(data []byte) error {

	var temp struct {
		Name   string `json:"name"`
		ID     int    `json:"id"`
		Type   string `json:"type"`
		Hp     int    `json:"hp"`
		Defend int    `json:"defend"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	a.ID = temp.ID
	a.Name = temp.Name
	a.Type = temp.Type
	a.Hp = temp.Hp
	a.Defend = temp.Defend

	return nil

}
