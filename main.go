package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/thoas/go-funk"
)

const ATCODER_API_SUBMISSION_URL = "https://kenkoooo.com/atcoder/atcoder-api/results?user=togatoga"

type AtCoderSubmission struct {
	ID            int     `json:"id"`
	EpochSecond   int     `json:"epoch_second"`
	ProblemID     string  `json:"problem_id"`
	ContestID     string  `json:"contest_id"`
	UserID        string  `json:"user_id"`
	Language      string  `json:"language"`
	Point         float64 `json:"point"`
	Length        int     `json:"length"`
	Result        string  `json:"result"`
	ExecutionTime int     `json:"execution_time"`
}

func isDirExist(path string) bool {
	info, _ := os.Stat(path)

	return info.IsDir()
}
func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

type Config struct {
	UserID         string `json:"user_id"`
	RepositoryPath string `json:"repository_path"`
	ServiceName    string `json:"service_name"`
}

func init() {

	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	configDir := filepath.Join(home, ".pgr")
	if !isDirExist(configDir) {
		err = os.MkdirAll(configDir, 0700)
		if err != nil {
			panic(err)
		}
	}

	configFile := filepath.Join(configDir, "config.json")
	if !isFileExist(configFile) {
		//initial config
		config := []Config{Config{"", "", "atcoder"}}
		jsonBytes, err := json.Marshal(config)
		if err != nil {
			panic(err)
		}
		json := string(jsonBytes)
		file, err := os.Create(filepath.Join(configDir, "config.json"))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		file.WriteString(json)
	}
}

func main() {
	resp, err := http.Get(ATCODER_API_SUBMISSION_URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var ss []AtCoderSubmission
	err = json.Unmarshal(bytes, &ss)
	if err != nil {
		panic(err)
	}

	//only ac
	ss = funk.Filter(ss, func(s AtCoderSubmission) bool {
		return s.Result == "AC"
	}).([]AtCoderSubmission)

	//rev sort by EpochSecond
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].EpochSecond > ss[j].EpochSecond
	})

	//filter latest submission for each problem
	v := map[string]struct{}{}
	ss = funk.Filter(ss, func(s AtCoderSubmission) bool {
		_, ok := v[s.ContestID+"_"+s.ProblemID]
		if ok {
			return false
		}
		v[s.ContestID+"_"+s.ProblemID] = struct{}{}
		return true
	}).([]AtCoderSubmission)
	fmt.Println(len(ss))
}