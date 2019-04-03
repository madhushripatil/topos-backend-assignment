# Building Footprint Data collector and Analyzer

This is a simple ETL process that allows users to collect and analyze Building footprints dataset. It is a REST endpoint supporting basic APIs like inserting, updating, deleting and querying Building footprints data.

Transformations like Splitting, Joining and Cleaning are performed on the provided dataset.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project.

### Prerequisites

1. #### The Go Programming Language

> Go is an open source programming language that makes it easy to build simple,
reliable, and efficient software.

### Download and Install

#### Binary Distributions

Official binary distributions are available at https://golang.org/dl/.

After downloading a binary release, visit https://golang.org/doc/install
or load [doc/install.html](./doc/install.html) in your web browser for installation
instructions.

#### Install From Source

If a binary distribution is not available for your combination of
operating system and architecture, visit
https://golang.org/doc/install/source or load [doc/install-source.html](./doc/install-source.html)
in your web browser for source installation instructions.

Set GOROOT, GOPOATH env variables

Check if the correct values are set using
$ go env

2. #### Mongo DB

MongoDB offers both an Enterprise and Community version of its powerful non-relational database. 

https://www.mongodb.com/download-center/community

Create a username and password to connect to MongoDB and enable auth while starting mongodb service.

Create a database to be use for the application.

3. #### Gorilla Mux Package for routing

[Gorilla Mux](https://github.com/gorilla/mux) - go get github.com/gorilla/mux

4. #### Mgo Package for MongoDB connections

[Rich MongoDB driver](https://gopkg.in/mgo.v2) - go get gopkg.in/mgo.v2

5. #### bson Package for BSON GO specification

[GO BSON Specification implementation](https://gopkg.in/mgo.v2/bson) - go get gopkg.in/mgo.v2/bson

6. #### joho godotenv Package for loading env file

[Application specific configuration](https://github.com/joho/godotenv) - go get github.com/joho/godotenv

## Deployment

1. Set PATH variables

``vi ~/.bashrc``

Set GOPATH, GOROOT, PATH variables

Set BUILDING_ENV variable to your application env filename, example: export BUILDING_ENV = environment

Save the bashrc file and exit

2. ``source ~/.bashrc``

Verify if all System variables are set by running the following command:

3. ``go env``

Also check if PATH and BUILDING_ENV variables are set by running the following commands:

``echo $PATH``

``echo $BUILDING_ENV``

Download the [Project zip](https://github.com/madhushripatil/topos-backend-assignment/archive/master.zip) to your $GOPATH/src directory.

4. Set your Database and REST Endpoint parameters in the [development.env](https://github.com/madhushripatil/topos-backend-assignment/blob/master/development.env) file provided.

There is a separate development.env file in the importData directory as well.

5. Run the following commands

``cd $GOPATH/src/topos-backend-assignment/importData``

``go build importCSVToMongo.go``

``./importCSVToMongo``

``go build importBuildingFeatTypeCSVToMongo.go``

``./importBuildingFeatTypeCSVToMongo``

``go build importBoroughCSVToMongo.go``

``./importBoroughCSVToMongo``

6. Run the following commands inside the $GOPATH/src/project directory:

``cd $GOPATH/src/topos-backend-assignment``

``go build``

``go run server.go``

The Server starts running.

You can now start making API calls.

### REST API Documentation

Refer to the [REST API Documentation](https://documenter.getpostman.com/view/2410794/S1EH21eE)

You may use Curl or Postman REST client to run the APIs provided.

## Built With

* [Go](https://golang.org/) - Open source programming language
* [Mongo DB](https://www.mongodb.com/what-is-mongodb) - Document Database with scalability and flexibility
* [GoLand IDE](https://www.jetbrains.com/go/?utm_expid=.qV9Irwa4SS-xPJHMhpNehw.0&utm_referrer=) A clever IDE to GO

## Future Scope

Advanced Transformations like 

[ ] Summarization

[ ] Aggregation

[ ] Data Validation and Integration can be done.
