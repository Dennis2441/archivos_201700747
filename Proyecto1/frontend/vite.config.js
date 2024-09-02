// vite.config.js  
import { defineConfig } from 'vite';  
import react from '@vitejs/plugin-react';  

// Configuración del servidor con proxy  
export default defineConfig({  
  plugins: [react()],  
  server: {  
    proxy: {  
      // Redirigir todas las solicitudes de /api al servidor Go  
      '/api': {  
        target: 'http://localhost:8080', // Cambiar este puerto si tu backend es diferente  
        changeOrigin: true,  
        rewrite: (path) => path.replace(/^\/api/, ''), // Reescribir la URL  
      },  
    },  
  },  
});