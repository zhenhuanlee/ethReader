package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ppool "./pool"
	"github.com/go-redis/redis"
)

const url = "https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&boolean=true&apikey=YourApiKeyToken&tag="

type block struct {
	JSONPRC string `json:""`
	ID      int    `json:""`
	RESULT  struct {
		GASUSED      string `json:"gasUsed"`
		HASH         string `json:"hash"`
		NONCE        string `json:"nonce"`
		NUMBER       string `json:"number"`
		TIMESTAMP    string `json:"timestamp"`
		TRANSACTIONS []struct {
			BLOCKHASH        string `json:"blockHash"`
			BLOCKNUMBER      string `json:"blockNumber"`
			FROM             string `json:"from"`
			GAS              string `json:"gas"`
			GASPRICE         string `json:"gasPrice"`
			HASH             string `json:"hash"`
			INPUT            string `json:"input"`
			NONCE            string `json:"nonce"`
			TO               string `json:"to"`
			TRANSACTIONINDEX string `json:"transactionIndex"`
			VALUE            string `json:"value"`
		} `json:"TRANSACTIONS"`
	} `json:"RESULT"`
}

func main() {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 1})
	pool := ppool.NewPool(10)
	from := 1000000

	for {
		res, err := http.Get(fmt.Sprintf("%s0x%x", url, from))
		from++

		if err != nil {
			fmt.Println(err)
			continue
		}

		pool.JobChan <- func() {
			foo := toJSON(res.Body)
			for _, tx := range foo.RESULT.TRANSACTIONS {
				str, _ := json.Marshal(tx)
				client.Set(tx.HASH, str, 0)
				fmt.Println(client.Get(tx.HASH))
			}
		}

	}
}

func toJSON(body io.ReadCloser) *block {
	b := new(block)
	json.NewDecoder(body).Decode(&b)
	// fmt.Printf("%+v\n", b)
	body.Close()
	return b
}
