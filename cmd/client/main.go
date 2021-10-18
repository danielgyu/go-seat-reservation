package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const localhost string = "http://localhost:8000"
const appjson string = "application/json"
const hallName string = "sarang"
const newUserCount int = 999

type TokenBody struct {
	Status string
	Token  string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func makeUrl(endpoint string) string {
	return fmt.Sprintf("%s%s", localhost, endpoint)
}

func main() {
	log.Println("started client with pid:", os.Getpid())
	deleteAllUsers()
	log.Println("succesfully deleted all users")

	for i := 0; i < newUserCount; i++ {
		createUser(strconv.Itoa(i))
	}
	log.Printf("successfully created %d users\n", newUserCount)

	var tokens []string
	for i := 0; i < newUserCount; i++ {
		token := getToken(strconv.Itoa(i))
		tokens = append(tokens, token)
	}
	log.Printf("successfully got %d tokens\n", newUserCount)

	start := time.Now()
	wg := new(sync.WaitGroup)
	for _, token := range tokens {
		wg.Add(1)
		go reserveSeat(wg, token, hallName)
	}
	wg.Wait()
	elapsed := int(time.Since(start).Seconds())
	log.Printf("successfully reserved %d seats\n", newUserCount)
	log.Printf("reserving %d seats took %d seconds\n", newUserCount, elapsed)
}

func deleteAllUsers() {
	url := makeUrl("/users")
	req, err := http.NewRequest("DELETE", url, nil)
	check(err)

	_, err = http.DefaultClient.Do(req)
	check(err)
}

func createUser(trailer string) {
	url := makeUrl("/signup")
	userpass := "user" + trailer
	body, _ := json.Marshal(map[string]string{
		"username": userpass,
		"password": userpass,
	})
	reqBody := bytes.NewBuffer(body)
	_, err := http.Post(url, appjson, reqBody)
	check(err)
}

func getToken(trailer string) (token string) {
	url := makeUrl("/login")
	userpass := "user" + trailer
	reqBody, _ := json.Marshal(map[string]string{
		"username": userpass,
		"password": userpass,
	})
	resBody := bytes.NewBuffer(reqBody)
	resp, err := http.Post(url, appjson, resBody)
	check(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	var tb TokenBody
	err = json.Unmarshal(body, &tb)
	check(err)

	token = tb.Token
	return
}

func reserveSeat(wg *sync.WaitGroup, token string, hallName string) {
	defer wg.Done()
	url := makeUrl(fmt.Sprintf("/reservation/%s", hallName))

	req, err := http.NewRequest("GET", url, nil)
	check(err)
	req.Header.Set("Authorization", token)

	_, err = http.DefaultClient.Do(req)
	check(err)
}
