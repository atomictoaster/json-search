package main

import (
    "fmt"
    "io/ioutil"
    "encoding/json"
    "strings"
    "strconv"
    "reflect"
    "time"
)

func parse_file(filename string) (bool, []map[string]interface{}) {

    start := time.Now()
    
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Printf("Could not read %v: %v\n", filename, err)
        return false, nil
    }

    // We expect an array, of maps, of strings to $something
    var json_data []map[string]interface{}

    if err := json.Unmarshal(bytes, &json_data); err != nil {
        fmt.Printf("Could not parse %v (either as json, or into the expected format): %v\n", filename, err)
        return false, nil
    }
    fmt.Printf("Parsed %v in %v\n", filename, time.Now().Sub(start))
    return true, json_data
}

func search_file(filename string, search_key string, search_value string) []map[string]interface{} {
    valid, json_data := parse_file(filename)
    if valid {
        return search_json(json_data, search_key, search_value)
    }
    return make([]map[string]interface{}, 0)
}

// Match a JSON map record based on a full string match of the value for a named key
// - An array value is considered to match if any element of the array matches
//   Currently all arrarys are assumed to have strings in them (almost certainly untrue IRL)
// - Other non-string values are coerced into string format before comparision
//   Converting the other way requires more complex validation and reporting or
//   malformed user input.
// - Since the JSON parser seems to represent Ints as Floats too, we drop all 
//   decimal places when converting from that type
func json_map_filter_full_string(search_key string, search_value string, element map[string]interface{}) bool {
    value := element[search_key]
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
}

func search_json(json_data []map[string]interface{}, search_key string, search_value string) []map[string]interface{} {

    start := time.Now()
    
    fmt.Printf("Searching %v records for entries with '%v' equal to '%v'\n", len(json_data), search_key, search_value)
    results := Filter(search_key, search_value, json_data, json_map_filter_full_string)

    fmt.Printf("%v record(s) found in %v\n", len(results), time.Now().Sub(start))
    return results
}

func Filter(search_key string, search_value string, input []map[string]interface{},
            filter_function func(string, string, map[string]interface{}) bool) []map[string]interface{} {
    results := make([]map[string]interface{}, 0)
    for _, element := range input {
        if filter_function(search_key, search_value, element) {
            results = append(results, element)
        }
    }
    return results
}

