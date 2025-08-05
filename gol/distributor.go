package gol

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"uk.ac.bris.cs/gameoflife/util"
)

///declare struct of channels distributor will communicate

type distributorChannels struct {
	events    chan<- Event
	ioCommand chan<- ioCommand
	ioIdle    <-chan bool

	////define distributor channels
	filename chan<- string //sends string into filename channel
	output   chan<- uint8  //sends uint8 into output channel
	input    <-chan uint8  //input recieves from uint8 channel
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels, keyPresses <-chan rune) {

	// TODO: Create a 2D slice to store the world.
	world := make([][]byte, p.ImageHeight)
	for i := range world {
		world[i] = make([]byte, p.ImageWidth)
	}

	// TODO: For all initially alive cells send a CellFlipped Event.
	///find initial alive cells from go routine startIo. read in the image
	c.ioCommand <- ioInput
	c.filename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight)}, "x")

	///store the image into a "D slice"
	for j := 0; j < p.ImageHeight; j++ {
		for k := 0; k < p.ImageWidth; k++ {
			val := <-c.input
			/* if val != 0 {
				fmt.Println("alive in", j, k) //for debug
			} */
			///initialise cell flipped to send even for alive cells
			///when image is loaded in
			loadinflip := CellFlipped{
				CompletedTurns: 0,
				Cell:           util.Cell{X: k, Y: j},
			}

			///send the cellflipped event into event channel
			c.events <- loadinflip

			///input state of cell into world slice
			world[j][k] = val
		}
	}

	var turn int
	///initialise cellflipped array to store what cells need to be
	///flipped at the end of turn

	var cellstoflip []util.Cell

	ticker := time.NewTicker(2 * time.Second)
	// TODO: Execute all turns of the Game of Life.
	///iterate through turns

	///make to sychronise
	var waitdata sync.WaitGroup

looper:
	for turn := 0; turn < p.Turns; turn++ {

		///to implement timer, keypressed and single thread
		select {
		case key := <-keyPresses:
			switch unicode.ToLower(key) {
			case 's':

				c.ioCommand <- ioOutput
				c.filename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight), strconv.Itoa(p.Turns)}, "x")

				for r := 0; r < p.ImageHeight; r++ {

					///iterate through x axis
					for s := 0; s < p.ImageWidth; s++ {

						c.output <- world[r][s]

					}

				}

			case 'p':

				state := StateChange{
					CompletedTurns: turn,
					NewState:       Paused,
				}
				c.events <- state
				//using to pause
				waitdata.Add(1)
				go func() {
				loop:
					for {
						select {
						case k := <-keyPresses:
							switch unicode.ToLower(k) {
							case 'p':
								state := StateChange{
									CompletedTurns: turn,
									NewState:       Executing,
								}
								c.events <- state

								fmt.Println("Continuing...")
								waitdata.Done()
								///finish pause
								break loop

							default:
							}
						}
					}
				}()
				waitdata.Wait()
			case 'q':
				state := StateChange{
					CompletedTurns: turn,
					NewState:       Quitting,
				}
				c.events <- state
				fmt.Println("quit")
				break looper
			default:
			}
			turn--
			///send to every worker
		case <-ticker.C:

			cellcounter := 0

			for y := 0; y < p.ImageHeight; y++ {
				for x := 0; x < p.ImageWidth; x++ {
					if world[y][x] != 0 {
						cellcounter++
					}
				}
			}
			alivecell := AliveCellsCount{
				CompletedTurns: turn,
				CellsCount:     cellcounter,
			}
			c.events <- alivecell
			turn--

			continue
			//signal for worker to work
		default:
			///iterate through y axis
			for l := 0; l < p.ImageHeight; l++ {

				///iterate through x axis
				for m := 0; m < p.ImageWidth; m++ {

					///var to store th no. of cells alive in the neighbourhood
					cellneighbours := 0

					///check neighbours of cell
					///iterate through y axis of neighbour
					for n := -1; n < 2; n++ {

						///iterate through x axis neighbour
						for q := -1; q < 2; q++ {

							///check how many alive cells around cell
							///ignore current cell
							if l+n == l && m+q == m {
								continue

								///account for edge cells, for example cell (0,0)
							} else if world[((l+n)+p.ImageHeight)%p.ImageHeight][((m+q)+p.ImageWidth)%p.ImageWidth] == 0xFF {

								///add to var iff alive cell in the neighbourhood
								cellneighbours++
							}

						}

					}

					///turn cell depending on rules
					/// kill cell if live cell has < 2 or >3 live neighbours
					if world[l][m] == 0xFF {
						if cellneighbours < 2 || cellneighbours > 3 {
							cellstoflip = append(cellstoflip, util.Cell{X: m, Y: l})
						}
					}
					/// if cell dead with 3 neighbours make live
					if world[l][m] == 0x00 {
						if cellneighbours == 3 {
							cellstoflip = append(cellstoflip, util.Cell{X: m, Y: l})
						}
					}

				}
			}

			///loop through array containing cells to flip and flip in the world slice
			for _, cell := range cellstoflip {
				world[cell.Y][cell.X] = world[cell.Y][cell.X] ^ 0xFF

				///initialise the cell flipped struct
				flip := CellFlipped{
					CompletedTurns: turn + 1,
					Cell:           util.Cell{X: cell.X, Y: cell.Y},
				}

				///send the cell flip event into channel
				c.events <- flip
			}
			///reset the array
			cellstoflip = nil

			///initialise turncomplete event
			completedturn := TurnComplete{
				CompletedTurns: turn + 1,
			}

			///send the complete turn event into event channel
			c.events <- completedturn
		}

	}

	// TODO: Send correct Events when required, e.g.
	//CellFlipped, TurnComplete and FinalTurnComplete.
	//See event.go for a list of all events.

	///initalise arrray to contain  cells alive
	var finalcellsalive []util.Cell

	c.ioCommand <- ioOutput
	c.filename <- strings.Join([]string{strconv.Itoa(p.ImageWidth), strconv.Itoa(p.ImageHeight), strconv.Itoa(p.Turns)}, "x")

	///iterate through y axis
	for r := 0; r < p.ImageHeight; r++ {

		///iterate through x axis
		for s := 0; s < p.ImageWidth; s++ {

			c.output <- world[r][s]

			///check if cell live then add to array
			if world[r][s] == 0xFF {
				finalcellsalive = append(finalcellsalive, util.Cell{X: s, Y: r})
			}

		}

	}

	///initialise final turn event
	finalturn := FinalTurnComplete{
		CompletedTurns: turn,
		Alive:          finalcellsalive,
	}

	///send the event into event channe;
	c.events <- finalturn

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}
	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
