# CineApp — API de Gestion de Cine

Backend RESTful desarrollado en **Go** con **Gin** y **MongoDB** para la gestion de un sistema de cine: peliculas, salas, funciones, usuarios y reservaciones.

---

## Estructura del Proyecto

```
cine-api/
├── cmd/
│   └── main.go               # Punto de entrada
├── config/
│   └── database.go           # Conexion a MongoDB
├── models/
│   ├── pelicula.go
│   ├── sala.go
│   ├── funcion.go
│   ├── usuario.go
│   └── reservacion.go
├── repositories/             # Acceso a datos (MongoDB)
│   ├── pelicula_repository.go
│   ├── sala_repository.go
│   ├── funcion_repository.go
│   ├── usuario_repository.go
│   └── reservacion_repository.go
├── services/                 # Logica de negocio
│   ├── pelicula_service.go
│   ├── sala_service.go
│   ├── funcion_service.go
│   ├── usuario_service.go
│   └── reservacion_service.go
├── handlers/                 # Controladores HTTP
│   ├── pelicula_handler.go
│   ├── sala_handler.go
│   ├── funcion_handler.go
│   ├── usuario_handler.go
│   └── reservacion_handler.go
├── routes/
│   └── routes.go             # Registro de rutas
├── consultas.md              # Documentacion de consultas MongoDB
├── .env.example
├── go.mod
└── README.md
```

---

## Variables de Entorno

Copia `.env.example` como `.env` y completa los valores:

```env
MONGODB_URI=mongodb+srv://<usuario>:<password>@cluster0.xxxxx.mongodb.net/
DB_NAME=cine_db
PORT=8080
GIN_MODE=release
```

### Variables en Railway

| Variable | Descripcion |
|---|---|
| `MONGODB_URI` | Connection string de MongoDB Atlas |
| `DB_NAME` | Nombre de la base de datos |
| `PORT` | Puerto del servidor (Railway lo asigna automaticamente) |
| `GIN_MODE` | `release` para produccion |

---

## Ejecucion Local

### Requisitos
- Go 1.21 o superior
- MongoDB Atlas (cuenta gratuita) o MongoDB local

```bash
# 1. Clonar el repositorio
git clone https://github.com/<tu-usuario>/cine-api.git
cd cine-api

# 2. Descargar dependencias (OBLIGATORIO en el primer build)
go mod tidy

# 3. Configurar variables de entorno
cp .env.example .env
# Editar .env con tus credenciales de MongoDB

# 4. Ejecutar
go run cmd/main.go
```

El servidor arranca en `http://localhost:8080`.

> **Nota:** `go mod tidy` descarga las dependencias declaradas en `go.mod`
> (Gin, MongoDB driver, godotenv) y genera el archivo `go.sum`.
> Es obligatorio ejecutarlo antes del primer `go run` o `go build`.

---

## Endpoints

### Peliculas `/api/v1/peliculas`
| Metodo | Ruta | Descripcion |
|---|---|---|
| GET | `/` | Listar todas (paginado) |
| POST | `/` | Crear pelicula |
| GET | `/:id` | Obtener por ID |
| PUT | `/:id` | Actualizar |
| DELETE | `/:id` | Eliminar (logico) |
| GET | `/genero/:genero` | Filtrar por genero |
| GET | `/reporte/popularidad` | Reporte de popularidad |

### Salas `/api/v1/salas`
| Metodo | Ruta | Descripcion |
|---|---|---|
| GET | `/` | Listar todas |
| POST | `/` | Crear sala |
| GET | `/:id` | Obtener por ID |
| GET | `/tipo/:tipo` | Filtrar por tipo |
| DELETE | `/:id` | Eliminar (logico) |

### Funciones `/api/v1/funciones`
| Metodo | Ruta | Descripcion |
|---|---|---|
| GET | `/` | Listar todas (paginado) |
| POST | `/` | Crear funcion |
| GET | `/:id` | Obtener por ID |
| GET | `/pelicula/:id` | Funciones de una pelicula |
| GET | `/fecha?desde=&hasta=` | Filtrar por rango de fechas |
| DELETE | `/:id` | Eliminar (logico) |
| GET | `/reporte/salas` | Reporte por sala |

### Usuarios `/api/v1/usuarios`
| Metodo | Ruta | Descripcion |
|---|---|---|
| GET | `/` | Listar todos (paginado) |
| POST | `/` | Registrar usuario |
| GET | `/:id` | Obtener por ID |
| PUT | `/:id` | Actualizar |
| DELETE | `/:id` | Eliminar (logico) |

### Reservaciones `/api/v1/reservaciones`
| Metodo | Ruta | Descripcion |
|---|---|---|
| POST | `/` | Crear reservacion |
| GET | `/:id` | Obtener por ID |
| GET | `/:id/detalle` | Obtener con detalle completo |
| GET | `/usuario/:id` | Reservaciones de un usuario |
| PUT | `/:id/cancelar` | Cancelar reservacion |
| GET | `/reporte/ingresos` | Reporte de ingresos |

---

## Ejemplos de Request

### Crear pelicula
```json
POST /api/v1/peliculas
{
  "titulo": "Dune: Parte Dos",
  "genero": ["Ciencia Ficcion", "Aventura"],
  "duracion": 166,
  "clasificacion": "B",
  "director": "Denis Villeneuve",
  "sinopsis": "Paul Atreides se une a los Fremen...",
  "reparto": [
    { "nombre": "Timothee Chalamet", "personaje": "Paul Atreides" },
    { "nombre": "Zendaya", "personaje": "Chani" }
  ]
}
```

### Crear reservacion
```json
POST /api/v1/reservaciones
{
  "usuario_id": "64a1b2c3d4e5f6a7b8c9d0e1",
  "funcion_id": "64a1b2c3d4e5f6a7b8c9d0e2",
  "asientos": [
    { "fila": "C", "numero": 5 },
    { "fila": "C", "numero": 6 }
  ]
}
```

---

## Deploy en Railway

1. Crear cuenta en [railway.com](https://railway.com)
2. New Project > Deploy from GitHub Repo
3. Seleccionar este repositorio
4. Agregar variables de entorno en el panel de Railway
5. Railway detecta Go automaticamente y realiza el build

---

## Tecnologias

- **Go 1.21**
- **Gin** — Framework HTTP
- **MongoDB Driver** — Driver oficial de MongoDB para Go
- **godotenv** — Manejo de variables de entorno
- **MongoDB Atlas** — Base de datos en la nube
- **Railway** — Plataforma de deploy
