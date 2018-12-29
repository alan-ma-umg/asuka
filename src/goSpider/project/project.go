package project

import (
	"net/http"
	"time"
)

type Project interface {
	// session
	NeedToLogin() bool
	IsLogin() bool
	LoginWithHeader() *http.Header

	// queue
	EnqueueFilter() bool

	// frequency
	NeedToPause() bool
	Throttle() time.Duration
}
