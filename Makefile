run:
	go mod download
	go build -o server ./back/cmd
	python -m webbrowser "file://$(shell pwd)/front/index.html"
	./server