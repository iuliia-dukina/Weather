package main

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	//страничку с погодой из гисметео взяла через курл через постман, потому что просто по ссылке ее не отдавали. ниже код сформированный постманом:
	url := "https://www.gismeteo.ru/weather-sankt-peterburg-4079/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	//
	req.Header.Add("authority", "www.gismeteo.ru")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("referer", "https://www.gismeteo.ru/weather-sankt-peterburg-4079/now/")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"104\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"104\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// ниже запускаем парсер странички:
	stringBody := string(body)                    // преобразуем массив байт в строку
	bodyIoReader := strings.NewReader(stringBody) // создаем io reader из строки
	doc, err := htmlquery.Parse(bodyIoReader)     // строим dom дерево.

	a := htmlquery.FindOne(doc, "//span[@class=\"unit unit_temperature_c\"][1]") // пакет htmlquery позволяет испорльзовать Xpath селекторы, поэтому используем его. Ищем первый встретившийся нам селектор [1] (span.unit.unit_temperature_c) и получаем один узел - FindOne, а не коллекцию
	innerText := htmlquery.InnerText(a)                                          // иннертекст - текст без тегов

	// Для того чтобы выложить в гит, нужно заменить токен бота на переменную окружения. В Windows>Свойства системы>Переменные среды задаем переменную TG_BOT_KEY со значением токена из телеграмм бота. Кроме этого устанавливаем переменное окружение через ide в настройках go build Weather>Edit Configurations>Environment. Проверить или разово установить можно в терминале через SET

	tgBotKey := os.Getenv("TG_BOT_KEY") // создаем переменную, чтобы ею заменить токен в url

	var x [4]string
	x[0] = "269892926"
	x[1] = "387585282"
	x[2] = "245240009"
	x[3] = "385599092"

	telegramMassivIDJson, _ := json.Marshal(x)
	fmt.Println(string(telegramMassivIDJson))

	//sv := []string{"269892926, 387585282, 245240009, 385599092"}
	//boolVar, _ := json.Marshal(sv)
	//fmt.Println(string(boolVar))

	var telegramMassivID []string
	_ = json.Unmarshal(telegramMassivIDJson, &telegramMassivID)

	message := []byte("Hello, Gophers!")
	err := os.WriteFile("testdata/hello", message, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// найти запрос, который получает входящий запрос для бота и вставить его в постман
	for i := 0; i < len(telegramMassivID); i++ {
		// x[i]

		url = "https://api.telegram.org/bot" + tgBotKey + "/sendMessage?chat_id=" + x[i] + "&text=" + strings.Trim(strings.Trim(innerText, " "), "\n")
		method = "GET"

		client = &http.Client{}
		req, err = http.NewRequest(method, url, nil)

		if err != nil {
			log.Fatal(err)
		}
		res, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))
	}
}
