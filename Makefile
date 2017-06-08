

build:
	go build -a -o json_search main.go search.go ui.go

run:
	go run main.go search.go ui.go

data:
	for n in 1000 1000000 1000000000; do echo "Creating organizations-$$n.json"; build-data.sh organizations-$$n.json $$n; done

test:
