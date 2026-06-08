# 📋 Documentación de Consultas — CineApp

Sistema de Gestión de Cine | Base de Datos: MongoDB  
Colecciones: `peliculas`, `salas`, `funciones`, `usuarios`, `reservaciones`

---

## CONSULTAS BÁSICAS

---

### 1. Inserción — Crear una película

**Descripción:** Inserta un nuevo documento en la colección `peliculas` con reparto embebido.  
**Objetivo de negocio:** Registrar una nueva película en cartelera.

**Request:**
```http
POST /api/v1/peliculas
Content-Type: application/json
```
```json
{
  "titulo": "Dune: Parte Dos",
  "genero": ["Ciencia Ficción", "Aventura"],
  "duracion": 166,
  "clasificacion": "B",
  "director": "Denis Villeneuve",
  "sinopsis": "Paul Atreides une fuerzas con los Fremen en su viaje de venganza.",
  "reparto": [
    { "nombre": "Timothée Chalamet", "personaje": "Paul Atreides" },
    { "nombre": "Zendaya", "personaje": "Chani" }
  ]
}
```

**Query MongoDB:**
```javascript
db.peliculas.insertOne({
  titulo: "Dune: Parte Dos",
  genero: ["Ciencia Ficción", "Aventura"],
  duracion: 166,
  clasificacion: "B",
  director: "Denis Villeneuve",
  sinopsis: "Paul Atreides une fuerzas con los Fremen en su viaje de venganza.",
  reparto: [
    { nombre: "Timothée Chalamet", personaje: "Paul Atreides" },
    { nombre: "Zendaya", personaje: "Chani" }
  ],
  activo: true,
  creado_en: new Date(),
  actualizado_en: new Date()
})
```

**Resultado esperado:**
```json
{
  "id": "64a1b2c3d4e5f6a7b8c9d0e1",
  "titulo": "Dune: Parte Dos",
  "genero": ["Ciencia Ficción", "Aventura"],
  "duracion": 166,
  "clasificacion": "B",
  "director": "Denis Villeneuve",
  "activo": true
}
```

---

### 2. Inserción — Registrar un usuario

**Descripción:** Crea un nuevo usuario en el sistema.  
**Objetivo de negocio:** Permitir que clientes se registren para hacer reservaciones.

**Request:**
```http
POST /api/v1/usuarios
Content-Type: application/json
```
```json
{
  "nombre": "María López",
  "email": "maria.lopez@email.com",
  "telefono": "9991234567",
  "password": "segura123"
}
```

**Query MongoDB:**
```javascript
db.usuarios.insertOne({
  nombre: "María López",
  email: "maria.lopez@email.com",
  telefono: "9991234567",
  password: "segura123",
  activo: true,
  creado_en: new Date()
})
```

**Resultado esperado:**
```json
{
  "id": "64a1b2c3d4e5f6a7b8c9d0e2",
  "nombre": "María López",
  "email": "maria.lopez@email.com",
  "telefono": "9991234567",
  "activo": true
}
```

---

### 3. Inserción — Crear una reservación

**Descripción:** Registra una reservación con asientos embebidos y actualiza asientos disponibles en la función.  
**Objetivo de negocio:** Permitir a un usuario comprar boletos para una función.

**Request:**
```http
POST /api/v1/reservaciones
Content-Type: application/json
```
```json
{
  "usuario_id": "64a1b2c3d4e5f6a7b8c9d0e2",
  "funcion_id": "64a1b2c3d4e5f6a7b8c9d0e3",
  "asientos": [
    { "fila": "C", "numero": 5 },
    { "fila": "C", "numero": 6 }
  ]
}
```

**Query MongoDB:**
```javascript
// Paso 1: Insertar reservación
db.reservaciones.insertOne({
  usuario_id: ObjectId("64a1b2c3d4e5f6a7b8c9d0e2"),
  funcion_id: ObjectId("64a1b2c3d4e5f6a7b8c9d0e3"),
  asientos: [
    { fila: "C", numero: 5 },
    { fila: "C", numero: 6 }
  ],
  total: 260.00,
  estado: "confirmada",
  creado_en: new Date()
})

// Paso 2: Descontar asientos disponibles
db.funciones.updateOne(
  { _id: ObjectId("64a1b2c3d4e5f6a7b8c9d0e3") },
  { $inc: { asientos_disponibles: -2 } }
)
```

**Resultado esperado:**
```json
{
  "id": "64a1b2c3d4e5f6a7b8c9d0e4",
  "usuario_id": "64a1b2c3d4e5f6a7b8c9d0e2",
  "funcion_id": "64a1b2c3d4e5f6a7b8c9d0e3",
  "asientos": [{"fila":"C","numero":5},{"fila":"C","numero":6}],
  "total": 260.00,
  "estado": "confirmada"
}
```

---

### 4. Actualización — Modificar datos de una película

**Descripción:** Actualiza campos específicos de una película usando `$set`.  
**Objetivo de negocio:** Corregir o actualizar información de la cartelera.

**Request:**
```http
PUT /api/v1/peliculas/64a1b2c3d4e5f6a7b8c9d0e1
Content-Type: application/json
```
```json
{
  "sinopsis": "Nueva sinopsis actualizada para la película.",
  "clasificacion": "B15"
}
```

**Query MongoDB:**
```javascript
db.peliculas.findOneAndUpdate(
  { _id: ObjectId("64a1b2c3d4e5f6a7b8c9d0e1") },
  { $set: {
      sinopsis: "Nueva sinopsis actualizada para la película.",
      clasificacion: "B15",
      actualizado_en: new Date()
    }
  },
  { returnDocument: "after" }
)
```

**Resultado esperado:** Documento con los campos modificados y `actualizado_en` actualizado.

---

### 5. Eliminación lógica — Desactivar una película

**Descripción:** Marca una película como inactiva en lugar de eliminarla físicamente.  
**Objetivo de negocio:** Conservar el historial sin mostrar la película en cartelera activa.

**Request:**
```http
DELETE /api/v1/peliculas/64a1b2c3d4e5f6a7b8c9d0e1
```

**Query MongoDB:**
```javascript
db.peliculas.updateOne(
  { _id: ObjectId("64a1b2c3d4e5f6a7b8c9d0e1") },
  { $set: { activo: false, actualizado_en: new Date() } }
)
```

**Resultado esperado:**
```json
{ "mensaje": "película eliminada correctamente" }
```

---

### 6. Búsqueda por ID

**Descripción:** Recupera un documento específico por su `_id`.  
**Objetivo de negocio:** Mostrar el detalle de una película en la app del cliente.

**Request:**
```http
GET /api/v1/peliculas/64a1b2c3d4e5f6a7b8c9d0e1
```

**Query MongoDB:**
```javascript
db.peliculas.findOne({
  _id: ObjectId("64a1b2c3d4e5f6a7b8c9d0e1"),
  activo: true
})
```

---

### 7. Listado general con paginación

**Descripción:** Lista todas las películas activas con paginación.  
**Objetivo de negocio:** Mostrar cartelera completa con carga progresiva.

**Request:**
```http
GET /api/v1/peliculas?page=1&limit=10
```

**Query MongoDB:**
```javascript
db.peliculas
  .find({ activo: true })
  .sort({ titulo: 1 })
  .skip(0)
  .limit(10)
```

---

### 8. Cancelar reservación (actualización de estado)

**Descripción:** Cambia el estado de una reservación a `cancelada` y devuelve los asientos.  
**Objetivo de negocio:** Procesar cancelaciones de boletos.

**Request:**
```http
PUT /api/v1/reservaciones/64a1b2c3d4e5f6a7b8c9d0e4/cancelar
```

**Query MongoDB:**
```javascript
// Paso 1: Cancelar reservación
db.reservaciones.updateOne(
  { _id: ObjectId("64a1b2c3d4e5f6a7b8c9d0e4") },
  { $set: { estado: "cancelada" } }
)

// Paso 2: Devolver asientos (incrementar)
db.funciones.updateOne(
  { _id: ObjectId("64a1b2c3d4e5f6a7b8c9d0e3") },
  { $inc: { asientos_disponibles: 2 } }
)
```

---

## CONSULTAS AVANZADAS

---

### 9. Filtro combinado — Películas por género y clasificación

**Descripción:** Busca películas que coincidan con género Y clasificación simultáneamente.  
**Objetivo de negocio:** Filtrar cartelera para recomendaciones segmentadas por audiencia.

**Request:**
```http
GET /api/v1/peliculas?genero=Acción&clasificacion=B
```

**Query MongoDB:**
```javascript
db.peliculas.find({
  genero: "Acción",
  clasificacion: "B",
  activo: true
}).sort({ titulo: 1 })
```

---

### 10. Búsqueda por texto — Título parcial

**Descripción:** Búsqueda flexible por fragmento del título usando expresión regular.  
**Objetivo de negocio:** Barra de búsqueda en la app del cliente.

**Query MongoDB:**
```javascript
db.peliculas.find({
  titulo: { $regex: "dune", $options: "i" },
  activo: true
})
```

---

### 11. Filtro por rango de fechas — Funciones disponibles

**Descripción:** Obtiene funciones dentro de un rango de fechas con asientos disponibles.  
**Objetivo de negocio:** Mostrar al cliente las funciones que puede reservar.

**Request:**
```http
GET /api/v1/funciones/fecha?desde=2025-06-01&hasta=2025-06-30
```

**Query MongoDB:**
```javascript
db.funciones.find({
  activa: true,
  asientos_disponibles: { $gt: 0 },
  fecha: {
    $gte: ISODate("2025-06-01T00:00:00Z"),
    $lte: ISODate("2025-06-30T23:59:59Z")
  }
}).sort({ fecha: 1 })
```

---

### 12. Paginación avanzada con conteo total

**Descripción:** Paginación con información del total de documentos.  
**Objetivo de negocio:** Componente de paginación en el frontend de administración.

**Query MongoDB:**
```javascript
// Datos de la página
const pagina = 2;
const limite = 10;
db.peliculas
  .find({ activo: true })
  .skip((pagina - 1) * limite)
  .limit(limite)

// Total de documentos
db.peliculas.countDocuments({ activo: true })
```

---

### 13. Aggregation — Reporte de ingresos por película

**Descripción:** Pipeline de agregación que calcula ingresos totales y boletos vendidos por película, relacionando reservaciones → funciones → películas.  
**Objetivo de negocio:** Reporte financiero para la gerencia del cine.

**Request:**
```http
GET /api/v1/reservaciones/reporte/ingresos
```

**Query MongoDB:**
```javascript
db.reservaciones.aggregate([
  { $match: { estado: "confirmada" } },
  {
    $lookup: {
      from: "funciones",
      localField: "funcion_id",
      foreignField: "_id",
      as: "funcion"
    }
  },
  { $unwind: "$funcion" },
  {
    $lookup: {
      from: "peliculas",
      localField: "funcion.pelicula_id",
      foreignField: "_id",
      as: "pelicula"
    }
  },
  { $unwind: "$pelicula" },
  {
    $group: {
      _id: "$pelicula.titulo",
      total_reservaciones: { $sum: 1 },
      ingresos_totales: { $sum: "$total" },
      total_boletos: { $sum: { $size: "$asientos" } }
    }
  },
  { $sort: { ingresos_totales: -1 } }
])
```

**Resultado esperado:**
```json
[
  {
    "_id": "Dune: Parte Dos",
    "total_reservaciones": 48,
    "ingresos_totales": 12480.00,
    "total_boletos": 96
  },
  {
    "_id": "Avengers: Endgame",
    "total_reservaciones": 35,
    "ingresos_totales": 9100.00,
    "total_boletos": 70
  }
]
```

---

### 14. Aggregation — Top 10 películas más populares

**Descripción:** Cuenta reservaciones por película usando lookup encadenado.  
**Objetivo de negocio:** Destacar las películas más exitosas en la cartelera.

**Request:**
```http
GET /api/v1/peliculas/reporte/popularidad
```

**Query MongoDB:**
```javascript
db.peliculas.aggregate([
  {
    $lookup: {
      from: "funciones",
      localField: "_id",
      foreignField: "pelicula_id",
      as: "funciones"
    }
  },
  {
    $lookup: {
      from: "reservaciones",
      localField: "funciones._id",
      foreignField: "funcion_id",
      as: "reservaciones"
    }
  },
  {
    $project: {
      titulo: 1,
      genero: 1,
      total_reservaciones: { $size: "$reservaciones" }
    }
  },
  { $sort: { total_reservaciones: -1 } },
  { $limit: 10 }
])
```

---

### 15. Aggregation — Funciones por sala con ingresos

**Descripción:** Agrupa funciones por sala y calcula el total de ingresos mediante $lookup + $group.  
**Objetivo de negocio:** Evaluar el desempeño y rentabilidad de cada sala.

**Request:**
```http
GET /api/v1/funciones/reporte/salas
```

**Query MongoDB:**
```javascript
db.funciones.aggregate([
  { $match: { activa: true } },
  {
    $lookup: {
      from: "salas",
      localField: "sala_id",
      foreignField: "_id",
      as: "sala"
    }
  },
  { $unwind: "$sala" },
  {
    $group: {
      _id: "$sala.nombre",
      total_funciones: { $sum: 1 },
      ingreso_potencial: { $sum: "$precio" }
    }
  },
  { $sort: { total_funciones: -1 } }
])
```

---

### 16. Lookup completo — Detalle de reservación

**Descripción:** Obtiene una reservación con todos sus datos relacionados: usuario, función, película y sala en una sola consulta.  
**Objetivo de negocio:** Vista de confirmación de compra con todos los detalles.

**Request:**
```http
GET /api/v1/reservaciones/64a1b2c3d4e5f6a7b8c9d0e4/detalle
```

**Query MongoDB:**
```javascript
db.reservaciones.aggregate([
  { $match: { _id: ObjectId("64a1b2c3d4e5f6a7b8c9d0e4") } },
  {
    $lookup: {
      from: "funciones",
      localField: "funcion_id",
      foreignField: "_id",
      as: "funcion"
    }
  },
  { $unwind: "$funcion" },
  {
    $lookup: {
      from: "peliculas",
      localField: "funcion.pelicula_id",
      foreignField: "_id",
      as: "pelicula"
    }
  },
  { $unwind: "$pelicula" },
  {
    $lookup: {
      from: "salas",
      localField: "funcion.sala_id",
      foreignField: "_id",
      as: "sala"
    }
  },
  { $unwind: "$sala" },
  {
    $lookup: {
      from: "usuarios",
      localField: "usuario_id",
      foreignField: "_id",
      as: "usuario"
    }
  },
  { $unwind: "$usuario" },
  { $project: { "usuario.password": 0 } }
])
```

**Resultado esperado:**
```json
{
  "_id": "64a1b2c3d4e5f6a7b8c9d0e4",
  "estado": "confirmada",
  "total": 260.00,
  "asientos": [{"fila":"C","numero":5},{"fila":"C","numero":6}],
  "usuario": { "nombre": "María López", "email": "maria.lopez@email.com" },
  "funcion": { "fecha": "2025-06-15T20:00:00Z", "precio": 130.00 },
  "pelicula": { "titulo": "Dune: Parte Dos", "duracion": 166 },
  "sala": { "nombre": "Sala IMAX 1", "tipo": "IMAX" }
}
```

---

### 17. Aggregation — Estadísticas generales del sistema

**Descripción:** Calcula estadísticas de ocupación: reservaciones por estado, asientos ocupados vs disponibles.  
**Objetivo de negocio:** Dashboard de administración con KPIs del negocio.

**Query MongoDB:**
```javascript
db.reservaciones.aggregate([
  {
    $group: {
      _id: "$estado",
      cantidad: { $sum: 1 },
      boletos: { $sum: { $size: "$asientos" } },
      ingresos: { $sum: "$total" }
    }
  }
])
```

---

### 18. Conteo — Total de usuarios registrados hoy

**Descripción:** Cuenta usuarios registrados en el día actual.  
**Objetivo de negocio:** Monitoreo diario de nuevos registros.

**Query MongoDB:**
```javascript
const hoy = new Date();
hoy.setHours(0, 0, 0, 0);
const manana = new Date(hoy);
manana.setDate(manana.getDate() + 1);

db.usuarios.countDocuments({
  activo: true,
  creado_en: { $gte: hoy, $lt: manana }
})
```
