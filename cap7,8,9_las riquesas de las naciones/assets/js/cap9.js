// Actualizar valores de los sliders
document.getElementById('riesgo-slider').addEventListener('input', function () {
    const value = this.value;
    document.getElementById('riesgo-value').textContent = value;
    document.getElementById('riesgo').value = value;
});

document.getElementById('competencia-slider').addEventListener('input', function () {
    const value = parseFloat(this.value).toFixed(2);
    document.getElementById('competencia-value').textContent = value;
    document.getElementById('competencia').value = value;
});

// Inicializar gráfico vacío
const ctx = document.getElementById('grafico').getContext('2d');
window.capitalChart = new Chart(ctx, {
    type: 'line',
    data: {
        labels: [],
        datasets: [
            {
                label: 'Evolución del Capital',
                data: [],
                borderColor: '#3498db',
                backgroundColor: 'rgba(52, 152, 219, 0.1)',
                yAxisID: 'y',
                fill: true,
                tension: 0.3,
                borderWidth: 3
            },
            {
                label: 'Retorno Anual',
                data: [],
                borderColor: '#2ecc71',
                backgroundColor: 'rgba(46, 204, 113, 0.1)',
                yAxisID: 'y1',
                type: 'bar',
                borderWidth: 2
            }
        ]
    },
    options: {
        responsive: true,
        interaction: { mode: 'index', intersect: false },
        scales: {
            y: {
                type: 'linear',
                position: 'left',
                title: {
                    text: 'Capital ($)',
                    display: true,
                    font: { weight: 'bold', size: 14 },
                    color: '#2c3e50'
                },
                beginAtZero: true,
                grid: {
                    color: 'rgba(0,0,0,0.05)'
                },
                ticks: {
                    callback: function (value) {
                        return '$' + value.toLocaleString('es-ES');
                    }
                }
            },
            y1: {
                type: 'linear',
                position: 'right',
                title: {
                    text: 'Retorno ($)',
                    display: true,
                    font: { weight: 'bold', size: 14 },
                    color: '#2c3e50'
                },
                grid: {
                    drawOnChartArea: false
                },
                beginAtZero: true,
                ticks: {
                    callback: function (value) {
                        return '$' + value.toLocaleString('es-ES');
                    }
                }
            }
        },
        plugins: {
            tooltip: {
                callbacks: {
                    label: function (context) {
                        let label = context.dataset.label || '';
                        if (label) label += ': ';
                        if (context.parsed.y !== null) {
                            label += new Intl.NumberFormat('es-ES', {
                                style: 'currency',
                                currency: 'USD'
                            }).format(context.parsed.y);
                        }
                        return label;
                    }
                }
            },
            legend: {
                position: 'top',
                labels: {
                    font: { size: 14 },
                    padding: 20
                }
            },
            title: {
                display: true,
                text: 'Evolución de la Inversión',
                font: { size: 18, weight: 'normal' },
                padding: {
                    top: 10,
                    bottom: 20
                }
            }
        }
    }
});

// Manejar el envío del formulario
document.getElementById('simulacion-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const data = {
        capital: parseFloat(document.getElementById('capital').value),
        sector: document.getElementById('sector').value,
        riesgo: parseFloat(document.getElementById('riesgo').value),
        competencia: parseFloat(document.getElementById('competencia').value),
        salario_medio: parseFloat(document.getElementById('salario_medio').value) || 0,
        pais: document.getElementById('pais').value,
    };

    const anios = document.getElementById('anios').value;

    try {
        // Mostrar estado de carga
        document.getElementById('results-section').style.display = 'block';
        document.querySelector('#result-table tbody').innerHTML = '<tr><td colspan="4" class="text-center py-4">Calculando simulación...</td></tr>';
        document.getElementById('salario-usado').textContent = 'Calculando...';
        document.getElementById('analisis-resultados').innerHTML = '<p>Analizando resultados según principios de Adam Smith...</p>';

        const response = await fetch(`/api/simular-retorno?anios=${anios}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Error en la simulación');
        }

        const result = await response.json();
        mostrarResultados(result);
    } catch (error) {
        alert(error.message || 'Error en la solicitud');
    }
});

function mostrarResultados(result) {
    // Actualizar salario usado
    document.getElementById('salario-usado').textContent =
        result.salario_medio.toLocaleString('es-ES', {
            style: 'currency',
            currency: 'USD'
        });

    // Preparar datos para tabla y gráfico
    const historial = result.historial;
    const labels = [];
    const capitalData = [];
    const retornoData = [];

    // Limpiar tabla
    const tableBody = document.querySelector('#result-table tbody');
    tableBody.innerHTML = '';

    // Verificar si hay datos
    if (!historial || historial.length === 0) {
        tableBody.innerHTML = '<tr><td colspan="4" class="text-center">No se recibieron datos de la simulación</td></tr>';
        return;
    }

    // Procesar datos
    let prevCapital = historial[0].Capital || historial[0].capital;
    historial.forEach((item, index) => {
        // Manejar campos en mayúsculas o minúsculas
        const año = item.Año !== undefined ? item.Año : item.año;
        const capital = item.Capital !== undefined ? item.Capital : item.capital;
        const retorno = item.Retorno !== undefined ? item.Retorno : item.retorno;

        // Calcular crecimiento porcentual
        let crecimiento = 0;
        if (index > 0) {
            crecimiento = ((capital - prevCapital) / prevCapital) * 100;
            prevCapital = capital;
        }

        labels.push(`Año ${año}`);
        capitalData.push(capital);
        retornoData.push(retorno);

        // Agregar fila a la tabla
        const row = document.createElement('tr');
        row.innerHTML = `
                <td>${año}</td>
                <td>${capital.toLocaleString('es-ES', { style: 'currency', currency: 'USD' })}</td>
                <td>${retorno.toLocaleString('es-ES', { style: 'currency', currency: 'USD' })}</td>
                <td>${index > 0 ? crecimiento.toFixed(2) + '%' : '-'}</td>
            `;
        tableBody.appendChild(row);
    });

    // Actualizar gráfico
    window.capitalChart.data.labels = labels;
    window.capitalChart.data.datasets[0].data = capitalData;
    window.capitalChart.data.datasets[1].data = retornoData;
    window.capitalChart.update();

    // Generar análisis
    generarAnalisis(result);
}

function generarAnalisis(result) {
    const historial = result.historial;
    if (!historial || historial.length === 0) return;

    const primerAño = historial[0];
    const ultimoAño = historial[historial.length - 1];
    const capitalInicial = primerAño.Capital || primerAño.capital;
    const capitalFinal = ultimoAño.Capital || ultimoAño.capital;
    const crecimiento = ((capitalFinal - capitalInicial) / capitalInicial) * 100;

    const riesgo = document.getElementById('riesgo').value;
    const competencia = document.getElementById('competencia').value;
    const sector = document.getElementById('sector').value;
    const salario = result.salario_medio;

    // Obtener nombre completo del sector
    const sectores = {
        agricultura: "Agricultura",
        manufactura: "Manufactura",
        comercio: "Comercio"
    };

    const sectorNombre = sectores[sector] || sector;

    // Análisis de riesgo
    let riesgoAnalisis = "";
    if (riesgo > 3) {
        riesgoAnalisis = "El alto nivel de riesgo seleccionado (nivel " + riesgo + ") ha contribuido a un mayor beneficio potencial, siguiendo el principio de Smith de que 'el beneficio es la recompensa por asumir riesgos'.";
    } else {
        riesgoAnalisis = "El moderado nivel de riesgo (nivel " + riesgo + ") ha generado beneficios más conservadores, de acuerdo con la teoría de Smith sobre la relación riesgo-beneficio.";
    }

    // Análisis de competencia
    let competenciaAnalisis = "";
    if (competencia > 0.7) {
        competenciaAnalisis = "La alta competencia en el mercado (" + competencia + ") ha reducido significativamente los beneficios, como Smith describió: 'La competencia reduce los beneficios al mínimo necesario'.";
    } else if (competencia > 0.4) {
        competenciaAnalisis = "La competencia moderada (" + competencia + ") ha permitido beneficios razonables, equilibrando el principio de Smith sobre los efectos de la competencia.";
    } else {
        competenciaAnalisis = "La baja competencia (" + competencia + ") ha permitido mayores beneficios, reflejando la situación de cuasi-monopolio que Smith analizó en su obra.";
    }

    // Análisis de sector
    let sectorAnalisis = "";
    if (sector === "comercio") {
        sectorAnalisis = "El sector comercial, como Smith observó, tiende a generar mayores beneficios debido a su naturaleza dinámica y oportunidades de crecimiento.";
    } else if (sector === "manufactura") {
        sectorAnalisis = "El sector manufacturero, según Smith, ofrece beneficios estables y consistentes, aunque generalmente menores que el comercio.";
    } else {
        sectorAnalisis = "El sector agrícola, como Smith notó, suele ofrecer beneficios más bajos pero estables, con menor volatilidad.";
    }

    // Análisis de salario
    let salarioAnalisis = "";
    if (salario > 30000) {
        salarioAnalisis = "Los altos salarios (" + salario.toLocaleString('es-ES', { style: 'currency', currency: 'USD' }) + ") han presionado a la baja los beneficios, confirmando la relación inversa entre salarios y beneficios que Smith describió.";
    } else {
        salarioAnalisis = "Los salarios moderados (" + salario.toLocaleString('es-ES', { style: 'currency', currency: 'USD' }) + ") han permitido mantener beneficios saludables, ilustrando la relación inversa entre salarios y beneficios.";
    }

    // Análisis general
    const analisisHTML = `
            <p>La inversión inicial de ${capitalInicial.toLocaleString('es-ES', { style: 'currency', currency: 'USD' })} 
            creció a ${capitalFinal.toLocaleString('es-ES', { style: 'currency', currency: 'USD' })} 
            en ${historial.length} años, un crecimiento del ${crecimiento.toFixed(2)}%.</p>
            
            <p>Según los principios de Adam Smith en el Capítulo 9:</p>
            <ul>
                <li>${riesgoAnalisis}</li>
                <li>${competenciaAnalisis}</li>
                <li>En el sector ${sectorNombre}, ${sectorAnalisis}</li>
                <li>${salarioAnalisis}</li>
            </ul>
            
            <div class="adamsmith-quote">
                "El beneficio es la recompensa al capitalista por asumir el riesgo y la incomodidad de invertir" - Adam Smith, 1776
            </div>
        `;

    document.getElementById('analisis-resultados').innerHTML = analisisHTML;
}

// Simular un clic al cargar la página (solo para demostración)
window.addEventListener('DOMContentLoaded', () => {
    setTimeout(() => {
        document.querySelector('button[type="submit"]').click();
    }, 1000);
});