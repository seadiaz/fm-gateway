const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function (app) {
    // Get the backend URL from environment variable, default to localhost:3000
    const backendUrl = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3000';

    console.log('Setting up proxy to:', backendUrl);

    app.use(
        '/api',
        createProxyMiddleware({
            target: backendUrl,
            changeOrigin: true,
            pathRewrite: {
                '^/api': '', // remove /api prefix when forwarding to backend
            },
        })
    );

    // Also proxy direct API calls (for backward compatibility)
    app.use(
        ['/companies', '/cafs', '/stamps', '/healthz'],
        createProxyMiddleware({
            target: backendUrl,
            changeOrigin: true,
        })
    );
}; 