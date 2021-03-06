package gobase

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const separator = ".##."

//EnvironmentPath ...
var EnvironmentPath = os.Getenv("GOPATH") + "/src/github.com/Harkishen-Singh/GO-file-database/"

func filter(text string) (string, bool) {

	if strings.Contains(text, separator) {
		return strings.Replace(text, separator, " ", -1), true
	}
	return "", false

}

func checkExistingDir(name string) bool {

	_, err := exec.Command("ls", name).Output()
	if err != nil {
		return false
	}
	return true
}

func makeDir(name string) bool {

	name = EnvironmentPath + "warehouse/" + name
	/**
	* returns the number of directories contained in the gien path [string]
	 */
	dirOccurence := func(path string) uint16 {
		var occ uint16
		for i := 0; i < len(path); i++ {
			if path[i] == '/' {
				occ++
			}
		}
		return occ
	}

	/**
	* creates new directories on name [string]
	 */
	createDirectoriesChain := func(name string) bool {
		if !checkExistingDir(name) {
			_, err := exec.Command("mkdir", "-p", name).Output()
			if err != nil {
				fmt.Println(err)
			}
			return true
		}

		return false
	}

	/**
	* controls the process of creation of chained directories on path [string]
	 */
	directoryController := func(path string) bool {
		tempPath := path + "/"
		createDirectoriesChain(tempPath)
		return true
	}

	if dirOccurence(name) == uint16(1) {
		if !checkExistingDir(name) {
			_, err := exec.Command("mkdir", name).Output()
			if err != nil {
				panic(err)
			}
		}
	} else {
		directoryController(name)
	}
	return true
}

func collectionStatus(collectionPath string) bool {

	warehouse() // check warehouse availability
	res := strings.LastIndex(collectionPath, "/")
	var subPath string
	if res != -1 {
		subPath = collectionPath[0:res]
		if !makeDir(subPath) {
			os.Exit(50005)
		}
	}
	expAddr := EnvironmentPath + "warehouse/" + subPath
	result, err := exec.Command("ls", expAddr).Output()
	if err != nil {
		return false
	}
	var existingFiles = strings.Split(string(result), "\n")
	var checkStatus bool

	for _, element := range existingFiles {
		if element == collectionPath[res+1:]+".data" {
			checkStatus = true
			break
		} else {
			checkStatus = false
		}
	}

	return checkStatus

}

func createCollection(address string) bool {

	address = EnvironmentPath + "warehouse/" + address + ".data"
	var ss = strings.LastIndex(address, "/")
	var pathx = address[:ss]
	exec.Command("mkdir", "-p", pathx).Output()
	_, err := exec.Command("touch", address).Output()

	if err != nil {
		fmt.Println("Error in createCollection Address: " + address)
		fmt.Println(err)
		return false
	}
	return true

}

//Retrive ...
func Retrive(address string) (string, string, bool) {

	var documentAvailable = collectionStatus(address)
	var data string
	var datatype string
	address = EnvironmentPath + "warehouse/" + address + ".data"
	if documentAvailable {
		openfile, err := ioutil.ReadFile(address)
		if err != nil {
			return "ERROR", "", false
		}
		data = string(openfile)
		datatype = data[:6]
		data = data[6:]
		return data, datatype, true
	}
	return "DOCUMENT_UNAVAILABLE", "", false

}

//CollectionsAvailable ...
func CollectionsAvailable(address string) ([]string, bool) {

	var existingCollections []string
	if address != "/" {
		path := EnvironmentPath + "warehouse/" + address
		response, err := exec.Command("ls", path).Output()
		if err != nil {
			fmt.Println("Error while looking for Collections, at Address: " + address)
			log.Fatal(err)
		}
		existingCollections = strings.Split(string(response), "\n")
	} else {
		response, err := exec.Command("ls", EnvironmentPath+"warehouse/").Output()
		if err != nil {
			fmt.Println("Error while looking for Collections, at Address: " + address)
			log.Fatal(err)
		}
		existingCollections = strings.Split(string(response), "\n")
	}
	return existingCollections, true

}

//Save ...
func Save(path string, data *string) bool {

	return saveCustom(path, *data, 12)

}

func saveCustom(path string, data string, pass uint16) bool {

	exists := collectionStatus(path)
	if exists == false {
		createCollection(path)
	}

	var address = EnvironmentPath + "warehouse/" + path + ".data"
	file, err := os.OpenFile(address, os.O_WRONLY, 0600)
	file.Seek(0, 0)
	file.Truncate(0)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data = CheckType(pass) + data

	_, err = file.WriteString(data)
	if err != nil {
		fmt.Println("Error occurred while writing the following data:")
		fmt.Println("\n" + data)
		fmt.Println("Address: warehouse/" + path)
		fmt.Println(err)
		return false
	}
	return true
}

//Delete ...
func Delete(path string) bool {

	path = EnvironmentPath + "warehouse/" + path + ".data"
	exec.Command("rm", path).Output()
	var x = strings.LastIndex(path, "/")

	cleanUPEmptyDirectories(path[:x])

	return true
}

func cleanUPEmptyDirectories(path string) {

	if path[len(path)-9:] != "warehouse" {
		fmt.Print("inside1")
		var tt = strings.LastIndex(path, "/")
		var xx = path[:tt]
		a, _ := exec.Command("ls", path).Output()
		var b = strings.Split(string(a), "\n")
		if xx[len(xx)-9:] != "warehouse" && len(b) != 0 {

			_, err := exec.Command("rm", "-R", xx).Output()
			if err != nil {
				fmt.Println(err)
			}
		}
		var temp = strings.LastIndex(path, "/")
		if path[len(path)-10:] != "warehouse" {
			cleanUPEmptyDirectories(path[:temp])
		}
	}
}

func warehouse() {

	db := "warehouse"
	resp, err := exec.Command("ls", EnvironmentPath).Output()
	if err != nil {
		panic(err)
	}
	var checkWarehouse bool
	var respStringArr = strings.Split(string(resp), "\n")
	for _, ele := range respStringArr {
		if ele == "warehouse" {
			checkWarehouse = true
			break
		}
	}
	if !checkWarehouse {
		exec.Command("mkdir", db).Output()
	}

}

//CheckType checks the type of data
func CheckType(pass uint16) string {
	var typeVar string

	switch pass {

	case 1:
		typeVar = "_uint8"

	case 2:
		typeVar = "__int8"

	case 3:
		typeVar = "uint16"

	case 4:
		typeVar = "_int16"

	case 5:
		typeVar = "uint32"

	case 6:
		typeVar = "_int32"

	case 7:
		typeVar = "uint64"

	case 8:
		typeVar = "_int64"

	case 9:
		typeVar = "___int"

	case 10:
		typeVar = "_flt32"

	case 11:
		typeVar = "_flt64"

	case 12:
		typeVar = "string"
	}

	return typeVar
}
