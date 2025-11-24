document.addEventListener('DOMContentLoaded', function () {
    const selectedProduct = JSON.parse(localStorage.getItem('selectedProduct'));
    console.log('cargamos selectedProduct del localStorage:', selectedProduct); // logcito sin stress
    if (!selectedProduct) {
        console.error('No product selected');
        return;
    }

    const urlParams = new URLSearchParams(window.location.search);
    const store = urlParams.get('store') || 'mewave';
    document.body.classList.add(store === 'celeste' ? 'celeste-product' : 'mewave-product');
    const storeName = store === 'celeste' ? 'Celeste' : 'Mewave';

    // Pedimos productos con el api pa llenar la sidebar
    apiClient.getAllProducts(storeName)
        .then(response => {
            const products = response.products || [];
            populateProductDetails(selectedProduct);
            populateSidebar(products, selectedProduct.category, selectedProduct.name, store);

            const minusBtn = document.querySelector('.quantity-btn:first-child');
            const plusBtn = document.querySelector('.quantity-btn:last-child');
            const quantityInput = document.querySelector('.quantity-input');
            const goToCartLink = document.querySelector('.go-to-cart');

            if (goToCartLink) {
                goToCartLink.classList.add('hidden');
                goToCartLink.setAttribute('aria-hidden', 'true');
            }

            minusBtn.addEventListener('click', () => {
                let qty = parseInt(quantityInput.textContent);
                if (qty > 1) {
                    quantityInput.textContent = qty - 1;
                }
            });

            plusBtn.addEventListener('click', () => {
                let qty = parseInt(quantityInput.textContent);
                if (qty < 10) {
                    quantityInput.textContent = qty + 1;
                }
            });

            // Funcion de agregar al carrito, calmada
            const addToCartBtn = document.querySelector('.add-to-cart');
            addToCartBtn.addEventListener('click', async () => {
                // Revisamos si esta logueado
                console.log('Checking auth state...');
                console.log('apiClient exists:', !!window.apiClient);
                if (window.apiClient) {
                    console.log('Is authenticated:', window.apiClient.isAuthenticated());
                    console.log('Token:', localStorage.getItem('accessToken'));
                }

                if (!apiClient.isAuthenticated()) {
                    console.log('User not authenticated, redirecting to login');
                    showNotification('Debes iniciar sesión para agregar productos al carrito', 'error');
                    window.location.href = 'login.html';
                    return;
                }

                const qty = parseInt(quantityInput.textContent);
                const selectedSize = document.querySelector('.size-option.active')?.textContent || '';

                try {
                    // Deshabilitamos boton mientras agrega
                    addToCartBtn.disabled = true;
                    addToCartBtn.innerHTML = '<i class="fa-solid fa-spinner fa-spin"></i>';

                    // Mandamos al carrito via API
                    console.log('Adding to cart:', {
                        product_id: selectedProduct.id,
                        quantity: qty,
                        size: selectedSize
                    });
                    await apiClient.addItemToCart(selectedProduct.id, qty, selectedSize);

                    // Mostramos mensaje de exito
                    addToCartBtn.innerHTML = '<i class="fa-solid fa-check"></i>';
                    setTimeout(() => {
                        addToCartBtn.innerHTML = '<i class="fa-solid fa-cart-plus"></i>';
                        addToCartBtn.disabled = false;
                    }, 1500);

                    showNotification('¡Producto añadido al carrito!', 'success');

                    if (goToCartLink) {
                        goToCartLink.classList.remove('hidden');
                        goToCartLink.removeAttribute('aria-hidden');
                    }

                    if (typeof window.refreshCartCount === 'function') {
                        window.refreshCartCount();
                    }
                } catch (error) {
                    console.error('Error adding to cart:', error);
                    showNotification('Error al agregar el producto al carrito. Por favor intenta de nuevo.', 'error');
                    addToCartBtn.innerHTML = '<i class="fa-solid fa-cart-plus"></i>';
                    addToCartBtn.disabled = false;
                }
            });
        })
        .catch(error => console.error('Error loading products:', error));
});

function populateProductDetails(product) {
    document.title = product.name + ' - Detail Page';

    const backBtn = document.querySelector('.back-nav a');
    const store = new URLSearchParams(window.location.search).get('store') || 'mewave';
    backBtn.href = store + '.html';

    const backImg = backBtn.querySelector('img');
    backImg.src = store === 'celeste' ? '../assets/icons/BackButtonCeleste.svg' : '../assets/icons/BackButtonMewave.svg';

    // Imagen grande del producto
    const productImage = document.querySelector('.product-image');
    const imageUrl = product.image_url ? ('../assets/images/' + product.image_url) : product.image;
    productImage.innerHTML = `<img src="${imageUrl}" alt="${product.name}" class="main-image">`;

    // Miniaturas
    const thumbnails = document.querySelector('.thumbnails');
    thumbnails.innerHTML = '';
    const thumbnail = document.createElement('div');
    thumbnail.className = 'thumbnail active';
    const thumbImg = document.createElement('img');
    thumbImg.src = imageUrl;
    thumbImg.alt = product.name;
    thumbnail.appendChild(thumbImg);
    thumbnails.appendChild(thumbnail);

    // Info del producto
    document.querySelector('h1').textContent = product.name;
    document.querySelector('.product-description').textContent = product.description;
    document.querySelector('.price').textContent = formatPrice(product.price);

    // Rating (las 5 estrellitas de siempre)
    const ratingDiv = document.querySelector('.rating');
    ratingDiv.innerHTML = '';
    for (let i = 0; i < 5; i++) {
        const star = document.createElement('i');
        star.className = 'fa-solid fa-star';
        ratingDiv.appendChild(star);
    }

    // Selector de color, lo escondemos por mientras
    const colorSelector = document.querySelector('.color-selector');
    colorSelector.style.display = 'none';

    // Selector de talla
    const sizesContainer = document.querySelector('.sizes');
    sizesContainer.innerHTML = '';
    product.sizes.forEach((size, index) => {
        const sizeOption = document.createElement('div');
        sizeOption.className = 'size-option' + (index === 0 ? ' active' : '');
        sizeOption.textContent = size;
        sizeOption.addEventListener('click', () => {
            document.querySelectorAll('.size-option').forEach(opt => opt.classList.remove('active'));
            sizeOption.classList.add('active');
        });
        sizesContainer.appendChild(sizeOption);
    });
}

function populateSidebar(products, category, currentProductName, store) {
    const title = document.querySelector('.sidebar-title');
    title.textContent = category.toUpperCase();
    title.className = 'sidebar-title' + (store === 'celeste' ? ' celeste-title' : '');

    const relatedProducts = products.filter(p => p.category === category && p.name !== currentProductName).slice(0, 3);

    const sidebar = document.querySelector('.sidebar');
    // Dejamos solo el titulo y limpiamos lo demas
    const titleElement = sidebar.querySelector('.sidebar-title');
    sidebar.innerHTML = '';
    sidebar.appendChild(titleElement);

    relatedProducts.forEach(product => {
        const cardContainer = document.createElement('div');
        cardContainer.className = 'sidebar-product-card';

        const img = document.createElement('img');
        const imageUrl = product.image_url ? ('../assets/images/' + product.image_url) : product.image;
        img.src = imageUrl;
        img.alt = product.name;
        img.className = 'sidebar-product-image';

        const info = document.createElement('div');
        info.className = 'sidebar-product-info';

        const name = document.createElement('div');
        name.className = 'sidebar-product-name';
        name.textContent = product.name;

        const price = document.createElement('div');
        price.className = 'sidebar-product-price';
        price.textContent = formatPrice(product.price);

        info.appendChild(name);
        info.appendChild(price);

        cardContainer.appendChild(img);
        cardContainer.appendChild(info);

        cardContainer.style.cursor = 'pointer';
        cardContainer.addEventListener('click', () => {
            localStorage.setItem('selectedProduct', JSON.stringify(product));
            window.location.href = 'product.html?store=' + store;
        });

        sidebar.appendChild(cardContainer);
    });
}

function formatPrice(price) {
    return '$' + (price / 1000).toFixed(0) + '.000';
}

// Notis sencillas para avisar cositas
function showNotification(message, type = 'info') {
    // Creamos el div de la noti
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;

    // Estilos al vuelo, cero stress
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

    // Pegamos al body
    document.body.appendChild(notification);

    // Se va solita en 3 segs
    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => {
            document.body.removeChild(notification);
        }, 300);
    }, 3000);
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
