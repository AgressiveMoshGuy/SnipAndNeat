package models

import (
	null "github.com/guregu/null"
)

type FeedBack struct {
	Rating     int8
	Unique     int8
	Fitted     null.Int
	Custom     null.Int
	Assortment int8
	Style      style
	Impression int8
	Quality    int8
	Refund     int8
}
