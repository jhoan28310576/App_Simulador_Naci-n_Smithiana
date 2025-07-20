let inflationChart;

// Productos por país (ajusta según tus datos reales)
const productosPorPais = {
    VE: [
        { id: "trigo", nombre: "Trigo" },
        { id: "ropa", nombre: "Ropa" },
        { id: "vivienda", nombre: "Vivienda" }
    ],
    CO: [
        { id: "herramientas", nombre: "Herramientas" }
    ]
};

// Cargar datos al iniciar
document.addEventListener('DOMContentLoaded', function () {
    cargarProductos();
    cargarDatosInflacion();
    actualizarSelectProductos();
});

// Función para cargar productos
async function cargarProductos() {
    document.getElementById('loadingProductos').style.display = 'block';
    try {
        const response = await fetch('/api/precios/productos');
        const data = await response.json();

        if (data.success) {
            mostrarProductos(data.productos);
        }
    } catch (error) {
        console.error('Error al cargar productos:', error);
    } finally {
        document.getElementById('loadingProductos').style.display = 'none';
    }
}

// Función para mostrar productos en tabla
function mostrarProductos(productos) {
    const container = document.getElementById('productosTable');
    let html = `
        <table class="table table-striped">
            <thead>
                <tr>
                    <th>Producto</th>
                    <th>Precio Real (horas)</th>
                    <th>Precio Nominal (USD)</th>
                    <th>Relación Nominal/Real</th>
                </tr>
            </thead>
            <tbody>
    `;

    for (const [id, producto] of Object.entries(productos)) {
        html += `
            <tr>
                <td><strong>${producto.nombre}</strong></td>
                <td>${producto.precio_real} horas</td>
                <td>$${producto.precio_nominal.toFixed(2)}</td>
                <td>${producto.relacion.toFixed(2)}</td>
            </tr>
        `;
    }

    html += '</tbody></table>';
    container.innerHTML = html;
}

// Función para cargar datos de inflación
async function cargarDatosInflacion() {
    try {
        const response = await fetch('/api/precios/inflacion');
        const data = await response.json();

        if (data.success) {
            mostrarEstadisticasInflacion(data.estadisticas);
            crearGraficoInflacion(data.estadisticas);
        }
    } catch (error) {
        console.error('Error al cargar datos de inflación:', error);
    }
}

// Función para mostrar estadísticas de inflación
function mostrarEstadisticasInflacion(estadisticas) {
    if (estadisticas.error) {
        document.getElementById('venezuelaStats').innerHTML = '<p class="text-warning">Error al cargar datos</p>';
        document.getElementById('colombiaStats').innerHTML = '<p class="text-warning">Error al cargar datos</p>';
        return;
    }

    const venezuela = estadisticas.venezuela;
    const colombia = estadisticas.colombia;
    const comparacion = estadisticas.comparacion;

    document.getElementById('venezuelaStats').innerHTML = `
        <p><strong>Factor de Inflación:</strong> ${venezuela.factor_inflacion.toFixed(2)}</p>
        <p><strong>Inflación Acumulada:</strong> ${((venezuela.factor_inflacion - 1) * 100).toFixed(1)}%</p>
    `;

    document.getElementById('colombiaStats').innerHTML = `
        <p><strong>Factor de Inflación:</strong> ${colombia.factor_inflacion.toFixed(2)}</p>
        <p><strong>Inflación Acumulada:</strong> ${((colombia.factor_inflacion - 1) * 100).toFixed(1)}%</p>
    `;
}

// Función para crear gráfico de inflación
function crearGraficoInflacion(estadisticas) {
    const ctx = document.getElementById('inflationChart').getContext('2d');

    if (inflationChart) {
        inflationChart.destroy();
    }

    const venezuelaData = estadisticas.venezuela.datos_recientes.map(d => d.valor);
    const colombiaData = estadisticas.colombia.datos_recientes.map(d => d.valor);
    const labels = estadisticas.venezuela.datos_recientes.map(d => d.ano);

    inflationChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Venezuela',
                data: venezuelaData,
                borderColor: '#ff6b6b',
                backgroundColor: 'rgba(255, 107, 107, 0.1)',
                tension: 0.1
            }, {
                label: 'Colombia',
                data: colombiaData,
                borderColor: '#4ecdc4',
                backgroundColor: 'rgba(78, 205, 196, 0.1)',
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: {
                    display: true,
                    text: 'Inflación Anual (%) - Últimos 5 años',
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

// Función para actualizar precios
async function actualizarPrecios(pais) {
    try {
        const response = await fetch(`/api/precios/actualizar/${pais}`);
        const data = await response.json();

        if (data.success) {
            alert(`Precios actualizados para ${pais}. Factor de inflación: ${data.factor_inflacion.toFixed(2)}`);
            cargarProductos();
        } else {
            alert('Error al actualizar precios: ' + data.error);
        }
    } catch (error) {
        console.error('Error al actualizar precios:', error);
        alert('Error al actualizar precios');
    }
}

// Función para calcular poder adquisitivo
async function calcularPoderAdquisitivo() {
    const producto = document.getElementById('productoSelect').value;
    const cantidad = document.getElementById('cantidadDinero').value;

    try {
        const response = await fetch(`/api/precios/poder-adquisitivo/${producto}/${cantidad}`);
        const data = await response.json();

        if (data.success) {
            document.getElementById('resultadoCalculadora').innerHTML = `
                <div class="alert alert-light">
                    <h5>Resultado:</h5>
                    <p><strong>Producto:</strong> ${data.producto}</p>
                    <p><strong>Cantidad de dinero:</strong> $${data.cantidad_dinero} ${data.moneda}</p>
                    <p><strong>Poder adquisitivo:</strong> ${data.poder_adquisitivo.toFixed(2)} unidades</p>
                </div>
            `;
        } else {
            alert('Error al calcular poder adquisitivo: ' + data.error);
        }
    } catch (error) {
        console.error('Error al calcular poder adquisitivo:', error);
        alert('Error al calcular poder adquisitivo');
    }
}

// Función para cargar historial
async function cargarHistorial() {
    const producto = document.getElementById('productoHistorial').value;
    if (!producto) return;
    try {
        const response = await fetch(`/api/precios/historial/${producto}`);
        const data = await response.json();
        if (data.success) {
            mostrarHistorial(data.historial, data.producto);
        } else {
            alert('Error al cargar historial: ' + data.error);
        }
    } catch (error) {
        console.error('Error al cargar historial:', error);
        alert('Error al cargar historial');
    }
}

// Función para mostrar historial
function mostrarHistorial(historial, producto) {
    console.log("Historial recibido:", historial, "Producto:", producto);
    if (!Array.isArray(historial)) historial = [];
    const container = document.getElementById('historialTable');

    if (historial.length === 0) {
        container.innerHTML = '<p class="text-muted">No hay historial disponible para este producto.</p>';
        return;
    }

    let html = `
        <h4>Historial de ${producto}</h4>
        <table class="table table-striped">
            <thead>
                <tr>
                    <th>Fecha</th>
                    <th>Precio Nominal</th>
                    <th>Factor Inflación</th>
                    <th>Precio Real</th>
                </tr>
            </thead>
            <tbody>
    `;

    historial.forEach(item => {
        const fecha = new Date(item.fecha).toLocaleDateString();
        html += `
            <tr>
                <td>${fecha}</td>
                <td>$${item.precio_nominal.toFixed(2)}</td>
                <td>${item.factor_inflacion.toFixed(2)}</td>
                <td>${item.precio_real.toFixed(2)} horas</td>
            </tr>
        `;
    });

    html += '</tbody></table>';
    container.innerHTML = html;
}

function actualizarSelectProductos() {
    const pais = document.getElementById('paisHistorial').value;
    const select = document.getElementById('productoHistorial');
    select.innerHTML = productosPorPais[pais].map(p => `<option value="${p.id}">${p.nombre}</option>`).join('');
    cargarHistorial(); // Carga el historial del primer producto del país seleccionado
}