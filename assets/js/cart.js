// Manejo del carrito pegado al backend, todo relax pa no enredarnos

let cartData = null;

document.addEventListener('DOMContentLoaded', async function () {
    // Dejamos el boton de pago apagado al inicio, por si las moscas
    setProceedButtonState(false);

    // Miramos si el user si esta logueado
    if (!apiClient.isAuthenticated()) {
        showLoginMessage();
        return;
    }

    await loadCart();
});

// Mensaje cute para los que no han iniciado sesion
function showLoginMessage() {
    const cartItemsContainer = document.getElementById('cart-items');
    cartItemsContainer.innerHTML = `
        <div style="text-align: center; padding: 40px;">
            <p style="font-size: 18px; margin-bottom: 20px;">Debes iniciar sesión para ver tu carrito</p>
            <a href="login.html" style="background: #FF69B4; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px;">Iniciar Sesión</a>
        </div>
    `;

    const orderSummary = document.querySelector('.order-summary');
    if (orderSummary) {
        orderSummary.innerHTML = `
            <div class="summary-total">
                <span>Total</span>
                <span>$0</span>
            </div>
        `;
    }

    setProceedButtonState(false);
}

// Cargamos carrito desde el backend, sin tanto rollo
async function loadCart() {
    try {
        const response = await apiClient.getCart();
        cartData = response;

        const cartItemsContainer = document.getElementById('cart-items');
        cartItemsContainer.innerHTML = '';

        if (!cartData.items || cartData.items.length === 0) {
            showEmptyCart();
            return;
        }

        // Pintamos cada item del carrito
        for (const item of cartData.items) {
            const cartItem = createCartItemElement(item);
            cartItemsContainer.appendChild(cartItem);
        }

        addEventListeners();
        updateOrderSummary();
    } catch (error) {
        console.error('Error loading cart:', error);
        showError('Error al cargar el carrito. Por favor intenta de nuevo.');
    }
}

// Mensaje de carrito vacio pa que vuelvan a la tienda
function showEmptyCart() {
    const cartItemsContainer = document.getElementById('cart-items');
    cartItemsContainer.innerHTML = `
        <div style="text-align: center; padding: 40px;">
            <p style="font-size: 18px; margin-bottom: 20px;">Tu carrito está vacío</p>
            <a href="../index.html" style="background: #E73873; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px;">Ir a las Tiendas</a>
        </div>
    `;

    updateOrderSummary();
}

// Armamos el HTML de cada item del carrito
function createCartItemElement(item) {
    const cartItem = document.createElement('div');
    cartItem.className = 'cart-item';
    cartItem.dataset.itemId = item.id;

    // Imagen del producto, si falla le mandamos placeholder
    const imageUrl = item.image_url
        ? (apiClient.baseURL + item.image_url)
        : '../assets/images/placeholder.png';

    cartItem.innerHTML = `
        <img src="${imageUrl}" alt="${item.product_name}" class="item-image" onerror="this.src='../assets/images/placeholder.png'">
        <div class="item-details">
            <h3>${item.product_name}</h3>
            <p class="item-price">${formatPrice(item.price)}</p>
            ${item.size ? `<p class="item-size">Talla: ${item.size}</p>` : ''}
            <div class="quantity-selector">
                <button class="quantity-btn minus-btn" data-item-id="${item.id}">-</button>
                <div class="quantity-input">${item.quantity}</div>
                <button class="quantity-btn plus-btn" data-item-id="${item.id}">+</button>
            </div>
        </div>
        <button class="delete-btn" data-item-id="${item.id}"><i class="fa-solid fa-trash"></i></button>
    `;

    return cartItem;
}

// Eventos de los items (sumar, restar, borrar)
function addEventListeners() {
    // Botones de cantidad
    document.querySelectorAll('.minus-btn').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            const itemId = e.currentTarget.dataset.itemId;
            const quantityInput = e.currentTarget.parentElement.querySelector('.quantity-input');
            const currentQty = parseInt(quantityInput.textContent);

            if (currentQty > 1) {
                await updateCartItemQuantity(itemId, currentQty - 1);
            }
        });
    });

    document.querySelectorAll('.plus-btn').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            const itemId = e.currentTarget.dataset.itemId;
            const quantityInput = e.currentTarget.parentElement.querySelector('.quantity-input');
            const currentQty = parseInt(quantityInput.textContent);

            if (currentQty < 10) {
                await updateCartItemQuantity(itemId, currentQty + 1);
            }
        });
    });

    // Botones de borrar
    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            const itemId = e.currentTarget.dataset.itemId;
            await removeFromCart(itemId);
        });
    });
}

// Actualizar cantidad en el carrito rapidito
async function updateCartItemQuantity(itemId, newQuantity) {
    try {
        // Buscamos el item pa sacar la talla
        const item = cartData.items.find(i => i.id === itemId);
        if (!item) return;

        await apiClient.updateCartItem(itemId, newQuantity, item.size || '');

        // Actualizamos el UI de una pa que se vea fluido
        const quantityInput = document.querySelector(`[data-item-id="${itemId}"]`).parentElement.querySelector('.quantity-input');
        quantityInput.textContent = newQuantity;

        // Recargamos todo para recalcular totales
        await loadCart();
    } catch (error) {
        console.error('Error updating cart item:', error);
        showError('Error al actualizar la cantidad. Por favor intenta de nuevo.');
    }
}

// Quitar item del carrito
async function removeFromCart(itemId) {
    try {
        await apiClient.removeItemFromCart(itemId);
        await loadCart();
    } catch (error) {
        console.error('Error removing item from cart:', error);
        showError('Error al eliminar el producto. Por favor intenta de nuevo.');
    }
}

// Actualizar el resumen del pedido
function updateOrderSummary() {
    const orderSummary = document.querySelector('.order-summary');
    if (!orderSummary) return;

    // Limpiamos items viejos, dejamos el total
    const existingItems = orderSummary.querySelectorAll('.summary-item:not(.summary-total)');
    existingItems.forEach(item => item.remove());

    const hr = orderSummary.querySelector('hr');
    const summaryTotal = orderSummary.querySelector('.summary-total');

    let total = 0;
    const hasItems = cartData && cartData.items && cartData.items.length > 0;

    if (hasItems) {
        // Metemos cada item al resumen
        cartData.items.forEach(item => {
            const itemTotal = item.price * item.quantity;
            total += itemTotal;

            const summaryItem = document.createElement('div');
            summaryItem.className = 'summary-item';
            summaryItem.innerHTML = `
                <span>${item.product_name} x${item.quantity}</span>
                <span>${formatPrice(itemTotal)}</span>
            `;
            orderSummary.insertBefore(summaryItem, hr);
        });
    }

    // Total final
    if (summaryTotal) {
        summaryTotal.querySelector('span:last-child').textContent = formatPrice(total);
    }

    setProceedButtonState(hasItems);
}

// Formato de precio COP con puntos
function formatPrice(price) {
    return '$' + price.toLocaleString('es-CO');
}

// Wrapper para errors
function showError(message) {
    showNotification(message, 'error');
}

// Notificacion casera para mostrar feedback
function showNotification(message, type = 'info') {
    // Creamos el div
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;

    // Estilos al vuelo
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

    // Se borra solito a los tres segs
    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => {
            document.body.removeChild(notification);
        }, 300);
    }, 3000);
}

// Habilitar o bloquear el boton de pagar segun haya cosas
function setProceedButtonState(enabled) {
    const proceedBtn = document.querySelector('.proceed-btn');
    if (!proceedBtn) return;

    if (enabled) {
        proceedBtn.classList.remove('disabled');
        proceedBtn.removeAttribute('aria-disabled');
        proceedBtn.removeAttribute('tabindex');
        proceedBtn.removeEventListener('click', preventProceedWhenDisabled);
    } else {
        proceedBtn.classList.add('disabled');
        proceedBtn.setAttribute('aria-disabled', 'true');
        proceedBtn.setAttribute('tabindex', '-1');
        proceedBtn.addEventListener('click', preventProceedWhenDisabled);
    }
}

function preventProceedWhenDisabled(event) {
    // No dejamos que navegue si esta deshabilitado
    event.preventDefault();
}

// Animaciones sencillas pa las notis
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
