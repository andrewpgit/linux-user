package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	userFile string = "/etc/passwd"
)

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	Name string `json:"name"`
	Pass string `json:"passwd"`
}

// Read json file and return slice of byte.
func ReadUsers(f string) []byte {

	jsonFile, err := os.Open(f)

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	data, _ := ioutil.ReadAll(jsonFile)
	return data
}

// Read file /etc/passwd and return slice of users
func ReadPasswd(f string) (list []string) {

	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := bufio.NewScanner(file)

	for r.Scan() {
		lines := r.Text()
		parts := strings.Split(lines, ":")
		list = append(list, parts[0])
	}
	return list
}

// Check if user on the host
func check(s []string, u string) bool {
	for _, w := range s {
		if u == w {
			return true
		}
	}
	return false
}

// Return securely generated random bytes

func CreateRandom(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
	return string(b)
}

// User is created by executing shell command useradd
func AddNewUser(u string, p string) bool {

	argUser := []string{"-d", "/tmp/henkel/home", "-s", "/sbin/nologin", u}
	argPass := []string{"-c", fmt.Sprintf("echo %s:%s | chpasswd", u, p)}
	fmt.Println("User add:", u, "=>", p)
	userCmd := exec.Command("useradd", argUser...)
	passCmd := exec.Command("/bin/sh", argPass...)
	if output, err := userCmd.Output(); err != nil {
		fmt.Println(err, "There was an error user add")
		return false
	} else {
		fmt.Println(string(output))

		if _, err := passCmd.Output(); err != nil {
			fmt.Println(err)
			return false
		}
		return true
	}
}

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Usage:", os.Args[0], "Name of json file")
		os.Exit(1)
	}
	NameOfFile := os.Args[1]

	d := ReadUsers(NameOfFile)

	var u Users
	json.Unmarshal(d, &u)

	userList := ReadPasswd(userFile)

	for i := range u.Users {

		c := check(userList, u.Users[i].Name)
		if c == false {
			encrypt := base64.StdEncoding.EncodeToString([]byte(u.Users[i].Pass + CreateRandom(5)))
			if info := AddNewUser(u.Users[i].Name, encrypt); info == true {
				fmt.Println("User was added:", u.Users[i].Name)
			}
		} else {
			fmt.Println("The user already exists:>", u.Users[i].Name)
		}
	}


}
