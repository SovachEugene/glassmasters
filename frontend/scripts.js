document.addEventListener('DOMContentLoaded', () => {
    // Константы
    const COOKIE_EXPIRATION_MINUTES = 10;
    const LOADER_MIN_DISPLAY_TIME_MS = 1000;
    const SWIPER_AUTOPLAY_DELAY = 3000;
    const SWIPER_SPEED = 1000;
    const TRANSLATION_API_URL = 'http://localhost:8080/api/translations';
    const LOCATION_API_URL = 'http://ip-api.com/json/';
    const FORM_API_URL = 'http://localhost:8080/api/clients';


    
    console.log('API URL:', TRANSLATION_API_URL);


    // Установка Cookie
    const setCookie = (name, value, minutes) => {
        const date = new Date();
        date.setTime(date.getTime() + minutes * 60 * 1000);
        document.cookie = `${name}=${value}; expires=${date.toUTCString()}; path=/`;
    };

    // Получение Cookie
    const getCookie = (name) => {
        const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
        return match ? match[2] : null;
    };

    // Обработчик бургер-меню
    const burgerMenu = document.querySelector('.burger-menu');
    const mobileMenu = document.querySelector('.mobile-menu');

    if (burgerMenu && mobileMenu) {
        burgerMenu.addEventListener('click', () => {
            burgerMenu.classList.toggle('active');
            mobileMenu.classList.toggle('active');
        });

        mobileMenu.addEventListener('click', (e) => {
            if (e.target.tagName === 'A') {
                burgerMenu.classList.remove('active');
                mobileMenu.classList.remove('active');
            }
        });
    }

    const langCircles = document.querySelectorAll('.lang-circle');
    const langElements = document.querySelectorAll('[data-lang-key]');
    const form = document.getElementById('contactForm');
    let currentLang = getCookie('selectedLang') || 'uk';
    let locationData = {};

    // Инициализация Swiper
    const swiper = new Swiper('.swiper-container', {
        loop: true,
        pagination: {
            el: '.swiper-pagination',
            clickable: true,
        },
        navigation: {
            nextEl: '.swiper-button-next',
            prevEl: '.swiper-button-prev',
        },
        autoplay: {
            delay: SWIPER_AUTOPLAY_DELAY,
            disableOnInteraction: false,
        },
        speed: SWIPER_SPEED,
    });

    // Функции нового лоадера
    const showLoader = () => {
        const loaderWrapper = document.querySelector('.loader-wrapper');
        if (loaderWrapper) loaderWrapper.style.display = 'flex';
    };

    const hideLoader = () => {
        const loaderWrapper = document.querySelector('.loader-wrapper');
        if (loaderWrapper) {
            setTimeout(() => {
                loaderWrapper.style.display = 'none';
            }, LOADER_MIN_DISPLAY_TIME_MS);
        }
    };

    // Загрузка переводов
    const loadTranslations = async (lang) => {
        showLoader();
        try {
            const response = await fetch(`${TRANSLATION_API_URL}?lang=${lang}`);
            if (!response.ok) {
                throw new Error('Ошибка загрузки переводов');
            }
            const translations = await response.json();
            applyTranslations(translations);
        } catch (error) {
            console.error('Ошибка загрузки переводов:', error);
        } finally {
            hideLoader();
        }
    };

    // Применение переводов
    const applyTranslations = (translations) => {
        langElements.forEach((el) => {
            const key = el.getAttribute('data-lang-key');
            if (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA') {
                el.setAttribute('placeholder', translations[key] || '');
            } else {
                el.textContent = translations[key] || '';
            }
        });
    };

    // Переключение языка
    const switchLanguage = (selectedLang) => {
        if (currentLang !== selectedLang) {
            currentLang = selectedLang;
            setCookie('selectedLang', currentLang, COOKIE_EXPIRATION_MINUTES);

            langCircles.forEach((circle) => circle.classList.remove('active'));
            const selectedCircle = document.querySelector(`[data-lang="${currentLang}"]`);
            if (selectedCircle) {
                selectedCircle.classList.add('active');
            }

            loadTranslations(currentLang);
        }
    };

    langCircles.forEach((circle) => {
        circle.addEventListener('click', () => {
            const selectedLang = circle.getAttribute('data-lang');
            switchLanguage(selectedLang);
        });
    });

    const setDefaultLanguage = () => {
        const defaultCircle = document.querySelector(`[data-lang="${currentLang}"]`);
        if (defaultCircle) {
            defaultCircle.classList.add('active');
        }
        loadTranslations(currentLang);
    };

    document.querySelectorAll('header nav a, .logo-link, .cta-button').forEach((anchor) => {
        anchor.addEventListener('click', (e) => {
            e.preventDefault();
            const targetId = anchor.getAttribute('href');
            const targetElement = document.querySelector(targetId);
            if (targetElement) {
                targetElement.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start',
                });
            }
        });
    });

    setDefaultLanguage();

    const fetchLocationData = async () => {
        try {
            const response = await fetch(LOCATION_API_URL);
            if (!response.ok) throw new Error('Ошибка получения данных о местоположении');
            locationData = await response.json();
        } catch (error) {
            console.error('Ошибка при получении местоположения:', error);
        }
    };

    fetchLocationData();

    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();

            const formData = {
                name: document.getElementById('name').value.trim(),
                email: document.getElementById('email').value.trim(),
                phone: document.getElementById('phone').value.trim(),
                message: document.getElementById('message').value.trim(),
            };

            if (!formData.name || !formData.email || !formData.phone || !formData.message) {
                alert('Пожалуйста, заполните все поля формы.');
                return;
            }

            showLoader();

            try {
                const response = await fetch(FORM_API_URL, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData),
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    console.error('Ошибка данных формы:', errorData);
                    alert('Ошибка отправки данных');
                } else {
                    alert('Данные успешно отправлены!');
                    form.reset();
                }
            } catch (error) {
                console.error('Ошибка при отправке данных:', error);
                alert('Ошибка отправки данных');
            } finally {
                hideLoader();
            }
        });
    }
});
