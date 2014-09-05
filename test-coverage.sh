#!/bin/bash

echo "mode: set" > acc.out
for Dir in $(find . -type d ); 
do
	if ls $Dir/*.go &> /dev/null;
	then
		returnval=`go test -coverprofile=profile.out $Dir`
		echo ${returnval}
		if [[ ${returnval} != *FAIL* ]]
		then
    		if [ -f profile.out ]
    		then
        		cat profile.out | grep -v "mode: set" >> acc.out 
    		fi
    	else
    		exit 1
    	fi	
    fi
done
if [ -n "$GOVERALLS_TOKEN" ]
then
    $HOME/gopath/bin/goveralls -coverprofile=acc.out -repotoken $GOVERALLS_TOKEN
fi	

rm -rf ./profile.out
rm -rf ./acc.out