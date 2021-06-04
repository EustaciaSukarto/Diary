CREATE TABLE user (
    ID int NOT NULL AUTO_INCREMENT,
    Fullname varchar(45),
    Birthday date,
    Email varchar(45) NOT NULL UNIQUE,
    Username varchar(45) NOT NULL UNIQUE,
    Password longtext NOT NULL,
    PRIMARY KEY (ID)
);
CREATE TABLE entry (
    ID int NOT NULL AUTO_INCREMENT,
    UserID int NOT NULL,
    Content longtext,
    PRIMARY KEY (ID),
    FOREIGN KEY (UserID) REFERENCES user(ID)
);
