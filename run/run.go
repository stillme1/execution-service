package run

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func RunCode(submissionId, problemId string) {

	submission_dir := submissionId + "/cpp"

	compiled := Compile(submission_dir+"/"+submissionId+".cpp", submission_dir+"/exec")
	if !compiled {
		println("Compilation Error")
		return
	} else {
		println("Compilation Success")
	}

	time_limit, memoryLimit, _, err := GetFiles(problemId, submissionId)
	timeLimit := float64(time_limit)
	if err != nil {
		println("Error getting files")
		return
	}
	// unzip submissionId.zip into submission_dir
	err = exec.Command("unzip", submissionId+".zip", "-d", submissionId).Run()
	if err != nil {
		fmt.Println("error unzipping:", err)
		return
	}
	meta_dir := submission_dir + "/meta"
	
	// Running docker contianer
	dir, _ := os.Getwd()
	err = exec.Command("docker", "run", "-e", "SUBMISSION_ID="+submissionId, "-e", "TIME_LIMIT="+fmt.Sprintf("%v", timeLimit), "-e", "MEMORY_LIMIT="+fmt.Sprintf("%d", memoryLimit), "-v", dir+"/"+submission_dir+":/"+submissionId, "my-image").Run()
	if err != nil {
		fmt.Println("error running docker:", err)
		return
	}

	files, err := os.Open(meta_dir)
	if err != nil {
		fmt.Println("error opening meta_dir:", err)
		return
	}
	defer files.Close()

	fileInfos, err := files.Readdir(-1)
	if err != nil {
		fmt.Println("error reading meta_dir:", err)
		return
	}
	for i, fileInfos := range fileInfos {
		metaFile, err := os.ReadFile(meta_dir + "/" + fileInfos.Name())
		if err != nil {
			fmt.Println("error reading meta_file:", err)
			break
		}
		metaInfo := strings.Split(string(metaFile), " ")
		// runtime as float32
		runTime, _ := strconv.ParseFloat(metaInfo[0], 32)
		// memory := metaInfo[1]
		exitCode := metaInfo[2]
		if exitCode != "0" {
			if (timeLimit - runTime) < 0.1 {
				println("timelimit exceeded")
			} else {
				println("runtime error")
			}
			break
		}
		// TODO Judge the output
		// .
		// .
		println("test case", i, "passed")
	}
}


func Compile(inpath, outpath string) bool {
	return exec.Command("g++", inpath, "-o", outpath).Run() == nil
}