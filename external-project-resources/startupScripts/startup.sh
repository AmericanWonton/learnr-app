#!/usr/bin/env bash

#This writes any output into a logfile in this same directory
function logger ()
{
    #Check to see if logger file has been created for today
    if [ -f "$FULLFILEPATH" ]; then
        #File exists, we can write here
        echo "File exists"
    else
        #File does NOT exist, create it
        echo "File does not exist, creating it rn..."
        sudo touch $FULLFILEPATH
    fi

    echo -e "$ADATE -- $AMESSAGE\n" >> $FULLFILEPATH
    
    return 0
}

#Get Current Date as a format for files
date=$(date '+%Y-%m-%d')
ADATE=$date
FULLFILENAME="startupLogger-$ADATE.log"
#FULLFILEPATH="/home/joek2/go-workspace/src/learnrapp/external-project-resources/startupScripts/$FULLFILENAME"
FULLFILEPATH="/home/ubuntu/startUpCronJob/logging/$FULLFILENAME"

echo $ADATE
echo $FULLFILENAME
echo $FULLFILEPATH


#Call Logger for debug print to start out with
AMESSAGE="We are starting the startupscript for today: $ADATE"
logger $AMESSAGE $FULLFILEPATH $ADATE

#Update some stuff
sudo apt update -y && sudo apt upgrade -y
sudo apt-get update -y && sudo apt-get upgrade -y
sudo apt autoremove -y

#See if docker containers are running; if they are, stop and delete them
sudo docker kill $(docker ps -q)
sudo docker rm -f $(docker ps -a -q)
sudo docker rmi $(docker images -q) -f

#Use Docker Credentials
sudo docker login --username americanwonton --password peanutdoggydoo111
#Pull all relevant docker images
sudo docker pull americanwonton/prodcrudproj:latest
sudo docker pull americanwonton/prodtextproj:latest
sudo docker pull americanwonton/prodlearnrproj:latest

#Run docker containers with their unique envioronment listings
#Careful sleep
sleep 5
#CRUD
sudo docker run -it --env-file /home/ubuntu/startUpCronJob/crud-env.list -d -p 4000:4000 americanwonton/prodcrudproj
#TEXT
#Careful sleep
sleep 5
sudo docker run -it --env-file /home/ubuntu/startUpCronJob/text-env.list -d -p 3000:3000 americanwonton/prodtextproj
#WEB
#Careful sleep
sleep 5
sudo docker run -it --env-file /home/ubuntu/startUpCronJob/web-env.list -d -p 80:8080 americanwonton/prodlearnrproj