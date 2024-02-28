// Package solver Find words in X by Y boggle grids and a set of results.
//
// A solver is created with the X and Y dimensions of the grid it searches, and
// is given an optional list of valid words. When given a boggle grid to solve,
// the solver looks in all paths through the grid to find words that are in the
// list of valid words. A set of found words is returned at the end of the
// search.
//
// The boggle grid to search is specified as a string with length equal to the
// grid size, which is the product of the X and Y dimensions that the solver
// was created with. The letter 'q' in the grid always represents "qu" when
// searching for words in the valid words list.
//
// For example: "qadfetriihkriflv" represents the 4x4 grid:
//
//	+---+---+---+---+
//	| Qu| A | D | F |
//	+---+---+---+---+
//	| E | T | R | I |
//	+---+---+---+---+
//	| I | H | K | R |
//	+---+---+---+---+
//	| I | F | L | V |
//	+---+---+---+---+
//
// This grid has 62 unique solutions using the default dictionary.
package solver
