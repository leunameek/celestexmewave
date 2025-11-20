document.addEventListener('DOMContentLoaded', function () {
    const selectedProduct = JSON.parse(localStorage.getItem('selectedProduct'));
    console.log('Loaded selectedProduct from localStorage:', selectedProduct); // Debug log
    if (!selectedProduct) {
        console.error('No product selected');
        return;
    }

    const urlParams = new URLSearchParams(window.location.search);
    const store = urlParams.get('store') || 'mewave';
    document.body.classList.add(store === 'celeste' ? 'celeste-product' : 'mewave-product');
    const storeName = store === 'celeste' ? 'Celeste' : 'Mewave';

    // Use API Client to fetch products for sidebar
    apiClient.getAllProducts(storeName)
        .then(response => {
            const products = response.products || [];
            populateProductDetails(selectedProduct);
            populateSidebar(products, selectedProduct.category, selectedProduct.name, store);

            const minusBtn = document.querySelector('.quantity-btn:first-child');
            const plusBtn = document.querySelector('.quantity-btn:last-child');
            const quantityInput = document.querySelector('.quantity-input');

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

            // Add to cart functionality
            const addToCartBtn = document.querySelector('.add-to-cart');
            addToCartBtn.addEventListener('click', async () => {
                // Check if user is authenticated
                console.log('Checking auth state...');
                console.log('apiClient exists:', !!window.apiClient);
                if (window.apiClient) {
                    console.log('Is authenticated:', window.apiClient.isAuthenticated());
                    console.log('Token:', localStorage.getItem('accessToken'));
                }

                if (!apiClient.isAuthenticated()) {
                    console.log('User not authenticated, redirecting to login');
                    alert('Debes iniciar sesión para agregar productos al carrito');
                    window.location.href = 'login.html';
                    return;
                }

                const qty = parseInt(quantityInput.textContent);
                const selectedSize = document.querySelector('.size-option.active')?.textContent || '';

                try {
                    // Disable button while adding
                    addToCartBtn.disabled = true;
                    addToCartBtn.innerHTML = '<i class="fa-solid fa-spinner fa-spin"></i>';

                    // Add item to cart via API
                    console.log('Adding to cart:', {
                        product_id: selectedProduct.id,
                        quantity: qty,
                        size: selectedSize
                    });
                    await apiClient.addItemToCart(selectedProduct.id, qty, selectedSize);

                    // Show success message
                    addToCartBtn.innerHTML = '<i class="fa-solid fa-check"></i>';
                    setTimeout(() => {
                        addToCartBtn.innerHTML = '<i class="fa-solid fa-cart-plus"></i>';
                        addToCartBtn.disabled = false;
                    }, 1500);

                    alert('¡Producto añadido al carrito!');
                } catch (error) {
                    console.error('Error adding to cart:', error);
                    alert('Error al agregar el producto al carrito. Por favor intenta de nuevo.');
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

    // Main product image
    const productImage = document.querySelector('.product-image');
    const imageUrl = product.image_url ? (apiClient.baseURL + product.image_url) : product.image;
    productImage.innerHTML = `<img src="${imageUrl}" alt="${product.name}" class="main-image">`;

    // Thumbnails
    const thumbnails = document.querySelector('.thumbnails');
    thumbnails.innerHTML = '';
    const thumbnail = document.createElement('div');
    thumbnail.className = 'thumbnail active';
    const thumbImg = document.createElement('img');
    thumbImg.src = imageUrl;
    thumbImg.alt = product.name;
    thumbnail.appendChild(thumbImg);
    thumbnails.appendChild(thumbnail);

    // Product info
    document.querySelector('h1').textContent = product.name;
    document.querySelector('.product-description').textContent = product.description;
    document.querySelector('.price').textContent = formatPrice(product.price);

    // Rating (default 5 stars)
    const ratingDiv = document.querySelector('.rating');
    ratingDiv.innerHTML = '';
    for (let i = 0; i < 5; i++) {
        const star = document.createElement('i');
        star.className = 'fa-solid fa-star';
        ratingDiv.appendChild(star);
    }

    // Color selector (hide by default)
    const colorSelector = document.querySelector('.color-selector');
    colorSelector.style.display = 'none';

    // Size selector
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
    // Keep only the title, remove other elements
    const titleElement = sidebar.querySelector('.sidebar-title');
    sidebar.innerHTML = '';
    sidebar.appendChild(titleElement);

    relatedProducts.forEach(product => {
        const cardContainer = document.createElement('div');
        cardContainer.className = 'sidebar-product-card';

        const img = document.createElement('img');
        const imageUrl = product.image_url ? (apiClient.baseURL + product.image_url) : product.image;
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