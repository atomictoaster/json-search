package main
import (
    "fmt"
    "io/ioutil"
    "os"
    "encoding/json"
    "strings"
    "flag"
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
    search_key := flag.String("key", "name", "JSON key to search for")
    search_value := flag.String("value", "Catalina Simpson", "Value which the specified key must contain")

    flag.Parse()

    println("Reading", *filename)
    bytes, err := ioutil.ReadFile(*filename)
    check(err)
    fmt.Print(string(bytes))

    println("Parsing", *filename)

    // We expect an array, of maps, of strings to $something
    var jdat []map[string]interface{}

    // TODO: Test behaviour of parser with malformed and "valid json
    // but not arranged like this" input
    if err := json.Unmarshal(bytes, &jdat); err != nil {
        panic(err)
    }

    println("Filtering based on", *search_key, "=", *search_value)
    enc := json.NewEncoder(os.Stdout)
    enc.Encode(Filter(jdat, func(v map[string]interface{}) bool {
        // TODO: Handle non-string fields, eg. arrays and even numbers
	name := v[*search_key].(string)
        return strings.HasPrefix(name, *search_value)
    }))

}
