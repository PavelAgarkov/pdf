FROM golang:1.21

RUN apt-get install -y git

CMD ["/bin/bash","-c","./prod.sh backend_build"]