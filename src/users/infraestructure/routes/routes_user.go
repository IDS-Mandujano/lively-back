package routes

import (
	"lively-backend/src/core/database"
	usecases "lively-backend/src/users/application/useCases"
	"lively-backend/src/users/infraestructure/controllers"
	"lively-backend/src/users/infraestructure/mysql"
	"net/http"

)

func SetupUserRoutes(mux *http.ServeMux) {

	userRepo := mysql.NewMySQLUserRepository(database.DB)
	registerUserUC := usecases.NewRegisterUserUseCase(userRepo)
	registerUserCtrl := controllers.NewRegisterUserController(registerUserUC)
	mux.HandleFunc("/api/users/register", registerUserCtrl.Handle)
	
	loginUserUC := usecases.NewLoginUserUseCase(userRepo)
	loginUserCtrl := controllers.NewLoginUserController(loginUserUC)
	mux.HandleFunc("/api/users/login", loginUserCtrl.Handle)
	
}
