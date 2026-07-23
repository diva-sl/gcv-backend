package controllers

import (
	"context"
	"net/http"
	"time"

	"gcv-backend/config"
	"gcv-backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetProjects lists all case studies in the GCV collection
func GetProjects(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.GetCollection("gcv")
	var projects []models.Project

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch showcases"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var project models.Project
		if err := cursor.Decode(&project); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed decoding database records"})
			return
		}
		projects = append(projects, project)
	}

	// Return empty array instead of null if empty
	if projects == nil {
		projects = []models.Project{}
	}

	c.JSON(http.StatusOK, projects)
}

// GetProjectByID retrieves a single project using the unique ProjectID
func GetProjectByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectID := c.Param("projectId")
	collection := config.GetCollection("gcv")

	var project models.Project
	err := collection.FindOne(ctx, bson.M{"projectId": projectID}).Decode(&project)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// CreateProject inserts a new case study
func CreateProject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newProject models.Project
	if err := c.ShouldBindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := config.GetCollection("gcv")

	// Ensure unique projectId check
	count, _ := collection.CountDocuments(ctx, bson.M{"projectId": newProject.ProjectID})
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A project with this project ID already exists"})
		return
	}

	_, err := collection.InsertOne(ctx, newProject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Project created successfully", "project": newProject})
}

// UpdateProject modifies an existing case study fields or mock images list
func UpdateProject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectID := c.Param("projectId")
	var updatedProject models.Project
	if err := c.ShouldBindJSON(&updatedProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := config.GetCollection("gcv")

	// Create bson updates document
	update := bson.M{
		"$set": bson.M{
			"title":         updatedProject.Title,
			"client":        updatedProject.Client,
			"category":      updatedProject.Category,
			"description":   updatedProject.Description,
			"outcome":       updatedProject.Outcome,
			"tags":          updatedProject.Tags,
			"image":         updatedProject.Image,
			"challenge":     updatedProject.Challenge,
			"solution":      updatedProject.Solution,
			"architecture":  updatedProject.Architecture,
			"metrics":       updatedProject.Metrics,
			"siteUrl":       updatedProject.SiteUrl,
			"adminUrl":      updatedProject.AdminUrl,
			"desktopMockup": updatedProject.DesktopMockup,
			"tabletMockup":  updatedProject.TabletMockup,
			"mobileMockup":  updatedProject.MobileMockup,
			"screenshots":   updatedProject.Screenshots,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"projectId": projectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project records"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Project updated successfully"})
}

// DeleteProject deletes a case study document
func DeleteProject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectID := c.Param("projectId")
	collection := config.GetCollection("gcv")

	result, err := collection.DeleteOne(ctx, bson.M{"projectId": projectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove project"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Project deleted successfully"})
}