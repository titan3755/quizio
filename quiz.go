package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"runtime"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/fatih/color"
	"github.com/imroc/req/v3"
	"github.com/inancgumus/screen"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/probandula/figlet4go"
	"github.com/tidwall/gjson"
)

type QuizQuestion struct {
	Question string 
	Options []interface{} 
	Correct int 
}

const welcomeText = "Welcome to Terminal quiz app! Here, you can participate in a quiz by reading data from an API or get quiz data from a JSON file in the same directory. Use the options below to select a JSON file, create a quiz-ready JSON file or use an API by mentioning the URL (note that if you use the API method, the response must be in the correct data format.)"
const jsonFormatText = `
[
	{
		"Question": "What is the age of the universe?",
		"Options": [
			"10.9 Bn Years",
			"8 Bn Years",
			"16.2 Bn Years",
			"13.8 Bn Years"
		],
		"Correct": 3
	}
]
`
const emptyJsonFileStartingText = "[]"
var appOptions = [...]string{"1. Use quiz JSON file", "2. Create quiz JSON file", "3. Use API quiz data", "4. Initialize quiz-writer", "5. View JSON data format", "6. Quit app"}
var reader = bufio.NewReader(os.Stdin)
var requestClient = req.C()
var menuInitTimes int = 0

func main() {
	for {
		resetTerminal()
		ClearTerminal()
		if menuInitTimes == 0 {
			welcomeASCIIText("Quiz-APP", "init")
			fmt.Println(welcomeText)
			fmt.Print("\n---Menu---\n")
		} else {
			welcomeASCIIText("QUIZ-MENU", "menu")
		}
		fmt.Printf("\n%v\n%v\n%v\n%v\n%v\n%v\n\n", appOptions[0], appOptions[1], appOptions[2], appOptions[3], appOptions[4], appOptions[5])
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
			questionWriter()
		case "5":
			jsonDataFormatOutput()
		case "6":
			quitApp()
		default:
			fmt.Println("Invalid Mode!")
		}
		menuInitTimes++
		endPrompt()
	}
}

func fileReadMode() {
	resetTerminal()
	color.Cyan("You have selected \"JSON file read mode\". Mention the file name of the JSON file in the same directory or the absolute path to a JSON file in a different directory to proceed.\n")
	fmt.Print("\nFile Name: ")
	text, _ := reader.ReadString('\n')
	data, err := os.ReadFile(strings.TrimSpace(text))
	if err != nil {
		color.Red("Something went wrong when reading the file! The specified file may not exist in the given directory or there might be something wrong with the path. \nerr: [%v]", err)
		return
	}
	questionInit(string(data))
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
	color.Cyan("You have selected \"API quiz data mode\". Mention the API URL in the input below and the quiz will be initialized automatically.\n")
	fmt.Print("\nAPI URL: ")
	text, _ := reader.ReadString('\n')
	color.Yellow("\nGetting data ...")
	resp, err := requestClient.R().Get(strings.TrimSpace(text))
	if err != nil {
		color.Red("Something went wrong when querying the API! The API may not provide data in the correct format or is unavailable or the URL may be invalid. \nerr: [%v]", err)
		return
	}
	questionInit(string(strings.TrimSpace(resp.String())))
}

func jsonDataFormatOutput() {
	resetTerminal()
	color.Cyan("You have selected \"View JSON data format mode\". You can see the JSON data format which is required for the quiz app to work properly below. Note that the \"correct\" key uses array indexes which means that the value is one less than the numerical position of the correct option in \"options\" key. You can add more quiz questions by appending more objects to the outermost array in the required format with \"question\", \"options\" and \"correct\" keys. The \"options\" key will be another array containing four options as strings (text enclosed in double quotes) and the options should be separated by commas.\n\n")
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

func welcomeASCIIText(txt string, state string) {
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	if state == "init" {
		options.FontColor = []figlet4go.Color{
			figlet4go.ColorCyan,
		}
	} else if state == "menu" {
		options.FontColor = []figlet4go.Color{
			figlet4go.ColorMagenta,
		}
	} else {
		options.FontColor = []figlet4go.Color{
			figlet4go.ColorWhite,
		}
	}
	renderStr, _ := ascii.RenderOpts(txt, options)
	fmt.Println(renderStr)
}

func questionWriter() {
	resetTerminal()
	color.Cyan("Welcome to question-writer! Here you can edit your newly created or existing JSON quiz file from within the app and you won't need any fancy editors to get the job done. Just specify the filename or path and you'll be all set!")
	fmt.Printf("\nFile name/path: ")
	text, _ := reader.ReadString('\n')
	dataFile, readErr := os.ReadFile(strings.TrimSpace(text))
	if readErr != nil {
		color.Red("Something went wrong when reading the file! The specified file may not exist in the given directory or there might be something wrong with the path. \nerr: [%v]", readErr)
		return
	}
	if len(string(dataFile)) == 0  {
		color.Yellow("Empty file detected, adding starter JSON to file ...")
		writeErr := os.WriteFile(strings.TrimSpace(text), []byte(emptyJsonFileStartingText), 0644)
		if writeErr != nil {
			color.Red("Something went wrong when writing to the file! The specified file may not exist in the given directory or there might be something wrong with the path. \nerr: [%v]", writeErr)
			return
		}
		color.Green("\nStarter JSON added successfully!")
	} else {
		color.Yellow("[WARNING] Non-empty file detected! Files of this type may have incorrect formatting or invalid JSON. If you think the file is OK, then continue, otherwise it's recommended to wipe the file clean.")
		fmt.Printf("\n%v", color.YellowString("Wipe file? (y/n) "))
		wipeResponse, _ := reader.ReadString('\n')
		if strings.TrimSpace(wipeResponse) == "y" {
			writeErr := os.WriteFile(strings.TrimSpace(text), []byte(emptyJsonFileStartingText), 0644)
			if writeErr != nil {
				color.Red("Something went wrong when writing to the file! The specified file may not exist in the given directory or there might be something wrong with the path. \nerr: [%v]", writeErr)
				return
			}	
			color.Green("\nFile has been emptied successfully!")
		} else {
			color.Yellow("\n[WARNING] Proceeding with existing file ...")
		}
	}
	color.Green("\nInitializing quiz-writer ...")
	time.Sleep(time.Second * 2)
	for {
		resetTerminal()
		color.New(color.BgWhite, color.FgBlack).Println(" Quiz-Writer (New Question) ")
		var questionOptions []interface{}
		fmt.Print("\nQuestion: ")
		questionResponse, _ := reader.ReadString('\n')
		for optionCount := 0; optionCount < 4; optionCount++ {
			fmt.Printf("Option %v: ", (optionCount + 1))
			optionResponse, _ := reader.ReadString('\n')
			questionOptions = append(questionOptions, strings.TrimSpace(optionResponse))
		}
		fmt.Print("Correct Answer (option number): ")
		correctResponse, _ := reader.ReadString('\n')
		convInt, errConv := strconv.Atoi(strings.TrimSpace(correctResponse))
		if errConv != nil {
			color.Red("Something went wrong when parsing response! The \"correct\" key must have non-number value. err: [%v]", errConv)
			color.Yellow("[WARNING] Resetting quiz-writer plz wait ...")
			time.Sleep(time.Second * 2)
			continue
		}
		quizQ := QuizQuestion{Question: strings.TrimSpace(questionResponse), Options: questionOptions, Correct: (convInt - 1)}
		readFileJson, readFileErr := os.ReadFile(strings.TrimSpace(text))
		if readFileErr != nil {
			color.Red("Something went wrong when reading the file! The specified file may not exist in the given directory or there might be something wrong with the path. \nerr: [%v]", readFileErr)
			color.Yellow("[WARNING] Resetting quiz-writer plz wait ...")
			time.Sleep(time.Second * 2)
			continue
		}
		var quiz []QuizQuestion
		errUnm := json.Unmarshal([]byte(readFileJson), &quiz)
		quiz = append(quiz, quizQ)
		if errUnm != nil {
			color.Red("Something went wrong when setting the JSON values! \nerr: [%v]", errUnm)
			color.Yellow("[WARNING] Resetting quiz-writer plz wait ...")
			time.Sleep(time.Second * 2)
			continue
		}
		result, errMarshal := json.Marshal(quiz)
		if errMarshal != nil {
			color.Red("Something went wrong when setting the JSON values! \nerr: [%v]", errMarshal)
			color.Yellow("[WARNING] Resetting quiz-writer plz wait ...")
			time.Sleep(time.Second * 2)
			continue
		}
		errWriteJson := os.WriteFile(strings.TrimSpace(text), result, 0644)
		if errWriteJson != nil {
			color.Red("Something went wrong when writing the new JSON array to the file! \nerr: [%v]", errWriteJson)
			color.Yellow("[WARNING] Resetting quiz-writer plz wait ...")
			time.Sleep(time.Second * 2)
			continue
		}
		color.Green("Successfully added new question!")
		fmt.Printf("\n%v", color.YellowString("Add another question? (y/n) "))
		continueResponse, _ := reader.ReadString('\n')
		if strings.TrimSpace(continueResponse) == "y" {
			continue
		} else {
			break
		}
	}
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
	if len(strings.TrimSpace(data)) == 0 || !isJSON(strings.TrimSpace(data))  {
		color.Red("Something is wrong with the input data! Check if the data in the JSON file is valid.")
		return
	}
	m, ok := gjson.Parse(data).Value().([]interface{})
	if !ok {
		color.Red("Something went wrong! Data may not be in the correct format. err: [%v]", ok)
		return
	}
	correctResponses := 0
	incorrectResponses := 0
	var correctQuestions []string = make([]string, 0)
	var incorrectQuestions []string = make([]string, 0)
	var answers []string = make([]string, 0)
	for l := 0; l < len(m); l++ {
		question := (m[l].(map[string]interface{}))["Question"]
		options := (m[l].(map[string]interface{}))["Options"]
		correct := (m[l].(map[string]interface{}))["Correct"]
		if question == nil || options == nil || correct == nil {
			color.Red("The data in the JSON file or data from API is not in the correct format! You can take a look at the correct data format for quiz questions by the 5th option in the menu. err: [invalid_data_format]")
			return
		}
	}
	numberOfQuestions := len(m)
	if numberOfQuestions == 0 {
		color.Red("The JSON data does not contain question objects! Modify the JSON file or change the API, add some question objects and retry. err: [file_is_empty]")
		return
	}
	for i := 0; i < len(m); i++ {
		resetTerminal()
		question := (m[i].(map[string]interface{}))["Question"].(string)
		options := (m[i].(map[string]interface{}))["Options"].([]interface{})
		correct := (m[i].(map[string]interface{}))["Correct"].(float64)
		fmt.Printf("%v\n", color.New(color.BgWhite, color.FgBlack).Sprintf("Quiz initiated %v", color.New(color.FgBlack).Sprintf(" ( Q: %v of %v ) ", (i + 1), len(m))))
		quizQuestionTableCreator(question)
		fmt.Print("\n")
		quizOptionTableCreator(options)
		fmt.Print("\n\nAnswer: (1/2/3/4) ")
		text, _ := reader.ReadString('\n')
		if result, err := strconv.Atoi(strings.TrimSpace(text)); err != nil {
			color.Red("Input is not a number! Considering input as incorrect answer ...")
			incorrectQuestions = append(incorrectQuestions, question)
			answers = append(answers, options[int(correct)].(string))
			incorrectResponses++
		} else if (result - 1) == int(correct) {
			color.Green("\nCorrect Answer!")
			correctQuestions = append(correctQuestions, question)
			answers = append(answers, options[int(correct)].(string))
			correctResponses++
		} else {
			color.Red("\nIncorrect Answer!")
			incorrectQuestions = append(incorrectQuestions, question)
			answers = append(answers, options[int(correct)].(string))
			incorrectResponses++
		}
		if (i + 1) < len(m) {
			fmt.Print(color.YellowString("\nProceed to next question? (y/n) "))
			proceedPrompt, _ := reader.ReadString('\n')
			if strings.TrimSpace(proceedPrompt) == "y" {
				color.Green("Proceeding ...")
			} else {
				break
			}
		}
		time.Sleep(time.Second * 1)
	}
	color.Green("Rendering answer sheet ...")
	time.Sleep(time.Second * 2)
	resetTerminal()
	ClearTerminal()
	fmt.Printf("%v\n", color.New(color.BgWhite, color.FgBlack).Sprintf("   Quiz Answer Table (may be larger than console window)  "))
	finalQuestionAnswersTableCreator(correctQuestions, incorrectQuestions, answers)
	fmt.Printf("\n%v: %v of %v, %v: %v of %v\n", color.New(color.BgGreen, color.FgBlack).Sprintf("Correct"), correctResponses, len(m), color.New(color.BgRed, color.FgBlack).Sprintf("Incorrect"), incorrectResponses, len(m))
}

func quizOptionTableCreator(options []interface{}) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.AppendHeader(table.Row{"#", "   Options   "})
	for j := 0; j < len(options); j++ {
		t.AppendRow(table.Row{fmt.Sprintf("%v", j + 1), fmt.Sprintf("%v\n", options[j])})
	}
	t.Render()
}

func quizQuestionTableCreator(question string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendRow(table.Row{"   " + color.New(color.FgHiWhite).Sprintf("%v", question) + "   "})
	t.AppendSeparator()
	t.Render()
}

func finalQuestionAnswersTableCreator(correct []string, incorrect []string, answers []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"#", "Question", "Answer", "State"})
	rowNum := 1
	for g := 0; g < len(correct); g++ {
		t.AppendRow([]interface{}{rowNum, correct[g], answers[rowNum - 1], "✅"})
		t.AppendSeparator()
		rowNum++
	}
	for x := 0; x < len(incorrect); x++ {
		t.AppendRow([]interface{}{rowNum, incorrect[x], answers[rowNum - 1], "❌"})
		t.AppendSeparator()
		rowNum++
	}
	t.Render()
}

func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func resetTerminal() {
	screen.Clear()
	screen.MoveTopLeft()
}

func runCmd(name string, arg ...string) {
    cmd := exec.Command(name, arg...)
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func ClearTerminal() {
    switch runtime.GOOS {
    case "darwin":
        runCmd("clear")
    case "linux":
        runCmd("clear")
    case "windows":
        runCmd("cmd", "/c", "cls")
    default:
        runCmd("clear")
    }
}