package solver

import (
	"bufio"
	"compress/gzip"
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/gammazero/deque"
	"github.com/gammazero/radixtree"
)

const defaultWords = "boggle_words.txt.gz"

//go:embed boggle_words.txt.gz
var wordsFile embed.FS

var adj = make([]int, 0, 8)

// qNode is a element of the queue constructed while searching word paths.
type qNode struct {
	parentSquare int
	parentTrie   *radixtree.Stepper
	seen         []int
}

// Solver implements the algorithm to find words in the Boggle grid.
//
// Solver searches all paths through a boggle grid, searching for words that
// occur in a given list of acceptable boggle words. The Solve() method can be
// used repeatedly to generate solutions for different boggle grids.
type Solver struct {
	cols int
	rows int
	rt   *radixtree.Tree
}

// New creates and initializes a Solver instance.
//
// This creates the internal trie for fast word lookup letter-by-letter. Words
// that begin with capital letters and words that are not within the specified
// length limits are filtered out.
//
// New takes the board dimensions xlen and ylen, a an optional file which can
// be gz compressed. If no file is specified, then the embedded words list is
// used.
//
// The maximum word length is the size of the board, and the minimum word
// length is 3 letters.
func New(xlen, ylen int, wordsPath string) (Solver, error) {
	if xlen < 1 || ylen < 1 {
		return Solver{}, errors.New("invalid board dimensions")
	}

	rt, err := loadWords(wordsPath, xlen*ylen, 3)
	if err != nil {
		return Solver{}, err
	}

	return Solver{
		cols: xlen,
		rows: ylen,
		rt:   rt,
	}, nil
}

// BoardSize return the size of the board (x * y).
func (s Solver) BoardSize() int {
	return s.cols * s.rows
}

// Dimensions returns the number of columns (x-size) and rows (y-size).
func (s Solver) Dimensions() (int, int) {
	return s.cols, s.rows
}

// WordCount returns the number of words read from the words file.
func (s Solver) WordCount() int {
	return s.rt.Len()
}

// Solve generates all solutions for the given Boggle grid.
//
// The grid argument is a string of X*Y characters, representing the letters in
// a Boggle grid, from top left to bottom right. This method returns a slice of
// the words that were found in the grid.
func (s Solver) Solve(grid string) ([]string, error) {
	if s.rt == nil {
		return nil, errors.New("failed to read words file")
	}
	if len(grid) != s.BoardSize() {
		if len(grid) < s.BoardSize() {
			return nil, errors.New("not enough letters for board")
		}
		return nil, errors.New("too many letters for board")
	}

	board := strings.ToLower(grid)
	words := make([]string, 0, 256)
	q := deque.New[qNode](s.BoardSize(), s.BoardSize())
	for initSq := 0; initSq < len(board); initSq++ {
		seen := make([]int, 1, 8)
		seen[0] = initSq
		stepper := s.rt.NewStepper()
		stepper.Next(board[initSq])
		q.PushBack(qNode{
			parentSquare: initSq,
			parentTrie:   stepper,
			seen:         seen,
		})
		for q.Len() != 0 {
			qn := q.PopFront()
			parentSq := qn.parentSquare
			parentTrie := qn.parentTrie
			seen = qn.seen
			sqAdj := calculateAdjacency(s.cols, s.rows, parentSq)
		AdjLoop:
			for _, curSq := range sqAdj {
				for i := range seen {
					if seen[i] == curSq {
						continue AdjLoop
					}
				}
				curNode := parentTrie.Copy()
				if !curNode.Next(board[curSq]) {
					continue
				}
				newSeen := make([]int, len(seen)+1)
				copy(newSeen, seen)
				newSeen[len(seen)] = curSq

				q.PushBack(qNode{
					parentSquare: curSq,
					parentTrie:   curNode,
					seen:         newSeen,
				})
				if item := curNode.Item(); item != nil {
					key := item.Key()
					if key[0] == 'q' {
						// Rehydrate q-words with 'u'.
						words = append(words, "qu"+key[1:])
					} else {
						words = append(words, key)
					}
				}
			}
		}
	}

	return uniqueSortedWords(words), nil
}

// Grid returns a printable string version of a X by Y boggle grid.
//
// The grid is given as a string of X*Y characters representing the letters in
// a boggle grid, from top left to bottom right.
func (s Solver) Grid(grid string) string {
	return GridString(grid, s.cols, s.rows)
}

func GridString(grid string, cols, rows int) string {
	if len(grid) != cols*rows {
		panic("number of letters in grid must equal cols * rows")
	}
	grid = strings.ToUpper(grid)
	gridChars := []byte(grid)

	line := make([]string, 0, cols+2)
	line = append(line, "")
	for i := 0; i < cols; i++ {
		line = append(line, "---")
	}
	line = append(line, "\n")
	hline := strings.Join(line, "+")

	gridLines := make([]string, 0, 2*rows+1)
	gridLines = append(gridLines, "")
	var yi int
	for y := 0; y < rows; y++ {
		yi = y * cols
		var cell byte
		for x := 0; x < cols; x++ {
			cell = gridChars[yi+x]
			if cell == 'Q' {
				line[1+x] = " Qu"
			} else {
				line[1+x] = fmt.Sprintf(" %c ", cell)
			}
		}
		gridLines = append(gridLines, strings.Join(line, "|"))
	}
	return strings.Join(append(gridLines, ""), hline)
}

// loadWords reads a file of words and creates a trie containing them. If no
// file name is specified then the embedded words list is loaded.
func loadWords(filePath string, maxLen, minLen int) (*radixtree.Tree, error) {
	var rdr io.Reader
	var gz bool
	if filePath == "" {
		f, err := wordsFile.Open(defaultWords)
		if err != nil {
			return nil, fmt.Errorf("solver: error opening words file: %s", err)
		}
		defer f.Close()
		rdr = f
		gz = true
	} else {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("solver: error opening words file: %s", err)
		}
		defer f.Close()
		rdr = f
		gz = strings.HasSuffix(filePath, ".gz")
	}
	if gz {
		var err error
		rdr, err = gzip.NewReader(rdr)
		if err != nil {
			return nil, fmt.Errorf("solver: error unzipping words file: %s", err)
		}
	}

	scanner := bufio.NewScanner(rdr)
	tree := radixtree.New()

	// Scan through line-dilimited words.
	for scanner.Scan() {
		word := scanner.Text()
		// Skip words that are too long or too short.
		if len(word) > maxLen || len(word) < minLen {
			continue
		}
		// Skip words that start with a capital letter.
		if int(word[0]) < 'a' {
			continue
		}
		// If word starts wit qu then remove u so that only q is mathced.
		if int(word[0]) == 'q' {
			// Skip words that start with q not followed by u.
			if int(word[1]) != 'u' {
				continue
			}
			word = "q" + word[2:]
		}

		tree.Put(word, nil)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("solver: error reading words file: %s", err)
	}

	return tree, nil
}

func uniqueSortedWords(words []string) []string {
	if len(words) == 0 {
		return words
	}
	sort.Sort(sort.StringSlice(words))
	unique := make([]string, 0, len(words))
	var prev string
	for _, w := range words {
		if w != prev {
			unique = append(unique, w)
			prev = w
		}
	}
	return unique
}

// calculateAdjacency calculates squares adjacent to the one given.
//
// Adjacent squares, up to eight, are calculated for the square specified by
// the x and y coordinates and are written to the given slice.
func calculateAdjacency(xlim, ylim, sq int) []int {
	// Current cell index = y * xlim + x
	y := sq / xlim
	x := sq - (y * xlim)
	var above, below int

	// Clear the adj slice.
	adj = adj[:0]

	// Look at row above current cell.
	if y-1 >= 0 {
		above = sq - xlim
		// Look to upper left.
		if x-1 >= 0 {
			adj = append(adj, above-1)
		}
		// Look above.
		adj = append(adj, above)
		// Look upper right.
		if x+1 < xlim {
			adj = append(adj, above+1)
		}
	}
	// Look at same row that current cell is on.
	// Look to left of current cell.
	if x-1 >= 0 {
		adj = append(adj, sq-1)
	}
	// Look to right of current cell.
	if x+1 < xlim {
		adj = append(adj, sq+1)
	}
	// Look at row below current cell.
	if y+1 < ylim {
		below = sq + xlim
		// Look to lower left.
		if x-1 >= 0 {
			adj = append(adj, below-1)
		}
		// Look below.
		adj = append(adj, below)
		// Look to lower rigth.
		if x+1 < xlim {
			adj = append(adj, below+1)
		}
	}
	return adj
}
