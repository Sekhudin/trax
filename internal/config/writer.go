package config

type Writer interface {
	Write() error
	File() string
}

type configwriter struct {
	file     string
	override bool
	safe     func(string) error
	force    func(string) error
}

func NewWriter(f string, o bool, safe, force func(string) error) Writer {
	return &configwriter{
		file:     f,
		override: o,
		safe:     safe,
		force:    force,
	}
}

func (c *configwriter) Write() error {
	if c.override {
		return c.force(c.file)
	}

	return c.safe(c.file)
}

func (c *configwriter) File() string {
	return c.file
}
