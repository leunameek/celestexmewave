/**
 * API Client for CelestexMewave Frontend
 * Handles all communication with the backend API
 */

class APIClient {
  constructor(baseURL = 'http://localhost:8080') {
    this.baseURL = baseURL;
    this.accessToken = localStorage.getItem('accessToken');
    this.refreshToken = localStorage.getItem('refreshToken');
  }

  /**
   * Set the base URL for API calls
   */
  setBaseURL(url) {
    this.baseURL = url;
  }

  /**
   * Get authorization headers
   */
  getHeaders(includeAuth = true) {
    const headers = {
      'Content-Type': 'application/json',
    };

    if (includeAuth && this.accessToken) {
      headers['Authorization'] = `Bearer ${this.accessToken}`;
    }

    return headers;
  }

  /**
   * Make a fetch request
   */
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

      // Handle 401 Unauthorized - try to refresh token
      if (response.status === 401 && this.refreshToken) {
        const refreshed = await this.refreshAccessToken();
        if (refreshed) {
          // Retry the original request with new token
          config.headers['Authorization'] = `Bearer ${this.accessToken}`;
          return fetch(url, config);
        }
      }

      if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw new Error(error.message || `HTTP ${response.status}`);
      }

      return response;
    } catch (error) {
      console.error('API Request Error:', error);
      throw error;
    }
  }

  /**
   * GET request
   */
  async get(endpoint, options = {}) {
    const response = await this.request(endpoint, {
      ...options,
      method: 'GET',
    });
    return response.json();
  }

  /**
   * POST request
   */
  async post(endpoint, data = {}, options = {}) {
    const response = await this.request(endpoint, {
      ...options,
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.json();
  }

  /**
   * PUT request
   */
  async put(endpoint, data = {}, options = {}) {
    const response = await this.request(endpoint, {
      ...options,
      method: 'PUT',
      body: JSON.stringify(data),
    });
    return response.json();
  }

  /**
   * DELETE request
   */
  async delete(endpoint, options = {}) {
    const response = await this.request(endpoint, {
      ...options,
      method: 'DELETE',
    });
    return response.json();
  }

  /**
   * AUTHENTICATION ENDPOINTS
   */

  /**
   * Register a new user
   */
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

  /**
   * Login user
   */
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

  /**
   * Refresh access token
   */
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

  /**
   * Logout user
   */
  async logout() {
    try {
      await this.post('/api/auth/logout', {}, { includeAuth: true });
    } catch (error) {
      console.error('Logout error:', error);
    }

    this.clearTokens();
  }

  /**
   * Request password reset
   */
  async requestPasswordReset(emailOrPhone) {
    return this.post('/api/auth/request-password-reset', {
      email_or_phone: emailOrPhone,
    }, { includeAuth: false });
  }

  /**
   * Verify reset code and update password
   */
  async verifyResetCode(emailOrPhone, resetCode, newPassword) {
    return this.post('/api/auth/verify-reset-code', {
      email_or_phone: emailOrPhone,
      reset_code: resetCode,
      new_password: newPassword,
    }, { includeAuth: false });
  }

  /**
   * USER ENDPOINTS
   */

  /**
   * Get user profile
   */
  async getUserProfile() {
    return this.get('/api/users/profile');
  }

  /**
   * Update user profile
   */
  async updateUserProfile(firstName, lastName, phone) {
    return this.put('/api/users/profile', {
      first_name: firstName,
      last_name: lastName,
      phone,
    });
  }

  /**
   * Get user orders
   */
  async getUserOrders(page = 1, limit = 10) {
    return this.get(`/api/users/orders?page=${page}&limit=${limit}`);
  }

  /**
   * Delete user profile
   */
  async deleteUserProfile() {
    return this.delete('/api/users/profile');
  }

  /**
   * PRODUCT ENDPOINTS
   */

  /**
   * Get all products
   */
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

  /**
   * Get product by ID
   */
  async getProductByID(productID) {
    return this.get(`/api/products/${productID}`);
  }

  /**
   * Get products by store
   */
  async getProductsByStore(storeID, page = 1, limit = 20) {
    return this.get(`/api/products/store/${storeID}?page=${page}&limit=${limit}`);
  }

  /**
   * Get products by category
   */
  async getProductsByCategory(category, page = 1, limit = 20) {
    return this.get(`/api/products/category/${category}?page=${page}&limit=${limit}`);
  }

  /**
   * CART ENDPOINTS
   */

  /**
   * Get cart
   */
  async getCart() {
    return this.get('/api/cart');
  }

  /**
   * Add item to cart
   */
  async addItemToCart(productID, quantity, size = '') {
    return this.post('/api/cart/items', {
      product_id: productID,
      quantity,
      size,
    });
  }

  /**
   * Update cart item
   */
  async updateCartItem(itemID, quantity, size = '') {
    return this.put(`/api/cart/items/${itemID}`, {
      quantity,
      size,
    });
  }

  /**
   * Remove item from cart
   */
  async removeItemFromCart(itemID) {
    return this.delete(`/api/cart/items/${itemID}`);
  }

  /**
   * Clear cart
   */
  async clearCart() {
    return this.delete('/api/cart');
  }

  /**
   * ORDER ENDPOINTS
   */

  /**
   * Create order
   */
  async createOrder(items, shippingAddress, billingAddress) {
    return this.post('/api/orders', {
      items,
      shipping_address: shippingAddress,
      billing_address: billingAddress,
    });
  }

  /**
   * Get order by ID
   */
  async getOrder(orderID) {
    return this.get(`/api/orders/${orderID}`);
  }

  /**
   * Get orders
   */
  async getOrders(page = 1, limit = 10) {
    return this.get(`/api/orders?page=${page}&limit=${limit}`);
  }

  /**
   * Process payment
   */
  async processPayment(orderID, paymentMethod, paymentDetails) {
    return this.post(`/api/orders/${orderID}/payment`, {
      payment_method: paymentMethod,
      payment_details: paymentDetails,
    });
  }

  /**
   * Get order confirmation
   */
  async getOrderConfirmation(orderID) {
    return this.get(`/api/orders/${orderID}/confirmation`);
  }

  /**
   * TOKEN MANAGEMENT
   */

  /**
   * Set tokens in storage
   */
  setTokens(accessToken, refreshToken) {
    this.accessToken = accessToken;
    this.refreshToken = refreshToken;
    localStorage.setItem('accessToken', accessToken);
    localStorage.setItem('refreshToken', refreshToken);
    localStorage.setItem('isLoggedIn', 'true');
  }

  /**
   * Clear tokens from storage
   */
  clearTokens() {
    this.accessToken = null;
    this.refreshToken = null;
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('isLoggedIn');
  }

  /**
   * Check if user is authenticated
   */
  isAuthenticated() {
    return !!this.accessToken;
  }

  /**
   * Get current user from token
   */
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

// Create global API client instance
console.log('Initializing API Client...');
window.apiClient = new APIClient(
  localStorage.getItem('apiBaseURL') || 'http://localhost:8080'
);
console.log('API Client initialized:', window.apiClient);
