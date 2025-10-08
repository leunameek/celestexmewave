document.addEventListener('DOMContentLoaded', function() {
  fetch('../assets/products/mewave.json')
    .then(response => response.json())
    .then(products => {
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
  img.src = product.image;
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
  const buttonImg = document.createElement('img');
  buttonImg.src = '../assets/icons/cart.svg';
  buttonImg.alt = 'Cart';
  button.appendChild(buttonImg);
  info.appendChild(button);

  card.appendChild(info);

  return card;
}

function formatPrice(price) {
  return '$' + (price / 1000).toFixed(0) + '.000';
}