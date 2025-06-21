package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type Event struct {
	OCRText   string `json:"ocrText"`
	RequestID string `json:"requestId"`
}

type Response struct {
	Category   string   `json:"category"`
	Confidence float64  `json:"confidence"`
	Keywords   []string `json:"keywords"`
	RequestID  string   `json:"requestId"`
}

var (
	bedrockClient *bedrockruntime.Client
	categories    = []string{"transfer", "lending", "otp", "ott", "login", "other"}
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	bedrockClient = bedrockruntime.NewFromConfig(cfg)
}

func handler(ctx context.Context, event Event) (Response, error) {
	// GenAI Prompt Engineering
	prompt := fmt.Sprintf(`Classify this banking error into exactly one of these categories: %s.
Text: %s
Respond ONLY with the category name.`, strings.Join(categories, ", "), event.OCRText)

	// Bedrock API Request (Claude Instant)
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String("anthropic.claude-instant-v1"),
		ContentType: aws.String("application/json"),
		Body: []byte(fmt.Sprintf(`{
			"prompt": "\n\nHuman:%s\n\nAssistant:",
			"max_tokens_to_sample": 30,
			"temperature": 0.5
		}`, prompt)),
	}

	output, err := bedrockClient.InvokeModel(ctx, input)
	if err != nil {
		return Response{}, fmt.Errorf("Bedrock error: %v", err)
	}

	var resp struct {
		Completion string `json:"completion"`
	}
	json.Unmarshal(output.Body, &resp)

	// Post-processing
	cleanedCat := strings.ToLower(strings.TrimSpace(resp.Completion))
	if !contains(categories, cleanedCat) {
		cleanedCat = "other"
	}

	return Response{
		Category:   cleanedCat,
		Confidence: 0.9, // Mock confidence
		Keywords:   extractKeywords(event.OCRText),
		RequestID:  event.RequestID,
	}, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func extractKeywords(text string) []string {
	// Implement your keyword extraction logic
	return []string{"error", "failed"}
}

func main() {
	// lambda.Start(handler)
	out, err := handler(context.Background(), Event{
		OCRText:   "10:09Chuyen dênTechcombankTCB - Ngân hàng TMCP Ky thuong ViêtNamSô tài khoán19037567041028!Không thê truy vân thông tinnguði nhânQuý khách có muôn tiêp tuc giao dich bångphuong thúc Chuyen khoán thuong duoc thuchiên qua hê thông thanh toán cúa NHNN (trong giòlàm viêc tù thú 2 dên thú 6) không?Thu laiChuyên khoán thuong",
		RequestID: "123",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)

}
