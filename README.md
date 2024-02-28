# bogglesolver
[![Go Reference](https://pkg.go.dev/badge/github.com/gammazero/bogglesolver.svg)](https://pkg.go.dev/github.com/gammazero/bogglesolver)

:abcd: CLI for solving Boggle puzzles of any size

This project provides a command line application for solving any size Boggle puzzles, and package for finding these solutions.

## Install

```shell
go install github.com/gammazero/bogglesolver
```

This installs the `bogglesolver` command into `$GOPATH/bin/`

## Run

To see instructions for use, run:
```
bogglesolver -help
```

### Example

```
> bogglesolver -grid qazwsxedcrfvtgby
Found 33 solutions for 4x4 grid in 78.959µs
+---+---+---+---+
| Qu| A | Z | W |
+---+---+---+---+
| S | X | E | D |
+---+---+---+---+
| C | R | F | V |
+---+---+---+---+
| T | G | B | Y |
+---+---+---+---+

derby             screw             grew              crew              
wert              defy              bred              vert              
verb              tref              axed              brew              
erg               sax               few               fez               
fed               zed               ers               ref               
rev               rex               sae               fer               
red               dex               dew               dev               
vex               wed               axe               zax               
qua               
```

If the `-grid` or `-rand` flag are specified a single solution is output. Otherwise, the user is interactively prompted for input.
