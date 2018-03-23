package gomorra

type ComputerStatGettable interface {
	// Get the number of cores in the system
	GetCores() (int, error)
	// Get how stressed the server is, this equates to load/cores
	GetLoadPercentage() (float32, error)
}