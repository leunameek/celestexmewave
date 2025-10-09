document.addEventListener('DOMContentLoaded', function() {
    const selectedProduct = JSON.parse(localStorage.getItem('selectedProduct'));
    if (!selectedProduct) {
        console.error('No product selected');
        return;
    }

    const urlParams = new URLSearchParams(window.location.search);
    const store = urlParams.get('store') || 'mewave';
    document.body.classList.add(store === 'celeste' ? 'celeste-product' : 'mewave-product');
    const jsonFile = store === 'celeste' ? 'celeste.json' : 'mewave.json';

    fetch('../assets/products/' + jsonFile)
        .then(response => response.json())
        .then(products => {
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

    const productImage = document.querySelector('.product-image');
    productImage.innerHTML = `<img src="${product.image}" alt="${product.name}" class="main-image">`;

    const thumbnails = document.querySelector('.thumbnails');
    thumbnails.innerHTML = `<div class="thumbnail active"><img src="${product.image}" alt="${product.name}"></div>`;

    document.querySelector('h1').textContent = product.name;
    document.querySelector('.product-description').textContent = product.description;
    document.querySelector('.price').textContent = formatPrice(product.price);

    const colorSelector = document.querySelector('.color-selector');
    colorSelector.style.display = 'none';

    const sizesContainer = document.querySelector('.sizes');
    sizesContainer.innerHTML = '';
    product.sizes.forEach((size, index) => {
        const sizeOption = document.createElement('div');
        sizeOption.className = 'size-option' + (index === 1 ? ' active' : '');
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
    title.textContent = 'Recomendados';
    title.className = 'sidebar-title' + (store === 'celeste' ? ' celeste-title' : '');

    const relatedProducts = products.filter(p => p.category === category && p.name !== currentProductName).slice(0, 2);

    const sidebar = document.querySelector('.sidebar');
    const titleElement = sidebar.querySelector('.sidebar-title');
    sidebar.innerHTML = '';
    sidebar.appendChild(titleElement);

    relatedProducts.forEach(product => {
        const img = createSidebarCard(product, store);
        sidebar.appendChild(img);
    });
}

function createSidebarCard(product, store) {
    const img = document.createElement('img');
    img.src = product.image;
    img.alt = product.name;
    img.style.width = '100%';
    img.style.height = '300px';
    img.style.objectFit = 'cover';
    img.style.borderRadius = '8px';
    img.style.cursor = 'pointer';
    img.addEventListener('click', () => {
        localStorage.setItem('selectedProduct', JSON.stringify(product));
        window.location.href = 'product.html?store=' + store;
    });
    return img;
}

function formatPrice(price) {
    return '$' + (price / 1000).toFixed(0) + '.000';
}