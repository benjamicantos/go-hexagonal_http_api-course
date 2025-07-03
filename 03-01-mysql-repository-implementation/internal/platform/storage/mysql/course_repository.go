package mysql

import (
	"context"
	"database/sql"
	"fmt"

	mooc "github.com/CodelyTV/go-hexagonal_http_api-course/03-01-mysql-repository-implementation/internal"
	"github.com/huandu/go-sqlbuilder"
)

// CourseRepository is a MySQL mooc.CourseRepository implementation.
type CourseRepository struct {
	db *sql.DB
}

// NewCourseRepository initializes a MySQL-based implementation of mooc.CourseRepository.
func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

// Save implements the mooc.CourseRepository interface.
func (r *CourseRepository) Save(ctx context.Context, course mooc.Course) error {
	courseSQLStruct := sqlbuilder.NewStruct(new(sqlCourse))
	query, args := courseSQLStruct.InsertInto(sqlCourseTable, sqlCourse{
		ID:       course.ID().String(),
		Name:     course.Name().String(),
		Duration: course.Duration().String(),
	}).Build()

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to persist course on database: %v", err)
	}

	return nil
}

// Get implements the mooc.CourseRepository interface.
func (r *CourseRepository) Get(ctx context.Context, id string) (mooc.Course, error) {
	courseSQLStruct := sqlbuilder.NewStruct(new(sqlCourse))
	query, args := courseSQLStruct.SelectFrom(sqlCourseTable).Where("id = ?", id).Build()

	row := r.db.QueryRowContext(ctx, query, args...)
	var sc sqlCourse
	if err := row.Scan(&sc.ID, &sc.Name, &sc.Duration); err != nil {
		if err == sql.ErrNoRows {
			return mooc.Course{}, fmt.Errorf("course with id %s not found", id)
		}
		return mooc.Course{}, fmt.Errorf("error querying course: %v", err)
	}

	course, err := mooc.NewCourse(sc.ID, sc.Name, sc.Duration)
	if err != nil {
		return mooc.Course{}, fmt.Errorf("error creating course domain object: %v", err)
	}

	return course, nil
}

// GetByName implements the mooc.CourseRepository interface.
func (r *CourseRepository) GetByName(ctx context.Context, name string) ([]mooc.Course, error) {
	courseSQLStruct := sqlbuilder.NewStruct(new(sqlCourse))
	query, args := courseSQLStruct.SelectFrom(sqlCourseTable).Where("name = ?", name).Build()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying courses by name: %v", err)
	}
	defer rows.Close()

	var courses []mooc.Course
	for rows.Next() {
		var sc sqlCourse
		if err := rows.Scan(&sc.ID, &sc.Name, &sc.Duration); err != nil {
			return nil, fmt.Errorf("error scanning course row: %v", err)
		}
		course, err := mooc.NewCourse(sc.ID, sc.Name, sc.Duration)
		if err != nil {
			return nil, fmt.Errorf("error creating course domain object: %v", err)
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over course rows: %v", err)
	}

	return courses, nil
}
