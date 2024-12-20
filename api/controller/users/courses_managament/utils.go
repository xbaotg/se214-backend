package coursesmanagament

import (
	"be/bootstrap"
	"be/models"
	"errors"

	// "gorm.io/gorm"
)

func CheckPrerequisite(app *bootstrap.App, user *models.User, course *models.Course) (string, error) {
	prerequisites := []models.PrerequisiteCourse{}
	if err := app.DB.Table("prerequisite_courses").Where("course_id = ?", course.CourseName).Find(&prerequisites).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		return "", err
	}
	
	// check if user has taken the prerequisite course
	prerequisitesCoursesName := []string{}
	for _, prerequisite := range prerequisites {
		prerequisitesCoursesName = append(prerequisitesCoursesName, prerequisite.PrerequisiteID)
	}

	app.Logger.Info().Msgf("prerequisitesCoursesName: %v", prerequisitesCoursesName)

	// check if user has taken the prerequisite course
	type Result struct {
		CourseName string
		Status     models.CoStatus
	}
	if len(prerequisitesCoursesName) == 0 {
		return "", nil
	}
	results := []Result{}
	if err := app.DB.Table("registered_courses").Joins("JOIN courses ON registered_courses.course_id = courses.id").Where(
		"registered_courses.user_id = ? and courses.course_name IN ?", user.ID, prerequisitesCoursesName,
		).Select("courses.course_name as course_name, registered_courses.status as status").Scan(&results).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		return "", err
	}
	if len(results) == 0 {
		return prerequisitesCoursesName[0], errors.New("Người dùng chưa đăng ký môn học tiên quyết")
	}
	app.Logger.Info().Msgf("results: %v", results)
	// check if user has taken the prerequisite course
	for _,courseName := range prerequisitesCoursesName {
		done := false
		for _, result := range results {
			if result.CourseName == courseName && result.Status == models.CoStatusDone {
				done = true
				break
			}
		}
		if !done {
			return courseName, errors.New("Nguời dùng chưa hoàn thành môn học tiên quyết")
		}
	}

	return "", nil
}
