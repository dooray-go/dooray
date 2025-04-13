package calendar

type Calendar struct {
	endPoint string
}

func NewDefaultCalendar() *Calendar {
	return &Calendar{
		endPoint: "https://api.dooray.com",
	}
}
func NewCalendar(endPoint string) *Calendar {
	return &Calendar{
		endPoint: endPoint,
	}
}
