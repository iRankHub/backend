  -- Drop all tables in reverse order of creation
    DROP TABLE IF EXISTS notification_metadata;
    DROP TABLE IF EXISTS RankingVisibility;
    DROP TABLE IF EXISTS VolunteerRatings;
    DROP TABLE IF EXISTS StudentTransfers;
    DROP TABLE IF EXISTS JudgeReviews;
    DROP TABLE IF EXISTS Communications;
    DROP TABLE IF EXISTS TeamScores;
    DROP TABLE IF EXISTS Schedules;
    DROP TABLE IF EXISTS DebateJudges;
    DROP TABLE IF EXISTS PairingHistory;
    DROP TABLE IF EXISTS SpeakerScores;
    DROP TABLE IF EXISTS Ballots;
    DROP TABLE IF EXISTS JudgeAssignments;
    DROP TABLE IF EXISTS Debates;
    DROP TABLE IF EXISTS Rounds;
    DROP TABLE IF EXISTS TeamMembers;
    DROP TABLE IF EXISTS Teams;
    DROP TABLE IF EXISTS TournamentInvitations;
    DROP TABLE IF EXISTS Volunteers;
    DROP TABLE IF EXISTS Students;
    DROP TABLE IF EXISTS Schools;
    DROP TABLE IF EXISTS Rooms;
    DROP TABLE IF EXISTS Tournaments;
    DROP TABLE IF EXISTS WebAuthnSessionData;
    DROP TABLE IF EXISTS WebAuthnCredentials;
    DROP TABLE IF EXISTS NotificationPreferences;
    DROP TABLE IF EXISTS Notifications;
    DROP TABLE IF EXISTS UserProfiles;
    DROP TABLE IF EXISTS VolunteerRatingTypes;
    DROP TABLE IF EXISTS CountryCodes;
    DROP TABLE IF EXISTS Leagues;
    DROP TABLE IF EXISTS TournamentFormats;
    DROP TABLE IF EXISTS Users;

    -- Drop sequences
    DROP SEQUENCE IF EXISTS idebate_volunteer_id_seq;
    DROP SEQUENCE IF EXISTS idebate_student_id_seq;

    -- Drop functions
    DROP FUNCTION IF EXISTS update_updated_at();
    DROP FUNCTION IF EXISTS update_user_counts();
    DROP FUNCTION IF EXISTS update_tournament_counts();
    DROP FUNCTION IF EXISTS calculate_team_average_rank();
    DROP FUNCTION IF EXISTS update_team_stats();
    DROP FUNCTION IF EXISTS handle_public_speaking_team();
    DROP FUNCTION IF EXISTS generate_idebate_school_id();
    DROP FUNCTION IF EXISTS generate_idebate_volunteer_id();
    DROP FUNCTION IF EXISTS generate_idebate_student_id();