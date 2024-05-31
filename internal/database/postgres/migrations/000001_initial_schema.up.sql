CREATE TABLE Users (
  UserID SERIAL PRIMARY KEY,
  FirstName VARCHAR(255) NOT NULL,
  LastName VARCHAR(255) NOT NULL,
  Email VARCHAR(255) UNIQUE NOT NULL,
  Password VARCHAR(255) NOT NULL,
  UserRole VARCHAR(50) NOT NULL,
  Bio TEXT,
  VerificationStatus BOOLEAN DEFAULT FALSE,
  ApprovalStatus BOOLEAN DEFAULT FALSE
);

CREATE TABLE Schools (
  SchoolID SERIAL PRIMARY KEY,
  UserID INTEGER NOT NULL REFERENCES Users(UserID),
  Name VARCHAR(255) NOT NULL,
  Country VARCHAR(255) NOT NULL,
  Province VARCHAR(255),
  District VARCHAR(255),
  SchoolType VARCHAR(255) NOT NULL,
  ContactPersonName VARCHAR(255) NOT NULL,
  ContactPersonNumber VARCHAR(20) NOT NULL,
  ContactEmail VARCHAR(255) NOT NULL,
  UniqueSchoolID VARCHAR(8) UNIQUE NOT NULL
);

CREATE TABLE Students (
  StudentID SERIAL PRIMARY KEY,
  UserID INTEGER NOT NULL REFERENCES Users(UserID),
  DateOfBirth DATE,
  SchoolID INTEGER REFERENCES Schools(SchoolID),
  UniqueStudentID VARCHAR(8) UNIQUE NOT NULL
);

CREATE TABLE Volunteers (
  VolunteerID SERIAL PRIMARY KEY,
  UserID INTEGER NOT NULL REFERENCES Users(UserID),
  DateOfBirth DATE,
  NationalID VARCHAR(255) NOT NULL,
  SchoolAttended VARCHAR(255) NOT NULL,
  GraduationYear INTEGER,
  RoleInterestedIn VARCHAR(255),
  SafeguardingCertificate BYTEA,
  UniqueVolunteerID VARCHAR(8) UNIQUE NOT NULL
);
		CREATE TABLE TournamentFormats (
			FormatID SERIAL PRIMARY KEY,
			FormatName VARCHAR(255) NOT NULL,
			Description TEXT
		);
		CREATE TABLE Tournaments (
			TournamentID SERIAL PRIMARY KEY,
			Name VARCHAR(255) NOT NULL,
			StartDate DATE NOT NULL,
			EndDate DATE NOT NULL,
			Location VARCHAR(255) NOT NULL,
			FormatID INTEGER NOT NULL REFERENCES TournamentFormats(FormatID)
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