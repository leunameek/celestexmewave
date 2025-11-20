// Profile page functionality
document.addEventListener('DOMContentLoaded', async () => {
    // Check if user is logged in
    const token = localStorage.getItem('accessToken');
    if (!token) {
        window.location.href = 'login.html';
        return;
    }

    // Load user profile data
    try {
        await loadUserProfile();
    } catch (error) {
        console.error('Failed to load profile:', error);
        showNotification('Error al cargar el perfil', 'error');
    }

    // Setup form handlers
    setupProfileUpdateForm();
    setupPasswordChangeForm();
    setupDeleteAccountForm();
});

// Load user profile data from API
async function loadUserProfile() {
    const user = await apiClient.get('/api/users/profile');

    // Populate form fields
    document.getElementById('first-name').value = user.first_name || '';
    document.getElementById('last-name').value = user.last_name || '';
    const phoneValue = user.phone ? user.phone.replace(/^\+?57/, '') : '';
    document.getElementById('phone').value = phoneValue;
    document.getElementById('email').value = user.email || '';

    // Email is typically not editable after registration
    document.getElementById('email').disabled = true;
}

// Setup profile update form
function setupProfileUpdateForm() {
    const updateBtn = document.getElementById('update-profile-btn');
    
    updateBtn.addEventListener('click', async (e) => {
        e.preventDefault();
        
        const firstName = document.getElementById('first-name').value.trim();
        const lastName = document.getElementById('last-name').value.trim();
        const phoneInput = document.getElementById('phone').value.trim();
        const phone = phoneInput ? '+57' + phoneInput : '';
        
        // Validation
        if (!firstName || !lastName) {
            showNotification('Por favor completa todos los campos requeridos', 'error');
            return;
        }
        
        if (phone && !validatePhone(phone)) {
            showNotification('Formato de teléfono inválido. Debe tener 10 dígitos después de +57', 'error');
            return;
        }
        
        // Show loading state
        updateBtn.disabled = true;
        updateBtn.textContent = 'Actualizando...';
        
        await apiClient.put('/api/users/profile', {
            first_name: firstName,
            last_name: lastName,
            phone: phone
        });

        showNotification('Perfil actualizado exitosamente', 'success');
    });
}

// Setup password change form
function setupPasswordChangeForm() {
    const changePasswordBtn = document.getElementById('change-password-btn');
    
    changePasswordBtn.addEventListener('click', async (e) => {
        e.preventDefault();
        
        const currentPassword = document.getElementById('current-password').value;
        const newPassword = document.getElementById('new-password').value;
        const confirmPassword = document.getElementById('confirm-password').value;
        
        // Validation
        if (!currentPassword || !newPassword || !confirmPassword) {
            showNotification('Por favor completa todos los campos de contraseña', 'error');
            return;
        }
        
        if (newPassword !== confirmPassword) {
            showNotification('Las contraseñas no coinciden', 'error');
            return;
        }
        
        if (newPassword.length < 8) {
            showNotification('La nueva contraseña debe tener al menos 8 caracteres', 'error');
            return;
        }
        
        // Show loading state
        changePasswordBtn.disabled = true;
        changePasswordBtn.textContent = 'Cambiando...';
        
        await apiClient.put('/api/users/change-password', {
            current_password: currentPassword,
            new_password: newPassword
        });

        showNotification('Contraseña cambiada exitosamente', 'success');

        // Clear password fields
        document.getElementById('current-password').value = '';
        document.getElementById('new-password').value = '';
        document.getElementById('confirm-password').value = '';
    });
}

// Setup delete account form
function setupDeleteAccountForm() {
    const deleteBtn = document.getElementById('delete-account-btn');

    deleteBtn.addEventListener('click', async (e) => {
        e.preventDefault();

        const confirmed = confirm('¿Estás seguro de que quieres eliminar tu cuenta? Esta acción no se puede deshacer.');
        if (!confirmed) return;

        const confirmedAgain = confirm('Esta es tu última oportunidad. ¿Realmente quieres eliminar tu cuenta permanentemente?');
        if (!confirmedAgain) return;

        deleteBtn.disabled = true;
        deleteBtn.textContent = 'Eliminando...';

        try {
            await apiClient.deleteUserProfile();

            showNotification('Cuenta eliminada exitosamente', 'success');

            // Logout and redirect
            setTimeout(() => {
                apiClient.logout();
                window.location.href = '../index.html';
            }, 2000);
        } catch (error) {
            console.error('Error deleting account:', error);
            showNotification('Error al eliminar la cuenta', 'error');
            deleteBtn.disabled = false;
            deleteBtn.textContent = 'Eliminar Cuenta';
        }
    });
}

// Validate phone number format
function validatePhone(phone) {
    // Colombian phone format: +57 followed by 10 digits
    const phoneRegex = /^\+57\d{10}$/;
    return phoneRegex.test(phone.replace(/\s/g, ''));
}

// Show notification to user
function showNotification(message, type = 'info') {
    // Create notification element
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;
    
    // Style the notification
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 15px 20px;
        background-color: ${type === 'success' ? '#4CAF50' : type === 'error' ? '#f44336' : '#2196F3'};
        color: white;
        border-radius: 5px;
        box-shadow: 0 2px 5px rgba(0,0,0,0.2);
        z-index: 10000;
        animation: slideIn 0.3s ease-out;
    `;
    
    // Add to page
    document.body.appendChild(notification);
    
    // Remove after 3 seconds
    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => {
            document.body.removeChild(notification);
        }, 300);
    }, 3000);
}

// Add CSS animations
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
    
    @keyframes slideOut {
        from {
            transform: translateX(0);
            opacity: 1;
        }
        to {
            transform: translateX(100%);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);