# aws
Proof of concept for aws

## 1 Configure
    go mod init git@github.com:NmoEvs/aws.git 

## 2 Download dependencies
    go get

## 3 Build
    go build -o ../bin

# Lambda
⚠️ need to zip exec from linux + chmod +x exec
    
    cd bin
    zip lambda.zip lambda