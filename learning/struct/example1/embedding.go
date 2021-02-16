package example1

type first struct{}

func (f first) share() string {
	return "shared for 1st"
}

type second struct{}

func (f second) share() string {
	return "shared for 2nd"
}

type Common struct {
	first first
	second
}

func (c *Common) F() string {
	return c.share()
}
