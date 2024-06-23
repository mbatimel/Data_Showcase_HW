package items

import "time"

type Item struct {
	key interface{}
	value interface{}
	expiration *time.Time
}
