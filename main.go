package main

import (
    "fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/model"
	"github.com/conductor-sdk/conductor-go/sdk/settings"

	"github.com/conductor-sdk/conductor-go/sdk/worker"
	"github.com/conductor-sdk/conductor-go/sdk/workflow/executor"
)

var (
	apiClient = client.NewAPIClient(
		authSettings(),
		httpSettings(),
	)
	taskRunner       = worker.NewTaskRunnerWithApiClient(apiClient)
	workflowExecutor = executor.NewWorkflowExecutor(apiClient)
)

func authSettings() *settings.AuthenticationSettings {
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	if key != "" && secret != "" {
		return settings.NewAuthenticationSettings(
			key,
			secret,
		)
	}

	return nil
}

func httpSettings() *settings.HttpSettings {
	url :=  os.Getenv("CONDUCTOR_SERVER_URL")
	if url == "" {
		log.Error("Error: CONDUCTOR_SERVER_URL env variable is not set")
		os.Exit(1)
	}

	return settings.NewHttpSettings(url)
}

func Greet(task *model.Task) (result interface{}, err error) {
	taskResult := model.NewTaskResultFromTask(task)
	
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Panic occurred in worker: %v", r)
			taskResult.Status = "FAILED"
			taskResult.ReasonForIncompletion = fmt.Sprintf("Task failed due to panic: %v", r)
			taskResult.OutputData = nil
			result = taskResult // Set the named return value
			err = fmt.Errorf("task failed due to panic: %v", r)
		}
	}()

	// panic("ahhh panicao")
	// Process the task
	greetingMsg := "Hello, " + fmt.Sprintf("%v", task.InputData["name"])
	
	taskResult.Status = "COMPLETED"
	taskResult.OutputData = map[string]interface{}{
		"greetings": greetingMsg,
	}
	
	return taskResult, nil
}

func main() {
	taskRunner.StartWorker("greet", Greet, 1, time.Millisecond*100)
    taskRunner.WaitWorkers();
}

