package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	h8HelperRand "github.com/novalagung/gubrak/v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type data_cetak struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func main() {
	//buat data randomize data water dan wind
	ticker := time.NewTicker(15 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				//ambil angka random
				water := h8HelperRand.RandomInt(1, 100)
				wind := h8HelperRand.RandomInt(1, 100) / 15
				var water_status string
				var wind_status string
				//kondisi water
				if water <= 5 {
					water_status = "Aman"
				} else if water > 5 && water <= 8 {
					water_status = "Siaga"
				} else if water > 8 {
					water_status = "Bahaya"
				}

				//kondisi wind
				if wind <= 6 {
					wind_status = "Aman"
				} else if water > 6 && water <= 15 {
					wind_status = "Siaga"
				} else if water > 15 {
					wind_status = "Bahaya"
				}

				//masukan data ke map
				data := map[string]interface{}{
					"water": water,
					"wind":  wind,
				}

				//ubah menjadi JSON
				requestJson, err := json.Marshal(data)
				client := &http.Client{}
				if err != nil {
					log.Fatalln(err)
				}

				//buat request dengan fungsi http.NewRequest
				req, err := http.NewRequest("POST",
					"https://jsonplaceholder.typicode.com/posts",
					bytes.NewBuffer(requestJson))
				req.Header.Set("Content-type", "application/json")
				if err != nil {
					log.Fatalln(err)
				}

				res, err := client.Do(req)
				if err != nil {
					log.Fatal(err)
				}

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Fatalln(err)
				}

				//unmarshall data agar hanya mengambil nilai water dan wind
				var cetak data_cetak
				if err := json.Unmarshal([]byte(string(body)), &cetak); err != nil {
					fmt.Println(err)
					return
				}

				//data kembali diconvert ke json sehingga hanya mengeluarkan output water dan wind
				cetakjson, err := json.Marshal(cetak)
				if err != nil {
					log.Fatalln(err)
				}

				log.Println(string(cetakjson))
				//print status berdasarkan kondisi
				fmt.Printf("status water: %s \n", water_status)
				fmt.Printf("status wind: %s \n \n", wind_status)
			}

		}
	}()
	<-done
}
