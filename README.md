# json-search

## Asumptions & Limitations

- A moderately modern linux (Centos 7.2 tested).
- Standard lanuage features are sufficiently well tested to not
  benefit from additional test here.
- Datasets are located in the same directory.
- Datasets consist of an array, of maps, of strings to $something
- Records within a dataset contain a sufficently uniform set of keys
  to use the first few records when providing sample values and key
  list.
- No explicit consideration of non-ascii characters was made for this
  version, though it is likely that if you can type it you can
  search for it.
- No fields or values will be the string '.help', '.quit', or '.back'
- No-one would create a directory with a '.json' suffix in the data
  directory.
- All fields are compared as strings and require full text
  equality in order to be considered a match.

## Requirements

Core

- golang-1.6.3-2.el7

Supplimental (build, dataset generation, etc)

- bash-4.2.46-20.el7_2
- make-3.82-21.el7
- words-3.0-22.el7

## Performance

Performance as measured on

> 8-Core, Intel(R) Xeon(R) CPU E3-1245 v5 @ 3.50GHz
> With ~28Gb free RAM

| Records        | Parse         | Filter       |
| -------------: |:-------------:|:------------:|
| 10             | 198.945µs     | 13.367µs     |
| 100            | 1.497425ms    | 64.775µs     |
| 1,000          | 13.289993ms   | 405.67µs     |
| 10,000         | 141.016349ms  | 3.99669ms    |
| 100,000        | 1.385676925s  | 41.508465ms  |
| 1,000,000      | 12.434854895s | 458.426874ms |

O(N) as expected and sufficiently fast since JSON files are not widely
considered replacements for a proper database.

## Build

```make``` or ```go build -a -o json_search main.go search.go ui.go```

## Execution

### Interactive mode

```# json_search [-data path_to_directory_with_json_files]```

Sample session:

```
JSON Search Tool

Please select a dataset to search, or '.quit' to exit:
   1) Organizations 2) Tickets 3) Users 
# 3

Parsed users.json in 1.607093ms

   '.help|.h' to see available fields
   '.back|.b' to go back
   '.quit|.q' to exit

Enter a term to search Users # .h

User records contain the following fields:
* last_login_at
* signature
* organization_id
* external_id
* created_at
* verified
* locale
* timezone
* tags
* _id
* email
* url
* name
* phone
* alias
* active
* shared
* suspended
* role

   '.help|.h' to see available fields
   '.back|.b' to go back
   '.quit|.q' to exit

Enter a term to search Users # _id

   '.help|.h' to see example values
   '.back|.b' to go back
   '.quit|.q' to exit

Enter a value to search Users[_id] # .b

   '.help|.h' to see available fields
   '.back|.b' to go back
   '.quit|.q' to exit

Enter a term to search Users # tags

   '.help|.h' to see example values
   '.back|.b' to go back
   '.quit|.q' to exit

Enter a value to search Users[tags] # .h

User records contain 'tags' values like:
* Springville
* Sutton
* Hartsville/Hartley
* Diaperville
* Foxworth
* Woodlands
* Herlong
* Henrietta
* Mulino
* Kenwood
* Wescosville
* Loyalhanna
* Gallina
* Glenshaw
* Rowe
* Babb

   '.help|.h' to see example values
   '.back|.b' to go back
   '.quit|.q' to exit

Enter a value to search Users[tags] # Mulino
Searching 75 records for entries with 'tags' equal to 'Mulino'
1 record(s) found in 84.34µs

**** Result[1/1]
                    _id:  3
                 active:  false
                  alias:  Miss Buck
             created_at:  2016-07-28T05:29:25 -10:00
                  email:  buckwagner@flotonic.com
            external_id:  85c599c1-ebab-474d-a4e6-32f1c06e8730
          last_login_at:  2013-02-07T05:53:38 -11:00
                 locale:  en-AU
                   name:  Ingrid Wagner
        organization_id:  104
                  phone:  9365-482-943
                   role:  end-user
                 shared:  false
              signature:  Don't Worry Be Happy!
              suspended:  false
                tags[0]:  Mulino
                tags[1]:  Kenwood
                tags[2]:  Wescosville
                tags[3]:  Loyalhanna
               timezone:  Trinidad and Tobago
                    url:  http://initech.zendesk.com/api/v2/users/3.json
               verified:  false

End of results

   '.help|.h' to see example values
   '.back|.b' to go back
   '.quit|.q' to exit

Enter a value to search Users[tags] # .quit
```

### Scripted session

Commands, one per line, can also be placed in a file and piped to the
json_search tool.

```
# cat sample
3
.help
_id
.back
tags
.help
Mulino
# cat sample | json_search
```

### Batch mode

To avoid the need for creating session files, the tool search
parameters can also be invoked from the command-line to facilitate
testing.

```
# json_search --help
Usage of json_search:
  -data string
    	Path to a directory containing valid JSON files (default "./")
  -filename string
    	File containing valid JSON
  -key string
    	JSON key to search for
  -value string
    	Value which the specified key must contain

# json_search -filename users.json -key tags -value Mulino
Parsed users.json in 1.72752ms
Searching 75 records for entries with 'tags' equal to 'Mulino'
1 record(s) found in 25.635µs

**** Result[1/1]
                    _id:  3
                 active:  false
                  alias:  Miss Buck
             created_at:  2016-07-28T05:29:25 -10:00
                  email:  buckwagner@flotonic.com
            external_id:  85c599c1-ebab-474d-a4e6-32f1c06e8730
          last_login_at:  2013-02-07T05:53:38 -11:00
                 locale:  en-AU
                   name:  Ingrid Wagner
        organization_id:  104
                  phone:  9365-482-943
                   role:  end-user
                 shared:  false
              signature:  Don't Worry Be Happy!
              suspended:  false
                tags[0]:  Mulino
                tags[1]:  Kenwood
                tags[2]:  Wescosville
                tags[3]:  Loyalhanna
               timezone:  Trinidad and Tobago
                    url:  http://initech.zendesk.com/api/v2/users/3.json
               verified:  false
```