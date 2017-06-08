package main
import (
    "fmt"
    "os"
    "strings"
    "bufio"
    "path/filepath"
    "strconv"
)

type dataset struct {
     title string
     path_to_file string
     json_data []map[string]interface{}
}

type uistate struct {
     phase int
     active_set *dataset
     key string
     value string
     datasets []dataset
     scanner *bufio.Scanner
}

func enter_interactive_loop(directory string) {
    scanner := bufio.NewScanner(os.Stdin)

    // TODO: Write something appropriate here
    println("Searchy searchy")

    datasets := find_datasets(directory)

    for {
        key, value, set := request_search_fields(datasets, scanner)
        if set != nil && set.json_data != nil {
            search_json(set.json_data, key, value)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Println(os.Stderr, "error:", err)
        os.Exit(1)
    }
}

func find_datasets(directory string) []dataset {

    // Assumes no-one will be malicious enough to create a direcory with a '.json' suffix
    path_elements := []string { directory, "*.json" }
    files, _ := filepath.Glob(strings.Join(path_elements, "/"))

    datasets := make([]dataset, 0)

    for _, path_to_file := range files {
        elements := strings.Split(path_to_file, "/")
        filename := elements[len(elements)-1]
        set := dataset{strings.Title(strings.TrimSuffix(filename, ".json")), path_to_file, nil}
        datasets = append(datasets, set)
    }
    return datasets
}

func unpack_dataset(set *dataset) bool {
    if set.json_data == nil {
        (*set).json_data = parse_file(set.path_to_file)
    }
    if set.json_data != nil {
        return true
    }    
    return false
}

func select_dataset(state uistate) uistate {
    fmt.Println("\nPlease select a dataset to search, or 'quit' to exit:")
    fmt.Printf("   ")
    for index, set := range state.datasets {
        if strings.Compare(set.title, "") != 0 {
            fmt.Printf("%v) %v ", index+1, set.title)
        }
    }
    fmt.Printf("\n# ")

    state.scanner.Scan()
    user_input := state.scanner.Text()
    if strings.Compare("quit", strings.ToLower(user_input)) == 0 {
        state.phase = -1
	return state
    }
    
    index, _ := strconv.ParseInt(user_input, 10, 32)
    if index > 0 && index <= int64(len(state.datasets)) {
        if unpack_dataset(&(state.datasets[index-1])) {
            state.active_set = &(state.datasets[index-1])
	    state.phase += 1
	    
        } else {
	    badset := state.datasets[index-1]
            fmt.Printf("Data source %v is corrupted. Please choose again.\n", badset.title)
            badset.title = ""
        }

    } else {
        fmt.Printf("Invalid selection: '%v'.\n", user_input)
    }

    return state
}

func select_field(state uistate) uistate {
    fmt.Printf("\nEnter a term to search for, '?' to see available fields, or '..' to go back\n%v # ", state.active_set.title)
    state.scanner.Scan()
    user_input := state.scanner.Text()
    if len(user_input) == 0 {
        // TODO: Could we ever wish to search all fields for a specific value?
        fmt.Printf("Invalid selection", user_input)
        
    } else if strings.Compare(user_input, "quit") == 0 {
        state.phase = -1

    } else if strings.Compare(user_input, "..") == 0 {
        state.phase -= 1

    } else if strings.Compare(user_input, "?") == 0 {
        if len(state.active_set.json_data) > 0 {
            // Assume for now that records are sufficiently uniform
                  fmt.Printf("\n%v records contain the following fields\n", strings.TrimSuffix(state.active_set.title, "s"))
            for key, _:= range state.active_set.json_data[0] { 
                  fmt.Printf("* %s\n", key)
            }
        } else {
            fmt.Printf("* No records found *\n")
        }

    } else {
        state.key = user_input
        state.phase += 1
    }

    return state
}

func select_value(state uistate) uistate {
    fmt.Printf("\nEnter a value to search for, '?' to see an example value, or '..' to go back\n%v[%v] # ", state.active_set.title, state.key)
    state.scanner.Scan()
    user_input := state.scanner.Text()
    if len(user_input) == 0 {
        // TODO: Add a confirmation here?
        fmt.Printf("Searching for empty '%v' fields", state.key)

    } else if strings.Compare(user_input, "quit") == 0 {
        state.phase = -1

    } else if strings.Compare(user_input, "?") == 0 {
        if len(state.active_set.json_data) > 0 {
            fmt.Printf("\n%v records contain the following fields\n", strings.TrimSuffix(state.active_set.title, "s"))
		
            // TODO: Handle complex types (arrays)
            for index, record := range state.active_set.json_data { 
                fmt.Printf("* %v\n", record[state.key])
                if index > 2 {
                        break
                }
            }
        } else {
            fmt.Printf("* No records found *\n")
        }

    } else if strings.Compare(user_input, "..") == 0 {
        state.phase -= 1

    } else {
        state.value = user_input
        state.phase += 1
    }
    return state
}


// TODO: Handle ctrl-d
func request_search_fields(datasets []dataset, scanner *bufio.Scanner) (string, string, *dataset) {

    state := uistate{0, nil, "", "", datasets, scanner}

    for {
        switch state.phase {
            case 0:
	        state = select_dataset(state)
	    case 1:
	        state = select_field(state)
	    case 2:
	        state = select_value(state)
	    case 3:
	        return state.key, state.value, state.active_set
	    case -1:
                os.Exit(0)
            default:
	        state.phase = 0
	}
    }
	       
}

