-- ==========================================
-- ROLES DE BASE DE DATOS - PROYECTO 3
-- ==========================================

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'rol_admin') THEN
        CREATE ROLE rol_admin;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'rol_gerente') THEN
        CREATE ROLE rol_gerente;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'rol_vendedor') THEN
        CREATE ROLE rol_vendedor;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'rol_bodega') THEN
        CREATE ROLE rol_bodega;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'rol_auditor') THEN
        CREATE ROLE rol_auditor;
    END IF;
END
$$;

-- Limpiar permisos públicos
REVOKE ALL ON ALL TABLES IN SCHEMA public FROM PUBLIC;
REVOKE ALL ON ALL SEQUENCES IN SCHEMA public FROM PUBLIC;

-- Permitir uso del esquema
GRANT USAGE ON SCHEMA public TO rol_admin;
GRANT USAGE ON SCHEMA public TO rol_gerente;
GRANT USAGE ON SCHEMA public TO rol_vendedor;
GRANT USAGE ON SCHEMA public TO rol_bodega;
GRANT USAGE ON SCHEMA public TO rol_auditor;

-- ==========================================
-- ROL ADMIN: acceso completo
-- ==========================================
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO rol_admin;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO rol_admin;

-- ==========================================
-- ROL GERENTE: lectura general y reportes
-- ==========================================
GRANT SELECT ON categoria, proveedor, cliente, empleado, producto, venta, detalle_venta, vista_ventas_detalladas TO rol_gerente;
GRANT SELECT ON app_usuario TO rol_gerente;

-- ==========================================
-- ROL VENDEDOR: consultar datos y registrar ventas
-- ==========================================
GRANT SELECT ON cliente, empleado, producto, categoria, proveedor TO rol_vendedor;
GRANT INSERT ON venta, detalle_venta TO rol_vendedor;
GRANT UPDATE (stock) ON producto TO rol_vendedor;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE venta_id_venta_seq TO rol_vendedor;

-- ==========================================
-- ROL BODEGA: manejo de inventario
-- ==========================================
GRANT SELECT, INSERT, UPDATE ON producto, categoria, proveedor TO rol_bodega;
GRANT DELETE ON producto TO rol_bodega;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE producto_id_producto_seq TO rol_bodega;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE categoria_id_categoria_seq TO rol_bodega;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE proveedor_id_proveedor_seq TO rol_bodega;

-- ==========================================
-- ROL AUDITOR: solo lectura
-- ==========================================
GRANT SELECT ON categoria, proveedor, cliente, empleado, producto, venta, detalle_venta, vista_ventas_detalladas TO rol_auditor;