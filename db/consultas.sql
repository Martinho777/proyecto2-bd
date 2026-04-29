-- 1. CONSULTAS CON JOIN

-- 1.1 Ventas con cliente y empleado
SELECT 
    v.id_venta,
    v.fecha,
    c.nombre AS cliente,
    e.nombre AS empleado
FROM venta v
JOIN cliente c ON v.id_cliente = c.id_cliente
JOIN empleado e ON v.id_empleado = e.id_empleado;


-- 1.2 Detalle de ventas con producto
SELECT 
    dv.id_venta,
    p.nombre AS producto,
    dv.cantidad,
    dv.precio_unitario,
    (dv.cantidad * dv.precio_unitario) AS subtotal
FROM detalle_venta dv
JOIN producto p ON dv.id_producto = p.id_producto;


-- 1.3 Productos con categoría y proveedor
SELECT 
    p.id_producto,
    p.nombre AS producto,
    c.nombre AS categoria,
    pr.nombre AS proveedor,
    p.precio,
    p.stock
FROM producto p
JOIN categoria c ON p.id_categoria = c.id_categoria
JOIN proveedor pr ON p.id_proveedor = pr.id_proveedor;


-- 2. CONSULTAS CON SUBQUERY

-- 2.1 Productos con stock mayor al promedio
SELECT 
    nombre,
    stock
FROM producto
WHERE stock > (
    SELECT AVG(stock)
    FROM producto
);


-- 2.2 Clientes que han realizado compras
SELECT 
    nombre
FROM cliente
WHERE id_cliente IN (
    SELECT id_cliente
    FROM venta
);


-- 3. CONSULTAS CON GROUP BY, HAVING Y AGREGACION

-- 3.1 Total vendido por producto
SELECT 
    p.nombre AS producto,
    SUM(dv.cantidad) AS total_unidades_vendidas,
    SUM(dv.cantidad * dv.precio_unitario) AS total_ingresos
FROM detalle_venta dv
JOIN producto p ON dv.id_producto = p.id_producto
GROUP BY p.nombre
HAVING SUM(dv.cantidad) > 1;


-- 4. CONSULTA CON CTE

-- 4.1 Total por venta
WITH totales_por_venta AS (
    SELECT 
        dv.id_venta,
        SUM(dv.cantidad * dv.precio_unitario) AS total
    FROM detalle_venta dv
    GROUP BY dv.id_venta
)
SELECT 
    t.id_venta,
    t.total
FROM totales_por_venta t
ORDER BY t.id_venta;


-- 5. CONSULTAS SOBRE VIEW

-- 5.1 Consulta de la vista de ventas detalladas
SELECT * 
FROM vista_ventas_detalladas;



-- 6. TRANSACCIONES

-- 6.1 Transacción correcta con COMMIT
BEGIN;

INSERT INTO venta (fecha, id_cliente, id_empleado)
VALUES ('2026-04-10 12:00:00', 1, 1);

INSERT INTO detalle_venta (id_venta, id_producto, cantidad, precio_unitario)
VALUES 
(6, 1, 2, 8.50),
(6, 3, 1, 4.00);

UPDATE producto
SET stock = stock - 2
WHERE id_producto = 1;

UPDATE producto
SET stock = stock - 1
WHERE id_producto = 3;

COMMIT;


-- 6.2 Transacción con error y ROLLBACK
BEGIN;

INSERT INTO venta (fecha, id_cliente, id_empleado)
VALUES ('2026-04-11 14:00:00', 2, 2);

INSERT INTO detalle_venta (id_venta, id_producto, cantidad, precio_unitario)
VALUES
(7, 2, 1, 8.00),
(7, 999, 1, 5.00);

UPDATE producto
SET stock = stock - 1
WHERE id_producto = 2;

ROLLBACK;