package mysql

import (
	"context"
	"errors"
	"testing"

	mooc "github.com/CodelyTV/go-hexagonal_http_api-course/03-02-repository-test/internal"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CourseRepository_Save_RepositoryError(t *testing.T) {
	courseID, courseName, courseDuration := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test Course", "10 months"
	course, err := mooc.NewCourse(courseID, courseName, courseDuration)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO courses (id, name, duration) VALUES (?, ?, ?)").
		WithArgs(courseID, courseName, courseDuration).
		WillReturnError(errors.New("something-failed"))

	repo := NewCourseRepository(db)

	err = repo.Save(context.Background(), course)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
}

func Test_CourseRepository_Save_Succeed(t *testing.T) {
	courseID, courseName, courseDuration := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test Course", "10 months"
	course, err := mooc.NewCourse(courseID, courseName, courseDuration)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO courses (id, name, duration) VALUES (?, ?, ?)").
		WithArgs(courseID, courseName, courseDuration).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := NewCourseRepository(db)

	err = repo.Save(context.Background(), course)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
}


func Test_CourseRepository_Get_Succeed(t *testing.T) {
	courseID, courseName, courseDuration := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test Course", "10 months"
	course, err := mooc.NewCourse(courseID, courseName, courseDuration)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)

	sqlMock.ExpectQuery(`SELECT courses\.id, courses\.name, courses\.duration FROM courses WHERE id = \?`).
		WithArgs(courseID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "duration"}).
		AddRow("37a0f027-15e6-47cc-a5d2-64183281087e", "Test Course", "10 months"))

	repo := NewCourseRepository(db)

	cID, err := mooc.NewCourseID(courseID)
	require.NoError(t, err)

	rsC, err := repo.GetCourse(context.Background(), cID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Equal(t, course, rsC)
}