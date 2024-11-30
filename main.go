package main

import (
 "Update-Game1/domain"
 "encoding/json"
 "fmt"
 "math/rand"
 "os"
 "sort"
 "strconv"
 "time"
)

const (
 totalPoints      int = 100
 pointPerQuestion int = 100
)

var id uint64 = 1

func main() {
    users := getUsers()
    for _, user := range users {
        if user.Id >= id {
            id = user.Id + 1
        }
    }
    fmt.Println("Welcome to game")
    
    for {
        menu()
        point := ""
        fmt.Scan(&point)
        switch point {
        case "1":
            user := play()
            users = getUsers()
            users = append(users, user)
            sortAndSave(users)
        case "2":
            users = getUsers()
            for _, user := range users {
                fmt.Printf("Id: %v Name: %s Time: %v \n", user.Id, user.Name, user.TimeSpent)
            }
        case "3":
            clearFile()
        case "4":
            return
        default:
            fmt.Println("Choose a valid option: 1, 2, 3, or 4")
        }
    }
}

func menu() {
    println("1. Startgame")
    println("2. Result")
    println("3. Clear results")
    println("4. Exit")
}

func play() domain.User {
    for i := 5; i >= 1; i-- {
        fmt.Printf("Game starts in: %v\n", i)
        time.Sleep(1 * time.Second)
    }
   
    startTime := time.Now()
    myPoint := 0
    operators := []string{"+", "-", "*", "/"}
    for myPoint < totalPoints {
        x, y := rand.Intn(100)+1, rand.Intn(100)+1
        operator := operators[rand.Intn(len(operators))]
      
        var correctAnswer int
        switch operator {
        case "+":
            correctAnswer = x + y
        case "-":
            correctAnswer = x - y
        case "*":
            correctAnswer = x * y
        case "/":
            if y == 0 {
                y = 1 // Avoid division by zero
            }
            correctAnswer = x / y
        }
        fmt.Printf("%v %s %v =", x, operator, y)
        ans := ""
        fmt.Scan(&ans)
      
        ansInt, err := strconv.Atoi(ans)
        if err != nil {
            fmt.Println("Please write an integer number")
        }else {

            if ansInt == correctAnswer {
                myPoint += pointPerQuestion
                fmt.Println("My points:", myPoint)
                fmt.Printf("Points remaining: %v\n", totalPoints-myPoint)
            }else {
                fmt.Println("Try again")
            }
        }
    }
    endTime := time.Now()
    timeSpent := endTime.Sub(startTime)
    fmt.Println("Game over, you win (time):", timeSpent)
    fmt.Println("Write your name: ")
   
    name := ""
    fmt.Scan(&name)
   
    user := domain.User{
        Id:        id,
        Name:      name,
        TimeSpent: timeSpent,
    }
    id++
    return user
}

func sortAndSave(users []domain.User) {
    sort.SliceStable(users, func(i, j int) bool {
     return users[i].TimeSpent < users[j].TimeSpent
    })
   
    file, err := os.OpenFile("users.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    if err != nil {
        fmt.Printf("sortAndSave -> os.OpenFile: %s\n", err)
        return
    }
    defer func(file *os.File) {
        err = file.Close()
        if err != nil {
            fmt.Printf("Error: %s", err)
        }
    }(file)
   
    encoder := json.NewEncoder(file)
    err = encoder.Encode(users)
    if err != nil {
        fmt.Printf("sortAndSave -> encoder.Encode: %s\n", err)
        return
    }
}

func getUsers() []domain.User {
    file, err := os.Open("users.json")
    if err != nil {
        if os.IsNotExist(err) {
            _, err = os.Create("users.json")
            if err != nil {
                fmt.Printf("getUsers -> os.Create: %s\n", err)
                return nil
            }
            return nil
        }
        fmt.Printf("getUsers -> os.Open: %s\n", err)
        return nil
    }
    defer file.Close()
   
    var users []domain.User
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&users)
    if err != nil {
        fmt.Printf("getUsers -> decoder.Decode: %s\n", err)
        return nil
    }
    return users
}

func clearFile() {
    err := os.WriteFile("users.json", []byte(""), 0644)
    if err != nil {
        fmt.Printf("clearFile -> os.WriteFile: %s\n", err)
    }else {
        fmt.Println("File cleared successfully.")
    }
}