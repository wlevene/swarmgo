package main

import (
	"fmt"
	"os"

	dotenv "github.com/joho/godotenv"
	wsarmgo "github.com/wlevene/wsarmgo"
	"github.com/wlevene/wsarmgo/llm"
)

func getWeather(args map[string]interface{}, contextVariables map[string]interface{}) wsarmgo.Result {
	location := args["location"].(string)
	time := "now"
	if t, ok := args["time"].(string); ok {
		time = t
	}
	return wsarmgo.Result{
		Success: true,
		Data:    fmt.Sprintf(`The temperature in %s is 65 degrees at %s.`, location, time),
	}
}

func sendEmail(args map[string]interface{}, contextVariables map[string]interface{}) wsarmgo.Result {
	recipient := args["recipient"].(string)
	subject := args["subject"].(string)
	body := args["body"].(string)
	fmt.Printf("Sending email...\nTo: %s\nSubject: %s\nBody: %s\n", recipient, subject, body)
	return wsarmgo.Result{
		Success: true,
		Data:    "Sent!",
	}
}

func main() {
	if err := dotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	client := wsarmgo.NewSwarm(os.Getenv("OPENAI_API_KEY"), llm.OpenAI)

	fnGetWeather, err := NewGetWeatherFn()
	if err != nil {
		fmt.Println("Error creating getWeather function:", err)
		os.Exit(1)
	}

	fnSendEmail, err := NewSendEmailFn()
	if err != nil {
		fmt.Println("Error creating sendEmail function:", err)
		os.Exit(1)
	}

	model := wsarmgo.LLM{
		Model:       "gpt-4",
		LLMProvider: "OPEN_AI",
		ApiKey:      os.Getenv("OPENAI_API_KEY"),
	}

	weatherAgent := wsarmgo.NewBaseAgent(
		"WeatherAgent",
		"You are a helpful weather assistant. Always respond in a natural, conversational way. When providing weather information, format it in a friendly manner rather than just returning raw data. For example, instead of showing JSON, say something like 'The temperature in [city] is [temp] degrees.'",
		model,
	)

	weatherAgent.AddFunction(fnGetWeather)
	weatherAgent.AddFunction(fnSendEmail)

	wsarmgo.RunDemoLoop(client, weatherAgent)
}

type sendEmailFn struct {
	wsarmgo.BaseFunction
}

func NewSendEmailFn() (*sendEmailFn, error) {
	fn := &sendEmailFn{}
	baseFn, err := wsarmgo.NewCustomFunction(fn)
	if err != nil {
		return nil, err
	}
	fn.BaseFunction = *baseFn
	fn.BaseFunction.SetFunction(sendEmail)
	return fn, nil
}

var _ wsarmgo.AgentFunction = (*sendEmailFn)(nil)

func (fn *sendEmailFn) GetID() string {
	return "sendEmailFn"
}

func (fn *sendEmailFn) GetName() string {
	return "sendEmail"
}

func (fn *sendEmailFn) GetDescription() string {
	return "Send an email to a recipient."
}

func (fn *sendEmailFn) GetParameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"recipient": map[string]interface{}{
				"type":        "string",
				"description": "The recipient's email address",
			},
			"subject": map[string]interface{}{
				"type":        "string",
				"description": "The subject of the email",
			},
			"body": map[string]interface{}{
				"type":        "string",
				"description": "The body of the email",
			},
		},
		"required": []interface{}{"recipient", "subject", "body"},
	}
}

type getWeatherFn struct {
	wsarmgo.BaseFunction
}

func NewGetWeatherFn() (*getWeatherFn, error) {
	fn := &getWeatherFn{}
	baseFn, err := wsarmgo.NewCustomFunction(fn)
	if err != nil {
		return nil, err
	}
	fn.BaseFunction = *baseFn
	fn.BaseFunction.SetFunction(getWeather)
	return fn, nil
}

var _ wsarmgo.AgentFunction = (*getWeatherFn)(nil)

func (fn *getWeatherFn) GetID() string {
	return "getWeatherFn"
}

func (fn *getWeatherFn) GetName() string {
	return "getWeather"
}

func (fn *getWeatherFn) GetDescription() string {
	return "Get the current weather in a given location. Location MUST be a city."
}

func (fn *getWeatherFn) GetParameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "The city to get the weather for",
			},
			"time": map[string]interface{}{
				"type":        "string",
				"description": "The time to get the weather for",
			},
		},
		"required": []interface{}{"location"},
	}
}
