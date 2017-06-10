SIZES=10 100 1000 10000 100000 1000000

build:
	go build -a -o json_search main.go search.go ui.go

run: build
	json_search

# Higher iterations take a _really_ long time (about a day)
perfdata:
	for n in $(SIZES); do echo "Creating organizations-$$n.json"; build-data.sh data/organizations-$$n.json $$n; done

test:

perf: 
	for n in $(SIZES); do json_search -filename data/organizations-$$n.json -key _id -value 678; done
