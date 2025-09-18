// Simplified JSON-only configuration system
// Configuration is loaded from /app-config.json at runtime

let configCache = null;
let configPromise = null;

/**
 * Load configuration from JSON file with caching
 * @returns {Promise<Object>} Configuration object
 */
async function loadConfig() {
  // Return cached config if available
  if (configCache) {
    return configCache;
  }

  // Return existing promise if already loading
  if (configPromise) {
    return configPromise;
  }

  // Start loading config
  configPromise = (async () => {
    const defaults = {
      apiBaseUrl: 'http://localhost:8080/api',
      apiTimeout: 30000,
    };

    try {
      console.log('📋 Loading configuration from /app-config.json...');
      const response = await fetch('/app-config.json?' + Date.now()); // Cache busting
      
      if (response.ok) {
        const jsonConfig = await response.json();
        configCache = { ...defaults, ...jsonConfig };
        console.log('✅ Configuration loaded from JSON:', configCache);
      } else {
        console.warn('⚠️ Failed to load app-config.json, using defaults');
        configCache = defaults;
      }
    } catch (error) {
      console.warn('⚠️ Error loading app-config.json:', error.message);
      console.log('📋 Using default configuration');
      configCache = defaults;
    }

    return configCache;
  })();

  return configPromise;
}

/**
 * Get current configuration (async)
 * @returns {Promise<Object>} Configuration object
 */
export async function getConfig() {
  return await loadConfig();
}

/**
 * Get configuration synchronously (returns null if not loaded yet)
 * @returns {Object|null} Configuration object or null
 */
export function getConfigSync() {
  return configCache;
}

/**
 * Clear configuration cache (useful for hot reloading)
 */
export function clearConfigCache() {
  configCache = null;
  configPromise = null;
}

// Immediate loading for better UX
loadConfig();

export default getConfig;