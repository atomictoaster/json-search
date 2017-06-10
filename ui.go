package main
import (
    "fmt"
    "os"
    "strings"
    "bufio"
    "path/filepath"
    "strconv"
    "sort"
)

type dataset struct {
     title string
     path_to_file string
     json_data []map[string]interface{}
}

type ui_search_state struct {
     phase int
     active_set *dataset
     key string
     value string
     datasets []dataset
     scanner *bufio.Scanner
}

var command_quit string = ".quit"
var command_back string = ".back"
var command_help string = ".help"
var command_done string = ".done"
var command_more string = ".more"

var command_short_quit string = ".q"
var command_short_back string = ".b"
var command_short_help string = ".h"
var command_short_done string = ".d"
var command_short_more string = ".m"

func scan_or_exit(scanner *bufio.Scanner) string {
    if scanner.Scan() == false && scanner.Err() == nil {
        os.Exit(0)
    }
    return scanner.Text()
}

// A simple, single dimension, finite state machine for handling the
// gathering of user input.
//
// Allows navigation between the various query phases and
// re-use of previously entered values - avoiding the need to re-enter
// the dataset and key each time which might be useful when
// searching for customers with hard-to-spell names .
//
func enter_interactive_loop(directory string) {
    scanner := bufio.NewScanner(os.Stdin)

    println("JSON Search Tool")

    datasets := find_datasets(directory)
    state := ui_search_state{0, nil, "", "", datasets, scanner}

    for {
        switch state.phase {
            case 0:
	        state = select_dataset(state)
	    case 1:
	        state = select_field(state)
	    case 2:
	        state = select_value(state)
	    case 3:
	        state = perform_search(state)
	    case -1:
                os.Exit(0)
            default:
	        state.phase = 0
	}
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
        valid, json_data := parse_file(set.path_to_file)
	if valid {
            (*set).json_data = json_data
	}
    }
    if set.json_data != nil {
        return true
    }    
    return false
}

func select_dataset(state ui_search_state) ui_search_state {
    fmt.Printf("\nPlease select a dataset to search, or '%v' to exit:\n", command_quit)
    fmt.Printf("   ")
    for index, set := range state.datasets {
        fmt.Printf("%v) %v ", index+1, set.title)
    }
    fmt.Printf("\n# ")

    user_input := scan_or_exit(state.scanner)
    fmt.Printf("\n")

    switch user_input {
        case command_quit, command_short_quit:
            state.phase = -1
        default:
            index, _ := strconv.ParseInt(user_input, 10, 32)
            if index > 0 && index <= int64(len(state.datasets)) {
                if unpack_dataset(&(state.datasets[index-1])) {
                    state.active_set = &(state.datasets[index-1])
        	    state.phase += 1
	    
                } else {
                    fmt.Printf("Data source %v is corrupted. Please choose again.\n", state.datasets[index-1].title)
                }

            } else {
                fmt.Printf("Invalid selection: '%v'.  Enter the number of a dataset instead.\n", user_input)
            }
    }
    return state
}

func prompt_for_input(prompt string, help string) {
    fmt.Printf("\n   '%v|%v' to see %v\n", command_help, command_short_help, help)
    fmt.Printf("   '%v|%v' to go back\n", command_back, command_short_back)
    fmt.Printf("   '%v|%v' to exit\n", command_quit, command_short_quit)
    fmt.Printf("\nEnter a %v to search ", prompt)
}


// NO attempt at field validation is attempted least the record fields vary.
// You gets what you ask for
func select_field(state ui_search_state) ui_search_state {
    prompt_for_input("term", "available fields")
    fmt.Printf("%v # ", state.active_set.title)

    user_input := scan_or_exit(state.scanner)
    switch user_input {
        case "":
            fmt.Printf("You must make a selection to continue\n")
        case command_quit, command_short_quit:
            state.phase = -1
        case command_back, command_short_back:
            state.phase -= 1
        case command_help, command_short_help:
            if len(state.active_set.json_data) > 0 {
                // Assume for now that records are sufficiently uniform
                fmt.Printf("\n%v records contain the following fields:\n", strings.TrimSuffix(state.active_set.title, "s"))
                keys := make([]string, 0)
                for key, _:= range state.active_set.json_data[0] { 
                    keys = append(keys, key)
                }
                sort.Strings(keys)
                for _, key := range keys {
                    fmt.Printf("* %s\n", key)
                }
            } else {
                fmt.Printf("* No records found *\n")
            }
        default:
            state.key = user_input
            state.phase += 1
    }

    return state
}

func select_value(state ui_search_state) ui_search_state {
    prompt_for_input("value", "example values")
    fmt.Printf("%v[%v] # ", state.active_set.title, state.key)

    user_input := scan_or_exit(state.scanner)
    switch user_input {
        case command_quit, command_short_quit:
            state.phase = -1
        case "":
            fmt.Printf("Searching for empty '%v' fields\n", state.key)
        case command_back, command_short_back:
            state.phase -= 1
        case command_help, command_short_help:
            if len(state.active_set.json_data) > 0 {
                fmt.Printf("\n%v records contain '%v' values like:\n", strings.TrimSuffix(state.active_set.title, "s"), state.key)
        	
                // TODO: Nicer handling complex types (arrays) 
                for index, record := range state.active_set.json_data {
                    switch typed_value := record[state.key].(type) {
                    case []interface{}:
                        for _, entry := range typed_value {
                            fmt.Printf("* %v\n", entry)
                        }
                    default:
                        fmt.Printf("* %v\n", typed_value)
                    }

                    if index > 2 {
                        break
                    }
                }

            } else {
                fmt.Printf("* No records found *\n")
            }
        default:
            state.value = user_input
            state.phase += 1
    }
    return state
}

func perform_search(state ui_search_state) ui_search_state {
    if state.active_set == nil || state.active_set.json_data == nil {
        // Prompt for a new set to be selected
        state.phase = 0
        return state
    }

    state.phase -= 1
    results := search_json(state.active_set.json_data, state.key, state.value)
    max := len(results)
    for index, record := range results {
        fmt.Printf("\n**** Result[%v/%v]\n", index+1, max)
        print_record(record)

        if index > 0 && (index + 1) % 10 == 0 {
            fmt.Printf("\n   %v|%v to stop\n", command_done, command_short_done)
            fmt.Printf("   %v|%v to exit\n", command_quit, command_short_quit)
            fmt.Printf("   or enter to continue\n")
            fmt.Printf("\n%v results remaining # ", max-index)

            user_input := scan_or_exit(state.scanner)
            switch user_input {
	        case command_done, command_short_done:
	            return state
	        case command_quit, command_short_quit:
                   state.phase = -1
                   return state
            }
        }
    }
    if max > 0 {
        fmt.Printf("\nEnd of results\n")
    }
    return state
}

func print_record(result map[string]interface{}) {
    // Sorting ensures we get a stable output (great for regression
    // tests) and aids humans to quickly skip over parts of the 
    // record they don't care about.
    //
    // Not exactly efficient to recompute every time,
    // but reliable if any additional fields ever sneak in
    keys := make([]string, 0)
    for key, _:= range result {
        keys = append(keys, key)
    }
    sort.Strings(keys)

    for _, key := range keys {
        switch typed_value := result[key].(type) {
        case []interface{}:
            for index, entry := range typed_value {
	        // Alignment will suffer if there are more than 9 array entries
                fmt.Printf("%20v[%v]:  %v\n", key, index, entry)
            }
        default:
            fmt.Printf("%23v:  %v\n", key, typed_value)
        }
    }
}
