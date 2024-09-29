package item

import "encoding/json"

type weapon struct {
	item
	Damage int
}

func (w *weapon) Use() {
}

func (w *weapon) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name   string `json:"name"`
		ID     int    `json:"id"`
		Damage int    `json:"damage"`
		Type   string `json:"type"`
	}{
		Name:   w.GetName(),
		ID:     w.GetID(),
		Damage: w.Damage,
		Type:   w.GetType(),
	})
}

func (w *weapon) UnmarshalJSON(data []byte) error {

	var temp struct {
		Name   string `json:"name"`
		ID     int    `json:"id"`
		Damage int    `json:"damage"`
		Type   string `json:"type"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	w.ID = temp.ID
	w.Name = temp.Name
	w.Damage = temp.Damage
	w.Type = temp.Type

	return nil

}
