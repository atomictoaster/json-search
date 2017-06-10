package main
import (
    "os"
    "fmt"
    "flag"
)

func main() {
    datadir := flag.String("data", "./", "Path to a directory containing valid JSON files")
    filename := flag.String("filename", "", "File containing valid JSON")
    search_key := flag.String("key", "", "JSON key to search for")
    search_value := flag.String("value", "", "Value which the specified key must contain")

    flag.Parse()
    if len(*filename) > 0 && len(*search_key) > 0 {
        results := search_file(*filename, *search_key, *search_value)
        for index , record := range results {
            fmt.Printf("\n**** Result[%v/%v]\n", index+1, len(results))
            print_record(record)
        }
        if len(results) > 0 {
           os.Exit(0)
        }
        os.Exit(1) // No result found
    }

    enter_interactive_loop(*datadir)
}
