
function copyPlainTextClipboard(val) {
    const textArea = document.createElement('textarea');
    textArea.value = val;
    document.body.appendChild(textArea);
    textArea.select();
    document.execCommand('copy');
    document.body.removeChild(textArea);

}

function singleChat() {
    return {
        hover: false,
        menuOpen: false,
        currentMenuIndex: 0,
        init() {
        },
        forceFocus($el) {
            const id = $el.closest("li").getAttribute("id");
            document.getElementById(id).focus();

        },
        onHoverIn() {
            this.hover = true;
        },
        onHoverOut() {
            this.hover = false;
            this.closeMenu();
        },
        onFocusIn() {
            this.hover = true;
        },
        onFocusOut() {
            this.hover = false;
        },
        toggleMenu() {
            if (this.menuOpen) {
                return this.closeMenu()
            }
            this.$refs.menu.focus()
            this.menuOpen = true

            const firstItem = this.$refs.items.querySelector('.dropdown-menu-item');
            if (firstItem) {
                firstItem.focus();
            }
        },
        closeMenu(focusAfter) {
            if (! this.menuOpen) return
            this.menuOpen = false
            focusAfter && focusAfter.focus()
        },
        menuIndexUp() {
            if (!this.menuOpen) return;
            this.currentIndex = this.currentIndex > 0 ? this.currentIndex - 1 : this.$refs.items.querySelectorAll('.dropdown-link').length - 1;
            this.menuFocusIndex();
        },
        menuIndexDown() {
            if (!this.menuOpen) return;
            const itemCount = this.$refs.items.querySelectorAll('.dropdown-menu-item').length;
            this.currentIndex = (this.currentIndex + 1) % itemCount;
            this.menuFocusIndex();
        },
        menuFocusIndex() {
            const items = this.$refs.items.querySelectorAll('.dropdown-menu-item');
            if (items.length > 0 && items[this.currentIndex]) {
                items[this.currentIndex].focus();
            }
        },
        handleKeydown(event, el) {
            switch (event.code) {
                case 'KeyE':
                    el.dispatchEvent(new Event('keyup-edit'));
                    break;
                case 'KeyC':
                    el.dispatchEvent(new Event('keyup-copy'));
                    break;
                case 'KeyD':
                    el.dispatchEvent(new Event('keyup-download'));
                    break;
                case 'Delete':
                    el.dispatchEvent(new Event('keyup-delete'));
                    break;
            }
        },
    }
}

function drop(e) {
    e.preventDefault();
    const input = document.getElementById('files-input');
    const dataTransfer = new DataTransfer();
    Array.from(e.dataTransfer.files).forEach(file => dataTransfer.items.add(file));
    input.files = dataTransfer.files;

    const event = new Event('change');
    input.dispatchEvent(event);

}


function Editor() {
    return {
        value: '',
        localStorage: false,
        localKey: 'editorText',
        init(autofocus = false, localStorage = false, localKey = 'editorText') {

            if (autofocus) {
                this.$refs.textarea.focus();
                this.$refs.textarea.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
            }

            if (localStorage) {
                this.localStorage = true;
                this.localKey = localKey;
                this.loadStoredText();
            } else {
                this.value = this.$refs.textarea.value;
            }

            setTimeout(() => {
                const event = new Event('input');
                const input = this.$refs.textarea;
                input.dispatchEvent(event);
            }, 10);

            this.handleInputEvent();

        },
        resize() {
            const textarea = this.$refs.textarea;
            const maxHeight = 300;
            this.$refs.textarea.style.height = 'auto';
            this.$refs.textarea.style.height = (textarea.scrollHeight > maxHeight ? maxHeight : textarea.scrollHeight) + 'px';
        },
        handleInputEvent() {
            this.value = this.value.trim() === '' ? '' : this.value;
            if (this.localStorage) {
                localStorage.setItem(this.localKey, this.value);
            }
        },
        loadStoredText() {
            const storedText = localStorage.getItem(this.localKey);
            if (storedText) {
                this.value = storedText;
                this.handleInputEvent();
            }
        },
        scrollToMessage($el) {
            const attribute = $el.getAttribute('data-message-id');
            const element = document.getElementById(`m-${attribute}`);
            if (element) {
                element.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
            }
        },
        fileUpload() {
            const input = document.getElementById('files-input');
            input.click();
        }
    }
}

function chatNavigation() {
    return {
        init() {
            // const firstItem = this.$refs.list.querySelector('.chat-item');
            // if (firstItem) {
            //     firstItem.focus();
            // }
        },

        handleKeydown(event) {
            const focusedItem = document.activeElement;
            let newItem = null;

            if (event.key === 'ArrowDown') {
                newItem = focusedItem.previousElementSibling || focusedItem;
            } else if (event.key === 'ArrowUp') {
                newItem = focusedItem.nextElementSibling || focusedItem;
            }

            if (newItem) {
                newItem.focus();
                newItem.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
            }
        },
    }
}

function chatList() {
    return {
        activeIndex: 0,
        init() {
            this.resize();

            const resizeObserver = new ResizeObserver((entries) => {
                for (let entry of entries) {
                    if (entry.target.id === 'chat-form') {
                        this.resize();
                    }
                }
            });

            resizeObserver.observe(document.getElementById('chat-form'));
            window.addEventListener('resize', this.resize);
        },
        handleKeydown(event) {
            switch (event.code) {
                case 'Escape':
                    document.querySelector('body').dispatchEvent(new Event('keyup-escape'));
                    break;
                case 'Enter':
                    document.querySelector('body').dispatchEvent(new Event('keyup-enter'));
                    break;
            }
        },
        scrollToBottom() {
            setTimeout(() => {
                // document.getElementById('chat-list').scrollTop = document.getElementById('chat-list').scrollHeight;
            }, 100);
        },
        resize() {
            const height = document.getElementById('chat-container-wrap').clientHeight - document.getElementById('chat-form').clientHeight - document.getElementById('chat-title').clientHeight;
            let margin = 40;
            if (window.innerWidth < 768) {
                margin = 12;
            }
            document.getElementById('chat-list').style.height = height - margin + 'px';
        }
    }
}

class Toast {

    constructor(level, message) {
        this.level = level;
        this.message = message;
    }

    show() {
        const now = new Date();
        const timestamp = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`;

        const toast = document.querySelector('#toast-template').cloneNode(true);
        toast.style.display = 'block';
        toast.classList.add("toast", `toast-${this.level}`);
        toast.removeAttribute("id");
        toast.setAttribute("role", "alert");

        toast.querySelector('.toast-title > span:first-child').textContent = timestamp;
        toast.querySelector('.toast-title > strong').textContent = this.level.substring(0, 3);
        toast.querySelector('.toast-title > span:last-child').textContent = this.message;

        document.getElementById("toast-container").style.display = 'block';
        const toastContainerBody = document.getElementById("toast-container-body");
        toastContainerBody.appendChild(toast);

        setTimeout(() => {
            toast.style.display = 'none';
            if (toastContainerBody.children.length === 0) {
                document.getElementById("toast-container").style.display = 'none';
            }
        }, 5000);
    }
}

function onMakeToast(e) {
    const toast = new Toast(e.detail.level, e.detail.message);
    toast.show();
}

function init() {

    if (history.scrollRestoration) {
        history.scrollRestoration = "manual";
    }

    document.body.addEventListener('htmx:configRequest', function(e) {
        e.detail.headers["X-Timezone"] = Intl.DateTimeFormat().resolvedOptions().timeZone
    });

    document.body.addEventListener('MessageEdited', function(e) {
        const prevItems = document.getElementsByClassName("message-edited")
        for (let i = 0; i < prevItems.length; i++) {
            prevItems[i].classList.remove("message-edited");
        }
        const item = document.getElementById("m-"+e.detail.target);
        item.classList.add("message-edited");
    });

    document.body.addEventListener('MessageEditCancelled', function(e) {
        const item = document.getElementById("m-"+e.detail.target);
        item.classList.remove("message-edited");
    });

    document.body.addEventListener('htmx:sseBeforeMessage', function(e) {
        const type = e.detail.type;

        if (type.startsWith("edited:")) {
            const event = new CustomEvent('MessageUpdated' + e.detail.data);
            document.body.dispatchEvent(event);
        }
        if (type.startsWith("deleted:")) {
            const element = document.getElementById("m-" + e.detail.data);
            if (element) {
                element.remove();
            }
        }
        if(type === "presence") {
            const event = new CustomEvent('PresenceUpdate', {detail: e.detail.data});
            document.body.dispatchEvent(event);

            console.log("presence update", e.detail.data);
        }
    });

    document.body.addEventListener("showToast", onMakeToast);

}

init();
