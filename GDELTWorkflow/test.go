package main

import (
	"bufio"
	//"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"io/ioutil"
	"log"
	//"time"
)

func pythonCall(progName string, inChannel chan <- string, workflowNumber string) {
	cmd := exec.Command("python3", progName, workflowNumber)
	out, err := cmd.CombinedOutput()
	log.Println(cmd.Run())

	if err != nil {
		fmt.Println(err)
		// Exit with status 3.
    os.Exit(3)
	}
	fmt.Println(string(out))
	//check if msg is legit
	msg := string(out)[:len(out)-1]
	//msg := ("Module Completed: " + progName)
	inChannel <- msg
}


func integratePythonCall(progName string, inChannel1 chan <- string, inChannel2 chan <- string, workflowNumber string) {
	cmd := exec.Command("python3", progName, workflowNumber)
	out, err := cmd.CombinedOutput()
	log.Println(cmd.Run())

	if err != nil {
		fmt.Println(err)
		// Exit with status 3.
    os.Exit(3)
	}
	fmt.Println(string(out))
	//check if msg is legit
	msg := string(out)[:len(out)-1]
	//msg := ("Module Completed: " + progName)
	inChannel1 <- msg
	inChannel2 <- msg
}


func simplePythonCall(progName string){
	cmd := exec.Command("python3", progName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os. Stderr
	log.Println(cmd.Run())
}

func messagePassing(inChannel <- chan string, outChannel chan <- string ){
	msg := <- inChannel
	outChannel <- msg
}
func integrateMessagePassing(inChannel1 <- chan string, inChannel2 <- chan string, outChannel chan <- string ){
	msg1 := <- inChannel1
	msg2 := <- inChannel2
	outChannel <- msg1 + msg2
}

func numOfFiles(folder string) int{
    files,_ := ioutil.ReadDir(folder)
    return len(files)
}

//reads a file and returns an array of comments beginning with ##
func readLines( progName string) [20]string{
		var commandsArray [20]string
    file, err := os.Open(progName)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    i := 0
    for scanner.Scan() {
        command := scanner.Text()
				//fmt.Println(len(command))
				if len(command) >1{
					//fmt.Println("dh")
					if command[0:2] == "##" {
						//fmt.Println(command[2:len(command)])
						commandsArray[i] = command[2:len(command)]
						i++
					}
				}
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

		return commandsArray
}

func main(){
	for i := 1; i<=2; i++{
		//check if input location is available
		fmt.Println((i))
		//c1 := "python -c from workflow import userScript; print userScript.inputDataset" + strconv.Itoa(i)
		cmd := exec.Command("python", "-c", "from workflow import userScript; print userScript.inputDataset" + strconv.Itoa(i))
		//cmd := exec.Command(c1)
		out, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Println(err)
			// Exit with status 3.
	    os.Exit(3)
		} else if out == nil{
			os.Exit(3)
		} else {
			//input dataset from disk
			//check if empty
			inputDataset := string(out)[:len(out)-1]
			fmt.Print(inputDataset+"\n")
		}


		//check if output location is available
		cmd1 := exec.Command("python", "-c", "from workflow import userScript; print userScript.outputLocation" + strconv.Itoa(i))
		out1, err1 := cmd1.CombinedOutput()

		if err1 != nil {
			fmt.Println(err1)
			// Exit with status 3.
			os.Exit(3)
		} else if out1 == nil{
			os.Exit(3)
		} else {
			//output dataset from disk
			//check if empty
			outputDataset := string(out1)[:len(out1)-1]
			fmt.Print(outputDataset+"\n")
		}
	}

	commandsArray := readLines("workflow/userScript.py")
	fmt.Println(commandsArray)

	//configurations
	simplePythonCall("workflow/parslConfig.py")


	//start module execution from here onwards
	inChannelModule0 := make(chan string, 1)
	outChannelModule0 := make(chan string, 1)
	go pythonCall("workflow/"+commandsArray[0], inChannelModule0,"1")
	//pythonCall("workflow/gdeltFileSelection/dataFilesIntegration.py", inChannelModule1)
	go messagePassing(inChannelModule0, outChannelModule0)
	fmt.Println(<-outChannelModule0)

	outChannelModule1 := make(chan string, 1)
	go pythonCall("workflow/"+commandsArray[1], outChannelModule0,"1")
	//pythonCall("workflow/gdeltFileSelection/countrySelection.py", inChannelModule1)
	go messagePassing(outChannelModule0, outChannelModule1)
	fmt.Println(<-outChannelModule1)

	outChannelModule2 := make(chan string, 1)
	go pythonCall("workflow/"+commandsArray[2], outChannelModule1, "1")
	//pythonCall("workflow/selection/selectUserDefinedColumns.py", outChannelModule1)
	go messagePassing(outChannelModule1, outChannelModule2)
	fmt.Println(<- outChannelModule2)

	outChannelModule3 := make(chan string, 1)
	go pythonCall("workflow/"+commandsArray[3], outChannelModule2, "1")
	//pythonCall("workflow/cleaning/dropUniqueColumns.py", outChannelModule2)
	go messagePassing(outChannelModule2, outChannelModule3)
	fmt.Println(<- outChannelModule3)

	outChannelModule4 := make(chan string, 1)
	go pythonCall("workflow/"+commandsArray[4], outChannelModule3, "1")
	//pythonCall("workflow/cleaning/dropColumnsCriteria.py", outChannelModule2)
	go messagePassing(outChannelModule3, outChannelModule4)
	fmt.Println(<- outChannelModule3)

	outChannelModule5 := make(chan string, 1)
	//pythonCall("workflow/cleaning/dropRowsCriteria.py", outChannelModule3)
	go pythonCall("workflow/"+commandsArray[5], outChannelModule4, "1")
	go messagePassing(outChannelModule4, outChannelModule5)
	fmt.Println(<- outChannelModule5)

	outChannelModule6 := make(chan string, 1)
	go pythonCall("workflow/"+commandsArray[6], outChannelModule5, "1")
	//pythonCall("workflow/cleaning/removeDuplicateRows.py", outChannelModule4)
	go messagePassing(outChannelModule5, outChannelModule6)
	fmt.Println(<- outChannelModule6)

	outChannelModule7 := make(chan string, 1)
	//pythonCall("workflow/cleaning/missingValuesMode.py", outChannelModule5)
	go pythonCall("workflow/"+commandsArray[7], outChannelModule6, "1")
	go messagePassing(outChannelModule6, outChannelModule7)
	fmt.Println(<- outChannelModule7)

	outChannelModule8 := make(chan string, 1)
	//pythonCall("workflow/transformation/combineColumns.py", outChannelModule8)
	go pythonCall("workflow/"+commandsArray[9], outChannelModule7, "1")
	go messagePassing(outChannelModule7, outChannelModule8)
	fmt.Println(<- outChannelModule8)


	inChannelModule21 := make(chan string, 1)
	outChannelModule21 := make(chan string, 1)
	pythonCall("workflow/"+commandsArray[2], inChannelModule21,"2")
	//pythonCall("workflow/selection/selectUserDefinedColumns.py", inChannelModule1)
	messagePassing(inChannelModule21, outChannelModule21)
	fmt.Println(<-outChannelModule21)

	outChannelModule22 := make(chan string, 1)
	pythonCall("workflow/"+commandsArray[3], outChannelModule21,"2")
	//pythonCall("workflow/selection/dropUniqueColumns.py", inChannelModule1)
	messagePassing(outChannelModule21, outChannelModule22)
	fmt.Println(<-outChannelModule22)

	outChannelModule23 := make(chan string, 1)
	pythonCall("workflow/"+commandsArray[6], outChannelModule22, "2")
	//pythonCall("workflow/cleaning/removeDuplicateRows.py", outChannelModule4)
	messagePassing(outChannelModule22, outChannelModule23)
	fmt.Println(<- outChannelModule23)

	outChannelModule24 := make(chan string, 1)
	//pythonCall("workflow/cleaning/missingValuesMode.py", outChannelModule5)
	pythonCall("workflow/"+commandsArray[7], outChannelModule23, "2")
	messagePassing(outChannelModule23, outChannelModule24)
	fmt.Println(<- outChannelModule24)

	outChannelModule25 := make(chan string, 1)
	//pythonCall("workflow/integrateLabels/addLabelColumn.py", outChannelModule5)
	pythonCall("workflow/"+commandsArray[10], outChannelModule24, "2")
	messagePassing(outChannelModule24, outChannelModule25)
	fmt.Println(<- outChannelModule25)

	outChannelModule26 := make(chan string, 1)
	//pythonCall("workflow/integrateLabels/assignCountryCode.py", outChannelModule5)
	pythonCall("workflow/"+commandsArray[11], outChannelModule25, "2")
	messagePassing(outChannelModule25, outChannelModule26)
	fmt.Println(<- outChannelModule26)

	outChannelModule27 := make(chan string, 1)
	//pythonCall("workflow/integrateLabels/splitDate.py", outChannelModule5)
	pythonCall("workflow/"+commandsArray[12], outChannelModule26, "2")
	messagePassing(outChannelModule26, outChannelModule27)
	fmt.Println(<- outChannelModule27)

	outChannelModule9 := make(chan string, 1)
	//pythonCall("workflow/integrateLabels/integrate.py", outChannelModule5)
	integratePythonCall("workflow/"+commandsArray[14], outChannelModule27, outChannelModule8, "1")
	integrateMessagePassing(outChannelModule27, outChannelModule8, outChannelModule9)
	fmt.Println(<- outChannelModule9)

	outChannelModule10 := make(chan string, 1)
	pythonCall("workflow/"+commandsArray[8], outChannelModule9, "1")
	//pythonCall("workflow/transformation/normalize.py", outChannelModule6)
	messagePassing(outChannelModule9, outChannelModule10)
	fmt.Println(<- outChannelModule10)

	outChannelModule11 := make(chan string, 1)
	//pythonCall("workflow/mining/randomForestClassification.py", outChannelModule5)
	pythonCall("workflow/"+commandsArray[15], outChannelModule10, "1")
	messagePassing(outChannelModule10, outChannelModule11)
	fmt.Println(<- outChannelModule11)
/*
	inChannelModule31
	outChannelModule31 := make(chan string, 1)
	//pythonCall("workflow/mining/randomForestClassification.py", outChannelModule5)
	pythonCall("workflow/"+commandsArray[15], outChannelModule10, "1")
	messagePassing(outChannelModule10, outChannelModule11)
	fmt.Println(<- outChannelModule11)
*/

}



/*

NEED TO CONNECT WITH RF WF
	outChannelModule29 := make(chan string, 1)
	//pythonCall("workflow/integrateLabels/integrate.py", outChannelModule5)
	pythonCall("workflow/"+commandsArray[5], outChannelModule28, "2")
	messagePassing(outChannelModule28, outChannelModule29)
	fmt.Println(<- outChannelModule29)

*/
