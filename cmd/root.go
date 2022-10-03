package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/docker"
	"github.com/luminous-gsm/fusion/event"
	"github.com/luminous-gsm/fusion/node"
	"github.com/luminous-gsm/fusion/router/rest"
	"github.com/luminous-gsm/fusion/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "fusion",
	Short: "Runs the fusion API server, allowing controller nodes from the Luminous console",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger, err := zap.NewDevelopment()
		if err != nil {
			log.Fatal("initial logger failed")
		}
		zap.ReplaceGlobals(logger)

		initConfig()
		initDirectories()
		initLogging()

		byteMarshal, _ := json.MarshalIndent(config.Get(), "", "    ")
		fmt.Print(string(byteMarshal))
		zap.S().Infow("configured config", "config", config.Get())
	},
	Run: rootRun,
}

func init() {
	rootCommand.AddCommand(configureCommand)
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatalf("failed to execute command: %s", err)
	}
}

func rootRun(cmd *cobra.Command, _ []string) {
	cfg := config.Get()

	// Create main context
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize the services and their event listeners
	event.InitEventBus().InitEventChannels()
	docker.InitDockerService(ctx).InitEventListeners()
	node.InitNodeService(ctx).InitEventListeners()

	// Create new server manager
	mgr, err := server.NewManager(ctx, cancel)
	if err != nil {
		zap.S().Fatalw("failed to initiate server manager",
			"error", err,
		)
	}

	// Create new router
	router := rest.NewRouter(mgr)

	// Run the HTTP server
	port := fmt.Sprintf("%v:%v", cfg.ApiHost, cfg.ApiPort)
	zap.S().Infow("started api server",
		"host", cfg.ApiHost,
		"port", cfg.ApiPort,
	)
	err = router.RunTLS(port, fmt.Sprintf("%vfusion.crt", config.Get().CertsDirectory), fmt.Sprintf("%vfusion.key", config.Get().CertsDirectory))
	if err != nil {
		zap.S().DPanicw("failed to start TLS server.", "error", err)
	}
}

func initConfig() {
	Configure()
}

func initDirectories() {
	// Create log file
	err := os.MkdirAll(config.Get().LogDirectory, os.ModePerm)
	if err != nil {
		zap.S().Errorw("cmd: cannot create logging directory", "error", err)
	}
	err = os.MkdirAll(config.Get().DataDirectory, os.ModePerm)
	if err != nil {
		zap.S().Errorw("cmd: cannot create data directory", "error", err)
	}
	err = os.MkdirAll(config.Get().RootDirectory, os.ModePerm)
	if err != nil {
		zap.S().Errorw("cmd: cannot create root directory", "error", err)
	}
}

func initLogging() {

	var (
		core                zapcore.Core
		loggerConfigEncoder zapcore.EncoderConfig
	)

	if config.Get().Debug {
		loggerConfigEncoder = zap.NewDevelopmentEncoderConfig()
		// Logging to the console
		loggerConfigEncoder.TimeKey = "timestamp"
		loggerConfigEncoder.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
		loggerConfigEncoder.FunctionKey = "function"
		loggerConfigEncoder.MessageKey = "message"
		loggerConfigEncoder.CallerKey = "caller"
		consoleEncoder := zapcore.NewConsoleEncoder(loggerConfigEncoder)
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		loggerConfigEncoder = zap.NewProductionEncoderConfig()
		// Logging to the console
		consoleEncoder := zapcore.NewConsoleEncoder(loggerConfigEncoder)
		loggerConfigEncoder.TimeKey = "timestamp"
		loggerConfigEncoder.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
		loggerConfigEncoder.FunctionKey = "function"
		loggerConfigEncoder.MessageKey = "message"
		loggerConfigEncoder.CallerKey = "caller"

		// Logging as JSON to file
		fileEncoder := zapcore.NewJSONEncoder(loggerConfigEncoder)
		logFile, err := os.OpenFile(config.Get().LogDirectory+"fusion.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_RDWR, 0644)
		if err != nil {
			log.Fatal("cmd: failed to create/open fusion log file")
		}
		writer := zapcore.AddSync(logFile)
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writer, zapcore.InfoLevel),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		)
	}

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	zap.ReplaceGlobals(logger)
	zap.S().Info("configured global logger")

	zap.S().Debug("running in debug mode")
}
