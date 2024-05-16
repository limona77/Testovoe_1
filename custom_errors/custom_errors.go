package custom_errors

import "errors"

var (
	ErrNotOpenYet       = errors.New("NotOpenYet")
	ErrYouShallNotPass  = errors.New("YouShallNotPass")
	ErrPlaceIsBusy      = errors.New("PlaceIsBusy")
	ErrClientUnknown    = errors.New("ClientUnknown")
	ErrICanWaitNoLonger = errors.New("ICanWaitNoLonger!")
	ErrCode             = errors.New("13")
	ErrActionNotExist   = errors.New("ActionNotExist")
)
