package gate

type Gate struct{}

func NewGate() *Gate {
	return new(Gate)
}

func (g *Gate) Run() {

}

func (g *Gate) Close() {

}
