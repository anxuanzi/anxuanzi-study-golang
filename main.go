package main

import (
	"az-ops/SystemInitialization"
	"github.com/gookit/gcli/v2"
)

func main() {
	geniusTonyOpsApp := gcli.NewApp()
	geniusTonyOpsApp.Name = "GeniusTony OPS Toolbox"
	geniusTonyOpsApp.Description = "[GeniusTony OPS Toolbox] GeniusTonyAn server operation and maintenance toolkit, professional, comprehensive and concise server rapid configuration and software deployment tools."
	geniusTonyOpsApp.Logo = gcli.Logo{
		Text:  "   _____                  _                 _______                                                 ____    _____     _____ \n  / ____|                (_)               |__   __|                             /\\                / __ \\  |  __ \\   / ____|\n | |  __    ___   _ __    _   _   _   ___     | |      ___    _ __    _   _     /  \\     _ __     | |  | | | |__) | | (___  \n | | |_ |  / _ \\ | '_ \\  | | | | | | / __|    | |     / _ \\  | '_ \\  | | | |   / /\\ \\   | '_ \\    | |  | | |  ___/   \\___ \\ \n | |__| | |  __/ | | | | | | | |_| | \\__ \\    | |    | (_) | | | | | | |_| |  / ____ \\  | | | |   | |__| | | |       ____) |\n  \\_____|  \\___| |_| |_| |_|  \\__,_| |___/    |_|     \\___/  |_| |_|  \\__, | /_/    \\_\\ |_| |_|    \\____/  |_|      |_____/ \n                                                                       __/ |                                                \n                                                                      |___/                                                 \n",
		Style: "info",
	}
	geniusTonyOpsApp.Version = "1.0.0"
	//debug
	//geniusTonyOpsApp.SetDebugMode()
	geniusTonyOpsApp.AddCommand(SystemInitialization.InitSystemCommand)
	geniusTonyOpsApp.Run()
}

//func main() {
//	Global.InitResources()
//	http.Handle("/", http.StripPrefix("/", http.FileServer(Global.GetStatikFS())))
//	http.ListenAndServe(":8787", nil)
//}
