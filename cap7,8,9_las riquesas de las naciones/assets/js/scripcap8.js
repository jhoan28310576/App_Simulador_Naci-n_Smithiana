let datosSimulacion = null;
let estadoActual = 'IA';
let anioActual = 2023;

function mostrarModalApi(key) {
    if (!datosSimulacion) return;
    let contenido = '';
    if (key === 'apiRaw1') contenido = JSON.stringify(datosSimulacion.apiRaw1, null, 2);
    if (key === 'apiRaw2') contenido = JSON.stringify(datosSimulacion.apiRaw2, null, 2);
    if (key.startsWith('apiRawDemanda')) {
        const anio = key.split('-')[1];
        contenido = JSON.stringify(datosSimulacion.apiRawDemanda[anio], null, 2);
    }
    document.getElementById('modalApiContent').innerText = contenido;
    const modal = new bootstrap.Modal(document.getElementById('modalApi'));
    modal.show();
}

function recargarSimulacion() {
    anioActual = document.getElementById('input-anio').value;
    estadoActual = document.getElementById('input-estado').value;

    if (!estadoActual) {
        mostrarMensaje('Por favor selecciona un estado', 'warning');
        return;
    }

    mostrarMensaje('Cargando datos...', 'info');
    cargarSimulacion();
}

function mostrarMensaje(texto, tipo) {
    const mensajeEl = document.getElementById('mensaje-estado');
    mensajeEl.className = `alert alert-${tipo}`;
    mensajeEl.innerText = texto;
    mensajeEl.style.display = 'block';

    if (tipo !== 'info') {
        setTimeout(() => {
            mensajeEl.style.display = 'none';
        }, 5000);
    }
}

function cargarSimulacion() {
    fetch(`/api/cap8/simulacion?anio=${anioActual}&estado=${estadoActual}`)
        .then(resp => resp.json())
        .then(data => {
            datosSimulacion = data;

            // Verificar si hay datos válidos
            if (data.valorProduccion === 0) {
                mostrarMensaje(`No se encontraron datos de producción de maíz para ${estadoActual} en ${anioActual}. Intenta con otro estado o año.`, 'warning');
                return;
            }

            mostrarMensaje('Datos cargados exitosamente', 'success');
            setTimeout(() => {
                document.getElementById('mensaje-estado').style.display = 'none';
            }, 3000);

            // Sección 1
            document.getElementById('salario-ajustado').innerText = `$${data.salarioAjustado.toFixed(2)}`;
            document.getElementById('valor-produccion').innerText = `$${data.valorProduccion.toLocaleString()}`;
            // Sección 2
            document.getElementById('num-ofertas').innerText = data.numOfertas;
            document.getElementById('valor-produccion-2').innerText = `$${data.valorProduccion.toLocaleString()}`;
            // Sección 3
            let html = '<table class="sim-data"><tr><th>Año</th><th>Producción</th><th>Demandas laborales</th><th>Ver datos API</th></tr>';
            data.demandaAnual.forEach(row => {
                html += `<tr>
                    <td>${row.anio}</td>
                    <td>$${row.valorProduccion.toLocaleString()}</td>
                    <td>${row.numDemandas}</td>
                    <td>
                        <button class="btn btn-secondary btn-sm" onclick="mostrarModalApi('apiRawDemanda-${row.anio}')">Ver datos</button>
                    </td>
                </tr>`;
            });
            html += '</table>';
            document.getElementById('demanda-anual-data').innerHTML = html;
            // Explicaciones
            document.getElementById('explicacion-salario').innerText =
                `El salario se calcula sumando un ajuste proporcional al valor de producción. Si la producción aumenta, el salario también.`;
            document.getElementById('explicacion-ofertas').innerText =
                `El número de ofertas laborales se calcula en función del valor de producción. Si la producción aumenta, hay más ofertas laborales.`;
            // Gráficos
            renderizarGraficos(data);
        })
        .catch(err => {
            mostrarMensaje('Error al cargar los datos. Verifica tu conexión e intenta nuevamente.', 'danger');
            document.getElementById('salario-ajustado').innerText = 'Error';
            document.getElementById('valor-produccion').innerText = 'Error';
            document.getElementById('num-ofertas').innerText = 'Error';
            document.getElementById('valor-produccion-2').innerText = 'Error';
            document.getElementById('demanda-anual-data').innerHTML = '<p class="text-danger">Error cargando datos</p>';
        });
}

// Renderizar gráficos usando Chart.js
let chartSalario = null, chartOfertas = null, chartDemanda = null;
function renderizarGraficos(data) {
    // Gráfico 1: Salario ajustado
    if (chartSalario) chartSalario.destroy();
    chartSalario = new Chart(document.getElementById('grafico-salario').getContext('2d'), {
        type: 'doughnut',
        data: {
            labels: ['Salario ajustado'],
            datasets: [{
                label: 'Salario ($)',
                data: [data.salarioAjustado, 100 - data.salarioAjustado],
                backgroundColor: [
                    '#667eea',
                    '#f8f9fa'
                ],
                borderColor: '#667eea',
                borderWidth: 2,
                cutout: '70%'
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false
                },
                tooltip: {
                    backgroundColor: 'rgba(0,0,0,0.8)',
                    titleColor: '#fff',
                    bodyColor: '#fff',
                    borderColor: '#667eea',
                    borderWidth: 1,
                    cornerRadius: 8,
                    displayColors: false,
                    callbacks: {
                        label: function (context) {
                            if (context.dataIndex === 0) {
                                return `Salario: $${context.parsed.toFixed(2)}`;
                            }
                            return null;
                        }
                    }
                }
            },
            animation: {
                animateRotate: true,
                animateScale: true,
                duration: 1500,
                easing: 'easeOutQuart'
            }
        }
    });

    // Gráfico 2: Ofertas laborales
    if (chartOfertas) chartOfertas.destroy();
    chartOfertas = new Chart(document.getElementById('grafico-ofertas').getContext('2d'), {
        type: 'bar',
        data: {
            labels: ['Ofertas laborales disponibles'],
            datasets: [{
                label: 'Ofertas',
                data: [data.numOfertas],
                backgroundColor: '#f093fb',
                borderColor: '#f093fb',
                borderWidth: 2,
                borderRadius: 8,
                borderSkipped: false,
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false
                },
                tooltip: {
                    backgroundColor: 'rgba(0,0,0,0.8)',
                    titleColor: '#fff',
                    bodyColor: '#fff',
                    borderColor: '#f093fb',
                    borderWidth: 1,
                    cornerRadius: 8,
                    displayColors: false,
                    callbacks: {
                        label: function (context) {
                            return `Ofertas: ${context.parsed.y}`;
                        }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    grid: {
                        color: 'rgba(0,0,0,0.1)',
                        drawBorder: false
                    },
                    ticks: {
                        color: '#666',
                        font: {
                            size: 12
                        }
                    }
                },
                x: {
                    grid: {
                        display: false
                    },
                    ticks: {
                        color: '#666',
                        font: {
                            size: 12
                        }
                    }
                }
            },
            animation: {
                duration: 1000,
                easing: 'easeOutQuart'
            }
        }
    });

    // Gráfico 3: Demanda laboral por año
    if (chartDemanda) chartDemanda.destroy();
    chartDemanda = new Chart(document.getElementById('grafico-demanda').getContext('2d'), {
        type: 'line',
        data: {
            labels: data.demandaAnual.map(row => row.anio),
            datasets: [{
                label: 'Demandas laborales',
                data: data.demandaAnual.map(row => row.numDemandas),
                backgroundColor: 'rgba(255, 99, 132, 0.2)',
                borderColor: '#ff6384',
                borderWidth: 3,
                fill: true,
                tension: 0.4,
                pointBackgroundColor: '#ff6384',
                pointBorderColor: '#fff',
                pointBorderWidth: 2,
                pointRadius: 6,
                pointHoverRadius: 8
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: true,
                    position: 'top',
                    labels: {
                        color: '#666',
                        font: {
                            size: 14,
                            weight: 'bold'
                        }
                    }
                },
                tooltip: {
                    backgroundColor: 'rgba(0,0,0,0.8)',
                    titleColor: '#fff',
                    bodyColor: '#fff',
                    borderColor: '#ff6384',
                    borderWidth: 1,
                    cornerRadius: 8,
                    displayColors: true
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    grid: {
                        color: 'rgba(0,0,0,0.1)',
                        drawBorder: false
                    },
                    ticks: {
                        color: '#666',
                        font: {
                            size: 12
                        }
                    }
                },
                x: {
                    grid: {
                        color: 'rgba(0,0,0,0.1)',
                        drawBorder: false
                    },
                    ticks: {
                        color: '#666',
                        font: {
                            size: 12
                        }
                    }
                }
            },
            animation: {
                duration: 1500,
                easing: 'easeOutQuart'
            }
        }
    });
}

// Inicializar con valores por defecto
cargarSimulacion();