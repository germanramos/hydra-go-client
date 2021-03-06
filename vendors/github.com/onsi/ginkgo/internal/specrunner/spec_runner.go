package specrunner

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo/config"
	"github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo/internal/leafnodes"
	"github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo/internal/spec"
	Writer "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo/internal/writer"
	"github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo/reporters"
	"github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo/types"

	"time"
)

type SpecRunner struct {
	description	string
	beforeSuiteNode	leafnodes.SuiteNode
	specs		*spec.Specs
	afterSuiteNode	leafnodes.SuiteNode
	reporters	[]reporters.Reporter
	startTime	time.Time
	suiteID		string
	runningSpec	*spec.Spec
	writer		Writer.WriterInterface
	config		config.GinkgoConfigType
	interrupted	bool
	lock		*sync.Mutex
}

func New(description string, beforeSuiteNode leafnodes.SuiteNode, specs *spec.Specs, afterSuiteNode leafnodes.SuiteNode, reporters []reporters.Reporter, writer Writer.WriterInterface, config config.GinkgoConfigType) *SpecRunner {
	return &SpecRunner{
		description:		description,
		beforeSuiteNode:	beforeSuiteNode,
		specs:			specs,
		afterSuiteNode:		afterSuiteNode,
		reporters:		reporters,
		writer:			writer,
		config:			config,
		suiteID:		randomID(),
		lock:			&sync.Mutex{},
	}
}

func (runner *SpecRunner) Run() bool {
	runner.reportSuiteWillBegin()
	go runner.registerForInterrupts()

	suitePassed := runner.runBeforeSuite()

	if suitePassed {
		suitePassed = runner.runSpecs()
	}

	runner.blockForeverIfInterrupted()

	suitePassed = runner.runAfterSuite() && suitePassed

	runner.reportSuiteDidEnd(suitePassed)

	return suitePassed
}

func (runner *SpecRunner) runBeforeSuite() bool {
	if runner.beforeSuiteNode == nil || runner.wasInterrupted() {
		return true
	}

	runner.writer.Truncate()
	conf := runner.config
	passed := runner.beforeSuiteNode.Run(conf.ParallelNode, conf.ParallelTotal, conf.SyncHost)
	if !passed {
		runner.writer.DumpOut()
	}
	runner.reportBeforeSuite(runner.beforeSuiteNode.Summary())
	return passed
}

func (runner *SpecRunner) runAfterSuite() bool {
	if runner.afterSuiteNode == nil {
		return true
	}

	runner.writer.Truncate()
	conf := runner.config
	passed := runner.afterSuiteNode.Run(conf.ParallelNode, conf.ParallelTotal, conf.SyncHost)
	if !passed {
		runner.writer.DumpOut()
	}
	runner.reportAfterSuite(runner.afterSuiteNode.Summary())
	return passed
}

func (runner *SpecRunner) runSpecs() bool {
	suiteFailed := false
	skipRemainingSpecs := false
	for _, spec := range runner.specs.Specs() {
		if runner.wasInterrupted() {
			return suiteFailed
		}
		if skipRemainingSpecs {
			spec.Skip()
		}
		runner.reportSpecWillRun(spec)

		if !spec.Skipped() && !spec.Pending() {
			runner.runningSpec = spec
			spec.Run()
			runner.runningSpec = nil
			if spec.Failed() {
				suiteFailed = true
			}
		} else if spec.Pending() && runner.config.FailOnPending {
			suiteFailed = true
		}

		runner.reportSpecDidComplete(spec)

		if spec.Failed() && runner.config.FailFast {
			skipRemainingSpecs = true
		}
	}

	return !suiteFailed
}

func (runner *SpecRunner) CurrentSpecSummary() (*types.SpecSummary, bool) {
	if runner.runningSpec == nil {
		return nil, false
	}

	return runner.runningSpec.Summary(runner.suiteID), true
}

func (runner *SpecRunner) registerForInterrupts() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	signal.Stop(c)
	runner.markInterrupted()
	go runner.registerForHardInterrupts()
	if runner.afterSuiteNode != nil {
		fmt.Fprintln(os.Stderr, "\nReceived interrupt.  Running AfterSuite...\n^C again to terminate immediately")
		runner.runAfterSuite()
	}
	runner.reportSuiteDidEnd(false)
	os.Exit(1)
}

func (runner *SpecRunner) registerForHardInterrupts() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Fprintln(os.Stderr, "\nReceived second interrupt.  Shutting down.")
	os.Exit(1)
}

func (runner *SpecRunner) blockForeverIfInterrupted() {
	runner.lock.Lock()
	interrupted := runner.interrupted
	runner.lock.Unlock()

	if interrupted {
		select {}
	}
}

func (runner *SpecRunner) markInterrupted() {
	runner.lock.Lock()
	defer runner.lock.Unlock()
	runner.interrupted = true
}

func (runner *SpecRunner) wasInterrupted() bool {
	runner.lock.Lock()
	defer runner.lock.Unlock()
	return runner.interrupted
}

func (runner *SpecRunner) reportSuiteWillBegin() {
	runner.startTime = time.Now()
	summary := runner.summary(true)
	for _, reporter := range runner.reporters {
		reporter.SpecSuiteWillBegin(runner.config, summary)
	}
}

func (runner *SpecRunner) reportBeforeSuite(summary *types.SetupSummary) {
	for _, reporter := range runner.reporters {
		reporter.BeforeSuiteDidRun(summary)
	}
}

func (runner *SpecRunner) reportAfterSuite(summary *types.SetupSummary) {
	for _, reporter := range runner.reporters {
		reporter.AfterSuiteDidRun(summary)
	}
}

func (runner *SpecRunner) reportSpecWillRun(spec *spec.Spec) {
	runner.writer.Truncate()

	summary := spec.Summary(runner.suiteID)
	for _, reporter := range runner.reporters {
		reporter.SpecWillRun(summary)
	}
}

func (runner *SpecRunner) reportSpecDidComplete(spec *spec.Spec) {
	summary := spec.Summary(runner.suiteID)
	for i := len(runner.reporters) - 1; i >= 1; i-- {
		runner.reporters[i].SpecDidComplete(summary)
	}

	if spec.Failed() {
		runner.writer.DumpOut()
	}

	runner.reporters[0].SpecDidComplete(summary)
}

func (runner *SpecRunner) reportSuiteDidEnd(success bool) {
	summary := runner.summary(success)
	summary.RunTime = time.Since(runner.startTime)
	for _, reporter := range runner.reporters {
		reporter.SpecSuiteDidEnd(summary)
	}
}

func (runner *SpecRunner) countSpecsSatisfying(filter func(ex *spec.Spec) bool) (count int) {
	count = 0

	for _, spec := range runner.specs.Specs() {
		if filter(spec) {
			count++
		}
	}

	return count
}

func (runner *SpecRunner) summary(success bool) *types.SuiteSummary {
	numberOfSpecsThatWillBeRun := runner.countSpecsSatisfying(func(ex *spec.Spec) bool {
		return !ex.Skipped() && !ex.Pending()
	})

	numberOfPendingSpecs := runner.countSpecsSatisfying(func(ex *spec.Spec) bool {
		return ex.Pending()
	})

	numberOfSkippedSpecs := runner.countSpecsSatisfying(func(ex *spec.Spec) bool {
		return ex.Skipped()
	})

	numberOfPassedSpecs := runner.countSpecsSatisfying(func(ex *spec.Spec) bool {
		return ex.Passed()
	})

	numberOfFailedSpecs := runner.countSpecsSatisfying(func(ex *spec.Spec) bool {
		return ex.Failed()
	})

	if runner.beforeSuiteNode != nil && !runner.beforeSuiteNode.Passed() {
		numberOfFailedSpecs = numberOfSpecsThatWillBeRun
	}

	return &types.SuiteSummary{
		SuiteDescription:	runner.description,
		SuiteSucceeded:		success,
		SuiteID:		runner.suiteID,

		NumberOfSpecsBeforeParallelization:	runner.specs.NumberOfOriginalSpecs(),
		NumberOfTotalSpecs:			len(runner.specs.Specs()),
		NumberOfSpecsThatWillBeRun:		numberOfSpecsThatWillBeRun,
		NumberOfPendingSpecs:			numberOfPendingSpecs,
		NumberOfSkippedSpecs:			numberOfSkippedSpecs,
		NumberOfPassedSpecs:			numberOfPassedSpecs,
		NumberOfFailedSpecs:			numberOfFailedSpecs,
	}
}
