// Sistema de Trueque Inteligente - JavaScript

class TruequeSystem {
    constructor() {
        this.usuarios = [];
        this.intercambios = [];
        this.ofertas = [];
        this.estadisticas = {};
        this.usuarioSeleccionado = null;
        
        this.init();
    }

    async init() {
        await this.cargarUsuarios();
        this.setupEventListeners();
        this.setupTabs();
        this.cargarEstadisticasGenerales();
        this.cargarOfertas();
    }

    setupEventListeners() {
        // Selector de usuario
        const usuarioSelect = document.getElementById('usuarioSelect');
        usuarioSelect.addEventListener('change', (e) => {
            this.usuarioSeleccionado = e.target.value;
            if (this.usuarioSeleccionado) {
                this.cargarIntercambios(this.usuarioSeleccionado);
            }
        });

        // Calculadora
        const calcularBtn = document.getElementById('calcularBtn');
        calcularBtn.addEventListener('click', () => this.calcularValor());

        // Modal
        const modal = document.getElementById('intercambioModal');
        const closeBtn = document.querySelector('.close');
        
        closeBtn.addEventListener('click', () => {
            modal.style.display = 'none';
        });

        window.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.style.display = 'none';
            }
        });
    }

    setupTabs() {
        const tabBtns = document.querySelectorAll('.tab-btn');
        const tabPanes = document.querySelectorAll('.tab-pane');

        tabBtns.forEach(btn => {
            btn.addEventListener('click', () => {
                const targetTab = btn.getAttribute('data-tab');
                
                // Remover clase active de todos los botones y paneles
                tabBtns.forEach(b => b.classList.remove('active'));
                tabPanes.forEach(p => p.classList.remove('active'));
                
                // Agregar clase active al botón clickeado y su panel correspondiente
                btn.classList.add('active');
                document.getElementById(targetTab).classList.add('active');
                
                // Cargar datos específicos de la tab
                this.cargarDatosTab(targetTab);
            });
        });
    }

    async cargarDatosTab(tabName) {
        switch(tabName) {
            case 'intercambios':
                if (this.usuarioSeleccionado) {
                    await this.cargarIntercambios(this.usuarioSeleccionado);
                }
                break;
            case 'ofertas':
                await this.cargarOfertas();
                break;
            case 'estadisticas':
                await this.cargarEstadisticas();
                break;
        }
    }

    async cargarUsuarios() {
        try {
            const response = await fetch('/api/users');
            const data = await response.json();
            
            if (data.success) {
                this.usuarios = data.users;
                this.popularSelectorUsuarios();
            }
        } catch (error) {
            console.error('Error cargando usuarios:', error);
            this.mostrarError('Error al cargar usuarios');
        }
    }

    popularSelectorUsuarios() {
        const select = document.getElementById('usuarioSelect');
        select.innerHTML = '<option value="">Selecciona un usuario...</option>';
        
        this.usuarios.forEach(usuario => {
            const option = document.createElement('option');
            option.value = usuario.id;
            option.textContent = `${usuario.nombre} (${usuario.rol})`;
            select.appendChild(option);
        });
    }

    async cargarIntercambios(usuarioId) {
        try {
            const container = document.getElementById('intercambiosContainer');
            container.innerHTML = '<div class="loading">Cargando intercambios...</div>';

            const response = await fetch(`/api/trueque/intercambios/${usuarioId}`);
            const data = await response.json();
            
            if (data.success) {
                this.intercambios = data.intercambios;
                this.mostrarIntercambios();
            }
        } catch (error) {
            console.error('Error cargando intercambios:', error);
            this.mostrarError('Error al cargar intercambios');
        }
    }

    mostrarIntercambios() {
        const container = document.getElementById('intercambiosContainer');
        
        if (this.intercambios.length === 0) {
            container.innerHTML = `
                <div class="no-data">
                    <i class="fas fa-info-circle"></i>
                    <p>No se encontraron intercambios viables para este usuario</p>
                </div>
            `;
            return;
        }

        container.innerHTML = this.intercambios.map(intercambio => {
            const usuarioOrigen = this.usuarios.find(u => u.id === intercambio.usuario_origen);
            const usuarioDestino = this.usuarios.find(u => u.id === intercambio.usuario_destino);
            
            return `
                <div class="intercambio-card" onclick="truequeSystem.mostrarDetallesIntercambio('${intercambio.id}')">
                    <div class="intercambio-header">
                        <span class="intercambio-id">${intercambio.id}</span>
                        <span class="intercambio-estado">${intercambio.estado}</span>
                    </div>
                    <div class="intercambio-body">
                        <div class="usuario-info">
                            <div class="usuario-nombre">${usuarioOrigen?.nombre || 'Usuario'}</div>
                            <div class="usuario-rol">${usuarioOrigen?.rol || 'N/A'}</div>
                            <div class="producto-info">
                                <div class="producto-nombre">${intercambio.producto_origen}</div>
                                <div class="producto-cantidad">${intercambio.cantidad_origen}</div>
                            </div>
                        </div>
                        <div class="intercambio-arrow">
                            <i class="fas fa-exchange-alt"></i>
                        </div>
                        <div class="usuario-info">
                            <div class="usuario-nombre">${usuarioDestino?.nombre || 'Usuario'}</div>
                            <div class="usuario-rol">${usuarioDestino?.rol || 'N/A'}</div>
                            <div class="producto-info">
                                <div class="producto-nombre">${intercambio.producto_destino}</div>
                                <div class="producto-cantidad">${intercambio.cantidad_destino}</div>
                            </div>
                        </div>
                    </div>
                    <div class="intercambio-footer">
                        <div class="valor-info">
                            Valor: ${intercambio.valor_origen.toFixed(1)}h ↔ ${intercambio.valor_destino.toFixed(1)}h
                        </div>
                        <div class="fecha-info">${intercambio.fecha_creacion}</div>
                    </div>
                </div>
            `;
        }).join('');
    }

    async cargarOfertas() {
        try {
            const container = document.getElementById('ofertasContainer');
            container.innerHTML = '<div class="loading">Cargando ofertas...</div>';

            const response = await fetch('/api/trueque/ofertas');
            const data = await response.json();
            
            if (data.success) {
                this.ofertas = data.ofertas;
                this.mostrarOfertas();
            }
        } catch (error) {
            console.error('Error cargando ofertas:', error);
            this.mostrarError('Error al cargar ofertas');
        }
    }

    mostrarOfertas() {
        const container = document.getElementById('ofertasContainer');
        
        if (this.ofertas.length === 0) {
            container.innerHTML = `
                <div class="no-data">
                    <i class="fas fa-info-circle"></i>
                    <p>No hay ofertas de trueque activas</p>
                </div>
            `;
            return;
        }

        container.innerHTML = this.ofertas.map(oferta => {
            const usuario = this.usuarios.find(u => u.id === oferta.usuario_id);
            
            return `
                <div class="oferta-card">
                    <div class="oferta-header">
                        <span class="oferta-id">${oferta.id}</span>
                        <span class="oferta-activa">${oferta.activa ? 'Activa' : 'Inactiva'}</span>
                    </div>
                    <div class="oferta-body">
                        <div class="producto-info">
                            <div class="producto-nombre">${oferta.producto_ofrece}</div>
                            <div class="producto-cantidad">${oferta.cantidad_ofrece}</div>
                        </div>
                        <div class="oferta-arrow">
                            <i class="fas fa-arrow-right"></i>
                        </div>
                        <div class="producto-info">
                            <div class="producto-nombre">${oferta.producto_busca}</div>
                            <div class="producto-cantidad">${oferta.cantidad_busca}</div>
                        </div>
                    </div>
                    <div class="oferta-footer">
                        <span>Por: ${usuario?.nombre || 'Usuario'}</span>
                        <span>Valor: ${oferta.valor_ofrece.toFixed(1)}h</span>
                    </div>
                </div>
            `;
        }).join('');
    }

    async cargarEstadisticas() {
        try {
            const response = await fetch('/api/trueque/estadisticas');
            const data = await response.json();
            
            if (data.success) {
                this.estadisticas = data.estadisticas;
                this.mostrarEstadisticas();
            }
        } catch (error) {
            console.error('Error cargando estadísticas:', error);
            this.mostrarError('Error al cargar estadísticas');
        }
    }

    mostrarEstadisticas() {
        // Resumen general
        const statsResumen = document.getElementById('statsResumen');
        statsResumen.innerHTML = `
            <div class="stat-item">
                <div class="value">${this.estadisticas.total_ofertas}</div>
                <div class="label">Total Ofertas</div>
            </div>
            <div class="stat-item">
                <div class="value">${this.estadisticas.valor_total_ofrecido.toFixed(1)}h</div>
                <div class="label">Valor Ofrecido</div>
            </div>
            <div class="stat-item">
                <div class="value">${this.estadisticas.valor_total_buscado.toFixed(1)}h</div>
                <div class="label">Valor Buscado</div>
            </div>
            <div class="stat-item">
                <div class="value">${this.estadisticas.balance_mercado.toFixed(1)}h</div>
                <div class="label">Balance</div>
            </div>
        `;

        // Productos ofrecidos
        const productosOfrecidos = document.getElementById('productosOfrecidos');
        productosOfrecidos.innerHTML = Object.entries(this.estadisticas.productos_ofrecidos || {})
            .map(([producto, count]) => `
                <div class="product-bar">
                    <span class="product-name">${producto}</span>
                    <span class="product-count">${count}</span>
                </div>
            `).join('');

        // Productos buscados
        const productosBuscados = document.getElementById('productosBuscados');
        productosBuscados.innerHTML = Object.entries(this.estadisticas.productos_buscados || {})
            .map(([producto, count]) => `
                <div class="product-bar">
                    <span class="product-name">${producto}</span>
                    <span class="product-count">${count}</span>
                </div>
            `).join('');
    }

    async cargarEstadisticasGenerales() {
        try {
            const response = await fetch('/api/stats');
            const data = await response.json();
            
            if (data.success) {
                document.getElementById('totalUsuarios').textContent = data.stats.total_usuarios;
            }

            const ofertasResponse = await fetch('/api/trueque/ofertas');
            const ofertasData = await ofertasResponse.json();
            
            if (ofertasData.success) {
                document.getElementById('totalOfertas').textContent = ofertasData.count;
            }

            const statsResponse = await fetch('/api/trueque/estadisticas');
            const statsData = await statsResponse.json();
            
            if (statsData.success) {
                document.getElementById('balanceMercado').textContent = 
                    statsData.estadisticas.balance_mercado.toFixed(1) + 'h';
            }
        } catch (error) {
            console.error('Error cargando estadísticas generales:', error);
        }
    }

    async calcularValor() {
        const producto = document.getElementById('productoSelect').value;
        const cantidad = document.getElementById('cantidadInput').value;

        if (!producto || !cantidad) {
            this.mostrarError('Por favor completa todos los campos');
            return;
        }

        try {
            const response = await fetch(`/api/trueque/valor/${producto}/${cantidad}`);
            const data = await response.json();
            
            if (data.success) {
                this.mostrarResultadoCalculadora(data);
            }
        } catch (error) {
            console.error('Error calculando valor:', error);
            this.mostrarError('Error al calcular el valor');
        }
    }

    mostrarResultadoCalculadora(data) {
        const container = document.getElementById('resultadoCalculadora');
        container.innerHTML = `
            <div class="calculation-result">
                <div class="product-info">
                    <div class="producto-nombre">${data.producto}</div>
                    <div class="producto-cantidad">${data.cantidad} unidades</div>
                </div>
                <div class="value-info">${data.valor.toFixed(1)}</div>
                <div class="unit-info">${data.unidad}</div>
            </div>
        `;
    }

    mostrarDetallesIntercambio(intercambioId) {
        const intercambio = this.intercambios.find(i => i.id === intercambioId);
        if (!intercambio) return;

        const usuarioOrigen = this.usuarios.find(u => u.id === intercambio.usuario_origen);
        const usuarioDestino = this.usuarios.find(u => u.id === intercambio.usuario_destino);

        const modalContent = document.getElementById('modalContent');
        modalContent.innerHTML = `
            <div class="intercambio-detalles">
                <div class="detalle-seccion">
                    <h3>Información del Intercambio</h3>
                    <p><strong>ID:</strong> ${intercambio.id}</p>
                    <p><strong>Estado:</strong> ${intercambio.estado}</p>
                    <p><strong>Fecha:</strong> ${intercambio.fecha_creacion}</p>
                </div>
                
                <div class="detalle-seccion">
                    <h3>Usuario Origen</h3>
                    <p><strong>Nombre:</strong> ${usuarioOrigen?.nombre || 'N/A'}</p>
                    <p><strong>Rol:</strong> ${usuarioOrigen?.rol || 'N/A'}</p>
                    <p><strong>Especialización:</strong> ${usuarioOrigen?.especializacion || 'N/A'}</p>
                    <p><strong>Ofrece:</strong> ${intercambio.cantidad_origen} ${intercambio.producto_origen} (${intercambio.valor_origen.toFixed(1)}h)</p>
                </div>
                
                <div class="detalle-seccion">
                    <h3>Usuario Destino</h3>
                    <p><strong>Nombre:</strong> ${usuarioDestino?.nombre || 'N/A'}</p>
                    <p><strong>Rol:</strong> ${usuarioDestino?.rol || 'N/A'}</p>
                    <p><strong>Especialización:</strong> ${usuarioDestino?.especializacion || 'N/A'}</p>
                    <p><strong>Ofrece:</strong> ${intercambio.cantidad_destino} ${intercambio.producto_destino} (${intercambio.valor_destino.toFixed(1)}h)</p>
                </div>
                
                <div class="detalle-seccion">
                    <h3>Análisis del Intercambio</h3>
                    <p><strong>Equilibrio de Valor:</strong> ${Math.abs(intercambio.valor_origen - intercambio.valor_destino) < 0.1 ? 'Equilibrado' : 'Desequilibrado'}</p>
                    <p><strong>Diferencia:</strong> ${Math.abs(intercambio.valor_origen - intercambio.valor_destino).toFixed(1)} horas</p>
                    <p><strong>Basado en:</strong> Teoría de valor-trabajo de Adam Smith</p>
                </div>
            </div>
        `;

        document.getElementById('intercambioModal').style.display = 'block';
    }

    mostrarError(mensaje) {
        // Implementar sistema de notificaciones de error
        console.error(mensaje);
        alert(mensaje);
    }
}

// Inicializar el sistema cuando el DOM esté listo
document.addEventListener('DOMContentLoaded', () => {
    window.truequeSystem = new TruequeSystem();
});

// Estilos adicionales para el modal
const style = document.createElement('style');
style.textContent = `
    .intercambio-detalles {
        max-height: 70vh;
        overflow-y: auto;
    }
    
    .detalle-seccion {
        margin-bottom: 20px;
        padding: 15px;
        background: #f8f9fa;
        border-radius: 8px;
    }
    
    .detalle-seccion h3 {
        color: #667eea;
        margin-bottom: 10px;
        font-size: 16px;
    }
    
    .detalle-seccion p {
        margin: 5px 0;
        color: #333;
    }
    
    .no-data {
        text-align: center;
        padding: 40px;
        color: #666;
    }
    
    .no-data i {
        font-size: 48px;
        margin-bottom: 15px;
        display: block;
        opacity: 0.5;
    }
`;
document.head.appendChild(style); 