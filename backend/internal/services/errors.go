package services

import "errors"

var ErrNotFound = errors.New("not found")


var InternalMarshalError = []byte(`{"error":"Internal marshal error"}`)
var InternalCacheError = []byte(`{"error":"Internal cache error"}`)
var OrderNotFoundJSON = []byte(`{"error":"Order not found"}`)
var InternalNatsError = []byte(`{"error":"Internal nats server error"}`)
