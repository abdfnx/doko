package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	httpClient "github.com/abdfnx/resto/client"

	"github.com/tidwall/gjson"
)

func GetLatest() string {
	url := "https://api.github.com/repos/abdfnx/doko/releases/latest"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Errorf("Error creating request: %s", err.Error())
	}

	client := httpClient.HttpClient()
	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error sending request: %s", err.Error())
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("Error reading response: %s", err.Error())
	}

	body := string(b)

	tag_name := gjson.Get(body, "tag_name")

	latestVersion := tag_name.String()

	return latestVersion
}
