document.addEventListener('DOMContentLoaded', async () => {
  const params = new URLSearchParams(window.location.search);
  const query = (params.get('q') || '').trim();

  const termSpans = document.querySelectorAll('.search-term');
  termSpans.forEach(span => span.textContent = query || '...'); // reflejamos el texto arriba para que se vea bacano

  const subtitle = document.querySelector('.search-subtitle');
  const resultsGrid = document.getElementById('search-results');
  const emptyState = document.getElementById('search-empty');

  if (!query) {
    subtitle.textContent = 'Escribe algo en la barra de búsqueda para empezar.';
    emptyState.classList.remove('hide');
    return;
  }

  subtitle.textContent = 'Buscando productos...';

  try {
    const stores = ['Celeste', 'Mewave'];
    const allProducts = [];

    for (const store of stores) {
      const response = await apiClient.getAllProducts(store, '', 0, 999999, 1, 200);
      const products = response.products || [];
      products.forEach(p => allProducts.push({ ...p, store }));
    }

    const q = query.toLowerCase();
    const filtered = allProducts.filter(p => {
      const fields = [p.name, p.category, p.description];
      return fields.some(f => typeof f === 'string' && f.toLowerCase().includes(q));
    });

    if (!filtered.length) {
      subtitle.textContent = '0 resultados';
      emptyState.classList.remove('hide');
      return;
    }

    subtitle.textContent = `${filtered.length} resultado${filtered.length === 1 ? '' : 's'}`;
    emptyState.classList.add('hide');
    resultsGrid.innerHTML = '';

    filtered.forEach(product => {
      const card = createSearchCard(product);
      resultsGrid.appendChild(card);
    });
  } catch (error) {
    console.error('Search error:', error);
    subtitle.textContent = 'No pudimos completar la búsqueda.';
    emptyState.classList.remove('hide');
  }
});

function createSearchCard(product) {
  const card = document.createElement('div');
  card.className = 'search-card';

  const imgWrap = document.createElement('div');
  imgWrap.className = 'search-card-image';
  const img = document.createElement('img');
  const imageUrl = product.image_url ? ('../assets/images/' + product.image_url) : product.image;
  img.src = imageUrl;
  img.alt = product.name;
  imgWrap.appendChild(img);

  const info = document.createElement('div');
  info.className = 'search-card-info';

  const tag = document.createElement('span');
  tag.className = 'search-card-tag';
  tag.textContent = product.store;
  const isCeleste = (product.store || '').toLowerCase() === 'celeste';
  tag.dataset.store = isCeleste ? 'celeste' : 'mewave';

  const title = document.createElement('h3');
  title.textContent = product.name;

  const category = document.createElement('p');
  category.className = 'search-card-category';
  category.textContent = product.category;

  const price = document.createElement('div');
  price.className = 'search-card-price';
  price.textContent = formatPrice(product.price);

  info.appendChild(tag);
  info.appendChild(title);
  info.appendChild(category);
  info.appendChild(price);

  card.appendChild(imgWrap);
  card.appendChild(info);

  card.addEventListener('click', () => {
    localStorage.setItem('selectedProduct', JSON.stringify(product));
    const storeParam = product.store && product.store.toLowerCase() === 'celeste' ? 'celeste' : 'mewave';
    window.location.href = `product.html?store=${storeParam}`;
  });

  return card;
}

function formatPrice(price) {
  return '$' + (price / 1000).toFixed(0) + '.000';
}
