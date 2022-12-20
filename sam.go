package main

import (
	"fmt"
	"log"
	"time"

	"github.com/arlsclu/notify"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

var notifier = notify.NewWeNotifier()

var (
	// control the repid limit of sending message
	limiter = make(chan struct{}, 5)
	// repid of check if event(s) happen
	checkEvery = 30 * time.Second

	cfg = config{}

	ls = make(map[string]chan struct{})
)

type configItem struct {
	Name  string `mapstructure:"name"`
	SpuID string `mapstructure:"spuID"`
}
type config struct {
	DetailItems []configItem `mapstructure:"detailItems"`
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	log.Printf("reading configi file success: %s ", viper.GetString("name"))

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("readconfig %+v ", cfg)
}

func main() {
	log.Println("server  start ... ")
	for _, v := range cfg.DetailItems {
		ls[v.Name] = limiter
	}
	go func() {
		for range time.Tick(10 * time.Minute) {
			for _, v := range cfg.DetailItems {
				<-ls[v.Name]
				log.Printf("清理了%s ", v.Name)
			}
		}
	}()

	c := cron.New(cron.WithSeconds())
	_, _ = c.AddFunc(fmt.Sprintf("@every %s", checkEvery.String()), listenSam)
	c.Run()
}

func listenSam() {
	for _, v := range cfg.DetailItems {
		register(v.Name, v.SpuID, checkerByDetail)
	}
}

func register(name string, spuID string, f func(spuID string) (bool, error)) {
	toSend, err := f(spuID)
	if err != nil {
		log.Println(err)
		notifier.Send(err.Error())
		return
	}
	if toSend {
		ls[name] <- struct{}{}
		msg := fmt.Sprintf("%s 极速达可以购买了", name)
		log.Printf("%s 条件成立", msg)
		notifier.Send(msg)
		return
	}
	log.Printf("%s 未监测到true", name)
}
