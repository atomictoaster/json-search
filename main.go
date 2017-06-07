package main
import (
    "fmt"
    "os"
    "strings"
    "flag"
    "bufio"
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

    scanner := bufio.NewScanner(os.Stdin)
    println("Searchy searchy")
    for {
        println("1 to search")
        println("2 to list fields") // TODO: Implement
        println("'quit' to exit")

        scanner.Scan()
        user_input := scanner.Text()

        switch strings.ToLower(user_input) {
           case "quit":
               os.Exit(0)
           case "1":
              filename, key, value := request_search_fields(scanner)
              search_file(filename, key, value)
           case "2":
              println("Not implemented yet")
           default:
              println("Invalid command:", user_input)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Println(os.Stderr, "error:", err)
        os.Exit(1)
    }
}

func request_search_fields(scanner *bufio.Scanner) (string, string, string) {

    filename := ""
    search_key := ""
    search_value := ""
    println("Search 1) Users 2) Tickets 3) Organizations")

    for len(filename) == 0 {
        scanner.Scan()
        user_input := scanner.Text()
        switch strings.ToLower(user_input) {
            case "quit":
                os.Exit(0)
            case "1":
                filename = "users.json"
            case "2":
                filename = "tickets.json"
            case "3":
                filename = "organizations.json"
            default:
                    println("Invalid command: '%v'. Try again", user_input)
        }
    }

    // TODO: Provide access to 'lookup search terms' from here. Eg. if '?' is entered.
    println("Enter a term to search for")
    for len(search_key) == 0 {
        scanner.Scan()
        user_input := scanner.Text()
        if len(user_input) == 0 {
            // TODO: Could we ever wish to search all fields for a specific value?
                println("Invalid command: '%v'. Try again", user_input)
        } else {
            // TODO: All keys seem to be lowercase. Valid assumption?
            search_key = strings.ToLower(user_input)
        }
    }

    println("Enter a value to search for")
    scanner.Scan()
    user_input := scanner.Text()
    if len(user_input) == 0 {
            println("Searching for empty '%v' fields", search_key)
    } else {
        search_value = user_input
    }
    return filename, search_key, search_value
}

