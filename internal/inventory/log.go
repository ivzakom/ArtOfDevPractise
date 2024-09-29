package inventory

import (
	"fmt"
	"golang_lessons/internal/item"
)

func (inv *Inventory) logOperation(logID int, item item.Item, quantity int, isAdd bool, err error) {

	textOperation := "remove"
	if isAdd {
		textOperation = "added"
	}

	textComment := ""
	errDisc := ""
	if err != nil {
		textComment = "Отказ в операции: "
		errDisc = fmt.Sprintf("\n %s", err.Error())
	}

	textLog := fmt.Sprintf("%s%d, %s has been %s to the inventory%s", textComment, quantity, item.GetName(), textOperation, errDisc)
	inv.Log[logID] = textLog

}
