/**
 * Cart Management with Backend Integration
 * Handles cart operations using the backend API
 */

let cartData = null;

document.addEventListener('DOMContentLoaded', async function () {
    // Check if user is authenticated
    if (!apiClient.isAuthenticated()) {
        showLoginMessage();
        return;
    }

    await loadCart();
});

/**
 * Show message for non-authenticated users
 */
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
}

/**
 * Load cart from backend
 */
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

        // Render each cart item
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

/**
 * Show empty cart message
 */
function showEmptyCart() {
    const cartItemsContainer = document.getElementById('cart-items');
    cartItemsContainer.innerHTML = `
        <div style="text-align: center; padding: 40px;">
            <p style="font-size: 18px; margin-bottom: 20px;">Tu carrito está vacío</p>
            <a href="../index.html" style="background: #FF69B4; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px;">Ir a las Tiendas</a>
        </div>
    `;

    updateOrderSummary();
}

/**
 * Create cart item HTML element
 */
function createCartItemElement(item) {
    const cartItem = document.createElement('div');
    cartItem.className = 'cart-item';
    cartItem.dataset.itemId = item.id;

    // Get product image (use API image_url or placeholder)
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

/**
 * Add event listeners to cart items
 */
function addEventListeners() {
    // Quantity buttons
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

    // Delete buttons
    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            const itemId = e.currentTarget.dataset.itemId;
            await removeFromCart(itemId);
        });
    });
}

/**
 * Update cart item quantity
 */
async function updateCartItemQuantity(itemId, newQuantity) {
    try {
        // Find the item to get its size
        const item = cartData.items.find(i => i.id === itemId);
        if (!item) return;

        await apiClient.updateCartItem(itemId, newQuantity, item.size || '');

        // Update UI immediately for better UX
        const quantityInput = document.querySelector(`[data-item-id="${itemId}"]`).parentElement.querySelector('.quantity-input');
        quantityInput.textContent = newQuantity;

        // Reload cart to get updated totals
        await loadCart();
    } catch (error) {
        console.error('Error updating cart item:', error);
        showError('Error al actualizar la cantidad. Por favor intenta de nuevo.');
    }
}

/**
 * Remove item from cart
 */
async function removeFromCart(itemId) {
    try {
        await apiClient.removeItemFromCart(itemId);
        await loadCart();
    } catch (error) {
        console.error('Error removing item from cart:', error);
        showError('Error al eliminar el producto. Por favor intenta de nuevo.');
    }
}

/**
 * Update order summary
 */
function updateOrderSummary() {
    const orderSummary = document.querySelector('.order-summary');
    if (!orderSummary) return;

    // Clear existing items (keep only the total)
    const existingItems = orderSummary.querySelectorAll('.summary-item:not(.summary-total)');
    existingItems.forEach(item => item.remove());

    const hr = orderSummary.querySelector('hr');
    const summaryTotal = orderSummary.querySelector('.summary-total');

    let total = 0;

    if (cartData && cartData.items && cartData.items.length > 0) {
        // Add each item to summary
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

    // Update total
    if (summaryTotal) {
        summaryTotal.querySelector('span:last-child').textContent = formatPrice(total);
    }
}

/**
 * Format price to Colombian Pesos
 */
function formatPrice(price) {
    return '$' + price.toLocaleString('es-CO');
}

/**
 * Show error message
 */
function showError(message) {
    const cartItemsContainer = document.getElementById('cart-items');
    const errorDiv = document.createElement('div');
    errorDiv.style.cssText = 'background: #ffebee; color: #c62828; padding: 12px; border-radius: 4px; margin: 10px 0;';
    errorDiv.textContent = message;
    cartItemsContainer.insertBefore(errorDiv, cartItemsContainer.firstChild);

    // Remove error after 5 seconds
    setTimeout(() => errorDiv.remove(), 5000);
}