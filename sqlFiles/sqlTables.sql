CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'localhost';

DROP database my_db;

CREATE database my_db;
USE my_db;

CREATE TABLE Passenger (ID varchar(20) NOT NULL PRIMARY KEY, FirstName VARCHAR(50), LastName VARCHAR(50), MobileNum CHAR(8), EmailAddress varchar(50) NOT NULL);
INSERT INTO Passenger (ID, FirstName, LastName, MobileNum, EmailAddress) VALUES ("P0001", "Marcus", "Hon", "12345678", "wrongemail@np.edu.sg"); 
INSERT INTO Passenger (ID, FirstName, LastName, MobileNum, EmailAddress) VALUES ("P0054", "Patino", "Fernandez", "74443345", "pat@np.edu.sg"); 

CREATE TABLE Driver (ID varchar (20) NOT NULL PRIMARY KEY, FirstName VARCHAR(50), LastName VARCHAR(50), EmailAddress varchar(50), MobileNum CHAR(8), DriverStatus enum('Available', 'Busy'), LicenceNumber INT NOT NULL);
INSERT INTO Driver (ID, FirstName, LastName, EmailAddress, MobileNum, DriverStatus, LicenceNumber) VALUES ("D0001", "Jimmy", "Neutron", "uncleJimmy@ride.com", "78783434", "Available", 3456789); 
INSERT INTO Driver (ID, FirstName, LastName, EmailAddress, MobileNum, DriverStatus, LicenceNumber) VALUES ("D0002", "Jiggle", "Paff", "jgpf@ride.com", "78783434", "Available", 1234567); 
INSERT INTO Driver (ID, FirstName, LastName, EmailAddress, MobileNum, DriverStatus, LicenceNumber) VALUES ("D0003", "Antoine", "Lima", "aLima@ride.com", "65434455", "Available", 4455667); 
INSERT INTO Driver (ID, FirstName, LastName, EmailAddress, MobileNum, DriverStatus, LicenceNumber) VALUES ("D0092", "Montana", "Hannah", "monhan@ride.com", "87774353", "Busy", 9874355); 

CREATE TABLE Ride (ID varchar(25) NOT NULL PRIMARY KEY, PassengerID VARCHAR(50), DriverID VARCHAR(50), TripStatus ENUM("In Progress", "Finished", "Open"), StartDateTime DateTime, EndDateTime DateTime, PostalCodePickUp varchar(10), PostalCodeDropOff varchar(10));
INSERT INTO Ride (ID, PassengerID, DriverID, TripStatus, StartDateTime, EndDateTime, PostalCodePickUp, PostalCodeDropOff) VALUES ("R0001", "P0005", "D0002", "Finished", "2022-05-20 10:01:00.999999", "2022-05-20 10:45:00.999999", "347666", "348996"); 