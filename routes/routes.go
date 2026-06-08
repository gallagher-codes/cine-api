package routes

import (
	"cine-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	// Peliculas
	ph := handlers.NewPeliculaHandler()
	peliculas := api.Group("/peliculas")
	{
		peliculas.POST("", ph.Crear)
		peliculas.GET("", ph.Listar)
		peliculas.GET("/reporte/popularidad", ph.ReportePopularidad)
		peliculas.GET("/genero/:genero", ph.ListarPorGenero)
		peliculas.GET("/:id", ph.ObtenerPorID)
		peliculas.PUT("/:id", ph.Actualizar)
		peliculas.DELETE("/:id", ph.Eliminar)
	}

	// Salas
	sh := handlers.NewSalaHandler()
	salas := api.Group("/salas")
	{
		salas.POST("", sh.Crear)
		salas.GET("", sh.Listar)
		salas.GET("/tipo/:tipo", sh.ListarPorTipo)
		salas.GET("/:id", sh.ObtenerPorID)
		salas.DELETE("/:id", sh.Eliminar)
	}

	// Funciones
	fh := handlers.NewFuncionHandler()
	funciones := api.Group("/funciones")
	{
		funciones.POST("", fh.Crear)
		funciones.GET("", fh.Listar)
		funciones.GET("/fecha", fh.ListarPorFecha)
		funciones.GET("/reporte/salas", fh.ReportePorSala)
		funciones.GET("/pelicula/:id", fh.ListarPorPelicula)
		funciones.GET("/:id", fh.ObtenerPorID)
		funciones.DELETE("/:id", fh.Eliminar)
	}

	// Usuarios
	uh := handlers.NewUsuarioHandler()
	usuarios := api.Group("/usuarios")
	{
		usuarios.POST("", uh.Registrar)
		usuarios.GET("", uh.Listar)
		usuarios.GET("/:id", uh.ObtenerPorID)
		usuarios.PUT("/:id", uh.Actualizar)
		usuarios.DELETE("/:id", uh.Eliminar)
	}

	// Reservaciones
	rh := handlers.NewReservacionHandler()
	reservaciones := api.Group("/reservaciones")
	{
		reservaciones.POST("", rh.Crear)
		reservaciones.GET("/reporte/ingresos", rh.ReporteIngresos)
		reservaciones.GET("/usuario/:id", rh.ListarPorUsuario)
		reservaciones.GET("/:id", rh.ObtenerPorID)
		reservaciones.GET("/:id/detalle", rh.ObtenerDetalle)
		reservaciones.PUT("/:id/cancelar", rh.Cancelar)
	}
}
