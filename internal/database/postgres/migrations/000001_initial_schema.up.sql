-- Create Users table first as it's referenced by many other tables
CREATE TABLE Users (
   UserID SERIAL PRIMARY KEY,
   WebAuthnUserID BYTEA UNIQUE,
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
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP
);

-- Biometric Credentials Table
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

-- Biometric Session Data
CREATE TABLE WebAuthnSessionData (
    UserID INTEGER PRIMARY KEY REFERENCES Users(UserID),
    SessionData BYTEA NOT NULL
);

-- Create UserProfiles table
CREATE TABLE UserProfiles (
   ProfileID SERIAL PRIMARY KEY,
   UserID INTEGER UNIQUE NOT NULL REFERENCES Users(UserID),
   Name VARCHAR(255) NOT NULL,
   UserRole VARCHAR(50) NOT NULL,
   Email VARCHAR(255) NOT NULL,
   Password VARCHAR(255) NOT NULL,
   Address VARCHAR(255),
   Phone VARCHAR(20),
   Bio TEXT,
   ProfilePicture BYTEA,
   VerificationStatus BOOLEAN DEFAULT FALSE
);

-- Create Notifications table
CREATE TABLE Notifications (
    NotificationID SERIAL PRIMARY KEY,
    UserID INTEGER NOT NULL REFERENCES Users(UserID),
    Type VARCHAR(50) NOT NULL,
    Message TEXT NOT NULL,
    IsRead BOOLEAN DEFAULT FALSE,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create NotificationPreferences table
CREATE TABLE NotificationPreferences (
    PreferenceID SERIAL PRIMARY KEY,
    UserID INTEGER NOT NULL REFERENCES Users(UserID),
    EmailNotifications BOOLEAN DEFAULT TRUE,
    EmailFrequency VARCHAR(20) DEFAULT 'daily' CHECK (EmailFrequency IN ('daily', 'weekly', 'monthly')),
    EmailDay INTEGER CHECK (EmailDay >= 1 AND EmailDay <= 7),
    EmailTime TIME,
    InAppNotifications BOOLEAN DEFAULT TRUE
);

-- Create TournamentFormats table
CREATE TABLE TournamentFormats (
   FormatID SERIAL PRIMARY KEY,
   FormatName VARCHAR(255) NOT NULL,
   Description TEXT,
   SpeakersPerTeam INTEGER NOT NULL,
   deleted_at TIMESTAMP
);

-- Create Leagues table
CREATE TABLE Leagues (
    LeagueID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    LeagueType VARCHAR(50) NOT NULL CHECK (LeagueType IN ('local', 'international')),
    Details JSONB NOT NULL DEFAULT '{}',
    deleted_at TIMESTAMP
);


-- Create Tournaments table
CREATE TABLE Tournaments (
    TournamentID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    StartDate DATE NOT NULL,
    EndDate DATE NOT NULL,
    Location VARCHAR(255) NOT NULL,
    FormatID INTEGER NOT NULL REFERENCES TournamentFormats(FormatID),
    LeagueID INTEGER REFERENCES Leagues(LeagueID),
    NumberOfPreliminaryRounds INTEGER NOT NULL,
    NumberOfEliminationRounds INTEGER NOT NULL,
    JudgesPerDebatePreliminary INTEGER NOT NULL,
    JudgesPerDebateElimination INTEGER NOT NULL,
    TournamentFee DECIMAL(10, 2) NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE CountryCodes (
    CountryName VARCHAR(255) PRIMARY KEY,
    IsoCode CHAR(3) NOT NULL UNIQUE
);

-- Create Schools table
CREATE TABLE Schools (
   SchoolID SERIAL PRIMARY KEY,
   iDebateSchoolID VARCHAR(35) UNIQUE,
   SchoolName VARCHAR(255) NOT NULL,
   Address VARCHAR(255) NOT NULL,
   Country VARCHAR(255),
   Province VARCHAR(255),
   District VARCHAR(255),
   ContactPersonID INTEGER NOT NULL REFERENCES Users(UserID),
   ContactEmail VARCHAR(255) NOT NULL UNIQUE,
   SchoolType VARCHAR(50) NOT NULL CHECK (SchoolType IN ('Private', 'Public', 'Government Aided', 'International'))
);

-- Create Students table
CREATE TABLE Students (
   StudentID SERIAL PRIMARY KEY,
   iDebateStudentID VARCHAR(20) UNIQUE,
   FirstName VARCHAR(255) NOT NULL,
   LastName VARCHAR(255) NOT NULL,
   Grade VARCHAR(10) NOT NULL,
   DateOfBirth DATE,
   Email VARCHAR(255) UNIQUE,
   Password VARCHAR(255) NOT NULL,
   SchoolID INTEGER NOT NULL REFERENCES Schools(SchoolID),
   UserID INTEGER NOT NULL REFERENCES Users(UserID)
);

-- Create Teams table
CREATE TABLE Teams (
   TeamID SERIAL PRIMARY KEY,
   Name VARCHAR(255) NOT NULL,
   SchoolID INTEGER NOT NULL REFERENCES Schools(SchoolID),
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID)
);

-- Create TeamMembers table
CREATE TABLE TeamMembers (
   TeamID INTEGER NOT NULL REFERENCES Teams(TeamID),
   StudentID INTEGER NOT NULL REFERENCES Students(StudentID),
   PRIMARY KEY (TeamID, StudentID)
);

-- Create Volunteers table
CREATE TABLE Volunteers (
   VolunteerID SERIAL PRIMARY KEY,
   iDebateVolunteerID VARCHAR(20) UNIQUE,
   FirstName VARCHAR(255) NOT NULL,
   LastName VARCHAR(255) NOT NULL,
   DateOfBirth DATE,
   Role VARCHAR(50) NOT NULL,
   GraduateYear INTEGER CHECK (GraduateYear >= 2000 AND GraduateYear <= EXTRACT(YEAR FROM CURRENT_DATE)),
   Password VARCHAR(255) NOT NULL,
   SafeGuardCertificate BOOLEAN DEFAULT FALSE,
   UserID INTEGER NOT NULL REFERENCES Users(UserID)
);

-- Create Rounds table
CREATE TABLE Rounds (
   RoundID SERIAL PRIMARY KEY,
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   RoundNumber INTEGER NOT NULL,
   IsEliminationRound BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create Rooms table
CREATE TABLE Rooms (
   RoomID SERIAL PRIMARY KEY,
   RoomName VARCHAR(255) NOT NULL,
   Location VARCHAR(255) NOT NULL,
   Capacity INTEGER NOT NULL
);

-- Create Debates table
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

-- Create JudgeAssignments table
CREATE TABLE JudgeAssignments (
   AssignmentID SERIAL PRIMARY KEY,
   VolunteerID INTEGER NOT NULL REFERENCES Volunteers(VolunteerID),
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   DebateID INTEGER NOT NULL REFERENCES Debates(DebateID)
);

-- Create Ballots table
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

-- Create TournamentCoordinators table
CREATE TABLE TournamentCoordinators (
   CoordinatorID SERIAL PRIMARY KEY,
   VolunteerID INTEGER NOT NULL REFERENCES Volunteers(VolunteerID),
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   AssignedDate DATE NOT NULL DEFAULT CURRENT_DATE
);

-- Create Schedules table
CREATE TABLE Schedules (
   ScheduleID SERIAL PRIMARY KEY,
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   RoundID INTEGER NOT NULL REFERENCES Rounds(RoundID),
   DebateID INTEGER NOT NULL REFERENCES Debates(DebateID),
   ScheduledTime TIMESTAMP NOT NULL
);

-- Create Results table
CREATE TABLE Results (
   ResultID SERIAL PRIMARY KEY,
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   TeamID INTEGER NOT NULL REFERENCES Teams(TeamID),
   Rank INTEGER,
   Points NUMERIC
);

-- Create RoomBookings table
CREATE TABLE RoomBookings (
   BookingID SERIAL PRIMARY KEY,
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   RoomID INTEGER NOT NULL REFERENCES Rooms(RoomID),
   StartTime TIMESTAMP NOT NULL,
   EndTime TIMESTAMP NOT NULL
);

-- Create Communications table
CREATE TABLE Communications (
   CommunicationID SERIAL PRIMARY KEY,
   UserID INTEGER NOT NULL REFERENCES Users(UserID),
   SchoolID INTEGER NOT NULL REFERENCES Schools(SchoolID),
   Type VARCHAR(50) NOT NULL,
   Content TEXT NOT NULL,
   Timestamp TIMESTAMP NOT NULL
);

-- Create JudgeReviews table
CREATE TABLE JudgeReviews (
   ReviewID SERIAL PRIMARY KEY,
   StudentID INTEGER NOT NULL REFERENCES Students(StudentID),
   JudgeID INTEGER NOT NULL REFERENCES Users(UserID),
   Rating NUMERIC NOT NULL,
   Comments TEXT
);

-- Create VolunteerRatingTypes table
CREATE TABLE VolunteerRatingTypes (
   RatingTypeID SERIAL PRIMARY KEY,
   RatingTypeName VARCHAR(255) NOT NULL
);

-- Create VolunteerRatings table
CREATE TABLE VolunteerRatings (
   RatingID SERIAL PRIMARY KEY,
   VolunteerID INTEGER NOT NULL REFERENCES Volunteers(VolunteerID),
   RatingTypeID INTEGER NOT NULL REFERENCES VolunteerRatingTypes(RatingTypeID),
   RatingScore NUMERIC NOT NULL,
   RatingComments TEXT,
   CumulativeRating NUMERIC
);

-- Create StudentRanks table
CREATE TABLE StudentRanks (
   RankID SERIAL PRIMARY KEY,
   StudentID INTEGER NOT NULL REFERENCES Students(StudentID),
   TournamentID INTEGER NOT NULL REFERENCES Tournaments(TournamentID),
   RankValue NUMERIC NOT NULL,
   RankComments TEXT
);

-- Create StudentTransfers table
CREATE TABLE StudentTransfers (
   TransferID SERIAL PRIMARY KEY,
   StudentID INTEGER NOT NULL REFERENCES Students(StudentID),
   FromSchoolID INTEGER REFERENCES Schools(SchoolID),
   ToSchoolID INTEGER NOT NULL REFERENCES Schools(SchoolID),
   TransferDate DATE NOT NULL,
   Reason VARCHAR(255)
);

-- Create Indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON Users(Email);
CREATE INDEX IF NOT EXISTS idx_users_status ON Users(Status);

CREATE INDEX IF NOT EXISTS idx_users_reset_token ON Users(reset_token);

CREATE INDEX IF NOT EXISTS idx_schools_contactpersonid ON Schools(ContactPersonID);
CREATE INDEX IF NOT EXISTS idx_schools_contactemail ON Schools(ContactEmail);

CREATE INDEX IF NOT EXISTS idx_students_email ON Students(Email);
CREATE INDEX IF NOT EXISTS idx_students_schoolid ON Students(SchoolID);
CREATE INDEX IF NOT EXISTS idx_students_userid ON Students(UserID);

CREATE INDEX IF NOT EXISTS idx_volunteers_userid ON Volunteers(UserID);

CREATE INDEX IF NOT EXISTS idx_notifications_userid ON Notifications(UserID);

INSERT INTO CountryCodes (IsoCode, CountryName) VALUES
('AFG', 'Afghanistan'),
('ALA', 'Aland Islands'),
('ALB', 'Albania'),
('DZA', 'Algeria'),
('ASM', 'American Samoa'),
('AND', 'Andorra'),
('AGO', 'Angola'),
('AIA', 'Anguilla'),
('ATA', 'Antarctica'),
('ATG', 'Antigua and Barbuda'),
('ARG', 'Argentina'),
('ARM', 'Armenia'),
('ABW', 'Aruba'),
('AUS', 'Australia'),
('AUT', 'Austria'),
('AZE', 'Azerbaijan'),
('BHS', 'Bahamas'),
('BHR', 'Bahrain'),
('BGD', 'Bangladesh'),
('BRB', 'Barbados'),
('BLR', 'Belarus'),
('BEL', 'Belgium'),
('BLZ', 'Belize'),
('BEN', 'Benin'),
('BMU', 'Bermuda'),
('BTN', 'Bhutan'),
('BOL', 'Bolivia'),
('BES', 'Bonaire, Sint Eustatius and Saba'),
('BIH', 'Bosnia and Herzegovina'),
('BWA', 'Botswana'),
('BVT', 'Bouvet Island'),
('BRA', 'Brazil'),
('IOT', 'British Indian Ocean Territory'),
('BRN', 'Brunei Darussalam'),
('BGR', 'Bulgaria'),
('BFA', 'Burkina Faso'),
('BDI', 'Burundi'),
('KHM', 'Cambodia'),
('CMR', 'Cameroon'),
('CAN', 'Canada'),
('CPV', 'Cape Verde'),
('CYM', 'Cayman Islands'),
('CAF', 'Central African Republic'),
('TCD', 'Chad'),
('CHL', 'Chile'),
('CHN', 'China'),
('CXR', 'Christmas Island'),
('CCK', 'Cocos (Keeling) Islands'),
('COL', 'Colombia'),
('COM', 'Comoros'),
('COG', 'Congo'),
('COD', 'Congo, The Democratic Republic of'),
('COK', 'Cook Islands'),
('CRI', 'Costa Rica'),
('CIV', 'Cote d''Ivoire'),
('HRV', 'Croatia'),
('CUB', 'Cuba'),
('CUW', 'Curaçao'),
('CYP', 'Cyprus'),
('CZE', 'Czechia'),
('DNK', 'Denmark'),
('DJI', 'Djibouti'),
('DMA', 'Dominica'),
('DOM', 'Dominican Republic'),
('ECU', 'Ecuador'),
('EGY', 'Egypt'),
('SLV', 'El Salvador'),
('GNQ', 'Equatorial Guinea'),
('ERI', 'Eritrea'),
('EST', 'Estonia'),
('ETH', 'Ethiopia'),
('FLK', 'Falkland Islands (Malvinas)'),
('FRO', 'Faroe Islands'),
('FJI', 'Fiji'),
('FIN', 'Finland'),
('FRA', 'France'),
('GUF', 'French Guiana'),
('PYF', 'French Polynesia'),
('ATF', 'French Southern Territories'),
('GAB', 'Gabon'),
('GMB', 'Gambia'),
('GEO', 'Georgia'),
('DEU', 'Germany'),
('GHA', 'Ghana'),
('GIB', 'Gibraltar'),
('GRC', 'Greece'),
('GRL', 'Greenland'),
('GRD', 'Grenada'),
('GLP', 'Guadeloupe'),
('GUM', 'Guam'),
('GTM', 'Guatemala'),
('GGY', 'Guernsey'),
('GIN', 'Guinea'),
('GNB', 'Guinea-Bissau'),
('GUY', 'Guyana'),
('HTI', 'Haiti'),
('HMD', 'Heard and Mc Donald Islands'),
('VAT', 'Holy See (Vatican City State)'),
('HND', 'Honduras'),
('HKG', 'Hong Kong'),
('HUN', 'Hungary'),
('ISL', 'Iceland'),
('IND', 'India'),
('IDN', 'Indonesia'),
('IRN', 'Iran, Islamic Republic of'),
('IRQ', 'Iraq'),
('IRL', 'Ireland'),
('IMN', 'Isle of Man'),
('ISR', 'Israel'),
('ITA', 'Italy'),
('JAM', 'Jamaica'),
('JPN', 'Japan'),
('JEY', 'Jersey'),
('JOR', 'Jordan'),
('KAZ', 'Kazakstan'),
('KEN', 'Kenya'),
('KIR', 'Kiribati'),
('PRK', 'Korea, Democratic People''s Republic of'),
('KOR', 'Korea, Republic of'),
('XKX', 'Kosovo'),
('KWT', 'Kuwait'),
('KGZ', 'Kyrgyzstan'),
('LAO', 'Lao, People''s Democratic Republic'),
('LVA', 'Latvia'),
('LBN', 'Lebanon'),
('LSO', 'Lesotho'),
('LBR', 'Liberia'),
('LBY', 'Libyan Arab Jamahiriya'),
('LIE', 'Liechtenstein'),
('LTU', 'Lithuania'),
('LUX', 'Luxembourg'),
('MAC', 'Macao'),
('MKD', 'Macedonia, The Former Yugoslav Republic Of'),
('MDG', 'Madagascar'),
('MWI', 'Malawi'),
('MYS', 'Malaysia'),
('MDV', 'Maldives'),
('MLI', 'Mali'),
('MLT', 'Malta'),
('MHL', 'Marshall Islands'),
('MTQ', 'Martinique'),
('MRT', 'Mauritania'),
('MUS', 'Mauritius'),
('MYT', 'Mayotte'),
('MEX', 'Mexico'),
('FSM', 'Micronesia, Federated States of'),
('MDA', 'Moldova, Republic of'),
('MCO', 'Monaco'),
('MNG', 'Mongolia'),
('MNE', 'Montenegro'),
('MSR', 'Montserrat'),
('MAR', 'Morocco'),
('MOZ', 'Mozambique'),
('MMR', 'Myanmar'),
('NAM', 'Namibia'),
('NRU', 'Nauru'),
('NPL', 'Nepal'),
('NLD', 'Netherlands'),
('NCL', 'New Caledonia'),
('NZL', 'New Zealand'),
('NIC', 'Nicaragua'),
('NER', 'Niger'),
('NGA', 'Nigeria'),
('NIU', 'Niue'),
('NFK', 'Norfolk Island'),
('MNP', 'Northern Mariana Islands'),
('NOR', 'Norway'),
('OMN', 'Oman'),
('PAK', 'Pakistan'),
('PLW', 'Palau'),
('PSE', 'Palestinian Territory, Occupied'),
('PAN', 'Panama'),
('PNG', 'Papua New Guinea'),
('PRY', 'Paraguay'),
('PER', 'Peru'),
('PHL', 'Philippines'),
('PCN', 'Pitcairn'),
('POL', 'Poland'),
('PRT', 'Portugal'),
('PRI', 'Puerto Rico'),
('QAT', 'Qatar'),
('SRB', 'Republic of Serbia'),
('REU', 'Reunion'),
('ROU', 'Romania'),
('RUS', 'Russia Federation'),
('RWA', 'Rwanda'),
('BLM', 'Saint Barthélemy'),
('SHN', 'Saint Helena'),
('KNA', 'Saint Kitts & Nevis'),
('LCA', 'Saint Lucia'),
('MAF', 'Saint Martin'),
('SPM', 'Saint Pierre and Miquelon'),
('VCT', 'Saint Vincent and the Grenadines'),
('WSM', 'Samoa'),
('SMR', 'San Marino'),
('STP', 'Sao Tome and Principe'),
('SAU', 'Saudi Arabia'),
('SEN', 'Senegal'),
('SYC', 'Seychelles'),
('SLE', 'Sierra Leone'),
('SGP', 'Singapore'),
('SXM', 'Sint Maarten'),
('SVK', 'Slovakia'),
('SVN', 'Slovenia'),
('SLB', 'Solomon Islands'),
('SOM', 'Somalia'),
('ZAF', 'South Africa'),
('SGS', 'South Georgia & The South Sandwich Islands'),
('SSD', 'South Sudan'),
('ESP', 'Spain'),
('LKA', 'Sri Lanka'),
('SDN', 'Sudan'),
('SUR', 'Suriname'),
('SJM', 'Svalbard and Jan Mayen'),
('SWZ', 'Swaziland'),
('SWE', 'Sweden'),
('CHE', 'Switzerland'),
('SYR', 'Syrian Arab Republic'),
('TWN', 'Taiwan, Province of China'),
('TJK', 'Tajikistan'),
('TZA', 'Tanzania, United Republic of'),
('THA', 'Thailand'),
('TLS', 'Timor-Leste'),
('TGO', 'Togo'),
('TKL', 'Tokelau'),
('TON', 'Tonga'),
('TTO', 'Trinidad and Tobago'),
('TUN', 'Tunisia'),
('TUR', 'Turkey'),
('TKM', 'Turkmenistan'),
('TCA', 'Turks and Caicos Islands'),
('TUV', 'Tuvalu'),
('UGA', 'Uganda'),
('UKR', 'Ukraine'),
('ARE', 'United Arab Emirates'),
('GBR', 'United Kingdom'),
('USA', 'United States of America'),
('UMI', 'United States Minor Outlying Islands'),
('URY', 'Uruguay'),
('UZB', 'Uzbekistan'),
('VUT', 'Vanuatu'),
('VEN', 'Venezuela'),
('VNM', 'Vietnam'),
('VGB', 'Virgin Islands, British'),
('VIR', 'Virgin Islands, U.S.'),
('WLF', 'Wallis and Futuna'),
('ESH', 'Western Sahara'),
('YEM', 'Yemen'),
('ZMB', 'Zambia'),
('ZWE', 'Zimbabwe');

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

-- functions and triggers to handle the Idebate IDs generation

-- function to generate the school ID
CREATE OR REPLACE FUNCTION generate_idebate_school_id()
RETURNS trigger AS $$
DECLARE
    country_code CHAR(3);
    province_code CHAR(2);
    district_letter CHAR(1);
    type_letters CHAR(2);
    random_number INT;
BEGIN
    -- Look up the country code
    SELECT IsoCode INTO country_code
    FROM CountryCodes
    WHERE CountryName ILIKE NEW.Country
    LIMIT 1;

    -- If no matching country found, use 'XXX'
    IF country_code IS NULL THEN
        country_code := 'XXX';
    END IF;

    -- Set province code
    IF country_code = 'RWA' THEN
        -- For Rwanda, use single letter province codes
        CASE
            WHEN NEW.Province ILIKE 'East%' THEN province_code := 'E';
            WHEN NEW.Province ILIKE 'West%' THEN province_code := 'W';
            WHEN NEW.Province ILIKE 'South%' THEN province_code := 'S';
            WHEN NEW.Province ILIKE 'North%' THEN province_code := 'N';
            WHEN NEW.Province ILIKE 'Kigali%' THEN province_code := 'K';
            ELSE province_code := 'X'; -- For unknown province in Rwanda
        END CASE;
    ELSE
        -- For other countries, use first two letters of the province
        province_code := UPPER(LEFT(NEW.Province, 2));
    END IF;

    -- Set district letter (first letter of district name)
    district_letter := UPPER(LEFT(NEW.District, 1));

    -- Set type letters based on SchoolType
    CASE
        WHEN NEW.SchoolType = 'Private' THEN type_letters := 'PV';
        WHEN NEW.SchoolType = 'Public' THEN type_letters := 'PB';
        WHEN NEW.SchoolType = 'Government Aided' THEN type_letters := 'GA';
        WHEN NEW.SchoolType = 'International' THEN type_letters := 'IN';
        ELSE type_letters := 'OT'; -- For other types (shouldn't occur due to CHECK constraint)
    END CASE;

    -- Generate random number (1 to 99999)
    random_number := floor(random() * 99999 + 1);

    -- Combine all parts to form the ID
    NEW.iDebateSchoolID := country_code || '-' || province_code || '-' || district_letter || '-' || type_letters || '-' || LPAD(random_number::TEXT, 5, '0');

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to automatically generate iDebate School ID
CREATE OR REPLACE TRIGGER set_idebate_school_id
BEFORE INSERT ON Schools
FOR EACH ROW
EXECUTE FUNCTION generate_idebate_school_id();



-- function to create IDs for Volunteers

CREATE SEQUENCE idebate_volunteer_id_seq START 1;

CREATE OR REPLACE FUNCTION generate_idebate_volunteer_id()
RETURNS trigger AS $$
BEGIN
  NEW.iDebateVolunteerID := 'VOLU' || LPAD(NEXTVAL('idebate_volunteer_id_seq')::TEXT, 6, '0');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- triggers to execute the function of volunteer_id before the insertion of the rest column details.

CREATE TRIGGER set_idebate_volunteer_id
BEFORE INSERT ON Volunteers
FOR EACH ROW
EXECUTE FUNCTION generate_idebate_volunteer_id();



-- For Students
CREATE SEQUENCE idebate_student_id_seq START 1;

CREATE OR REPLACE FUNCTION generate_idebate_student_id()
RETURNS trigger AS $$
BEGIN
  NEW.iDebateStudentID := 'STUD' || LPAD(NEXTVAL('idebate_student_id_seq')::TEXT, 6, '0');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- trigger to execute the function before insertion .

CREATE TRIGGER set_idebate_student_id
BEFORE INSERT ON Students
FOR EACH ROW
EXECUTE FUNCTION generate_idebate_student_id();