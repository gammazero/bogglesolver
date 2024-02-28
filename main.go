// Interactive command-line application to solve Boggle grids of any size.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gammazero/bogglesolver/solver"
)

func main() {
	var grid string
	xLen := flag.Int("x", 4, "width (X-length) of board")
	yLen := flag.Int("y", 4, "height (Y-length) of board")
	flag.StringVar(&grid, "grid", "", "populate grid with these characters (X*Y length) and exit")
	random := flag.Bool("rand", false, "populate grid with randomly generated characters and exit")
	quiet := flag.Bool("q", false, "do not display grid")
	veryQuiet := flag.Bool("qq", false, "do not display grid or solutions")
	words := flag.String("words", "", "optional file containing valid words separated by newline")
	flag.Parse()

	fmt.Printf("board size (X=%d Y=%d): %d\n", *xLen, *yLen, *xLen**yLen)
	if grid != "" {
		fmt.Println("grid:", grid)
	}
	fmt.Println("rand:", *random)
	fmt.Println("quiet:", *quiet)
	fmt.Println("veryQuiet:", *veryQuiet)
	fmt.Println("words file:", *words)

	var quietLevel int
	if *veryQuiet {
		quietLevel = 2
	} else if *quiet {
		quietLevel = 1
	}

	err := runBoard(grid, *words, *xLen, *yLen, quietLevel, *random)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// runBoard loops getting grid data and finding solutions for that grid.
func runBoard(grid, wordsFile string, xlen, ylen, quietLevel int, random bool) error {
	sol, err := solver.New(xlen, ylen, wordsFile)
	if err != nil {
		return err
	}
	if random {
		grid = randomGrid(sol.BoardSize())
	}
	ever := true
	boardSize := sol.BoardSize()
	for ever {
		if grid == "" {
			grid, err = readGridFromUser(boardSize)
			if err != nil {
				return err
			}
			if grid == "" {
				break
			}
		} else {
			ever = false
		}

		start := time.Now()
		words, err := sol.Solve(grid)
		if err != nil {
			return err
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
		grid = ""
	}
	return nil
}

var rnd *rand.Rand

func randomGrid(size int) string {
	if rnd == nil {
		rnd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	}

	const a = 97
	grid := make([]byte, size)
	for i := 0; i < size; i++ {
		n := rnd.Intn(25)
		grid[i] = byte(a + n)
	}
	return string(grid)
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
func readGridFromUser(boardSize int) (string, error) {
	consReader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nEnter %d letters into boggle grid or * for random: ", boardSize)
	var grid string
	var valid bool
	for {
		input, err := consReader.ReadString('\n')
		if err != nil {
			return "", errors.New("error reading input")
		}
		input = strings.TrimRight(input, "\n")
		if len(input) == 0 {
			return "", nil
		}
		if len(input) == 1 && strings.HasPrefix(input, "*") {
			return randomGrid(boardSize), nil
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
		grid = grid[:boardSize]
	}

	return grid, nil
}
