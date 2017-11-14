package commandbuilder

import (
	"time"

	"github.com/myteksi/hystrix-go/hystrix"
)

// CommandBuilder builder for constructing new command

type CommandBuilder struct {
	commandName            string
	timeout                int
	commandGroup           string
	maxConcurrentRequests  int
	requestVolumeThreshold int
	sleepWindow            int
	errorPercentThreshold  int
	// for more details refer - https://github.com/Netflix/Hystrix/wiki/Configuration#maxqueuesize
	queueSizeRejectionThreshold *int
}

// New Create new command
func New(commandName string) *CommandBuilder {
	return &CommandBuilder{
		commandName:            commandName,
		timeout:                hystrix.DefaultTimeout,
		commandGroup:           commandName,
		maxConcurrentRequests:  hystrix.DefaultMaxConcurrent,
		requestVolumeThreshold: hystrix.DefaultVolumeThreshold,
		sleepWindow:            hystrix.DefaultSleepWindow,
		errorPercentThreshold:  hystrix.DefaultErrorPercentThreshold,
	}
}

// WithTimeout modify timeout
func (cb *CommandBuilder) WithTimeout(timeoutInMs int) *CommandBuilder {
	if timeoutInMs == 0 {
		timeoutInMs = hystrix.DefaultTimeout
	}
	cb.timeout = timeoutInMs
	return cb
}

// WithCommandGroup modify commandGroup
func (cb *CommandBuilder) WithCommandGroup(commandGroup string) *CommandBuilder {
	if commandGroup == "" {
		commandGroup = cb.commandName
	}
	cb.commandGroup = commandGroup
	return cb
}

// WithMaxConcurrentRequests modify max concurrent requests
// if not already set, this will also set the queue size as 5 times the max concurrent requests
func (cb *CommandBuilder) WithMaxConcurrentRequests(maxConcurrentRequests int) *CommandBuilder {
	if maxConcurrentRequests == 0 {
		maxConcurrentRequests = hystrix.DefaultMaxConcurrent
	}
	cb.maxConcurrentRequests = maxConcurrentRequests
	if cb.queueSizeRejectionThreshold == nil {
		queueSize := 5 * cb.maxConcurrentRequests
		cb.queueSizeRejectionThreshold = &queueSize
	}
	return cb
}

// WithRequestVolumeThreshold modify request volume threshold
func (cb *CommandBuilder) WithRequestVolumeThreshold(requestVolThreshold int) *CommandBuilder {
	if requestVolThreshold == 0 {
		requestVolThreshold = hystrix.DefaultVolumeThreshold
	}
	cb.requestVolumeThreshold = requestVolThreshold
	return cb
}

// WithSleepWindow modify sleep window
func (cb *CommandBuilder) WithSleepWindow(sleepWindow int) *CommandBuilder {
	if sleepWindow == 0 {
		sleepWindow = hystrix.DefaultSleepWindow
	}
	cb.sleepWindow = sleepWindow
	return cb
}

// WithErrorPercentageThreshold modify error percentage threshold
func (cb *CommandBuilder) WithErrorPercentageThreshold(errPercentThreshold int) *CommandBuilder {
	if errPercentThreshold == 0 {
		errPercentThreshold = hystrix.DefaultErrorPercentThreshold
	}
	cb.errorPercentThreshold = errPercentThreshold
	return cb
}

// WithQueueSize modify queue size
func (cb *CommandBuilder) WithQueueSize(queueSize int) *CommandBuilder {
	if queueSize == 0 {
		zeroQueueSize := 0
		cb.queueSizeRejectionThreshold = &zeroQueueSize
	}
	cb.queueSizeRejectionThreshold = &queueSize
	return cb
}

// Build the command setting, Use hystrix.Initialize for setup
func (cb *CommandBuilder) Build() *hystrix.Settings {
	if cb.queueSizeRejectionThreshold == nil {
		cb.queueSizeRejectionThreshold = &hystrix.DefaultQueueSizeRejectionThreshold
	}
	return &hystrix.Settings{
		CommandName:                 cb.commandName,
		QueueSizeRejectionThreshold: *cb.queueSizeRejectionThreshold,
		ErrorPercentThreshold:       cb.errorPercentThreshold,
		CommandGroup:                cb.commandGroup,
		Timeout:                     time.Duration(cb.timeout * 1000000),
		MaxConcurrentRequests:       cb.maxConcurrentRequests,
		RequestVolumeThreshold:      uint64(cb.requestVolumeThreshold),
		SleepWindow:                 time.Duration(cb.sleepWindow * 1000000),
	}
}
