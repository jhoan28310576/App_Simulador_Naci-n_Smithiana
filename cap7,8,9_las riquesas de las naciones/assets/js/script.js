// Actualizar valores de los sliders
document.getElementById('droughtSeverity').addEventListener('input', function () {
    document.getElementById('droughtSeverityValue').textContent = this.value + '%';
});

document.getElementById('affectedArea').addEventListener('input', function () {
    document.getElementById('affectedAreaValue').textContent = this.value + '%';
});

async function runSimulation() {
    const year = document.getElementById('year').value;
    const states = Array.from(document.getElementById('states').selectedOptions).map(option => option.value);
    const droughtSeverity = parseFloat(document.getElementById('droughtSeverity').value) / 100;
    const affectedArea = parseFloat(document.getElementById('affectedArea').value) / 100;

    // Mostrar loading
    document.getElementById('loading').style.display = 'block';
    document.getElementById('results').style.display = 'none';
    document.getElementById('error').style.display = 'none';

    try {
        const response = await fetch('/api/drought-simulation', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                year: year,
                states: states,
                drought_severity: droughtSeverity,
                affected_area: affectedArea
            })
        });

        const data = await response.json();

        if (response.ok) {
            displayResults(data.result);
        } else {
            throw new Error(data.error || 'Error en la simulación');
        }
    } catch (error) {
        document.getElementById('error').textContent = 'Error: ' + error.message;
        document.getElementById('error').style.display = 'block';
    } finally {
        document.getElementById('loading').style.display = 'none';
    }
}

function displayResults(result) {
    // Formatear números
    const formatNumber = (num) => {
        if (num >= 1e9) return (num / 1e9).toFixed(2) + ' B';
        if (num >= 1e6) return (num / 1e6).toFixed(2) + ' M';
        if (num >= 1e3) return (num / 1e3).toFixed(2) + ' K';
        return num.toFixed(2);
    };

    const formatPercentage = (num) => (num * 100).toFixed(2) + '%';
    const formatCurrency = (num) => '$' + formatNumber(num);

    // Actualizar métricas principales
    document.getElementById('originalProduction').textContent = formatNumber(result.original_production) + ' bushels';
    document.getElementById('simulatedProduction').textContent = formatNumber(result.simulated_production) + ' bushels';
    document.getElementById('productionLoss').textContent = formatNumber(result.production_loss) + ' bushels';
    document.getElementById('priceIncrease').textContent = formatPercentage(result.price_increase);
    document.getElementById('economicImpact').textContent = formatCurrency(result.economic_impact);

    // Mostrar datos por estado
    const statesGrid = document.getElementById('statesGrid');
    statesGrid.innerHTML = '';

    result.states.forEach(state => {
        const stateCard = document.createElement('div');
        stateCard.className = 'state-card';
        stateCard.innerHTML = `
                    <div class="state-name">${state.state}</div>
                    <div class="state-production">${formatNumber(state.production)} ${state.unit}</div>
                `;
        statesGrid.appendChild(stateCard);
    });

    document.getElementById('results').style.display = 'block';
}

// Ejecutar simulación inicial al cargar la página
window.addEventListener('load', function () {
    setTimeout(runSimulation, 1000);
});