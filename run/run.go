package run

import (
	"os"
	"os/exec"
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
		os.RemoveAll(submission_dir)
		return
	}

	err = os.Mkdir(error_dir, 0777)
	if err != nil {
		os.RemoveAll(submission_dir)
		return
	}

	err = os.Mkdir(meta_dir, 0777)
	if err != nil {
		os.RemoveAll(submission_dir)
		return
	}

	// copy the file to submission directory
	// Later going to be fetched from database
	err = exec.Command("cp", "-r", input_dir, submission_dir).Run()
	if err != nil {
		panic(err)
	}

	// // time docker run -e SUBMISSION_ID="1" - TIME_LIMIT=2 -e MEMORY_LIMIT=100 -v /home/sahil/dockerTest/cpp:/1 my-image
	dir, _ := os.Getwd()
	err = exec.Command("docker", "run", "-e", "SUBMISSION_ID="+submissionId, "-e", "TIME_LIMIT=2", "-e", "MEMORY_LIMIT=100", "-v", dir+"/"+submission_dir+":/"+submissionId, "my-image").Run()
	if err != nil {
		return
	}
	//

	os.RemoveAll(submission_dir)
	println("exiting RunCode")
}
