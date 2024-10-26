-- Create tables with foreign key dependencies
CREATE TABLE UserProfiles (
   ProfileID SERIAL PRIMARY KEY,
   UserID INTEGER UNIQUE NOT NULL REFERENCES Users(UserID),
   Name VARCHAR(255) NOT NULL,
   UserRole VARCHAR(50) NOT NULL,
   Email VARCHAR(255) NOT NULL,
   Password VARCHAR(255) NOT NULL,
   Gender VARCHAR(10) CHECK (Gender IN ('male', 'female', 'non-binary')),
   Address VARCHAR(255),
   Phone VARCHAR(20),
   Bio TEXT,
   ProfilePicture VARCHAR(2048),
   VerificationStatus BOOLEAN DEFAULT FALSE
);



CREATE TABLE NotificationPreferences (
    PreferenceID SERIAL PRIMARY KEY,
    UserID INTEGER NOT NULL REFERENCES Users(UserID),
    EmailNotifications BOOLEAN DEFAULT TRUE,
    EmailFrequency VARCHAR(20) DEFAULT 'daily' CHECK (EmailFrequency IN ('daily', 'weekly', 'monthly')),
    EmailDay INTEGER CHECK (EmailDay >= 1 AND EmailDay <= 7),
    EmailTime TIME,
    InAppNotifications BOOLEAN DEFAULT TRUE
);

CREATE TABLE WebAuthnCredentials (
    ID SERIAL PRIMARY KEY,
    UserID INTEGER NOT NULL REFERENCES Users(UserID),
    CredentialID BYTEA NOT NULL,
    PublicKey BYTEA NOT NULL,
    AttestationType VARCHAR(255) NOT NULL,
    AAGUID BYTEA NOT NULL,
    SignCount BIGINT NOT NULL,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE WebAuthnSessionData (
    UserID INTEGER PRIMARY KEY REFERENCES Users(UserID),
    SessionData BYTEA NOT NULL
);

CREATE TABLE Tournaments (
    TournamentID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    StartDate TIMESTAMP NOT NULL,
    EndDate TIMESTAMP NOT NULL,
    Location VARCHAR(255) NOT NULL,
    FormatID INTEGER NOT NULL REFERENCES TournamentFormats(FormatID),
    LeagueID INTEGER REFERENCES Leagues(LeagueID),
    CoordinatorID INTEGER NOT NULL REFERENCES Users(UserID),
    NumberOfPreliminaryRounds INTEGER NOT NULL,
    NumberOfEliminationRounds INTEGER NOT NULL,
    JudgesPerDebatePreliminary INTEGER NOT NULL,
    JudgesPerDebateElimination INTEGER NOT NULL,
    TournamentFee DECIMAL(10, 2) NOT NULL,
    ImageUrl VARCHAR(2048),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    yesterday_total_count INT DEFAULT 0,
    yesterday_upcoming_count INT DEFAULT 0,
    yesterday_active_debaters_count INT DEFAULT 0
);

CREATE TABLE Rooms (
   RoomID SERIAL PRIMARY KEY,
   RoomName VARCHAR(255) NOT NULL,
   Location VARCHAR(255) NOT NULL,
   Capacity INTEGER NOT NULL,
   TournamentID INTEGER REFERENCES Tournaments(TournamentID)
);

CREATE TABLE Schools (
   SchoolID SERIAL PRIMARY KEY,
   iDebateSchoolID VARCHAR(35) UNIQUE,
   SchoolName VARCHAR(255) NOT NULL,
   Address VARCHAR(255) NOT NULL,
   Country VARCHAR(255),
   Province VARCHAR(255),
   District VARCHAR(255),
   ContactPersonID INTEGER NOT NULL REFERENCES Users(UserID),
   ContactPersonNationalID VARCHAR(50),
   ContactEmail VARCHAR(255) NOT NULL UNIQUE,
   SchoolEmail VARCHAR(255) NOT NULL UNIQUE,
   SchoolType VARCHAR(50) NOT NULL CHECK (SchoolType IN ('Private', 'Public', 'Government Aided', 'International'))
);

CREATE TABLE Students (
   StudentID SERIAL PRIMARY KEY,
   iDebateStudentID VARCHAR(20) UNIQUE,
   FirstName VARCHAR(255) NOT NULL,
   LastName VARCHAR(255) NOT NULL,
   Gender VARCHAR(10) CHECK (Gender IN ('male', 'female', 'non-binary')),
   Grade VARCHAR(10) NOT NULL,
   DateOfBirth DATE,
   Email VARCHAR(255) UNIQUE,
   Password VARCHAR(255) NOT NULL,
   SchoolID INTEGER NOT NULL REFERENCES Schools(SchoolID),
   UserID INTEGER NOT NULL REFERENCES Users(UserID)
);

CREATE TABLE Volunteers (
   VolunteerID SERIAL PRIMARY KEY,
   iDebateVolunteerID VARCHAR(20) UNIQUE,
   FirstName VARCHAR(255) NOT NULL,
   LastName VARCHAR(255) NOT NULL,
   Gender VARCHAR(10) CHECK (Gender IN ('male', 'female', 'non-binary')),
   DateOfBirth DATE,
   NationalID VARCHAR(50),
   Role VARCHAR(50) NOT NULL,
   GraduateYear INTEGER CHECK (GraduateYear >= 2000 AND GraduateYear <= EXTRACT(YEAR FROM CURRENT_DATE)),
   Password VARCHAR(255) NOT NULL,
   SafeGuardCertificate VARCHAR(2048),
   HasInternship BOOLEAN DEFAULT FALSE,
   IsEnrolledInUniversity BOOLEAN DEFAULT FALSE,
   UserID INTEGER NOT NULL REFERENCES Users(UserID),
   yesterday_rounds_judged INT DEFAULT 0,
   yesterday_tournaments_attended INT DEFAULT 0,
   yesterday_upcoming_tournaments INT DEFAULT 0
);

CREATE TABLE TournamentInvitations (
    InvitationID SERIAL PRIMARY KEY,
    TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
    InviteeID VARCHAR(35) NOT NULL, -- This will store the iDebate ID
    InviteeRole VARCHAR(20) NOT NULL CHECK (InviteeRole IN ('school', 'volunteer', 'student')),
    Status VARCHAR(20) NOT NULL CHECK (Status IN ('pending', 'accepted', 'rejected')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ReminderSentAt TIMESTAMP
);

CREATE TABLE Teams (
   TeamID SERIAL PRIMARY KEY,
   Name VARCHAR(255) NOT NULL,
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   TotalWins INTEGER DEFAULT 0,
   TotalSpeakerPoints NUMERIC(8,2) DEFAULT 0,
   AverageRank NUMERIC(5,2) DEFAULT 0
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


CREATE TABLE Debates (
   DebateID SERIAL PRIMARY KEY,
   RoundID INTEGER NOT NULL,
   RoundNumber INTEGER NOT NULL,
   IsEliminationRound BOOLEAN NOT NULL DEFAULT FALSE,
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   Team1ID INTEGER NOT NULL REFERENCES Teams(TeamID),
   Team2ID INTEGER NOT NULL REFERENCES Teams(TeamID),
   StartTime TIMESTAMP NOT NULL,
   EndTime TIMESTAMP,
   RoomID INTEGER NOT NULL REFERENCES Rooms(RoomID),
   Status VARCHAR(50) NOT NULL DEFAULT 'scheduled'
);

-- Create JudgeAssignments table
CREATE TABLE JudgeAssignments (
    AssignmentID SERIAL PRIMARY KEY,
    TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
    JudgeID INTEGER NOT NULL REFERENCES Users(UserID),
    DebateID INTEGER NOT NULL REFERENCES Debates(DebateID),
    RoundNumber INTEGER NOT NULL,
    IsElimination BOOLEAN NOT NULL DEFAULT FALSE,
    IsHeadJudge BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE(TournamentID, JudgeID, RoundNumber, IsElimination)
);

CREATE TABLE JudgeFeedback (
    FeedbackID SERIAL PRIMARY KEY,
    JudgeID INTEGER REFERENCES Users(UserID),
    StudentID INTEGER REFERENCES Students(StudentID),
    DebateID INTEGER REFERENCES Debates(DebateID),
    ClarityRating NUMERIC(5,2) CHECK (ClarityRating BETWEEN 0 AND 100),
    ConstructivenessRating NUMERIC(5,2) CHECK (ConstructivenessRating BETWEEN 0 AND 100),
    TimelinessRating NUMERIC(5,2) CHECK (TimelinessRating BETWEEN 0 AND 100),
    FairnessRating NUMERIC(5,2) CHECK (FairnessRating BETWEEN 0 AND 100),
    EngagementRating NUMERIC(5,2) CHECK (EngagementRating BETWEEN 0 AND 100),
    AverageRating NUMERIC(5,2) GENERATED ALWAYS AS (
        (ClarityRating + ConstructivenessRating + TimelinessRating + FairnessRating + EngagementRating) / 5
    ) STORED,
    TextFeedback TEXT,
    IsRead BOOLEAN DEFAULT FALSE,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Ballots (
   BallotID SERIAL PRIMARY KEY,
   DebateID INTEGER NOT NULL REFERENCES Debates(DebateID),
   JudgeID INTEGER NOT NULL,
   Team1TotalScore NUMERIC,
   Team1Feedback TEXT,
   Team2TotalScore NUMERIC,
   Team2Feedback TEXT,
   RecordingStatus VARCHAR(20) NOT NULL DEFAULT 'pending',
   Verdict VARCHAR(255) NOT NULL DEFAULT 'pending',
   last_updated_by INT REFERENCES Users(UserID),
   last_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   head_judge_submitted BOOLEAN DEFAULT FALSE
);


-- Create SpeakerScores table
CREATE TABLE SpeakerScores (
    ScoreID SERIAL PRIMARY KEY,
    BallotID INTEGER NOT NULL REFERENCES Ballots(BallotID),
    SpeakerID INTEGER NOT NULL REFERENCES Students(StudentID),
    SpeakerRank INTEGER NOT NULL,
    SpeakerPoints NUMERIC(5,2) NOT NULL,
    Feedback TEXT,
    UNIQUE(BallotID, SpeakerID),
    IsRead BOOLEAN DEFAULT FALSE
);

-- Create PairingHistory table
CREATE TABLE PairingHistory (
    HistoryID SERIAL PRIMARY KEY,
    TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
    Team1ID INTEGER NOT NULL REFERENCES Teams(TeamID),
    Team2ID INTEGER NOT NULL REFERENCES Teams(TeamID),
    RoundNumber INTEGER NOT NULL,
    IsElimination BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE(TournamentID, Team1ID, Team2ID, RoundNumber, IsElimination)
);

CREATE TABLE IF NOT EXISTS DebateJudges (
    DebateID INTEGER NOT NULL REFERENCES Debates(DebateID),
    JudgeID INTEGER NOT NULL REFERENCES Volunteers(VolunteerID),
    PRIMARY KEY (DebateID, JudgeID)
);

CREATE TABLE Schedules (
   ScheduleID SERIAL PRIMARY KEY,
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   RoundID INTEGER NOT NULL REFERENCES Rounds(RoundID),
   DebateID INTEGER NOT NULL REFERENCES Debates(DebateID),
   ScheduledTime TIMESTAMP NOT NULL
);


CREATE TABLE TeamScores (
    ScoreID SERIAL PRIMARY KEY,
    TeamID INTEGER REFERENCES Teams(TeamID),
    DebateID INTEGER REFERENCES Debates(DebateID),
    TotalScore NUMERIC(5,2),
    Rank INTEGER,
    IsElimination BOOLEAN
);

CREATE TABLE Communications (
   CommunicationID SERIAL PRIMARY KEY,
   UserID INTEGER NOT NULL REFERENCES Users(UserID),
   SchoolID INTEGER NOT NULL REFERENCES Schools(SchoolID),
   Type VARCHAR(50) NOT NULL,
   Content TEXT NOT NULL,
   Timestamp TIMESTAMP NOT NULL
);

CREATE TABLE VolunteerRatings (
   RatingID SERIAL PRIMARY KEY,
   VolunteerID INTEGER NOT NULL REFERENCES Volunteers(VolunteerID),
   RatingTypeID INTEGER NOT NULL REFERENCES VolunteerRatingTypes(RatingTypeID),
   RatingScore NUMERIC NOT NULL,
   RatingComments TEXT,
   CumulativeRating NUMERIC
);

CREATE TABLE StudentTransfers (
   TransferID SERIAL PRIMARY KEY,
   StudentID INTEGER NOT NULL REFERENCES Students(StudentID),
   FromSchoolID INTEGER REFERENCES Schools(SchoolID),
   ToSchoolID INTEGER NOT NULL REFERENCES Schools(SchoolID),
   TransferDate DATE NOT NULL,
   Reason VARCHAR(255)
);

CREATE TABLE TournamentExpenses (
    ExpenseID SERIAL PRIMARY KEY,
    TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
    FoodExpense DECIMAL(10, 2) NOT NULL DEFAULT 0,
    TransportExpense DECIMAL(10, 2) NOT NULL DEFAULT 0,
    PerDiemExpense DECIMAL(10, 2) NOT NULL DEFAULT 0,
    AwardingExpense DECIMAL(10, 2) NOT NULL DEFAULT 0,
    StationaryExpense DECIMAL(10, 2) NOT NULL DEFAULT 0,
    OtherExpenses DECIMAL(10, 2) NOT NULL DEFAULT 0,
    TotalExpense DECIMAL(10, 2) GENERATED ALWAYS AS (
        FoodExpense + TransportExpense + PerDiemExpense +
        AwardingExpense + StationaryExpense + OtherExpenses
    ) STORED,
    Currency VARCHAR(3) NOT NULL DEFAULT 'RWF',
    Notes TEXT,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CreatedBy INTEGER REFERENCES Users(UserID),
    UpdatedBy INTEGER REFERENCES Users(UserID)
);

CREATE TABLE SchoolTournamentRegistrations (
    RegistrationID SERIAL PRIMARY KEY,
    SchoolID INTEGER NOT NULL REFERENCES Schools(SchoolID),
    TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
    PlannedTeamsCount INTEGER NOT NULL,
    ActualTeamsCount INTEGER,
    AmountPerTeam DECIMAL(10, 2) NOT NULL,
    TotalAmount DECIMAL(10, 2) GENERATED ALWAYS AS (PlannedTeamsCount * AmountPerTeam) STORED,
    DiscountAmount DECIMAL(10, 2) DEFAULT 0,
    ActualPaidAmount DECIMAL(10, 2),
    PaymentStatus VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (PaymentStatus IN ('pending', 'partial', 'paid', 'cancelled')),
    PaymentDate TIMESTAMP,
    Currency VARCHAR(3) NOT NULL DEFAULT 'RWF',
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CreatedBy INTEGER REFERENCES Users(UserID),
    UpdatedBy INTEGER REFERENCES Users(UserID),
    UNIQUE(SchoolID, TournamentID)
);