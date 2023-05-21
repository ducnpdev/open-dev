package main

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
)

func CreateSession() *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           *aws.String("my"),
		Config: aws.Config{
			Region: aws.String("ap-southeast-1"), // London
		},
	}))
	return sess
}

var textractSession *textract.Textract

func init() {
	textractSession = textract.New(CreateSession())
}
func main() {
	file, err := ioutil.ReadFile("front.jpg")
	if err != nil {
		panic(err)
	}
	textractSession.DetectDocumentTextRequest(&textract.DetectDocumentTextInput{
		// Document: ,

	})
	resp, err := textractSession.DetectDocumentText(&textract.DetectDocumentTextInput{
		Document: &textract.Document{
			Bytes: file,
		},
	})
	if err != nil {
		panic(err)
	}
	for i := 1; i < len(resp.Blocks); i++ {
		a := *resp.Blocks[i]
		fmt.Println(a)
		if *resp.Blocks[i].BlockType == "LINE" {
			fmt.Println(*resp.Blocks[i].Text)
		} else if *resp.Blocks[i].BlockType == "WORD" {
			fmt.Println(*resp.Blocks[i].Text)
		} else {
			fmt.Println(*resp.Blocks[i].BlockType)
			fmt.Println(*resp.Blocks[i].Text)
		}
	}
}
