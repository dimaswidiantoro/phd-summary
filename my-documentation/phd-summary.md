---
layout: default
title: PhD Summary
---

Sure! Here’s a comprehensive and well-structured documentation for your PhD Summary app, integrating all the necessary steps, code examples, and troubleshooting tips.

---

# PhD Summary Application Documentation

## Introduction

This documentation provides a comprehensive guide to creating a full-stack application for managing PhD summaries using Go, MongoDB, and React.js. It includes backend setup with Go, MongoDB for data storage, and a React.js frontend.

## Table of Contents

1. [Setting Up the Backend with Go](#setting-up-the-backend-with-go)
   - [Defining Data Models](#defining-data-models)
   - [Implementing CRUD Operations](#implementing-crud-operations)
2. [Setting Up the Frontend with React](#setting-up-the-frontend-with-react)
   - [Creating a React App](#creating-a-react-app)
   - [Creating Components](#creating-components)
   - [API Service](#api-service)
3. [Running the Application](#running-the-application)
4. [Common Issues and Troubleshooting](#common-issues-and-troubleshooting)
   - [Jekyll Installation Issues](#jekyll-installation-issues)
   - [Liquid Syntax Errors](#liquid-syntax-errors)
5. [Conclusion](#conclusion)

## Setting Up the Backend with Go

### Defining Data Models

Create a file named `main.go` and define the data models.

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
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
```

### Implementing CRUD Operations

Add the CRUD operation functions to `main.go`.

```go
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

{% raw %}
func GetTags(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("content-type", "application/json")
    collection := client.Database("PhDSummary").Collection("Chapters")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

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
{% raw %}

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

    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    })

    handler := c.Handler(router)
    log.Fatal(http.ListenAndServe(":8000", handler))
}
```

## Setting Up the Frontend with React

### Creating a React App

1. **Create a React App**:

   ```sh
   npx create-react-app phd-summary-frontend
   cd phd-summary-frontend
   ```

2. **Install Axios**:

   ```sh
   npm install axios
   ```

3. **Install React Router**:

   ```sh
   npm install react-router-dom
   ```

### Creating Components

Create the necessary React components for your application.

#### AddChapter Component

**src/components/AddChapter.js**:

```jsx
import React, { useState } from 'react';
import { createChapter } from '../services/api';

const AddChapter = ({ onAddChapter }) => {
  const [chapterTitle, setChapterTitle] = useState('');
  const [tags, setTags] = useState('');
  const [subsections, setSubsections] = useState([{ subsectionTitle: '', findings: [{ findingDescription: '', supportingAuthors: [] }] }]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const newChapter = { chapterTitle, tags: tags.split(','), subsections };
    await createChapter(newChapter);
    setChapterTitle('');
    setTags('');
    setSubsections([{ subsectionTitle: '', findings: [{ findingDescription: '', supportingAuthors: [] }] }]);
    onAdd

Chapter();
  };

  return (
    <div>
      <h1>Add Chapter</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Chapter Title:</label>
          <input type="text" value={chapterTitle} onChange={(e) => setChapterTitle(e.target.value)} required />
        </div>
        <div>
          <label>Tags (comma separated):</label>
          <input type="text" value={tags} onChange={(e) => setTags(e.target.value)} required />
        </div>
        <div>
          <label>Subsections:</label>
          {subsections.map((subsection, sIdx) => (
            <div key={sIdx}>
              <input type="text" placeholder="Subsection Title" value={subsection.subsectionTitle}
                onChange={(e) => {
                  const newSubsections = [...subsections];
                  newSubsections[sIdx].subsectionTitle = e.target.value;
                  setSubsections(newSubsections);
                }} required />
              {subsection.findings.map((finding, fIdx) => (
                <div key={fIdx}>
                  <input type="text" placeholder="Finding Description" value={finding.findingDescription}
                    onChange={(e) => {
                      const newSubsections = [...subsections];
                      newSubsections[sIdx].findings[fIdx].findingDescription = e.target.value;
                      setSubsections(newSubsections);
                    }} required />
                  <input type="text" placeholder="Supporting Authors (comma separated)"
                    onChange={(e) => {
                      const newSubsections = [...subsections];
                      newSubsections[sIdx].findings[fIdx].supportingAuthors = e.target.value.split(',');
                      setSubsections(newSubsections);
                    }} required />
                </div>
              ))}
            </div>
          ))}
        </div>
        <button type="submit">Add Chapter</button>
      </form>
    </div>
  );
};

export default AddChapter;
```

#### ChaptersList Component

**src/components/ChaptersList.js**:

```jsx
import React, { useState } from 'react';
import { Link } from 'react-router-dom';

const ChaptersList = ({ chapters, tags }) => {
  const [filterTag, setFilterTag] = useState('');

  const handleFilterChange = (e) => {
    setFilterTag(e.target.value);
  };

  const filteredChapters = chapters.filter(chapter =>
    filterTag === '' || (chapter.tags && chapter.tags.includes(filterTag))
  );

  return (
    <div>
      <h2>Chapters</h2>
      <div>
        <label>Filter by tag:</label>
        <select onChange={handleFilterChange} value={filterTag}>
          <option value="">All</option>
          {tags.map((tag, index) => (
            <option key={index} value={tag}>{tag}</option>
          ))}
        </select>
      </div>
      <ul>
        {filteredChapters.map((chapter) => (
          <li key={chapter.id}>
            <h3>{chapter.chapterTitle}</h3>
            <p><strong>Tags:</strong> {chapter.tags ? chapter.tags.join(', ') : 'No tags'}</p>
            {chapter.subsections.map((subsection, sIdx) => (
              <div key={sIdx}>
                <h4>{subsection.subsectionTitle}</h4>
                {subsection.findings.map((finding, fIdx) => (
                  <div key={fIdx} style={{ marginLeft: '20px' }}>
                    <p><strong>Finding:</strong> {finding.findingDescription}</p>
                    <p><strong>Supporting Authors:</strong> {finding.supportingAuthors.join(', ')}</p>
                  </div>
                ))}
              </div>
            ))}
            <Link to={`/edit/${chapter.id}`}>Edit</Link>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default ChaptersList;
```

#### EditChapter Component

**src/components/EditChapter.js**:

```jsx
import React, { useEffect, useState } from 'react';
import { getChapter, updateChapter } from '../services/api';
import { useParams, useNavigate } from 'react-router-dom';

const EditChapter = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [chapterTitle, setChapterTitle] = useState('');
  const [subsections, setSubsections] = useState([]);

  useEffect(() => {
    const fetchChapter = async () => {
      const chapter = await getChapter(id);
      setChapterTitle(chapter.chapterTitle);
      setSubsections(chapter.subsections);
    };
    fetchChapter();
  }, [id]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const updatedChapter = { chapterTitle, subsections };
    await updateChapter(id, updatedChapter);
    navigate('/');
  };

  return (
    <div>
      <h1>Edit Chapter</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Chapter Title:</label>
          <input type="text" value={chapterTitle} onChange={(e) => setChapterTitle(e.target.value)} required />
        </div>
        <div>
          <label>Subsections:</label>
          {subsections.map((subsection, sIdx) => (
            <div key={sIdx}>
              <input type="text" placeholder="Subsection Title" value={subsection.subsectionTitle}
                onChange={(e) => {
                  const newSubsections = [...subsections];
                  newSubsections[sIdx].subsectionTitle = e.target.value;
                  setSubsections(newSubsections);
                }} required />
              {subsection.findings.map((finding, fIdx) => (
                <div key={fIdx} style={{ marginLeft: '20px' }}>
                  <input type="text" placeholder="Finding Description" value={finding.findingDescription}
                    onChange={(e) => {
                      const newSubsections = [...subsections];
                      newSubsections[sIdx].findings[fIdx].findingDescription = e.target.value;
                      setSubsections(newSubsections);
                    }} required />
                  <input type="text" placeholder="Supporting Authors (comma separated)"
                    onChange={(e) => {
                      const newSubsections = [...subsections];
                      newSubsections[sIdx].findings[fIdx].supportingAuthors = e.target.value.split(',');
                      setSubsections(newSubsections);
                    }} required />
                </div>
              ))}
            </div>
          ))}
        </div>
        <button type="submit">Update Chapter</button>
      </form>
    </div>
  );
};

export default EditChapter;
```

### API Service

Create a file `src/services/api.js` to handle API requests.

**src/services/api.js**:

```javascript
import axios from 'axios';

const API_URL = 'http://localhost:8000';

export const getChapters = async () => {
  const response = await axios.get(`${API_URL}/chapters`);
  return response.data;
};

export const getChapter = async (id) => {
  const response = await axios.get(`${API_URL}/chapter/${id}`);
  return response.data;
};

export const createChapter = async (chapter) => {
  const response = await axios.post(`${API_URL}/chapter`, chapter);
  return response.data;
};

export const updateChapter = async (id, chapter) => {
  const response = await axios.put(`${API_URL}/chapter/${id}`, chapter);
  return response.data;
};

export const getTags = async () => {
  const response = await axios.get(`${API_URL}/tags`);
  return response.data;
};
```

### App Component

**src/App.js**:

```javascript
import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import './App.css';
import ChaptersList from './components/ChaptersList';
import AddChapter from './components/AddChapter';
import EditChapter from './components/EditChapter';
import { getChapters, getTags } from './services/api';

function App() {
  const [chapters, setChapters] = useState([]);
  const [tags, setTags] = useState([]);

  const fetchChapters = async () => {
    const chaptersData = await getChapters();
    setChapters(chaptersData);
  };

  const fetchTags = async () => {
    const tagsData = await getTags();
    setTags(tagsData);
  };

  useEffect(() => {
    fetchChapters();
    fetchTags();
  }, []);

  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <h1>PhD Summary</h1>
        </header>
        <Routes>
          <Route path="/" element={
            <>
              <AddChapter onAddChapter={fetchChapters} />
              <ChaptersList chapters={chapters} tags={tags} />
            </>
          } />
          <Route path="/edit/:id" element={<EditChapter />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
```

## Running the Application

### Backend

1. **Ensure MongoDB is running**:

   ```sh
   mongod
   ```

2. **Run the Go backend**:

   ```sh
   go run main.go
   ```

### Frontend

1. **Start the React development server**:

   ```sh
   npm start
   ```

2. **Open your browser and navigate to

 `http://localhost:3000`**.

## Common Issues and Troubleshooting

### Jekyll Installation Issues

#### Gem::FilePermissionError

If you encounter a `Gem::FilePermissionError`, use a Ruby version manager like `rbenv` to avoid permission issues.

1. **Install `rbenv` and `ruby-build`**:
   ```sh
   brew install rbenv ruby-build
   ```

2. **Set Up `rbenv`**:
   ```sh
   echo 'eval "$(rbenv init -)"' >> ~/.zshrc
   source ~/.zshrc
   ```

3. **Install a Newer Version of Ruby**:
   ```sh
   rbenv install 3.1.0
   rbenv global 3.1.0
   ```

4. **Install Jekyll and Bundler**:
   ```sh
   gem install jekyll bundler
   ```

### Liquid Syntax Errors

If you encounter Liquid syntax errors due to the use of `{{` and `}}` in your code blocks, wrap the affected code with `{% raw %}` and `{% endraw %}`.

#### Example:

```markdown
{% raw %}
```go
pipeline := mongo.Pipeline{
    {{"$unwind", bson.D{{"path", "$tags"}}}},
    {{"$group", bson.D{{"_id", nil}, {"tags", bson.D{{"$addToSet", "$tags"}}}}}},
    {{"$project", bson.D{{"_id", 0}, {"tags", 1}}}},
}
```
{% endraw %}
```

## Conclusion

By following this guide, you can create a full-stack application to manage PhD summaries using Go, MongoDB, and React.js. This documentation also provides solutions to common issues encountered during setup and configuration, ensuring a smooth development experience.

---

You can use this documentation as a Markdown file in your Jekyll site. Ensure you include the `{% raw %}` and `{% endraw %}` tags around code blocks that might be misinterpreted by Liquid to avoid syntax errors. If you encounter any further issues, feel free to ask for more assistance.