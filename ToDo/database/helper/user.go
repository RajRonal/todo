package helper

import (
	"ToDo/database"
	"ToDo/models"
	"database/sql"
	"github.com/lib/pq"
	"time"
)

func LoginUser(username string) (*models.AddLogin, error) {
	SQL := `SElECT id,password from user_details where username=$1`
	var pass models.AddLogin
	err := database.Db.Get(&pass, SQL, username)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &pass, nil
}

func CreateUser(user models.CreateUser) (string, error) {
	SQL := `INSERT INTO user_details (name, email,username,password) 	
			VALUES ($1, $2,$3,$4) 
			returning id;`
	var userID string
	err := database.Db.Get(&userID, SQL, user.Name, user.Email, user.Username, user.Password)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func AddTask(task, id string, status models.Status) (string, error) {
	SQL := `INSERT INTO todo (task,Status, id) 	
			VALUES ($1, $2,$3)  
			returning id;`
	var userID string
	err := database.Db.Get(&userID, SQL, task, status, id)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func CreateSession(ID string, ExpiredAt time.Time) (string, error) {
	SQL := `INSERT INTO sessions (id,expired_at ) 	
			VALUES ($1, $2)  
			returning session_id;`
	var sessId string
	err := database.Db.Get(&sessId, SQL, ID, ExpiredAt)
	if err != nil {
		return "", err
	}
	return sessId, nil
}

func SessionExist(sessionID string) (bool, error) {

	var isExpired bool
	query := `SELECT count(*) > 0
			  FROM sessions 
			WHERE session_id=$1  and expired_at >now() and archived_at is null`
	checkSessionErr := database.Db.Get(&isExpired, query, sessionID)
	if checkSessionErr != nil {
		return isExpired, checkSessionErr
	}
	return isExpired, nil
}

func DeleteSession(sessionID string) error {
	currentTime := time.Now()
	sql := `UPDATE sessions
			  SET archived_at=$1,expired_at=now()
			  WHERE session_id=$2`
	_, err := database.Db.Exec(sql, currentTime, sessionID)

	return err
}

func DeleteTask(taskid string) error {
	sql := `UPDATE todo
			  SET archived_at=now()
			  WHERE task_id=$1`
	_, err := database.Db.Exec(sql, taskid)
	if err != nil {
		return err
	}

	return nil
}

func GetAllTask(userid, searchText string, status pq.StringArray, isStatus bool, pageNo, taskSize int) (models.PaginatedTask, error) {
	var data models.PaginatedTask
	SQL := `WITH userTask AS (SELECT  count(*) over () total_count,task_id, task, status
			FROM todo
			WHERE id = $1 
			  AND task ILIKE '%' || $2 ||'%'  
			  AND ($4 OR status= ANY($3))
			  AND archived_at is null)
			 SELECT  total_count,task_id ,task, status from userTask 
			        LIMIT $5
					OFFSET $6
					`
	user := make([]models.PaginatedData, 0)
	err := database.Db.Select(&user, SQL, userid, searchText, status, !isStatus, taskSize, pageNo*taskSize)
	if len(user) == 0 {
		return data, err
	}
	data.TotalCount = user[0].TotalCount
	data.Data = user
	return data, err

}

func UpdateTasks(taskid, task string) error {
	SQL := `UPDATE todo
			  SET task=$1
			  WHERE task_id=$2 and archived_at is null`
	_, err := database.Db.Exec(SQL, task, taskid)
	if err != nil {
		return err
	}
	return nil
}

func MarkComplete(taskID string) error {

	SQL := `update todo
            set is_complete=$1
            where task_id=$2`
	_, err := database.Db.Exec(SQL, true, taskID)
	return err
}
