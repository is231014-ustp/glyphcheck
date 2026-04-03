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

var defaultSuspiciousCharacters = map[rune]struct{}{

	0x01C3: {}, //   LATIN LETTER RETROFLEX CLICK
	0xA778: {}, //	 LATIN SMALL LETTER UM
	0xA78C: {}, //   LATIN SMALL LETTER SALTILLO
	0x0181: {}, //	 LATIN CAPITAL LETTER B WITH HOOK
	0x018A: {}, //	 LATIN CAPITAL LETTER D WITH HOOK
	0x01A4: {}, //	 LATIN CAPITAL LETTER P WITH HOOK
	0x01AC: {}, //	 LATIN CAPITAL LETTER T WITH HOOK
	0x01B3: {}, //	 LATIN CAPITAL LETTER Y WITH HOOK
	0x0149: {}, //	 LATIN SMALL LETTER N PRECEDED BY APOSTROPHE
	0xFF2F: {}, //	 FULLWIDTH LATIN CAPITAL LETTER O
	0x0196: {}, //	 LATIN CAPITAL LETTER IOTA
	0x01C0: {}, //	 LATIN LETTER DENTAL CLICK
	0xFF29: {}, //	 FULLWIDTH LATIN CAPITAL LETTER I
	0xFF4C: {}, //	 FULLWIDTH LATIN SMALL LETTER L
	0x01A7: {}, //	 LATIN CAPITAL LETTER TONE TWO
	0x01B7: {}, //	 LATIN CAPITAL LETTER EZH
	0x021C: {}, //	 LATIN CAPITAL LETTER YOGH
	0xA76A: {}, //	 LATIN CAPITAL LETTER ET
	0xA7AB: {}, //	 LATIN CAPITAL LETTER REVERSED OPEN E
	0x01BC: {}, //	 LATIN CAPITAL LETTER TONE FIVE
	0x0222: {}, //	 LATIN CAPITAL LETTER OU
	0x0223: {}, //	 LATIN SMALL LETTER OU
	0xA76E: {}, //	 LATIN CAPITAL LETTER CON
	0x0241: {}, //	 LATIN CAPITAL LETTER GLOTTAL STOP
	0x0294: {}, //	 LATIN LETTER GLOTTAL STOP
	0xFF21: {}, //	 FULLWIDTH LATIN CAPITAL LETTER A
	0xA7B4: {}, //	 LATIN CAPITAL LETTER BETA
	0xFF22: {}, //	 FULLWIDTH LATIN CAPITAL LETTER B
	0xFF23: {}, //	 FULLWIDTH LATIN CAPITAL LETTER C
	0x0187: {}, //	 LATIN CAPITAL LETTER C WITH HOOK
	0x00C7: {}, //	 LATIN CAPITAL LETTER C WITH CEDILLA
	0x00D0: {}, //	 LATIN CAPITAL LETTER ETH
	0x0110: {}, //	 LATIN CAPITAL LETTER D WITH STROKE
	0x0189: {}, //	 LATIN CAPITAL LETTER AFRICAN D
	0xFF25: {}, //	 FULLWIDTH LATIN CAPITAL LETTER E
	0x0246: {}, //	 LATIN CAPITAL LETTER E WITH STROKE
	0xA798: {}, //	 LATIN CAPITAL LETTER F WITH STROKE
	0x0191: {}, //	 LATIN CAPITAL LETTER F WITH HOOK
	0x0193: {}, //	 LATIN CAPITAL LETTER G WITH HOOK
	0x01E4: {}, //	 LATIN CAPITAL LETTER G WITH STROKE
	0xFF28: {}, //	 FULLWIDTH LATIN CAPITAL LETTER H
	0x2C67: {}, //	 LATIN CAPITAL LETTER H WITH DESCENDER
	0x0197: {}, //	 LATIN CAPITAL LETTER I WITH STROKE
	0x019A: {}, //	 LATIN SMALL LETTER L WITH BAR
	0xA7B2: {}, //	 LATIN CAPITAL LETTER J WITH CROSSED-TAIL
	0xFF2A: {}, //	 FULLWIDTH LATIN CAPITAL LETTER J
	0x0248: {}, //	 LATIN CAPITAL LETTER J WITH STROKE
	0xFF2B: {}, //	 FULLWIDTH LATIN CAPITAL LETTER K
	0x0198: {}, //	 LATIN CAPITAL LETTER K WITH HOOK
	0x2C69: {}, //	 LATIN CAPITAL LETTER K WITH DESCENDER
	0xA740: {}, //	 LATIN CAPITAL LETTER K WITH STROKE
	0x0141: {}, //	 LATIN CAPITAL LETTER L WITH STROKE
	0xFF2D: {}, //	 FULLWIDTH LATIN CAPITAL LETTER M
	0xFF2E: {}, //	 FULLWIDTH LATIN CAPITAL LETTER N
	0x019D: {}, //	 LATIN CAPITAL LETTER N WITH LEFT HOOK
	0x01A0: {}, //	 LATIN CAPITAL LETTER O WITH HORN
	0x00D8: {}, //	 LATIN CAPITAL LETTER O WITH STROKE
	0x01FE: {}, //	 LATIN CAPITAL LETTER O WITH STROKE AND ACUTE
	0xFF30: {}, //	 FULLWIDTH LATIN CAPITAL LETTER P
	0x01A6: {}, //	 LATIN LETTER YR
	0xFF33: {}, //	 FULLWIDTH LATIN CAPITAL LETTER S
	0xFF34: {}, //	 FULLWIDTH LATIN CAPITAL LETTER T
	0x01AE: {}, //	 LATIN CAPITAL LETTER T WITH RETROFLEX HOOK
	0x0166: {}, //	 LATIN CAPITAL LETTER T WITH STROKE
	0x023E: {}, //	 LATIN CAPITAL LETTER T WITH DIAGONAL STROKE
	0x0244: {}, //	 LATIN CAPITAL LETTER U BAR
	0xA7B3: {}, //	 LATIN CAPITAL LETTER CHI
	0xFF38: {}, //	 FULLWIDTH LATIN CAPITAL LETTER X
	0xFF39: {}, //	 FULLWIDTH LATIN CAPITAL LETTER Y
	0x024E: {}, //	 LATIN CAPITAL LETTER Y WITH STROKE
	0xFF3A: {}, //	 FULLWIDTH LATIN CAPITAL LETTER Z
	0x0224: {}, //	 LATIN CAPITAL LETTER Z WITH HOOK
	0x01B5: {}, //	 LATIN CAPITAL LETTER Z WITH STROKE
	0x0251: {}, //	 LATIN SMALL LETTER ALPHA
	0xFF41: {}, //	 FULLWIDTH LATIN SMALL LETTER A
	0x0184: {}, //	 LATIN CAPITAL LETTER TONE SIX
	0x0253: {}, //	 LATIN SMALL LETTER B WITH HOOK
	0x0180: {}, //	 LATIN SMALL LETTER B WITH STROKE
	0x1D04: {}, //	 LATIN LETTER SMALL CAPITAL C
	0xFF43: {}, //	 FULLWIDTH LATIN SMALL LETTER C
	0x00E7: {}, //	 LATIN SMALL LETTER C WITH CEDILLA
	0x023C: {}, //	 LATIN SMALL LETTER C WITH STROKE
	0x0257: {}, //	 LATIN SMALL LETTER D WITH HOOK
	0xAB32: {}, //	 LATIN SMALL LETTER BLACKLETTER E
	0xFF45: {}, //	 FULLWIDTH LATIN SMALL LETTER E
	0x0247: {}, //	 LATIN SMALL LETTER E WITH STROKE
	0x0192: {}, //	 LATIN SMALL LETTER F WITH HOOK
	0x1E9D: {}, //	 LATIN SMALL LETTER LONG S WITH HIGH STROKE
	0xA799: {}, //	 LATIN SMALL LETTER F WITH STROKE
	0xAB35: {}, //	 LATIN SMALL LETTER LENIS F
	0x017F: {}, //	 LATIN SMALL LETTER LONG S
	0x1D6E: {}, //	 LATIN SMALL LETTER F WITH MIDDLE TILDE
	0x018D: {}, //	 LATIN SMALL LETTER TURNED DELTA
	0x0261: {}, //	 LATIN SMALL LETTER SCRIPT G
	0x1D83: {}, //	 LATIN SMALL LETTER G WITH PALATAL HOOK
	0xFF47: {}, //	 FULLWIDTH LATIN SMALL LETTER G
	0x0260: {}, //	 LATIN SMALL LETTER G WITH HOOK
	0x01E5: {}, //	 LATIN SMALL LETTER G WITH STROKE
	0xFF48: {}, //	 FULLWIDTH LATIN SMALL LETTER H
	0x0266: {}, //	 LATIN SMALL LETTER H WITH HOOK
	0x0127: {}, //	 LATIN SMALL LETTER H WITH STROKE
	0x0131: {}, //	 LATIN SMALL LETTER DOTLESS I
	0x0269: {}, //	 LATIN SMALL LETTER IOTA
	0x026A: {}, //	 LATIN LETTER SMALL CAPITAL I
	0xFF49: {}, //	 FULLWIDTH LATIN SMALL LETTER I
	0x1D7C: {}, //	 LATIN SMALL LETTER IOTA WITH STROKE
	0xFF4A: {}, //	 FULLWIDTH LATIN SMALL LETTER J
	0x0249: {}, //	 LATIN SMALL LETTER J WITH STROKE
	0x0199: {}, //	 LATIN SMALL LETTER K WITH HOOK
	0x026D: {}, //	 LATIN SMALL LETTER L WITH RETROFLEX HOOK
	0x026B: {}, //	 LATIN SMALL LETTER L WITH MIDDLE TILDE
	0x0142: {}, //	 LATIN SMALL LETTER L WITH STROKE
	0x0271: {}, //	 LATIN SMALL LETTER M WITH HOOK
	0x0273: {}, //	 LATIN SMALL LETTER N WITH RETROFLEX HOOK
	0x014B: {}, //	 LATIN SMALL LETTER ENG
	0x019E: {}, //	 LATIN SMALL LETTER N WITH LONG RIGHT LEG
	0x1D70: {}, //	 LATIN SMALL LETTER N WITH MIDDLE TILDE
	0x1D0F: {}, //	 LATIN LETTER SMALL CAPITAL O
	0x1D11: {}, //	 LATIN SMALL LETTER SIDEWAYS O
	0xAB3D: {}, //	 LATIN SMALL LETTER BLACKLETTER O
	0xFF4F: {}, //	 FULLWIDTH LATIN SMALL LETTER O
	0x01A1: {}, //	 LATIN SMALL LETTER O WITH HORN
	0x0275: {}, //	 LATIN SMALL LETTER BARRED O
	0xA74B: {}, //	 LATIN SMALL LETTER O WITH LONG STROKE OVERLAY
	0x00F8: {}, //	 LATIN SMALL LETTER O WITH STROKE
	0xAB3E: {}, //	 LATIN SMALL LETTER BLACKLETTER O WITH STROKE
	0x00FE: {}, //	 LATIN SMALL LETTER THORN
	0x01BF: {}, //	 LATIN LETTER WYNN
	0xFF50: {}, //	 FULLWIDTH LATIN SMALL LETTER P
	0x01A5: {}, //	 LATIN SMALL LETTER P WITH HOOK
	0x1D7D: {}, //	 LATIN SMALL LETTER P WITH STROKE
	0x02A0: {}, //	 LATIN SMALL LETTER Q WITH HOOK
	0xAB47: {}, //	 LATIN SMALL LETTER R WITHOUT HANDLE
	0xAB48: {}, //	 LATIN SMALL LETTER DOUBLE R
	0x027D: {}, //	 LATIN SMALL LETTER R WITH TAIL
	0x027C: {}, //	 LATIN SMALL LETTER R WITH LONG LEG
	0x1D72: {}, //	 LATIN SMALL LETTER R WITH MIDDLE TILDE
	0x024D: {}, //	 LATIN SMALL LETTER R WITH STROKE
	0x01BD: {}, //	 LATIN SMALL LETTER TONE FIVE
	0xA731: {}, //	 LATIN LETTER SMALL CAPITAL S
	0xFF53: {}, //	 FULLWIDTH LATIN SMALL LETTER S	
	0x1D74: {}, //	 LATIN SMALL LETTER S WITH MIDDLE TILDE
	0x01AD: {}, //	 LATIN SMALL LETTER T WITH HOOK
	0x1D75: {}, //	 LATIN SMALL LETTER T WITH MIDDLE TILDE
	0x0167: {}, //	 LATIN SMALL LETTER T WITH STROKE
	0x028B: {}, //	 LATIN SMALL LETTER V WITH HOOK
	0x1D1C: {}, //	 LATIN LETTER SMALL CAPITAL U
	0xA79F: {}, //	 LATIN SMALL LETTER VOLAPUK UE
	0xAB4E: {}, //	 LATIN SMALL LETTER U WITH SHORT RIGHT LEG
	0xAB52: {}, //	 LATIN SMALL LETTER U WITH LEFT HOOK
	0x1D20: {}, //	 LATIN LETTER SMALL CAPITAL V
	0xFF56: {}, //	 FULLWIDTH LATIN SMALL LETTER V
	0x026F: {}, //	 LATIN SMALL LETTER TURNED M
	0x1D21: {}, //	 LATIN LETTER SMALL CAPITAL W
	0xA761: {}, //	 LATIN SMALL LETTER VY
	0xFF58: {}, //	 FULLWIDTH LATIN SMALL LETTER X
	0x0263: {}, //	 LATIN SMALL LETTER GAMMA
	0x028F: {}, //	 LATIN LETTER SMALL CAPITAL Y
	0x1D8C: {}, //	 LATIN SMALL LETTER V WITH PALATAL HOOK
	0x1EFF: {}, //	 LATIN SMALL LETTER Y WITH LOOP
	0xAB5A: {}, //	 LATIN SMALL LETTER Y WITH SHORT RIGHT LEG
	0xFF59: {}, //	 FULLWIDTH LATIN SMALL LETTER Y
	0x01B4: {}, //	 LATIN SMALL LETTER Y WITH HOOK
	0x024F: {}, //	 LATIN SMALL LETTER Y WITH STROKE
	0x1D22: {}, //	 LATIN LETTER SMALL CAPITAL Z
	0x1D76: {}, //	 LATIN SMALL LETTER Z WITH MIDDLE TILDE
	0x01B6: {}, //	 LATIN SMALL LETTER Z WITH STROKE
	0x00DE: {}, //	 LATIN CAPITAL LETTER THORN
	0x1E9E: {}, //	 LATIN CAPITAL LETTER SHARP S
	0xA7B5: {}, //	 LATIN SMALL LETTER BETA
	0x0138: {}, //	 LATIN SMALL LETTER KRA
	0x1D0B: {}, //	 LATIN LETTER SMALL CAPITAL K
	0x0272: {}, //	 LATIN SMALL LETTER N WITH LEFT HOOK

}


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
			fmt.Printf("::error :: file=%s,line=%d,col=%d disallowed unicode character %U (%c)\n",fileName, lineIndex, runeIndex + 1, currentRune, currentRune);

		} else {

			var runeIsSuspicious bool;
			runeIsSuspicious = isSuspiciousCheck(currentRune, currentSuspiciousList);

			if runeIsSuspicious {

				fmt.Printf("::warning :: file=%s,line=%d,col=%d suspicious unicode character %U (%c)\n",fileName, lineIndex, runeIndex + 1, currentRune, currentRune);

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

			fmt.Printf("Unknown script in config: %s\n", scriptName);
			continue;

		}

		currentAllowList.Scripts = append(currentAllowList.Scripts, scriptRangeTable);

	}

	var categoryName string;
	
	for _, categoryName = range config.Allowed.Categories {

		var categoryRangeTable *unicode.RangeTable;
		categoryRangeTable = categoryNameToRangeTable(categoryName);

		if categoryRangeTable == nil {

			fmt.Printf("Unknown Category in config: %s\n", categoryName);
			continue;
		}
	
		currentAllowList.Categories = append(currentAllowList.Categories, categoryRangeTable);

	}


	return currentAllowList;

}

func buildSuspiciousList (config Config) suspiciousList {

	var currentSuspiciousList suspiciousList;
	currentSuspiciousList.Enabled = config.Suspicious.Enabled;
	
	var defaultConfigMerge map[rune]struct{};
	defaultConfigMerge = make(map[rune]struct{});

	var currentRune rune;

	for currentRune = range defaultSuspiciousCharacters {

		defaultConfigMerge[currentRune] = struct{}{};

	}


	var configMap map[rune]struct{};
	configMap = characterStringSliceToMap(config.Suspicious.Characters);


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
