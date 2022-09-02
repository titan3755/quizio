package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"github.com/probandula/figlet4go"
)

const welcomeText = "Welcome to Terminal quiz app! Here, you can participate in a quiz by reading data from an API or get quiz data from a JSON file in the same directory. Use the options below to select a JSON file, create a quiz-ready JSON file or use an API by mentioning the URL (note that if you use the API method, the response must be in the correct data format.)"
const jsonFormatText = `
[
	{
		question: "What is the age of the universe?",
		options: [
			"10.9 Bn Years",
			"8 Bn Years",
			"16.2 Bn Years",
			"13.8 Bn Years"
		],
		correct: 3
	}
]
`
var appOptions = [...]string{"1. Use quiz JSON file", "2. Create quiz JSON file", "3. Use API quiz data (not available yet)", "4. View JSON data format", "5. Quit app"}
var reader = bufio.NewReader(os.Stdin)

func main() {
	resetTerminal()
	welcomeASCIIText("Quiz-APP")
	fmt.Println(welcomeText)
	for {
		fmt.Print("\n---Menu---\n")
		fmt.Printf("\n%v\n%v\n%v\n%v\n%v\n\n", appOptions[0], appOptions[1], appOptions[2], appOptions[3], appOptions[4])
		fmt.Print("--> ")
		text, _ := reader.ReadString('\n')
		switch userResponse := strings.TrimSpace(text); userResponse {
		case "1":
			fileReadMode()
		case "2":
			fileWriteMode()
		case "3":
			apiMode()
		case "4":
			jsonDataFormatOutput()
		case "5":
			quitApp()
		default:
			fmt.Println("Invalid Mode!")
		}
		endPrompt()
	}
}

func fileReadMode() {
	resetTerminal()
	color.Cyan("You have selected \"JSON file read mode\". Mention the file name of the JSON file in the same directory to proceed.\n")
	fmt.Print("\nFile Name: ")
	text, _ := reader.ReadString('\n')
	data, err := os.ReadFile(strings.TrimSpace(text))
	if err != nil {
		color.Red("Something went wrong when reading the file! The specified file may not exist in the current directory or the name is wrong. \nerr: [%v]", err)
		return
	}
	fmt.Printf("\nContents of file: %v\n", string(data))
}

func fileWriteMode() {
	resetTerminal()
	color.Cyan("You have selected \"JSON file write mode\". Mention the file name of the JSON file that's to be created in the same directory to proceed.\n")
	fmt.Print("\nFile Name: ")
	text, _ := reader.ReadString('\n')
	file, err := os.Create(strings.TrimSpace(text) + ".json")
	if err != nil {
		color.Red("Something went wrong when creating the file! The specified file may already exist in the current directory or the name is invalid. \nerr: [%v]", err)
		return
	}
	color.Green("File created successfully!")
	fmt.Print("\nAdd quiz questions to new file? (y/n) ")
	quizPromptResponse, _ := reader.ReadString('\n') 
	if strings.TrimSpace(quizPromptResponse) == "y" {
		questionWriter()
	}
	defer file.Close()
}

func apiMode() {
	resetTerminal()
}

func jsonDataFormatOutput() {
	resetTerminal()
	color.Cyan("You have selected \"View JSON data format mode\". You can see the JSON data format which is required for the quiz app to work properly below. Note that the \"correct\" key uses array indexes which means that the value is one less than the numerical position of the correct option in \"options\" key.\n\n")
	err := quick.Highlight(os.Stdout, jsonFormatText, "", "", "monokai")
	if err != nil {
		color.Red("Something went wrong during syntax highlighting! \nerr: [%v]", err)
		return
	}
}

func quitApp() {
	resetTerminal()
	color.Red("Quitting App ...")
	time.Sleep(999999999)
	os.Exit(0)
}

func welcomeASCIIText(txt string) {
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorCyan,
	}
	renderStr, _ := ascii.RenderOpts(txt, options)
	fmt.Println(renderStr)
}

func questionWriter() {
	resetTerminal()
	color.Cyan("Welcome to question-writer! Here you can edit your newly created or existing JSON quiz file from within the app and you won't need any fancy editors to get the job done. Just specify the filename or path and you'll be all set!")
	fmt.Printf("File name/path: ")
	text, _ := reader.ReadString('\n')
	fmt.Println("\nPath: " + text)
}

func endPrompt() {
	fmt.Print(color.YellowString("\nReturn to menu? (y/n) "))
	text, _ := reader.ReadString('\n')
	if strings.TrimSpace(text) != "y" {
		quitApp()
	}
	resetTerminal()
}

func questionInit(data string) {
	resetTerminal()
}

func resetTerminal() {
	screen.Clear()
	screen.MoveTopLeft()
}