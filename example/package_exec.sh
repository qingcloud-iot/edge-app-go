#! /bin/bash

cd checker
GOOS=linux go build -o app_example_checker checker.go
cd ..

mkdir app_example_checker
mv checker/app_example_checker ./app_example_checker/

tar zcvf app_example_checker.tar.gz app_example_checker

rm -rf app_example_checker
