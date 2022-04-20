package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/morzhanov/go-termui-chess/internal/board"
	"github.com/morzhanov/go-termui-chess/internal/ui"
)

func main() {
	b := board.New()
	renderer, err := ui.NewRenderer(b)
	if err != nil {
		log.Fatal(err)
	}
	go renderer.Start()
	time.Sleep(time.Second * 1)

	for {
		fmt.Println("Enter your move...")
		var reader = bufio.NewReader(os.Stdin)
		move, _ := reader.ReadString('\n')
		renderer.Update(move)
	}
}
