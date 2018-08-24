package main

import (
	"fmt"
	ipfs "github.com/ipfs/go-ipfs-api"
	"time"
	"os"
	"bufio"
)

func subChannel(subscribtion *ipfs.PubSubSubscription, c chan string) {
	for {
		records, _ := subscribtion.Next()
		content := records.Data()
		peerID := records.From()
		if content != nil {
			msg := time.Now().String() + "\n from <" + peerID.Pretty()[:5] + "***>: " + string(content) + "\n"
			c <- msg
		}
	}
}

func printch(c chan string) {
	for i := range c {
		fmt.Println(i)
	}

}

func joinChannel(ipfShell *ipfs.Shell, channel string) {
	ch := make(chan string)
	subscribtion, _ := ipfShell.PubSubSubscribe(channel)
	go subChannel(subscribtion, ch)
	printch(ch)
}

func uploadFile(ipfShell *ipfs.Shell, filePath string) string {
	f, _ := os.Stat(filePath)
	var result string
	if f.IsDir() {
		result, _ = ipfShell.AddDir(filePath)

	} else {
		f, _ := os.Open(filePath)
		result, _ = ipfShell.Add(f)
	}
	return result
}
func main() {

	//var channelName string;

	//flag.StringVar(&channelName, "c", "hello", "channel that use")
	//flag.Parse()

	channelName := "hello"
	host := "http://192.168.1.140"
	apiPort := "5001"
	gatePort := "8080"
	ipfShell := ipfs.NewShell(host + ":" + apiPort)
	go joinChannel(ipfShell, channelName)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		if len(s) > 4 && s[:4] == "@pf@" {
			fileHash := uploadFile(ipfShell, s[4:])
			s = host + ":" + gatePort + "/ipfs/" + fileHash
		}
		ipfShell.PubSubPublish(channelName, s)
	}
	if err := scanner.Err(); err != nil {
		os.Exit(1)
	}
}
