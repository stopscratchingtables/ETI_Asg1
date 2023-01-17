# ETI Assignment 1 - Ride Sharing Platform 🚗
## Table of Contents  
[Introduction](#introduction)

[To-Do-List](#todolist)

[Instructions for setting up](#instructions)

[Architecture Diagram](#arcdiagram)

[Micro Services](#microservices)

[Console Functions](#confunc)

## Introduction

<a name="introduction"/>
Ride Sharing platform is a platform where passengers can book a ride and will get allocated a driver. The driver will be able to end the ride as well. 

This assignment uses MySQL, GoLang and API Implementation with the use of Micro-Services through a Domain-Driven-Design

<a name="todolist"/>

## Check-List of features

#### General

- [x] Passanger Login/Logout
- [x] Driver Login/Logout
- [x] Login/Logout Validation
- [ ] Implement GraphSQL Code
- [ ] Implement a front-end interface

#### Passengers

- [x] Update Account Details
- [x] View Ride History
- [x] Book a Ride
- [x] Validation

#### Driver

- [x] Update Account Details (except identification number)
- [x] End Ride
- [x] Validation

<a name="instructions"/>

## Instructions on setting up and running microservices

1. Run the SQL code in "sqlTables.sql" on MySQL
2. Run the golang code in "driversAPI.go" in "microservice - drivers" at port 1000
3. Run the golang code in "rideAPI.go" in "microservice - ride" at port 5000
4. Run the golang code in "passangersAPI.go" in "microservice - passengers" at port 6969
5. Run "main.go" in consoleFolder using command prompt terminal

<a name="arcdiagram"/>

## Domain Driven Design Architecture Diagrams

The domain of this application is the Ride Sharing Platform.
The subdomains are the Driver, Passenger and Ride. The driver is the core subdomain out of all the three as it is the most essential to ensure the Ride-Sharing Platform is running. Ride and passenger are supporting subdomains as they are not as essential to the processes of a ride-sharing platform but instead, acts as a supoort for the core subdomain "Driver".

![Ride Sharing Platform Domain Driven Design Diagram]([Ride-Sharing Platform Domain Driven Design.png](https://github.com/stopscratchingtables/ETI_Asg1/blob/main/Ride-Sharing%20Platform%20Domain%20Driven%20Design.png?raw=true))

<a name="microservices"/>

## Micro-Services

Each micro-service has been assigned it's own table in the "my_db" database.
The Passenger information is allocated in the Passanger Table, driver information in the Driver Table and ride information in the Ride Table.

Although traditionally different micro-services are assigned different databases, within the context of this assignment, the domain design of ride-sharing platform will not be broken down into different databases and tables instead. This is because there is not many tables required for each micro-service compared to the traditional numerous tables that a micro-service might have. Therefore, 

#### Passenger Micro-Service

Description:
This service performs most of the required tasks and duties of a passenger in the ride-sharing platform. Necessary information is obtained from the passenger and ride tables from this database.

1. GET Passenger (From Passenger Table)
2. POST Passenger (From Passenger Table)
3. PUT Passenger (From Passenger Table)
4. DELETE Passenger (From Passenger Table)
5. GET Available Rides (From Ride Table)

#### Driver Micro-Service

Description:
This service performs most of the required tasks and duties of a driver in the ride-sharing platform. Necessary information is obtained from the driver and ride tables from this database.

1. GET Driver (From Driver Table)
2. POST Driver (From Driver Table)
3. PUT Driver (From Driver Table)
4. DELETE Driver (From Driver Table)
5. GET Available Rides (From Ride Table)

#### Ride Micro-Service

Description:
This service performs most of the required tasks and duties of a passenger in the ride-sharing platform. Necessary information is obtained from the ride tables from the database. This micro-service works in conjunction with ride and driver micro-services to get the necessary information required from each other.

1. GET Rides (From Ride Table)
2. POST Rides (From Ride Table)
3. PUT Ride (From Ride Table)
4. DELETE Ride (From Ride Table)

<a name="confunc"/>

## Console Functions

#### Passenger functions
- Book Ride
- View Riding History
- Update Account Details

#### Driver functions
- End Ride
- Update Account Details


