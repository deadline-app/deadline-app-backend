package models

import "github.com/go-pg/pg/v10"

type Card struct {
	ID        int64  `json:"id"`
	Subject   string `json:"subject"`
	Task_name string `json:"task_name"`
}

func CreateCard(db *pg.DB, req *Card) (*Card, error) {
	_, err := db.Model(req).Insert()
	if err != nil {
		return nil, err
	}

	card := &Card{}

	err = db.Model(card).
		Where("card.id = ?", req.ID).
		Select()

	return card, err
}

func GetCard(db *pg.DB, cardID string) (*Card, error) {
	card := &Card{}

	err := db.Model(card).
		Where("card.id = ?", cardID).
		Select()

	return card, err
}

func GetAllCards(db *pg.DB) ([]*Card, error) {
	cards := make([]*Card, 0)

	err := db.Model(&cards).
		Select()

	return cards, err
}
