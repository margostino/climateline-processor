package cache

import (
	"github.com/margostino/climateline-processor/domain"
	"net/http"
)

func Delete(writer *http.ResponseWriter, db *map[string]domain.Item) {
	(*writer).WriteHeader(http.StatusOK)
	*db = make(map[string]domain.Item)
}
