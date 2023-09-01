package domain

import "fmt"

type Alias struct {
	ID     uint64
	User   string
	Domain string
}

func (a *Alias) GetEmail() string {
	return fmt.Sprintf("%s@%s", a.User, a.Domain)
}
