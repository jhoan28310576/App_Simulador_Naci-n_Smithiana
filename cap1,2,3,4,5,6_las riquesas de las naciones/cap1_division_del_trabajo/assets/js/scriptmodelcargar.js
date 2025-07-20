
// Función para calcular la productividad según la especialización
function calculateProductivity(user) {
    const baseProductivity = user.productividad;
    const hasSpecialization = user.especializacion !== null;
    const multiplier = hasSpecialization ? 10 : 1;
    return baseProductivity * multiplier;
  }
  
  // Función para mostrar el modal con la teoría
  function showTheory(user) {
    console.log('Mostrando teoría para:', user) // Para debugging
    const modal = document.getElementById('theoryModal')
    const content = document.getElementById('theoryContent')
    const hasSpecialization = user.especializacion !== null
    const productivity = calculateProductivity(user)
  
    content.innerHTML = `
            <div class="productivity-info">
              <h4>Análisis de Productividad</h4>
              <p><strong>Usuario:</strong> ${user.nombre}</p>
              <p><strong>Rol:</strong> ${user.rol}</p>
              <p><strong>Especialización:</strong> ${user.especializacion || 'Sin especialización'}</p>
              <p><strong>Productividad Base:</strong> ${user.productividad}</p>
              <p><strong>Multiplicador por Especialización:</strong> ${hasSpecialization ? '10x' : '1x'}</p>
              <p><strong>Productividad Total:</strong> ${productivity}</p>
              <p><strong>Teoría de Adam Smith:</strong></p>
              <p>${hasSpecialization
        ? 'Este trabajador está especializado, lo que multiplica su productividad por 10 según la teoría de la división del trabajo de Adam Smith.'
        : 'Este trabajador no está especializado, por lo que su productividad no se multiplica según la teoría de la división del trabajo.'}</p>
            </div>
          `
  
    modal.style.display = 'block'
  }
  
  // Función para cerrar el modal
  function closeModal() {
    const modal = document.getElementById('theoryModal');
    modal.style.display = 'none';
  }
  
  // Función para cargar los usuarios
  async function loadUsers() {
    const loadingDiv = document.getElementById('loading')
    const errorDiv = document.getElementById('error')
    const tableBody = document.getElementById('usersTableBody')
  
    try {
      const response = await fetch('/api/users')
      const data = await response.json()
  
      if (data.success) {
        loadingDiv.style.display = 'none'
        tableBody.innerHTML = ''
  
        data.users.forEach((user) => {
          const row = document.createElement('tr')
          // Crear el botón de teoría
          const theoryButton = document.createElement('button')
          theoryButton.className = 'theory-btn'
          theoryButton.textContent = 'Teoría'
          theoryButton.onclick = () => showTheory(user)
  
          row.innerHTML = `
                  <td>${user.id}</td>
                  <td>${user.nombre}</td>
                  <td>${user.rol}</td>
                  <td>${user.especializacion || 'Sin especialización'}</td>
                  <td>${user.productividad}</td>
                  <td></td>
                `
          tableBody.appendChild(row)
          // Agregar el botón a la última celda
          row.lastElementChild.appendChild(theoryButton)
        })
      } else {
        throw new Error(data.error || 'Error al cargar los usuarios')
      }
    } catch (error) {
      loadingDiv.style.display = 'none'
      errorDiv.style.display = 'block'
      errorDiv.textContent = error.message
    }
  }
  
  // Cerrar el modal cuando se hace clic en la X
  document.querySelector('.close-modal').addEventListener('click', closeModal);
  
  // Cerrar el modal cuando se hace clic fuera de él
  window.addEventListener('click', (event) => {
    const modal = document.getElementById('theoryModal');
    if (event.target === modal) {
      closeModal();
    }
  });
  
  // Cargar usuarios cuando la página se cargue
  document.addEventListener('DOMContentLoaded', loadUsers)