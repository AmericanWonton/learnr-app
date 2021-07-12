#!/bin/sh
#See this for additional scripting notes: https://www.shellscript.sh/
#Set this to readable/writeable when created: chmod a+rx update-script.sh
#Start announcing the script
echo "Doing test script"
#Update some stuff
sudo apt-get update -y && sudo apt-get upgrade -y
sudo apt autoremove -y
#See if docker containers are running; if they are, stop and delete them
sudo docker kill $(docker ps -q)
sudo docker rm -f $(docker ps -a -q)
sudo docker rmi $(docker images -q) -f
#Pull the latest image to run, then run it
sudo docker login --username americanwonton --password peanutdoggydoo111
sudo docker pull americanwonton/prodlearnrproj:latest
sudo docker run -it --env-file ./env.list -d -p 80:8080 americanwonton/prodcrudproj