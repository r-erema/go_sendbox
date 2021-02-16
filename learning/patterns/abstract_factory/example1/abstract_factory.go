package example1

type AbstractLiquid interface {
	GetVolume() float64
}
type CocaCola struct {
	volume float64
}

func (c *CocaCola) GetVolume() float64 {
	return c.volume
}

type Pepsi struct {
	volume float64
}

func (p *Pepsi) GetVolume() float64 {
	return p.volume
}

type AbstractBottle interface {
	PourLiquid(liquid AbstractLiquid)
	GetVolume() float64
	GetLiquidVolume() float64
}
type CocaColaBottle struct {
	liquid AbstractLiquid
	volume float64
}

func (c *CocaColaBottle) PourLiquid(liquid AbstractLiquid) {
	c.liquid = liquid
}
func (c *CocaColaBottle) GetVolume() float64 {
	return c.volume
}
func (c *CocaColaBottle) GetLiquidVolume() float64 {
	return c.liquid.GetVolume()
}

type PepsiBottle struct {
	liquid AbstractLiquid
	volume float64
}

func (p *PepsiBottle) PourLiquid(liquid AbstractLiquid) {
	p.liquid = liquid
}
func (p *PepsiBottle) GetVolume() float64 {
	return p.volume
}
func (p *PepsiBottle) GetLiquidVolume() float64 {
	return p.liquid.GetVolume()
}

type AbstractFactory interface {
	CreateLiquid(volume float64) AbstractLiquid
	CreateBottle(volume float64) AbstractBottle
}
type CocaColaFactory struct{}

func (c *CocaColaFactory) CreateLiquid(volume float64) AbstractLiquid {
	return &CocaCola{volume: volume}
}
func (c *CocaColaFactory) CreateBottle(volume float64) AbstractBottle {
	return &CocaColaBottle{volume: volume}
}

type PepsiFactory struct{}

func (p *PepsiFactory) CreateLiquid(volume float64) AbstractLiquid {
	return &Pepsi{volume: volume}
}
func (p *PepsiFactory) CreateBottle(volume float64) AbstractBottle {
	return &PepsiBottle{volume: volume}
}
