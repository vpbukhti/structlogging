install:
	go install golang.org/x/tools/cmd/present@latest

demo:
	go run . 

present:
	present presentation.slide