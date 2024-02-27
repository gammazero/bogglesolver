/*
Find words in X by Y boggle grids and display results interactively.

This script uses the solver module to generate solutions to the boggle grids
entered by a user.  The solver's internal dictionary is created once when the
object is initialized.  It is then reused for subsequent solution searches.

The user is prompted to input a string of x*y characters, representing the
letters in a X by Y Boggle grid.  Use the letter 'q' to represent "qu".

For example: "qadfetriihkriflv" represents the 4x4 grid:
+---+---+---+---+
| Qu| A | D | F |
+---+---+---+---+
| E | T | R | I |
+---+---+---+---+
| I | H | K | R |
+---+---+---+---+
| I | F | L | V |
+---+---+---+---+

This grid has 62 unique solutions using the default dictionary.

Display help to see usage infomation: boggle --help
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gammazero/bogglesolver/solver"
)

const (
	defaultWords = "boggle_dict.txt.gz"
)

// runBoard loops getting grid data and finding solutions for that grid.
func runBoard(wordsFile string, xlen, ylen, quietLevel int, bench bool) {
	sol, err := solver.New(xlen, ylen, wordsFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	boardSize := sol.BoardSize()

	var grid string
	var c rune
	for {
		if bench {
			sbgrid := make([]rune, 0, boardSize)
			for i := 0; i < boardSize; i++ {
				if c == 26 {
					c = 0
				}
				sbgrid = append(sbgrid, 'a'+c)
				c++
			}
			grid = string(sbgrid)
		} else {
			grid = readGridFromUser(boardSize)
		}

		if grid == "" {
			break
		}

		start := time.Now()
		words, err := sol.Solve(grid)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		elapsed := time.Since(start)

		if len(words) == 0 {
			continue
		}

		fmt.Printf("\nFound %d solutions for %dx%d grid in %s\n", len(words), xlen, ylen, elapsed)

		if quietLevel < 2 {
			if quietLevel < 1 {
				fmt.Print(sol.Grid(grid))
			}
			showWords(words)
		}

		if bench {
			time.Sleep(time.Second)
		}
	}
}

// showWords prints words in four columns.
func showWords(words []string) {
	// Sort words by lenght.
	sort.Slice(words, func(i, j int) bool { return len(words[i]) > len(words[j]) })
	for i, w := range words {
		if i%4 == 0 {
			fmt.Println("")
		}
		fmt.Printf("%-18s", w)
	}
	fmt.Println("")
}

// readGridFromUser reads input from user, rejecting invalid characters.
func readGridFromUser(boardSize int) string {
	consReader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nEnter %d letters into boggle grid: ", boardSize)
	var grid string
	var valid bool
	for {
		input, err := consReader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading input")
			os.Exit(1)
		}
		input = strings.TrimRight(input, "\n")
		if len(input) == 0 {
			return ""
		}
		input = strings.ToLower(input)
		valid = true
		for _, c := range input {
			if c < 'a' || c > 'z' {
				fmt.Fprintln(os.Stderr, "input contains invalid cahracters")
				valid = false
				break
			}
		}
		if valid {
			grid = grid + input
		}
		if len(grid) >= boardSize {
			break
		}
		fmt.Printf("\n%d more letters needed: ", boardSize-len(grid))
	}

	if len(grid) > boardSize {
		grid = grid[0:boardSize]
	}

	return grid
}

func main() {
	var xLen = flag.Int("x", 4, "width (X-length) of board")
	var yLen = flag.Int("y", 4, "height (Y-length) of board")
	var bench = flag.Bool("b", false, "run benchmark test")
	var quiet = flag.Bool("q", false, "do not display grid")
	var veryQuiet = flag.Bool("qq", false, "do not display grid or solutions")
	var words = flag.String("words", defaultWords,
		"file containing valid words, separated by newline")
	var grid = flag.String("grid", "", "grid letters (must be X*Y length)")
	flag.Parse()

	fmt.Printf("board size (X=%d Y=%d): %d\n", *xLen, *yLen, *xLen**yLen)
	fmt.Println("grid:", *grid)
	fmt.Println("words file:", *words)
	fmt.Println("bench:", *bench)
	fmt.Println("quiet:", *quiet)
	fmt.Println("veryQuiet:", *veryQuiet)

	var quietLevel int
	if *veryQuiet {
		quietLevel = 2
	} else if *quiet {
		quietLevel = 1
	}

	runBoard(*words, *xLen, *yLen, quietLevel, *bench)

}
