// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package models

import (
	"database/sql"
	"time"
)

type Ballot struct {
	Ballotid              int32          `json:"ballotid"`
	Debateid              int32          `json:"debateid"`
	Judgeid               int32          `json:"judgeid"`
	Team1debaterascore    sql.NullString `json:"team1debaterascore"`
	Team1debateracomments sql.NullString `json:"team1debateracomments"`
	Team1debaterbscore    sql.NullString `json:"team1debaterbscore"`
	Team1debaterbcomments sql.NullString `json:"team1debaterbcomments"`
	Team1debatercscore    sql.NullString `json:"team1debatercscore"`
	Team1debaterccomments sql.NullString `json:"team1debaterccomments"`
	Team1totalscore       sql.NullString `json:"team1totalscore"`
	Team2debaterascore    sql.NullString `json:"team2debaterascore"`
	Team2debateracomments sql.NullString `json:"team2debateracomments"`
	Team2debaterbscore    sql.NullString `json:"team2debaterbscore"`
	Team2debaterbcomments sql.NullString `json:"team2debaterbcomments"`
	Team2debatercscore    sql.NullString `json:"team2debatercscore"`
	Team2debaterccomments sql.NullString `json:"team2debaterccomments"`
	Team2totalscore       sql.NullString `json:"team2totalscore"`
}

type Communication struct {
	Communicationid int32     `json:"communicationid"`
	Userid          int32     `json:"userid"`
	Schoolid        int32     `json:"schoolid"`
	Type            string    `json:"type"`
	Content         string    `json:"content"`
	Timestamp       time.Time `json:"timestamp"`
}

type Debate struct {
	Debateid     int32        `json:"debateid"`
	Roundid      int32        `json:"roundid"`
	Tournamentid int32        `json:"tournamentid"`
	Team1id      int32        `json:"team1id"`
	Team2id      int32        `json:"team2id"`
	Starttime    time.Time    `json:"starttime"`
	Endtime      sql.NullTime `json:"endtime"`
	Roomid       int32        `json:"roomid"`
	Status       string       `json:"status"`
}

type Internationalleaguedetail struct {
	Leagueid  int32          `json:"leagueid"`
	Continent sql.NullString `json:"continent"`
	Country   sql.NullString `json:"country"`
}

type Judgeassignment struct {
	Assignmentid int32 `json:"assignmentid"`
	Volunteerid  int32 `json:"volunteerid"`
	Tournamentid int32 `json:"tournamentid"`
	Debateid     int32 `json:"debateid"`
}

type Judgereview struct {
	Reviewid  int32          `json:"reviewid"`
	Studentid int32          `json:"studentid"`
	Judgeid   int32          `json:"judgeid"`
	Rating    string         `json:"rating"`
	Comments  sql.NullString `json:"comments"`
}

type League struct {
	Leagueid   int32  `json:"leagueid"`
	Name       string `json:"name"`
	Leaguetype string `json:"leaguetype"`
}

type Localleaguedetail struct {
	Leagueid int32          `json:"leagueid"`
	Province sql.NullString `json:"province"`
	District sql.NullString `json:"district"`
}

type Result struct {
	Resultid     int32          `json:"resultid"`
	Tournamentid int32          `json:"tournamentid"`
	Teamid       int32          `json:"teamid"`
	Rank         sql.NullInt32  `json:"rank"`
	Points       sql.NullString `json:"points"`
}

type Room struct {
	Roomid   int32  `json:"roomid"`
	Roomname string `json:"roomname"`
	Location string `json:"location"`
	Capacity int32  `json:"capacity"`
}

type Roombooking struct {
	Bookingid    int32     `json:"bookingid"`
	Tournamentid int32     `json:"tournamentid"`
	Roomid       int32     `json:"roomid"`
	Starttime    time.Time `json:"starttime"`
	Endtime      time.Time `json:"endtime"`
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
	Schoolid        int32          `json:"schoolid"`
	Schoolname      string         `json:"schoolname"`
	Address         string         `json:"address"`
	Country         sql.NullString `json:"country"`
	Province        sql.NullString `json:"province"`
	District        sql.NullString `json:"district"`
	Contactpersonid int32          `json:"contactpersonid"`
	Contactemail    string         `json:"contactemail"`
	Schooltype      string         `json:"schooltype"`
}

type Student struct {
	Studentid   int32          `json:"studentid"`
	Firstname   string         `json:"firstname"`
	Lastname    string         `json:"lastname"`
	Grade       string         `json:"grade"`
	Dateofbirth sql.NullTime   `json:"dateofbirth"`
	Email       sql.NullString `json:"email"`
	Password    string         `json:"password"`
	Schoolid    int32          `json:"schoolid"`
	Userid      int32          `json:"userid"`
}

type Studentrank struct {
	Rankid       int32          `json:"rankid"`
	Studentid    int32          `json:"studentid"`
	Tournamentid int32          `json:"tournamentid"`
	Rankvalue    string         `json:"rankvalue"`
	Rankcomments sql.NullString `json:"rankcomments"`
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
	Teamid       int32  `json:"teamid"`
	Name         string `json:"name"`
	Schoolid     int32  `json:"schoolid"`
	Tournamentid int32  `json:"tournamentid"`
}

type Teammember struct {
	Teamid    int32 `json:"teamid"`
	Studentid int32 `json:"studentid"`
}

type Tournament struct {
	Tournamentid int32         `json:"tournamentid"`
	Name         string        `json:"name"`
	Startdate    time.Time     `json:"startdate"`
	Enddate      time.Time     `json:"enddate"`
	Location     string        `json:"location"`
	Formatid     int32         `json:"formatid"`
	Leagueid     sql.NullInt32 `json:"leagueid"`
}

type Tournamentcoordinator struct {
	Coordinatorid int32  `json:"coordinatorid"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Userid        int32  `json:"userid"`
}

type Tournamentformat struct {
	Formatid    int32          `json:"formatid"`
	Formatname  string         `json:"formatname"`
	Description sql.NullString `json:"description"`
}

type User struct {
	Userid              int32          `json:"userid"`
	Name                string         `json:"name"`
	Email               string         `json:"email"`
	Password            string         `json:"password"`
	Userrole            string         `json:"userrole"`
	Verificationstatus  sql.NullBool   `json:"verificationstatus"`
	Approvalstatus      sql.NullBool   `json:"approvalstatus"`
	TwoFactorSecret     sql.NullString `json:"two_factor_secret"`
	TwoFactorEnabled    sql.NullBool   `json:"two_factor_enabled"`
	FailedLoginAttempts sql.NullInt32  `json:"failed_login_attempts"`
	LastLoginAttempt    sql.NullTime   `json:"last_login_attempt"`
	ResetToken          sql.NullString `json:"reset_token"`
	ResetTokenExpires   sql.NullTime   `json:"reset_token_expires"`
	BiometricToken      sql.NullString `json:"biometric_token"`
}

type Userprofile struct {
	Profileid      int32          `json:"profileid"`
	Userid         int32          `json:"userid"`
	Address        sql.NullString `json:"address"`
	Phone          sql.NullString `json:"phone"`
	Bio            sql.NullString `json:"bio"`
	Profilepicture []byte         `json:"profilepicture"`
}

type Volunteer struct {
	Volunteerid          int32         `json:"volunteerid"`
	Firstname            string        `json:"firstname"`
	Lastname             string        `json:"lastname"`
	Dateofbirth          sql.NullTime  `json:"dateofbirth"`
	Role                 string        `json:"role"`
	Graduateyear         sql.NullInt32 `json:"graduateyear"`
	Password             string        `json:"password"`
	Safeguardcertificate sql.NullBool  `json:"safeguardcertificate"`
	Userid               int32         `json:"userid"`
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
	Ratingtypeid   int32  `json:"ratingtypeid"`
	Ratingtypename string `json:"ratingtypename"`
}
