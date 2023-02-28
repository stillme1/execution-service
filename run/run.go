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
		println("Error getting meta_files")
		return
	}
	// unzip submissionId.zip into submission_dir
	err = exec.Command("unzip", submissionId+".zip", "-d", submissionId).Run()
	if err != nil {
		println("error unzipping:", err)
		return
	}
	meta_dir := submission_dir + "/meta"
	checker_dir := submission_dir + "/checker"
	verdict_dir := submission_dir + "/verdict"
	
	// Running docker contianer
	dir, _ := os.Getwd()
	err = exec.Command("docker", "run", "-e", "SUBMISSION_ID="+submissionId, "-e", "TIME_LIMIT="+fmt.Sprintf("%v", timeLimit), "-e", "MEMORY_LIMIT="+fmt.Sprintf("%d", memoryLimit), "-v", dir+"/"+submission_dir+":/"+submissionId, "my_imm").Run()
	if err != nil {
		println("error running docker:", err)
		return
	}

	meta_files, err := os.Open(meta_dir)
	if err != nil {
		println("error opening meta_dir:", err)
		return
	}
	input_files, err := os.Open(submission_dir + "/input")
	if err != nil {
		println("error opening meta_dir:", err)
		return
	}
	defer meta_files.Close()
	defer input_files.Close()

	meta_info_files, err := meta_files.Readdir(-1)
	if err != nil {
		println("error reading meta_dir:", err)
		return
	}
	for _, meta_info_files := range meta_info_files {
		metaFile, err := os.ReadFile(meta_dir + "/" + meta_info_files.Name())
		if err != nil {
			println("error reading meta_file:", err)
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
	}
	// compare output with checker
	input_files_info, err := input_files.Readdir(-1)
	if err != nil {
		println("error reading input_dir:", err.Error())
		return
	}
	for i, input_files_info := range input_files_info {
		input_file := submission_dir + "/input/" + input_files_info.Name()
		outut_file := submission_dir + "/output/" + input_files_info.Name()
		answer_file := submission_dir + "/answer/" + input_files_info.Name()
		verdict_file := verdict_dir + "/" + input_files_info.Name()

		// Judging the output
		err := exec.Command("bash","-c","./"+checker_dir+" "+ input_file+" "+outut_file+" "+answer_file+" 2> "+verdict_file).Run()
		if err != nil {
			println("Judge error on test case:", i+1, err.Error())
			return;
		}
		verdict, err := os.ReadFile(verdict_file)
		if err != nil {
			println("error reading verdict_file:", err)
			break
		}
		if string(verdict)[:2] != "ok" {
			println("wrong answer on test case", i+1)
			return
		} else {
			println("test case", i+1, "passed")
		}
	}
	println("AC")
}


func Compile(inpath, outpath string) bool {
	return exec.Command("g++", inpath, "-o", outpath).Run() == nil
}