-- Drop existing triggers first
DROP TRIGGER IF EXISTS update_users_updated_at ON Users;
DROP TRIGGER IF EXISTS set_idebate_school_id ON Schools;
DROP TRIGGER IF EXISTS set_idebate_volunteer_id ON Volunteers;
DROP TRIGGER IF EXISTS set_idebate_student_id ON Students;
DROP TRIGGER IF EXISTS calculate_team_average_rank_trigger ON Ballots;
DROP TRIGGER IF EXISTS update_team_stats_trigger ON Ballots;
DROP TRIGGER IF EXISTS handle_public_speaking_team_trigger ON Ballots;
DROP TRIGGER IF EXISTS update_tournament_expenses_modtime ON TournamentExpenses;
DROP TRIGGER IF EXISTS update_school_tournament_registrations_modtime ON SchoolTournamentRegistrations;

-- Create triggers

CREATE TRIGGER update_tournament_expenses_modtime
    BEFORE UPDATE ON TournamentExpenses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_school_tournament_registrations_modtime
    BEFORE UPDATE ON SchoolTournamentRegistrations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON Users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER set_idebate_school_id
    BEFORE INSERT ON Schools
    FOR EACH ROW
    EXECUTE FUNCTION generate_idebate_school_id();

CREATE TRIGGER set_idebate_volunteer_id
    BEFORE INSERT ON Volunteers
    FOR EACH ROW
    EXECUTE FUNCTION generate_idebate_volunteer_id();

CREATE TRIGGER set_idebate_student_id
    BEFORE INSERT ON Students
    FOR EACH ROW
    EXECUTE FUNCTION generate_idebate_student_id();

CREATE TRIGGER calculate_team_average_rank_trigger
    AFTER UPDATE OF RecordingStatus ON Ballots
    FOR EACH ROW
    EXECUTE FUNCTION calculate_team_average_rank();

CREATE TRIGGER update_team_stats_trigger
    AFTER UPDATE OF RecordingStatus ON Ballots
    FOR EACH ROW
    EXECUTE FUNCTION update_team_stats();

CREATE TRIGGER handle_public_speaking_team_trigger
    AFTER UPDATE OF RecordingStatus ON Ballots
    FOR EACH ROW
    EXECUTE FUNCTION handle_public_speaking_team();