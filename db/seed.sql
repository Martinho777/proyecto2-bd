-- ==========================================
-- 25 REGISTROS POR TABLA
-- ==========================================

-- =========================
-- CATEGORIA
-- =========================
INSERT INTO categoria (id_categoria, nombre, descripcion) VALUES
(1, 'Bebidas', 'Productos líquidos para consumo'),
(2, 'Snacks', 'Productos empacados para consumo rápido'),
(3, 'Limpieza', 'Productos de higiene y limpieza'),
(4, 'Lácteos', 'Productos derivados de la leche'),
(5, 'Panadería', 'Productos de pan y repostería'),
(6, 'Enlatados', 'Productos conservados en lata'),
(7, 'Cereales', 'Cereales y granolas'),
(8, 'Congelados', 'Productos almacenados en congelación'),
(9, 'Carnes frías', 'Jamones, salchichas y embutidos'),
(10, 'Frutas', 'Frutas frescas'),
(11, 'Verduras', 'Verduras frescas'),
(12, 'Cuidado personal', 'Productos de uso e higiene personal'),
(13, 'Abarrotes', 'Productos básicos de despensa'),
(14, 'Galletas', 'Galletas dulces y saladas'),
(15, 'Dulces', 'Caramelos, chocolates y dulces'),
(16, 'Pastas', 'Pastas y fideos'),
(17, 'Salsas', 'Salsas y aderezos'),
(18, 'Mascotas', 'Productos para mascotas'),
(19, 'Bebidas energéticas', 'Bebidas energéticas y deportivas'),
(20, 'Café y té', 'Café, té e infusiones'),
(21, 'Harinas', 'Harinas y mezclas'),
(22, 'Huevos', 'Productos avícolas'),
(23, 'Quesos', 'Quesos y derivados'),
(24, 'Helados', 'Postres y helados'),
(25, 'Papelería', 'Artículos básicos de papelería');

-- =========================
-- PROVEEDOR
-- =========================
INSERT INTO proveedor (id_proveedor, nombre, telefono, correo) VALUES
(1, 'Distribuidora Central', '5555-1001', 'central@correo.com'),
(2, 'Alimentos del Norte', '5555-1002', 'norte@correo.com'),
(3, 'Fresh Market Supply', '5555-1003', 'fresh@correo.com'),
(4, 'Lácteos Unidos', '5555-1004', 'lacteos@correo.com'),
(5, 'Pan Express', '5555-1005', 'pan@correo.com'),
(6, 'Conservas Premium', '5555-1006', 'conservas@correo.com'),
(7, 'Granos Selectos', '5555-1007', 'granos@correo.com'),
(8, 'Frío Total', '5555-1008', 'frio@correo.com'),
(9, 'Embutidos Modernos', '5555-1009', 'embutidos@correo.com'),
(10, 'Campo Fresco', '5555-1010', 'campo@correo.com'),
(11, 'Verdes del Valle', '5555-1011', 'verdes@correo.com'),
(12, 'BioCare GT', '5555-1012', 'biocare@correo.com'),
(13, 'Despensa Uno', '5555-1013', 'despensa@correo.com'),
(14, 'Galletería Real', '5555-1014', 'galletas@correo.com'),
(15, 'Dulcería Nacional', '5555-1015', 'dulces@correo.com'),
(16, 'Pastas La Casa', '5555-1016', 'pastas@correo.com'),
(17, 'Sabor y Salsa', '5555-1017', 'salsas@correo.com'),
(18, 'Pet Home', '5555-1018', 'pets@correo.com'),
(19, 'Energy Plus', '5555-1019', 'energy@correo.com'),
(20, 'Café Montaña', '5555-1020', 'cafe@correo.com'),
(21, 'Molinos del Sur', '5555-1021', 'molinos@correo.com'),
(22, 'Avícola Maya', '5555-1022', 'avicola@correo.com'),
(23, 'Quesos Selectos', '5555-1023', 'quesos@correo.com'),
(24, 'Frío Dulce', '5555-1024', 'helados@correo.com'),
(25, 'Office Market', '5555-1025', 'papeleria@correo.com');

-- =========================
-- CLIENTE
-- =========================
INSERT INTO cliente (id_cliente, nombre, telefono, correo) VALUES
(1, 'Juan Pérez', '4444-1001', 'juan@email.com'),
(2, 'María López', '4444-1002', 'maria@email.com'),
(3, 'Carlos García', '4444-1003', 'carlos@email.com'),
(4, 'Ana Morales', '4444-1004', 'ana@email.com'),
(5, 'Luis Fernández', '4444-1005', 'luis@email.com'),
(6, 'Andrea Ruiz', '4444-1006', 'andrea@email.com'),
(7, 'José Martínez', '4444-1007', 'jose@email.com'),
(8, 'Patricia Gómez', '4444-1008', 'patricia@email.com'),
(9, 'Miguel Castro', '4444-1009', 'miguel@email.com'),
(10, 'Daniela Herrera', '4444-1010', 'daniela@email.com'),
(11, 'Oscar Mejía', '4444-1011', 'oscar@email.com'),
(12, 'Lucía Ramírez', '4444-1012', 'lucia@email.com'),
(13, 'Roberto Fuentes', '4444-1013', 'roberto@email.com'),
(14, 'Karla Díaz', '4444-1014', 'karla@email.com'),
(15, 'Fernando Reyes', '4444-1015', 'fernando@email.com'),
(16, 'Paola Torres', '4444-1016', 'paola@email.com'),
(17, 'Diego Mendoza', '4444-1017', 'diego@email.com'),
(18, 'Valeria Soto', '4444-1018', 'valeria@email.com'),
(19, 'Héctor Cifuentes', '4444-1019', 'hector@email.com'),
(20, 'Camila Navarro', '4444-1020', 'camila@email.com'),
(21, 'Ricardo Pineda', '4444-1021', 'ricardo@email.com'),
(22, 'Monserrat León', '4444-1022', 'monserrat@email.com'),
(23, 'Kevin Alvarado', '4444-1023', 'kevin@email.com'),
(24, 'Natalia Barrios', '4444-1024', 'natalia@email.com'),
(25, 'Esteban Orellana', '4444-1025', 'esteban@email.com');

-- =========================
-- EMPLEADO
-- =========================
INSERT INTO empleado (id_empleado, nombre, puesto, correo) VALUES
(1, 'Sofía Hernández', 'Cajera', 'sofia@tienda.com'),
(2, 'Pedro Ramírez', 'Vendedor', 'pedro@tienda.com'),
(3, 'Elena Castro', 'Supervisora', 'elena@tienda.com'),
(4, 'Diego Ruiz', 'Cajero', 'diego@tienda.com'),
(5, 'Laura Méndez', 'Vendedora', 'laura@tienda.com'),
(6, 'Mario Aguilar', 'Cajero', 'mario@tienda.com'),
(7, 'Claudia Pérez', 'Vendedora', 'claudia@tienda.com'),
(8, 'Jorge Molina', 'Bodeguero', 'jorge@tienda.com'),
(9, 'Andrea López', 'Supervisora', 'andreal@tienda.com'),
(10, 'Samuel Herrera', 'Cajero', 'samuel@tienda.com'),
(11, 'Paola García', 'Vendedora', 'paolag@tienda.com'),
(12, 'Cristian Morales', 'Bodeguero', 'cristian@tienda.com'),
(13, 'Verónica Díaz', 'Cajera', 'veronica@tienda.com'),
(14, 'Raúl Gómez', 'Vendedor', 'raul@tienda.com'),
(15, 'Mónica León', 'Cajera', 'monica@tienda.com'),
(16, 'Andrés Ponce', 'Supervisora', 'andres@tienda.com'),
(17, 'Gabriela Soto', 'Vendedora', 'gabriela@tienda.com'),
(18, 'Ricardo Barrios', 'Cajero', 'ricardob@tienda.com'),
(19, 'Tatiana Fuentes', 'Bodeguera', 'tatiana@tienda.com'),
(20, 'Nicolás Reyes', 'Vendedor', 'nicolas@tienda.com'),
(21, 'Julieta Paz', 'Cajera', 'julieta@tienda.com'),
(22, 'Mauricio Castañeda', 'Bodeguero', 'mauricio@tienda.com'),
(23, 'Silvia Alvarado', 'Supervisora', 'silvia@tienda.com'),
(24, 'Daniel Rivas', 'Cajero', 'danielr@tienda.com'),
(25, 'Mariana Ochoa', 'Vendedora', 'mariana@tienda.com');

-- =========================
-- USUARIOS DE APLICACIÓN
-- =========================
INSERT INTO app_usuario (nombre, correo, password, rol) VALUES
('Administrador General', 'admin@tienda.com', 'secret', 'admin'),
('Gerente de Tienda', 'gerente@tienda.com', 'secret', 'gerente'),
('Vendedor Principal', 'vendedor@tienda.com', 'secret', 'vendedor'),
('Encargado de Bodega', 'bodega@tienda.com', 'secret', 'bodega'),
('Auditor Interno', 'auditor@tienda.com', 'secret', 'auditor');

-- =========================
-- PRODUCTO
-- =========================
INSERT INTO producto (id_producto, nombre, descripcion, precio, stock, id_categoria, id_proveedor) VALUES
(1, 'Coca Cola 355ml', 'Bebida gaseosa', 8.50, 100, 1, 1),
(2, 'Pepsi 355ml', 'Bebida gaseosa', 8.00, 90, 1, 1),
(3, 'Tortrix', 'Snack de maíz', 4.00, 150, 2, 2),
(4, 'Churritos', 'Snack salado', 4.50, 130, 2, 2),
(5, 'Cloro 1L', 'Producto de limpieza', 12.00, 50, 3, 3),
(6, 'Detergente 500g', 'Detergente en polvo', 18.00, 60, 3, 3),
(7, 'Leche Entera 1L', 'Leche líquida', 10.50, 80, 4, 4),
(8, 'Yogurt Natural', 'Yogurt de vaso', 6.75, 70, 4, 4),
(9, 'Pan francés', 'Unidad de pan francés', 1.00, 200, 5, 5),
(10, 'Pan dulce', 'Pieza de pan dulce', 3.00, 120, 5, 5),
(11, 'Atún en lata', 'Atún en aceite', 11.25, 95, 6, 6),
(12, 'Corn Flakes', 'Cereal de maíz', 24.00, 55, 7, 7),
(13, 'Pizza congelada', 'Pizza lista para horno', 38.50, 35, 8, 8),
(14, 'Jamón de pavo', 'Paquete de jamón rebanado', 21.00, 45, 9, 9),
(15, 'Manzana roja', 'Fruta fresca por unidad', 2.50, 180, 10, 10),
(16, 'Tomate', 'Verdura fresca por libra', 5.50, 140, 11, 11),
(17, 'Shampoo 400ml', 'Cuidado capilar', 27.00, 48, 12, 12),
(18, 'Arroz 1kg', 'Arroz blanco', 9.75, 110, 13, 13),
(19, 'Galletas de vainilla', 'Galletas dulces', 7.25, 85, 14, 14),
(20, 'Chocolate barra', 'Chocolate con leche', 6.50, 75, 15, 15),
(21, 'Espagueti 400g', 'Pasta larga', 8.20, 92, 16, 16),
(22, 'Ketchup 500g', 'Salsa de tomate', 13.50, 58, 17, 17),
(23, 'Concentrado perro', 'Alimento para perro 2kg', 42.00, 33, 18, 18),
(24, 'Bebida energética', 'Bebida energética 473ml', 14.00, 67, 19, 19),
(25, 'Café molido 250g', 'Café tostado y molido', 29.00, 52, 20, 20);

-- =========================
-- VENTA
-- =========================
INSERT INTO venta (id_venta, fecha, id_cliente, id_empleado) VALUES
(1, '2026-04-01 10:15:00', 1, 1),
(2, '2026-04-01 11:30:00', 2, 2),
(3, '2026-04-02 09:45:00', 3, 1),
(4, '2026-04-02 15:20:00', 4, 4),
(5, '2026-04-03 17:10:00', 5, 2),
(6, '2026-04-04 08:35:00', 6, 3),
(7, '2026-04-04 12:50:00', 7, 5),
(8, '2026-04-05 14:05:00', 8, 6),
(9, '2026-04-05 16:25:00', 9, 7),
(10, '2026-04-06 09:10:00', 10, 8),
(11, '2026-04-06 11:40:00', 11, 9),
(12, '2026-04-06 18:15:00', 12, 10),
(13, '2026-04-07 10:05:00', 13, 11),
(14, '2026-04-07 13:45:00', 14, 12),
(15, '2026-04-08 08:20:00', 15, 13),
(16, '2026-04-08 12:30:00', 16, 14),
(17, '2026-04-08 17:55:00', 17, 15),
(18, '2026-04-09 09:50:00', 18, 16),
(19, '2026-04-09 14:10:00', 19, 17),
(20, '2026-04-10 10:35:00', 20, 18),
(21, '2026-04-10 15:40:00', 21, 19),
(22, '2026-04-11 08:45:00', 22, 20),
(23, '2026-04-11 12:25:00', 23, 21),
(24, '2026-04-12 16:05:00', 24, 22),
(25, '2026-04-12 18:30:00', 25, 23);

-- =========================
-- DETALLE_VENTA
-- 50 registros
-- =========================
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
(5, 10, 4, 3.00),

(6, 11, 2, 11.25),
(6, 18, 1, 9.75),

(7, 12, 1, 24.00),
(7, 19, 3, 7.25),

(8, 13, 1, 38.50),
(8, 24, 2, 14.00),

(9, 14, 2, 21.00),
(9, 16, 3, 5.50),

(10, 15, 8, 2.50),
(10, 20, 2, 6.50),

(11, 17, 1, 27.00),
(11, 22, 1, 13.50),

(12, 18, 3, 9.75),
(12, 21, 2, 8.20),

(13, 23, 1, 42.00),
(13, 24, 1, 14.00),

(14, 25, 1, 29.00),
(14, 9, 10, 1.00),

(15, 1, 1, 8.50),
(15, 11, 2, 11.25),

(16, 2, 2, 8.00),
(16, 7, 1, 10.50),

(17, 3, 4, 4.00),
(17, 19, 2, 7.25),

(18, 4, 2, 4.50),
(18, 10, 3, 3.00),

(19, 5, 1, 12.00),
(19, 17, 1, 27.00),

(20, 6, 1, 18.00),
(20, 22, 2, 13.50),

(21, 8, 3, 6.75),
(21, 15, 5, 2.50),

(22, 12, 1, 24.00),
(22, 20, 3, 6.50),

(23, 13, 1, 38.50),
(23, 14, 1, 21.00),

(24, 16, 4, 5.50),
(24, 18, 2, 9.75),

(25, 21, 2, 8.20),
(25, 25, 1, 29.00);

-- =========================
-- AJUSTE DE SECUENCIAS
-- =========================
SELECT setval('categoria_id_categoria_seq', (SELECT MAX(id_categoria) FROM categoria));
SELECT setval('proveedor_id_proveedor_seq', (SELECT MAX(id_proveedor) FROM proveedor));
SELECT setval('cliente_id_cliente_seq', (SELECT MAX(id_cliente) FROM cliente));
SELECT setval('empleado_id_empleado_seq', (SELECT MAX(id_empleado) FROM empleado));
SELECT setval('producto_id_producto_seq', (SELECT MAX(id_producto) FROM producto));
SELECT setval('venta_id_venta_seq', (SELECT MAX(id_venta) FROM venta));