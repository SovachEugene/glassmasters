document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("loginForm");
    const errorMessage = document.getElementById("errorMessage");

    loginForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const username = document.getElementById("username").value.trim();
        const password = document.getElementById("password").value.trim();

        try {
            const response = await fetch("http://localhost:8082/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ username, password }),
            });

            if (!response.ok) {
                const data = await response.json();
                errorMessage.textContent = data.message || "Login failed.";
                errorMessage.classList.add("active");
                return;
            }

            // Успешный вход
            const data = await response.json();
            localStorage.setItem("authToken", data.token); // Сохраняем токен
            window.location.href = "crm.html"; // Переход на главную страницу CRM
        } catch (error) {
            errorMessage.textContent = "Error: Unable to login.";
            errorMessage.classList.add("active");
            console.error(error);
        }
    });
});
