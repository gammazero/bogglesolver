package solver

import (
	"fmt"
	"testing"
)

const wordsFile = "boggle_dict.txt.gz"

func TestLoadWords(t *testing.T) {
	rt, err := loadWords("_not_here_", 16, 3)
	if err == nil {
		t.Fatal("failed to catch bad file")
	}

	// Load from embedded file.
	rt, err = loadWords("", 16, 3)
	if rt == nil {
		t.Fatal("expected trie")
	}
	if rt.Len() < 1 {
		t.Fatal("expected more words")
	}
	fmt.Println("Loaded", rt.Len(), "words from embedded dictionary")

	// Load from external file.
	rt, err = loadWords("", 16, 3)
	if rt == nil {
		t.Fatal("expected trie")
	}
	if rt.Len() < 1 {
		t.Fatal("expected more words")
	}
	fmt.Println("Loaded", rt.Len(), "words from", wordsFile)
}

func TestCalcAdjacency(t *testing.T) {
	// Test corners
	sq := 0
	adj = calculateAdjacency(4, 4, sq)
	//fmt.Println("adj:", adj)
	if len(adj) != 3 || adj[0] != 1 || adj[1] != 4 || adj[2] != 5 {
		t.Error("wrong adjacency for square", sq)
	}

	sq = 3
	adj = calculateAdjacency(4, 4, sq)
	//fmt.Println("adj:", adj)
	if len(adj) != 3 || adj[0] != 2 || adj[1] != 6 || adj[2] != 7 {
		t.Error("wrong adjacency for square", sq)
	}

	sq = 12
	adj = calculateAdjacency(4, 4, sq)
	//fmt.Println("adj:", adj)
	if len(adj) != 3 || adj[0] != 8 || adj[1] != 9 || adj[2] != 13 {
		t.Error("wrong adjacency for square", sq)
	}

	sq = 15
	adj = calculateAdjacency(4, 4, sq)
	//fmt.Println("adj:", adj)
	if len(adj) != 3 || adj[0] != 10 || adj[1] != 11 || adj[2] != 14 {
		t.Error("wrong adjacency for square", sq)
	}

	// Test edge
	sq = 1
	adj = calculateAdjacency(4, 4, sq)
	//fmt.Println("adj:", adj)
	if len(adj) != 5 || adj[0] != 0 || adj[1] != 2 || adj[2] != 4 || adj[3] != 5 || adj[4] != 6 {
		t.Error("wrong adjacency for square", sq)
	}

	// Test center
	sq = 5
	adj = calculateAdjacency(4, 4, sq)
	//fmt.Println("adj:", adj)
	if len(adj) != 8 || adj[0] != 0 || adj[1] != 1 || adj[2] != 2 || adj[3] != 4 || adj[4] != 6 || adj[5] != 8 || adj[6] != 9 || adj[7] != 10 {
		t.Error("wrong adjacency for square", sq)
	}

}

func TestUniqueSortedWords(t *testing.T) {
	words := []string{"gamma", "delta", "alpha", "beta", "zeta", "delta", "delta"}
	usw := uniqueSortedWords(words)
	if len(usw) != 5 {
		t.Fatal("wrong number of unique words")
	}

	for i, w := range []string{"alpha", "beta", "delta", "gamma", "zeta"} {
		if w != usw[i] {
			t.Fatal("words not sorted")
		}
	}
}

func TestSolverBadNew(t *testing.T) {
	_, err := New(4, 5, "_not_here_")
	if err == nil {
		t.Fatal("failed to catch bad file")
	}

	_, err = New(-4, 5, "")
	if err == nil {
		t.Fatal("failed to catch negative dimension")
	}

	_, err = New(4, 0, "")
	if err == nil {
		t.Fatal("failed to catch zero dimension")
	}
}

func TestGrid(t *testing.T) {
	gs := GridString("abcdefghi", 3, 3)
	expect := "+---+---+---+\n" +
		"| A | B | C |\n" +
		"+---+---+---+\n" +
		"| D | E | F |\n" +
		"+---+---+---+\n" +
		"| G | H | I |\n" +
		"+---+---+---+\n"
	if gs != expect {
		t.Error("did not get expected grid string")
	}
}

func TestSolver(t *testing.T) {
	s, err := New(4, 5, "")
	if err != nil {
		t.Fatal(err)
	}

	xlen, ylen := s.Dimensions()
	if xlen != 4 || ylen != 5 {
		t.Fatal("incorrect board dimensions")
	}

	if s.BoardSize() != xlen*ylen {
		t.Fatal("incorrect board size")
	}

	if s.WordCount() < 1 {
		t.Fatal("expected more words")
	}

	grid := "qadfetriihkriflv"
	words, err := s.Solve(grid)
	if err == nil {
		t.Error("failed to catch missing letters")
	}

	grid = "qadfetriihkriflvqadfetriihkriflv"
	words, err = s.Solve(grid)
	if err == nil {
		t.Error("failed to catch too many letters")
	}

	grid = "qadfetriihkriflvctor"
	words, err = s.Solve(grid)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Found %d solutions for %dx%d grid:\n", len(words), xlen, ylen)
	fmt.Println(s.Grid(grid))
	if len(words) != 80 {
		t.Fatal("wrong number of solutions")
	}
	for _, w := range words {
		fmt.Print(w, " ")
	}
	fmt.Println("")
}

func genGrid(boardSize int) string {
	var c rune
	sbgrid := make([]rune, 0, boardSize)
	for i := 0; i < boardSize; i++ {
		if c == 26 {
			c = 0
		}
		sbgrid = append(sbgrid, 'a'+c)
		c++
	}
	return string(sbgrid)
}

func BenchmarkSolver(b *testing.B) {
	const xlen = 50
	const ylen = 50
	s, _ := New(xlen, ylen, "")
	grid := genGrid(s.BoardSize())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Solve(grid)
	}
}

func BenchmarkSolverPrecomp(b *testing.B) {
	const xlen = 50
	const ylen = 50
	s, _ := New(xlen, ylen, "")
	grid := genGrid(s.BoardSize())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Solve(grid)
	}
}
