#! /bin/bash

cd checker
#GOOS=linux go build -o demo_checker main.go
GOOS=linux GOARCH=arm64 go build -o app_example_checker checker.go
cd ..

mkdir app_example_checker
mv checker/app_example_checker ./app_example_checker/

tar zcvf app_example_checker.tar.gz app_example_checker

scp app_example_checker.tar.gz root@139.198.21.191:/var/www/html

rm -rf app_example_checker
