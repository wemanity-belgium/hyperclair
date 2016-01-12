FROM golang:1.5.2

RUN git config --global user.email "garciagonzalez.julien@gmail.com"
RUN git config --global user.name "jgsqware"

RUN go get -v github.com/spf13/cobra/cobra
