package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Campeon struct {
	ID             int      `json:"id"`
	Nombre         string   `json:"nombre"`
	Rol            string   `json:"rol"`
	Dificultad     int      `json:"dificultad"`
	Region         string   `json:"region"`
	Recurso        string   `json:"recurso"`
	AñoLanzamiento int      `json:"año_lanzamiento"`
	Habilidades    []string `json:"habilidades"`
}

type RespuestaError struct {
	Error   bool   `json:"error"`
	Mensaje string `json:"mensaje"`
}

func responderError(w http.ResponseWriter, codigo int, mensaje string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(RespuestaError{
		Error:   true,
		Mensaje: mensaje,
	})
}

func cargarCampeones() ([]Campeon, error) {
	datos, err := os.ReadFile("champions.json")
	if err != nil {
		return nil, err
	}

	var campeones []Campeon
	err = json.Unmarshal(datos, &campeones)
	return campeones, err
}

func guardarCampeones(campeones []Campeon) error {
	datos, err := json.MarshalIndent(campeones, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("champions.json", datos, 0644)
}

func validarCampeon(c Campeon) string {
	if c.Nombre == "" {
		return "El nombre es obligatorio"
	}
	if c.Rol == "" {
		return "El rol es obligatorio"
	}
	if c.Dificultad < 1 || c.Dificultad > 3 {
		return "La dificultad debe estar entre 1 y 3"
	}
	if c.AñoLanzamiento < 2000 {
		return "El año de lanzamiento no es válido"
	}
	if len(c.Habilidades) == 0 {
		return "Debe tener al menos una habilidad"
	}
	return ""
}

func obtenerTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	campeones, err := cargarCampeones()
	if err != nil {
		responderError(w, 500, "Error al cargar los datos")
		return
	}

	// Filtros combinados
	rol := r.URL.Query().Get("rol")
	region := r.URL.Query().Get("region")

	var filtrados []Campeon

	for _, c := range campeones {
		if rol != "" && strings.ToLower(c.Rol) != strings.ToLower(rol) {
			continue
		}
		if region != "" && strings.ToLower(c.Region) != strings.ToLower(region) {
			continue
		}
		filtrados = append(filtrados, c)
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(filtrados)
}

func obtenerPorID(w http.ResponseWriter, id int) {
	w.Header().Set("Content-Type", "application/json")

	campeones, err := cargarCampeones()
	if err != nil {
		responderError(w, 500, "Error al cargar datos")
		return
	}

	for _, c := range campeones {
		if c.ID == id {
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "  ")
			encoder.Encode(c)
			return
		}
	}

	responderError(w, 404, "Campeón no encontrado")
}

func crear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var nuevo Campeon
	if err := json.NewDecoder(r.Body).Decode(&nuevo); err != nil {
		responderError(w, 400, "JSON inválido")
		return
	}

	if mensaje := validarCampeon(nuevo); mensaje != "" {
		responderError(w, 400, mensaje)
		return
	}

	campeones, err := cargarCampeones()
	if err != nil {
		responderError(w, 500, "Error al cargar datos")
		return
	}

	nuevo.ID = campeones[len(campeones)-1].ID + 1
	campeones = append(campeones, nuevo)

	if err := guardarCampeones(campeones); err != nil {
		responderError(w, 500, "Error al guardar datos")
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(nuevo)
}

func actualizar(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	var actualizado Campeon
	if err := json.NewDecoder(r.Body).Decode(&actualizado); err != nil {
		responderError(w, 400, "JSON inválido")
		return
	}

	if mensaje := validarCampeon(actualizado); mensaje != "" {
		responderError(w, 400, mensaje)
		return
	}

	campeones, err := cargarCampeones()
	if err != nil {
		responderError(w, 500, "Error al cargar datos")
		return
	}

	for i, c := range campeones {
		if c.ID == id {
			actualizado.ID = id
			campeones[i] = actualizado
			guardarCampeones(campeones)

			json.NewEncoder(w).Encode(actualizado)
			return
		}
	}

	responderError(w, 404, "Campeón no encontrado")
}

func eliminar(w http.ResponseWriter, id int) {
	w.Header().Set("Content-Type", "application/json")

	campeones, err := cargarCampeones()
	if err != nil {
		responderError(w, 500, "Error al cargar datos")
		return
	}

	var nuevos []Campeon
	encontrado := false

	for _, c := range campeones {
		if c.ID == id {
			encontrado = true
			continue
		}
		nuevos = append(nuevos, c)
	}

	if !encontrado {
		responderError(w, 404, "Campeón no encontrado")
		return
	}

	guardarCampeones(nuevos)

	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Campeón eliminado correctamente",
	})
}

func main() {
	http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {

		partes := strings.Split(r.URL.Path, "/")

		// PATH PARAMETER
		if len(partes) == 4 {
			id, _ := strconv.Atoi(partes[3])

			switch r.Method {
			case http.MethodGet:
				obtenerPorID(w, id)
			case http.MethodPut:
				actualizar(w, r, id)
			case http.MethodDelete:
				eliminar(w, id)
			default:
				responderError(w, 405, "Método no permitido")
			}
			return
		}

		switch r.Method {
		case http.MethodGet:
			obtenerTodos(w, r)
		case http.MethodPost:
			crear(w, r)
		default:
			responderError(w, 405, "Método no permitido")
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "24730"
	}

	fmt.Println("Servidor ejecutándose en el puerto " + port)
	http.ListenAndServe(":"+port, nil)
}
