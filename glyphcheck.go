//Test characters:
// а = Cyrillic 'a'; U+0430
// ρ = Greek 'roh'; U+03C1

package main

import "fmt"
import "os"
import "bufio"
import "io"
import "golang.org/x/text/unicode/norm"
import "unicode/utf8"
import "unicode"
import "strings"
import "io/fs"
import "path/filepath"

type allowList struct {

	Scripts []*unicode.RangeTable
	Categories []*unicode.RangeTable
	Characters map[rune]struct{}

}

var excludedDirectories = map[string]struct{}{

	".git":{},

}

var allowedExtensions = map[string]struct{}{

	".go":{},

}

var devAllowList = allowList{

	Scripts: []*unicode.RangeTable{
		
		unicode.Latin,
		unicode.Common,

	},

	Categories: []*unicode.RangeTable{

		unicode.Letter,
		unicode.Number,
		unicode.Punct,

	},

	Characters: map[rune]struct{}{

		'\n':{},
		'\r':{},
		'\t':{},
		' ':{},
		'=':{},
		'+':{},
		'<':{},
		'>':{},

	},

}

func main() {

	var root string = ".";

	if scanDirectory(root) {
	
		os.Exit(1)

	}

}

func readFile(fileName string, currentAllowList allowList) bool {

	var file *os.File;
	var errOpenFile error;


	file, errOpenFile = os.Open(fileName);
	if errOpenFile != nil {

		fmt.Println(errOpenFile);
		return true;

	}
	defer file.Close();

	var reader *bufio.Reader;
	reader = bufio.NewReader(file);

	var readBytes []byte;
	var readBytesError error;

	var violationInFile bool = false;

	var lineIndex int = 1;

	for {

		readBytes, readBytesError = reader.ReadBytes('\n');

		if readBytesError != nil && readBytesError != io.EOF {

			fmt.Println(readBytesError);
			return true;

		}

		if readBytesError == io.EOF && len(readBytes) == 0 {

			break;

		}

		if checkLine(readBytes, lineIndex, fileName, currentAllowList) {

			violationInFile = true;

		}

		lineIndex++;

	}

	return violationInFile;

}

func checkLine(readBytes []byte, lineIndex int, fileName string, currentAllowList allowList) bool {

	var normalizedBytes []byte;	
	normalizedBytes = norm.NFC.Bytes(readBytes);

	var runeIndex int = 0;

	var violationInLine bool = false;

	for i := 0; i < len(normalizedBytes); {

		var currentRune rune;
		var runeSize int;
		
		currentRune, runeSize = utf8.DecodeRune(normalizedBytes[i:])			
		if currentRune == utf8.RuneError && runeSize == 1 {

			fmt.Printf("Invalid UTF-8 encoding 0x%X at offset %d\n", normalizedBytes[i], i);
			violationInLine = true;
			runeIndex++;
			i += runeSize;

			continue;

		}
		

		var runeIsAllowed bool = isAllowed(currentRune, currentAllowList);
		
		if !runeIsAllowed {

			violationInLine = true;	
			fmt.Printf("File: %s; Line: %d; Column: %d; Byte Offset: %d; Character: %c; Unicode: %U\n",fileName, lineIndex, runeIndex + 1, i, currentRune, currentRune);

		}
	
		runeIndex++;
		i += runeSize;	
		

	}

	return violationInLine;

}

func getAllowList() allowList {

	return devAllowList;

}

func isAllowed(runeToCheck rune, currentAllowList allowList) bool {

	var runeExists bool;
	_, runeExists = currentAllowList.Characters[runeToCheck]

	if runeExists {

		return true;

	}

	if !unicode.In(runeToCheck, currentAllowList.Scripts...) {

		return false;

	}

	if !unicode.In(runeToCheck, currentAllowList.Categories...) {

		return false;

	}

	return true;

}

func scanDirectory (root string) bool {

	var currentAllowList allowList;
	currentAllowList = getAllowList();


	var violationInDirectory bool = false;

	var walkDirError error;

	walkDirError = filepath.WalkDir(root, func(path string, currentDirectory fs.DirEntry, err error) error{

		if err != nil {

			fmt.Printf("Error while scanning Directories: %s: %v\n", path, err);
			violationInDirectory = true;
			return nil;

		}

		var name string;
		name = currentDirectory.Name();
		
		if currentDirectory.IsDir() {

			var dirIsExcluded bool;
			_, dirIsExcluded = excludedDirectories[name];

			if dirIsExcluded {

				return filepath.SkipDir;

			}
			
			return nil;

		}

		if !currentDirectory.Type().IsRegular() {

			return nil;

		}


		var extension string;
		extension = strings.ToLower(filepath.Ext(name))

		var extensionIsAllowed bool;
		_, extensionIsAllowed = allowedExtensions[extension];

		if !extensionIsAllowed {

			return nil;

		}

		if readFile(path, currentAllowList) {

			violationInDirectory = true;

		}

		return nil;


	})

	if walkDirError != nil {

		fmt.Printf("WalkDir failed: %v\n", walkDirError);
		violationInDirectory = true;

	}

	return violationInDirectory


}
