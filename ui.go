package main
import (
    "fmt"
    "os"
    "strings"
    "bufio"
    "path/filepath"
    "strconv"
)

var json_data []map[string]interface{}

func enter_interactive_loop(directory string) {
    last_filename := ""

    scanner := bufio.NewScanner(os.Stdin)
    println("Searchy searchy")
    for {
        filename, key, value := request_search_fields(directory, scanner)
        if strings.Compare(filename, last_filename) != 0 {
       	  // The lowest of low bars...
	  // If the dataset didn't change, don't waste cycles
            json_data = parse_file(filename)
  	  last_filename = filename
        }
        search_json(json_data, key, value)
    }
    if err := scanner.Err(); err != nil {
        fmt.Println(os.Stderr, "error:", err)
        os.Exit(1)
    }
}

type dataset struct {
     title string
     file_with_path string
}

func find_datasets(directory string) []dataset {

    // Assumes no-one will be malicious enough to create a direcory with a '.json' suffix
    path_elements := []string { directory, "*.json" }
    files, _ := filepath.Glob(strings.Join(path_elements, "/"))

    datasets := make([]dataset, 0)

    for _, path_to_file := range files {
    	elements := strings.Split(path_to_file, "/")
	filename := elements[len(elements)-1]
	set := dataset{strings.Title(strings.TrimSuffix(filename, ".json")), path_to_file}
	datasets = append(datasets, set)
    }
    return datasets
}

func request_search_fields(directory string, scanner *bufio.Scanner) (string, string, string) {

    title := ""
    filename := ""
    search_key := ""
    search_value := ""

    datasets := find_datasets(directory)
    for len(filename) == 0 {
        fmt.Println("\nPlease select a dataset to search, or 'quit' to exit:")
	fmt.Printf("   ")
        for index, set := range datasets {
            fmt.Printf("%v) %v ", index+1, set.title)
        }
        fmt.Printf("\n# ")

        scanner.Scan()
        user_input := scanner.Text()
        if strings.Compare("quit", strings.ToLower(user_input)) == 0 {
	    os.Exit(0)
	}

	index, _ := strconv.ParseInt(user_input, 10, 32)
	if index > 0 && index <= int64(len(datasets)) {
	    title = datasets[index-1].title
	    filename = datasets[index-1].file_with_path
	} else {
            fmt.Printf("Invalid selection: '%v'. Try again", user_input)
        }
    }

    for len(search_key) == 0 {
        fmt.Printf("\nEnter a term to search for, or '?' to see available fields\n%v # ", title)
        scanner.Scan()
        user_input := scanner.Text()
        if len(user_input) == 0 {
            // TODO: Could we ever wish to search all fields for a specific value?
            fmt.Printf("Invalid selection: '%v'. Try again", user_input)
        } else if strings.Compare(user_input, "?") == 0 {
            // TODO: Provide access to 'lookup search terms' from here.
        } else {
            search_key = user_input
        }
    }

    fmt.Printf("\nEnter a value to search for, or '?' to see an example value\n%v[%v] # ", title, search_key)
    scanner.Scan()
    user_input := scanner.Text()
    if len(user_input) == 0 {
        // TODO: Add a confirmation here?
        println("Searching for empty '%v' fields", search_key)
    } else if strings.Compare(user_input, "?") == 0 {
        // TODO: Show an example for that field
    } else {
        search_value = user_input
    }
    return filename, search_key, search_value
}

