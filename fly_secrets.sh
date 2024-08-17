#!/bin/bash

if [ ! -f .env ]; then
  echo ".env file not found!"
  exit 1
fi

command="fly secrets set"

while IFS= read -r line; do
  if [[ "$line" =~ ^#.*$ ]] || [[ -z "$line" ]]; then
    continue
  fi

  var_name=$(echo "$line" | cut -d '=' -f 1)
  var_value=$(echo "$line" | cut -d '=' -f 2-)

  if [[ -z "$var_value" ]]; then
    continue
  fi

  command+=" ${var_name}=${var_value}"
done < .env

eval "$command"
