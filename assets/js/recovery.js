document.addEventListener('DOMContentLoaded', function() {
    const step1 = document.getElementById('step1');
    const step2 = document.getElementById('step2');

    const sendCodeBtn = document.getElementById('send-code-btn');
    const savePasswordBtn = document.getElementById('save-password-btn');

    const recoveryEmail = document.getElementById('recovery-email');
    const recoveryCode = document.getElementById('recovery-code');
    const newPassword = document.getElementById('new-password');
    const confirmPassword = document.getElementById('confirm-password');
    const messageBox = document.getElementById('recovery-message');

    let currentIdentifier = '';

    function showMessage(text, type = 'info') {
        if (!messageBox) return;
        messageBox.textContent = text;
        messageBox.style.display = 'block';
        messageBox.style.color = type === 'error' ? '#c62828' : '#2e7d32';
    }

    function hideMessage() {
        if (messageBox) {
            messageBox.style.display = 'none';
            messageBox.textContent = '';
        }
    }

    sendCodeBtn.addEventListener('click', async function() {
        hideMessage();
        const identifier = recoveryEmail.value.trim();
        if (!identifier) {
            showMessage('Por favor, ingresa tu email o teléfono.', 'error');
            return;
        }

        try {
            sendCodeBtn.disabled = true;
            sendCodeBtn.textContent = 'Enviando...';
            currentIdentifier = identifier;
            await apiClient.requestPasswordReset(identifier);
            showMessage('Código enviado. Revisa tu correo o teléfono.');
            step1.style.display = 'none';
            step2.style.display = 'block';
        } catch (error) {
            console.error('Recovery send code error:', error);
            showMessage(error.message || 'No pudimos enviar el código.', 'error');
        } finally {
            sendCodeBtn.disabled = false;
            sendCodeBtn.textContent = 'Enviar código';
        }
    });

    savePasswordBtn.addEventListener('click', async function() {
        hideMessage();
        const code = recoveryCode.value.trim();
        const password = newPassword.value.trim();
        const confirm = confirmPassword.value.trim();

        if (!code || !password || !confirm) {
            showMessage('Por favor, completa todos los campos.', 'error');
            return;
        }
        if (password !== confirm) {
            showMessage('Las contraseñas no coinciden.', 'error');
            return;
        }
        if (!currentIdentifier) {
            showMessage('Primero solicita el código.', 'error');
            return;
        }

        try {
            savePasswordBtn.disabled = true;
            savePasswordBtn.textContent = 'Guardando...';
            await apiClient.verifyResetCode(currentIdentifier, code, password);
            showMessage('Contraseña actualizada correctamente.');
            setTimeout(() => {
                window.location.href = 'login.html';
            }, 1200);
        } catch (error) {
            console.error('Recovery save error:', error);
            showMessage(error.message || 'No pudimos actualizar la contraseña.', 'error');
        } finally {
            savePasswordBtn.disabled = false;
            savePasswordBtn.textContent = 'Guardar';
        }
    });
});
