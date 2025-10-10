document.addEventListener('DOMContentLoaded', function() {
    loadCart();
});

function loadCart() {
    const cart = JSON.parse(localStorage.getItem('cart')) || [];
    const cartItemsContainer = document.getElementById('cart-items');

    cartItemsContainer.innerHTML = '';

    cart.forEach((item, index) => {
        const cartItem = document.createElement('div');
        cartItem.className = 'cart-item';
        cartItem.innerHTML = `
            <img src="${item.product.image}" alt="Product Image" class="item-image">
            <div class="item-details">
                <h3>${item.product.name}</h3>
                <p>${formatPrice(item.product.price)}</p>
                <div class="quantity-selector">
                    <button class="quantity-btn minus-btn">-</button>
                    <div class="quantity-input">${item.quantity}</div>
                    <button class="quantity-btn plus-btn">+</button>
                </div>
            </div>
            <button class="delete-btn" data-index="${index}"><i class="fa-solid fa-trash"></i></button>
        `;
        cartItemsContainer.appendChild(cartItem);
    });

    addEventListeners();
    updateOrderSummary();
}

function addEventListeners() {
    document.querySelectorAll('.quantity-selector').forEach(selector => {
        const minusBtn = selector.querySelector('.minus-btn');
        const plusBtn = selector.querySelector('.plus-btn');
        const quantityInput = selector.querySelector('.quantity-input');

        minusBtn.addEventListener('click', () => {
            let qty = parseInt(quantityInput.textContent);
            if (qty > 1) {
                quantityInput.textContent = qty - 1;
                updateCartQuantity(selector, qty - 1);
                updateOrderSummary();
            }
        });

        plusBtn.addEventListener('click', () => {
            let qty = parseInt(quantityInput.textContent);
            if (qty < 10) {
                quantityInput.textContent = qty + 1;
                updateCartQuantity(selector, qty + 1);
                updateOrderSummary();
            }
        });
    });

    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const index = parseInt(e.currentTarget.dataset.index);
            removeFromCart(index);
        });
    });
}

function updateCartQuantity(selector, newQty) {
    const cartItem = selector.closest('.cart-item');
    const index = Array.from(cartItem.parentNode.children).indexOf(cartItem);
    let cart = JSON.parse(localStorage.getItem('cart')) || [];
    cart[index].quantity = newQty;
    localStorage.setItem('cart', JSON.stringify(cart));
}

function removeFromCart(index) {
    let cart = JSON.parse(localStorage.getItem('cart')) || [];
    cart.splice(index, 1);
    localStorage.setItem('cart', JSON.stringify(cart));
    loadCart(); // Reload to update indices
}

function updateOrderSummary() {
    const cart = JSON.parse(localStorage.getItem('cart')) || [];
    let total = 0;

    // Clear existing summary items except total
    const summaryItems = document.querySelectorAll('.summary-item');
    summaryItems.forEach(item => {
        if (!item.classList.contains('summary-total')) {
            item.remove();
        }
    });

    const orderSummary = document.querySelector('.order-summary');
    const hr = orderSummary.querySelector('hr');
    const summaryTotal = document.querySelector('.summary-total');

    cart.forEach(item => {
        const name = item.product.name;
        const price = item.product.price;
        const qty = item.quantity;
        const itemTotal = price * qty;

        const summaryItem = document.createElement('div');
        summaryItem.className = 'summary-item';
        summaryItem.innerHTML = `
            <span>${name} x${qty}</span>
            <span>${formatPrice(itemTotal)}</span>
        `;
        orderSummary.insertBefore(summaryItem, hr);

        total += itemTotal;
    });

    summaryTotal.querySelector('span:last-child').textContent = formatPrice(total);
}

function formatPrice(price) {
    return '$' + (price / 1000).toFixed(0) + '.000';
}