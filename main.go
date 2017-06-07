package main
import (
    "fmt"
    "io/ioutil"
    "os"
    "encoding/json"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    println("Opening file")
    bytes, err := ioutil.ReadFile("./users.json")
    check(err)
    fmt.Print(string(bytes))

    println("Parsing JSON")
    var jdat interface{}
    if err := json.Unmarshal(bytes, &jdat); err != nil {
        panic(err)
    }

    //println("Native print")
    //fmt.Println(jdat)

    println("JSON print")
    enc := json.NewEncoder(os.Stdout)
    enc.Encode(jdat)
}
