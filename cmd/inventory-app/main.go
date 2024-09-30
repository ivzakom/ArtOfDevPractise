package main

import (
	"artOfDevPractise/internal/config"
	"artOfDevPractise/internal/inventory"
	"artOfDevPractise/internal/user"
	"artOfDevPractise/internal/user/db"
	"artOfDevPractise/pkg/client/mongoDB"
	"artOfDevPractise/pkg/logging"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	w.Write([]byte(fmt.Sprintf("Hello %s", name)))
}

func main() {

	logger := logging.GetLogger()
	logger.Info("Starting server")

	cfg := config.GetConfig()
	cfgMongo := cfg.MongoDB
	MongoDBCient, err := mongoDB.NewCient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username, cfgMongo.Password, cfgMongo.Database, cfgMongo.Auth_db)
	if err != nil {
		panic(err)
	}

	storage := db.NewStorage(MongoDBCient, cfgMongo.Collection, logger)

	user1 := user.User{
		ID:           "",
		Email:        "izakom@mail.ru",
		Username:     "Igor",
		PasswordHash: "1234",
	}

	user1ID, err := storage.Create(context.Background(), user1)

	logger.Info(user1ID)

	user2 := user.User{
		ID:           "",
		Email:        "izakom2@mail.ru",
		Username:     "Igor2",
		PasswordHash: "12345",
	}

	user2ID, err := storage.Create(context.Background(), user2)
	if err != nil {
		panic(err)
	}

	logger.Info(user2ID)

	inv := inventory.NewInventory()
	logger.Info("start read data")
	inv.ReadData()
	logger.Info("end read data")

	router := httprouter.New()

	handlerInv := inventory.NewHandler(inv, logger)
	handlerInv.Register(router)

	handlerUser := user.NewHandler(logger)
	handlerUser.Register(router)

	start(router, cfg)

	//filters := []func(item item.Item, quantity int) error{
	//	filters.FilterNotNull,
	//}
	//
	//sword, err := item.NewItem(item.WeaponType, item.ItemConfig{
	//	ID:     1,
	//	Name:   "sword",
	//	Damage: 5,
	//})
	//
	//if err != nil {
	//	fmt.Println("Error creating item:", err)
	//	return
	//}
	//
	//shild, err := item.NewItem(item.ArmorType, item.ItemConfig{
	//	ID:     2,
	//	Name:   "shild",
	//	Defend: 10,
	//	Hp:     50,
	//})
	//
	//if err != nil {
	//	fmt.Println("Error creating item:", err)
	//	return
	//}
	//
	//arr := []item.Item{sword, shild}
	//
	//inv.AddItems(filters, arr, 5)
	//
	//inv.ShowInventory(true)
	//
	//inv.AddItems(filters, []item.Item{shild}, 10)
	//
	//inv.ShowInventory(true)
	//inv.ShowInventory(false)
	//
	//err = inv.RemoveItems(filters, []item.Item{sword, shild}, 50)
	//
	//inv.ShowInventory(true)
	//
	//for i, e := range inv.Log {
	//	fmt.Printf("%d - %s \n", i, e)
	//}
	//
	//if err != nil {
	//	fmt.Printf("Error: %s\n", err.Error())
	//}
	//
	//err = inv.SaveData()
	//if err != nil {
	//	fmt.Printf("SeveError: %s\n", err.Error())
	//}
	//
	//inv.Clear()
	//
	//inv.ShowInventory(true)
	//
	//inv.UseItem(shild)
	//fmt.Println(shild)

}

func start(router *httprouter.Router, cfg *config.Config) {

	logger := logging.GetLogger()
	logger.Info("start server")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket ")
		socketPath := path.Join(appdir, "app.sock")
		logger.Debugf("socket path: %s", socketPath)

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)

	} else {
		logger.Info("listen unix")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
		logger.Info(fmt.Sprintf("start is lissening %s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
	}

	if listenErr != nil {
		panic(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	server.Serve(listener)

}
