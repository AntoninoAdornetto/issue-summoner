#!/bin/bash

# Variables
name="World"
current_date=$(date +"%Y-%m-%d")

echo "Hello, $name!"
echo "Today's date is $current_date
and I like to do da dance
"

# @TEST_ANNOTATION iterate through 10
for i in {1..5}; do
    echo "Iteration $i"
done

greet() {
    echo "Welcome, $1!"
}

greet "User"

exit 0
