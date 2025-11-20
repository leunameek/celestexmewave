document.addEventListener('DOMContentLoaded', async function() {
    const celesteLogo = document.getElementById('celeste-logo');
    const mewaveLogo = document.getElementById('mewave-logo');
    
    // Find user icon - look for any link with user.svg
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
        // Function to check authentication status
        function checkAuthStatus() {
            const isAuth = typeof apiClient !== 'undefined' && apiClient.isAuthenticated();
            console.log('Auth check:', isAuth, 'apiClient exists:', typeof apiClient !== 'undefined');
            return isAuth;
        }

        // Function to create dropdown content based on auth status
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
            
            // Refresh dropdown content on each click to check current auth status
            const newContent = createDropdownContent();
            console.log('Dropdown content:', newContent);
            dropdown.innerHTML = newContent;
            
            // Toggle dropdown visibility
            dropdown.style.display = dropdown.style.display === 'none' ? 'block' : 'none';
            
            // Attach event listeners to dropdown links
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
                        // For other links, close dropdown and allow navigation
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
});
