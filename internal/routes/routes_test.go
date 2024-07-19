package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"githubECS/models"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.Repository{}, &models.Commit{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	SetupRouter(r, db)
	return r
}

func TestSearchRepos(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	r := setupTestRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/search?query=test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetRepositories(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := models.Repository{
		FullName:    "test-repo",
		Description: "A test repository",
	}
	db.Create(&repo)

	r := setupTestRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/repositories/test-repo", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var repos []models.Repository
	err = db.Find(&repos).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, repos)
}

func TestGetCommits(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := models.Repository{
		FullName: "test-repo",
	}
	db.Create(&repo)

	commit := models.Commit{
		Message:      "Initial commit",
		AuthorName:   "John Doe",
		AuthorEmail:  "john.doe@example.com",
		URL:          "https://github.com/test/test-repo/commit/1",
		RepositoryID: repo.ID,
	}
	db.Create(&commit)

	r := setupTestRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/repositories/test-repo/commits", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var commits []models.Commit
	err = db.Where("repository_id = ?", repo.ID).Find(&commits).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, commits)
}

func TestSearchByLanguage(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := models.Repository{
		FullName: "test-repo",
		Language: "go",
	}
	db.Create(&repo)

	r := setupTestRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?language=go", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var repos []models.Repository
	err = db.Where("LOWER(language) = ?", "go").Find(&repos).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, repos)
}

func TestGetTopRepositories(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo1 := models.Repository{
		FullName:   "test-repo1",
		StarsCount: 100,
	}
	repo2 := models.Repository{
		FullName:   "test-repo2",
		StarsCount: 200,
	}
	db.Create(&repo1)
	db.Create(&repo2)

	r := setupTestRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/top?n=1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var repos []models.Repository
	err = db.Order("stars_count desc").Limit(1).Find(&repos).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, repos)
	assert.Equal(t, repo2.FullName, repos[0].FullName)
}
