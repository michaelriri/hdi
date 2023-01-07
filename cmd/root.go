package cmd

import (
	"fmt"
	"hdi/openai"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var verboseCount int

func setupLogging(level int) zerolog.Logger {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	switch level {
	case 0:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	return log
}

var rootCmd = &cobra.Command{
	Use:   "hdi {question}...",
	Short: "A GPT-3 powered command line assistant",
	Long: `hdi ("how do I?") is a GPT-3 powered command line assistant.

It uses the OpenAI API to answer questions like 'How do I remove all docker containers?'
The GPT-3 model's answer is then shown in your shell.
Info on the OpenAI API can be found at: https://openai.com/api`,

	Example: `hdi remove all docker containers
	docker rm $(docker ps -aq)`,
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		argString := strings.Join(args, " ")
		log := setupLogging(verboseCount)
		log.Debug().Msg("Arguments: " + argString)
		log.Debug().Msgf("verboseCount: %d", verboseCount)

		secret, ok := os.LookupEnv("OPENAI_API_KEY")
		if !ok {
			log.Fatal().Msg("OPENAI_API_KEY not set")
		}
		log.Info().Str("secret", secret).Msg("using OPENAI_API_KEY")

		question := strings.Join(args, " ")
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
		s.Prefix = "  "
		s.Suffix = "  loading..."
		s.Color("magenta", "bold")
		s.Start()
		res, err := gpt3.Completion(&request)
		s.Stop()

		if err != nil {
			log.Fatal().Err(err)
		}
		log.Debug().Msgf("response: %+v\n", res)
		log.Debug().Msgf("response.choices: %+v\n", res.Choices)

		if len(res.Choices) == 0 {
			log.Fatal().Msg("no choices returned")
		}
		fmt.Println(res.Choices[0].Text)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().CountVarP(&verboseCount, "verbose", "v", "verbose output")
}
