package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "encoding/json"
    "strings"
    "strconv"
    "reflect"
    "time"
)

func parse_file(filename string) (bool, []map[string]interface{}) {

    start := time.Now()
    
    bytes, err := ioutil.ReadFile(filename)
    check(err)

    // We expect an array, of maps, of strings to $something
    var json_data []map[string]interface{}

    if err := json.Unmarshal(bytes, &json_data); err != nil {
        fmt.Printf("Could not parse %v: %v\n", filename, err)
	return false, nil
    }
    fmt.Printf("Parsed %v in %v\n", filename, time.Now().Sub(start))
    return true, json_data
}

func search_file(filename string, search_key string, search_value string) int {
    valid, json_data := parse_file(filename)
    if valid {
        return search_json(json_data, search_key, search_value)
    }
    return 0
}

func search_json(json_data []map[string]interface{}, search_key string, search_value string) int {

    start := time.Now()
    
    fmt.Printf("Filtering %v records based on '%v'='%v'\n", len(json_data), search_key, search_value)
    results := Filter(json_data, func(v map[string]interface{}) bool {
        value := v[search_key]
        if value == nil && len(search_value) == 0 {
            // Searching for empty values
            return true
        } else if value == nil {
            return false
        }
        string_value := ""

        switch typed_value := value.(type) {
        case bool:
            string_value = strconv.FormatBool(value.(bool));
        case int64:
            string_value = strconv.FormatInt(value.(int64), 10);
        case float64:
            string_value = strconv.FormatFloat(value.(float64), 'f', 0, 64);
        case string:
            string_value = value.(string)
	case []interface{}:
            for _, entry := range typed_value {
	        // TODO: Handle other types too
                if strings.Compare(entry.(string), search_value) == 0 {
                    return true
	        }
            }
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
        // TODO: This output is awful
        enc := json.NewEncoder(os.Stdout)
        enc.Encode(results)
    }
    fmt.Printf("%v record(s) found in %v\n", num_results, time.Now().Sub(start))
    return num_results
}

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
