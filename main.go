package main
import (
    "os"
    "flag"
)

func main() {
    filename := flag.String("filename", "./users.json", "File containing valid JSON")
    search_key := flag.String("key", "", "JSON key to search for")
    search_value := flag.String("value", "Catalina Simpson", "Value which the specified key must contain")

    flag.Parse()
    if len(*filename) > 0 && len(*search_key) > 0 {
        if search_file(*filename, *search_key, *search_value) > 0 {
           os.Exit(0)
        }
        os.Exit(1) // No result found
    }

    enter_interactive_loop()
}
