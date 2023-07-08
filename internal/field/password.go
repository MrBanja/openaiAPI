package field

type Password string

func (p Password) String() string {
	return "***"
}

func (p Password) Reveal() string {
	return string(p)
}
