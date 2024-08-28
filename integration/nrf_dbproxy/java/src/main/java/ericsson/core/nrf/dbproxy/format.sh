#!/bin/bash

rm *orig*
function format_dir() {
    for file in `ls $1 | grep -v format.sh`
    do
    	if [ -d $1"/"$file ]
    	then
		rm $1"/"$file"/"*orig*
    		format_dir $1"/"$file
    	else
   		astyle --style=linux $1"/"$file
    	fi
    done
}
format_dir .
function delete_internal_files() {
    for file in `ls $1 | grep -v format.sh`
    do
    	if [ -d $1"/"$file ]
    	then
		rm $1"/"$file"/"*orig*
    		delete_internal_files $1"/"$file
    	fi
    done
}
delete_internal_files .
rm *orig*
