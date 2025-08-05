# Conways Game Of Life Project part of a CS module

# Game Rules 

The Game of Life is a cellular automaton devised by John Horton Conway.
The game is built on a 2D matrix where pixels can be 'alive' or 'dead'. The game relies soley on the starting state and no inputs.
Each cell in the grid interacts with its eight surrounding neighbors—the cells directly adjacent horizontally, vertically, and diagonally. When the grid updates for the next time step, the following changes happen for each cell:

- A living cell with fewer than two living neighbors dies, simulating underpopulation.
- A living cell with two or three living neighbors remains alive without change.
- A living cell with more than three living neighbors dies, simulating overcrowding.
- A dead cell with exactly three living neighbors becomes alive, simulating reproduction.

Key p is to pause
Key s is to save the current sdl window
Key q is to quit

# Output

Using the SDL, the automaton is visualised and creates interesting patterns.

