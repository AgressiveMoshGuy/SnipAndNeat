package ozon

import (
	"SnipAndNeat/app/config"
	storage "SnipAndNeat/app/strorage"
	models "SnipAndNeat/generated"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/diphantxm/ozon-api-client/ozon"
	ozon_cli "github.com/diphantxm/ozon-api-client/ozon"
	"github.com/rs/zerolog"
)

type Ozon interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	VientoProductsInfo(ctx context.Context, in *ozon.GetListOfProductsParams) (*ozon.GetListOfProductsResponse, error)
	ListProducts(ctx context.Context) ([]*models.VientoProduct, error)
}

type OzonAPI struct {
	log         zerolog.Logger
	initialDate time.Time
	client      *ozon_cli.Client
	db          *storage.DB
	memcache    *memcache.Client
}

// {
// 	"Host": "https://api-seller.ozon.ru",
// 	"Client-Id": 284574,
// 	"Api-Key": "aab10a66-435c-4166-89aa-4aba8d2fa2a1",
// 	"Content-Type": "application/json"
//}

func New(ctx context.Context, cfg *config.Config, d *storage.DB) (*OzonAPI, error) {
	mClient := memcache.New(cfg.OzonClientConfig.MemcacheHost)

	opts := []ozon.ClientOption{
		ozon.WithAPIKey(cfg.OzonClientConfig.APIKey),
		ozon.WithClientId(cfg.OzonClientConfig.ClientID),
		ozon.WithHttpClient(&http.Client{
			Transport: &http.Transport{
				IdleConnTimeout:       cfg.OzonClientConfig.Timeout,
				ResponseHeaderTimeout: cfg.OzonClientConfig.Timeout,
				ExpectContinueTimeout: cfg.OzonClientConfig.Timeout,
			},
		}),
	}

	ozonClient := ozon_cli.NewClient(opts...)

	return &OzonAPI{
		log:      zerolog.New(os.Stdout),
		db:       d,
		client:   ozonClient,
		memcache: mClient,
	}, nil
}

func (z *OzonAPI) Start(ctx context.Context) error {
	return nil
}

func (z *OzonAPI) Stop(ctx context.Context) error {
	return nil
}
