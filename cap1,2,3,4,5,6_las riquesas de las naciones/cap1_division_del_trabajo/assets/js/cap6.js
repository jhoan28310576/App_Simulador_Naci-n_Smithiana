let productos = [];
let productoActual = null;
let historialChart = null;

// Cargar productos al iniciar
document.addEventListener('DOMContentLoaded', function () {
    cargarProductos();
    cargarIndicadoresEconomicos();
});

// Función para cargar productos
async function cargarProductos() {
    try {
        const response = await fetch('/api/cap6/productos');
        const data = await response.json();

        if (data.success) {
            productos = data.productos;
            mostrarListaProductos();
            actualizarSelectProductos();
        }
    } catch (error) {
        console.error('Error al cargar productos:', error);
    }
}

// Función para mostrar lista de productos
function mostrarListaProductos() {
    const container = document.getElementById('listaProductos');
    let html = '<div class="table-responsive"><table class="table">';
    html += '<thead><tr><th>Producto</th><th>Precio Mercado</th><th>Precio Natural</th><th>Oferta</th><th>Demanda</th><th>País</th><th>Acciones</th></tr></thead><tbody>';

    productos.forEach(producto => {
        html += `
            <tr>
                <td><strong>${producto.nombre}</strong></td>
                <td>$${producto.precio_mercado.toFixed(2)}</td>
                <td>$${producto.precio_natural.toFixed(2)}</td>
                <td>${producto.oferta}</td>
                <td>${producto.demanda}</td>
                <td>${producto.pais}</td>
                <td>
                    <button class="btn btn-sm btn-primary" onclick="seleccionarProducto('${producto.id}')">
                        <i class="fas fa-eye"></i> Ver
                    </button>
                </td>
            </tr>
        `;
    });

    html += '</tbody></table></div>';
    container.innerHTML = html;
}

// Función para actualizar select de productos
function actualizarSelectProductos() {
    const select = document.getElementById('productoSelect');
    select.innerHTML = '<option value="">Selecciona un producto...</option>';

    productos.forEach(producto => {
        select.innerHTML += `<option value="${producto.id}">${producto.nombre} (${producto.pais})</option>`;
    });
}

// Función para cargar indicadores económicos
async function cargarIndicadoresEconomicos() {
    try {
        const response = await fetch('/api/cap6/indicadores');
        const data = await response.json();

        if (data.success) {
            mostrarIndicadoresEconomicos(data.indicadores);
        }
    } catch (error) {
        console.error('Error al cargar indicadores:', error);
    }
}

// Función para mostrar indicadores económicos
function mostrarIndicadoresEconomicos(indicadores) {
    const container = document.getElementById('indicadoresEconomicos');
    let html = '';

    for (const [pais, datos] of Object.entries(indicadores)) {
        html += `
            <div class="mb-3">
                <h6>${pais.charAt(0).toUpperCase() + pais.slice(1)}</h6>
                <p class="mb-1"><small>Inflación: ${datos.inflacion_anual}%</small></p>
                <p class="mb-1"><small>PIB per cápita: $${datos.pib_per_capita}</small></p>
                <p class="mb-1"><small>Desempleo: ${datos.tasa_desempleo}%</small></p>
                <p class="mb-0"><small>Salario mínimo: $${datos.salario_minimo}</small></p>
            </div>
        `;
    }

    container.innerHTML = html;
}

// Función para cargar producto específico
async function cargarProducto() {
    const productoId = document.getElementById('productoSelect').value;
    if (!productoId) return;

    await seleccionarProducto(productoId);
}

// Función para seleccionar producto
async function seleccionarProducto(productoId) {
    try {
        const response = await fetch(`/api/cap6/producto/${productoId}`);
        const data = await response.json();

        if (data.success) {
            productoActual = data.producto;
            mostrarProducto(productoActual);
            document.getElementById('productoSelect').value = productoId;
        }
    } catch (error) {
        console.error('Error al cargar producto:', error);
    }
}

// Función para mostrar producto
function mostrarProducto(producto) {
    // Mostrar secciones
    document.getElementById('infoProducto').style.display = 'block';
    document.getElementById('componentesPrecio').style.display = 'block';
    document.getElementById('simuladorMercado').style.display = 'block';
    document.getElementById('historialPrecios').style.display = 'block';

    // Actualizar información básica
    document.getElementById('precioNatural').textContent = producto.precio_natural.toFixed(2);
    document.getElementById('precioMercado').textContent = producto.precio_mercado.toFixed(2);
    document.getElementById('diferenciaPrecio').textContent = (producto.precio_mercado - producto.precio_natural).toFixed(2);
    document.getElementById('oferta').textContent = producto.oferta;
    document.getElementById('demanda').textContent = producto.demanda;
    document.getElementById('pais').textContent = producto.pais;

    // Actualizar componentes
    document.getElementById('salarios').textContent = producto.componentes.salarios.toFixed(2);
    document.getElementById('beneficios').textContent = producto.componentes.beneficios.toFixed(2);
    document.getElementById('rentas').textContent = producto.componentes.rentas.toFixed(2);

    // Calcular porcentajes
    const total = producto.precio_mercado;
    const porcentajeSalarios = (producto.componentes.salarios / total) * 100;
    const porcentajeBeneficios = (producto.componentes.beneficios / total) * 100;
    const porcentajeRentas = (producto.componentes.rentas / total) * 100;

    // Actualizar barras de progreso
    document.getElementById('salariosBar').style.width = porcentajeSalarios + '%';
    document.getElementById('salariosBar').textContent = porcentajeSalarios.toFixed(1) + '%';

    document.getElementById('beneficiosBar').style.width = porcentajeBeneficios + '%';
    document.getElementById('beneficiosBar').textContent = porcentajeBeneficios.toFixed(1) + '%';

    document.getElementById('rentasBar').style.width = porcentajeRentas + '%';
    document.getElementById('rentasBar').textContent = porcentajeRentas.toFixed(1) + '%';

    // Actualizar valores de simulación
    document.getElementById('nuevaOferta').value = producto.oferta;
    document.getElementById('nuevaDemanda').value = producto.demanda;

    // Crear gráfico de historial
    crearGraficoHistorial(producto);
}

// Función para crear gráfico de historial
function crearGraficoHistorial(producto) {
    const ctx = document.getElementById('historialChart').getContext('2d');

    if (historialChart) {
        historialChart.destroy();
    }

    const labels = [];
    for (let i = 0; i < producto.historial_precios.natural.length; i++) {
        labels.push(`Período ${i + 1}`);
    }

    historialChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Precio Natural',
                data: producto.historial_precios.natural,
                borderColor: '#4ecdc4',
                backgroundColor: 'rgba(78, 205, 196, 0.1)',
                tension: 0.1
            }, {
                label: 'Precio de Mercado',
                data: producto.historial_precios.mercado,
                borderColor: '#ff6b6b',
                backgroundColor: 'rgba(255, 107, 107, 0.1)',
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: {
                    display: true,
                    text: 'Evolución de Precios',
                    color: 'white'
                },
                legend: {
                    labels: {
                        color: 'white'
                    }
                }
            },
            scales: {
                y: {
                    ticks: {
                        color: 'white'
                    },
                    grid: {
                        color: 'rgba(255, 255, 255, 0.1)'
                    }
                },
                x: {
                    ticks: {
                        color: 'white'
                    },
                    grid: {
                        color: 'rgba(255, 255, 255, 0.1)'
                    }
                }
            }
        }
    });
}

// Función para simular mercado
async function simularMercado() {
    if (!productoActual) {
        alert('Selecciona un producto primero');
        return;
    }

    const nuevaOferta = parseInt(document.getElementById('nuevaOferta').value);
    const nuevaDemanda = parseInt(document.getElementById('nuevaDemanda').value);

    if (!nuevaOferta || !nuevaDemanda) {
        alert('Ingresa valores válidos para oferta y demanda');
        return;
    }

    try {
        const response = await fetch(`/api/cap6/simular/${productoActual.id}?oferta=${nuevaOferta}&demanda=${nuevaDemanda}`);
        const data = await response.json();

        if (data.success) {
            productoActual = data.producto;
            mostrarProducto(productoActual);

            // Mostrar resultado de simulación
            const cambios = data.cambios;
            document.getElementById('resultadoSimulacion').innerHTML = `
                <div class="alert alert-success">
                    <h6>Simulación Completada</h6>
                    <p><strong>Precio anterior:</strong> $${cambios.precio_anterior.toFixed(2)}</p>
                    <p><strong>Precio nuevo:</strong> $${cambios.precio_nuevo.toFixed(2)}</p>
                    <p><strong>Cambio:</strong> $${(cambios.precio_nuevo - cambios.precio_anterior).toFixed(2)}</p>
                </div>
            `;
        }
    } catch (error) {
        console.error('Error al simular mercado:', error);
        alert('Error al simular mercado');
    }
}