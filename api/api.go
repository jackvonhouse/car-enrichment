package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-faker/faker/v4"
	"log"
	"math/rand/v2"
	"net/http"
)

type Owner struct {
	Name       string `json:"name" faker:"first_name"`
	Surname    string `json:"surname" faker:"last_name"`
	Patronymic string `json:"patronymic" faker:"first_name"`
}

type Car struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark" faker:"word"`
	Model  string `json:"model" faker:"word"`
	Year   int    `json:"year" faker:"oneof: 15, 27, 61"`
	Owner  Owner  `json:"owner"`
}

func main() {
	http.Handle("/info", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		regNum := queries.Get("regNum")

		i := rand.Int()
		if i%2 == 0 {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		log.Println(regNum)

		car := Car{}

		if err := faker.FakeData(&car); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		car.RegNum = regNum

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(car); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}))

	fmt.Println(http.ListenAndServe(":9999", nil))
}
