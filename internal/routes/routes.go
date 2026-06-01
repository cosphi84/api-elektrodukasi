package routes

import (
	"net/http"

	"elektrod/internal/handlers"
	"elektrod/internal/middleware"
	"elektrod/internal/services"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// Initialize services
	userService := services.NewUserService(db)
	categoryService := services.NewCategoryService(db)
	tagService := services.NewTagService(db)
	articleService := services.NewArticleService(db)
	commentService := services.NewCommentService(db)
	projectService := services.NewProjectService(db)
	authService := services.NewAuthService(userService)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	tagHandler := handlers.NewTagHandler(tagService)
	articleHandler := handlers.NewArticleHandler(articleService)
	commentHandler := handlers.NewCommentHandler(commentService)
	projectHandler := handlers.NewProjectHandler(projectService)
	authHandler := handlers.NewAuthHandler(authService)

	// Public routes (no authentication required)
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("GET /api/articles", articleHandler.ListArticles)
	mux.HandleFunc("GET /api/articles/{id}", articleHandler.GetArticle)
	mux.HandleFunc("GET /api/articles/slug/{slug}", articleHandler.GetArticleBySlug)
	mux.HandleFunc("GET /api/categories/{categoryId}/articles", articleHandler.GetArticlesByCategory)
	mux.HandleFunc("GET /api/categories", categoryHandler.ListCategories)
	mux.HandleFunc("GET /api/categories/{id}", categoryHandler.GetCategory)
	mux.HandleFunc("GET /api/tags", tagHandler.ListTags)
	mux.HandleFunc("GET /api/tags/{id}", tagHandler.GetTag)
	mux.HandleFunc("GET /api/articles/{articleId}/comments", commentHandler.GetCommentsByArticle)
	mux.HandleFunc("GET /api/comments/{id}", commentHandler.GetComment)
	mux.HandleFunc("GET /api/comments/{id}/replies", commentHandler.GetReplies)
	mux.HandleFunc("GET /api/projects", projectHandler.ListProjects)
	mux.HandleFunc("GET /api/projects/{id}", projectHandler.GetProject)
	mux.HandleFunc("GET /api/projects/slug/{slug}", projectHandler.GetProjectBySlug)

	// Protected routes - Admin only
	adminUserCreatePost := middleware.Chain(userHandler.CreateUser, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminUserList := middleware.Chain(userHandler.ListUsers, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminUserGet := middleware.Chain(userHandler.GetUser, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminUserUpdate := middleware.Chain(userHandler.UpdateUser, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminUserDelete := middleware.Chain(userHandler.DeleteUser, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminUserDeactivate := middleware.Chain(userHandler.DeactivateUser, middleware.AuthMiddleware, middleware.RequireAdmin)

	mux.HandleFunc("POST /api/users", adminUserCreatePost)
	mux.HandleFunc("GET /api/users", adminUserList)
	mux.HandleFunc("GET /api/users/{id}", adminUserGet)
	mux.HandleFunc("PUT /api/users/{id}", adminUserUpdate)
	mux.HandleFunc("DELETE /api/users/{id}", adminUserDelete)
	mux.HandleFunc("PATCH /api/users/{id}/deactivate", adminUserDeactivate)

	// Protected routes - Admin only for categories
	adminCategoryCreate := middleware.Chain(categoryHandler.CreateCategory, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminCategoryUpdate := middleware.Chain(categoryHandler.UpdateCategory, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminCategoryDelete := middleware.Chain(categoryHandler.DeleteCategory, middleware.AuthMiddleware, middleware.RequireAdmin)

	mux.HandleFunc("POST /api/categories", adminCategoryCreate)
	mux.HandleFunc("PUT /api/categories/{id}", adminCategoryUpdate)
	mux.HandleFunc("DELETE /api/categories/{id}", adminCategoryDelete)

	// Protected routes - Admin only for tags
	adminTagCreate := middleware.Chain(tagHandler.CreateTag, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminTagUpdate := middleware.Chain(tagHandler.UpdateTag, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminTagDelete := middleware.Chain(tagHandler.DeleteTag, middleware.AuthMiddleware, middleware.RequireAdmin)

	mux.HandleFunc("POST /api/tags", adminTagCreate)
	mux.HandleFunc("PUT /api/tags/{id}", adminTagUpdate)
	mux.HandleFunc("DELETE /api/tags/{id}", adminTagDelete)

	// Protected routes - Admin only for articles
	adminArticleCreate := middleware.Chain(articleHandler.CreateArticle, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminArticleUpdate := middleware.Chain(articleHandler.UpdateArticle, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminArticlePublish := middleware.Chain(articleHandler.PublishArticle, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminArticleAddTag := middleware.Chain(articleHandler.AddTag, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminArticleRemoveTag := middleware.Chain(articleHandler.RemoveTag, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminArticleDelete := middleware.Chain(articleHandler.DeleteArticle, middleware.AuthMiddleware, middleware.RequireAdmin)

	mux.HandleFunc("POST /api/articles", adminArticleCreate)
	mux.HandleFunc("PUT /api/articles/{id}", adminArticleUpdate)
	mux.HandleFunc("PATCH /api/articles/{id}/publish", adminArticlePublish)
	mux.HandleFunc("POST /api/articles/{id}/tags", adminArticleAddTag)
	mux.HandleFunc("DELETE /api/articles/{id}/tags/{tagId}", adminArticleRemoveTag)
	mux.HandleFunc("DELETE /api/articles/{id}", adminArticleDelete)

	// Protected routes - Users and Admin can write comments
	userCommentCreate := middleware.Chain(commentHandler.CreateComment, middleware.AuthMiddleware, middleware.InferAuthorizationFromMethod(middleware.ResourceComment))
	userCommentUpdate := middleware.Chain(commentHandler.UpdateComment, middleware.AuthMiddleware, middleware.InferAuthorizationFromMethod(middleware.ResourceComment))
	userCommentDelete := middleware.Chain(commentHandler.DeleteComment, middleware.AuthMiddleware, middleware.InferAuthorizationFromMethod(middleware.ResourceComment))
	adminCommentApprove := middleware.Chain(commentHandler.ApproveComment, middleware.AuthMiddleware, middleware.RequireAdmin)

	mux.HandleFunc("POST /api/comments", userCommentCreate)
	mux.HandleFunc("PUT /api/comments/{id}", userCommentUpdate)
	mux.HandleFunc("DELETE /api/comments/{id}", userCommentDelete)
	mux.HandleFunc("PATCH /api/comments/{id}/approve", adminCommentApprove)

	// Protected routes - Admin only for projects
	adminProjectCreate := middleware.Chain(projectHandler.CreateProject, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminProjectUpdate := middleware.Chain(projectHandler.UpdateProject, middleware.AuthMiddleware, middleware.RequireAdmin)
	adminProjectDelete := middleware.Chain(projectHandler.DeleteProject, middleware.AuthMiddleware, middleware.RequireAdmin)
	userProjectList := middleware.Chain(projectHandler.GetProjectsByOwner, middleware.AuthMiddleware, middleware.RequireUser)

	mux.HandleFunc("POST /api/projects", adminProjectCreate)
	mux.HandleFunc("PUT /api/projects/{id}", adminProjectUpdate)
	mux.HandleFunc("DELETE /api/projects/{id}", adminProjectDelete)
	mux.HandleFunc("GET /api/users/{ownerId}/projects", userProjectList)

	return mux
}
