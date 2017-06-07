package main
import (
    "fmt"
    "io/ioutil"
    "os"
    "encoding/json"
    "strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func Any(vs []map[string]interface{}, f func(map[string]interface{}) bool) bool {
    for _, v := range vs {
        if f(v) {
            return true
        }
    }
    return false
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
    filename := "./users.json"
    search_key := "name"
    search_value := "Catalina Simpson"
    
    println("Reading", filename)
    bytes, err := ioutil.ReadFile(filename)
    check(err)
    fmt.Print(string(bytes))

    println("Parsing", filename)

    // We expect an array, of maps, of strings to $something
    var jdat []map[string]interface{}

    // TODO: Test behaviour of parser with malformed and "valid json
    // but not arranged like this" input
    if err := json.Unmarshal(bytes, &jdat); err != nil {
        panic(err)
    }

    println("Searching for", search_key, "=", search_value)
    fmt.Println(Any(jdat, func(v map[string]interface{}) bool {
	name := v[search_key].(string)
        return strings.HasPrefix(name, search_value)
    }))

    println("Filtering based on", search_key, "=", search_value)
    enc := json.NewEncoder(os.Stdout)
    enc.Encode(Filter(jdat, func(v map[string]interface{}) bool {
	name := v[search_key].(string)
        return strings.HasPrefix(name, search_value)
    }))

}
