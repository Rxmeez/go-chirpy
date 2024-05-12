run:
	if [ -f database.json ]; then rm database.json; fi
	go run .

.PHONY: 
	run