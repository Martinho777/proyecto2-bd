package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *pgxpool.Pool
var dbORM *gorm.DB

var sesiones = map[string]Sesion{}
var sesionesMu sync.RWMutex

type Producto struct {
	IDProducto  int     `json:"id_producto" gorm:"column:id_producto;primaryKey;autoIncrement"`
	Nombre      string  `json:"nombre" gorm:"column:nombre"`
	Descripcion string  `json:"descripcion" gorm:"column:descripcion"`
	Precio      float64 `json:"precio" gorm:"column:precio"`
	Stock       int     `json:"stock" gorm:"column:stock"`
	IDCategoria int     `json:"id_categoria" gorm:"column:id_categoria"`
	IDProveedor int     `json:"id_proveedor" gorm:"column:id_proveedor"`
}

func (Producto) TableName() string {
	return "producto"
}

type VentaDetallada struct {
	IDVenta        int       `json:"id_venta"`
	Fecha          time.Time `json:"fecha"`
	Cliente        string    `json:"cliente"`
	Empleado       string    `json:"empleado"`
	Producto       string    `json:"producto"`
	Cantidad       int       `json:"cantidad"`
	PrecioUnitario float64   `json:"precio_unitario"`
	Subtotal       float64   `json:"subtotal"`
}

type ReporteProducto struct {
	Producto              string  `json:"producto"`
	TotalUnidadesVendidas int     `json:"total_unidades_vendidas"`
	TotalIngresos         float64 `json:"total_ingresos"`
}

type Cliente struct {
	IDCliente int    `json:"id_cliente"`
	Nombre    string `json:"nombre"`
	Telefono  string `json:"telefono"`
	Correo    string `json:"correo"`
}

type VentaRequest struct {
	IDCliente  int `json:"id_cliente"`
	IDEmpleado int `json:"id_empleado"`
	IDProducto int `json:"id_producto"`
	Cantidad   int `json:"cantidad"`
}

type VentaResponse struct {
	Mensaje        string  `json:"mensaje"`
	IDVenta        int     `json:"id_venta"`
	IDProducto     int     `json:"id_producto"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
}

type LoginRequest struct {
	Correo   string `json:"correo"`
	Password string `json:"password"`
}

type UsuarioSesion struct {
	IDUsuario int    `json:"id_usuario"`
	Nombre    string `json:"nombre"`
	Correo    string `json:"correo"`
	Rol       string `json:"rol"`
}

type LoginResponse struct {
	Mensaje string `json:"mensaje"`
	Token   string `json:"token"`
	Nombre  string `json:"nombre"`
	Correo  string `json:"correo"`
	Rol     string `json:"rol"`
}

type Sesion struct {
	Usuario UsuarioSesion
	Creada  time.Time
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://proy3:secret@127.0.0.1:5433/tienda_db?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	db, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Error creando pool de conexión: %v", err)
	}

	if err = db.Ping(ctx); err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}

	fmt.Println("Conexión a PostgreSQL exitosa")

	dbORM, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error conectando con GORM: %v", err)
	}

	fmt.Println("Conexión a PostgreSQL con GORM exitosa")

	http.HandleFunc("/", enableCORS(inicioHandler))

	http.HandleFunc("/login", enableCORS(loginHandler))
	http.HandleFunc("/logout", enableCORS(logoutHandler))
	http.HandleFunc("/me", enableCORS(meHandler))

	http.HandleFunc("/productos", enableCORS(productosHandler))
	http.HandleFunc("/ventas-detalladas", enableCORS(ventasDetalladasHandler))
	http.HandleFunc("/reporte/productos", enableCORS(reporteProductosHandler))
	http.HandleFunc("/clientes", enableCORS(clientesHandler))
	http.HandleFunc("/ventas", enableCORS(ventasHandler))

	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func inicioHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respuesta := map[string]string{
		"mensaje": "Backend funcionando correctamente",
	}
	json.NewEncoder(w).Encode(respuesta)
}

func generarToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func obtenerSesion(r *http.Request) (UsuarioSesion, bool) {
	auth := r.Header.Get("Authorization")

	if !strings.HasPrefix(auth, "Bearer ") {
		return UsuarioSesion{}, false
	}

	token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))

	sesionesMu.RLock()
	defer sesionesMu.RUnlock()

	sesion, existe := sesiones[token]
	if !existe {
		return UsuarioSesion{}, false
	}

	return sesion.Usuario, true
}

func exigirRol(w http.ResponseWriter, r *http.Request, rolesPermitidos ...string) bool {
	usuario, ok := obtenerSesion(r)
	if !ok {
		http.Error(w, "Debes iniciar sesión", http.StatusUnauthorized)
		return false
	}

	if usuario.Rol == "admin" {
		return true
	}

	for _, rol := range rolesPermitidos {
		if usuario.Rol == rol {
			return true
		}
	}

	http.Error(w, "No tienes permiso para esta operación", http.StatusForbidden)
	return false
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	var usuario UsuarioSesion

	err = db.QueryRow(context.Background(), `
		SELECT id_usuario, nombre, correo, rol
		FROM app_usuario
		WHERE correo = $1 AND password = $2
	`, req.Correo, req.Password).Scan(
		&usuario.IDUsuario,
		&usuario.Nombre,
		&usuario.Correo,
		&usuario.Rol,
	)

	if err != nil {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	token, err := generarToken()
	if err != nil {
		http.Error(w, "Error generando sesión", http.StatusInternalServerError)
		return
	}

	sesionesMu.Lock()
	sesiones[token] = Sesion{
		Usuario: usuario,
		Creada:  time.Now(),
	}
	sesionesMu.Unlock()

	json.NewEncoder(w).Encode(LoginResponse{
		Mensaje: "Login correcto",
		Token:   token,
		Nombre:  usuario.Nombre,
		Correo:  usuario.Correo,
		Rol:     usuario.Rol,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	auth := r.Header.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))

	if token != "" {
		sesionesMu.Lock()
		delete(sesiones, token)
		sesionesMu.Unlock()
	}

	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Logout correcto",
	})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usuario, ok := obtenerSesion(r)
	if !ok {
		http.Error(w, "No hay sesión activa", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(usuario)
}

func listarProductos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var productos []Producto

	err := dbORM.Order("id_producto").Find(&productos).Error
	if err != nil {
		http.Error(w, "Error consultando productos con ORM: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error ORM listando productos:", err)
		return
	}

	json.NewEncoder(w).Encode(productos)
}

func crearProducto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p Producto
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = dbORM.Create(&p).Error
	if err != nil {
		http.Error(w, "Error creando producto con ORM: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func actualizarProducto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p Producto
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if p.IDProducto == 0 {
		http.Error(w, "Debes enviar id_producto", http.StatusBadRequest)
		return
	}

	resultado := dbORM.Model(&Producto{}).
		Where("id_producto = ?", p.IDProducto).
		Updates(map[string]interface{}{
			"nombre":       p.Nombre,
			"descripcion":  p.Descripcion,
			"precio":       p.Precio,
			"stock":        p.Stock,
			"id_categoria": p.IDCategoria,
			"id_proveedor": p.IDProveedor,
		})

	if resultado.Error != nil {
		http.Error(w, "Error actualizando producto con ORM: "+resultado.Error.Error(), http.StatusInternalServerError)
		return
	}

	if resultado.RowsAffected == 0 {
		http.Error(w, "No se encontró el producto", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Producto actualizado correctamente con ORM",
	})
}

func eliminarProducto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Debes enviar el id del producto", http.StatusBadRequest)
		return
	}

	resultado := dbORM.Where("id_producto = ?", id).Delete(&Producto{})

	if resultado.Error != nil {
		http.Error(w, "Error eliminando producto con ORM: "+resultado.Error.Error(), http.StatusInternalServerError)
		return
	}

	if resultado.RowsAffected == 0 {
		http.Error(w, "No se encontró el producto", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Producto eliminado correctamente con ORM",
	})
}

func productosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if !exigirRol(w, r, "gerente", "vendedor", "bodega", "auditor") {
			return
		}
		listarProductos(w, r)

	case http.MethodPost:
		if !exigirRol(w, r, "bodega") {
			return
		}
		crearProducto(w, r)

	case http.MethodPut:
		if !exigirRol(w, r, "bodega") {
			return
		}
		actualizarProducto(w, r)

	case http.MethodDelete:
		if !exigirRol(w, r, "bodega") {
			return
		}
		eliminarProducto(w, r)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func ventasDetalladasHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !exigirRol(w, r, "gerente", "auditor") {
		return
	}

	ctx := context.Background()

	tx, err := db.Begin(ctx)
	if err != nil {
		http.Error(w, "Error iniciando transacción para ventas detalladas: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `CALL sp_ventas_detalladas('cursor_ventas_detalladas')`)
	if err != nil {
		http.Error(w, "Error ejecutando stored procedure de ventas detalladas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := tx.Query(ctx, `FETCH ALL FROM cursor_ventas_detalladas`)
	if err != nil {
		http.Error(w, "Error leyendo cursor de ventas detalladas: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ventas []VentaDetallada

	for rows.Next() {
		var v VentaDetallada

		err := rows.Scan(
			&v.IDVenta,
			&v.Fecha,
			&v.Cliente,
			&v.Empleado,
			&v.Producto,
			&v.Cantidad,
			&v.PrecioUnitario,
			&v.Subtotal,
		)

		if err != nil {
			http.Error(w, "Error leyendo ventas detalladas: "+err.Error(), http.StatusInternalServerError)
			return
		}

		ventas = append(ventas, v)
	}

	if rows.Err() != nil {
		http.Error(w, "Error recorriendo ventas detalladas: "+rows.Err().Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		http.Error(w, "Error confirmando transacción de ventas detalladas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ventas)
}

func reporteProductosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !exigirRol(w, r, "gerente", "auditor") {
		return
	}

	ctx := context.Background()

	tx, err := db.Begin(ctx)
	if err != nil {
		http.Error(w, "Error iniciando transacción para reporte: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `CALL sp_reporte_productos('cursor_reporte_productos')`)
	if err != nil {
		http.Error(w, "Error ejecutando stored procedure de reporte: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := tx.Query(ctx, `FETCH ALL FROM cursor_reporte_productos`)
	if err != nil {
		http.Error(w, "Error leyendo cursor de reporte: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reporte []ReporteProducto

	for rows.Next() {
		var rp ReporteProducto

		err := rows.Scan(
			&rp.Producto,
			&rp.TotalUnidadesVendidas,
			&rp.TotalIngresos,
		)

		if err != nil {
			http.Error(w, "Error leyendo reporte de productos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		reporte = append(reporte, rp)
	}

	if rows.Err() != nil {
		http.Error(w, "Error recorriendo reporte de productos: "+rows.Err().Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		http.Error(w, "Error confirmando transacción de reporte: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reporte)
}

func clientesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if !exigirRol(w, r, "gerente", "vendedor", "auditor") {
			return
		}
		listarClientes(w, r)

	case http.MethodPost:
		if !exigirRol(w, r, "vendedor") {
			return
		}
		crearCliente(w, r)

	case http.MethodPut:
		if !exigirRol(w, r, "vendedor") {
			return
		}
		actualizarCliente(w, r)

	case http.MethodDelete:
		if !exigirRol(w, r, "vendedor") {
			return
		}
		eliminarCliente(w, r)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func ventasHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if !exigirRol(w, r, "vendedor") {
			return
		}
		crearVenta(w, r)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func crearVenta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req VentaRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.IDCliente == 0 || req.IDEmpleado == 0 || req.IDProducto == 0 || req.Cantidad <= 0 {
		http.Error(w, "Debes enviar id_cliente, id_empleado, id_producto y cantidad válida", http.StatusBadRequest)
		return
	}

	var idVenta int
	var precioUnitario float64
	var subtotal float64

	err = db.QueryRow(context.Background(), `
		CALL sp_registrar_venta($1, $2, $3, $4, NULL, NULL, NULL)
	`, req.IDCliente, req.IDEmpleado, req.IDProducto, req.Cantidad).Scan(
		&idVenta,
		&precioUnitario,
		&subtotal,
	)

	if err != nil {
		http.Error(w, "Error registrando venta con stored procedure: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := VentaResponse{
		Mensaje:        "Venta registrada correctamente con stored procedure",
		IDVenta:        idVenta,
		IDProducto:     req.IDProducto,
		Cantidad:       req.Cantidad,
		PrecioUnitario: precioUnitario,
		Subtotal:       subtotal,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func listarClientes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT id_cliente, nombre, telefono, correo
		FROM cliente
		ORDER BY id_cliente;
	`

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Error consultando clientes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clientes []Cliente

	for rows.Next() {
		var c Cliente
		err := rows.Scan(&c.IDCliente, &c.Nombre, &c.Telefono, &c.Correo)
		if err != nil {
			http.Error(w, "Error leyendo clientes: "+err.Error(), http.StatusInternalServerError)
			return
		}
		clientes = append(clientes, c)
	}

	json.NewEncoder(w).Encode(clientes)
}

func crearCliente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var c Cliente
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = db.QueryRow(context.Background(), `
		CALL sp_crear_cliente($1, $2, $3, NULL)
	`, c.Nombre, c.Telefono, c.Correo).Scan(&c.IDCliente)

	if err != nil {
		http.Error(w, "Error creando cliente con stored procedure: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func eliminarCliente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idTexto := r.URL.Query().Get("id")
	if idTexto == "" {
		http.Error(w, "Debes enviar el id del cliente", http.StatusBadRequest)
		return
	}

	idCliente, err := strconv.Atoi(idTexto)
	if err != nil {
		http.Error(w, "El id del cliente debe ser numérico", http.StatusBadRequest)
		return
	}

	var filasAfectadas int

	err = db.QueryRow(context.Background(), `
		CALL sp_eliminar_cliente($1, NULL)
	`, idCliente).Scan(&filasAfectadas)

	if err != nil {
		http.Error(w, "Error eliminando cliente con stored procedure: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje":          "Cliente eliminado correctamente con stored procedure",
		"filas_afectadas":  filasAfectadas,
	})
}

func actualizarCliente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var c Cliente
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if c.IDCliente == 0 {
		http.Error(w, "Debes enviar id_cliente", http.StatusBadRequest)
		return
	}

	var filasAfectadas int

	err = db.QueryRow(context.Background(), `
		CALL sp_actualizar_cliente($1, $2, $3, $4, NULL)
	`, c.IDCliente, c.Nombre, c.Telefono, c.Correo).Scan(&filasAfectadas)

	if err != nil {
		http.Error(w, "Error actualizando cliente con stored procedure: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje":          "Cliente actualizado correctamente con stored procedure",
		"filas_afectadas":  filasAfectadas,
	})
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}