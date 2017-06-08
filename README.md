# json-search

## Asumptions & Limitations

- A moderately modern linux (Centos 7.2 tested).
- Datasets are located in the same directory.
- Records within a dataset contain a sufficently uniform set of keys
  to use the first few records when providing sample values and key
  list.
- No explicit consideration of non-ascii characters was made for this
  version.
- No fields or values will be named 'quit', '?', or '..'
- No-one would create a directory with a '.json' suffix in the data
  directory.

## Requirements

Core

- golang-1.6.3-2.el7

Supplimental (build, dataset generation, etc)

- bash-4.2.46-20.el7_2
- make-3.82-21.el7
- words-3.0-22.el7

## Performance

- Search 1,000 records:         
- Search 1,000,000 records:     
- Search 1,000,000,000 records: 

## Build

```make``` or ```go build -a -o json_search main.go search.go ui.go```

## Execution

### Interactive mode

```# json_search```

Sample session:

```
JSON Search Tool

Please select a dataset to search, or 'quit' to exit:
   1) Organizations 2) Tickets 3) Users 
# 3
Parsed users.json in 2.737131ms

Enter a term to search for:
   '?' to see available fields,
   '..' to go back
   'quit' to exit
Users # ?

User records contain the following fields
* last_login_at
* suspended
* external_id
* name
* alias
* _id
* locale
* organization_id
* email
* tags
* url
* active
* timezone
* phone
* signature
* role
* created_at
* verified
* shared

Enter a term to search for:
   '?' to see available fields,
   '..' to go back
   'quit' to exit
Users # _id

Enter a value to search for:
   '?' to see example values,
   '..' to go back
   'quit' to exit
Users[_id] # ..

Enter a term to search for:
   '?' to see available fields,
   '..' to go back
   'quit' to exit
Users # tags

Enter a value to search for:
   '?' to see example values,
   '..' to go back
   'quit' to exit
Users[tags] # ?

User records contain values like:
* [Springville Sutton Hartsville/Hartley Diaperville]
* [Foxworth Woodlands Herlong Henrietta]
* [Mulino Kenwood Wescosville Loyalhanna]
* [Gallina Glenshaw Rowe Babb]

Enter a value to search for:
   '?' to see example values,
   '..' to go back
   'quit' to exit
Users[tags] # Mulino
Filtering 75 records based on 'tags'='Mulino'
[{"_id":3,"active":false,"alias":"Miss Buck","created_at":"2016-07-28T05:29:25 -10:00","email":"buckwagner@flotonic.com","external_id":"85c599c1-ebab-474d-a4e6-32f1c06e8730","last_login_at":"2013-02-07T05:53:38 -11:00","locale":"en-AU","name":"Ingrid Wagner","organization_id":104,"phone":"9365-482-943","role":"end-user","shared":false,"signature":"Don't Worry Be Happy!","suspended":false,"tags":["Mulino","Kenwood","Wescosville","Loyalhanna"],"timezone":"Trinidad and Tobago","url":"http://initech.zendesk.com/api/v2/users/3.json","verified":false}]
1 record(s) found in 220.473µs

Please select a dataset to search, or 'quit' to exit:
   1) Organizations 2) Tickets 3) Users 
# quit
```

### Batch mode

```# json_search -filename users.json -key tags -value Mulino```

```
Parsed users.json in 4.130999ms
Filtering 75 records based on 'tags'='Mulino'
[{"_id":3,"active":false,"alias":"Miss Buck","created_at":"2016-07-28T05:29:25 -10:00","email":"buckwagner@flotonic.com","external_id":"85c599c1-ebab-474d-a4e6-32f1c06e8730","last_login_at":"2013-02-07T05:53:38 -11:00","locale":"en-AU","name":"Ingrid Wagner","organization_id":104,"phone":"9365-482-943","role":"end-user","shared":false,"signature":"Don't Worry Be Happy!","suspended":false,"tags":["Mulino","Kenwood","Wescosville","Loyalhanna"],"timezone":"Trinidad and Tobago","url":"http://initech.zendesk.com/api/v2/users/3.json","verified":false}]
1 record(s) found in 201.074µs
```

