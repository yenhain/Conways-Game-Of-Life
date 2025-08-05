package gol

// Params provides the details of how to run the Game of Life and which image to load.
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run(p Params, events chan<- Event, keyPresses <-chan rune) {

	//initialises channels
	ioCommand := make(chan ioCommand)
	ioIdle := make(chan bool)
	///initialise channels
	filename := make(chan string)
	output := make(chan uint8)
	input := make(chan uint8)
	///

	///initialise struct
	distributorChannels := distributorChannels{
		events,
		ioCommand,
		ioIdle,
		filename,
		output,
		input,
	}

	go distributor(p, distributorChannels, keyPresses)

	ioChannels := ioChannels{
		command:  ioCommand,
		idle:     ioIdle,
		filename: filename,
		output:   output,
		input:    input,
	}
	go startIo(p, ioChannels)
}
