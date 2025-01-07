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

	fmt.Println("Sending email...: ", args)
	recipient := args["recipient"].(string)
	subject := args["subject"].(string)
	body := args["body"].(string)
	fmt.Println("Sending email...")
	fmt.Printf("To: %s\nSubject: %s\nBody: %s\n", recipient, subject, body)
	return wsarmgo.Result{
		Success: true,
		Data:    "Sent!",
	}
}

func main() {
	dotenv.Load()

	client := wsarmgo.NewSwarm(os.Getenv("OPENAI_API_KEY"), llm.OpenAI)

	fn_getweather := NewGetWeatherFn()
	fn_sendmail := NewSendEmalFn()

	model := wsarmgo.LLM{
		Model:       "gpt-4",
		LLMProvider: "OPEN_AI",
		ApiKey:      os.Getenv("OPENAI_API_KEY"),
	}
	weatherAgent := wsarmgo.NewBaseAgent("WeatherAgent", "You are a helpful weather assistant. Always respond in a natural, conversational way. When providing weather information, format it in a friendly manner rather than just returning raw data. For example, instead of showing JSON, say something like 'The temperature in [city] is [temp] degrees.'", model)
	weatherAgent.AddFunction(fn_getweather)
	weatherAgent.AddFunction(fn_sendmail)

	wsarmgo.RunDemoLoop(client, weatherAgent)
}

type sendEmailFn struct {
	wsarmgo.BaseFunction
}

func NewSendEmalFn() *sendEmailFn {
	fn := &sendEmailFn{}
	fn.BaseFunction = *wsarmgo.NewCustomFunction(fn)
	fn.BaseFunction.SetFunction(sendEmail)
	return fn
}

var _ wsarmgo.AgentFunction = (*sendEmailFn)(nil)

func (fn *sendEmailFn) GetName() string {
	return "sendEmailFn"
}
func (fn *sendEmailFn) GetDescription() string {
	return "Send an email to a recipient."
}
func (fn *sendEmailFn) GetParameters() map[string]interface{} {
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

type getWeatherFn struct {
	wsarmgo.BaseFunction
}

func NewGetWeatherFn() *getWeatherFn {
	fn := &getWeatherFn{}
	fn.BaseFunction = *wsarmgo.NewCustomFunction(fn)
	fn.BaseFunction.SetFunction(getWeather)
	return fn
}

var _ wsarmgo.AgentFunction = (*getWeatherFn)(nil)

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
