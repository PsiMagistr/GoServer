import { socket_events } from "./socket_events.js";
import { criticalController } from "./modal_controllers/criticalController.js";

// Убедитесь, что basePath определен в вашем проекте (например, импортирован из конфига)
const basePath = window.basePath || "";

export const network = {
    socket: null,
    reconnectInterval: 1000,       // Начальная задержка (1 секунда)
    maxReconnectInterval: 30000,   // Максимальная задержка (30 секунд)
    reconnectTimer: null,
    isManuallyClosed: false,       // Флаг намеренного закрытия (Logout)
    isDisplaced: false,            // Флаг вытеснения новой вкладкой
    isConnecting: false,           // Флаг процесса подключения

    connect() {
        // Если уже в процессе подключения — выходим
        if (this.isConnecting) return;

        // Если сокет уже полностью открыт, повторно не коннектимся
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            return;
        }

        this.isConnecting = true;
        this.isManuallyClosed = false;
        this.isDisplaced = false;

        const protocol = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}${basePath}/ws`;

        console.log("[Network] Попытка подключения к:", protocol);
        
        try {
            this.socket = new WebSocket(protocol);
        } catch (e) {
            this.isConnecting = false;
            console.error("[Network] Ошибка создания WebSocket:", e);
            this.scheduleReconnect();
            return;
        }

        this.socket.onopen = () => {
            console.log("[Network] Сетевой узел активен.");
            this.isConnecting = false;
            this.reconnectInterval = 1000; // Сбрасываем таймер задержки при успехе
        };

        this.socket.onmessage = (event) => {
            if (this.isDisplaced) return;

            try {
                const msg = JSON.parse(event.data);

                // ОБРАБОТКА ВЫТЕСНЕНИЯ (Событие от сервера)
                if (msg.type === "connection_replaced") {
                    console.error("[Network] Сессия прервана: открыта новая вкладка.");
                    this.isDisplaced = true;
                    this.socket.close();
                    
                    if (typeof modalManager !== 'undefined') {
                        modalManager.hide();
                    }
                    if (typeof engine !== 'undefined') {
                        engine.stopMainLoop();
                        engine.stopBattleLoop();
                    }
                    criticalController.open({});
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

            if (this.isDisplaced) {
                console.warn("[Network] Соединение закрыто из-за вытеснения другой вкладкой. Авто-реконнект отключен.");
                return;
            }

            console.warn(`[Network] Связь с сервером потеряна (Код: ${event.code})`);

            if (!this.isManuallyClosed) {
                this.scheduleReconnect();
            }
        };

        this.socket.onerror = (error) => {
            this.isConnecting = false;
            console.error("[Network] Ошибка сокета:", error);
            // Событие onerror всегда сопровождается вызовом onclose, поэтому реконнект сработает там
        };
    },

    scheduleReconnect() {
        if (this.reconnectTimer || this.isDisplaced || this.isManuallyClosed) return;

        console.log(`[Network] Переподключение через ${Math.round(this.reconnectInterval)}ms...`);

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
            this.socket.send(typeof data === "string" ? data : JSON.stringify(data));
        } else {
            console.warn("[Network] Попытка отправки сообщения при закрытом сокете.");
        }
    },

    /**
     * Метод для внешней проверки и восстановления связи (при разблокировке экрана / возврате сети)
     */
    checkConnection() {
        if (this.isManuallyClosed || this.isDisplaced) return;

        // Если сокет не открыт или завис в состоянии CONNECTING — принудительно сбрасываем и переподключаем
        if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
            console.log("[Network] Проверка связи: сокет неактивен, принудительное восстановление.");
            if (this.socket) {
                try { this.socket.close(); } catch (e) {}
            }
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

// Обработка восстановления системного интернета (Wi-Fi / 4G)
window.addEventListener("online", () => {
    network.checkConnection();
});