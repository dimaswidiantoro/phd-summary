package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Chapter struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ChapterTitle string             `bson:"chapterTitle" json:"chapterTitle"`
	Subsections  []Subsection       `bson:"subsections" json:"subsections"`
	Tags         []string           `bson:"tags" json:"tags"`
	ImageURL     string             `bson:"imageURL" json:"imageURL"`
}

type Subsection struct {
	SubsectionTitle string    `bson:"subsectionTitle" json:"subsectionTitle"`
	Findings        []Finding `bson:"findings" json:"findings"`
}

type Finding struct {
	FindingDescription string   `bson:"findingDescription" json:"findingDescription"`
	SupportingAuthors  []string `bson:"supportingAuthors" json:"supportingAuthors"`
}

var client *mongo.Client

func CreateChapter(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var chapter Chapter
	_ = json.NewDecoder(request.Body).Decode(&chapter)
	if chapter.Tags == nil {
		chapter.Tags = []string{}
	}
	collection := client.Database("PhDSummary").Collection("Chapters")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.InsertOne(ctx, chapter)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(response).Encode(result)
}

func GetChapter(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var chapter Chapter
	collection := client.Database("PhDSummary").Collection("Chapters")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&chapter)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(response).Encode(chapter)
}

func GetChapters(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var chapters []Chapter
	collection := client.Database("PhDSummary").Collection("Chapters")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var chapter Chapter
		cursor.Decode(&chapter)
		chapters = append(chapters, chapter)
	}
	if err := cursor.Err(); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(response).Encode(chapters)
}

func UpdateChapter(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var chapter Chapter
	_ = json.NewDecoder(request.Body).Decode(&chapter)
	if chapter.Tags == nil {
		chapter.Tags = []string{}
	}
	collection := client.Database("PhDSummary").Collection("Chapters")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": chapter})
	if err != nil {
		log.Println("Error updating chapter:", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Chapter updated successfully")
	json.NewEncoder(response).Encode(chapter)
}

func GetTags(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	collection := client.Database("PhDSummary").Collection("Chapters")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Aggregation pipeline to get unique tags
	pipeline := mongo.Pipeline{
		{{"$unwind", bson.D{{"path", "$tags"}}}},
		{{"$group", bson.D{{"_id", nil}, {"tags", bson.D{{"$addToSet", "$tags"}}}}}},
		{{"$project", bson.D{{"_id", 0}, {"tags", 1}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(result) > 0 {
		json.NewEncoder(response).Encode(result[0]["tags"])
	} else {
		json.NewEncoder(response).Encode([]string{})
	}
}

func UploadImage(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	file, handler, err := request.FormFile("image")
	if err != nil {
		log.Println("Error retrieving the file:", err)
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure the images directory exists
	os.MkdirAll("images", os.ModePerm)

	// Create a temporary file within our images directory
	filePath := fmt.Sprintf("images/%s", handler.Filename)
	tempFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating the file:", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Read all the contents of the uploaded file into a byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading the file:", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write this byte array to our temporary file
	tempFile.Write(fileBytes)

	imageURL := fmt.Sprintf("/images/%s", handler.Filename)
	log.Println("Uploaded file:", imageURL)

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(map[string]string{
		"imageURL": imageURL,
	})
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/chapter", CreateChapter).Methods("POST")
	router.HandleFunc("/chapter/{id}", GetChapter).Methods("GET")
	router.HandleFunc("/chapters", GetChapters).Methods("GET")
	router.HandleFunc("/chapter/{id}", UpdateChapter).Methods("PUT")
	router.HandleFunc("/tags", GetTags).Methods("GET")
	router.HandleFunc("/upload", UploadImage).Methods("POST")

	// Serve static files (images)
	fs := http.FileServer(http.Dir("images"))
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", fs))

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}
