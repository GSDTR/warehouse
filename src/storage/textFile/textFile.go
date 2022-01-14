package TextFile

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Read_key_value_map(path string, delimiter string) map[string]string {
	resultedMap := make(map[string]string)

	file, err := os.Open(path)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		curString := scanner.Text()
		stringSplit := strings.Split(curString, delimiter)
//		fmt.Println(stringSplit)
		resultedMap[stringSplit[0]] = stringSplit[1]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return resultedMap
}