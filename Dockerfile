FROM golang:1.21

RUN apt-get install -y git

CMD ["/bin/bash","-c","./help.sh go_build"]