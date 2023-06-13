#!/bin/zsh

app_name="polar-savannah-16447"

# Loop over each line in .env
while IFS= read -r line
do
    # Split the line into name and value
    name="${line%=*}"
    value="${line#*=}"
    
    # Set the variable in Heroku
    heroku config:set "$name=$value" --app "$app_name"
done < .env
