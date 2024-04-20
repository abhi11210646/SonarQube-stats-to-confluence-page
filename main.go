package main

func main() {
	// sonarConfig := GetSonarConfig()

	// for _, projectKey := range SonarConfig.Projects {
	data := SonarStats("app.webmail")
	// }

	// getByPageId()
	updateByPageId(data)

	// fmt.Println(generetaeHTML(data))
}
