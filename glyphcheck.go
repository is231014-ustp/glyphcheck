package main

import "fmt"
import "os"
import "log"
import "bufio"
import "io"

func main() {

var fileName string = "glyphcheck.go"
var file *os.File;
var errOpenFile error;

file, errOpenFile = os.Open(fileName);
if errOpenFile != nil {

	log.Fatal(errOpenFile);

}
defer file.Close();

var reader *bufio.Reader;
reader = bufio.NewReader(file);

var test []byte;
var testString string;
var testByteReadError error;

for {

	test, testByteReadError = reader.ReadBytes('\n');

	if testByteReadError != nil && testByteReadError != io.EOF {

		log.Fatal(testByteReadError);

	}

	if testByteReadError == io.EOF && len(test) == 0 {

		break

	}

	testString = string(test);

	fmt.Print(testString);

}

}
