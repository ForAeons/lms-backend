// nolint
package main

import (
	"fmt"
	"lms-backend/cmd/seed/helper"
	"lms-backend/internal/app"
	"lms-backend/internal/dataaccess/user"
	"lms-backend/internal/database"
	"lms-backend/internal/model"
	"lms-backend/internal/policy/abilities"
	"log"
	"time"

	"github.com/go-loremipsum/loremipsum"
	"gorm.io/gorm"
)

func main() {
	var err error

	// Load environment variables and connect to database
	err = app.LoadEnvAndConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	db := database.GetDB()

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	fmt.Println("Seeding database...")
	fmt.Println("Seeding users and people...")
	err = seedUsersAndPeople(tx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Seeding books...")
	err = seedBooks(tx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Seeding roles...")
	err = seedRolesAbilities(tx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Seeding complete!")
}

func seedUsersAndPeople(db *gorm.DB) error {
	var count int64

	result := db.Model(&model.User{}).Count(&count)
	if result.Error != nil {
		return result.Error
	}

	if count >= 100 {
		return nil
	}

	loremIpsumGenerator := loremipsum.New()

	// Generate 100 users
	users := make([]model.User, 10)
	for i := 1; i <= 10; i++ {
		users[i-1] = model.User{
			Email:             fmt.Sprintf("user%d@gmail.com", i),
			EncryptedPassword: "P4ssw0rd!",
			Person: &model.Person{
				FirstName: loremIpsumGenerator.Word(),
				LastName:  loremIpsumGenerator.Word(),
			},
		}
	}

	return db.Create(&users).Error
}

func seedBooks(db *gorm.DB) error {
	var count int64

	result := db.Model(&model.Book{}).Count(&count)
	if result.Error != nil {
		return result.Error
	}

	if count >= 1000 {
		return nil
	}

	loremIpsumGenerator := loremipsum.New()

	books := make([]model.Book, 3000)
	for i := 1; i <= 3000; i++ {
		books[i-1] = model.Book{
			Title:           loremIpsumGenerator.Words(helper.RandInt(4, 11)),
			Author:          loremIpsumGenerator.Words(helper.RandInt(2, 5)),
			ISBN:            helper.GenerateISBN13(),
			Publisher:       loremIpsumGenerator.Words(helper.RandInt(4, 7)),
			PublicationDate: helper.RandomDate(time.Now().AddDate(-10, 0, 0), time.Now()),
			Genre:           loremIpsumGenerator.Words(helper.RandInt(1, 3)),
			Language:        loremIpsumGenerator.Words(helper.RandInt(1, 3)),
		}
	}

	return db.Create(&books).Error
}

func seedRolesAbilities(db *gorm.DB) error {
	roles := []model.Role{
		{
			Name: "System Admin",
		},
		{
			Name: "Library Admin",
		},
		{
			Name: "Staff",
		},
		{
			Name: "Basic",
		},
	}

	// Create roles
	if err := db.Create(&roles).Error; err != nil {
		return err
	}

	abilities := []model.Ability{
		abilities.CanManageAll,

		abilities.CanReadAuditLog,
		abilities.CanCreateAuditLog,

		abilities.CanUpdateUser,
		abilities.CanUpdateRole,

		abilities.CanCreatePerson,
		abilities.CanUpdatePerson,

		abilities.CanReadBook,
		abilities.CanCreateBook,
		abilities.CanUpdateBook,
		abilities.CanDeleteBook,
		abilities.CanBorrowBook,
		abilities.CanReturnBook,
		abilities.CanRenewBook,
		abilities.CanManageBookRecords,
	}

	// Create abilities
	if err := db.Create(&abilities).Error; err != nil {
		return err
	}

	const (
		T = true
		F = false
	)

	var (
		rolesAbilitiesMap = [][4]bool{
			{T, F, F, F}, // Manage all

			{T, T, F, F}, // Read audit log
			{T, T, F, F}, // Create audit log

			{T, T, T, F}, // Update user
			{T, T, F, F}, // Update role

			{T, T, T, F}, // Create person
			{T, T, T, F}, // Update person

			{T, T, T, T}, // Read book
			{T, T, F, F}, // Create book
			{T, T, F, F}, // Update book
			{T, T, F, F}, // Delete book

			{T, T, T, T}, // Borrow book
			{T, T, T, T}, // Return book
			{T, T, T, T}, // Renew book
			{T, T, T, F}, // Manage book records
		}
	)

	// Assign abilities to roles
	for _, role := range roles {
		for i, ability := range abilities {
			if !rolesAbilitiesMap[i][role.ID-1] {
				continue
			}

			if err := db.Exec(
				fmt.Sprintf("%s %s %s %s",
					"INSERT INTO role_abilities (role_id, ability_id)",
					"SELECT",
					"(SELECT id FROM roles WHERE name = ?),",
					"(SELECT id FROM abilities WHERE name = ?)",
				),
				role.Name, ability.Name,
			).Error; err != nil {
				return err
			}
		}
	}

	// Assign roles to users
	for i := int64(1); i <= 4; i++ {
		_, err := user.UpdateRoles(db, i, []int64{i})
		if err != nil {
			return err
		}
	}

	return nil
}
