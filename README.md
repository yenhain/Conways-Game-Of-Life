# Conways Game Of Life Project part of a CS module

# Game Rules 

The Game of Life is a cellular automaton devised by John Horton Conway.
The game is built on a 2D matrix where pixels can be 'alive' or 'dead'. The game relies solely on the starting state and no inputs.
Each cell in the grid interacts with its eight surrounding neighbours—the cells directly adjacent horizontally, vertically, and diagonally. When the grid updates for the next time step, the following changes happen for each cell:

- A living cell with fewer than two living neighbours dies, simulating underpopulation.
- A living cell with two or three living neighbours remains alive without change.
- A living cell with more than three living neighbours dies, simulating overcrowding.
- A dead cell with exactly three living neighbours becomes alive, simulating reproduction.

- Key p is to pause
- Key s is to save the current sdl window
- Key q is to quit

# Output

Using the SDL, the automaton is visualised and creates interesting patterns.

# Example
[![Watch the demo video](https://github.com/yenhain/Conways-Game-Of-Life/raw/main/Screenshot%20from%202025-08-05%2018-24-19.png)](https://github.com/yenhain/Conways-Game-Of-Life/raw/main/Screencast%20from%202025-08-05%2017-48-49.mp4)
