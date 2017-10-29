package main

import (
	"fmt"
	"log"
	"math"
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
	food   = []string{"Grapes", "Melon", "Watermelon", "Tangerine", "Lemon", "Banana", "Pineapple", "Red Apple", "Green Apple", "Pear", "Peach", "Cherries", "Strawberry", "Kiwi Fruit", "Tomato", "Avocado", "Eggplant", "Potato", "Carrot", "Ear of Corn", "Hot Pepper", "Cucumber", "Broccoli", "Mushroom", "Peanuts", "Chestnut", "Bread", "Croissant", "Baguette Bread", "Pancakes", "Cheese Wedge", "Meat on Bone", "Poultry Leg", "Cut of Meat", "Bacon", "Hamburger", "French Fries", "Pizza", "Hot Dog", "Sandwich", "Taco", "Burrito", "Green Salad", "Popcorn", "Canned Food", "Bento Box", "Rice Cracker", "Rice Ball", "Cooked Rice", "Curry Rice", "Spaghetti", "Roasted Sweet Potato", "Oden", "Sushi", "Fried Shrimp", "Fish Cake With Swirl", "Dango", "Dumpling", "Fortune Cookie", "Soft Ice Cream", "Ice Cream", "Doughnut", "Cookie", "Birthday Cake", "Shortcake", "Pie", "Chocolate Bar", "Candy", "Lollipop", "Honey"}
	emojis = []rune{'🍇', '🍈', '🍉', '🍊', '🍋', '🍌', '🍍', '🍎', '🍏', '🍐', '🍑', '🍒', '🍓', '🥝', '🍅', '🥑', '🍆', '🥔', '🥕', '🌽', '🌶', '🥒', '🥦', '🍄', '🥜', '🌰', '🍞', '🥐', '🥖', '🥞', '🧀', '🍖', '🍗', '🥩', '🥓', '🍔', '🍟', '🍕', '🌭', '🥪', '🌮', '🌯', '🥗', '🍿', '🥫', '🍱', '🍘', '🍙', '🍚', '🍛', '🍝', '🍠', '🍢', '🍣', '🍤', '🍥', '🍡', '🥟', '🥠', '🍦', '🍨', '🍩', '🍪', '🎂', '🍰', '🥧', '🍫', '🍬', '🍭', '🍯'}
)

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
		angle := rand.Float64() * 2 * math.Pi
		radius := MaxDistance * rand.Float64()
		latitude := CracowLatitude + math.Sin(angle)*radius
		longitude := CracowLongitude + math.Cos(angle)*radius
		ind := rand.Intn(len(food))
		offer := model.Offer{
			Latitude:    &latitude,
			Longitude:   &longitude,
			Title:       food[ind],
			Picture:     string(emojis[ind]),
			Expiration:  uint64(time.Now().Unix() + int64(rand.Intn(1000*60*60*24*7))),
			PickupStart: uint64(time.Now().Unix()),
			PickupEnd:   uint64(time.Now().Unix() + int64(rand.Intn(1000*60*60*24*7))),
			Description: fmt.Sprintf("I have %d grams of %s to give away", int(rand.Float64()*1000), food[ind]),
		}
		if err := offer.Add(app.Database); err != nil {
			logger.Fatal(err)
		}
	}
}
