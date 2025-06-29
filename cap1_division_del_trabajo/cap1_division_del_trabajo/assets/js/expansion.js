// Módulo de Expansión de Mercado - JavaScript

class ExpansionSystem {
    constructor() {
        this.market = {
            radius: 10,
            population: 50,
            specializationLevel: 1,
            marketValue: 1000,
            infrastructure: {
                roads: 0,
                ports: 0,
                markets: 0,
                warehouses: 0
            },
            history: []
        };
        
        this.specializations = {
            1: ['Agricultor', 'Herrero', 'Carpintero'],
            2: ['Agricultor', 'Herrero', 'Carpintero', 'Alfarero', 'Tejedor', 'Panadero', 'Carnicero', 'Pescador'],
            3: ['Agricultor', 'Herrero', 'Carpintero', 'Alfarero', 'Tejedor', 'Panadero', 'Carnicero', 'Pescador', 
                'Sastre', 'Zapatero', 'Joyero', 'Pintor', 'Escultor', 'Músico', 'Médico', 'Abogado', 'Maestro', 'Sacerdote'],
            4: ['Agricultor', 'Herrero', 'Carpintero', 'Alfarero', 'Tejedor', 'Panadero', 'Carnicero', 'Pescador', 
                'Sastre', 'Zapatero', 'Joyero', 'Pintor', 'Escultor', 'Músico', 'Médico', 'Abogado', 'Maestro', 'Sacerdote',
                'Ingeniero', 'Arquitecto', 'Contador', 'Banquero', 'Comerciante', 'Marinero', 'Soldado', 'Policía', 'Bombero',
                'Periodista', 'Escritor', 'Actor', 'Cocinero', 'Camarero', 'Taxista', 'Conductor', 'Mecánico', 'Electricista'],
            5: ['Agricultor', 'Herrero', 'Carpintero', 'Alfarero', 'Tejedor', 'Panadero', 'Carnicero', 'Pescador', 
                'Sastre', 'Zapatero', 'Joyero', 'Pintor', 'Escultor', 'Músico', 'Médico', 'Abogado', 'Maestro', 'Sacerdote',
                'Ingeniero', 'Arquitecto', 'Contador', 'Banquero', 'Comerciante', 'Marinero', 'Soldado', 'Policía', 'Bombero',
                'Periodista', 'Escritor', 'Actor', 'Cocinero', 'Camarero', 'Taxista', 'Conductor', 'Mecánico', 'Electricista',
                'Programador', 'Diseñador', 'Marketing', 'Ventas', 'Recursos Humanos', 'Logística', 'Investigador', 'Científico',
                'Profesor Universitario', 'Piloto', 'Astrónomo', 'Biólogo', 'Químico', 'Físico', 'Matemático', 'Historiador', 'Arqueólogo']
        };
        
        this.nextSpecializations = {
            1: ['Alfarero', 'Tejedor'],
            2: ['Sastre', 'Zapatero', 'Joyero'],
            3: ['Ingeniero', 'Arquitecto', 'Contador'],
            4: ['Programador', 'Diseñador', 'Marketing'],
            5: ['Investigador', 'Científico', 'Profesor Universitario']
        };
        
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.setupTabs();
        this.updateDisplay();
        this.updateMarketVisualization();
        this.updateSpecializationLevels();
        this.initializeCharts();
    }

    setupEventListeners() {
        // Sliders de infraestructura
        const sliders = ['roads', 'ports', 'markets', 'warehouses'];
        sliders.forEach(type => {
            const slider = document.getElementById(`${type}Input`);
            const valueSpan = document.getElementById(`${type}Value`);
            
            slider.addEventListener('input', (e) => {
                this.market.infrastructure[type] = parseInt(e.target.value);
                valueSpan.textContent = e.target.value;
                this.calculateMarketExpansion();
            });
        });

        // Botones de acción
        document.getElementById('expandBtn').addEventListener('click', () => {
            this.expandMarket();
        });

        document.getElementById('resetBtn').addEventListener('click', () => {
            this.resetMarket();
        });

        // Modal
        const modal = document.getElementById('theoryModal');
        const closeBtn = document.querySelector('.close');
        
        closeBtn.addEventListener('click', () => {
            modal.style.display = 'none';
        });

        window.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.style.display = 'none';
            }
        });

        // Botón para abrir modal de teoría
        document.querySelector('.theory-quote').addEventListener('click', () => {
            modal.style.display = 'block';
        });
    }

    setupTabs() {
        const tabBtns = document.querySelectorAll('.tab-btn');
        const tabPanes = document.querySelectorAll('.tab-pane');

        tabBtns.forEach(btn => {
            btn.addEventListener('click', () => {
                const targetTab = btn.getAttribute('data-tab');
                
                tabBtns.forEach(b => b.classList.remove('active'));
                tabPanes.forEach(p => p.classList.remove('active'));
                
                btn.classList.add('active');
                document.getElementById(targetTab).classList.add('active');
                
                if (targetTab === 'workers') {
                    this.updateWorkersChart();
                }
            });
        });
    }

    calculateMarketExpansion() {
        const { roads, ports, markets, warehouses } = this.market.infrastructure;
        
        // Calcular nuevo radio basado en infraestructura
        let newRadius = 10; // Radio base
        newRadius += roads * 5; // Cada carretera añade 5km
        newRadius += ports * 8; // Cada puerto añade 8km
        newRadius += markets * 3; // Cada mercado añade 3km
        newRadius += warehouses * 2; // Cada almacén añade 2km
        
        this.market.radius = newRadius;
        
        // Calcular nueva población
        this.market.population = Math.floor(50 + (newRadius - 10) * 10);
        
        // Calcular nuevo nivel de especialización
        this.market.specializationLevel = Math.min(5, Math.max(1, Math.floor(newRadius / 20)));
        
        // Calcular nuevo valor del mercado
        this.market.marketValue = Math.floor(1000 + (newRadius - 10) * 100 + this.market.population * 5);
        
        this.updateDisplay();
        this.updateMarketVisualization();
        this.updateSpecializationLevels();
    }

    expandMarket() {
        const expansion = {
            timestamp: new Date(),
            radius: this.market.radius,
            population: this.market.population,
            specializationLevel: this.market.specializationLevel,
            infrastructure: { ...this.market.infrastructure }
        };
        
        this.market.history.push(expansion);
        this.updateTimeline();
        this.updateWorkersChart();
        
        // Mostrar notificación
        this.showNotification('¡Mercado expandido exitosamente!', 'success');
    }

    resetMarket() {
        this.market = {
            radius: 10,
            population: 50,
            specializationLevel: 1,
            marketValue: 1000,
            infrastructure: {
                roads: 0,
                ports: 0,
                markets: 0,
                warehouses: 0
            },
            history: []
        };
        
        // Resetear sliders
        const sliders = ['roads', 'ports', 'markets', 'warehouses'];
        sliders.forEach(type => {
            const slider = document.getElementById(`${type}Input`);
            const valueSpan = document.getElementById(`${type}Value`);
            slider.value = 0;
            valueSpan.textContent = '0';
        });
        
        this.updateDisplay();
        this.updateMarketVisualization();
        this.updateSpecializationLevels();
        this.updateTimeline();
        this.updateWorkersChart();
        
        this.showNotification('Mercado reiniciado', 'info');
    }

    updateDisplay() {
        document.getElementById('marketRadius').textContent = this.market.radius;
        document.getElementById('marketPopulation').textContent = this.market.population;
        document.getElementById('specializationLevel').textContent = this.market.specializationLevel;
        document.getElementById('marketValue').textContent = this.market.marketValue.toLocaleString();
        
        // Actualizar métricas de comercio
        document.getElementById('tradeRoutes').textContent = this.market.infrastructure.roads + this.market.infrastructure.ports;
        document.getElementById('tradeVolume').textContent = Math.floor(this.market.marketValue / 100);
        document.getElementById('marketReach').textContent = this.market.radius;
        
        // Actualizar métricas de producción
        const productivity = Math.min(100, 60 + this.market.specializationLevel * 10);
        const efficiency = Math.min(100, 50 + this.market.infrastructure.roads * 5);
        const quality = Math.min(100, 70 + this.market.infrastructure.warehouses * 5);
        
        document.getElementById('productivityValue').textContent = productivity + '%';
        document.getElementById('efficiencyValue').textContent = efficiency + '%';
        document.getElementById('qualityValue').textContent = quality + '%';
        
        document.getElementById('productivityBar').style.width = productivity + '%';
        document.getElementById('efficiencyBar').style.width = efficiency + '%';
        document.getElementById('qualityBar').style.width = quality + '%';
    }

    updateMarketVisualization() {
        // Actualizar anillos del mercado
        const rings = document.querySelectorAll('.ring');
        rings.forEach((ring, index) => {
            const level = index + 1;
            if (this.market.specializationLevel >= level) {
                ring.classList.add('active');
            } else {
                ring.classList.remove('active');
            }
        });
        
        // Actualizar marcadores de infraestructura
        this.updateInfrastructureMarkers();
    }

    updateInfrastructureMarkers() {
        const markersContainer = document.getElementById('infrastructureMarkers');
        markersContainer.innerHTML = '';
        
        const { roads, ports, markets, warehouses } = this.market.infrastructure;
        
        // Crear marcadores para carreteras
        for (let i = 0; i < roads; i++) {
            this.createMarker('road', 'fa-road', i * 30);
        }
        
        // Crear marcadores para puertos
        for (let i = 0; i < ports; i++) {
            this.createMarker('port', 'fa-ship', i * 60);
        }
        
        // Crear marcadores para mercados
        for (let i = 0; i < markets; i++) {
            this.createMarker('market', 'fa-store', i * 45);
        }
        
        // Crear marcadores para almacenes
        for (let i = 0; i < warehouses; i++) {
            this.createMarker('warehouse', 'fa-warehouse', i * 75);
        }
    }

    createMarker(type, icon, angle) {
        const marker = document.createElement('div');
        marker.className = `marker ${type}`;
        marker.innerHTML = `<i class="fas ${icon}"></i>`;
        
        const radius = 120;
        const radian = (angle * Math.PI) / 180;
        const x = Math.cos(radian) * radius + 150;
        const y = Math.sin(radian) * radius + 150;
        
        marker.style.left = x + 'px';
        marker.style.top = y + 'px';
        
        document.getElementById('infrastructureMarkers').appendChild(marker);
    }

    updateSpecializationLevels() {
        // Actualizar indicadores de nivel
        const indicators = document.querySelectorAll('.indicator');
        indicators.forEach((indicator, index) => {
            const level = index + 1;
            if (this.market.specializationLevel >= level) {
                indicator.classList.add('active');
            } else {
                indicator.classList.remove('active');
            }
        });
        
        // Actualizar tarjetas de nivel
        const levelCards = document.querySelectorAll('.level-card');
        levelCards.forEach((card, index) => {
            const level = index + 1;
            if (this.market.specializationLevel >= level) {
                card.classList.add('active');
            } else {
                card.classList.remove('active');
            }
        });
        
        // Actualizar listas de especializaciones
        this.updateSpecializationLists();
    }

    updateSpecializationLists() {
        const currentContainer = document.getElementById('currentSpecializations');
        const nextContainer = document.getElementById('nextSpecializations');
        
        // Especializaciones actuales
        const currentSpecs = this.specializations[this.market.specializationLevel] || [];
        currentContainer.innerHTML = currentSpecs.map(spec => 
            `<div class="specialization-item">${spec}</div>`
        ).join('');
        
        // Próximas especializaciones
        const nextSpecs = this.nextSpecializations[this.market.specializationLevel] || [];
        nextContainer.innerHTML = nextSpecs.map(spec => 
            `<div class="specialization-item">${spec}</div>`
        ).join('');
    }

    updateTimeline() {
        const timeline = document.getElementById('expansionTimeline');
        timeline.innerHTML = '';
        
        this.market.history.slice(-5).forEach(expansion => {
            const item = document.createElement('div');
            item.className = 'timeline-item';
            item.innerHTML = `
                <i class="fas fa-expand-arrows-alt"></i>
                <div class="timeline-content">
                    <div class="timeline-title">Expansión del Mercado</div>
                    <div class="timeline-description">
                        Radio: ${expansion.radius}km | Población: ${expansion.population} | 
                        Especialización: Nivel ${expansion.specializationLevel}
                    </div>
                </div>
                <div class="timeline-time">${expansion.timestamp.toLocaleTimeString()}</div>
            `;
            timeline.appendChild(item);
        });
    }

    initializeCharts() {
        this.updateWorkersChart();
    }

    updateWorkersChart() {
        const ctx = document.getElementById('workersChart');
        if (!ctx) return;
        
        // Destruir gráfico existente si hay uno
        if (this.workersChart) {
            this.workersChart.destroy();
        }
        
        const levels = [1, 2, 3, 4, 5];
        const data = levels.map(level => {
            const specs = this.specializations[level] || [];
            return specs.length;
        });
        
        this.workersChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: ['Aldea', 'Pueblo', 'Ciudad', 'Metrópolis', 'Imperio'],
                datasets: [{
                    label: 'Número de Especializaciones',
                    data: data,
                    borderColor: '#667eea',
                    backgroundColor: 'rgba(102, 126, 234, 0.1)',
                    borderWidth: 3,
                    fill: true,
                    tension: 0.4
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: false
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        title: {
                            display: true,
                            text: 'Especializaciones'
                        }
                    },
                    x: {
                        title: {
                            display: true,
                            text: 'Nivel de Mercado'
                        }
                    }
                }
            }
        });
    }

    showNotification(message, type = 'info') {
        // Crear notificación temporal
        const notification = document.createElement('div');
        notification.className = `notification notification-${type}`;
        notification.innerHTML = `
            <i class="fas fa-${type === 'success' ? 'check-circle' : 'info-circle'}"></i>
            <span>${message}</span>
        `;
        
        // Estilos para la notificación
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            background: ${type === 'success' ? '#28a745' : '#17a2b8'};
            color: white;
            padding: 15px 20px;
            border-radius: 8px;
            box-shadow: 0 4px 15px rgba(0,0,0,0.2);
            z-index: 1000;
            display: flex;
            align-items: center;
            gap: 10px;
            animation: slideIn 0.3s ease;
        `;
        
        document.body.appendChild(notification);
        
        // Remover después de 3 segundos
        setTimeout(() => {
            notification.style.animation = 'slideOut 0.3s ease';
            setTimeout(() => {
                document.body.removeChild(notification);
            }, 300);
        }, 3000);
    }
}

// Estilos adicionales para animaciones
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from { transform: translateX(100%); opacity: 0; }
        to { transform: translateX(0); opacity: 1; }
    }
    
    @keyframes slideOut {
        from { transform: translateX(0); opacity: 1; }
        to { transform: translateX(100%); opacity: 0; }
    }
    
    .theory-quote {
        cursor: pointer;
        transition: all 0.3s ease;
    }
    
    .theory-quote:hover {
        transform: scale(1.02);
        text-shadow: 0 2px 10px rgba(255,255,255,0.3);
    }
`;
document.head.appendChild(style);

// Inicializar el sistema cuando el DOM esté listo
document.addEventListener('DOMContentLoaded', () => {
    window.expansionSystem = new ExpansionSystem();
}); 