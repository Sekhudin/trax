package outmock

type Color struct {
	RedCalled    bool
	YellowCalled bool
	GreenCalled  bool
	BlueCalled   bool
	CyanCalled   bool
	GrayCalled   bool
	BoldCalled   bool
}

func (c *Color) Red(v ...any) string {
	c.RedCalled = true
	return ""
}

func (c *Color) Yellow(v ...any) string {
	c.YellowCalled = true
	return ""
}

func (c *Color) Green(v ...any) string {
	c.GreenCalled = true
	return ""
}

func (c *Color) Blue(v ...any) string {
	c.BlueCalled = true
	return ""
}

func (c *Color) Cyan(v ...any) string {
	c.CyanCalled = true
	return ""
}

func (c *Color) Gray(v ...any) string {
	c.GrayCalled = true
	return ""
}

func (c *Color) Bold(v ...any) string {
	c.BoldCalled = true
	return ""
}

type Out struct {
	SuccesCalled bool
	InfoCalled   bool
	WarnCalled   bool
	ErrorCalled  bool
	CauseCalled  bool

	AsFlatCalled bool
	AsJsonCalled bool

	JSONErr error
	FlatErr error
}

func (o *Out) Success(scope, msg string) {
	o.SuccesCalled = true
}

func (o *Out) Info(scope, msg string) {
	o.InfoCalled = true
}

func (o *Out) Warn(scope, msg string) {
	o.WarnCalled = true
}

func (o *Out) Error(scope, msg string) {
	o.ErrorCalled = true
}

func (o *Out) Cause(scope, msg string) {
	o.CauseCalled = true
}

func (o *Out) AsFlat(prefix string, data map[string]any) error {
	o.AsFlatCalled = true
	return o.FlatErr
}

func (o *Out) AsJSON(data map[string]any) error {
	o.AsJsonCalled = true
	return o.JSONErr
}
