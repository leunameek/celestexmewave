document.addEventListener('DOMContentLoaded', function() {
    const celesteLogo = document.getElementById('celeste-logo');
    const mewaveLogo = document.getElementById('mewave-logo');
    const userIcon = document.querySelector('.icons a[href*="login.html"] img, .icons a[href*="profile.html"] img');
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

    const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true';

    if (userIcon) {
        const dropdown = document.createElement('div');
        dropdown.className = 'user-dropdown';
        dropdown.innerHTML = `
            <div class="dropdown-content">
                ${isLoggedIn ? `
                    <a href="${basePath}profile.html">Ir al perfil</a>
                    <a href="#" id="logout">Cerrar sesión</a>
                ` : `
                    <a href="${basePath}profile.html">Ir al perfil</a>
                    <a href="${basePath}login.html">Iniciar sesión</a>
                `}
            </div>
        `;
        dropdown.style.display = 'none';
        userIcon.parentElement.appendChild(dropdown);

        userIcon.style.cursor = 'pointer';
        userIcon.addEventListener('click', (e) => {
            e.preventDefault();
            dropdown.style.display = dropdown.style.display === 'none' ? 'block' : 'none';
        });

        document.addEventListener('click', (e) => {
            if (!userIcon.contains(e.target) && !dropdown.contains(e.target)) {
                dropdown.style.display = 'none';
            }
        });

        const logoutLink = dropdown.querySelector('#logout');
        if (logoutLink) {
            logoutLink.addEventListener('click', () => {
                localStorage.removeItem('isLoggedIn');
                window.location.href = basePath + 'index.html';
            });
        }
    }
});