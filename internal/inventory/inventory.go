package inventory

import (
	flt "artOfDevPractise/internal/filters"
	itm "artOfDevPractise/internal/item"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

type Inventory struct {
	Log   map[int]string
	mu    sync.Mutex
	Items map[itm.Item]int
}

func NewInventory() *Inventory {
	inv := &Inventory{
		Items: make(map[itm.Item]int),
		Log:   make(map[int]string),
	}
	return inv
}

func (inv *Inventory) ShowInventory(ascending bool) {

	type itemQuantity struct {
		Item     itm.Item
		Quantity int
	}

	var itemWithQuantity []itemQuantity

	for item, quantity := range inv.Items {
		itemWithQuantity = append(itemWithQuantity, itemQuantity{item, quantity})
	}

	sort.Slice(itemWithQuantity, func(i, j int) bool {
		if ascending {
			return itemWithQuantity[i].Quantity < itemWithQuantity[j].Quantity
		} else {
			return itemWithQuantity[i].Quantity > itemWithQuantity[j].Quantity
		}
	})

	for _, itemQuantity := range itemWithQuantity {
		text := fmt.Sprintf("Item: %s, Quantity: %d", itemQuantity.Item.GetName(), itemQuantity.Quantity)
		fmt.Println(text)
	}
	fmt.Println()
}

func (inv *Inventory) StructureInventory(ascending bool) string {

	type itemQuantity struct {
		Item     itm.Item
		Quantity int
	}

	var itemWithQuantity []itemQuantity

	for item, quantity := range inv.Items {
		itemWithQuantity = append(itemWithQuantity, itemQuantity{item, quantity})
	}

	sort.Slice(itemWithQuantity, func(i, j int) bool {
		if ascending {
			return itemWithQuantity[i].Quantity < itemWithQuantity[j].Quantity
		} else {
			return itemWithQuantity[i].Quantity > itemWithQuantity[j].Quantity
		}
	})

	text := ""
	for _, itemQuantity := range itemWithQuantity {
		text += fmt.Sprintf("Item: %s, Quantity: %d \n", itemQuantity.Item.GetName(), itemQuantity.Quantity)
	}
	return text
}

func (inv *Inventory) AddItems(filters []func(item itm.Item, quantity int) error, items []itm.Item, quantity int) error {

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	initialLogID := len(inv.Log)
	recordResult := make(chan itm.Item)
	for i, item := range items {

		itemErr := inv.applyFilters(filters, item, inv.Items[item]+quantity)
		if itemErr != nil {
			err = errors.Join(err, itemErr)
			inv.logOperation(initialLogID+i+1, item, quantity, true, itemErr)
		} else {
			go inv.addItem(ctx, item, quantity, recordResult)
		}

	}

	for i := range items {
		inv.logOperation(initialLogID+i+1, <-recordResult, quantity, true, nil)
	}

	return flt.NewFilterError(err, "Filter errors")
}

func (inv *Inventory) addItem(ctx context.Context, item itm.Item, quantity int, recordResult chan itm.Item) {

	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(*itm.ErrorItemDoesNotExist); ok {
				inv.Items[item] = quantity
			} else {
				panic(err)
			}
		}
		recordResult <- item
	}()

	inv.mu.Lock()
	defer inv.mu.Unlock()

	time.Sleep(10 * time.Second) // предположим, что добавление новых элементов это долгий процесс

	select {
	case <-ctx.Done():
		fmt.Println("Operation timed out")
		return
	default:
	}

	inv.changeQantity(item, quantity)
}

func (inv *Inventory) changeQantity(item itm.Item, quantity int) {

	if _, exist := inv.Items[item]; !exist {
		panic(&itm.ErrorItemDoesNotExist{"item does not exist", nil})
	}

	inv.Items[item] += quantity
}

func (inv *Inventory) RemoveItems(filters []func(item itm.Item, quantity int) error, items []itm.Item, quantity int) error {
	var err error
	initialLogID := len(inv.Log)
	recordResult := make(chan itm.Item, len(items)) // Буферизованный канал

	// Список предметов, которые успешно прошли фильтры
	var filteredItems []itm.Item

	for i, item := range items {
		itemErr := inv.applyFilters(filters, item, inv.Items[item]-quantity)
		if itemErr != nil {
			err = errors.Join(err, itemErr)
			inv.logOperation(initialLogID+i+1, item, quantity, false, itemErr)
		} else {
			filteredItems = append(filteredItems, item)
		}
	}

	// Обработка только тех предметов, которые прошли фильтры
	for _, itm := range filteredItems {
		go inv.removeItem(itm, quantity, recordResult)
	}

	// Логирование для всех обработанных предметов
	for i := 0; i < len(filteredItems); i++ {
		inv.logOperation(initialLogID+i+1, <-recordResult, quantity, false, nil)
	}

	return flt.NewFilterError(err, "Filter errors")
}

func (inv *Inventory) removeItem(item itm.Item, quantity int, recordResult chan itm.Item) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in removeItem: %v\n", r)
		}
		recordResult <- item
	}()

	inv.mu.Lock()
	defer inv.mu.Unlock()

	if _, exist := inv.Items[item]; exist {
		inv.Items[item] -= quantity
		if inv.Items[item] <= 0 {
			delete(inv.Items, item)
		}
	}
}

func (inv *Inventory) applyFilters(filters []func(item itm.Item, quantity int) error, item itm.Item, quantity int) error {

	var err error

	for _, filter := range filters {

		filterErr := filter(item, quantity)
		if filterErr != nil {
			err = errors.Join(err, filterErr)
		}
	}

	return err

}

func (inv *Inventory) UseItem(item itm.Item) {
	fmt.Printf("Using %s %s \n", item.GetType(), item.GetName())
	item.Use() // Вызов общего метода Use
}

func (inv *Inventory) Clear() {
	inv.Items = make(map[itm.Item]int)
}

type itemWithQuantity struct {
	ItemType   string         `json:"itemType"`
	ItemConfig itm.ItemConfig `json:"ItemConfig"`
	Quantity   int            `json:"quantity"`
}

func (inv *Inventory) SaveData() error {

	var itemsList []itemWithQuantity

	for item, quantity := range inv.Items {
		itemsList = append(itemsList, itemWithQuantity{
			ItemType:   item.GetType(),
			ItemConfig: itm.GetItemConfig(item),
			Quantity:   quantity,
		})
	}

	var err error

	itemsData, jsonErr := json.Marshal(itemsList)
	if jsonErr != nil {
		return jsonErr
	}

	err = os.WriteFile("save.txt", itemsData, 0644)

	return err

}

func (inv *Inventory) ReadData() error {

	var err error

	var itemsList []itemWithQuantity

	dataJson, err := os.ReadFile("save.txt")
	if err != nil {
		return err
	}

	err = json.Unmarshal(dataJson, &itemsList)
	if err != nil {
		return err
	}

	for _, itemWithQuantity := range itemsList {

		item, err := itm.NewItem(itemWithQuantity.ItemType, itemWithQuantity.ItemConfig)
		if err != nil {
			continue
		}
		inv.Items[item] = itemWithQuantity.Quantity
	}

	return err

}
