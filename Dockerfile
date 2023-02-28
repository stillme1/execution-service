FROM python

RUN apt update
RUN apt install -y time

# take script.py as command
COPY script.py .
CMD ["/bin/bash", "-c", "python ./script.py"]