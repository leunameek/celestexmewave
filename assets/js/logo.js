document.addEventListener('DOMContentLoaded', function() {
    const celesteLogo = document.getElementById('celeste-logo');
    const mewaveLogo = document.getElementById('mewave-logo');
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
});