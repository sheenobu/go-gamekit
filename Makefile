

examples: bin
	go build -o bin/dragdemo ./_examples/dragdemo
	go build -o bin/dragdropdemo ./_examples/dragdropdemo
	go build -o bin/helloworld ./_examples/helloworld

bin:
	mkdir -p bin
