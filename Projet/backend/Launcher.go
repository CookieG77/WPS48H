package backend

import (
	pages "Proj48h/backend/pages"
	"Proj48h/functions"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
)

var Port = "8080"

// LaunchWebApp launch the web application.
// It will listen on the port 8080 by default.
// If the argument '--port [port]' is provided, it will listen on the specified port.
// If the port is not a valid number, the program will exit.
func LaunchWebApp() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				functions.RmTempDir()
				functions.ClearCmd()
				os.Exit(1)
			}
		}
	}()

	// Managing the arguments
	Args := functions.GetArgs()
	if slices.Contains(Args, "--port") || slices.Contains(Args, "-p") {
		// Check if the port is provided
		if len(Args) > slices.Index(Args, "--port")+1 {
			rawProposedPort := Args[slices.Index(Args, "--port")+1]
			proposedPort, err := strconv.Atoi(rawProposedPort)

			if err != nil || proposedPort < 1 || proposedPort > 65535 {
				functions.FatalPrintln("The port provided is not a valid number. Please provide a valid port number. (1-65535)")
			}
			Port = rawProposedPort
		} else if len(Args) > slices.Index(Args, "-p")+1 {
			rawProposedPort := Args[slices.Index(Args, "-p")+1]
			proposedPort, err := strconv.Atoi(rawProposedPort)

			if err != nil || proposedPort < 1 || proposedPort > 65535 {
				functions.FatalPrintln("The port provided is not a valid number. Please provide a valid port number. (1-65535)")
			}
			Port = rawProposedPort
		}
	}
	if slices.Contains(Args, "--show-info") || slices.Contains(Args, "-si") {
		functions.ShouldLogInfo = true
	}

	// Create the temporary directory
	functions.MkTempDir()

	// Managing the static files
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./statics/css"))))
	http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./statics/img"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./statics/js"))))

	// Managing the pages
	http.HandleFunc("/", pages.HomePage)
	http.HandleFunc("/report", pages.ReportPage)
	http.HandleFunc("/download", pages.DownloadFile)
	http.HandleFunc("/send-mail-for-download", pages.SendByMail)

	// Set the port to listen on and initialize the mail service with the configuration file
	finalPort := ""
	if Port != "" {
		finalPort = ":" + Port
	} else {
		finalPort = ":8080"
	}
	functions.InitMail("MailConfig.json")

	functions.SuccessPrintf("Server started at -> http://localhost%s\n", finalPort)

	// Launch the server
	if err := http.ListenAndServe(finalPort, nil); err != nil {
		panic(err)
	}
}
