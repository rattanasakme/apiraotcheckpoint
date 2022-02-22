FROM golang:1.16 as final

WORKDIR /app
EXPOSE 80
EXPOSE 443


COPY go.mod .
COPY go.sum .
RUN go mod download

# Change timezone to local time
ENV TZ=Asia/Bangkok
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

RUN go build -o ./out/dist .

CMD ./out/dist



