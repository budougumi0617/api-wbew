package main

import (
    "fmt"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/mmcdole/gofeed"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/sheets/v4"
)

type Info struct{
    Name string `json:"name"`
    Image string `json:"image"`
    Title string `json:"title"`
    Url string `json:"url"`
}

type Credential struct{
    Type string `json:"type"`
    Project_id string `json:"project_id"`
    Private_key_id string `json:"private_key_id"`
    Private_key string `json:"private_key"`
    Client_email string `json:"client_email"`
    Client_id string `json:"client_id"`
    Auth_uri string `json:"auth_uri"`
    Token_uri string `json:"token_uri"`
    Auth_provider_x509_cert_url string `json:"auth_provider_x509_cert_url"`
    Client_x509_cert_url string `json:"client_x509_cert_url"`
}

type Infos []Info

func httpClient(data []byte) (*http.Client, error) {

    conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/spreadsheets")
    if err != nil {
        return nil, err
    }

    return conf.Client(oauth2.NoContext), nil
}

func main() {

    spreadsheetId := "1b7vrA9DYNUYLvLMZVA2buZ8lJWWjhDLLS3sz6CKHWDQ"

    var credential = Credential{
        Type : os.Getenv("TYPE"),
        Token_uri : os.Getenv("TOKEN_URI"),
        Project_id : os.Getenv("PROJECT_ID"),
        Private_key_id : os.Getenv("PRIVATE_KEY_ID"),
        Private_key : strings.Replace(os.Getenv("PRIVATE_KEY"), "\\n", "\n", -1),
        Client_x509_cert_url : os.Getenv("CLIENT_X509_CERT_URL"),
        Client_id : os.Getenv("CLIENT_ID"),
        Client_email : os.Getenv("CLIENT_EMAIL"),
        Auth_uri : os.Getenv("AUTH_URI"),
        Auth_provider_x509_cert_url : os.Getenv("AUTH_PROVIDER_X509_CERT_URL"),
    }

    jsonBytes, err := json.Marshal(credential)
    if err != nil {
        fmt.Println("JSON Marshal error:", err)
        return
    }

    client, err := httpClient(jsonBytes)
    if err != nil {
        fmt.Println("JSON Marshal error:", err)
        return
    }

    sheetService, err := sheets.New(client)
    if err != nil {
        log.Fatalf("Unable to retrieve Sheets Client %v", err)
    }

    _, err = sheetService.Spreadsheets.Get(spreadsheetId).Do()
    if err != nil {
        log.Fatalf("Unable to get Spreadsheets. %v", err)
    }

    readRange := "data!A1:D"
    resp, err := sheetService.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
    if err != nil {
        log.Fatalf("Unable to retrieve data from sheet: %v", err)
    }

    fp := gofeed.NewParser()
    var info = Info{}
    var infos Infos


    if len(resp.Values) == 0 {
        fmt.Println("No data found.")
    } else {
            for _, row := range resp.Values {

                feed, _ := fp.ParseURL(row[1].(string))

                items := feed.Items
        
                info.Name = row[0].(string)
                info.Image =  row[2].(string)
                info.Title = items[0].Title
                info.Url = items[0].Link
                infos = append(infos, info)
            }
            outputJson, err := json.Marshal(infos)
            if err != nil {
                panic(err)
            }
            fmt.Println(string(outputJson))
        }
}