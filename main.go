package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

var lengthRates = map[string]float64{
	"millimeter": 0.001,
	"centimeter": 0.01,
	"meter":      1.0,
	"kilometer":  1000.0,
	"inch":       0.0254,
	"foot":       0.3048,
	"yard":       0.9144,
	"mile":       1609.34,
}

var weightRates = map[string]float64{
	"milligram": 0.001,
	"gram":      1.0,
	"kilogram":  1000.0,
	"ounce":     28.3495,
	"pound":     453.592,
}

func convert(value float64, from string, to string, rates map[string]float64) float64 {
	baseValue := value * rates[from]

	return baseValue / rates[to]
}

func convertTemperature(value float64, from string, to string) float64 {
	var celsius float64

	switch from {
	case "Celsius":
		celsius = value
	case "Fahrenheit":
		celsius = (value - 32) * 5 / 9
	case "Kelvin":
		celsius = value - 273.15
	}

	switch to {
	case "Celsius":
		return celsius
	case "Fahrenheit":
		return (celsius * 9 / 5) + 32
	case "Kelvin":
		return +273.15
	}

	return celsius
}

var htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Unit Converter</title>
	<style>
		body { font-family: 'Comic Sans MS', cursive, sans-serif; display: flex; justify-content: center; margin-top: 50px; background-color: #f4f4f4; }
		.card { border: 3px solid #222; border-radius: 8px; padding: 40px; width: 400px; background-color: white; }
		h1 { margin-top: 0; font-size: 28px; }
		.nav { display: flex; gap: 20px; margin-bottom: 30px; font-size: 18px; font-weight: bold; }
		.nav a { text-decoration: none; color: #222; padding-bottom: 2px; }
		.nav a.active { color: #4285f4; border-bottom: 3px solid #4285f4; }
		
		label { display: block; font-weight: bold; margin-bottom: 8px; margin-top: 20px; font-size: 18px; }
		input, select { width: 100%; box-sizing: border-box; padding: 12px; border: 2px solid #222; border-radius: 6px; font-family: inherit; font-size: 16px; background-color: white;}
		
		button { padding: 10px 30px; border: 2px solid #222; border-radius: 6px; font-family: inherit; font-size: 18px; font-weight: bold; background-color: white; cursor: pointer; margin-top: 25px; }
		button:hover { background-color: #eee; }
		
		.result-title { font-weight: bold; font-size: 20px; margin-bottom: 10px; }
		.result-value { font-size: 40px; font-weight: bold; margin-bottom: 30px; }
	</style>
</head>
<body>
	<div class="card">
		<h1>Unit Converter</h1>
		
		<div class="nav">
			<a href="/length" class="{{if eq .Category "length"}}active{{end}}">Length</a>
			<a href="/weight" class="{{if eq .Category "weight"}}active{{end}}">Weight</a>
			<a href="/temperature" class="{{if eq .Category "temperature"}}active{{end}}">Temperature</a>
		</div>

		{{if .Result}}
			<div class="result-title">Result of your calculation</div>
			<div class="result-value">{{.Result}}</div>
			
			<form method="GET" action="/{{.Category}}">
				<button type="submit">Reset</button>
			</form>
		{{else}}
			<form method="POST">
				<label>Enter the {{.Category}} to convert</label>
				<input type="number" step="any" name="value" required>

				<label>Unit to Convert from</label>
				<select name="from">
					{{range .Units}} <option value="{{.}}">{{.}}</option> {{end}}
				</select>

				<label>Unit to Convert to</label>
				<select name="to">
					{{range .Units}} <option value="{{.}}">{{.}}</option> {{end}}
				</select>

				<button type="submit">Convert</button>
			</form>
		{{end}}
	</div>
</body>
</html>
`

var tmpl = template.Must(template.New("webpage").Parse(htmlTemplate))

type PageData struct {
	Category string
	Units    []string
	Result   string
}

func lengthHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Category: "length",
		Units:    []string{"millimeter", "centimeter", "kilometer", "inch", "foot", "yard", "mile"},
	}

	if r.Method == http.MethodPost {
		valStr := r.FormValue("value")
		from := r.FormValue("from")
		to := r.FormValue("to")

		value, _ := strconv.ParseFloat(valStr, 64)

		resultVal := convert(value, from, to, lengthRates)

		data.Result = fmt.Sprintf("%g %s = %g %s", value, from, resultVal, to)
	}

	tmpl.Execute(w, data)
}

func weightHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Category: "weight",
		Units:    []string{"milligram", "gram", "kilogram", "ounce", "pound"},
	}

	if r.Method == http.MethodPost {
		val, _ := strconv.ParseFloat(r.FormValue("value"), 64)
		from, to := r.FormValue("from"), r.FormValue("to")
		res := convert(val, from, to, weightRates)
		data.Result = fmt.Sprintf("%g %s = %g %s", val, from, res, to)
	}

	tmpl.Execute(w, data)
}

func tempHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Category: "temperature",
		Units: []string{
			"Celsius",
			"Fahrenheit",
			"Kelvin",
		},
	}

	if r.Method == http.MethodPost {
		val, _ := strconv.ParseFloat(r.FormValue("value"), 64)
		from, to := r.FormValue("from"), r.FormValue("to")
		res := convertTemperature(val, from, to)
		data.Result = fmt.Sprintf("%g %s = %g %s", val, from, res, to)
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", lengthHandler)
	http.HandleFunc("/length", lengthHandler)
	http.HandleFunc("/weight", weightHandler)
	http.HandleFunc("/temperature", tempHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
