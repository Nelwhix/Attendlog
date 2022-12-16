package Controllers

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	Semester string `valid:"required"`
	Name     string `valid:"required"`
	Code     string `valid:"required"`
}

type User struct {
	Username string `valid:"alpha,required"`
	Password string `valid:"alpha,required"`
}

type Record struct {
	gorm.Model
	Course string 
	Name string `valid:"required"`
	Matric string `valid:"numeric,required"`
	Signature string
}