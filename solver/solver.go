// Solver finds all words in X by Y boggle grids.
package solver

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/gammazero/queue"
	"io"
	"os"
	"sort"
	"strings"
)

// qNode is a element of the queue constructed while searching word paths.
type qNode struct {
	parentSquare int
	prefix       string
	parentTrie   *Trie
	seen         []int
}

func newQNode(parentSquare int, prefix string, parentTrie *Trie, seen []int) *qNode {
	return &qNode{parentSquare, prefix, parentTrie, seen}
}

// Solver implements the algorithm to find words in the Boggle grid.
//
// Solver uses an external words file as a dictionary of acceptable boggle
// words.  When an instance of Solver is created, it sets up an internal
// dictionary to look up valid boggle answers.  The Solve() method can be used
// repeatedly to generate solutions for different boggle grids.
type BoggleSolver struct {
	rows, cols, boardSize int
	root                  *Trie
	wordCount             int
	adjacency             [][]int
}

// NewSolver creates and initializes a Solver instance.
//
// This creates the internal trie for fast word lookup letter-by-letter.  Words
// that begin with capital letters and words that are not within the specified
// length limits are filtered out.
//
// NewSolver takes the board dimensions xlen and ylen, a file (optionally gz
// compressed) and flag specifying whether or not to use a pre-calculated
// adjacency matrix (uses more space to save some time).
//
// The maximum word length is the size of the board, and the minimum word
// length is 3 letters.
func NewSolver(xlen, ylen int, wordsFile string, preCalcAdjacency bool) (*BoggleSolver, error) {
	if xlen < 1 || ylen < 1 {
		return nil, errors.New("invalid board dimensions")
	}

	rt, wc, err := loadWords(wordsFile, xlen*ylen, 3)
	if err != nil {
		return nil, err
	}

	solver := BoggleSolver{
		cols:      xlen,
		rows:      ylen,
		boardSize: xlen * ylen,
		root:      rt,
		wordCount: wc}
	if preCalcAdjacency {
		solver.adjacency = calculateAdjacencyMatrix(xlen, ylen)
	}
	return &solver, nil
}

// BoardSize return the size of the board (x * y).
func (s *BoggleSolver) BoardSize() int {
	return s.boardSize
}

// Dimensions returns the number of columns (x-size) and rows (y-size).
func (s *BoggleSolver) Dimensions() (int, int) {
	return s.cols, s.rows
}

// WordCount returns the number of words read from the words file.
func (s *BoggleSolver) WordCount() int {
	return s.wordCount
}

// Solve generates all solutions for the given Boggle grid.
//
// The grid argument is a string of X*Y characters, representing the letters in
// a Boggle grid, from top left to bottom right.  This method returns a slice
// of the words that were found in the grid.
func (s *BoggleSolver) Solve(grid string) ([]string, error) {
	if s.root == nil {
		return nil, errors.New("failed to read words file")
	}
	if len(grid) != s.boardSize {
		if len(grid) < s.boardSize {
			return nil, errors.New("not enough letters for board")
		}
		return nil, errors.New("too many letters for board")
	}

	board := []rune(strings.ToLower(grid))
	trie := s.root
	words := make([]string, 0, 32)
	q := queue.New(s.boardSize)
	adj := make([]int, 0, 8)
	sqAdj := adj
	var adjCount int
	for initSq, c := range board {
		seen := make([]int, 0, 8)
		seen = append(seen, initSq)
		cstr := string(c)
		qn := newQNode(initSq, cstr, trie.Child(c), seen)
		q.Push(qn)
		for !q.Empty() {
			qn = q.Pop().(*qNode)
			parentSq := qn.parentSquare
			prefix := qn.prefix
			parentTrie := qn.parentTrie
			seen = qn.seen
			if s.adjacency == nil {
				sqAdj = calculateAdjacency(s.cols, s.rows, parentSq, adj)
			} else {
				sqAdj = s.adjacency[parentSq]
			}
			adjCount = len(sqAdj)
			for a := 0; a < adjCount; a++ {
				curSq := sqAdj[a]
				hasCur := false
				for _, x := range seen {
					if x == curSq {
						hasCur = true
						break
					}
				}
				if hasCur {
					continue
				}
				c = board[curSq]
				curNode := parentTrie.Child(c)
				if curNode == nil {
					continue
				}
				cstr = prefix + string(c)
				newSeen := make([]int, 0, len(seen)+1)
				newSeen = append(newSeen, seen...)
				newSeen = append(newSeen, curSq)
				newNode := newQNode(curSq, cstr, curNode, newSeen)
				q.Push(newNode)
				if curNode.IsWord() {
					if cstr[0] == 'q' {
						// Rehydrate q-words with 'u'.
						words = append(words, "qu"+cstr[1:])
					} else {
						words = append(words, cstr)
					}
				}
			}
		}
	}

	return uniqueSortedWords(words), nil
}

// GridString returns a printable string version of a X by Y boggle grid.
//
// The grid is given as a string of X*Y characters representing the letters in
// a boggle grid, from top left to bottom right.
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

// loadWords reads a file of words and creates a trie containing them.
func loadWords(wordsFile string, maxLen, minLen int) (*Trie, int, error) {
	f, err := os.Open(wordsFile)
	if err != nil {
		return nil, 0, fmt.Errorf("solver: error opening words file: %s", err)
	}
	defer f.Close()

	var rdr io.Reader
	if strings.HasSuffix(wordsFile, ".gz") {
		rdr, err = gzip.NewReader(f)
		if err != nil {
			return nil, 0, fmt.Errorf("solver: error unzipping words file:", err)
		}
	} else {
		rdr = f
	}
	scanner := bufio.NewScanner(rdr)
	root := NewTrie()
	wordCount := 0
	var word string

	// Scan through line-dilimited words.
	for scanner.Scan() {
		word = scanner.Text()
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

		root.Insert(word)
		wordCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, fmt.Errorf("solver: error reading words file:", err)
	}

	return root, wordCount, nil
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

// caculateAdjacency creates the adjacency matrix for any board dimensions.
//
// An array of adjacent squares, up to eight, is calculated for each square on
// the board.  The dimensions are given by the xlim and ylim parameters.
func calculateAdjacencyMatrix(xlim, ylim int) [][]int {
	// adjList is an array of slices of int.
	adjList := make([][]int, ylim*xlim)
	for sq := 0; sq < xlim*ylim; sq++ {
		// adj holds adjacent squares, up to 8.
		// Store a slice adjacent squares, for each square in board.
		adjList[sq] = calculateAdjacency(xlim, ylim, sq, make([]int, 0, 8))
	}
	return adjList
}

// calculateAdjacency calculates squares adjacent to the one given.
//
// An array of adjacent squares, up to eight, in calculated for the square
// specified by the x and y coordinates, and are written to the given slice.
func calculateAdjacency(xlim, ylim, sq int, adj []int) []int {
	// Current cell index = y * xlim + x
	y := sq / xlim
	x := sq - (y * xlim)
	var above, below int

	// Clear the adj slice.
	adj = adj[0:0]

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
