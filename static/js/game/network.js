import { socket_events } from "./socket_events.js";

export const network = {
    socket: null,
    reconnectInterval: 1000,      // Начальная задержка (1 секунда)
    maxReconnectInterval: 30000, // Максимальная задержка (30 секунд)
    reconnectTimer: null,
    isManuallyClosed: false,     // Флаг намеренного закрытия (Logout)
    isDisplaced: false,          // Флаг вытеснения новой вкладкой (защита от петли)
    isConnecting: false,         // Флаг процесса подключения

    connect() {
        // Если мы уже подключаемся или сокет открыт — выходим
        if (this.isConnecting || (this.socket && this.socket.readyState === WebSocket.OPEN)) {
            return;
        }

        this.isConnecting = true;
        this.isManuallyClosed = false;
        this.isDisplaced = false; // Сбрасываем флаг при новой попытке подключения

        // Собираем URL, используя твой basePath
        const protocol = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}${basePath}/ws`;

        console.log("[Network] Попытка подключения к:", protocol);
        this.socket = new WebSocket(protocol);

        this.socket.onopen = () => {
            console.log("[Network] Сетевой узел активен.");
            this.isConnecting = false;
            this.reconnectInterval = 1000; // Сбрасываем таймер задержки при успехе
        };

        this.socket.onmessage = (event) => {
            // Если эта вкладка помечена как вытесненная, игнорируем любые сообщения
            if (this.isDisplaced) return;

            try {
                const msg = JSON.parse(event.data);

                // ОБРАБОТКА ВЫТЕСНЕНИЯ (Событие от сервера)
                if (msg.type === "connection_replaced") {
                    console.error("[Network] Сессия прервана: открыта новая вкладка.");
                    this.isDisplaced = true;
                    this.socket.close(); // Физически закрываем текущую трубу                    
                    // Показываем оверлей блокировки (нужно добавить в HTML/Templates)                   
                    //alert("Вкладка неактивна: обнаружено новое подключение к аккаунту.");
                    if (typeof modalManager !== 'undefined') {
                        modalManager.hide();
                    }

                    // 2. ОСТАНАВЛИВАЕМ ОСНОВНОЙ ЦИКЛ (на случай, если модалки не было)
                    if (typeof engine !== 'undefined') {
                        engine.stopMainLoop();
                        engine.stopBattleLoop();
                    }
                    return;
                }

                // Передаем сообщение в диспетчер событий
                if (socket_events[msg.type]) {
                    socket_events[msg.type](msg);
                }
            } catch (e) {
                console.error("[Network] Ошибка парсинга сообщения:", e);
            }
        };

        this.socket.onclose = (event) => {
            this.isConnecting = false;

            // Если нас вытеснил сервер — реконнект ЗАПРЕЩЕН
            if (this.isDisplaced) {
                console.warn("[Network] Соединение закрыто из-за вытеснения другой вкладкой. Авто-реконнект отключен.");
                return;
            }

            console.warn(`[Network] Связь с сервером потеряна (Код: ${event.code})`);

            // Реконнект только если закрытие было аварийным
            if (!this.isManuallyClosed) {
                this.scheduleReconnect();
            }
        };

        this.socket.onerror = (error) => {
            this.isConnecting = false;
            console.error("[Network] Ошибка сокета:", error);
        };
    },

    scheduleReconnect() {
        if (this.reconnectTimer || this.isDisplaced) return;

        console.log(`[Network] Переподключение через ${this.reconnectInterval}ms...`);

        this.reconnectTimer = setTimeout(() => {
            this.reconnectTimer = null;
            this.connect();
        }, this.reconnectInterval);

        // Экспоненциальный Backoff с добавлением случайности (jitter)
        this.reconnectInterval = Math.min(
            this.maxReconnectInterval,
            this.reconnectInterval * 1.5 + Math.random() * 1000
        );
    },

    send(data = {}) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(data));
        } else {
            console.warn("[Network] Попытка отправки сообщения при закрытом сокете.");
        }
    },

    /**
     * Метод для внешней проверки и восстановления связи
     */
    checkConnection() {
        // Если вытеснены или закрыли сами — не восстанавливаем автоматически
        if (this.isManuallyClosed || this.isDisplaced) return;

        if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
            console.log("[Network] Проверка связи: требуется восстановление.");
            this.connect();
        }
    },

    disconnect() {
        this.isManuallyClosed = true;
        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
        }
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }
};

// --- Глобальные слушатели состояния среды ---

// Обработка пробуждения смартфона и активации вкладки
document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "visible") {
        network.checkConnection();
    }
});

// Обработка восстановления системного интернета
window.addEventListener("online", () => {
    network.checkConnection();
});