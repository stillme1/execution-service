package run

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func RunCode(submissionId, problemId, input_dir string) {

	submission_dir := submissionId

	compiled := Compile(submission_dir+"/"+submissionId+".cpp", submission_dir+"/exec")
	if !compiled {
		println("Compilation Error")
		return
	} else {
		println("Compilation Success")
	}

	output_dir := submissionId + "/output"
	error_dir := submissionId + "/error"
	meta_dir := submissionId + "/meta"

	err := os.Mkdir(output_dir, 0777)
	if err != nil {
		return
	}

	err = os.Mkdir(error_dir, 0777)
	if err != nil {
		return
	}

	err = os.Mkdir(meta_dir, 0777)
	if err != nil {
		return
	}

	// copy the file to submission output_dir
	// Later going to be fetched from database
	timeLimit := 2.0
	memoryLimit := 100
	err = exec.Command("cp", "-r", input_dir, submission_dir).Run()
	if err != nil {
		fmt.Println("error copying input_dir:", err)
		return
	}

	// // time docker run -e SUBMISSION_ID="1" - TIME_LIMIT=2 -e MEMORY_LIMIT=100 -v /home/sahil/dockerTest/cpp:/1 my-image
	dir, _ := os.Getwd()
	err = exec.Command("docker", "run", "-e", "SUBMISSION_ID="+submissionId, "-e", "TIME_LIMIT="+fmt.Sprintf("%v", timeLimit), "-e", "MEMORY_LIMIT="+fmt.Sprintf("%d", memoryLimit), "-v", dir+"/"+submission_dir+":/"+submissionId, "my-image").Run()
	if err != nil {
		fmt.Println("error running docker:", err)
		return
	}

	// Judge the output
	// Judge(submissionId, problemId)

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
		println("test case", i, "passed")
	}
	println("exiting RunCode")
}
