-- ==========================================
-- STORED PROCEDURES - PROYECTO 3
-- ==========================================

-- 1. Crear cliente con parámetros de entrada/salida y manejo de excepciones
CREATE OR REPLACE PROCEDURE sp_crear_cliente(
    IN p_nombre VARCHAR,
    IN p_telefono VARCHAR,
    IN p_correo VARCHAR,
    OUT o_id_cliente INT
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO cliente (nombre, telefono, correo)
    VALUES (p_nombre, p_telefono, p_correo)
    RETURNING id_cliente INTO o_id_cliente;

EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Error al crear cliente: %', SQLERRM;
END;
$$;


-- 2. Actualizar cliente
CREATE OR REPLACE PROCEDURE sp_actualizar_cliente(
    IN p_id_cliente INT,
    IN p_nombre VARCHAR,
    IN p_telefono VARCHAR,
    IN p_correo VARCHAR,
    OUT o_filas_afectadas INT
)
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE cliente
    SET nombre = p_nombre,
        telefono = p_telefono,
        correo = p_correo
    WHERE id_cliente = p_id_cliente;

    GET DIAGNOSTICS o_filas_afectadas = ROW_COUNT;

    IF o_filas_afectadas = 0 THEN
        RAISE EXCEPTION 'No existe el cliente con id %', p_id_cliente;
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Error al actualizar cliente: %', SQLERRM;
END;
$$;


-- 3. Eliminar cliente
CREATE OR REPLACE PROCEDURE sp_eliminar_cliente(
    IN p_id_cliente INT,
    OUT o_filas_afectadas INT
)
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM cliente
    WHERE id_cliente = p_id_cliente;

    GET DIAGNOSTICS o_filas_afectadas = ROW_COUNT;

    IF o_filas_afectadas = 0 THEN
        RAISE EXCEPTION 'No existe el cliente con id %', p_id_cliente;
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Error al eliminar cliente: %', SQLERRM;
END;
$$;


-- 4. Registrar venta con transacción explícita y rollback
CREATE OR REPLACE PROCEDURE sp_registrar_venta(
    IN p_id_cliente INT,
    IN p_id_empleado INT,
    IN p_id_producto INT,
    IN p_cantidad INT,
    OUT o_id_venta INT,
    OUT o_precio_unitario NUMERIC,
    OUT o_subtotal NUMERIC
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_stock INT;
BEGIN
    IF p_cantidad <= 0 THEN
        ROLLBACK;
        RAISE EXCEPTION 'La cantidad debe ser mayor que cero';
    END IF;

    SELECT precio, stock
    INTO o_precio_unitario, v_stock
    FROM producto
    WHERE id_producto = p_id_producto;

    IF NOT FOUND THEN
        ROLLBACK;
        RAISE EXCEPTION 'No existe el producto con id %', p_id_producto;
    END IF;

    IF v_stock < p_cantidad THEN
        ROLLBACK;
        RAISE EXCEPTION 'Stock insuficiente. Stock disponible: %, cantidad solicitada: %', v_stock, p_cantidad;
    END IF;

    INSERT INTO venta (fecha, id_cliente, id_empleado)
    VALUES (NOW(), p_id_cliente, p_id_empleado)
    RETURNING id_venta INTO o_id_venta;

    INSERT INTO detalle_venta (id_venta, id_producto, cantidad, precio_unitario)
    VALUES (o_id_venta, p_id_producto, p_cantidad, o_precio_unitario);

    UPDATE producto
    SET stock = stock - p_cantidad
    WHERE id_producto = p_id_producto;

    o_subtotal := o_precio_unitario * p_cantidad;

    COMMIT;
END;
$$;


-- 5. Reporte de productos vendidos usando cursor
CREATE OR REPLACE PROCEDURE sp_reporte_productos(
    IN p_cursor REFCURSOR
)
LANGUAGE plpgsql
AS $$
BEGIN
    OPEN p_cursor FOR
        SELECT 
            p.nombre AS producto,
            COALESCE(SUM(dv.cantidad), 0) AS total_unidades_vendidas,
            COALESCE(SUM(dv.cantidad * dv.precio_unitario), 0) AS total_ingresos
        FROM detalle_venta dv
        JOIN producto p ON dv.id_producto = p.id_producto
        GROUP BY p.nombre
        HAVING SUM(dv.cantidad) > 1
        ORDER BY total_unidades_vendidas DESC, producto;
END;
$$;


-- 6. Ventas detalladas usando cursor
CREATE OR REPLACE PROCEDURE sp_ventas_detalladas(
    IN p_cursor REFCURSOR
)
LANGUAGE plpgsql
AS $$
BEGIN
    OPEN p_cursor FOR
        SELECT 
            id_venta,
            fecha,
            cliente,
            empleado,
            producto,
            cantidad,
            precio_unitario,
            subtotal
        FROM vista_ventas_detalladas
        ORDER BY id_venta, producto;
END;
$$;