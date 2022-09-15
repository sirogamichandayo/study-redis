//go:generate mockgen -source=$GOFILE -destination=./mock/$GOFILE
package redTime

import "time"

type ITime interface {
	Now() time.Time
}
