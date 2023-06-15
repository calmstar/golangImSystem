package main

func main() {
	s := NewServer("127.0.0.1", 9999)
	s.Start()
}
