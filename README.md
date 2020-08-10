
# go-simple-webapp

go-simple-webapp is an example of a basic crud webapp


 ## Status
[![Build Status](https://travis-ci.com/bishy999/go-simple-webapp.svg?branch=master)](https://travis-ci.com/bishy999/go-simple-webapp)
[![Go Report Card](https://goreportcard.com/badge/github.com/bishy999/go-simple-webapp)](https://goreportcard.com/report/github.com/bishy999/go-simple-webapp)
[![GoDoc](https://godoc.org/github.com/bishy999/go-simple-webapp/pkg/tag?status.svg)](https://godoc.org/github.com/bishy999/go-simple-webapp/pkg/app)
![GitHub Repo size](https://img.shields.io/github/repo-size/bishy999/go-simple-webapp)
[![GitHub Tag](https://img.shields.io/github/tag/bishy999/go-simple-webapp.svg)](https://github.com/bishy999/go-simple-webapp/releases/latest)
[![GitHub Activity](https://img.shields.io/github/commit-activity/m/bishy999/go-simple-webapp)](https://github.com/bishy999/go-simple-webapp)
[![GitHub Contributors](https://img.shields.io/github/contributors/bishy999/go-simple-webapp)](https://github.com/bishy999/go-simple-webapp)
[![codecov](https://codecov.io/gh/bishy999/go-simple-webapp/branch/master/graph/badge.svg)](https://codecov.io/gh/bishy999/go-simple-webapp)



# pull and run image

```
sudo docker run --name=go-simple-webapp -d -p 8080:8080 bishy999/go-simple-webapp:1.15
```

# sql image for this
```
see https://github.com/bishy999/mysql-simple-db
``` 


# check app is working via browser/cli

```
http://localhost:8080
curl http://localhost:8080
```


####################################################################
#          build and deploy simple golang webapp with docker       #
####################################################################


<u>Note:</u>  This is automatically done via travis (see .travis.yml) but manual steps are listed here for reference<br>


# build image (don't use cache)

```
sudo docker build --no-cache -t go-simple-webapp .
```


# list images

```
sudo docker images
```
 
# tag image

```
sudo docker tag <image ID>  <docker hub username>/<repository><image name>:<version label or tag>

e.g.

sudo docker tag go-simple-webapp bishy999/go-simple-webapp:1.15
```


# give terminal your docker hub credentials

```
sudo docker login
```


# push image to docker hub

```
docker push <docker hub username>/<repository><image name>

e.g. on Mac

sudo docker push bishy999/go-simple-webapp:1.15
```


# check docker hub

```
image used in example here is stored on docker hub e.g https://hub.docker.com/r/bishy999/go-simple-webapp
```

# create a container from your image and run it
 
```
sudo docker run --name=go-simple-webapp -d -p 8080:8080 bishy999/go-simple-webapp:1.15

```




# API Interaction
``

SwaggerUI is accessible via accessing localhost:8080/swaggerui/

To use the API you will need to generate a JWT token (JSON Web Tokens)

Select Authorize and add token

Bearer <token> e.g. Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTU5NDk4NTYsImlhdCI6MTU5NTk0NjI1NiwidXNlciI6ImJpc2h5OTk5QGhvdG1haWwuY29tIn0.asfclX0eGC6Ve-XvW2PL1UpGOtbibSIuiUdg5Xpk-T0


## Contributing

We love pull requests! Please see the [contribution guidelines](CONTRIBUTING.md).

