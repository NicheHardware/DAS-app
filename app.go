package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	StartCM01Manager()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) LogPrintln(log string) (len int) {
	len, _ = fmt.Println(log)
	return len
}

func (a *App) StartSim(count int) {
	PauseExternalAcquisition()
	StopSim()
	StartSim(count)
}

func (a *App) StopSim() {
	StopSim()
	ResumeExternalAcquisition()
}

func (a *App) NewDataNotify() {
	runtime.EventsEmit(a.ctx, "new-data", nodes)
}
