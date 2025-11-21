// Pagina de perfil, todo tranqui
document.addEventListener('DOMContentLoaded', async () => {
    // Revisamos si esta logueado o no
    const token = localStorage.getItem('accessToken');
    if (!token) {
        window.location.href = 'login.html';
        return;
    }

    // Cargamos datos del perfil
    try {
        await loadUserProfile();
    } catch (error) {
        console.error('Failed to load profile:', error);
        showNotification('Error al cargar el perfil', 'error');
    }

    // Preparamos los formularios y demas cositas
    setupProfileUpdateForm();
    setupPasswordChangeForm();
    setupDeleteAccountForm();
    setupOrdersDisplay();
});

// Traer datos del perfil desde la API
async function loadUserProfile() {
    const user = await apiClient.get('/api/users/profile');

    // Rellenamos los campos
    document.getElementById('first-name').value = user.first_name || '';
    document.getElementById('last-name').value = user.last_name || '';
    const phoneValue = user.phone ? user.phone.replace(/^\+?57/, '') : '';
    document.getElementById('phone').value = phoneValue;
    document.getElementById('email').value = user.email || '';

    // El email no se toca despues de registrarse
    document.getElementById('email').disabled = true;
}

// Formulario para actualizar perfil
function setupProfileUpdateForm() {
    const updateBtn = document.getElementById('update-profile-btn');
    
    updateBtn.addEventListener('click', async (e) => {
        e.preventDefault();
        
        const firstName = document.getElementById('first-name').value.trim();
        const lastName = document.getElementById('last-name').value.trim();
        const phoneInput = document.getElementById('phone').value.trim();
        const phone = phoneInput ? '+57' + phoneInput : '';
        
        // Validaciones basicas
        if (!firstName || !lastName) {
            showNotification('Por favor completa todos los campos requeridos', 'error');
            return;
        }
        
        if (phone && !validatePhone(phone)) {
            showNotification('Formato de teléfono inválido. Debe tener 10 dígitos después de +57', 'error');
            return;
        }
        
        // Ponemos estado de cargando
        updateBtn.disabled = true;
        updateBtn.textContent = 'Actualizando...';
        
        await apiClient.put('/api/users/profile', {
            first_name: firstName,
            last_name: lastName,
            phone: phone
        });

        showNotification('Perfil actualizado exitosamente', 'success');
    });
}

// Form para cambiar clave
function setupPasswordChangeForm() {
    const changePasswordBtn = document.getElementById('change-password-btn');
    
    changePasswordBtn.addEventListener('click', async (e) => {
        e.preventDefault();
        
        const currentPassword = document.getElementById('current-password').value;
        const newPassword = document.getElementById('new-password').value;
        const confirmPassword = document.getElementById('confirm-password').value;
        
        // Validamos cositas
        if (!currentPassword || !newPassword || !confirmPassword) {
            showNotification('Por favor completa todos los campos de contraseña', 'error');
            return;
        }
        
        if (newPassword !== confirmPassword) {
            showNotification('Las contraseñas no coinciden', 'error');
            return;
        }
        
        if (newPassword.length < 8) {
            showNotification('La nueva contraseña debe tener al menos 8 caracteres', 'error');
            return;
        }
        
        // Estado cargando
        changePasswordBtn.disabled = true;
        changePasswordBtn.textContent = 'Cambiando...';
        
        try {
            const response = await apiClient.put('/api/users/change-password', {
                current_password: currentPassword,
                new_password: newPassword
            });

            showNotification(response.message || 'Contraseña cambiada exitosamente', 'success');
        } catch (error) {
            // Mostramos el error que llegue
            const errorMessage = error.message || 'Error al cambiar la contraseña';
            showNotification(errorMessage, 'error');
            changePasswordBtn.disabled = false;
            changePasswordBtn.textContent = 'Cambiar contraseña';
            return; // No limpiamos si fallo
        }

        // Limpiamos campos de passwd
        document.getElementById('current-password').value = '';
        document.getElementById('new-password').value = '';
        document.getElementById('confirm-password').value = '';
    });
}

// Form para eliminar la cuenta sin dramas
function setupDeleteAccountForm() {
    const deleteBtn = document.getElementById('delete-account-btn');
    const modal = createConfirmModal();

    deleteBtn.addEventListener('click', async (e) => {
        e.preventDefault();

        const confirmed = await modal.confirm('Eliminar cuenta', '¿Estás seguro de que quieres eliminar tu cuenta? Esta acción no se puede deshacer.');
        if (!confirmed) return;

        deleteBtn.disabled = true;
        deleteBtn.textContent = 'Eliminando...';

        try {
            await apiClient.deleteUserProfile();

            showNotification('Cuenta eliminada exitosamente', 'success');

            // Logout y mandamos al home
            setTimeout(() => {
                apiClient.logout();
                window.location.href = '../index.html';
            }, 2000);
        } catch (error) {
            console.error('Error deleting account:', error);
            showNotification('Error al eliminar la cuenta', 'error');
            deleteBtn.disabled = false;
            deleteBtn.textContent = 'Eliminar Cuenta';
        }
    });
}

// Modalito simple pa confirmar cosas heavy
function createConfirmModal() {
    const overlay = document.createElement('div');
    overlay.className = 'confirm-overlay hidden';

    const dialog = document.createElement('div');
    dialog.className = 'confirm-dialog';

    const titleEl = document.createElement('h3');
    const messageEl = document.createElement('p');

    const actions = document.createElement('div');
    actions.className = 'confirm-actions';

    const cancelBtn = document.createElement('button');
    cancelBtn.type = 'button';
    cancelBtn.textContent = 'Cancelar';
    cancelBtn.className = 'confirm-btn cancel';

    const acceptBtn = document.createElement('button');
    acceptBtn.type = 'button';
    acceptBtn.textContent = 'Sí, eliminar';
    acceptBtn.className = 'confirm-btn danger';

    actions.appendChild(cancelBtn);
    actions.appendChild(acceptBtn);
    dialog.appendChild(titleEl);
    dialog.appendChild(messageEl);
    dialog.appendChild(actions);
    overlay.appendChild(dialog);
    document.body.appendChild(overlay);

    let resolver;

    const close = () => {
        overlay.classList.add('hidden');
    };

    cancelBtn.addEventListener('click', () => {
        close();
        resolver(false);
    });

    acceptBtn.addEventListener('click', () => {
        close();
        resolver(true);
    });

    return {
        confirm(title, message) {
            titleEl.textContent = title;
            messageEl.textContent = message;
            overlay.classList.remove('hidden');
            return new Promise((resolve) => {
                resolver = resolve;
            });
        }
    };
}

// Validar formato del telefono
function validatePhone(phone) {
    // Formato Colombiano: +57 y 10 digitos pegados
    const phoneRegex = /^\+57\d{10}$/;
    return phoneRegex.test(phone.replace(/\s/g, ''));
}

// Notificacion casera
function showNotification(message, type = 'info') {
    // Creamos el elemento
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;
    
    // Estilito rapido
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 15px 20px;
        background-color: ${type === 'success' ? '#4CAF50' : type === 'error' ? '#ff4081' : '#2196F3'};
        color: white;
        border-radius: 8px;
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        border: none;
        font-family: 'Montserrat', sans-serif;
        font-size: 14px;
        font-weight: 500;
        z-index: 10000;
        animation: slideIn 0.3s ease-out;
        max-width: 400px;
        word-wrap: break-word;
    `;
    
    // Lo pegamos al body
    document.body.appendChild(notification);
    
    // Se borra despues de 3 segs
    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => {
            document.body.removeChild(notification);
        }, 300);
    }, 3000);
}

// Mostrar pedidos en la vista
let currentPage = 1;
const ordersPerPage = 5;

async function setupOrdersDisplay() {
    await loadOrders(currentPage);

    // Paginacion basica
    document.getElementById('prev-page').addEventListener('click', () => {
        if (currentPage > 1) {
            currentPage--;
            loadOrders(currentPage);
        }
    });

    document.getElementById('next-page').addEventListener('click', () => {
        currentPage++;
        loadOrders(currentPage);
    });

    // Boton para completar pagos pendientes
    document.addEventListener('click', (e) => {
        if (e.target.classList.contains('complete-payment-btn')) {
            const orderId = e.target.getAttribute('data-order-id');
            window.location.href = `pedido.html?order_id=${orderId}`;
        }
    });
}

async function loadOrders(page) {
    const container = document.getElementById('orders-container');
    const pagination = document.getElementById('orders-pagination');

    try {
        const response = await apiClient.getOrders(page, ordersPerPage);

        if (response.orders && response.orders.length > 0) {
            displayOrders(response.orders);
            updatePagination(response.total, page);
            pagination.style.display = 'flex';
        } else {
            container.innerHTML = '<div class="no-orders">No tienes pedidos aún.</div>';
            pagination.style.display = 'none';
        }
    } catch (error) {
        console.error('Error loading orders:', error);
        container.innerHTML = '<div class="no-orders">Error al cargar los pedidos.</div>';
        pagination.style.display = 'none';
    }
}

function displayOrders(orders) {
    const container = document.getElementById('orders-container');
    container.innerHTML = '';

    orders.forEach(order => {
        const orderElement = createOrderElement(order);
        container.appendChild(orderElement);
    });
}

function createOrderElement(order) {
    const orderDiv = document.createElement('div');
    orderDiv.className = 'order-item';

    const statusClass = order.payment_status === 'completed' ? order.status : 'failed';

    orderDiv.innerHTML = `
        <div class="order-header">
            <div class="order-id">Pedido #${order.id.slice(0, 8)}</div>
            <div class="order-date">${new Date(order.created_at).toLocaleDateString('es-CO')}</div>
            <div class="order-status ${statusClass}">${getStatusText(order.status, order.payment_status)}</div>
            ${(order.payment_status === 'failed' || order.payment_status === 'pending') ? `<button class="btn-secondary complete-payment-btn" data-order-id="${order.id}">Completar Pago</button>` : ''}
        </div>
        <div class="order-items-list">
            ${order.items.map(item => `
                <div class="order-item-detail">
                    <div class="item-info">
                        <div class="item-name">${item.product_name}</div>
                        <div class="item-details">Cantidad: ${item.quantity}, Talla: ${item.size}</div>
                    </div>
                    <div class="item-price">$${formatPriceColombian(item.unit_price)}</div>
                </div>
            `).join('')}
        </div>
        <div class="order-total">Total: $${formatPriceColombian(order.total_amount)}</div>
    `;

    return orderDiv;
}

function getStatusText(status, paymentStatus) {
    if (paymentStatus !== 'completed') {
        return 'Pago Fallido';
    }

    switch (status) {
        case 'pending': return 'Pendiente';
        case 'confirmed': return 'Confirmado';
        case 'shipped': return 'Enviado';
        case 'delivered': return 'Entregado';
        default: return status;
    }
}

function updatePagination(total, currentPage) {
    const totalPages = Math.ceil(total / ordersPerPage);
    const pageInfo = document.getElementById('page-info');
    const prevBtn = document.getElementById('prev-page');
    const nextBtn = document.getElementById('next-page');

    pageInfo.textContent = `Página ${currentPage} de ${totalPages}`;

    prevBtn.disabled = currentPage <= 1;
    nextBtn.disabled = currentPage >= totalPages;
}

// Precio al estilo CO con punticos
function formatPriceColombian(price) {
    const intPrice = Math.floor(price);
    return intPrice.toLocaleString('es-CO').replace(/,/g, '.');
}

// Animaciones basicas pa las notis
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }

    @keyframes slideOut {
        from {
            transform: translateX(0);
            opacity: 1;
        }
        to {
            transform: translateX(100%);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);
