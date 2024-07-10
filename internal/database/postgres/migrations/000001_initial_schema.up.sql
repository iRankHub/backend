CREATE TABLE Users (
   UserID SERIAL PRIMARY KEY,
   Name VARCHAR(255) NOT NULL,
   Email VARCHAR(255) UNIQUE NOT NULL,
   Password VARCHAR(255) NOT NULL,
   UserRole VARCHAR(50) NOT NULL,
   Status VARCHAR(20) DEFAULT 'pending' CHECK (Status IN ('pending', 'approved', 'rejected')),
   VerificationStatus BOOLEAN DEFAULT FALSE,
   DeactivatedAt TIMESTAMP,
   two_factor_secret VARCHAR(32),
   two_factor_enabled BOOLEAN DEFAULT FALSE,
   failed_login_attempts INTEGER DEFAULT 0,
   last_login_attempt TIMESTAMP,
   last_logout TIMESTAMP,
   reset_token VARCHAR(64),
   reset_token_expires TIMESTAMP,
   biometric_token VARCHAR(64),
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE UserProfiles (
   ProfileID SERIAL PRIMARY KEY,
   UserID INTEGER UNIQUE NOT NULL REFERENCES Users(UserID),
   Name VARCHAR(255) NOT NULL,
   UserRole VARCHAR(50) NOT NULL,
   Email VARCHAR(255) NOT NULL,
   Address VARCHAR(255),
   Phone VARCHAR(20),
   Bio TEXT,
   ProfilePicture BYTEA,
   VerificationStatus BOOLEAN DEFAULT FALSE
);

-- Notifications table
CREATE TABLE Notifications (
    NotificationID SERIAL PRIMARY KEY,
    UserID INTEGER NOT NULL REFERENCES Users(UserID),
    Type VARCHAR(50) NOT NULL,
    Message TEXT NOT NULL,
    IsRead BOOLEAN DEFAULT FALSE,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- NotificationPreferences table
CREATE TABLE NotificationPreferences (
    PreferenceID SERIAL PRIMARY KEY,
    UserID INTEGER NOT NULL REFERENCES Users(UserID),
    EmailNotifications BOOLEAN DEFAULT TRUE,
    EmailFrequency VARCHAR(20) DEFAULT 'daily' CHECK (EmailFrequency IN ('daily', 'weekly', 'monthly')),
    EmailDay INTEGER CHECK (EmailDay >= 1 AND EmailDay <= 7),
    EmailTime TIME,
    InAppNotifications BOOLEAN DEFAULT TRUE
);

-- Create Indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON Users(Email);
CREATE INDEX IF NOT EXISTS idx_users_status ON Users(Status);
CREATE INDEX IF NOT EXISTS idx_users_biometric_token ON Users(biometric_token);
CREATE INDEX IF NOT EXISTS idx_users_reset_token ON Users(reset_token);

CREATE INDEX IF NOT EXISTS idx_schools_contactpersonid ON Schools(ContactPersonID);
CREATE INDEX IF NOT EXISTS idx_schools_contactemail ON Schools(ContactEmail);

CREATE INDEX IF NOT EXISTS idx_students_email ON Students(Email);
CREATE INDEX IF NOT EXISTS idx_students_schoolid ON Students(SchoolID);
CREATE INDEX IF NOT Exists idx_students_userid ON Students(UserID);

CREATE INDEX IF NOT EXISTS idx_volunteers_userid ON Volunteers(UserID);

CREATE INDEX IF NOT EXISTS idx_notifications_userid ON Notifications(UserID);

-- Create Triggers
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON Users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TABLE TournamentFormats (
   FormatID SERIAL PRIMARY KEY,
   FormatName VARCHAR(255) NOT NULL,
   Description TEXT
);

CREATE TABLE Leagues (
    LeagueID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    LeagueType VARCHAR(50) NOT NULL CHECK (LeagueType IN ('local', 'international'))
);

CREATE TABLE LocalLeagueDetails (
    LeagueID INTEGER PRIMARY KEY REFERENCES Leagues(LeagueID),
    Province VARCHAR(255),
    District VARCHAR(255)
);

CREATE TABLE InternationalLeagueDetails (
    LeagueID INTEGER PRIMARY KEY REFERENCES Leagues(LeagueID),
    Continent VARCHAR(255),
    Country VARCHAR(255)
);

CREATE TABLE Tournaments (
    TournamentID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    StartDate DATE NOT NULL,
    EndDate DATE NOT NULL,
    Location VARCHAR(255) NOT NULL,
    FormatID INTEGER NOT NULL REFERENCES TournamentFormats(FormatID),
    LeagueID INTEGER REFERENCES Leagues(LeagueID)
);

CREATE TABLE Schools (
   SchoolID SERIAL PRIMARY KEY,
   SchoolName VARCHAR(255) NOT NULL,
   Address VARCHAR(255) NOT NULL,
   Country VARCHAR(255),
   Province VARCHAR(255),
   District VARCHAR(255),
   ContactPersonID INTEGER NOT NULL REFERENCES Users(UserID),
   ContactEmail VARCHAR(255) NOT NULL UNIQUE,
   SchoolType VARCHAR(50) NOT NULL
);

CREATE TABLE Students (
   StudentID SERIAL PRIMARY KEY,
   FirstName VARCHAR(255) NOT NULL,
   LastName VARCHAR(255) NOT NULL,
   Grade VARCHAR(10) NOT NULL,
   DateOfBirth DATE,
   Email VARCHAR(255) UNIQUE,
   Password VARCHAR(255) NOT NULL,
   SchoolID INTEGER NOT NULL REFERENCES Schools(SchoolID),
   UserID INTEGER NOT NULL REFERENCES Users(UserID)
);

CREATE TABLE Teams (
   TeamID SERIAL PRIMARY KEY,
   Name VARCHAR(255) NOT NULL,
   SchoolID INTEGER NOT NULL REFERENCES Schools(SchoolID),
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID)
);

CREATE TABLE TeamMembers (
   TeamID INTEGER NOT NULL REFERENCES Teams(TeamID),
   StudentID INTEGER NOT NULL REFERENCES Students(StudentID),
   PRIMARY KEY (TeamID, StudentID)
);

CREATE TABLE Volunteers (
   VolunteerID SERIAL PRIMARY KEY,
   FirstName VARCHAR(255) NOT NULL,
   LastName VARCHAR(255) NOT NULL,
   DateOfBirth DATE,
   Role VARCHAR(50) NOT NULL,
   GraduateYear INTEGER CHECK (GraduateYear >= 2000 AND GraduateYear <= EXTRACT(YEAR FROM CURRENT_DATE)),
   Password VARCHAR(255) NOT NULL,
   SafeGuardCertificate BOOLEAN DEFAULT FALSE,
   UserID INTEGER NOT NULL REFERENCES Users(UserID)
);

CREATE TABLE Rounds (
   RoundID SERIAL PRIMARY KEY,
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   RoundNumber INTEGER NOT NULL,
   IsEliminationRound BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE Rooms (
   RoomID SERIAL PRIMARY KEY,
   RoomName VARCHAR(255) NOT NULL,
   Location VARCHAR(255) NOT NULL,
   Capacity INTEGER NOT NULL
);

CREATE TABLE Debates (
   DebateID SERIAL PRIMARY KEY,
   RoundID INTEGER NOT NULL REFERENCES Rounds(RoundID),
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   Team1ID INTEGER NOT NULL REFERENCES Teams(TeamID),
   Team2ID INTEGER NOT NULL REFERENCES Teams(TeamID),
   StartTime TIMESTAMP NOT NULL,
   EndTime TIMESTAMP,
   RoomID INTEGER NOT NULL REFERENCES Rooms(RoomID),
   Status VARCHAR(50) NOT NULL
);

CREATE TABLE JudgeAssignments (
   AssignmentID SERIAL PRIMARY KEY,
   VolunteerID INTEGER NOT NULL REFERENCES Volunteers(VolunteerID),
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   DebateID INTEGER NOT NULL REFERENCES Debates(DebateID)
);

CREATE TABLE Ballots (
   BallotID SERIAL PRIMARY KEY,
   DebateID INTEGER NOT NULL REFERENCES Debates(DebateID),
   JudgeID INTEGER NOT NULL REFERENCES JudgeAssignments(AssignmentID),
   Team1DebaterAScore NUMERIC,
   Team1DebaterAComments TEXT,
   Team1DebaterBScore NUMERIC,
   Team1DebaterBComments TEXT,
   Team1DebaterCScore NUMERIC,
   Team1DebaterCComments TEXT,
   Team1TotalScore NUMERIC,
   Team2DebaterAScore NUMERIC,
   Team2DebaterAComments TEXT,
   Team2DebaterBScore NUMERIC,
   Team2DebaterBComments TEXT,
   Team2DebaterCScore NUMERIC,
   Team2DebaterCComments TEXT,
   Team2TotalScore NUMERIC
);

CREATE TABLE TournamentCoordinators (
   CoordinatorID SERIAL PRIMARY KEY,
   FirstName VARCHAR(255) NOT NULL,
   LastName VARCHAR(255) NOT NULL,
   Email VARCHAR(255) UNIQUE NOT NULL,
   Password VARCHAR(255) NOT NULL,
   UserID INTEGER NOT NULL REFERENCES Users(UserID)
);

CREATE TABLE Schedules (
   ScheduleID SERIAL PRIMARY KEY,
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   RoundID INTEGER NOT NULL REFERENCES Rounds(RoundID),
   DebateID INTEGER NOT NULL REFERENCES Debates(DebateID),
   ScheduledTime TIMESTAMP NOT NULL
);

CREATE TABLE Results (
   ResultID SERIAL PRIMARY KEY,
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   TeamID INTEGER NOT NULL REFERENCES Teams(TeamID),
   Rank INTEGER,
   Points NUMERIC
);

CREATE TABLE RoomBookings (
   BookingID SERIAL PRIMARY KEY,
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   RoomID INTEGER NOT NULL REFERENCES Rooms(RoomID),
   StartTime TIMESTAMP NOT NULL,
   EndTime TIMESTAMP NOT NULL
);

CREATE TABLE Communications (
   CommunicationID SERIAL PRIMARY KEY,
   UserID INTEGER NOT NULL REFERENCES Users(UserID),
   SchoolID INTEGER NOT NULL REFERENCES Schools(SchoolID),
   Type VARCHAR(50) NOT NULL,
   Content TEXT NOT NULL,
   Timestamp TIMESTAMP NOT NULL
);

CREATE TABLE JudgeReviews (
   ReviewID SERIAL PRIMARY KEY,
   StudentID INTEGER NOT NULL REFERENCES Students(StudentID),
   JudgeID INTEGER NOT NULL REFERENCES Users(UserID),
   Rating NUMERIC NOT NULL,
   Comments TEXT
);

CREATE TABLE VolunteerRatingTypes (
   RatingTypeID SERIAL PRIMARY KEY,
   RatingTypeName VARCHAR(255) NOT NULL
);

CREATE TABLE VolunteerRatings (
   RatingID SERIAL PRIMARY KEY,
   VolunteerID INTEGER NOT NULL REFERENCES Volunteers(VolunteerID),
   RatingTypeID INTEGER NOT NULL REFERENCES VolunteerRatingTypes(RatingTypeID),
   RatingScore NUMERIC NOT NULL,
   RatingComments TEXT,
   CumulativeRating NUMERIC
);

CREATE TABLE StudentRanks (
   RankID SERIAL PRIMARY KEY,
   StudentID INTEGER NOT NULL REFERENCES Students(StudentID),
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   RankValue NUMERIC NOT NULL,
   RankComments TEXT
);

CREATE TABLE StudentTransfers (
   TransferID SERIAL PRIMARY KEY,
   StudentID INTEGER NOT NULL REFERENCES Students(StudentID),
   FromSchoolID INTEGER REFERENCES Schools(SchoolID),
   ToSchoolID INTEGER NOT NULL REFERENCES Schools(SchoolID),
   TransferDate DATE NOT NULL,
   Reason VARCHAR(255)
);