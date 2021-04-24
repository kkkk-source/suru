#!/bin/bash
#docker run -it suru --apiUrl="url" --fromEmail "email" --fromEmailKey "secret" --toEmail "email"
go run main.go broker.go bestbuy.go  --apiUrl "url" --fromEmail "email" --fromEmailKey "secret" --toEmail "email"
