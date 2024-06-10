package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Food struct to hold food information
type Food struct {
	ID       int
	Name     string
	Calories int
	Fat      float64
	Vegan    bool
	Diseases string
	Taste    string
	Age      string
	Score    float64
}

// Load data from CSV file
func loadFoodData() ([]Food, error) {
	file, err := os.Open("data.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	foods := make([]Food, 0, len(records)-1) // Exclude header row
	for _, record := range records[1:] {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		calories, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, err
		}
		fat, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}
		vegan := strings.ToLower(record[4]) == "tidak"
		diseases := record[5]
		taste := record[6]
		age := record[7]

		foods = append(foods, Food{
			ID:       id,
			Name:     record[1],
			Calories: calories,
			Fat:      fat,
			Vegan:    vegan,
			Diseases: diseases,
			Taste:    taste,
			Age:      age,
		})
	}

	return foods, nil
}

// Get user preferences
func getUserPreferences() (int, int, bool) {
	var maxCalories, minProtein int
	var veganPreference string

	fmt.Println("Enter your maximum calorie intake per food item:")
	fmt.Scan(&maxCalories)

	fmt.Println("Enter your minimum protein requirement per food item:")
	fmt.Scan(&minProtein)

	fmt.Println("Do you prefer vegan options? (yes/no):")
	fmt.Scan(&veganPreference)

	isVegan := strings.ToLower(veganPreference) == "yes"

	return maxCalories, minProtein, isVegan
}

// Normalize and calculate score for each food item
func calculateScores(foods []Food, maxCalories, minProtein int, isVegan bool) {
	var maxCaloriesVal, maxFatVal float64

	// Determine maximum values for normalization
	for _, food := range foods {
		if float64(food.Calories) > maxCaloriesVal {
			maxCaloriesVal = float64(food.Calories)
		}
		if food.Fat > maxFatVal {
			maxFatVal = food.Fat
		}
	}

	// Calculate score for each food item
	for i := range foods {
		score := 0.0

		// Normalize calories (assuming less is better)
		if maxCaloriesVal != 0 {
			score += (maxCaloriesVal - float64(foods[i].Calories)) / maxCaloriesVal
		}

		// Normalize fat (assuming less is better)
		if maxFatVal != 0 {
			score += (maxFatVal - foods[i].Fat) / maxFatVal
		}

		// Add vegan preference score
		if isVegan && foods[i].Vegan {
			score += 1.0
		}

		foods[i].Score = score
	}
}

// Filter and recommend foods based on user preferences
func recommendFoods(foods []Food, maxCalories, minProtein int, isVegan bool) []Food {
	var recommendations []Food

	calculateScores(foods, maxCalories, minProtein, isVegan)

	for _, food := range foods {
		if food.Calories <= maxCalories && food.Fat <= 5 && (!isVegan || (isVegan && food.Vegan)) {
			recommendations = append(recommendations, food)
		}
	}

	// Sort recommendations by score in descending order
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	return recommendations
}

// Display recommended foods
func displayRecommendations(recommendations []Food) {
	if len(recommendations) == 0 {
		fmt.Println("No food items match your preferences.")
	} else {
		fmt.Println("Recommended food items for you:")
		for _, food := range recommendations {
			fmt.Printf("- %s: %d calories, %.1f g fat, Score: %.2f\n", food.Name, food.Calories, food.Fat, food.Score)
		}
	}
}

func main() {
	foods, err := loadFoodData()
	if err != nil {
		fmt.Println("Error loading food data:", err)
		return
	}

	maxCalories, minProtein, isVegan := getUserPreferences()
	recommendations := recommendFoods(foods, maxCalories, minProtein, isVegan)

	displayRecommendations(recommendations)
}
