import { defineConfig } from 'vite';

export default defineConfig({
    server: {
        port: 3000, // Задайте нужный порт
    },
    define: {
        'process.env': process.env, // Импортируйте переменные окружения
    },
});
