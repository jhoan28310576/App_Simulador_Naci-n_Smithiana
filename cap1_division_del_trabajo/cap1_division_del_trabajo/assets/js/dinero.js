// JS para Simulador del Dinero (Capítulo 4)

document.addEventListener('DOMContentLoaded', () => {
    cargarValoresMonedas();
    cargarUsuarios();
    inicializarGraficas();
    setupConversion();
});

// Cargar valores actuales de monedas y llenar tabla
async function cargarValoresMonedas() {
    try {
        const res = await fetch('/api/monedas/valores');
        const data = await res.json();
        if (!data.success) throw new Error('No se pudo obtener valores de monedas');
        const tbody = document.getElementById('tablaValoresMonedas');
        tbody.innerHTML = '';
        for (const key of ['oro', 'plata', 'cobre']) {
            if (data.valores[key]) {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${data.valores[key].nombre}</td>
                    <td>${data.valores[key].valor_relativo_oro}</td>
                `;
                tbody.appendChild(row);
            }
        }
    } catch (err) {
        mostrarError('Error al cargar valores de monedas');
    }
}

// Cargar usuarios y mostrar sus saldos en monedas
async function cargarUsuarios() {
    try {
        const res = await fetch('/api/users');
        const data = await res.json();
        if (!data.success) throw new Error('No se pudo obtener usuarios');
        const tbody = document.getElementById('tablaUsuariosMonedas');
        tbody.innerHTML = '';
        data.users.forEach(u => {
            const inv = u.inventario;
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${u.nombre}</td>
                <td>${inv.oro ? inv.oro : 0}</td>
                <td>${inv.plata ? inv.plata : 0}</td>
                <td>${inv.cobre ? inv.cobre : 0}</td>
                <td>${inv.dinero ? inv.dinero : 0}</td>
            `;
            tbody.appendChild(row);
        });
    } catch (err) {
        mostrarError('Error al cargar usuarios');
    }
}

// Configurar formulario de conversión
function setupConversion() {
    const form = document.getElementById('conversionForm');
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const cantidad = parseFloat(document.getElementById('cantidadInput').value);
        const origen = document.getElementById('origenSelect').value;
        const destino = document.getElementById('destinoSelect').value;
        const resultadoDiv = document.getElementById('conversionResultado');
        resultadoDiv.textContent = '';
        if (isNaN(cantidad) || cantidad < 0) {
            mostrarError('Cantidad inválida');
            return;
        }
        if (origen === destino) {
            resultadoDiv.textContent = 'Selecciona monedas diferentes para convertir.';
            return;
        }
        try {
            const res = await fetch(`/api/monedas/convertir/${cantidad}/${origen}/${destino}`);
            const data = await res.json();
            if (!data.success) throw new Error('No se pudo convertir');
            resultadoDiv.textContent = `${cantidad} ${origen} equivalen a ${data.resultado.toFixed(4)} ${destino}`;
        } catch (err) {
            mostrarError('Error en la conversión');
        }
    });
}

// Inicializar y cargar gráficas de historial
function inicializarGraficas() {
    cargarHistorialYGraficar('oro', 'graficaOro', '#FFD700');
    cargarHistorialYGraficar('plata', 'graficaPlata', '#C0C0C0');
    cargarHistorialYGraficar('cobre', 'graficaCobre', '#B87333');
}

async function cargarHistorialYGraficar(moneda, canvasId, color) {
    try {
        const res = await fetch(`/api/monedas/historial/${moneda}`);
        const data = await res.json();
        if (!data.success) throw new Error('No se pudo obtener historial');
        const ctx = document.getElementById(canvasId).getContext('2d');
        // Si ya existe una gráfica previa, destrúyela
        if (ctx.chartInstance) {
            ctx.chartInstance.destroy();
        }
        ctx.chartInstance = new Chart(ctx, {
            type: 'line',
            data: {
                labels: data.historial.map((_, i) => `T${i+1}`),
                datasets: [{
                    label: `Valor de ${moneda}`,
                    data: data.historial,
                    borderColor: color,
                    backgroundColor: color + '22',
                    fill: true,
                    tension: 0.3
                }]
            },
            options: {
                responsive: true,
                plugins: { legend: { display: false } },
                scales: {
                    y: { beginAtZero: true }
                }
            }
        });
    } catch (err) {
        mostrarError(`Error al cargar historial de ${moneda}`);
    }
}

// Mostrar errores de forma clara
function mostrarError(msg) {
    alert(msg);
} 