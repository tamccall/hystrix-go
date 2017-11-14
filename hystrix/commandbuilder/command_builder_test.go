package commandbuilder

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/myteksi/hystrix-go/hystrix"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCommandBuilderDefaults(t *testing.T) {
	Convey("given a command configured for a default queue", t, func() {
		commandSetting := New("command1").Build()
		hystrix.Initialize(commandSetting)

		Convey("reading the timeout should be the same", func() {
			circuits := hystrix.GetCircuitSettings()
			So(circuits["command1"].QueueSizeRejectionThreshold, ShouldEqual, hystrix.DefaultQueueSizeRejectionThreshold)
			So(circuits["command1"].ErrorPercentThreshold, ShouldEqual, hystrix.DefaultErrorPercentThreshold)
			So(circuits["command1"].Timeout.Nanoseconds()/1000000, ShouldEqual, hystrix.DefaultTimeout)
			So(circuits["command1"].SleepWindow.Nanoseconds()/1000000, ShouldEqual, hystrix.DefaultSleepWindow)
		})
	})
}

func TestCommandBuilderInvalidInput(t *testing.T) {
	Convey("given a command configured for a default queue", t, func() {
		commandSetting := New("command1").WithErrorPercentageThreshold(0).WithSleepWindow(0).Build()
		hystrix.Initialize(commandSetting)

		Convey("reading the timeout should be the same", func() {
			circuits := hystrix.GetCircuitSettings()
			So(circuits["command1"].QueueSizeRejectionThreshold, ShouldEqual, hystrix.DefaultQueueSizeRejectionThreshold)
			So(circuits["command1"].ErrorPercentThreshold, ShouldEqual, hystrix.DefaultErrorPercentThreshold)
			So(circuits["command1"].Timeout.Nanoseconds()/1000000, ShouldEqual, hystrix.DefaultTimeout)
			So(circuits["command1"].SleepWindow.Nanoseconds()/1000000, ShouldEqual, hystrix.DefaultSleepWindow)
		})
	})
}

func TestCommandBuilderWithCommandGroup(t *testing.T) {
	Convey("given a command configured for a default queue", t, func() {
		commandSetting := New("command2").WithMaxConcurrentRequests(41).WithCommandGroup("service1").Build()
		hystrix.Initialize(commandSetting)

		Convey("reading the timeout should be the same", func() {
			circuits := hystrix.GetCircuitSettings()
			So(circuits["command2"].MaxConcurrentRequests, ShouldEqual, 41)
			So(circuits["command2"].QueueSizeRejectionThreshold, ShouldEqual, 41*5)
			So(circuits["command2"].CommandGroup, ShouldEqual, "service1")
			So(circuits["command2"].ErrorPercentThreshold, ShouldEqual, hystrix.DefaultErrorPercentThreshold)
			So(circuits["command2"].Timeout.Nanoseconds()/1000000, ShouldEqual, hystrix.DefaultTimeout)
			So(circuits["command2"].SleepWindow.Nanoseconds()/1000000, ShouldEqual, hystrix.DefaultSleepWindow)
		})
	})
}

func TestCommandBuilderNoQueue(t *testing.T) {
	Convey("given a command configured for a default queue", t, func() {
		commandSetting := New("command3").WithQueueSize(0).Build()
		hystrix.Initialize(commandSetting)

		Convey("reading the timeout should be the same", func() {
			circuits := hystrix.GetCircuitSettings()
			So(circuits["command3"].QueueSizeRejectionThreshold, ShouldEqual, 0)
			So(circuits["command3"].ErrorPercentThreshold, ShouldEqual, hystrix.DefaultErrorPercentThreshold)
			So(circuits["command3"].Timeout.Nanoseconds()/1000000, ShouldEqual, hystrix.DefaultTimeout)
			So(circuits["command3"].SleepWindow.Nanoseconds()/1000000, ShouldEqual, hystrix.DefaultSleepWindow)
		})
	})
}

func TestOverflowWithoutQueue(t *testing.T) {
	defer hystrix.Flush()

	Convey("testing for overflow without queue", t, func() {
		command := New("command1").WithQueueSize(0).WithMaxConcurrentRequests(10).WithTimeout(1000).Build()
		hystrix.Initialize(command)

		maxConcurrencyErr := int32(0)
		success := int32(0)
		completedAll := make(chan struct{})
		totalExecution := int32(0)

		for i := 0; i < 25; i++ {
			go func(idx int) {
				resChan := make(chan *struct{}, 1)
				errChan := hystrix.Go("command1", func() error {
					time.Sleep(300 * time.Millisecond)
					resChan <- &struct{}{}
					return nil
				}, nil)

				var err error
				select {
				case err = <-errChan:
				case _ = <-resChan:
					err = nil
				}
				if err == nil {
					atomic.AddInt32(&success, 1)
				}
				if err == hystrix.ErrMaxConcurrency {
					atomic.AddInt32(&maxConcurrencyErr, 1)
				}
				total := atomic.AddInt32(&totalExecution, 1)

				if total == 25 {
					close(completedAll)
				}
			}(i)
			time.Sleep(10 * time.Millisecond)
		}
		Convey("number of success, timeout maxConn is correct", func() {
			<-completedAll
			So(success, ShouldEqual, 10)
			So(maxConcurrencyErr, ShouldEqual, 15)

		})

	})
}
