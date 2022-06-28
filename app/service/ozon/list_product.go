package ozon

import (
	"bufio"
	"context"
	"math"
	"os"
	"strconv"
	"strings"

	models "SnipAndNeat/generated"

	"github.com/diphantxm/ozon-api-client/ozon"
	"github.com/samber/lo"
)

const tablePath = "../../temp_scropt/prices.csv"

// получает все продукты из озона и заполняет цены из таблицы, сформированной на Python
// цены из прайса Виенто
func (z *OzonAPI) ListProducts(ctx context.Context) ([]*models.VientoProduct, error) {
	r, err := z.client.Products().GetListOfProducts(ctx, &ozon.GetListOfProductsParams{})
	if err != nil {
		return nil, err
	}

	var producIDs []int64
	lo.Map(r.Result.Items, func(v ozon.GetListOfProductsResultItem, _ int) []int64 {
		producIDs = append(producIDs, v.ProductId)
		return producIDs
	})
	result, err := z.client.Products().ListProductsByIDs(ctx, &ozon.ListProductsByIDsParams{
		ProductId: producIDs,
	})

	var tempMap = make(map[string]*models.VientoProduct)
	for _, v := range result.Result.Items {
		tempMap[v.Barcode] = &models.VientoProduct{
			ProductID: models.NewOptInt64(v.Id),
			OfferID:   models.NewOptString(v.OfferId),
			Barcode:   models.NewOptString(v.Barcode),
		}
	}

	priceFile, err := os.Open(tablePath)
	if err != nil {
		return nil, err
	}
	defer priceFile.Close()

	priceMap := make(map[string]*models.VientoProduct)
	for scanner := bufio.NewScanner(priceFile); scanner.Scan(); {
		fields := strings.Split(scanner.Text(), ",")
		if len(fields) != 3 {
			continue
		}

		price, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			continue
		}

		priceMap[fields[2]] = &models.VientoProduct{
			Barcode: models.NewOptString(fields[2]),
			Price:   models.NewOptFloat64(math.Round(price)),
		}
	}

	for k, v := range tempMap {
		if p, ok := priceMap[k]; ok {
			v.IsFBOVisible = models.NewOptBool(true)
			v.Barcode = p.Barcode
			v.Price = p.Price
		} else {
			v.IsFBOVisible = models.NewOptBool(false)
		}
	}

	return lo.Values(tempMap), nil
}

// grasp
// solid
// gop
