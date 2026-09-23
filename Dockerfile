from golang:latest as builder
label name="Golang"
workdir /app 
copy go.mod ./
run go mod download
copy . . 
run go build -o html-go .


from alpine:latest 
workdir /app 
copy --from=builder /app/html-go ./
copy --from=builder /app/templates ./templates
copy --from=builder /app/style ./style
run adduser -D rajsoni
user rajsoni
expose 8080 
cmd ["./html-go"]
