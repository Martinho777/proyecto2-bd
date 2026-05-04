package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

type Producto struct {
	IDProducto  int     `json:"id_producto"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
	IDCategoria int     `json:"id_categoria"`
	IDProveedor int     `json:"id_proveedor"`
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

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://proy2:secret@127.0.0.1:5433/tienda_db?sslmode=disable"
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

	http.HandleFunc("/", enableCORS(inicioHandler))
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

func listarProductos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT id_producto, nombre, descripcion, precio, stock, id_categoria, id_proveedor
		FROM producto
		ORDER BY id_producto;
	`

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Error consultando productos", http.StatusInternalServerError)
		log.Println("Error en query de productos:", err)
		return
	}
	defer rows.Close()

	var productos []Producto

	for rows.Next() {
		var p Producto
		err := rows.Scan(
			&p.IDProducto,
			&p.Nombre,
			&p.Descripcion,
			&p.Precio,
			&p.Stock,
			&p.IDCategoria,
			&p.IDProveedor,
		)
		if err != nil {
			http.Error(w, "Error leyendo productos", http.StatusInternalServerError)
			log.Println("Error escaneando producto:", err)
			return
		}
		productos = append(productos, p)
	}

	if rows.Err() != nil {
		http.Error(w, "Error recorriendo productos", http.StatusInternalServerError)
		log.Println("Error en rows:", rows.Err())
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

	query := `
		INSERT INTO producto (nombre, descripcion, precio, stock, id_categoria, id_proveedor)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id_producto;
	`

	err = db.QueryRow(
		context.Background(),
		query,
		p.Nombre,
		p.Descripcion,
		p.Precio,
		p.Stock,
		p.IDCategoria,
		p.IDProveedor,
	).Scan(&p.IDProducto)

	if err != nil {
		http.Error(w, "Error creando producto: "+err.Error(), http.StatusInternalServerError)
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

	query := `
		UPDATE producto
		SET nombre = $1, descripcion = $2, precio = $3, stock = $4, id_categoria = $5, id_proveedor = $6
		WHERE id_producto = $7
	`

	commandTag, err := db.Exec(
		context.Background(),
		query,
		p.Nombre,
		p.Descripcion,
		p.Precio,
		p.Stock,
		p.IDCategoria,
		p.IDProveedor,
		p.IDProducto,
	)

	if err != nil {
		http.Error(w, "Error actualizando producto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "No se encontró el producto", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Producto actualizado correctamente",
	})
}

func eliminarProducto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Debes enviar el id del producto", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM producto WHERE id_producto = $1`

	commandTag, err := db.Exec(context.Background(), query, id)
	if err != nil {
		http.Error(w, "Error eliminando producto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "No se encontró el producto", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Producto eliminado correctamente",
	})
}

func productosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listarProductos(w, r)
	case http.MethodPost:
		crearProducto(w, r)
	case http.MethodPut:
		actualizarProducto(w, r)
	case http.MethodDelete:
		eliminarProducto(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func ventasDetalladasHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT id_venta, fecha, cliente, empleado, producto, cantidad, precio_unitario, subtotal
		FROM vista_ventas_detalladas
		ORDER BY id_venta, producto;
	`

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Error consultando ventas detalladas", http.StatusInternalServerError)
		log.Println("Error en query de ventas detalladas:", err)
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
			http.Error(w, "Error leyendo ventas detalladas", http.StatusInternalServerError)
			log.Println("Error escaneando venta detallada:", err)
			return
		}
		ventas = append(ventas, v)
	}

	if rows.Err() != nil {
		http.Error(w, "Error recorriendo ventas detalladas", http.StatusInternalServerError)
		log.Println("Error en rows:", rows.Err())
		return
	}

	json.NewEncoder(w).Encode(ventas)
}

func reporteProductosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT 
			p.nombre AS producto,
			SUM(dv.cantidad) AS total_unidades_vendidas,
			SUM(dv.cantidad * dv.precio_unitario) AS total_ingresos
		FROM detalle_venta dv
		JOIN producto p ON dv.id_producto = p.id_producto
		GROUP BY p.nombre
		HAVING SUM(dv.cantidad) > 1
		ORDER BY total_unidades_vendidas DESC, producto;
	`

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Error consultando reporte de productos: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error en query de reporte de productos:", err)
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
			log.Println("Error escaneando reporte de productos:", err)
			return
		}
		reporte = append(reporte, rp)
	}

	if rows.Err() != nil {
		http.Error(w, "Error recorriendo reporte de productos: "+rows.Err().Error(), http.StatusInternalServerError)
		log.Println("Error en rows:", rows.Err())
		return
	}

	json.NewEncoder(w).Encode(reporte)
}

func clientesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listarClientes(w, r)
	case http.MethodPost:
		crearCliente(w, r)
	case http.MethodPut:
		actualizarCliente(w, r)
	case http.MethodDelete:
		eliminarCliente(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func ventasHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
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

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		http.Error(w, "Error iniciando transacción: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	var precio float64
	var stock int

	err = tx.QueryRow(ctx, `
		SELECT precio, stock
		FROM producto
		WHERE id_producto = $1
	`, req.IDProducto).Scan(&precio, &stock)
	if err != nil {
		http.Error(w, "Error obteniendo producto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if stock < req.Cantidad {
		http.Error(w, "Stock insuficiente", http.StatusBadRequest)
		return
	}

	var idVenta int
	err = tx.QueryRow(ctx, `
		INSERT INTO venta (fecha, id_cliente, id_empleado)
		VALUES (NOW(), $1, $2)
		RETURNING id_venta
	`, req.IDCliente, req.IDEmpleado).Scan(&idVenta)
	if err != nil {
		http.Error(w, "Error creando venta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO detalle_venta (id_venta, id_producto, cantidad, precio_unitario)
		VALUES ($1, $2, $3, $4)
	`, idVenta, req.IDProducto, req.Cantidad, precio)
	if err != nil {
		http.Error(w, "Error creando detalle de venta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(ctx, `
		UPDATE producto
		SET stock = stock - $1
		WHERE id_producto = $2
	`, req.Cantidad, req.IDProducto)
	if err != nil {
		http.Error(w, "Error actualizando stock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		http.Error(w, "Error confirmando transacción: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := VentaResponse{
		Mensaje:        "Venta registrada correctamente",
		IDVenta:        idVenta,
		IDProducto:     req.IDProducto,
		Cantidad:       req.Cantidad,
		PrecioUnitario: precio,
		Subtotal:       precio * float64(req.Cantidad),
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

	query := `
		INSERT INTO cliente (nombre, telefono, correo)
		VALUES ($1, $2, $3)
		RETURNING id_cliente;
	`

	err = db.QueryRow(context.Background(), query, c.Nombre, c.Telefono, c.Correo).Scan(&c.IDCliente)
	if err != nil {
		http.Error(w, "Error creando cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func eliminarCliente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Debes enviar el id del cliente", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM cliente WHERE id_cliente = $1`

	commandTag, err := db.Exec(context.Background(), query, id)
	if err != nil {
		http.Error(w, "Error eliminando cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "No se encontró el cliente", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Cliente eliminado correctamente",
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

	query := `
		UPDATE cliente
		SET nombre = $1, telefono = $2, correo = $3
		WHERE id_cliente = $4
	`

	commandTag, err := db.Exec(context.Background(), query, c.Nombre, c.Telefono, c.Correo, c.IDCliente)
	if err != nil {
		http.Error(w, "Error actualizando cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "No se encontró el cliente", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Cliente actualizado correctamente",
	})
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}