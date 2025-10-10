document.addEventListener('DOMContentLoaded', function() {
    const step1 = document.getElementById('step1');
    const step2 = document.getElementById('step2');
    const step3 = document.getElementById('step3');

    const sendCodeBtn = document.getElementById('send-code-btn');
    const verifyCodeBtn = document.getElementById('verify-code-btn');
    const savePasswordBtn = document.getElementById('save-password-btn');

    const recoveryEmail = document.getElementById('recovery-email');
    const recoveryCode = document.getElementById('recovery-code');
    const newPassword = document.getElementById('new-password');
    const confirmPassword = document.getElementById('confirm-password');

    sendCodeBtn.addEventListener('click', function() {
        if (recoveryEmail.value.trim() === '') {
            alert('Por favor, ingresa tu email.');
            return;
        }
        alert('Código enviado a tu email.');
        step1.style.display = 'none';
        step2.style.display = 'block';
    });

    verifyCodeBtn.addEventListener('click', function() {
        if (recoveryCode.value.trim() === '') {
            alert('Por favor, ingresa el código.');
            return;
        }
        alert('Código verificado.');
        step2.style.display = 'none';
        step3.style.display = 'block';
    });

    savePasswordBtn.addEventListener('click', function() {
        if (newPassword.value.trim() === '' || confirmPassword.value.trim() === '') {
            alert('Por favor, completa todos los campos.');
            return;
        }
        if (newPassword.value !== confirmPassword.value) {
            alert('Las contraseñas no coinciden.');
            return;
        }
        alert('Contraseña actualizada exitosamente.');
        window.location.href = 'login.html';
    });
});