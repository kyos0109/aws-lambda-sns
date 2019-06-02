package main

import (
    "context"
    "fmt"
    "encoding/json"
    "errors"
    "log"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"
)

const lineNotifyURL = "https://notify-api.line.me/api/notify"

type LineInfo struct {
    Token   string
    Message string
    Debug   bool
}

type codeDeployReturn struct {
    EventTriggerName    string `json:"eventTriggerName"`
    DeploymentId        string `json:"deploymentId"`
    ApplicationName     string `json:"applicationName, omitempty"`
    DeploymentGroupName string `json:"deploymentGroupName, omitempty"`
    Status              string `json:"status, omitempty"`
    InstanceStatus		string `json:instanceStatus, omitempty`
    // ErrorInformation    string `json:"errorInformation"`
}

func main() {
    lambda.Start(handler)
}


func handler(ctx context.Context, snsEvent events.SNSEvent) {

    for _, record := range snsEvent.Records {
        snsRecord := record.SNS

        info := LineInfo{
            Token:   os.Getenv("TOKEN"),
            Message: convertMessage(snsRecord.Message),
            Debug:   getBoolEnv("DEBUG"),
        }

        if err := send(info); err != nil {
            log.Fatal(err.Error())
        }

        fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)
    }
}

func convertMessage(msg string) string {

    var codeMsg codeDeployReturn
    json.Unmarshal([]byte(msg), &codeMsg)
    if m, err := json.MarshalIndent(codeMsg, "", "\t"); err == nil {
        return string(m)
    }
    return ""
}

func send(l LineInfo) error {

    if l.Token == "" || l.Message == "" || len(l.Message) == 0 {
        return errors.New("error env: PLUGIN_TOKEN or PLUGIN_MESSAGE is empty.")
    }

    data := url.Values{}
    data.Add("message", l.Message)

    req, err := http.NewRequest(
        "POST",
        lineNotifyURL,
        strings.NewReader(data.Encode()),
    )

    if err != nil {
        return errors.New("error request : " + err.Error())
    }

    req.Header.Add("Authorization", "Bearer "+l.Token)
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

    client := &http.Client{}
    resp, err := client.Do(req)

    if err != nil {
        return errors.New("error response: " + err.Error())
    }

    defer resp.Body.Close()

    if resp.StatusCode == 200 {
        log.Println("send...OK")
    }

    if l.Debug {
        log.Println(resp)
    }

    return nil
}

func getBoolEnv(key string) bool {
    if v, ok := os.LookupEnv(key); ok {
        if strings.ToLower(v) == "true" {
            return true
        }
    }
    return false
}