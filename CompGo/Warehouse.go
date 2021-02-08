package main

type Warehouse struct{
	capacity int
	demandAvailability int
}

func newWarehouse(capacity int) (*Warehouse){
	return &Warehouse{capacity, capacity}
}
func (warehouse *Warehouse) decrementAvailability(difference int){
	warehouse.demandAvailability = warehouse.demandAvailability-difference
}

func (warehouse *Warehouse) incrementAvailability(sum int){
	warehouse.demandAvailability = warehouse.demandAvailability + sum
}

func (warehouse *Warehouse) getAvailability() int{
	return warehouse.demandAvailability
}

func (warehouse *Warehouse) getCapacity() int{
	return warehouse.capacity
}

func (warehouse *Warehouse) available() bool{
	return warehouse.demandAvailability>0
}