package display

type Display struct{}

func New() (d *Display)   { return &Display{} }
func (d *Display) Init()  {}
func (d *Display) Clear() {}
