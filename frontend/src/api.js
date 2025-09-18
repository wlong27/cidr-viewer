import { getConfig, getConfigSync } from './config.js';

// API client functions with JSON configuration
export const api = {
  // Analyze multiple CIDR ranges with categorization
  async analyzeCIDRs(requestData, options = {}) {
    const config = await getConfig();
    const { signal, timeout = config.apiTimeout } = options;
    const controller = signal ? null : new AbortController();
    const finalSignal = signal || controller?.signal;
    
    const timeoutId = setTimeout(() => {
      if (controller) {
        controller.abort();
      }
    }, timeout);
    
    try {
      const response = await fetch(`${config.apiBaseUrl}/analyze`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
        signal: finalSignal,
      });
      
      clearTimeout(timeoutId);
      
      if (!response.ok) {
        throw new Error(`API Error: ${response.status} ${response.statusText}`);
      }
      
      return response.json();
    } catch (error) {
      clearTimeout(timeoutId);
      throw error;
    }
  },

  // Validate a single CIDR
  async validateCIDR(cidr, options = {}) {
    const config = await getConfig();
    const { signal, timeout = config.apiTimeout } = options;
    const controller = signal ? null : new AbortController();
    const finalSignal = signal || controller?.signal;
    
    const timeoutId = setTimeout(() => {
      if (controller) {
        controller.abort();
      }
    }, timeout);
    
    try {
      const response = await fetch(`${config.apiBaseUrl}/validate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ cidr }),
        signal: finalSignal,
      });
      
      clearTimeout(timeoutId);
      
      if (!response.ok) {
        throw new Error(`API Error: ${response.statusText}`);
      }
      
      return response.json();
    } catch (error) {
      clearTimeout(timeoutId);
      throw error;
    }
  },

  // Health check
  async healthCheck(options = {}) {
    const config = await getConfig();
    const { signal, timeout = config.apiTimeout } = options;
    const controller = signal ? null : new AbortController();
    const finalSignal = signal || controller?.signal;
    
    const timeoutId = setTimeout(() => {
      if (controller) {
        controller.abort();
      }
    }, timeout);
    
    try {
      const response = await fetch(`${config.apiBaseUrl}/health`, {
        signal: finalSignal,
      });
      
      clearTimeout(timeoutId);
      
      if (!response.ok) {
        throw new Error(`API Error: ${response.statusText}`);
      }
      
      return response.json();
    } catch (error) {
      clearTimeout(timeoutId);
      throw error;
    }
  },

  // Get the current API base URL (useful for debugging)
  async getBaseUrl() {
    const config = await getConfig();
    return config.apiBaseUrl;
  },

  // Get current configuration
  async getCurrentConfig() {
    return await getConfig();
  }
};