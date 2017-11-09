package main
import "fmt"

type Cell struct {
    value string
    possibles string
    solved bool
}

var board [81]Cell

func parse(s string) {
  for i := 0; i < len(s); i++ {
    c := string(s[i])
    if c == "." {
      board[i].possibles = "123456789"
      board[i].solved = false
    } else {
      board[i].value = c
      board[i].solved = true
    }
  }
}

func printBoard() {
  for i := 0; i < len(board); i++ {
    if board[i].solved {
      fmt.Print(board[i].value)
    } else {
      fmt.Print(".")
    }

    if i > 0 {
      j := i + 1
      if j % 9 == 0 {
        fmt.Println()
        if j == 27 || j == 54 {
          fmt.Println("===|===|===")
        }
      } else if j % 3 == 0 {
        fmt.Print("|")
      }
    }
  }
  fmt.Println()
}

func main() {
    parse("..6..7..8..1.3....25......9..7.58...9.......1...14.7..8......16....9.4..4..5..8..")
    printBoard()
}
