package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/maciekmm/HackYeah/app"
	"github.com/maciekmm/HackYeah/model"
)

const (
	CracowLatitude  = 50.06143
	CracowLongitude = 19.93658
	//MaxDistance in degrees
	MaxDistance = 0.3
)

var (
	food = []string{"Artichoke", "Arugula", "Asparagus", "Green, Purple, White", "Avocado", "Bamboo Shoots", "Bean Sprouts", "Beans- see Bean List", "Beet", "Belgian Endive", "Bell Pepper", "Bitter Melon/Bitter Gourd", "Bok Choy/Bok Choi/Pak Choy", "Broccoli", "Brussels Sprouts", "Burdock Root/Gobo", "Cabbage", "Green, Red, Savoy", "Calabash", "Capers", "Carrot", "Cassava/Yuca", "Cauliflower", "Celery", "Celery Root/Celeriac", "Celtuce", "Chayote", "Chinese Broccoli/Kai-lan", "Corn/Maize", "Baby Corn/Candle Corn", "Cucumber", "English Cucumber", "Gherkin", "Pickling Cucumbers", "Daikon Radish", "Edamame", "Eggplant/Aubergine", "Elephant Garlic", "Endive", "Curly/Frisee", "Escarole", "Fennel", "Fiddlehead", "Galangal", "Garlic", "Ginger", "Grape Leaves", "Green Beans/String Beans/Snap Beans", "Wax Beans", "Greens", "Amaranth Leaves/Chinese Spinach", "Beet Greens", "Collard Greens", "Dandelion Greens", "Kale", "Kohlrabi Greens", "Mustard Greens", "Rapini", "Spinach", "Swiss Chard", "Turnip Greens", "Hearts of Palm", "Horseradish", "Jerusalem Artichoke/Sunchokes", "Jícama", "Kale", "Curly", "Lacinato", "Ornamental", "Kohlrabi", "Leeks", "Lemongrass", "Lettuce", "Butterhead- Bibb, Boston", "Iceberg", "Leaf- Green Leaf, Red Leaf", "Romaine", "Lotus Root", "Lotus Seed", "Mushrooms- see Mushroom List", "Napa Cabbage", "Nopales", "Okra", "Olive", "Onion", "Green Onions/Scallions", "Parsley", "Parsley Root", "Parsnip", "Peas", "green peas", "snow peas", "sugar snap peas", "Peppers- see Peppers List", "Plantain", "Potato", "Pumpkin", "Purslane", "Radicchio", "Radish", "Rutabaga", "Sea Vegetables- see Sea Vegetable List", "Shallots", "Spinach", "Squash- see Squash List", "Sweet Potato", "Swiss Chard", "Taro", "Tomatillo", "Tomato", "Turnip", "Water Chestnut", "Water Spinach", "Watercress", "Winter Melon", "Yams", "Zucchini"}
)

func nextSign() float64 {
	r := rand.Intn(2)
	if r > 0 {
		return -1
	}
	return 1
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Invalid number of arguments.")
		os.Exit(1)
	}
	number, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("First argument should be a count")
		os.Exit(1)
	}

	logger := log.New(os.Stdout, "HackYeah!", log.Ldate|log.Lshortfile)
	app := &app.Application{Logger: logger}

	err = app.Init()
	if err != nil {
		logger.Fatal(err)
	}

	for ; number > 0; number-- {
		latitude := CracowLatitude + nextSign()*MaxDistance*rand.Float64()
		longitude := CracowLongitude + nextSign()*MaxDistance*rand.Float64()
		offer := model.Offer{
			Latitude:    &latitude,
			Longitude:   &longitude,
			Title:       food[rand.Intn(len(food))],
			Expiration:  uint64(time.Now().Unix() + int64(rand.Intn(1000*60*60*24*7))),
			PickupStart: uint64(time.Now().Unix()),
			PickupEnd:   uint64(time.Now().Unix() + int64(rand.Intn(1000*60*60*24*7))),
			Description: "test",
		}
		if err := offer.Add(app.Database); err != nil {
			logger.Fatal(err)
		}
	}
}
