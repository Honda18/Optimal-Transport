package main 

type Factory struct{
	production int
	availableOutput int
}

func newFactory(production int)(*Factory){
	return &Factory{production, production}
}

func (factory *Factory) decrementOutput(difference int) {
	factory.availableOutput=(factory.availableOutput-difference)
}

func (factory *Factory) incrementOutput(add int) {
	factory.availableOutput=(factory.availableOutput+add)
}

func(factory *Factory) getOutput() int{
	return factory.availableOutput
}

func(factory *Factory) getProduction() int{
	return factory.production
}

func(factory *Factory) available() bool{
	return factory.availableOutput>0
}



