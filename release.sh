CGO_ENABLED=1 go build -race -o ./biehdc.webapp.applister.linux.withrace
CGO_ENABLED=0 go build -o ./biehdc.webapp.applister.linux
GOOS=windows go build -o ./biehdc.webapp.applister.exe
