import { defineConfig } from 'vite';

export default defineConfig({
  server: {
    host: '0.0.0.0', // Позволяет принимать подключения извне контейнера
    port: 3000, // Указывает порт для разработки
  },
  define: {
    'process.env': {}, // Позволяет использовать process.env для переменных окружения
  },
});
