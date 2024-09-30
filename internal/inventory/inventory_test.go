package inventory

import (
	fl "artOfDevPractise/internal/filters"
	. "artOfDevPractise/internal/item"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAddItems(t *testing.T) {

	inv := NewInventory()
	filters := []func(item Item, quantity int) error{
		fl.FilterNotNull,
	}

	shild, _ := NewItem(WeaponType, ItemConfig{
		ID:     2,
		Name:   "shild",
		Defend: 10,
		Hp:     50,
	})

	_ = inv.AddItems(filters, []Item{shild}, 10)

	Convey("Given an inventory", t, func() {
		//So(inv.Items[shild], ShouldEqual, 10)
		//
		//filterErr := fl.NewFilterError(nil, "")
		//
		//filterErr = inv.RemoveItems(filters, []Item{shild}, 5)
		//So(filterErr, ShouldBeNil)
		//
		//err = inv.RemoveItems(filters, []Item{shild}, 5)
		//So(err, ShouldBeNil)
		//
		//err = inv.RemoveItems(filters, []Item{shild}, 50)
		//So(err, ShouldNotBeNil)

	})
}
