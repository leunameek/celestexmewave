document.addEventListener('DOMContentLoaded', function () {
  // Use API Client to fetch products
  apiClient.getAllProducts('Mewave')
    .then(response => {
      const products = response.products || [];
      console.log('Fetched products for Mewave:', products); // Debug log
      const categories = {};
      products.forEach(product => {
        if (!categories[product.category]) {
          categories[product.category] = [];
        }
        categories[product.category].push(product);
      });

      const main = document.querySelector('main');

      Object.keys(categories).forEach(category => {
        const section = document.createElement('section');
        section.className = 'categoria';

        const h2 = document.createElement('h2');
        h2.textContent = category;
        section.appendChild(h2);

        const grid = document.createElement('div');
        grid.className = 'productos-grid';

        categories[category].forEach(product => {
          const card = createProductCard(product);
          grid.appendChild(card);
        });

        section.appendChild(grid);
        main.appendChild(section);
      });
    })
    .catch(error => console.error('Error loading products:', error));
});

function createProductCard(product) {
  const card = document.createElement('div');
  card.className = 'producto-card';

  const img = document.createElement('img');
  const imageUrl = product.image_url ? (apiClient.baseURL + product.image_url) : product.image;
  img.src = imageUrl;
  img.alt = product.name;
  card.appendChild(img);

  const info = document.createElement('div');
  info.className = 'info';

  const p = document.createElement('p');
  p.textContent = product.category;
  info.appendChild(p);

  const h3 = document.createElement('h3');
  h3.textContent = product.name;
  info.appendChild(h3);

  const precio = document.createElement('span');
  precio.className = 'precio';
  precio.textContent = formatPrice(product.price);
  info.appendChild(precio);

  const button = document.createElement('button');
  const svg = `<svg width="28" height="28" viewBox="0 0 28 28" xmlns="http://www.w3.org/2000/svg" class="cart-icon">
                <path d="M9.91671 22.1667C9.57059 22.1667 9.23224 22.2693 8.94446 22.4616C8.65667 22.6539 8.43237 22.9272 8.29992 23.247C8.16746 23.5667 8.13281 23.9186 8.20033 24.2581C8.26786 24.5975 8.43453 24.9094 8.67927 25.1541C8.92401 25.3988 9.23583 25.5655 9.5753 25.633C9.91477 25.7006 10.2666 25.6659 10.5864 25.5335C10.9062 25.401 11.1795 25.1767 11.3718 24.8889C11.5641 24.6011 11.6667 24.2628 11.6667 23.9167C11.6667 23.4525 11.4823 23.0074 11.1541 22.6792C10.826 22.351 10.3808 22.1667 9.91671 22.1667ZM22.1667 18.6667H8.16671C7.85729 18.6667 7.56054 18.5437 7.34175 18.325C7.12296 18.1062 7.00004 17.8094 7.00004 17.5C7.00004 17.1906 7.12296 16.8938 7.34175 16.675C7.56054 16.4562 7.85729 16.3333 8.16671 16.3333H18.0731C18.8331 16.3309 19.5718 16.0822 20.1786 15.6246C20.7854 15.1669 21.2275 14.525 21.4387 13.7949L23.2884 7.32072C23.338 7.14708 23.3466 6.9643 23.3136 6.78676C23.2806 6.60922 23.2068 6.44177 23.0981 6.29759C22.9894 6.1534 22.8486 6.03643 22.687 5.95587C22.5254 5.87531 22.3473 5.83336 22.1667 5.83333H7.86229C7.62077 5.15359 7.17556 4.56483 6.5873 4.14728C5.99905 3.72974 5.29636 3.50371 4.575 3.5H3.50004C3.19062 3.5 2.89388 3.62292 2.67508 3.84171C2.45629 4.0605 2.33337 4.35725 2.33337 4.66667C2.33337 4.97609 2.45629 5.27283 2.67508 5.49162C2.89388 5.71042 3.19062 5.83333 3.50004 5.83333H4.575C4.82826 5.83422 5.07442 5.91712 5.27662 6.06962C5.47883 6.22212 5.62619 6.43602 5.69665 6.67928L5.87809 7.31481L5.87837 7.32072L7.79244 14.0199C6.90176 14.1158 6.08173 14.5494 5.50118 15.2317C4.92064 15.9139 4.62379 16.7928 4.67175 17.6873C4.71972 18.5819 5.10885 19.424 5.75901 20.0402C6.40918 20.6565 7.27089 21 8.16671 21H22.1667C22.4761 21 22.7729 20.8771 22.9917 20.6583C23.2105 20.4395 23.3334 20.1428 23.3334 19.8333C23.3334 19.5239 23.2105 19.2272 22.9917 19.0084C22.7729 18.7896 22.4761 18.6667 22.1667 18.6667ZM20.6201 8.16667L19.1953 13.1535C19.125 13.3969 18.9776 13.611 18.7752 13.7637C18.5729 13.9163 18.3266 13.9992 18.0731 14H10.2135L9.91607 12.9591L8.54753 8.16667H20.6201ZM19.25 22.1667C18.9039 22.1667 18.5656 22.2693 18.2778 22.4616C17.99 22.6539 17.7657 22.9272 17.6333 23.247C17.5008 23.5667 17.4661 23.9186 17.5337 24.2581C17.6012 24.5975 17.7679 24.9094 18.0126 25.1541C18.2573 25.3988 18.5692 25.5655 18.9086 25.633C19.2481 25.7006 19.6 25.6659 19.9197 25.5335C20.2395 25.401 20.5128 25.1767 20.7051 24.8889C20.8974 24.6011 21 24.2628 21 23.9167C21 23.4525 20.8157 23.0074 20.4875 22.6792C20.1593 22.351 19.7142 22.1667 19.25 22.1667Z"/>
              </svg>`;
  button.innerHTML = svg;
  info.appendChild(button);

  card.appendChild(info);

  card.style.cursor = 'pointer';
  card.addEventListener('click', () => {
    localStorage.setItem('selectedProduct', JSON.stringify(product));
    window.location.href = 'product.html?store=mewave';
  });

  return card;
}

function formatPrice(price) {
  return '$' + (price / 1000).toFixed(0) + '.000';
}