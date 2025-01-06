document.addEventListener("DOMContentLoaded", () => {
    const tableBody = document.querySelector("#clientsTable tbody");

    // Функция для загрузки данных клиентов
    async function fetchClients() {
        try {
            const response = await fetch("http://localhost:8082/api/clients");
            if (!response.ok) {
                throw new Error("Ошибка при загрузке данных клиентов");
            }

            const clients = await response.json();
            renderClients(clients);
        } catch (error) {
            console.error("Ошибка:", error);
        }
    }

    // Функция для отображения данных клиентов в таблице
    function renderClients(clients) {
        tableBody.innerHTML = ""; // Очищаем таблицу перед добавлением новых данных
        clients.forEach((client) => {
            const row = document.createElement("tr");

            row.innerHTML = `
                <td>${client.id}</td>
                <td>${client.name}</td>
                <td>${client.email}</td>
                <td>${client.phone}</td>
                <td>${client.message}</td>
                <td>${client.city || "—"}</td>
                <td>${client.country || "—"}</td>
                <td>
                    ${client.processed ? "Да" : 
                    `<button class="process-button" data-id="${client.id}">Обработано</button>`}
                </td>
            `;

            tableBody.appendChild(row);
        });

        attachEventListeners();
    }

    // Функция для добавления обработчиков событий на кнопки
    function attachEventListeners() {
        const processButtons = document.querySelectorAll(".process-button");
        processButtons.forEach((button) => {
            button.addEventListener("click", async (event) => {
                const clientId = event.target.dataset.id;

                try {
                    const response = await fetch(`http://localhost:8082/api/clients/${clientId}/processed`, {
                        method: "POST",
                    });

                    if (!response.ok) {
                        throw new Error("Ошибка при обработке клиента");
                    }

                    // Обновляем список клиентов после успешного выполнения
                    await fetchClients();
                } catch (error) {
                    console.error("Ошибка:", error);
                }
            });
        });
    }

    // Загрузка данных при загрузке страницы
    fetchClients();
});
