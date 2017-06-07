package main
import (
    "fmt"
    "io/ioutil"
    "os"
    "encoding/json"
    "strings"
    "flag"
    "bufio"
    "strconv"
    "reflect"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func Filter(vs []map[string]interface{}, f func(map[string]interface{}) bool) []map[string]interface{} {
    vsf := make([]map[string]interface{}, 0)
    for _, v := range vs {
        if f(v) {
            vsf = append(vsf, v)
        }
    }
    return vsf
}

func main() {

    filename := flag.String("filename", "./users.json", "File containing valid JSON")
    search_key := flag.String("key", "", "JSON key to search for")
    search_value := flag.String("value", "Catalina Simpson", "Value which the specified key must contain")

    flag.Parse()
    if len(*filename) > 0 && len(*search_key) > 0 {
        if search(*filename, *search_key, *search_value) > 0 {
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
              search(filename, key, value)
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

func search(filename string, search_key string, search_value string) int {
    // TODO: Test with really big data sets to see if its worth caching the file contents.
    //       Assumes repeated searches of the same dataset in interactive mode 

    println("Reading", filename)
    bytes, err := ioutil.ReadFile(filename)
    check(err)

    println("Parsing", filename)

    // We expect an array, of maps, of strings to $something
    var jdat []map[string]interface{}

    // TODO: Test behaviour of parser with malformed and "valid json,
    // but not arranged like this" input
    if err := json.Unmarshal(bytes, &jdat); err != nil {
        panic(err)
    }

    println("Filtering", len(jdat), "records based on", search_key, "=", search_value)
    results := Filter(jdat, func(v map[string]interface{}) bool {
        value := v[search_key]
        if value == nil && len(search_value) == 0 {
            // Searching for empty values
            return true
        } else if value == nil {
            return false
        }
        string_value := ""

        switch value.(type) {
        case bool:
            string_value = strconv.FormatBool(value.(bool));
        case int64:
            string_value = strconv.FormatInt(value.(int64), 10);
        case float64:
            string_value = strconv.FormatFloat(value.(float64), 'f', 0, 64);
        case string:
            string_value = value.(string)
	case []interface {}:
	    // TODO: Handle arrays
	    return false
        default:
            unhandled := reflect.Indirect(reflect.ValueOf(value))
            fmt.Println("Unhandled type", unhandled.Kind(), "in filter for key", search_key)
            return false
        }

        if strings.Compare(string_value, search_value) == 0 {
            return true
        }
        return false
    })
    num_results := len(results)
    if num_results > 0 {
        enc := json.NewEncoder(os.Stdout)
        enc.Encode(results)
    }
    println(num_results, "record(s) found")
    return num_results
}
