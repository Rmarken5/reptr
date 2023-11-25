package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rmarken/reptr/service/internal/logic"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strconv"
	"strings"
	"time"
)

const uri = "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000"

type (
	menu struct {
		logic logic.Controller
	}
)

func main() {
	ctx := context.Background()
	log := zerolog.New(os.Stdout).With().Str("program", "crud tester").Logger()
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Panic().Err(err).Msg("while connecting to mongo")
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	db := client.Database("deck")

	l := logic.New(log, db)

	m := menu{
		l,
	}

	m.start()
}

func (m *menu) start() {
	for {
		choice := m.promptMenu()
		switch choice {
		case 1:
			m.createDeck()
		case 2:
			m.createCard()
		case 3:
			m.updateCard()
		case 4:
			m.getDecks()
		case 9:
			os.Exit(1)
		}

	}
}

func (m *menu) promptMenu() int {
	reader := bufio.NewReader(os.Stdin)
	sb := strings.Builder{}
	sb.WriteString("What would you like to do?\n")
	sb.WriteString("1. Create Deck\n")
	sb.WriteString("2. Create Card\n")
	sb.WriteString("3. Update Card\n")
	sb.WriteString("4. Get Decks\n")
	sb.WriteString("9. Exit\n")
	sb.WriteString("Choice: ")
	_, err := os.Stdout.WriteString(sb.String())
	if err != nil {
		panic(err)
	}
	choice, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	choice = strings.Trim(choice, "\n")
	toI, err := strconv.Atoi(choice)
	if err != nil {
		os.Stdout.WriteString("Invalid choice: " + choice + "\n")
		return m.promptMenu()
	}
	return toI
}

func (m *menu) createDeck() {
	reader := bufio.NewReader(os.Stdin)

	_, err := os.Stdout.WriteString("Name your deck: ")
	if err != nil {
		panic(err)
	}

	deckName, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	deckName = strings.Trim(deckName, "\n")

	id, err := m.logic.CreateDeck(context.TODO(), models.Deck{ID: uuid.NewString(), Name: deckName, CreatedAt: time.Now()})
	if err != nil {
		os.Stdout.WriteString(err.Error())
	}
	_, err = os.Stdout.WriteString("\nDeck Created with id: " + id + "\n")
	if err != nil {
		panic(err)
	}

}

func (m *menu) createCard() {
	reader := bufio.NewReader(os.Stdin)
	_, err := os.Stdout.WriteString("Type of card (1 for basic, 2 for multiple choice): ")

	t, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	t = strings.Trim(t, "\n")

	tToI, err := strconv.Atoi(t)
	if err != nil {
		panic("invalid choice")
	}

	_, err = os.Stdout.WriteString("\nFront of Card: ")
	if err != nil {
		panic(err)
	}
	front, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	front = strings.Trim(front, "\n")

	_, err = os.Stdout.WriteString("\nBack of Card: ")
	if err != nil {
		panic(err)
	}
	back, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	back = strings.Trim(back, "\n")

	_, err = os.Stdout.WriteString("\nID of deck to add card to: ")
	if err != nil {
		panic(err)
	}
	deckID, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	deckID = strings.Trim(deckID, "\n")

	err = m.logic.AddCardToDeck(context.TODO(), deckID, models.Card{
		ID:        uuid.NewString(),
		Front:     front,
		Back:      back,
		Kind:      models.Type(tToI - 1),
		CreatedAt: time.Now(),
	})
	if err != nil {
		os.Stdout.WriteString(err.Error())
	}
	_, err = os.Stdout.WriteString("\nCard Created.\n")
	if err != nil {
		panic(err)
	}
}

func (m *menu) updateCard() {
	reader := bufio.NewReader(os.Stdin)
	_, err := os.Stdout.WriteString("ID of the card you want to update: ")

	cardID, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	cardID = strings.Trim(cardID, "\n")

	_, err = os.Stdout.WriteString("\nFront of Card (enter to skip): ")
	if err != nil {
		panic(err)
	}
	front, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	front = strings.Trim(front, "\n")

	_, err = os.Stdout.WriteString("\nBack of Card (enter to skip): ")
	if err != nil {
		panic(err)
	}
	back, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	back = strings.Trim(back, "\n")

	err = m.logic.UpdateCard(context.TODO(), models.Card{
		ID:        cardID,
		Front:     front,
		Back:      back,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		os.Stdout.WriteString(err.Error())
	}
	_, err = os.Stdout.WriteString("\nCard Created.\n")
	if err != nil {
		panic(err)
	}
}

func (m *menu) getDecks() {
	now := time.Now()
	decks, err := m.logic.GetDecks(context.TODO(), time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil, 0, 0)
	if err != nil {
		os.Stdout.WriteString(err.Error())
		return
	}
	sb := strings.Builder{}

	for _, deck := range decks {
		sb.WriteString(fmt.Sprintf("%+v\n", deck))
	}

	os.Stdout.WriteString(sb.String())

}

func (m *menu) getGroups() {
	decks, err := m.logic.GetGroups(context.TODO(), time.Now().Truncate(time.Hour), nil, 0, 0)

	if err != nil {
		os.Stdout.WriteString(err.Error())
		return
	}
	sb := strings.Builder{}

	for _, deck := range decks {
		sb.WriteString(fmt.Sprintf("%+v\n", deck))
	}

	os.Stdout.WriteString(sb.String())
}
