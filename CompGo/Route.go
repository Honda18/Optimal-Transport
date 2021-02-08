package main

type Route struct{
	factory *Factory 
	warehouse *Warehouse
	cost int
	unitsAllocated int
	tested bool
}

func newRoute(factory *Factory, warehouse *Warehouse, cost int) (*Route){
	return &Route{factory, warehouse, cost, 0, false}
}

func (route *Route) allocateToRoute(units int){
	route.factory.decrementOutput(units)
	route.warehouse.decrementAvailability(units)
	route.unitsAllocated+=units
}

func (route *Route) decrementRoute(units int){
	route.factory.incrementOutput(units)
	route.warehouse.incrementAvailability(units)
	route.unitsAllocated-=units
}

func (route *Route) getCost() int{
	return route.cost
}

func (route *Route) getWarehouse() (*Warehouse){
	return route.warehouse
}

func (route *Route) getFactory() (*Factory){
	return route.factory
}

func (route *Route) getUnitsAllocated() int{
	return route.unitsAllocated
}

func (route *Route) unused() bool{
	return route.unitsAllocated==0
}

func (route *Route) beenTested() bool{
	return route.tested
}	

func (route *Route) setTestedFalse(){
	route.tested= false
}

func (route *Route) setTestedTrue(){
	route.tested= true
}







