package games

const Tic_Tac_Toe_Template = `
interface Move {
	x: number;
	y: number;
}

/**
At each move, the function gets following parameters.

state:
-------
Example:
x - x 
- 0 -
- - -

Here player 1 has placed x at (0,0) and (0, 2)
And  player 2 has placed 0 at (1, 1)
'-' represents empty slot. 

moves:
------
Possible moves i.e. empty slots list of type Move.

Returns:
--------
Index of one of the possible move. Value should be one of the index in array moves.
**/
function Play(state: string[][], moves: Move[]): number {
	return Math.floor(Math.random() * moves.length);
}

`

const Tic_Tac_Toe_Driver_Code = `
//user function
{{.}}

playerMove = Play(null, [ { x: 1, y: 2 } ]);
console.write(playerMove);
`
