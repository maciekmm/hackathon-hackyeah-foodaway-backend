package model

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/maciekmm/HackYeah/utils"
)

var (
	ErrOfferInternal           = errors.New("internal error")
	ErrOfferIDInvalid          = errors.New("invalid id")
	ErrOfferTitleInvalid       = errors.New("invalid title")
	ErrOfferDescriptionInvalid = errors.New("invalid description")
	ErrOfferLocationInvalid    = errors.New("invalid pickup location")
	ErrOfferPickupInvalid      = errors.New("invalid pickup duration")
)

type Offer struct {
	gorm.Model
	UserID uint `json:"user_id,omitempty"`

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Picture     string `json:"picture,omitempty"`
	Expiration  uint64 `json:"expiration,omitempty"`

	Longitude   *float64 `json:"longitude,omitempty"`
	Latitude    *float64 `json:"latitude,omitempty"`
	PickupStart uint64   `json:"pickup_start,omitempty"`
	PickupEnd   uint64   `json:"pickup_end,omitempty"`
}

func (o *Offer) Add(db *gorm.DB) error {
	errors := []error{}
	if len(o.Title) == 0 {
		errors = append(errors, ErrOfferTitleInvalid)
	}

	if len(o.Description) == 0 {
		errors = append(errors, ErrOfferDescriptionInvalid)
	}

	if o.Longitude == nil || o.Latitude == nil {
		errors = append(errors, ErrOfferLocationInvalid)
	}

	if o.PickupStart == 0 || o.PickupEnd == 0 {
		errors = append(errors, ErrOfferPickupInvalid)
	}

	if len(errors) > 0 {
		return utils.NewErrorResponse(errors...)
	}

	if res := db.Create(&o); res.Error != nil {
		return &utils.ErrorResponse{
			Errors:      []string{ErrOfferInternal.Error()},
			DebugErrors: []string{res.Error.Error()},
		}
	}
	return nil
}
