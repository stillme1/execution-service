import os
import resource
import subprocess

# Set the resource limits
cpu_time_limit = (int)(os.environ.get("TIME_LIMIT")) # in seconds
memory_limit = 1024 * 1024 * (int)(os.environ.get("MEMORY_LIMIT"))  # in MegaBytes

# Get the submission ID from the SUBMISSION_ID environment variable
submission_id = os.environ.get("SUBMISSION_ID")

# Set the paths to the input, output, meta, error and executable files
input_dir = f"/{submission_id}/input"       # input file path
output_dir = f"/{submission_id}/output"     # output file path
meta_dir = f"/{submission_id}/meta"         # timetaken memoryused exitcode file path
error_dir = f"/{submission_id}/error"       # stderr file path
executable_file = f"./{submission_id}/exec" # executable file path

# Find all input files in the input directory
input_files = [f for f in os.listdir(input_dir) if os.path.isfile(os.path.join(input_dir, f))]

# Run the executable for each input file and write the output to the output directory
for input_file in input_files:
    input_path = os.path.join(input_dir, input_file)
    output_path = os.path.join(output_dir, input_file)
    meta_path = os.path.join(meta_dir, input_file)
    error_path = os.path.join(error_dir, input_file)

    with open(output_path, "w") as output_file, open(input_path, "r") as input_file:
        # Set the resource limits for the execution
        resource.setrlimit(resource.RLIMIT_CPU, (cpu_time_limit, cpu_time_limit))
        resource.setrlimit(resource.RLIMIT_DATA, (memory_limit, memory_limit))

        start_time = resource.getrusage(resource.RUSAGE_CHILDREN).ru_utime
        with open(error_path, "w") as error_file:
            process = subprocess.run(
                ["/usr/bin/time",'--format="%M"',
                executable_file],
                stdin=input_file,
                stdout=output_file,
                stderr=error_file
            )
        end_time = resource.getrusage(resource.RUSAGE_CHILDREN).ru_utime

        with open(error_path, "r") as error_file:
            max_rss = error_file.readlines()[-1].replace("\"", "").replace("\n", "")

        if(process.returncode == 0):
            with open(meta_path, "w") as meta_file:
                meta_file.write(str(end_time-start_time)+" "+str(max_rss)+" 0")
        else:
            with open(meta_path, "w") as meta_file:
                meta_file.write(str(end_time-start_time)+" "+str(max_rss)+" "
                                 +str(process.returncode))
            break