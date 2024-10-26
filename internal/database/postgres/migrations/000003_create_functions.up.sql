-- For Students
CREATE SEQUENCE idebate_student_id_seq START 1;
CREATE SEQUENCE idebate_volunteer_id_seq START 1;

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to update historical data
CREATE OR REPLACE FUNCTION update_user_counts()
RETURNS VOID AS $$
BEGIN
  UPDATE Users
  SET yesterday_approved_count = (SELECT COUNT(*) FROM Users WHERE Status = 'approved' AND deleted_at IS NULL);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_tournament_counts()
RETURNS VOID AS $$
BEGIN
  UPDATE Tournaments
  SET
    yesterday_total_count = (SELECT COUNT(*) FROM Tournaments WHERE deleted_at IS NULL),
    yesterday_upcoming_count = (SELECT COUNT(*) FROM Tournaments WHERE deleted_at IS NULL AND StartDate BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '30 days'),
    yesterday_active_debaters_count = (
      SELECT COUNT(DISTINCT tm.StudentID)
      FROM TeamMembers tm
      JOIN Teams t ON tm.TeamID = t.TeamID
      JOIN Students s ON tm.StudentID = s.StudentID
      JOIN Tournaments tour ON t.TournamentID = tour.TournamentID
      WHERE tour.deleted_at IS NULL
    )
  WHERE TournamentID = (SELECT MIN(TournamentID) FROM Tournaments);
END;
$$ LANGUAGE plpgsql;

-- Function to update volunteer counts
CREATE OR REPLACE FUNCTION update_volunteer_counts()
RETURNS VOID AS $$
BEGIN
  UPDATE Volunteers v
  SET
    yesterday_rounds_judged = (
      SELECT COUNT(DISTINCT d.DebateID)
      FROM JudgeAssignments ja
      JOIN Debates d ON ja.DebateID = d.DebateID
      WHERE ja.JudgeID = v.UserID
    ),
    yesterday_tournaments_attended = (
      SELECT COUNT(DISTINCT t.TournamentID)
      FROM JudgeAssignments ja
      JOIN Debates d ON ja.DebateID = d.DebateID
      JOIN Tournaments t ON d.TournamentID = t.TournamentID
      WHERE ja.JudgeID = v.UserID
    ),
    yesterday_upcoming_tournaments = (
      SELECT COUNT(DISTINCT t.TournamentID)
      FROM TournamentInvitations ti
      JOIN Tournaments t ON ti.TournamentID = t.TournamentID
      WHERE ti.InviteeID = v.iDebateVolunteerID
        AND ti.Status = 'accepted'
        AND t.StartDate > CURRENT_DATE
    );
END;
$$ LANGUAGE plpgsql;


-- Corrected trigger function for calculating team average rank
CREATE OR REPLACE FUNCTION calculate_team_average_rank()
RETURNS TRIGGER AS $$
DECLARE
    avg_rank FLOAT;
    debate_record RECORD;
BEGIN
    IF NEW.RecordingStatus = 'Recorded' THEN
        -- Get the debate information
        SELECT * INTO debate_record FROM Debates WHERE DebateID = NEW.DebateID;

        -- Calculate average rank for Team 1
        SELECT AVG(ss.SpeakerRank)::FLOAT INTO avg_rank
        FROM SpeakerScores ss
        JOIN TeamMembers tm ON ss.SpeakerID = tm.StudentID
        WHERE tm.TeamID = debate_record.Team1ID AND ss.BallotID = NEW.BallotID;

        UPDATE TeamScores
        SET Rank = avg_rank
        WHERE TeamID = debate_record.Team1ID AND DebateID = NEW.DebateID;

        -- Calculate average rank for Team 2
        SELECT AVG(ss.SpeakerRank)::FLOAT INTO avg_rank
        FROM SpeakerScores ss
        JOIN TeamMembers tm ON ss.SpeakerID = tm.StudentID
        WHERE tm.TeamID = debate_record.Team2ID AND ss.BallotID = NEW.BallotID;

        UPDATE TeamScores
        SET Rank = avg_rank
        WHERE TeamID = debate_record.Team2ID AND DebateID = NEW.DebateID;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Corrected trigger function for updating team stats
CREATE OR REPLACE FUNCTION update_team_stats()
RETURNS TRIGGER AS $$
DECLARE
    debate_record RECORD;
BEGIN
    IF NEW.RecordingStatus = 'Recorded' THEN
        -- Get the debate information
        SELECT * INTO debate_record FROM Debates WHERE DebateID = NEW.DebateID;

        -- Update stats for Team 1
        WITH team_stats AS (
            SELECT
                t.TeamID,
                t.TournamentID,
                COUNT(CASE WHEN b.Verdict = t.Name THEN 1 ELSE NULL END) AS TotalWins,
                AVG(ts.Rank) AS AvgRank,
                SUM(ts.TotalScore::NUMERIC) AS TotalSpeakerPoints
            FROM
                Teams t
            JOIN
                Debates d ON (t.TeamID = d.Team1ID OR t.TeamID = d.Team2ID) AND t.TournamentID = d.TournamentID
            JOIN
                Ballots b ON d.DebateID = b.DebateID
            JOIN
                TeamScores ts ON t.TeamID = ts.TeamID AND d.DebateID = ts.DebateID
            WHERE
                t.TeamID = debate_record.Team1ID AND t.TournamentID = debate_record.TournamentID
            GROUP BY
                t.TeamID, t.TournamentID
        )
        UPDATE Teams
        SET
            TotalWins = team_stats.TotalWins,
            AverageRank = team_stats.AvgRank,
            TotalSpeakerPoints = team_stats.TotalSpeakerPoints
        FROM
            team_stats
        WHERE
            Teams.TeamID = team_stats.TeamID AND Teams.TournamentID = team_stats.TournamentID;

        -- Update stats for Team 2
        WITH team_stats AS (
            SELECT
                t.TeamID,
                t.TournamentID,
                COUNT(CASE WHEN b.Verdict = t.Name THEN 1 ELSE NULL END) AS TotalWins,
                AVG(ts.Rank) AS AvgRank,
                SUM(ts.TotalScore::NUMERIC) AS TotalSpeakerPoints
            FROM
                Teams t
            JOIN
                Debates d ON (t.TeamID = d.Team1ID OR t.TeamID = d.Team2ID) AND t.TournamentID = d.TournamentID
            JOIN
                Ballots b ON d.DebateID = b.DebateID
            JOIN
                TeamScores ts ON t.TeamID = ts.TeamID AND d.DebateID = ts.DebateID
            WHERE
                t.TeamID = debate_record.Team2ID AND t.TournamentID = debate_record.TournamentID
            GROUP BY
                t.TeamID, t.TournamentID
        )
        UPDATE Teams
        SET
            TotalWins = team_stats.TotalWins,
            AverageRank = team_stats.AvgRank,
            TotalSpeakerPoints = team_stats.TotalSpeakerPoints
        FROM
            team_stats
        WHERE
            Teams.TeamID = team_stats.TeamID AND Teams.TournamentID = team_stats.TournamentID;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Corrected trigger function to handle the special case for "Public Speaking" teams
CREATE OR REPLACE FUNCTION handle_public_speaking_team()
RETURNS TRIGGER AS $$
DECLARE
    debate_record RECORD;
BEGIN
    IF NEW.RecordingStatus = 'Recorded' THEN
        -- Get the debate information
        SELECT * INTO debate_record FROM Debates WHERE DebateID = NEW.DebateID;

        -- Set rank to 99 for "Public Speaking" teams
        UPDATE TeamScores ts
        SET Rank = 99
        WHERE ts.DebateID = NEW.DebateID
          AND ts.TeamID IN (
              SELECT t.TeamID
              FROM Teams t
              WHERE t.TeamID IN (debate_record.Team1ID, debate_record.Team2ID)
                AND t.Name = 'Public Speaking'
          );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to generate the school ID
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

CREATE OR REPLACE FUNCTION generate_idebate_student_id()
RETURNS trigger AS $$
BEGIN
  NEW.iDebateStudentID := 'STUD' || LPAD(NEXTVAL('idebate_student_id_seq')::TEXT, 6, '0');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to create IDs for Volunteers

CREATE OR REPLACE FUNCTION generate_idebate_volunteer_id()
RETURNS trigger AS $$
BEGIN
  NEW.iDebateVolunteerID := 'VOL' || LPAD(NEXTVAL('idebate_volunteer_id_seq')::TEXT, 6, '0');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.UpdatedAt = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';