package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/package/models"
	"log"
)

func AddNotificationToDB(db *sql.DB, postId, userId int, comId int64) error {
	stmt, err := db.Prepare("INSERT INTO NotificationComments(UserId,PostId,PostAuthorId,CommentId)VALUES(?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		return err
	}
	post, err := GetOnePost(db, postId)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userId, postId, post.UserId, comId)
	if err != nil {
		// fmt.Println(err)
		return err
	}
	return nil
}

func GetNotificationComments(db *sql.DB, userId int) ([]models.NotificationComments, error) {
	NotifComments := []models.NotificationComments{}
	row, err := db.Query("SELECT UserId,PostId,PostAuthorId,CommentId FROM NotificationComments WHERE PostAuthorId = ? ORDER BY Id DESC", userId)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		NotifComment := models.NotificationComments{}
		err = row.Scan(&NotifComment.UserId, &NotifComment.PostId, &NotifComment.PostAuthorId, &NotifComment.CommentId)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		NotifComments = append(NotifComments, NotifComment)
	}
	return NotifComments, nil
}

func DeleteOneNotificationComment(db *sql.DB, comId, PostAuthorId int) error {
	stmt, err := db.Prepare("DELETE FROM NotificationComments WHERE PostAuthorId = ? AND CommentId = ?")
	if err != nil {
		return err
	}
	row, err := stmt.Exec(PostAuthorId, comId)
	if err != nil {
		return err
	}
	res, _ := row.RowsAffected()
	if res == 0 {
		return errors.New("not found")
	}
	return nil
}

func AddRateNotificationToDB(db *sql.DB, postId, UserId int, rate string) error {
	post, err := GetOnePost(db, postId)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO PostRateNotification(UserId,PostId,PostAuthorId,Rate)VALUES(?,?,?,?)")
	_, err = stmt.Exec(UserId, postId, post.UserId, rate)
	if err != nil {
		return err
	}
	return nil
}

func GetRateNotification(db *sql.DB, UserId int) ([]models.RateNotification, error) {
	NotifRate := []models.RateNotification{}
	row, err := db.Query("SELECT Id,UserId,PostId,PostAuthorId,Rate FROM PostRateNotification WHERE PostAuthorId = ? ORDER BY Id DESC", UserId)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		notif := models.RateNotification{}
		err = row.Scan(&notif.Id, &notif.UserId, &notif.PostId, &notif.PostAuthorId, &notif.Rate)
		if err != nil {
			return nil, err
		}
		NotifRate = append(NotifRate, notif)
	}
	return NotifRate, nil
}

func DeleteRateNotification(db *sql.DB, UserId, Id int) error {
	stmt, err := db.Prepare("DELETE FROM PostRateNotification WHERE PostAuthorId = ? AND Id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(UserId, Id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllNotification(db *sql.DB, userId int) {
	stmt, err := db.Prepare("DELETE FROM NotificationComments WHERE PostAuthorId = ?")
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(userId)
	if err != nil {
		log.Println(err)
	}
	stmt1, err := db.Prepare("DELETE FROM PostRateNotification WHERE PostAuthorId = ?")
	if err != nil {
		log.Println(err)
	}
	_, err = stmt1.Exec(userId)
	if err != nil {
		log.Println(err)
	}
}

func DeleteNotificationIfPostDeleted(db *sql.DB, postId int) {
	stmt, err := db.Prepare("DELETE FROM NotificationComments WHERE PostId = ?")
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(postId)
	if err != nil {
		log.Println(err)
	}
	stmt1, err := db.Prepare("DELETE FROM PostRateNotification WHERE PostId = ?")
	if err != nil {
		log.Println(err)
	}
	_, err = stmt1.Exec(postId)
	if err != nil {
		log.Println(err)
	}
}

func DeleteNotificationIfCommentDeleted(db *sql.DB, comId int) {
	stmt, err := db.Prepare("DELETE FROM NotificationComments WHERE CommentId = ?")
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(comId)
	if err != nil {
		log.Println(err)
	}
}
