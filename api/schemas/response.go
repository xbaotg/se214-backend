package schemas

import (
	"be/models"

	"github.com/google/uuid"
)

type CreateCourseResponse struct {
	ID                uuid.UUID  `json:"id"`
	CourseTeacherID   uuid.UUID  `json:"course_teacher_id"`
	CourseTeacherName string     `json:"course_teacher_name"`
	CourseDepartment  uuid.UUID  `json:"course_department"`
	CourseName        string     `json:"course_name"`
	CourseFullname    string     `json:"course_fullname"`
	CourseCredit      int32      `json:"course_credit"`
	CourseYear        int32      `json:"course_year"`
	CourseSemester    int32      `json:"course_semester"`
	CourseStartShift  int32      `json:"course_start_shift"`
	CourseEndShift    int32      `json:"course_end_shift"`
	CourseDay         models.Day `json:"course_day"`
	MaxEnroll         int32      `json:"max_enroll"`
	CurrentEnroll     int32      `json:"current_enroll"`
	CourseRoom        string     `json:"course_room"`
}
