INSERT INTO categoria (nombre, descripcion) VALUES
('Bebidas', 'Productos líquidos para consumo'),
('Snacks', 'Productos empacados para consumo rápido'),
('Limpieza', 'Productos de higiene y limpieza'),
('Lácteos', 'Productos derivados de la leche'),
('Panadería', 'Productos de pan y repostería');

INSERT INTO proveedor (nombre, telefono, correo) VALUES
('Distribuidora Central', '5555-1001', 'central@correo.com'),
('Alimentos del Norte', '5555-1002', 'norte@correo.com'),
('Fresh Market Supply', '5555-1003', 'fresh@correo.com'),
('Lácteos Unidos', '5555-1004', 'lacteos@correo.com'),
('Pan Express', '5555-1005', 'pan@correo.com');

INSERT INTO cliente (nombre, telefono, correo) VALUES
('Juan Pérez', '4444-1001', 'juan@email.com'),
('María López', '4444-1002', 'maria@email.com'),
('Carlos García', '4444-1003', 'carlos@email.com'),
('Ana Morales', '4444-1004', 'ana@email.com'),
('Luis Fernández', '4444-1005', 'luis@email.com');

INSERT INTO empleado (nombre, puesto, correo) VALUES
('Sofía Hernández', 'Cajera', 'sofia@tienda.com'),
('Pedro Ramírez', 'Vendedor', 'pedro@tienda.com'),
('Elena Castro', 'Supervisora', 'elena@tienda.com'),
('Diego Ruiz', 'Cajero', 'diego@tienda.com'),
('Laura Méndez', 'Vendedora', 'laura@tienda.com');

INSERT INTO producto (nombre, descripcion, precio, stock, id_categoria, id_proveedor) VALUES
('Coca Cola 355ml', 'Bebida gaseosa', 8.50, 100, 1, 1),
('Pepsi 355ml', 'Bebida gaseosa', 8.00, 90, 1, 1),
('Tortrix', 'Snack de maíz', 4.00, 150, 2, 2),
('Churritos', 'Snack salado', 4.50, 130, 2, 2),
('Cloro 1L', 'Producto de limpieza', 12.00, 50, 3, 3),
('Detergente 500g', 'Detergente en polvo', 18.00, 60, 3, 3),
('Leche Entera 1L', 'Leche líquida', 10.50, 80, 4, 4),
('Yogurt Natural', 'Yogurt de vaso', 6.75, 70, 4, 4),
('Pan francés', 'Unidad de pan francés', 1.00, 200, 5, 5),
('Pan dulce', 'Pieza de pan dulce', 3.00, 120, 5, 5);

INSERT INTO venta (fecha, id_cliente, id_empleado) VALUES
('2026-04-01 10:15:00', 1, 1),
('2026-04-01 11:30:00', 2, 2),
('2026-04-02 09:45:00', 3, 1),
('2026-04-02 15:20:00', 4, 4),
('2026-04-03 17:10:00', 5, 2);

INSERT INTO detalle_venta (id_venta, id_producto, cantidad, precio_unitario) VALUES
(1, 1, 2, 8.50),
(1, 3, 1, 4.00),
(2, 7, 2, 10.50),
(2, 9, 6, 1.00),
(3, 5, 1, 12.00),
(3, 6, 1, 18.00),
(4, 2, 3, 8.00),
(4, 4, 2, 4.50),
(5, 8, 2, 6.75),
(5, 10, 4, 3.00);