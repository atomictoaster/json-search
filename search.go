package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "encoding/json"
    "strings"
    "strconv"
    "reflect"
)

func parse_file(filename string) []map[string]interface{} {
    println("Reading", filename)
    bytes, err := ioutil.ReadFile(filename)
    check(err)

    println("Parsing", filename)

    // We expect an array, of maps, of strings to $something
    var json_data []map[string]interface{}

    // TODO: Test behaviour of parser with malformed and "valid json,
    // but not arranged like this" input
    // TODO: Also, slightly more graceful error handling
    if err := json.Unmarshal(bytes, &json_data); err != nil {
        panic(err)
    }
    return json_data
}

func search_file(filename string, search_key string, search_value string) int {
    json_data := parse_file(filename)
    return search_json(json_data, search_key, search_value)
}

func search_json(json_data []map[string]interface{}, search_key string, search_value string) int {

    println("Filtering", len(json_data), "records based on", search_key, "=", search_value)
    results := Filter(json_data, func(v map[string]interface{}) bool {
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
