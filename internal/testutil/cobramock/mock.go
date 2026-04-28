package cobramock

import "errors"

type FlagBroken struct{}

func (b *FlagBroken) String() string {
	return ""
}

func (b *FlagBroken) Set(string) error {
	return errors.New("flag error")
}

func (b *FlagBroken) Type() string {
	return "broken"
}
