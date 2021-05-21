package main

import (
	"database/sql"
	"errors"
	"strings"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
)

type Course struct {
	Id          string `json:"Courseid"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Status      string `json:"Status"`
}

func (c *Course) createCourse() error {

	results, err := myDb.Exec("INSERT INTO course VALUES (?,?,?,?)", c.Id, c.Title, c.Description, c.Status)

	if err != nil {
		return err
	} else {
		rows, _ := results.RowsAffected()

		if rows > 0 {
			return nil
		}
	}

	return errors.New("unable to create course")
}

func getCourse(id string) (*Course, error) {

	var c Course
	err := myDb.QueryRow("SELECT * FROM course WHERE id = ? LIMIT 1", id).Scan(&c.Id, &c.Title, &c.Description, &c.Status)
	switch {
	case err == sql.ErrNoRows:
		return nil, errCourseNotFound
	case err != nil:
		return nil, err
	default:
		return &c, nil
	}
}

func getCourses(status string) ([]*Course, error) {

	var sqlQuery string

	switch {
	case status == "inactive":
		sqlQuery = "Select * FROM course WHERE status = 'inactive' ORDER BY Id ASC"
	case status == "active":
		sqlQuery = "Select * FROM course WHERE status = 'active' ORDER BY Id ASC"
	default:
		sqlQuery = "Select * FROM course ORDER BY Id ASC"
	}

	results, err := myDb.Query(sqlQuery)

	if err != nil {
		return nil, err
	}

	var courses = []*Course{}

	for results.Next() {
		// map this type to the record in the table
		var c Course

		err = results.Scan(&c.Id, &c.Title, &c.Description, &c.Status)
		if err != nil {
			return nil, err
		}

		courses = append(courses, &c)
	}

	return courses, nil
}

func (c *Course) deleteCourse() error {
	results, err := myDb.Exec("DELETE FROM course WHERE ID=?", c.Id)

	if err != nil {
		return err
	} else {
		rows, _ := results.RowsAffected()

		if rows > 0 {
			return nil
		}
	}

	return errors.New("unable to delete course")
}

func (c *Course) updateCourse() error {
	results, err := myDb.Exec("UPDATE course SET Title=?, Description=?, Status=? WHERE Id=?", c.Title, c.Description, c.Status, c.Id)

	if err != nil {
		return err
	} else {
		rows, _ := results.RowsAffected()

		if rows > 0 {
			return nil
		}
	}

	return errors.New("unable to update course")
}

func (c *Course) validateFields() error {

	fErr := c.fieldsCheck()

	if fErr != nil {
		return fErr
	}

	sErr := c.dbSchemaCheck()

	if sErr != nil {
		return sErr
	}

	return nil
}

func (c *Course) fieldsCheck() error {

	if strings.TrimSpace(c.Id) == "" || strings.TrimSpace(c.Title) == "" || strings.TrimSpace(c.Description) == "" || strings.TrimSpace(c.Status) == "" {
		return errInvalidCourseInfo
	}

	if strings.ToLower(c.Status) != "active" && strings.ToLower(c.Status) != "inactive" {
		return errInvalidCourseStatus
	}

	return nil
}

func (c *Course) dbSchemaCheck() error {

	if utf8.RuneCountInString(c.Id) > 20 {
		return errInvalidCourseIdLength
	}

	if utf8.RuneCountInString(c.Title) > 255 {
		return errInvalidCourseTitleLength
	}

	if utf8.RuneCountInString(c.Description) > 255 {
		return errInvalidCourseDescriptionLength
	}

	return nil
}
