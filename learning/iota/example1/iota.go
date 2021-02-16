package example1

import "errors"

const (
	monday int = iota + 1
	tuesday
	wednesday
	thursday
	friday
	saturday
	sunday
)

const (
	Monday    = "Monday"
	Tuesday   = "Tuesday"
	Wednesday = "Wednesday"
	Thursday  = "Thursday"
	Friday    = "Friday"
	Saturday  = "Saturday"
	Sunday    = "Sunday"
)

func dayNumber(day string) (int, error) {
	switch day {
	case Monday:
		return monday, nil
	case Tuesday:
		return tuesday, nil
	case Wednesday:
		return wednesday, nil
	case Thursday:
		return thursday, nil
	case Friday:
		return friday, nil
	case Saturday:
		return saturday, nil
	case Sunday:
		return sunday, nil
	default:
		return -1, errors.New("wrong day passed")
	}
}
