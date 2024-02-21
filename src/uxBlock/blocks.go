// Package uxBlock provides building blocks for UX and communication with a user.
package uxBlock

import (
	"context"

	"github.com/zeropsio/zcli/src/logger"
	"github.com/zeropsio/zcli/src/uxBlock/styles"
)

//go:generate go run --mod=mod github.com/golang/mock/mockgen -source=$GOFILE -destination=$PWD/mocks/$GOFILE -package=mocks

type UxBlocks interface {
	LogDebug(message string)
	PrintInfo(line styles.Line)
	PrintWarning(line styles.Line)
	PrintError(line styles.Line)
	Table(body *TableBody, auxOptions ...TableOption)
	Select(ctx context.Context, tableBody *TableBody, auxOptions ...SelectOption) ([]int, error)
	Prompt(
		ctx context.Context,
		message string,
		choices []string,
		auxOptions ...PromptOption,
	) (int, error)
	RunSpinners(ctx context.Context, spinners []*Spinner, auxOptions ...SpinnerOption) func()
}

type uxBlocks struct {
	outputLogger    logger.Logger
	debugFileLogger logger.Logger
	isTerminal      bool
	terminalWidth   int

	// ctxCancel is used to cancel the context of the command.
	// Bubbles package that we use to render graphic components steals the signal handler.
	// In case that I want to cancel the context of a running component, e.g. spinner, the main context is not canceled.
	// Therefore, we need to pass the cancel function to the uxBlocks.
	// If the graphic component is canceled, we cancel the main context.
	ctxCancel context.CancelFunc
}

func NewBlock(
	outputLogger logger.Logger,
	debugFileLogger logger.Logger,
	isTerminal bool,
	terminalWidth int,
	ctxCancel context.CancelFunc,
) *uxBlocks {
	// safety check
	if ctxCancel == nil {
		ctxCancel = func() {}
	}

	return &uxBlocks{
		outputLogger:    outputLogger,
		debugFileLogger: debugFileLogger,
		isTerminal:      isTerminal,
		terminalWidth:   terminalWidth,
		ctxCancel:       ctxCancel,
	}
}
