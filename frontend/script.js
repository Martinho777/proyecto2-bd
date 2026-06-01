const API = "http://localhost:8080";

let clienteEditandoId = null;
let productoEditandoId = null;

let token = localStorage.getItem("token") || "";
let usuario = JSON.parse(localStorage.getItem("usuario") || "null");

function tieneRol(...roles) {
  if (!usuario) return false;

  const rolActual = usuario.rol || usuario.Rol || usuario.role || "";

  if (rolActual === "admin") return true;

  return roles.includes(rolActual);
}

async function apiFetch(url, options = {}) {
  const headers = {
    ...(options.headers || {})
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const res = await fetch(url, {
    ...options,
    headers
  });

  if (res.status === 401) {
    cerrarSesionLocal();
    throw new Error("Sesión inválida o expirada. Inicia sesión de nuevo.");
  }

  if (!res.ok) {
    const texto = await res.text();
    throw new Error(texto || "Error en la solicitud");
  }

  return res;
}

function mostrarApp() {
  document.getElementById("loginSection").classList.add("oculto");
  document.getElementById("sesionSection").classList.remove("oculto");
  document.getElementById("appSection").classList.remove("oculto");

  const nombre = usuario.nombre || usuario.Nombre || "";
  const correo = usuario.correo || usuario.Correo || "";
  const rol = usuario.rol || usuario.Rol || usuario.role || "";

  usuario = {
    nombre,
    correo,
    rol
  };

  localStorage.setItem("usuario", JSON.stringify(usuario));

  document.getElementById("usuarioInfo").textContent =
    `${usuario.nombre} (${usuario.correo}) - Rol: ${usuario.rol}`;

  aplicarPermisosUI();
  cargarDatosPermitidos();
}

function mostrarLogin() {
  document.getElementById("loginSection").classList.remove("oculto");
  document.getElementById("sesionSection").classList.add("oculto");
  document.getElementById("appSection").classList.add("oculto");
}

function cerrarSesionLocal() {
  token = "";
  usuario = null;
  localStorage.removeItem("token");
  localStorage.removeItem("usuario");
  mostrarLogin();
}

function aplicarPermisosUI() {
  const puedeVerClientes = tieneRol("gerente", "vendedor", "auditor");
  const puedeEditarClientes = tieneRol("vendedor");

  const puedeVerProductos = tieneRol("gerente", "vendedor", "bodega", "auditor");
  const puedeEditarProductos = tieneRol("bodega");

  const puedeVender = tieneRol("vendedor");
  const puedeVerReportes = tieneRol("gerente", "auditor");

  document.getElementById("seccionClientes").classList.toggle("oculto", !puedeVerClientes);
  document.getElementById("formCliente").classList.toggle("oculto", !puedeEditarClientes);

  document.getElementById("seccionProductos").classList.toggle("oculto", !puedeVerProductos);
  document.getElementById("formProducto").classList.toggle("oculto", !puedeEditarProductos);

  document.getElementById("seccionVenta").classList.toggle("oculto", !puedeVender);
  document.getElementById("seccionReporte").classList.toggle("oculto", !puedeVerReportes);
  document.getElementById("seccionVentasDetalladas").classList.toggle("oculto", !puedeVerReportes);
}

async function cargarDatosPermitidos() {
  try {
    if (tieneRol("gerente", "vendedor", "auditor")) {
      await cargarClientes();
    }

    if (tieneRol("gerente", "vendedor", "bodega", "auditor")) {
      await cargarProductos();
    }

    if (tieneRol("gerente", "auditor")) {
      await cargarReporteProductos();
      await cargarVentasDetalladas();
    }
  } catch (err) {
    alert(err.message);
  }
}

function textoURL(valor) {
  return encodeURIComponent(valor ?? "");
}

function leerTextoURL(valor) {
  return decodeURIComponent(valor ?? "");
}

document.getElementById("formLogin").addEventListener("submit", async (e) => {
  e.preventDefault();

  const body = {
    correo: document.getElementById("loginCorreo").value,
    password: document.getElementById("loginPassword").value
  };

  try {
    const res = await fetch(`${API}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body)
    });

    if (!res.ok) {
      throw new Error("Correo o contraseña incorrectos");
    }

    const data = await res.json();

    token = data.token;
    usuario = {
      nombre: data.nombre,
      correo: data.correo,
      rol: data.rol
    };

    localStorage.setItem("token", token);
    localStorage.setItem("usuario", JSON.stringify(usuario));

    e.target.reset();
    mostrarApp();
  } catch (err) {
    alert(err.message);
  }
});

document.getElementById("btnLogout").addEventListener("click", async () => {
  try {
    if (token) {
      await apiFetch(`${API}/logout`, {
        method: "POST"
      });
    }
  } catch (err) {
    console.log(err.message);
  }

  cerrarSesionLocal();
});

document.getElementById("formVenta").addEventListener("submit", async (e) => {
  e.preventDefault();

  const body = {
    id_cliente: parseInt(document.getElementById("ventaCliente").value),
    id_empleado: parseInt(document.getElementById("ventaEmpleado").value),
    id_producto: parseInt(document.getElementById("ventaProducto").value),
    cantidad: parseInt(document.getElementById("ventaCantidad").value)
  };

  try {
    const res = await apiFetch(`${API}/ventas`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body)
    });

    const data = await res.json();

    alert(`Venta registrada correctamente. ID venta: ${data.id_venta}`);
    e.target.reset();

    if (tieneRol("gerente", "vendedor", "bodega", "auditor")) {
      cargarProductos();
    }

    if (tieneRol("gerente", "auditor")) {
      cargarReporteProductos();
      cargarVentasDetalladas();
    }
  } catch (err) {
    alert(err.message);
  }
});

async function cargarClientes() {
  try {
    const res = await apiFetch(`${API}/clientes`);
    const data = await res.json();

    const puedeEditar = tieneRol("vendedor");

    let html = `
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Nombre</th>
            <th>Teléfono</th>
            <th>Correo</th>
            ${puedeEditar ? "<th>Acciones</th>" : ""}
          </tr>
        </thead>
        <tbody>
    `;

    data.forEach(c => {
      html += `
        <tr>
          <td>${c.id_cliente}</td>
          <td>${c.nombre}</td>
          <td>${c.telefono}</td>
          <td>${c.correo}</td>
          ${puedeEditar ? `
            <td class="acciones">
              <button onclick="editarCliente(${c.id_cliente}, '${textoURL(c.nombre)}', '${textoURL(c.telefono)}', '${textoURL(c.correo)}')">Editar</button>
              <button onclick="eliminarCliente(${c.id_cliente})">Eliminar</button>
            </td>
          ` : ""}
        </tr>
      `;
    });

    html += `</tbody></table>`;
    document.getElementById("clientes").innerHTML = html;
  } catch (err) {
    document.getElementById("clientes").innerHTML = `<p class="error">${err.message}</p>`;
  }
}

function editarCliente(id, nombre, telefono, correo) {
  clienteEditandoId = id;

  document.getElementById("clienteNombre").value = leerTextoURL(nombre);
  document.getElementById("clienteTelefono").value = leerTextoURL(telefono);
  document.getElementById("clienteCorreo").value = leerTextoURL(correo);
  document.getElementById("btnCliente").textContent = "Actualizar cliente";
}

async function cargarProductos() {
  try {
    const res = await apiFetch(`${API}/productos`);
    const data = await res.json();

    const puedeEditar = tieneRol("bodega");

    let html = `
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Nombre</th>
            <th>Descripción</th>
            <th>Precio</th>
            <th>Stock</th>
            <th>Categoría</th>
            <th>Proveedor</th>
            ${puedeEditar ? "<th>Acciones</th>" : ""}
          </tr>
        </thead>
        <tbody>
    `;

    data.forEach(p => {
      html += `
        <tr>
          <td>${p.id_producto}</td>
          <td>${p.nombre}</td>
          <td>${p.descripcion}</td>
          <td>${p.precio}</td>
          <td>${p.stock}</td>
          <td>${p.id_categoria}</td>
          <td>${p.id_proveedor}</td>
          ${puedeEditar ? `
            <td class="acciones">
              <button onclick="editarProducto(${p.id_producto}, '${textoURL(p.nombre)}', '${textoURL(p.descripcion)}', ${p.precio}, ${p.stock}, ${p.id_categoria}, ${p.id_proveedor})">Editar</button>
              <button onclick="eliminarProducto(${p.id_producto})">Eliminar</button>
            </td>
          ` : ""}
        </tr>
      `;
    });

    html += `</tbody></table>`;
    document.getElementById("productos").innerHTML = html;
  } catch (err) {
    document.getElementById("productos").innerHTML = `<p class="error">${err.message}</p>`;
  }
}

function editarProducto(id, nombre, descripcion, precio, stock, idCategoria, idProveedor) {
  productoEditandoId = id;

  document.getElementById("productoNombre").value = leerTextoURL(nombre);
  document.getElementById("productoDescripcion").value = leerTextoURL(descripcion);
  document.getElementById("productoPrecio").value = precio;
  document.getElementById("productoStock").value = stock;
  document.getElementById("productoCategoria").value = idCategoria;
  document.getElementById("productoProveedor").value = idProveedor;
  document.getElementById("btnProducto").textContent = "Actualizar producto";
}

async function cargarReporteProductos() {
  try {
    const res = await apiFetch(`${API}/reporte/productos`);
    const data = await res.json();

    let html = `
      <table>
        <thead>
          <tr>
            <th>Producto</th>
            <th>Total unidades vendidas</th>
            <th>Total ingresos</th>
          </tr>
        </thead>
        <tbody>
    `;

    data.forEach(r => {
      html += `
        <tr>
          <td>${r.producto}</td>
          <td>${r.total_unidades_vendidas}</td>
          <td>${r.total_ingresos}</td>
        </tr>
      `;
    });

    html += `</tbody></table>`;
    document.getElementById("reporteProductos").innerHTML = html;
  } catch (err) {
    document.getElementById("reporteProductos").innerHTML = `<p class="error">${err.message}</p>`;
  }
}

async function cargarVentasDetalladas() {
  try {
    const res = await apiFetch(`${API}/ventas-detalladas`);
    const data = await res.json();

    let html = `
      <table>
        <thead>
          <tr>
            <th>ID Venta</th>
            <th>Fecha</th>
            <th>Cliente</th>
            <th>Empleado</th>
            <th>Producto</th>
            <th>Cantidad</th>
            <th>Precio Unitario</th>
            <th>Subtotal</th>
          </tr>
        </thead>
        <tbody>
    `;

    data.forEach(v => {
      html += `
        <tr>
          <td>${v.id_venta}</td>
          <td>${new Date(v.fecha).toLocaleString()}</td>
          <td>${v.cliente}</td>
          <td>${v.empleado}</td>
          <td>${v.producto}</td>
          <td>${v.cantidad}</td>
          <td>${v.precio_unitario}</td>
          <td>${v.subtotal}</td>
        </tr>
      `;
    });

    html += `</tbody></table>`;
    document.getElementById("ventasDetalladas").innerHTML = html;
  } catch (err) {
    document.getElementById("ventasDetalladas").innerHTML = `<p class="error">${err.message}</p>`;
  }
}

document.getElementById("formCliente").addEventListener("submit", async (e) => {
  e.preventDefault();

  const body = {
    nombre: document.getElementById("clienteNombre").value,
    telefono: document.getElementById("clienteTelefono").value,
    correo: document.getElementById("clienteCorreo").value
  };

  let metodo = "POST";
  let mensajeOk = "Cliente agregado correctamente";

  if (clienteEditandoId !== null) {
    body.id_cliente = clienteEditandoId;
    metodo = "PUT";
    mensajeOk = "Cliente actualizado correctamente";
  }

  try {
    await apiFetch(`${API}/clientes`, {
      method: metodo,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body)
    });

    alert(mensajeOk);
    e.target.reset();
    clienteEditandoId = null;
    document.getElementById("btnCliente").textContent = "Agregar cliente";
    cargarClientes();
  } catch (err) {
    alert(err.message);
  }
});

document.getElementById("formProducto").addEventListener("submit", async (e) => {
  e.preventDefault();

  const body = {
    nombre: document.getElementById("productoNombre").value,
    descripcion: document.getElementById("productoDescripcion").value,
    precio: parseFloat(document.getElementById("productoPrecio").value),
    stock: parseInt(document.getElementById("productoStock").value),
    id_categoria: parseInt(document.getElementById("productoCategoria").value),
    id_proveedor: parseInt(document.getElementById("productoProveedor").value)
  };

  let metodo = "POST";
  let mensajeOk = "Producto agregado correctamente";

  if (productoEditandoId !== null) {
    body.id_producto = productoEditandoId;
    metodo = "PUT";
    mensajeOk = "Producto actualizado correctamente";
  }

  try {
    await apiFetch(`${API}/productos`, {
      method: metodo,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body)
    });

    alert(mensajeOk);
    e.target.reset();
    productoEditandoId = null;
    document.getElementById("btnProducto").textContent = "Agregar producto";
    cargarProductos();
  } catch (err) {
    alert(err.message);
  }
});

async function eliminarCliente(id) {
  if (!confirm("¿Deseas eliminar este cliente?")) return;

  try {
    await apiFetch(`${API}/clientes?id=${id}`, {
      method: "DELETE"
    });

    alert("Cliente eliminado correctamente");
    cargarClientes();
  } catch (err) {
    alert(err.message);
  }
}

async function eliminarProducto(id) {
  if (!confirm("¿Deseas eliminar este producto?")) return;

  try {
    await apiFetch(`${API}/productos?id=${id}`, {
      method: "DELETE"
    });

    alert("Producto eliminado correctamente");
    cargarProductos();
  } catch (err) {
    alert(err.message);
  }
}

async function verificarSesionGuardada() {
  if (!token || !usuario) {
    mostrarLogin();
    return;
  }

  try {
    const res = await apiFetch(`${API}/me`);
    const data = await res.json();

    usuario = {
      nombre: data.nombre,
      correo: data.correo,
      rol: data.rol
    };

    localStorage.setItem("usuario", JSON.stringify(usuario));
    mostrarApp();
  } catch (err) {
    cerrarSesionLocal();
  }
}

verificarSesionGuardada();