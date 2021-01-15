FROM golang

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN mkdir texts

EXPOSE 8080

RUN go build -o main .

CMD ./main