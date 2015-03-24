package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

//TODO: http://blog.csdn.net/xxxxxx91116/article/details/7971134
// 添加其他获取种子的方法

func getTorrentXunlei(hash string) {
	a := hash[0:2]
	b := hash[len(hash)-2 : len(hash)]
	fmt.Println(a, b)

	resp, _ := http.Get(fmt.Sprintf("http://bt.box.n0808.com/%s/%s/%s.torrent", a, b, hash))
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		f, _ := os.Create(hash + ".torrent")
		defer f.Close()

		_, err := io.Copy(f, resp.Body)
		if err != nil {
			fmt.Printf("Write %s.torrent failed: %s\n", hash, err)
		}
	}
}

func getTorrent(hash string, wg *sync.WaitGroup) {
	defer wg.Done()
	getTorrentXunlei(hash)
}

func main() {

	hashs := os.Args[1:]

	var wg sync.WaitGroup
	wg.Add(len(hashs))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	for _, hash := range hashs {
		go getTorrent(hash, &wg)
	}

	go func() {
		s := <-c
		fmt.Printf("Got os signal: %v, exit\n", s)
		os.Exit(1)
	}()

	wg.Wait()
	fmt.Println("All done")
}
