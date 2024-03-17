const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = {
  devServer: {
    port: 3000, // Укажите порт, если он отличается
    proxy: {
      '/v1': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
    },
  },
};
