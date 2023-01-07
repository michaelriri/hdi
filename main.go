package main

import (
	"fmt"
	"hdi/openai"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rs/zerolog"
)

func main() {

	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	secret, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		log.Fatal().Msg("OPENAI_API_KEY not set")
	}
	log.Info().Str("secret", secret).Msg("using OPENAI_API_KEY")

	question := strings.Join(os.Args[1:], " ")
	log.Info().Str("question", "[how do i] "+question).Msg("question for gpt-3 (openai)")

	osTarget := runtime.GOOS
	prompt := fmt.Sprintf(TEMPLATE_PROMPT, osTarget, osTarget, question)

	request := openai.CompletionRequest{
		Model:       "text-davinci-003",
		Prompt:      prompt,
		MaxTokens:   100,
		Temperature: 0,
		TopP:        1,
		N:           1,
		Stream:      false,
		Logprobs:    0,
		Stop:        "\\n",
	}

	log.Debug().Msgf("request: %+v\n", request)

	gpt3 := openai.NewOpenAI(secret)

	s := spinner.New(spinner.CharSets[1], 100*time.Millisecond)
	s.Prefix = "  "           // Append text before the spinner
	s.Suffix = "  loading..." // Append text after the spinner
	s.Color("magenta", "bold")
	s.Start()
	res, err := gpt3.Completion(&request)
	s.Stop()

	if err != nil {
		log.Fatal().Err(err)
	}
	fmt.Println(res)
	fmt.Println(res.Choices[0].Text)
}
