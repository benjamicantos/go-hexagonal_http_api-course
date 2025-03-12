package mysql

import (
	"context"
	"database/sql"
	"fmt"

	mooc "github.com/CodelyTV/go-hexagonal_http_api-course/03-02-repository-test/internal"
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


// Save implements the mooc.CourseRepository interface.
func (r *CourseRepository) GetCourse(ctx context.Context, ID mooc.CourseID) (mooc.Course, error) {
    row := r.db.QueryRowContext(ctx, "SELECT courses.id, courses.name, courses.duration FROM courses WHERE id = ?", ID.String())

    var courseID, courseName, courseDuration string
    err := row.Scan(&courseID, &courseName, &courseDuration)
    if err != nil {
        return mooc.Course{}, fmt.Errorf("error trying to get course from database: %w", err)
    }

    return mooc.NewCourse(courseID, courseName, courseDuration)	
}