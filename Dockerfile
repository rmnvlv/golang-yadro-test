FROM golang:latest
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o ./app .
CMD ["./app" "--path=/test_file.txt"]