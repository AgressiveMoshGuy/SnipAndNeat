// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// GetVientoProducts implements getVientoProducts operation.
//
// Get viento products.
//
// POST /viento/products
func (UnimplementedHandler) GetVientoProducts(ctx context.Context) (r GetVientoProductsRes, _ error) {
	return r, ht.ErrNotImplemented
}