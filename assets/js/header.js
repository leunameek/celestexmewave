document.addEventListener('DOMContentLoaded', async function() {
    const celesteLogo = document.getElementById('celeste-logo');
    const mewaveLogo = document.getElementById('mewave-logo');
    const cartLink = document.querySelector('.icons a[href*="cart.html"]');
    if (cartLink) {
        cartLink.classList.add('cart-link');
    }

    
    // Buscamos el icono del user, el del monito
    const userIconLink = document.querySelector('.icons a[href*="login.html"], .icons a[href*="profile.html"]');
    const userIcon = userIconLink ? userIconLink.querySelector('img') : null;
    const basePath = window.location.pathname.includes('/pages/') ? '' : 'pages/';

    if (celesteLogo) {
        celesteLogo.style.cursor = 'pointer';
        celesteLogo.addEventListener('click', () => {
            window.location.href = basePath + 'celeste.html';
        });
    }

    if (mewaveLogo) {
        mewaveLogo.style.cursor = 'pointer';
        mewaveLogo.addEventListener('click', () => {
            window.location.href = basePath + 'mewave.html';
        });
    }

    if (userIcon && userIconLink) {
        // Chequeo rapido pa ver si esta logueado
        function checkAuthStatus() {
            const isAuth = typeof apiClient !== 'undefined' && apiClient.isAuthenticated();
            //console.log('Auth check:', isAuth, 'apiClient exists:', typeof apiClient !== 'undefined');
            return isAuth;
        }

        // Armamos el dropdown segun si esta o no logueado
        function createDropdownContent() {
            const isLoggedIn = checkAuthStatus();
            if (isLoggedIn) {
                return `
                    <div class="dropdown-content">
                        <a href="${basePath}profile.html">Ir al perfil</a>
                        <a href="#" id="logout">Cerrar sesión</a>
                    </div>
                `;
            } else {
                return `
                    <div class="dropdown-content">
                        <a href="${basePath}login.html">Iniciar sesión</a>
                        <a href="${basePath}register.html">Crear cuenta</a>
                    </div>
                `;
            }
        }

        const dropdown = document.createElement('div');
        dropdown.className = 'user-dropdown';
        dropdown.innerHTML = createDropdownContent();
        dropdown.style.display = 'none';
        userIconLink.parentElement.appendChild(dropdown);

        userIcon.style.cursor = 'pointer';
        userIconLink.addEventListener('click', (e) => {
            e.preventDefault();
            e.stopPropagation();
            
            // Refrescamos el contenido cada click, pa no quedar desactualizados
            const newContent = createDropdownContent();
            console.log('Dropdown content:', newContent);
            dropdown.innerHTML = newContent;
            
            // Mostrar u ocultar el menu a lo basico
            dropdown.style.display = dropdown.style.display === 'none' ? 'block' : 'none';
            
            // Pegamos listeners a los links del menu
            const dropdownLinks = dropdown.querySelectorAll('a');
            dropdownLinks.forEach(link => {
                link.addEventListener('click', async (e) => {
                    if (link.id === 'logout') {
                        e.preventDefault();
                        e.stopPropagation();
                        try {
                            if (typeof apiClient !== 'undefined') {
                                await apiClient.logout();
                            }
                        } catch (error) {
                            console.error('Logout error:', error);
                        }
                        window.location.href = basePath + 'login.html';
                    } else {
                        // Para los demas links, cerramos y dejamos ir
                        dropdown.style.display = 'none';
                    }
                });
            });
        });

        document.addEventListener('click', (e) => {
            if (!userIconLink.contains(e.target) && !dropdown.contains(e.target)) {
                dropdown.style.display = 'none';
            }
        });
    }

    // Badge del carrito, pa ver cuantas cositas hay
    let cartBadge;
    function ensureCartBadge() {
        if (!cartLink) return;
        cartBadge = cartLink.querySelector('.cart-count');
        if (!cartBadge) {
            cartBadge = document.createElement('span');
            cartBadge.className = 'cart-count';
            cartLink.appendChild(cartBadge);
        }
    }

    function setCartCount(count) {
        if (!cartLink) return;
        ensureCartBadge();
        if (count > 0) {
            cartBadge.textContent = count;
            cartBadge.classList.add('show');
        } else {
            cartBadge.textContent = '';
            cartBadge.classList.remove('show');
        }
    }

    async function refreshCartCount() {
        if (typeof apiClient === 'undefined') return;
        try {
            const cart = await apiClient.getCart();
            const count = (cart.items || []).reduce((sum, item) => sum + item.quantity, 0);
            setCartCount(count);
        } catch (err) {
            console.error('Cart count error:', err);
        }
    }

    window.refreshCartCount = refreshCartCount;
    window.setCartCount = setCartCount;

    await refreshCartCount();

    // Busqueda global, sin complicarse
    const searchBar = document.querySelector('.search-bar');
    if (searchBar) {
        const searchInput = searchBar.querySelector('input');
        const searchButton = searchBar.querySelector('button');
        const basePath = window.location.pathname.includes('/pages/') ? '' : 'pages/';

        const executeSearch = () => {
            const query = searchInput.value.trim();
            if (!query) return;
            window.location.href = `${basePath}search.html?q=${encodeURIComponent(query)}`;
        };

        searchInput.addEventListener('keydown', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                executeSearch();
            }
        });

        if (searchButton) {
            searchButton.addEventListener('click', (e) => {
                e.preventDefault();
                executeSearch();
            });
        }
    }
});
