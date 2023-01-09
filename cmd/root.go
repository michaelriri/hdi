package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/michaelriri/hdi/openai"

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
	Use:   "hdi <question>...",
	Short: "A GPT-3 powered command line assistant",

	Example: `$ hdi remove all docker containers
$ hdi execute a python script and end up in a repl
$ hdi view the contents of a file`,

	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		argString := strings.Join(args, " ")
		log := setupLogging(verboseCount)
		log.Debug().Msg("Arguments: " + argString)
		log.Debug().Msgf("verboseCount: %d", verboseCount)

		secret, ok := os.LookupEnv("OPENAI_API_KEY")
		if !ok {
			fmt.Println("Required environment variable OPENAI_API_KEY is not set")
			fmt.Println("See https://beta.openai.com/docs/api-reference/authentication for more information")
			os.Exit(1)
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
