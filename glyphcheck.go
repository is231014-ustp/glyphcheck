//Test characters:
// а = Cyrillic 'a'; U+0430
// ρ = Greek 'roh'; U+03C1
// ɑ = Latin small alpha
// ᴀ = Latin small capital 'A'

package main

import "fmt"
import "os"
import "bufio"
import "io"
import "unicode/utf8"
import "unicode"
import "strings"
import "io/fs"
import "path/filepath"

import "golang.org/x/text/unicode/norm"
import "gopkg.in/yaml.v3"

type allowList struct {

	Scripts []*unicode.RangeTable
	Categories []*unicode.RangeTable
	Characters map[rune]struct{}

}

type suspiciousList struct {

	Enabled bool
	Characters map[rune]struct{}

}

var defaultSuspiciousCharacters = []rune{'ɑ'}


type Config struct {

	Root string `yaml:"root"`
	AllowedExtensions []string `yaml:"allowed_extensions"`
	ExcludedDirectories []string `yaml:"excluded_directories"`
	Allowed struct {
		Scripts []string `yaml:"scripts"`
		Categories []string `yaml:"categories"`	
		Characters []string `yaml:"characters"`
	} `yaml:"allowed"`
	Suspicious struct{
		Enabled bool `yaml:"enabled"`
		Characters []string `yaml:"characters"`

	} `yaml:"suspicious"`

}

func main() {

	var configName string = ".glyphcheck.yaml";

	var currentConfig Config;
	var currentConfigLoadError error;

	currentConfig, currentConfigLoadError = loadConfig(configName);

	if currentConfigLoadError != nil {

		fmt.Println(currentConfigLoadError);
		os.Exit(1);

	}

	var currentConfigInvalid error;

	currentConfigInvalid = validateConfig(currentConfig);

	if currentConfigInvalid != nil {

		fmt.Print(currentConfigInvalid);
		os.Exit(1);

	}

	var excludedDirectories map[string]struct{};
	excludedDirectories = directoryStringSliceToMap(currentConfig.ExcludedDirectories);
	
	var allowedExtensions map[string]struct{};
	allowedExtensions = extensionStringSliceToMap(currentConfig.AllowedExtensions);

	var currentAllowList allowList;
	currentAllowList = buildAllowList(currentConfig);

	var currentSuspiciousList suspiciousList;
	currentSuspiciousList = buildSuspiciousList(currentConfig);

	if scanDirectory(currentConfig.Root, currentAllowList, excludedDirectories, allowedExtensions, currentSuspiciousList) {
	
		os.Exit(1)

	}

}

func readFile(fileName string, currentAllowList allowList, currentSuspiciousList suspiciousList) bool {

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

		if checkLine(readBytes, lineIndex, fileName, currentAllowList, currentSuspiciousList) {

			violationInFile = true;

		}

		lineIndex++;

	}

	return violationInFile;

}

func checkLine(readBytes []byte, lineIndex int, fileName string, currentAllowList allowList, currentSuspiciousList suspiciousList) bool {

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
		

		var runeIsAllowed bool; 
		runeIsAllowed = isAllowed(currentRune, currentAllowList);

		if !runeIsAllowed {

			violationInLine = true;	
			fmt.Printf("[VIOLATION] File: %s; Line: %d; Column: %d; disallowed unicode character: %U (%c)\n",fileName, lineIndex, runeIndex + 1, currentRune, currentRune);

		} else {

			var runeIsSuspicious bool;
			runeIsSuspicious = isSuspiciousCheck(currentRune, currentSuspiciousList);

			if runeIsSuspicious {

				fmt.Printf("[WARNING] File: %s; Line: %d; Column: %d; suspicious unicode character: %U (%c)\n",fileName, lineIndex, runeIndex + 1, currentRune, currentRune);

			}
		}
	
		runeIndex++;
		i += runeSize;	

	}

	return violationInLine;

}

func isAllowed(runeToCheck rune, currentAllowList allowList) bool {

	var runeExists bool;
	_, runeExists = currentAllowList.Characters[runeToCheck];

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

func isSuspiciousCheck(runeToCheck rune, currentSuspiciousList suspiciousList) bool {

	var runeExists bool;

	if !currentSuspiciousList.Enabled {

		return false;

	}

	_, runeExists = currentSuspiciousList.Characters[runeToCheck];

	return runeExists;
}

func scanDirectory (root string, currentAllowList allowList, excludedDirectories map[string]struct{}, allowedExtensions map[string]struct{}, currentSuspiciousList suspiciousList) bool {

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

		if readFile(path, currentAllowList, currentSuspiciousList) {

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

func loadConfig(configPath string) (Config, error) {

	var currentConfig Config;

	currentConfig.Root = ".";
	currentConfig.AllowedExtensions = []string{".go"};
	currentConfig.ExcludedDirectories = []string{".git"};

	currentConfig.Allowed.Scripts = []string{"Latin", "Common"};
	currentConfig.Allowed.Categories = []string{"Letter", "Number", "Punct"};
	currentConfig.Allowed.Characters = []string{"\n","\r","\t", " ", "=", "+", "-", "<", ">"};

	currentConfig.Suspicious.Enabled = true;

	var configFileData []byte;
	var configFileError error;

	configFileData, configFileError = os.ReadFile(configPath);

	if configFileError != nil {

		return currentConfig, configFileError;

	}

	var yamlConfigFileError error;

	yamlConfigFileError = yaml.Unmarshal(configFileData, &currentConfig);

	if yamlConfigFileError != nil {

		return currentConfig, yamlConfigFileError;

	}


	return currentConfig, nil;

}

func validateConfig(configToCheck Config) error {

	if strings.TrimSpace(configToCheck.Root) == "" {

		return fmt.Errorf("root can not bee empty\n");

	}

	if len(configToCheck.AllowedExtensions) == 0 {

		return fmt.Errorf("extensions must contain at least one entry\n");

	}

	if len(configToCheck.Allowed.Scripts) == 0 {

		return fmt.Errorf("allowed.scripts can not be empty\n");

	}

	if len(configToCheck.Allowed.Categories) == 0 {

		return fmt.Errorf("allowed.categories can not be emtpy\n");

	}

	return nil;

}

func directoryStringSliceToMap (directories []string) (map[string]struct{}) {

	var directoryMap map[string]struct{};
	directoryMap = make(map[string]struct{});

	var directoryString string;

	for _, directoryString = range directories {

		directoryString = strings.TrimSpace(directoryString);

		if directoryString == "" {

			continue;

		}

		directoryMap[directoryString] = struct {}{};

	}

	return directoryMap;

}

func extensionStringSliceToMap (extensions []string) (map[string]struct{}) {

	var extensionMap map[string]struct{};
	extensionMap = make(map[string]struct{});

	var extensionString string;

	for _, extensionString = range extensions {

		extensionString = normalizeExtension(extensionString);

		if extensionString == "" {

			continue;

		}

		extensionMap[extensionString] = struct {}{};

	}

	return extensionMap;

}

func normalizeExtension (extension string) string {

	extension = strings.ToLower(strings.TrimSpace(extension));

	if extension == "" {

		return "";

	}

	if !strings.HasPrefix(extension, ".") {

		extension = "." + extension

	}

	return extension;

}

func buildAllowList (config Config) allowList {

	var currentAllowList allowList;
	
	currentAllowList.Characters = characterStringSliceToMap(config.Allowed.Characters);

	var scriptName string;

	for _, scriptName = range config.Allowed.Scripts {

		var scriptRangeTable *unicode.RangeTable;
		scriptRangeTable = scriptNameToRangeTable(scriptName);

		if scriptRangeTable == nil {

			fmt.Printf("Unkown script in config: %s\n", scriptName);
			continue;

		}

		currentAllowList.Scripts = append(currentAllowList.Scripts, scriptRangeTable);

	}

	var categoryName string;
	
	for _, categoryName = range config.Allowed.Categories {

		var categoryRangeTable *unicode.RangeTable;
		categoryRangeTable = categoryNameToRangeTable(categoryName);

		if categoryRangeTable == nil {

			fmt.Printf("Unkown Category in config: %s\n", categoryName);
			continue;
		}
	
		currentAllowList.Categories = append(currentAllowList.Categories, categoryRangeTable);

	}


	return currentAllowList;

}

func runeSliceToMap(runeSlice []rune) map[rune]struct{} {

	var runeMap map[rune]struct{};
	runeMap = make(map[rune]struct{});

	var currentRune rune;
	
	for _, currentRune = range runeSlice {

		runeMap[currentRune] = struct{}{};

	}

	return runeMap;

}

func buildSuspiciousList (config Config) suspiciousList {

	var currentSuspiciousList suspiciousList;
	currentSuspiciousList.Enabled = config.Suspicious.Enabled;
	
	var defaultConfigMerge map[rune]struct{};
	defaultConfigMerge = runeSliceToMap(defaultSuspiciousCharacters);

	var configMap map[rune]struct{};
	configMap = characterStringSliceToMap(config.Suspicious.Characters);

	var currentRune rune;

	for currentRune = range configMap {

		defaultConfigMerge[currentRune] = struct{}{};

	}

	currentSuspiciousList.Characters = defaultConfigMerge;

	return currentSuspiciousList;

}

func characterStringSliceToMap (characters []string) (map[rune]struct{}) {

	var characterMap map[rune]struct{};
	characterMap = make(map[rune]struct{});

	var characterString string;

	for _, characterString = range characters {

		var characterRune rune;

		for _, characterRune = range characterString {

			characterMap[characterRune] = struct{}{};

		}

	}

	return characterMap;

}

func categoryNameToRangeTable(categoryName string) *unicode.RangeTable {

	categoryName = strings.TrimSpace(categoryName);

	switch categoryName {

	case "Letter":
		return unicode.Letter;
	case "Number":
		return unicode.Number;
	case "Punct":
		return unicode.Punct;

	}

	return nil
}

func scriptNameToRangeTable(scriptName string) *unicode.RangeTable {

	scriptName = strings.TrimSpace(scriptName);	

	switch scriptName {

	case "Latin":
		return unicode.Latin;
	case "Common":
		return unicode.Common;

	}

	return nil
}
