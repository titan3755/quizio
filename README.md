# QuizIO

QuizIO is a feature-rich terminal quiz app made with the Go programming language! You can use the program to easily take part in a quiz with data provided from a JSON file or API. You can create and edit the JSON file from a user-friendly interface from within the program.

## Usage

QuizIO is easier to use than you think! Most of the features you need will be found within the app. The app has five "modes" in a menu using which you can achieve different tasks.

### App Modes:
* `JSON file read`
* `JSON file write`
* `API data`
* `Quiz writer`
* `Preview data format`

#### `JSON file read mode:`
This mode can be used to read the quiz data from a JSON file by providing the relative or absolute path of the file in the input prompt. The quiz will be initiated automatically after the file data has been read.
#### `JSON file write mode:`
This mode can be used to create a quiz data file by providing the name of the JSON file in the input prompt. An option to edit the newly created file will come up.
#### `API data mode:`
This mode can be used to get the quiz data from an API by providing the URL in the input prompt. The API must provide the quiz data in the correct format. The quiz will be initiated automatically after the file data has been read.
#### `Quiz writer mode:`
This mode can be used to modify a newly created or existing JSON file from within the app and you will not need to manually modify a JSON file from another editor. An option to wipe the provided JSON file will come as existing JSON files may not contain data in the required format. You can also proceed with the existing JSON file.
#### `Preview data format mode:`
This is not exactly a "mode" but a display of the required data format that is mandatory in the input JSON files or API data. If the JSON data is not in the format showed in this mode, the quiz will not be initiated and an error will occur.
## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. Report encountered bugs as soon as possible to get a quick fix.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
