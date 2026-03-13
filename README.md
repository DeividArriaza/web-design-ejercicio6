# API de Campeones de League of Legends

## 📌 Descripción

Esta es una API REST desarrollada en **Go** que permite gestionar información sobre campeones de League of Legends.

La API permite:

- Obtener todos los campeones
- Obtener un campeón por ID (path parameter)
- Filtrar campeones por rol y región (query parameters combinados)
- Crear nuevos campeones
- Actualizar campeones existentes
- Eliminar campeones
- Persistencia real en archivo JSON
- Manejo estructurado de errores en formato JSON

El servidor se ejecuta en el puerto correspondiente al carnet: **24730**.

---

## 🚀 Tecnologías utilizadas

- Go (net/http – librería estándar)
- JSON para persistencia
- Docker
- HTTP REST

---

## 🐳 Ejecución con Docker

### 1️⃣ Construir la imagen

```bash
docker build -t lol-api .
