SIZES=10 100 1000 10000 100000 1000000
TESTS=missing-value unhandled-type single-record invalid
build:
	go build -a -o json_search main.go search.go ui.go

run: build
	json_search

# Higher iterations take a _really_ long time (about a day)
perfdata:
	for n in $(SIZES); do echo "Creating organizations-$$n.json"; build-data.sh data/organizations-$$n.json $$n; done

test:
	for t in $(TESTS); do json_search -filename data/$$t.json -key tags -value Veguita; done | sed -e 's/found in.*//' -e 's/Parsed .* in .*//' > data/parser.out
	diff -u data/parser.{expected,out}
	cat data/session | json_search | sed -e 's/found in.*//' -e 's/Parsed .* in .*//' > data/session.out
	diff -u data/session.{expected,out}
	json_search -filename users.json -key signature -value "Unicode ist spaß" | sed -e 's/found in.*//' -e 's/Parsed .* in .*//'  > data/unicode.out
	diff -u data/unicode.{expected,out}

perf: 
	for n in $(SIZES); do json_search -filename data/organizations-$$n.json -key _id -value 678; done
