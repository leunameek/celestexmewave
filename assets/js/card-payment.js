document.addEventListener('DOMContentLoaded', function() {
  const paymentMethods = document.getElementById('payment-methods');
  const cardVisual = document.getElementById('card-visual');
  const cardForm = document.getElementById('card-form');
  const card = document.getElementById('card');
  const changeMethodBtn = document.getElementById('change-method-btn');

  // Miramos si esto es pa terminar un pago que ya estaba empezado
  const urlParams = new URLSearchParams(window.location.search);
  const existingOrderId = urlParams.get('order_id');

  if (existingOrderId) {
    // Escondemos envio porque ya viene listo ese pedido
    document.getElementById('shipping-form').style.display = 'none';
    document.querySelector('.checkout-content').style.justifyContent = 'center';

    // Cargamos el pedido y mostramos los metodos chill
    loadOrderForPayment(existingOrderId);
  }

  // Eleccion del metodo de pago, sin complique
  paymentMethods.addEventListener('click', function(e) {
    const method = e.target.closest('svg')?.getAttribute('data-method');
    if (method === 'mastercard' || method === 'visa') {
      // Ocultamos los logos
      paymentMethods.style.display = 'none';
      // Mostramos el cambio, la tarjetica y el form
      changeMethodBtn.style.display = 'block';
      cardVisual.style.display = 'block';
      cardForm.style.display = 'block';

      // Fondo de la tarjeta segun el metodo, para el flow
      if (method === 'visa') {
        card.querySelector('.card-front').style.background = 'linear-gradient(135deg, #1a1f71 0%, #003d82 100%)';
      } else if (method === 'mastercard') {
        card.querySelector('.card-front').style.background = 'linear-gradient(135deg, #eb001b 0%, #f79e1b 100%)';
      }
    }
  });

  // Cambiar de metodo si nos arrepentimos
  changeMethodBtn.addEventListener('click', function() {
    // Escondemos todo
    changeMethodBtn.style.display = 'none';
    cardVisual.style.display = 'none';
    cardForm.style.display = 'none';
    // Volvemos a mostrar logos
    paymentMethods.style.display = 'flex';
    // Limpiamos el form para no dejar basura
    document.getElementById('card-number-input').value = '';
    document.getElementById('card-name-input').value = '';
    document.getElementById('card-expiry-input').value = '';
    document.getElementById('card-cvv-input').value = '';
    cardNumberDisplay.textContent = '#### #### #### ####';
    cardNameDisplay.textContent = 'NOMBRE COMPLETO';
    cardExpiryDisplay.textContent = 'MM/YY';
    cardCvvDisplay.textContent = 'CVV';
    card.classList.remove('card-flip');
  });

  // Escuchas para actualizar en vivo
  const cardNumberInput = document.getElementById('card-number-input');
  const cardNameInput = document.getElementById('card-name-input');
  const cardExpiryInput = document.getElementById('card-expiry-input');
  const cardCvvInput = document.getElementById('card-cvv-input');

  const cardNumberDisplay = document.getElementById('card-number');
  const cardNameDisplay = document.getElementById('card-name');
  const cardExpiryDisplay = document.getElementById('card-expiry');
  const cardCvvDisplay = document.getElementById('card-cvv');

  // Formatear numero de tarjeta todo bonito
  cardNumberInput.addEventListener('input', function(e) {
    let value = e.target.value.replace(/\s+/g, '').replace(/[^0-9]/gi, '');
    let formatted = value.replace(/(.{4})/g, '$1 ').trim();
    e.target.value = formatted;
    cardNumberDisplay.textContent = formatted || '#### #### #### ####';
  });

  // Actualizar el nombre de la tarjeta
  cardNameInput.addEventListener('input', function(e) {
    cardNameDisplay.textContent = e.target.value.toUpperCase() || 'NOMBRE COMPLETO';
  });

  // Formatear y validar fecha de expi jeje
  cardExpiryInput.addEventListener('input', function(e) {
    let value = e.target.value.replace(/\D/g, '');
    if (value.length >= 2) {
      value = value.slice(0, 2) + '/' + value.slice(2, 4);
    }
    e.target.value = value;

    // Validamos MM/YY a lo basico
    if (value.length === 5) {
      const mm = parseInt(value.slice(0, 2));
      const yy = parseInt('20' + value.slice(3, 5));
      const currentYear = new Date().getFullYear();
      const currentMonth = new Date().getMonth() + 1; // 1-12 pues

      if (mm < 1 || mm > 12) {
        e.target.setCustomValidity('Mes inválido (01-12)');
        e.target.style.borderColor = 'red';
      } else if (yy < currentYear || (yy === currentYear && mm <= currentMonth)) {
        e.target.setCustomValidity('Fecha de expiración inválida');
        e.target.style.borderColor = 'red';
      } else {
        e.target.setCustomValidity('');
        e.target.style.borderColor = '#ccc';
      }
    } else {
      e.target.setCustomValidity('');
      e.target.style.borderColor = '#ccc';
    }

    cardExpiryDisplay.textContent = value || 'MM/YY';
  });

  // Actualizar el CVV
  cardCvvInput.addEventListener('input', function(e) {
    let value = e.target.value.replace(/\D/g, '');
    e.target.value = value;
    cardCvvDisplay.textContent = value || 'CVV';
  });

  // Volteamos la tarjeta si tocan el CVV
  cardCvvInput.addEventListener('focus', function() {
    card.classList.add('card-flip');
  });

  cardCvvInput.addEventListener('blur', function() {
    card.classList.remove('card-flip');
  });

  // Confirmar y pagar, sin tanto susto
  const confirmBtn = document.getElementById('confirm-btn');
  confirmBtn.addEventListener('click', async function(e) {
    e.preventDefault();

    // Recolectamos los datos de la tarjeta
    const cardNumber = cardNumberInput.value.replace(/\s/g, '');
    const cardHolder = cardNameInput.value;
    const expiry = cardExpiryInput.value;
    const cvv = cardCvvInput.value;

    if (!cardNumber || !cardHolder || !expiry || !cvv) {
      showNotification('Por favor complete toda la información de la tarjeta', 'error');
      return;
    }

    // Validar la longitud del numero (13-19 digitos)
    if (cardNumber.length < 13 || cardNumber.length > 19) {
      showNotification('El número de tarjeta debe tener entre 13 y 19 dígitos', 'error');
      return;
    }

    // Validar longitud del CVV
    if (cvv.length < 3 || cvv.length > 4) {
      showNotification('El CVV debe tener entre 3 y 4 dígitos', 'error');
      return;
    }

    // Validar formato de fecha
    const expiryRegex = /^(0[1-9]|1[0-2])\/\d{2}$/;
    if (!expiryRegex.test(expiry)) {
      showNotification('El formato de fecha de expiración debe ser MM/YY', 'error');
      return;
    }

    const expiryParts = expiry.split('/');
    const expiryMonth = parseInt(expiryParts[0]);
    const expiryYear = 2000 + parseInt(expiryParts[1]);

    try {
      let orderId;

      if (existingOrderId) {
        // Si ya habia pedido, solo pagamos
        orderId = existingOrderId;
      } else {
        // Buscamos session o creamos una
        let sessionId = localStorage.getItem('sessionId');
        if (!sessionId) {
          sessionId = generateSessionId();
          localStorage.setItem('sessionId', sessionId);
        }

        // Sacamos los datos de envio
        const shippingForm = document.getElementById('shipping-form');
        const shippingData = new FormData(shippingForm);
        const shippingInfo = {
          name: shippingData.get('shipping-name') || '',
          phone: shippingData.get('shipping-phone') || '',
          email: shippingData.get('shipping-email') || '',
          city: shippingData.get('shipping-city') || '',
          address: shippingData.get('shipping-address') || '',
          address2: shippingData.get('shipping-address2') || '',
          postalCode: shippingData.get('shipping-postal-code') || '',
          notes: shippingData.get('shipping-notes') || '',
        };

        // Validamos que no falte nada de envio
        if (!shippingInfo.name || !shippingInfo.phone || !shippingInfo.email || !shippingInfo.city || !shippingInfo.address) {
          showNotification('Por favor complete toda la información de envío', 'error');
          return;
        }

        // Creamos el pedido en la API
        const orderResponse = await apiClient.createOrder(sessionId, shippingInfo);
        orderId = orderResponse.id;
      }

      // Llamamos el pago
      const paymentResponse = await apiClient.processPayment(orderId, {
        number: cardNumber,
        holder: cardHolder,
        expiryMonth: expiryMonth,
        expiryYear: expiryYear,
        cvv: cvv,
      });

      if (paymentResponse.payment_status === 'completed') {
        // Aviso de exito todo happy
        showNotification('¡Pago procesado exitosamente!', 'success');
        // Redirigimos a la confirmacion con una pausita
        setTimeout(() => {
          window.location.href = `confirmation.html?order_id=${orderId}`;
        }, 1500);
      } else {
        showNotification('Pago fallido: ' + paymentResponse.message, 'error');
      }
    } catch (error) {
      console.error('Order/Payment error:', error);
      showNotification('Error al procesar el pago: ' + error.message, 'error');
    }
  });
});

async function loadOrderForPayment(orderId) {
  try {
    const order = await apiClient.getOrder(orderId);

    // Cambiamos titulo o mostramos daticos del pedido
    const title = document.querySelector('h1') || document.querySelector('title');
    if (title) {
      document.title = `Completar Pago - Pedido ${orderId.slice(0, 8)}`;
    }

    // Resumen rapido del pedido
    const paymentCard = document.querySelector('.payment-card h3');
    paymentCard.textContent = `Completar Pago - Pedido #${orderId.slice(0, 8)}`;

    // Metemos el resumen antes de los metodos de pago
    const orderSummary = document.createElement('div');
    orderSummary.className = 'order-summary';
    orderSummary.innerHTML = `
      <div class="order-summary-header">
        <h4>Resumen del Pedido</h4>
        <div class="order-total">Total: $${formatPriceColombian(order.total_amount)}</div>
      </div>
      <div class="order-items-summary">
        ${order.items.map(item => `<div>${item.product_name} x${item.quantity}</div>`).join('')}
      </div>
    `;

    const paymentCardDiv = document.querySelector('.payment-card');
    paymentCardDiv.insertBefore(orderSummary, document.getElementById('payment-methods'));

    // Mostramos los metodos de pago
    paymentMethods.style.display = 'flex';

  } catch (error) {
    console.error('Error loading order:', error);
    showNotification('Error al cargar el pedido', 'error');
  }
}

function generateSessionId() {
  return 'session_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
}

// Formato de precio a la criolla
function formatPriceColombian(price) {
  const intPrice = Math.floor(price);
  return intPrice.toLocaleString('es-CO').replace(/,/g, '.');
}

// Notificaciones sencillas hechas a mano
function showNotification(message, type = 'info', duration = 5000) {
  // Quitamos las que ya hay
  const existing = document.querySelectorAll('.notification');
  existing.forEach(el => el.remove());

  // Creamos la nueva
  const notification = document.createElement('div');
  notification.className = `notification ${type}`;
  notification.textContent = message;

  // La pegamos al body
  document.body.appendChild(notification);

  // Se va solita despues de un rato
  setTimeout(() => {
    notification.classList.add('fade-out');
    setTimeout(() => {
      if (notification.parentNode) {
        notification.parentNode.removeChild(notification);
      }
    }, 300);
  }, duration);
}
