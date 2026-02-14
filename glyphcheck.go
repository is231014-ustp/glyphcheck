package main

import "fmt"
import "os"
import "log"
import "bufio"
import "io"
import "golang.org/x/text/unicode/norm"
import "unicode/utf8"


func main() {

	var fileName string;
	fileName = "glyphcheck.go"

	readFile(fileName);


}

func readFile(fileName string) {

	var file *os.File;
	var errOpenFile error;


	file, errOpenFile = os.Open(fileName);
	if errOpenFile != nil {

		log.Fatal(errOpenFile);

	}
	defer file.Close();

	var reader *bufio.Reader;
	reader = bufio.NewReader(file);

	var readBytes []byte;
	var readBytesError error;

	var lineIndex int = 1;

	for {

		readBytes, readBytesError = reader.ReadBytes('\n');

		if readBytesError != nil && readBytesError != io.EOF {

			log.Fatal(readBytesError);

		}

		if readBytesError == io.EOF && len(readBytes) == 0 {

			break;

		}

		checkLine(readBytes, lineIndex);
		lineIndex ++;

	}

}

func checkLine(readBytes []byte, lineIndex int) {

	var normalizedBytes []byte;	
	normalizedBytes = norm.NFC.Bytes(readBytes);

	var runeIndex int = 0;

	for i := 0; i < len(normalizedBytes); {

		var currentRune rune;
		var runeSize int;
		
		currentRune, runeSize = utf8.DecodeRune(normalizedBytes[i:])			
		if currentRune == utf8.RuneError && runeSize == 1 {

			fmt.Printf("Invalid UTF-8 encoding 0x%X at offset %d\n", normalizedBytes[i], i);
			runeIndex ++;
			i += runeSize;
			continue;

		}

		fmt.Printf("Line: %d; Column: %d; Character: %c; Unicode: %U\n", lineIndex, runeIndex + 1 , currentRune, currentRune);
		
		runeIndex ++;
		i += runeSize;	
		

	}


}
