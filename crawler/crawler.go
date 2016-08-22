package main
import(
  "fmt"
  "flag"
  "os"
)

func main() {
  flag.Parse()
  args := flag.Args()
  if len(args) < 1 {
    fmt.Println("please specify parameter")
    os.Exit(1)
  }
}
