const API = "http://localhost:8080";
let clienteEditandoId = null;
let productoEditandoId = null;

document.getElementById("formVenta").addEventListener("submit", async (e) => {
  e.preventDefault();

  const body = {
    id_cliente: parseInt(document.getElementById("ventaCliente").value),
    id_empleado: parseInt(document.getElementById("ventaEmpleado").value),
    id_producto: parseInt(document.getElementById("ventaProducto").value),
    cantidad: parseInt(document.getElementById("ventaCantidad").value)
  };

  const res = await fetch(`${API}/ventas`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body)
  });

  const data = await res.json().catch(() => null);

  if (res.ok) {
    alert(`Venta registrada correctamente. ID venta: ${data.id_venta}`);
    e.target.reset();
    cargarProductos();
    cargarReporteProductos();
    cargarVentasDetalladas();
  } else {
    alert(data?.mensaje || data || "Error al registrar venta");
  }
});

async function cargarClientes() {
  const res = await fetch(`${API}/clientes`);
  const data = await res.json();

  let html = `
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Nombre</th>
          <th>Teléfono</th>
          <th>Correo</th>
          <th>Acciones</th>
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
        <td class="acciones">
          <button onclick="editarCliente(${c.id_cliente}, '${c.nombre}', '${c.telefono}', '${c.correo}')">Editar</button>
          <button onclick="eliminarCliente(${c.id_cliente})">Eliminar</button>
        </td>
      </tr>
    `;
  });

  html += `</tbody></table>`;
  document.getElementById("clientes").innerHTML = html;
}

function editarCliente(id, nombre, telefono, correo) {
  clienteEditandoId = id;

  document.getElementById("clienteNombre").value = nombre;
  document.getElementById("clienteTelefono").value = telefono;
  document.getElementById("clienteCorreo").value = correo;
  document.getElementById("btnCliente").textContent = "Actualizar cliente";
}

async function cargarProductos() {
  const res = await fetch(`${API}/productos`);
  const data = await res.json();

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
          <th>Acciones</th>
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
        <td class="acciones">
          <button onclick="editarProducto(${p.id_producto}, '${p.nombre}', '${p.descripcion}', ${p.precio}, ${p.stock}, ${p.id_categoria}, ${p.id_proveedor})">Editar</button>
          <button onclick="eliminarProducto(${p.id_producto})">Eliminar</button>
        </td>
      </tr>
    `;
  });

  html += `</tbody></table>`;
  document.getElementById("productos").innerHTML = html;
}

function editarProducto(id, nombre, descripcion, precio, stock, idCategoria, idProveedor) {
  productoEditandoId = id;

  document.getElementById("productoNombre").value = nombre;
  document.getElementById("productoDescripcion").value = descripcion;
  document.getElementById("productoPrecio").value = precio;
  document.getElementById("productoStock").value = stock;
  document.getElementById("productoCategoria").value = idCategoria;
  document.getElementById("productoProveedor").value = idProveedor;
  document.getElementById("btnProducto").textContent = "Actualizar producto";
}

async function cargarReporteProductos() {
  const res = await fetch(`${API}/reporte/productos`);
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
}

async function cargarVentasDetalladas() {
  const res = await fetch(`${API}/ventas-detalladas`);
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

  const res = await fetch(`${API}/clientes`, {
    method: metodo,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body)
  });

  if (res.ok) {
    alert(mensajeOk);
    e.target.reset();
    clienteEditandoId = null;
    document.getElementById("btnCliente").textContent = "Agregar cliente";
    cargarClientes();
  } else {
    alert("Error al guardar cliente");
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

  const res = await fetch(`${API}/productos`, {
    method: metodo,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body)
  });

  if (res.ok) {
    alert(mensajeOk);
    e.target.reset();
    productoEditandoId = null;
    document.getElementById("btnProducto").textContent = "Agregar producto";
    cargarProductos();
  } else {
    alert("Error al guardar producto");
  }
});

async function eliminarCliente(id) {
  if (!confirm("¿Deseas eliminar este cliente?")) return;

  const res = await fetch(`${API}/clientes?id=${id}`, {
    method: "DELETE"
  });

  if (res.ok) {
    alert("Cliente eliminado correctamente");
    cargarClientes();
  } else {
    alert("No se pudo eliminar el cliente");
  }
}

async function eliminarProducto(id) {
  if (!confirm("¿Deseas eliminar este producto?")) return;

  const res = await fetch(`${API}/productos?id=${id}`, {
    method: "DELETE"
  });

  if (res.ok) {
    alert("Producto eliminado correctamente");
    cargarProductos();
  } else {
    alert("No se pudo eliminar el producto");
  }
}

cargarClientes();
cargarProductos();
cargarReporteProductos();
cargarVentasDetalladas();