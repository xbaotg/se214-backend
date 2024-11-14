package coursesmanagament

import (
	"be/bootstrap"
	"be/models"
	"errors"

	"gorm.io/gorm"
)

func CheckPrerequisite(app *bootstrap.App, user *models.User, course *models.Course) error {
	prerequisite := models.PrerequisiteCourse{CourseID: course.ID}

	if err := app.DB.First(&prerequisite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// if course has no prerequisite
			return nil
		}

		app.Logger.Error().Err(err).Msg(err.Error())
		return err
	}

	// check if user has taken the prerequisite course
	userCourse := models.RegisteredCourse{UserID: user.ID, CourseID: prerequisite.PrerequisiteID}
	if err := app.DB.First(&userCourse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user has not taken the prerequisite course")
		}

		app.Logger.Error().Err(err).Msg(err.Error())
		return err
	}

	if userCourse.Status != models.CoStatusDone {
		return errors.New("user has not completed the prerequisite course")
	}

	return nil
}
