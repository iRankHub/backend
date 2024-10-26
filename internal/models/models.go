// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Ballot struct {
	Ballotid           int32          `json:"ballotid"`
	Debateid           int32          `json:"debateid"`
	Judgeid            int32          `json:"judgeid"`
	Team1totalscore    sql.NullString `json:"team1totalscore"`
	Team1feedback      sql.NullString `json:"team1feedback"`
	Team2totalscore    sql.NullString `json:"team2totalscore"`
	Team2feedback      sql.NullString `json:"team2feedback"`
	Recordingstatus    string         `json:"recordingstatus"`
	Verdict            string         `json:"verdict"`
	LastUpdatedBy      sql.NullInt32  `json:"last_updated_by"`
	LastUpdatedAt      sql.NullTime   `json:"last_updated_at"`
	HeadJudgeSubmitted sql.NullBool   `json:"head_judge_submitted"`
}

type Communication struct {
	Communicationid int32     `json:"communicationid"`
	Userid          int32     `json:"userid"`
	Schoolid        int32     `json:"schoolid"`
	Type            string    `json:"type"`
	Content         string    `json:"content"`
	Timestamp       time.Time `json:"timestamp"`
}

type Countrycode struct {
	Countryname string `json:"countryname"`
	Isocode     string `json:"isocode"`
}

type Debate struct {
	Debateid           int32        `json:"debateid"`
	Roundid            int32        `json:"roundid"`
	Roundnumber        int32        `json:"roundnumber"`
	Iseliminationround bool         `json:"iseliminationround"`
	Tournamentid       int32        `json:"tournamentid"`
	Team1id            int32        `json:"team1id"`
	Team2id            int32        `json:"team2id"`
	Starttime          time.Time    `json:"starttime"`
	Endtime            sql.NullTime `json:"endtime"`
	Roomid             int32        `json:"roomid"`
	Status             string       `json:"status"`
}

type Debatejudge struct {
	Debateid int32 `json:"debateid"`
	Judgeid  int32 `json:"judgeid"`
}

type Judgeassignment struct {
	Assignmentid  int32 `json:"assignmentid"`
	Tournamentid  int32 `json:"tournamentid"`
	Judgeid       int32 `json:"judgeid"`
	Debateid      int32 `json:"debateid"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
	Isheadjudge   bool  `json:"isheadjudge"`
}

type Judgefeedback struct {
	Feedbackid             int32          `json:"feedbackid"`
	Judgeid                sql.NullInt32  `json:"judgeid"`
	Studentid              sql.NullInt32  `json:"studentid"`
	Debateid               sql.NullInt32  `json:"debateid"`
	Clarityrating          sql.NullString `json:"clarityrating"`
	Constructivenessrating sql.NullString `json:"constructivenessrating"`
	Timelinessrating       sql.NullString `json:"timelinessrating"`
	Fairnessrating         sql.NullString `json:"fairnessrating"`
	Engagementrating       sql.NullString `json:"engagementrating"`
	Averagerating          sql.NullString `json:"averagerating"`
	Textfeedback           sql.NullString `json:"textfeedback"`
	Isread                 sql.NullBool   `json:"isread"`
	Createdat              sql.NullTime   `json:"createdat"`
}

type League struct {
	Leagueid   int32           `json:"leagueid"`
	Name       string          `json:"name"`
	Leaguetype string          `json:"leaguetype"`
	Details    json.RawMessage `json:"details"`
	DeletedAt  sql.NullTime    `json:"deleted_at"`
}

type Notification struct {
	Notificationid int32          `json:"notificationid"`
	Userid         int32          `json:"userid"`
	Type           string         `json:"type"`
	Message        string         `json:"message"`
	Recipientemail sql.NullString `json:"recipientemail"`
	Subject        sql.NullString `json:"subject"`
	Isread         sql.NullBool   `json:"isread"`
	Createdat      sql.NullTime   `json:"createdat"`
}

type Notificationpreference struct {
	Preferenceid       int32          `json:"preferenceid"`
	Userid             int32          `json:"userid"`
	Emailnotifications sql.NullBool   `json:"emailnotifications"`
	Emailfrequency     sql.NullString `json:"emailfrequency"`
	Emailday           sql.NullInt32  `json:"emailday"`
	Emailtime          sql.NullTime   `json:"emailtime"`
	Inappnotifications sql.NullBool   `json:"inappnotifications"`
}

type Pairinghistory struct {
	Historyid     int32 `json:"historyid"`
	Tournamentid  int32 `json:"tournamentid"`
	Team1id       int32 `json:"team1id"`
	Team2id       int32 `json:"team2id"`
	Roundnumber   int32 `json:"roundnumber"`
	Iselimination bool  `json:"iselimination"`
}

type Room struct {
	Roomid       int32         `json:"roomid"`
	Roomname     string        `json:"roomname"`
	Location     string        `json:"location"`
	Capacity     int32         `json:"capacity"`
	Tournamentid sql.NullInt32 `json:"tournamentid"`
}

type Round struct {
	Roundid            int32 `json:"roundid"`
	Tournamentid       int32 `json:"tournamentid"`
	Roundnumber        int32 `json:"roundnumber"`
	Iseliminationround bool  `json:"iseliminationround"`
}

type Schedule struct {
	Scheduleid    int32     `json:"scheduleid"`
	Tournamentid  int32     `json:"tournamentid"`
	Roundid       int32     `json:"roundid"`
	Debateid      int32     `json:"debateid"`
	Scheduledtime time.Time `json:"scheduledtime"`
}

type School struct {
	Schoolid                int32          `json:"schoolid"`
	Idebateschoolid         sql.NullString `json:"idebateschoolid"`
	Schoolname              string         `json:"schoolname"`
	Address                 string         `json:"address"`
	Country                 sql.NullString `json:"country"`
	Province                sql.NullString `json:"province"`
	District                sql.NullString `json:"district"`
	Contactpersonid         int32          `json:"contactpersonid"`
	Contactpersonnationalid sql.NullString `json:"contactpersonnationalid"`
	Contactemail            string         `json:"contactemail"`
	Schoolemail             string         `json:"schoolemail"`
	Schooltype              string         `json:"schooltype"`
}

type Schooltournamentregistration struct {
	Registrationid    int32          `json:"registrationid"`
	Schoolid          int32          `json:"schoolid"`
	Tournamentid      int32          `json:"tournamentid"`
	Plannedteamscount int32          `json:"plannedteamscount"`
	Actualteamscount  sql.NullInt32  `json:"actualteamscount"`
	Amountperteam     string         `json:"amountperteam"`
	Totalamount       sql.NullString `json:"totalamount"`
	Discountamount    sql.NullString `json:"discountamount"`
	Actualpaidamount  sql.NullString `json:"actualpaidamount"`
	Paymentstatus     string         `json:"paymentstatus"`
	Paymentdate       sql.NullTime   `json:"paymentdate"`
	Currency          string         `json:"currency"`
	Createdat         sql.NullTime   `json:"createdat"`
	Updatedat         sql.NullTime   `json:"updatedat"`
	Createdby         sql.NullInt32  `json:"createdby"`
	Updatedby         sql.NullInt32  `json:"updatedby"`
}

type Speakerscore struct {
	Scoreid       int32          `json:"scoreid"`
	Ballotid      int32          `json:"ballotid"`
	Speakerid     int32          `json:"speakerid"`
	Speakerrank   int32          `json:"speakerrank"`
	Speakerpoints string         `json:"speakerpoints"`
	Feedback      sql.NullString `json:"feedback"`
	Isread        sql.NullBool   `json:"isread"`
}

type Student struct {
	Studentid        int32          `json:"studentid"`
	Idebatestudentid sql.NullString `json:"idebatestudentid"`
	Firstname        string         `json:"firstname"`
	Lastname         string         `json:"lastname"`
	Gender           sql.NullString `json:"gender"`
	Grade            string         `json:"grade"`
	Dateofbirth      sql.NullTime   `json:"dateofbirth"`
	Email            sql.NullString `json:"email"`
	Password         string         `json:"password"`
	Schoolid         int32          `json:"schoolid"`
	Userid           int32          `json:"userid"`
}

type Studenttransfer struct {
	Transferid   int32          `json:"transferid"`
	Studentid    int32          `json:"studentid"`
	Fromschoolid sql.NullInt32  `json:"fromschoolid"`
	Toschoolid   int32          `json:"toschoolid"`
	Transferdate time.Time      `json:"transferdate"`
	Reason       sql.NullString `json:"reason"`
}

type Team struct {
	Teamid             int32          `json:"teamid"`
	Name               string         `json:"name"`
	Tournamentid       int32          `json:"tournamentid"`
	Totalwins          sql.NullInt32  `json:"totalwins"`
	Totalspeakerpoints sql.NullString `json:"totalspeakerpoints"`
	Averagerank        sql.NullString `json:"averagerank"`
}

type Teammember struct {
	Teamid    int32 `json:"teamid"`
	Studentid int32 `json:"studentid"`
}

type Teamscore struct {
	Scoreid       int32          `json:"scoreid"`
	Teamid        sql.NullInt32  `json:"teamid"`
	Debateid      sql.NullInt32  `json:"debateid"`
	Totalscore    sql.NullString `json:"totalscore"`
	Rank          sql.NullInt32  `json:"rank"`
	Iselimination sql.NullBool   `json:"iselimination"`
}

type Tournament struct {
	Tournamentid                 int32          `json:"tournamentid"`
	Name                         string         `json:"name"`
	Startdate                    time.Time      `json:"startdate"`
	Enddate                      time.Time      `json:"enddate"`
	Location                     string         `json:"location"`
	Formatid                     int32          `json:"formatid"`
	Leagueid                     sql.NullInt32  `json:"leagueid"`
	Coordinatorid                int32          `json:"coordinatorid"`
	Numberofpreliminaryrounds    int32          `json:"numberofpreliminaryrounds"`
	Numberofeliminationrounds    int32          `json:"numberofeliminationrounds"`
	Judgesperdebatepreliminary   int32          `json:"judgesperdebatepreliminary"`
	Judgesperdebateelimination   int32          `json:"judgesperdebateelimination"`
	Tournamentfee                string         `json:"tournamentfee"`
	Imageurl                     sql.NullString `json:"imageurl"`
	CreatedAt                    sql.NullTime   `json:"created_at"`
	UpdatedAt                    sql.NullTime   `json:"updated_at"`
	DeletedAt                    sql.NullTime   `json:"deleted_at"`
	YesterdayTotalCount          sql.NullInt32  `json:"yesterday_total_count"`
	YesterdayUpcomingCount       sql.NullInt32  `json:"yesterday_upcoming_count"`
	YesterdayActiveDebatersCount sql.NullInt32  `json:"yesterday_active_debaters_count"`
}

type Tournamentexpense struct {
	Expenseid         int32          `json:"expenseid"`
	Tournamentid      int32          `json:"tournamentid"`
	Foodexpense       string         `json:"foodexpense"`
	Transportexpense  string         `json:"transportexpense"`
	Perdiemexpense    string         `json:"perdiemexpense"`
	Awardingexpense   string         `json:"awardingexpense"`
	Stationaryexpense string         `json:"stationaryexpense"`
	Otherexpenses     string         `json:"otherexpenses"`
	Totalexpense      sql.NullString `json:"totalexpense"`
	Currency          string         `json:"currency"`
	Notes             sql.NullString `json:"notes"`
	Createdat         sql.NullTime   `json:"createdat"`
	Updatedat         sql.NullTime   `json:"updatedat"`
	Createdby         sql.NullInt32  `json:"createdby"`
	Updatedby         sql.NullInt32  `json:"updatedby"`
}

type Tournamentformat struct {
	Formatid        int32          `json:"formatid"`
	Formatname      string         `json:"formatname"`
	Description     sql.NullString `json:"description"`
	Speakersperteam int32          `json:"speakersperteam"`
	DeletedAt       sql.NullTime   `json:"deleted_at"`
}

type Tournamentinvitation struct {
	Invitationid   int32        `json:"invitationid"`
	Tournamentid   int32        `json:"tournamentid"`
	Inviteeid      string       `json:"inviteeid"`
	Inviteerole    string       `json:"inviteerole"`
	Status         string       `json:"status"`
	CreatedAt      sql.NullTime `json:"created_at"`
	UpdatedAt      sql.NullTime `json:"updated_at"`
	Remindersentat sql.NullTime `json:"remindersentat"`
}

type User struct {
	Userid                 int32          `json:"userid"`
	Webauthnuserid         []byte         `json:"webauthnuserid"`
	Name                   string         `json:"name"`
	Gender                 sql.NullString `json:"gender"`
	Email                  string         `json:"email"`
	Password               string         `json:"password"`
	Userrole               string         `json:"userrole"`
	Status                 sql.NullString `json:"status"`
	Verificationstatus     sql.NullBool   `json:"verificationstatus"`
	Deactivatedat          sql.NullTime   `json:"deactivatedat"`
	TwoFactorSecret        sql.NullString `json:"two_factor_secret"`
	TwoFactorEnabled       sql.NullBool   `json:"two_factor_enabled"`
	FailedLoginAttempts    sql.NullInt32  `json:"failed_login_attempts"`
	LastLoginAttempt       sql.NullTime   `json:"last_login_attempt"`
	LastLogout             sql.NullTime   `json:"last_logout"`
	ResetToken             sql.NullString `json:"reset_token"`
	ResetTokenExpires      sql.NullTime   `json:"reset_token_expires"`
	CreatedAt              sql.NullTime   `json:"created_at"`
	UpdatedAt              sql.NullTime   `json:"updated_at"`
	DeletedAt              sql.NullTime   `json:"deleted_at"`
	YesterdayApprovedCount sql.NullInt32  `json:"yesterday_approved_count"`
}

type Userprofile struct {
	Profileid          int32          `json:"profileid"`
	Userid             int32          `json:"userid"`
	Name               string         `json:"name"`
	Userrole           string         `json:"userrole"`
	Email              string         `json:"email"`
	Password           string         `json:"password"`
	Gender             sql.NullString `json:"gender"`
	Address            sql.NullString `json:"address"`
	Phone              sql.NullString `json:"phone"`
	Bio                sql.NullString `json:"bio"`
	Profilepicture     sql.NullString `json:"profilepicture"`
	Verificationstatus sql.NullBool   `json:"verificationstatus"`
}

type Volunteer struct {
	Volunteerid                  int32          `json:"volunteerid"`
	Idebatevolunteerid           sql.NullString `json:"idebatevolunteerid"`
	Firstname                    string         `json:"firstname"`
	Lastname                     string         `json:"lastname"`
	Gender                       sql.NullString `json:"gender"`
	Dateofbirth                  sql.NullTime   `json:"dateofbirth"`
	Nationalid                   sql.NullString `json:"nationalid"`
	Role                         string         `json:"role"`
	Graduateyear                 sql.NullInt32  `json:"graduateyear"`
	Password                     string         `json:"password"`
	Safeguardcertificate         sql.NullString `json:"safeguardcertificate"`
	Hasinternship                sql.NullBool   `json:"hasinternship"`
	Isenrolledinuniversity       sql.NullBool   `json:"isenrolledinuniversity"`
	Userid                       int32          `json:"userid"`
	YesterdayRoundsJudged        sql.NullInt32  `json:"yesterday_rounds_judged"`
	YesterdayTournamentsAttended sql.NullInt32  `json:"yesterday_tournaments_attended"`
	YesterdayUpcomingTournaments sql.NullInt32  `json:"yesterday_upcoming_tournaments"`
}

type Volunteerrating struct {
	Ratingid         int32          `json:"ratingid"`
	Volunteerid      int32          `json:"volunteerid"`
	Ratingtypeid     int32          `json:"ratingtypeid"`
	Ratingscore      string         `json:"ratingscore"`
	Ratingcomments   sql.NullString `json:"ratingcomments"`
	Cumulativerating sql.NullString `json:"cumulativerating"`
}

type Volunteerratingtype struct {
	Ratingtypeid int32          `json:"ratingtypeid"`
	Category     sql.NullString `json:"category"`
}

type Webauthncredential struct {
	ID              int32        `json:"id"`
	Userid          int32        `json:"userid"`
	Credentialid    []byte       `json:"credentialid"`
	Publickey       []byte       `json:"publickey"`
	Attestationtype string       `json:"attestationtype"`
	Aaguid          []byte       `json:"aaguid"`
	Signcount       int64        `json:"signcount"`
	Createdat       sql.NullTime `json:"createdat"`
}

type Webauthnsessiondatum struct {
	Userid      int32  `json:"userid"`
	Sessiondata []byte `json:"sessiondata"`
}
