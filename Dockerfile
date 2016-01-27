FROM golang:1.5.2

RUN git config user.email "garciagonzalez.julien@gmail.com"
RUN git config user.name "wemanity-belgium"

RUN go get -v github.com/spf13/cobra/cobra
