package controllers

import (
	"academy-go-q12021/models"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func HelloWorld(ctx *fiber.Ctx) error {
	return ctx.Send([]byte("Hello World!"))
}

func GeneratePokemonCSV(ctx *fiber.Ctx) error {
	pokemonName := ctx.Params("name")

	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + pokemonName)

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	if resp.StatusCode == 404 {
		return ctx.JSON(fiber.Map{
			"message": "Invalid pokemon name",
		})
	} else if resp.StatusCode != 200 {
		return ctx.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	defer resp.Body.Close()

	var pokemon models.Pokemon

	err = json.NewDecoder(resp.Body).Decode(&pokemon)

	filepath := "./generated/output.csv"
	file, err := os.Create(filepath)

	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)

	writer.Write([]string{
		"ID",
		"Name",
		"Type",
		"Picture",
		"Abilities",
	})

	var types, abilities string

	for _, value := range pokemon.Types {
		types += value.Type.Name + ", "
	}

	for _, value := range pokemon.Abilities {
		abilities += value.Ability.Name + ", "
	}

	writer.Write([]string{
		strconv.Itoa(pokemon.ID),
		pokemon.Name,
		types,
		pokemon.Sprites.FrontDefault,
		abilities,
	})

	writer.Flush()
	_ = file.Close()

	/*return ctx.JSON(fiber.Map{
		"message": "Everything went fine...",
	})*/

	return ctx.Download(filepath)
}

func GetTop10Pokemons(ctx *fiber.Ctx) error {

	/**
	1. Reach PokeAPI w/pokemons data
	2. Build CSV from previous data
	3. Download CSV
	*/

	//1. Collecting data
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/pikachu")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	/*err = json.NewDecoder(resp.Body).Decode(&pokemon)
	//err := json.Unmarshal(resp.Body, &pokemon)

	if err != nil {
		return err
	}*/

	body, err := ioutil.ReadAll(resp.Body)

	//fmt.Printf("%s", body)

	//pokemon := Pokemon{}
	var pokemon map[string]interface{}

	err = json.Unmarshal(body, &pokemon)

	if err != nil {
		return err
	}

	for key, value := range pokemon {
		if key == "name" {
			fmt.Println(key, value)
		}
	}

	fmt.Println("Pokemon name: " + pokemon["name"].(string))
	fmt.Println("Pokemon id: ", pokemon["id"].(float64))
	//fmt.Println("abilities: ", pokemon["abilities"].(map[interface{}])

	//2. Generate File
	filepath := "./generated/output.csv"

	//if err := generateCSV(filepath); err != nil {
	//	return err
	//}

	file, err := os.Create(filepath)

	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)

	//defer writer.Flush()

	writer.Write([]string{
		"ID",
		"Name",
		"Type",
		"Picture",
		"Abilities",
	})

	id := pokemon["id"].(float64)
	name := pokemon["name"].(string)
	//types := pokemon["types"].(map[string]interface{})
	pokemon_types := ""
	pictures := pokemon["sprites"].(map[string]interface{})
	abilities := ""

	//for _, value := range types {
	//	pokemon_types += value["type"]["name"]
	//}

	writer.Write([]string{
		fmt.Sprintf("%.0f", id),
		name,
		pokemon_types,
		pictures["front_default"].(string),
		abilities,
	})

	writer.Flush()
	_ = file.Close()

	return ctx.Download(filepath)

	//return ctx.JSON(fiber.Map{
	//	"message": "Here you have 10 pokemons!",
	//})
}

//func generateCSV(filepath string) error {
//
//}

func WriteCSV(ctx *fiber.Ctx) error {

	form, _ := ctx.MultipartForm()

	/*if err != nil || err != err{}; {
		return ctx.JSON(fiber.Map{
			"message": "Error uploading files...",
			"error": err,
		})
	}*/

	//form.File["paramName"] gets an array of files
	file := form.File["csvFile"][0]

	//Saving file to x folder

	ctx.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename))

	lines, err := readCSV("./uploads/" + file.Filename)

	if err != nil {
		panic(err)
	}

	/*for _, line := range lines {
		fmt.Println(line[0] + ": " + line[1])
	}*/

	//matrix := make([][]string, 0)

	for _, line := range lines {
		for key, _ := range line {
			fmt.Println(line[key])
			//fmt.Printf("%T\n", line)

			//length := len(line)
			//var array [length]string

		}
	}

	return ctx.JSON(fiber.Map{
		"message": "So far so good...",
		//"form":    form.File,
		//"file": files,
	})

}

func ReadCSV(ctx *fiber.Ctx) error {

	readCSV("./test.csv")

	return ctx.JSON(fiber.Map{
		"message": "CSV readed",
	})
}

func readCSV(filePath string) ([][]string, error) {

	f, err := os.Open(filePath)

	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	//Read file into variables
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	fmt.Println(lines)

	return lines, nil
}

func Import(ctx *fiber.Ctx) error {
	form, _ := ctx.MultipartForm()

	file := form.File["csvFile"][0]

	ctx.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename))

	lines, err := readCSV("./uploads/" + file.Filename)

	if err != nil {
		return err
	}

	pokemon := make([]models.Input, 0)

	for i, line := range lines {
		//for key, _ := range line {
		//	fmt.Println(key, line[key])
		//}
		if i == 0 {
			continue
		}

		//fmt.Println("Line: ", line, i)
		pokemon = append(pokemon, models.Input{
			ID:        line[0],
			Name:      line[1],
			Types:     line[2],
			Picture:   line[3],
			Abilities: line[4],
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "imported successfully",
		"data":    pokemon,
	})
}
