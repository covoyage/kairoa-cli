package cmd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var mockCmd = &cobra.Command{
	Use:   "mock",
	Short: "Mock data generator",
	Long:  `Generate mock data for testing and development.`,
}

var firstNames = []string{"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda", "William", "Elizabeth"}
var lastNames = []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez"}
var domains = []string{"gmail.com", "yahoo.com", "hotmail.com", "outlook.com", "example.com"}
var companies = []string{"Acme Corp", "Globex", "Soylent Corp", "Initech", "Umbrella Corp", "Hooli", "Vehement Capital"}
var jobTitles = []string{"Software Engineer", "Product Manager", "Data Scientist", "Designer", "DevOps Engineer", "QA Engineer"}
var cities = []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego"}
var countries = []string{"USA", "UK", "Canada", "Australia", "Germany", "France", "Japan", "China"}

var mockUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Generate mock user data",
	RunE: func(cmd *cobra.Command, args []string) error {
		count, _ := cmd.Flags().GetInt("count")
		format, _ := cmd.Flags().GetString("format")

		rand.Seed(time.Now().UnixNano())

		type User struct {
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Email     string `json:"email"`
			Phone     string `json:"phone"`
			Age       int    `json:"age"`
			City      string `json:"city"`
			Country   string `json:"country"`
		}

		var users []User
		for i := 0; i < count; i++ {
			first := firstNames[rand.Intn(len(firstNames))]
			last := lastNames[rand.Intn(len(lastNames))]
			users = append(users, User{
				FirstName: first,
				LastName:  last,
				Email:     fmt.Sprintf("%s.%s@%s", strings.ToLower(first), strings.ToLower(last), domains[rand.Intn(len(domains))]),
				Phone:     fmt.Sprintf("+1-555-%03d-%04d", rand.Intn(1000), rand.Intn(10000)),
				Age:       18 + rand.Intn(50),
				City:      cities[rand.Intn(len(cities))],
				Country:   countries[rand.Intn(len(countries))],
			})
		}

		outputMockData(users, format)
		return nil
	},
}

var mockEmployeeCmd = &cobra.Command{
	Use:   "employee",
	Short: "Generate mock employee data",
	RunE: func(cmd *cobra.Command, args []string) error {
		count, _ := cmd.Flags().GetInt("count")
		format, _ := cmd.Flags().GetString("format")

		rand.Seed(time.Now().UnixNano())

		type Employee struct {
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Email     string `json:"email"`
			Company   string `json:"company"`
			JobTitle  string `json:"jobTitle"`
			Salary    int    `json:"salary"`
		}

		var employees []Employee
		for i := 0; i < count; i++ {
			first := firstNames[rand.Intn(len(firstNames))]
			last := lastNames[rand.Intn(len(lastNames))]
			employees = append(employees, Employee{
				FirstName: first,
				LastName:  last,
				Email:     fmt.Sprintf("%s.%s@%s", strings.ToLower(first), strings.ToLower(last), domains[rand.Intn(len(domains))]),
				Company:   companies[rand.Intn(len(companies))],
				JobTitle:  jobTitles[rand.Intn(len(jobTitles))],
				Salary:    50000 + rand.Intn(150000),
			})
		}

		outputMockData(employees, format)
		return nil
	},
}

var mockAddressCmd = &cobra.Command{
	Use:   "address",
	Short: "Generate mock address data",
	RunE: func(cmd *cobra.Command, args []string) error {
		count, _ := cmd.Flags().GetInt("count")
		format, _ := cmd.Flags().GetString("format")

		rand.Seed(time.Now().UnixNano())

		streets := []string{"Main St", "Oak Ave", "Maple Rd", "Cedar Ln", "Pine Dr", "Elm St", "Washington Ave"}
		states := []string{"CA", "NY", "TX", "FL", "IL", "PA", "OH", "GA"}

		type Address struct {
			Street  string `json:"street"`
			City    string `json:"city"`
			State   string `json:"state"`
			ZipCode string `json:"zipCode"`
			Country string `json:"country"`
		}

		var addresses []Address
		for i := 0; i < count; i++ {
			addresses = append(addresses, Address{
				Street:  fmt.Sprintf("%d %s", 100+rand.Intn(9900), streets[rand.Intn(len(streets))]),
				City:    cities[rand.Intn(len(cities))],
				State:   states[rand.Intn(len(states))],
				ZipCode: fmt.Sprintf("%05d", rand.Intn(100000)),
				Country: "USA",
			})
		}

		outputMockData(addresses, format)
		return nil
	},
}

var mockProductCmd = &cobra.Command{
	Use:   "product",
	Short: "Generate mock product data",
	RunE: func(cmd *cobra.Command, args []string) error {
		count, _ := cmd.Flags().GetInt("count")
		format, _ := cmd.Flags().GetString("format")

		rand.Seed(time.Now().UnixNano())

		productNames := []string{"Laptop", "Phone", "Tablet", "Monitor", "Keyboard", "Mouse", "Headphones", "Camera"}
		categories := []string{"Electronics", "Accessories", "Computers", "Mobile"}

		type Product struct {
			ID       string  `json:"id"`
			Name     string  `json:"name"`
			Category string  `json:"category"`
			Price    float64 `json:"price"`
			InStock  bool    `json:"inStock"`
		}

		var products []Product
		for i := 0; i < count; i++ {
			products = append(products, Product{
				ID:       fmt.Sprintf("PROD-%04d", i+1),
				Name:     productNames[rand.Intn(len(productNames))],
				Category: categories[rand.Intn(len(categories))],
				Price:    float64(int((10+rand.Float64()*990)*100)) / 100,
				InStock:  rand.Float32() > 0.2,
			})
		}

		outputMockData(products, format)
		return nil
	},
}

func outputMockData(data interface{}, format string) {
	switch format {
	case "json":
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(jsonData))
	default:
		// Simple text output
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(jsonData))
	}
}

func init() {
	rootCmd.AddCommand(mockCmd)
	mockCmd.AddCommand(mockUserCmd)
	mockCmd.AddCommand(mockEmployeeCmd)
	mockCmd.AddCommand(mockAddressCmd)
	mockCmd.AddCommand(mockProductCmd)

	mockUserCmd.Flags().IntP("count", "c", 5, "Number of records to generate")
	mockUserCmd.Flags().StringP("format", "f", "json", "Output format (json)")

	mockEmployeeCmd.Flags().IntP("count", "c", 5, "Number of records to generate")
	mockEmployeeCmd.Flags().StringP("format", "f", "json", "Output format (json)")

	mockAddressCmd.Flags().IntP("count", "c", 5, "Number of records to generate")
	mockAddressCmd.Flags().StringP("format", "f", "json", "Output format (json)")

	mockProductCmd.Flags().IntP("count", "c", 5, "Number of records to generate")
	mockProductCmd.Flags().StringP("format", "f", "json", "Output format (json)")
}
