// Todo lo serio de las peticiones pasa por aki, sin tanto protocolo

class APIClient {
  constructor(baseURL = 'http://localhost:8080') {
    this.baseURL = baseURL;
    this.accessToken = localStorage.getItem('accessToken');
    this.refreshToken = localStorage.getItem('refreshToken');
  }

  // Cambio rapdio de base URL cuando toque
  setBaseURL(url) {
    this.baseURL = url;
  }

  // Armamos headers y metemos el auth si toca, sin lio
  getHeaders(includeAuth = true) {
    const headers = {
      'Content-Type': 'application/json',
    };

    if (includeAuth && this.accessToken) {
      headers['Authorization'] = `Bearer ${this.accessToken}`;
    }

    return headers;
  }

  // Wrapper de fetch pa no repetir codigoxd
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      ...options,
      headers: {
        ...this.getHeaders(options.includeAuth !== false),
        ...options.headers,
      },
    };

    try {
      const response = await fetch(url, config);

      // Si nos botan con 401, probamos refrescar token tranqui
      if (response.status === 401 && this.refreshToken) {
        const refreshed = await this.refreshAccessToken();
        if (refreshed) {
          // Reintentamos la solicitud con token nuevo, suave
          config.headers['Authorization'] = `Bearer ${this.accessToken}`;
          return fetch(url, config);
        }
      }

      if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        const message = error.error || error.message || `HTTP ${response.status}`;
        throw new Error(message);
      }

      return response;
    } catch (error) {
      console.error('API Request Error:', error);
      throw error;
    }
  }

  // GETcito chill
  async get(endpoint, options = {}) {
    const response = await this.request(endpoint, {
      ...options,
      method: 'GET',
    });
    return response.json();
  }

  // POST relajado
  async post(endpoint, data = {}, options = {}) {
    const response = await this.request(endpoint, {
      ...options,
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.json();
  }

  // PUT pa mandar updates
  async put(endpoint, data = {}, options = {}) {
    const response = await this.request(endpoint, {
      ...options,
      method: 'PUT',
      body: JSON.stringify(data),
    });
    return response.json();
  }

  // DELETE para borrar cositas
  async delete(endpoint, options = {}) {
    const response = await this.request(endpoint, {
      ...options,
      method: 'DELETE',
    });
    return response.json();
  }

  // ENDPOINTS DE AUTH
  // Registro chill
  async register(email, phone, firstName, lastName, password) {
    const data = await this.post('/api/auth/register', {
      email,
      phone,
      first_name: firstName,
      last_name: lastName,
      password,
    }, { includeAuth: false });

    if (data.access_token) {
      this.setTokens(data.access_token, data.refresh_token);
    }

    return data;
  }

  // Login basico
  async login(emailOrPhone, password) {
    const data = await this.post('/api/auth/login', {
      email: emailOrPhone,
      password,
    }, { includeAuth: false });

    if (data.access_token) {
      this.setTokens(data.access_token, data.refresh_token);
    }

    return data;
  }

  // Refresh de token cuando se vence
  async refreshAccessToken() {
    try {
      const data = await this.post('/api/auth/refresh-token', {
        refresh_token: this.refreshToken,
      }, { includeAuth: false });

      if (data.access_token) {
        this.accessToken = data.access_token;
        localStorage.setItem('accessToken', data.access_token);
        return true;
      }
      return false;
    } catch (error) {
      console.error('Token refresh failed:', error);
      this.logout();
      return false;
    }
  }

  // Logout y limpiamos todo
  async logout() {
    try {
      await this.post('/api/auth/logout', {}, { includeAuth: true });
    } catch (error) {
      console.error('Logout error:', error);
    }

    this.clearTokens();
  }

  // Pedir codigo para reset de clave
  async requestPasswordReset(emailOrPhone) {
    const payload = emailOrPhone.includes('@')
      ? { email: emailOrPhone }
      : { phone: emailOrPhone };

    return this.post('/api/auth/request-password-reset', payload, { includeAuth: false });
  }

  // Verificar codigo y guardar clave nueva
  async verifyResetCode(emailOrPhone, resetCode, newPassword) {
    const payload = emailOrPhone.includes('@')
      ? { email: emailOrPhone }
      : { phone: emailOrPhone };

    return this.post('/api/auth/verify-reset-code', {
      ...payload,
      reset_code: resetCode,
      new_password: newPassword,
    }, { includeAuth: false });
  }

  // ENDPOINTS DE USUARIO
  // Perfil del usuario
  async getUserProfile() {
    return this.get('/api/users/profile');
  }

  // Actualizar perfil
  async updateUserProfile(firstName, lastName, phone) {
    return this.put('/api/users/profile', {
      first_name: firstName,
      last_name: lastName,
      phone,
    });
  }

  // Pedidos del usuario
  async getUserOrders(page = 1, limit = 10) {
    return this.get(`/api/users/orders?page=${page}&limit=${limit}`);
  }

  // Borrar cuenta
  async deleteUserProfile() {
    return this.delete('/api/users/profile');
  }

  // ENDPOINTS DE PRODUCTO
  // Traer productos segun filtros
  async getAllProducts(store = '', category = '', minPrice = 0, maxPrice = 999999, page = 1, limit = 20) {
    const params = new URLSearchParams({
      store,
      category,
      min_price: minPrice,
      max_price: maxPrice,
      page,
      limit,
    });
    return this.get(`/api/products?${params.toString()}`);
  }

  // Un producto puntual
  async getProductByID(productID) {
    return this.get(`/api/products/${productID}`);
  }

  // Productos de una tienda
  async getProductsByStore(storeID, page = 1, limit = 20) {
    return this.get(`/api/products/store/${storeID}?page=${page}&limit=${limit}`);
  }

  // Productos por categoria
  async getProductsByCategory(category, page = 1, limit = 20) {
    return this.get(`/api/products/category/${category}?page=${page}&limit=${limit}`);
  }

  // ENDPOINTS DEL CARRITO
  // Ver carrito
  async getCart() {
    return this.get('/api/cart');
  }

  // Agregar item al carrito
  async addItemToCart(productID, quantity, size = '') {
    return this.post('/api/cart/items', {
      product_id: productID,
      quantity,
      size,
    });
  }

  // Actualizar item del carrito
  async updateCartItem(itemID, quantity, size = '') {
    return this.put(`/api/cart/items/${itemID}`, {
      quantity,
      size,
    });
  }

  // Quitar item del carrito
  async removeItemFromCart(itemID) {
    return this.delete(`/api/cart/items/${itemID}`);
  }

  // Vaciar carrito
  async clearCart() {
    return this.delete('/api/cart');
  }

  // ENDPOINTS DE PEDIDOS
  // Crear pedido
  async createOrder(sessionId, shippingInfo) {
    return this.post('/api/orders', {
      session_id: sessionId,
      shipping_name: shippingInfo.name,
      shipping_phone: shippingInfo.phone,
      shipping_email: shippingInfo.email,
      shipping_city: shippingInfo.city,
      shipping_address: shippingInfo.address,
      shipping_address2: shippingInfo.address2,
      shipping_postal_code: shippingInfo.postalCode,
      shipping_notes: shippingInfo.notes,
    });
  }

  // Pedido por id
  async getOrder(orderID) {
    return this.get(`/api/orders/${orderID}`);
  }

  // Pedidos paginados
  async getOrders(page = 1, limit = 10) {
    return this.get(`/api/orders?page=${page}&limit=${limit}`);
  }

  // Pago de un pedido
  async processPayment(orderID, cardDetails) {
    return this.post(`/api/orders/${orderID}/payment`, {
      card_number: cardDetails.number,
      card_holder: cardDetails.holder,
      expiry_month: cardDetails.expiryMonth,
      expiry_year: cardDetails.expiryYear,
      cvv: cardDetails.cvv,
    });
  }

  // Confirmacion del pedido
  async getOrderConfirmation(orderID) {
    return this.get(`/api/orders/${orderID}/confirmation`);
  }

  // TOKENS Y SESION
  // Guardar tokens
  setTokens(accessToken, refreshToken) {
    this.accessToken = accessToken;
    this.refreshToken = refreshToken;
    localStorage.setItem('accessToken', accessToken);
    localStorage.setItem('refreshToken', refreshToken);
    localStorage.setItem('isLoggedIn', 'true');
  }

  // Limpiar tokens
  clearTokens() {
    this.accessToken = null;
    this.refreshToken = null;
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('isLoggedIn');
  }

  // Saber si el user sigue logueado
  isAuthenticated() {
    return !!this.accessToken;
  }

  // Leer datos del token actual
  getCurrentUser() {
    if (!this.accessToken) return null;

    try {
      const payload = JSON.parse(atob(this.accessToken.split('.')[1]));
      return {
        id: payload.sub,
        email: payload.email,
        firstName: payload.first_name,
        lastName: payload.last_name,
      };
    } catch (error) {
      console.error('Error parsing token:', error);
      return null;
    }
  }
}

// Instanciamos el cliente global pa no repetir codigo
console.log('Initializing API Client...');
window.apiClient = new APIClient(
  localStorage.getItem('apiBaseURL') || 'http://localhost:8080'
);
console.log('API Client initialized:', window.apiClient);
