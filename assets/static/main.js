const { Marked } = window.marked;
const { markedHighlight } = window.markedHighlight;
const { markedAlert } = window.markedAlert;
const { markedPlaintify } = window.markedPlaintify;

const marked = new Marked(
    markedHighlight({
        langPrefix: 'hljs language-',
        highlight(code, lang, info) {
            const language = hljs.getLanguage(lang) ? lang : 'plaintext';
            return hljs.highlight(code, {language}).value;
        }
    }),
    window.markedAlert(),
);

const plaintextMark = new Marked(
    window.markedPlaintify(),
);

marked.use({
    gfm: true,
    breaks: true,
});

plaintextMark.use({
    gfm: true,
    breaks: true,
});

function renderChat(val) {
    return marked.parse(val);
}

function copyPlainTextClipboard(val) {
    const textArea = document.createElement('textarea');
    textArea.value = plaintextMark.parse(val);
    document.body.appendChild(textArea);
    textArea.select();
    document.execCommand('copy');
    document.body.removeChild(textArea);

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

function markdownEditor(name) {
    return {
        editorValue: '',
        renderedMarkdown: '',
        files: [],
        localKey: name + 'Text',
        dragging: false,
        filesInput: undefined,
        textInput: undefined,
        init() {
            this.loadStoredText();
            this.filesInput = document.getElementById('files-input');
            this.textInput = document.getElementById('message-input');
        },
        handleFileUpload(e) {
            console.log("e.detail.xhr.response", e.detail.xhr.response);
        },
        handleFileDelete(e) {
            const index = this.files.findIndex(file => file.id === e.detail.id);
            this.files.splice(index, 1);
        },
        resize() {
            const textarea = this.textInput;
            const maxHeight = 700;
            textarea.style.height = 'auto';
            textarea.style.height = (textarea.scrollHeight > maxHeight ? maxHeight : textarea.scrollHeight) + 'px';
        },
        handleInputEvent() {
            this.editorValue = this.editorValue.trim() === '' ? '' : this.editorValue;
            localStorage.setItem(this.localKey, this.editorValue);
        },
        loadStoredText() {
            const storedText = localStorage.getItem(this.localKey);
            if (storedText) {
                this.editorValue = storedText;
                this.handleInputEvent();
                setTimeout(() => {
                    resetResizer();
                }, 10);
            }
        },
        handleUploadButtonClick() {
            this.filesInput.click();
        }
    }
}

function Editor() {
    return {
        value: '',
        init(autofocus = false) {
            this.value = this.$refs.textarea.value;
            setTimeout(() => {
                const event = new Event('input');
                const input = this.$refs.textarea;
                input.dispatchEvent(event);
            }, 10);

            if (autofocus) {
                this.$refs.textarea.focus();
            }
        },
        resize() {
            const textarea = this.$refs.textarea;
            const maxHeight = 700;
            this.$refs.textarea.style.height = 'auto';
            this.$refs.textarea.style.height = (textarea.scrollHeight > maxHeight ? maxHeight : textarea.scrollHeight) + 'px';
        },
        handleInputEvent() {
            this.value = this.value.trim() === '' ? '' : this.value;
        },
    }
}

function resetResizer() {
    const event = new Event('input');
    const input = document.getElementById('message-input');
    input.dispatchEvent(event);
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

            if (event.key === 'ArrowUp') {
                newItem = focusedItem.previousElementSibling || focusedItem;
            } else if (event.key === 'ArrowDown') {
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
        sessionDialog: false,
        dialogClicks: 0,
        showProgressBar: false,
        progressInterval: null,
        activeIndex: 0,
        init() {
            this.resize();
            this.scrollToBottom();

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
                document.getElementById('chat-list').scrollTop = document.getElementById('chat-list').scrollHeight;
            }, 100);
        },
        resize() {
            const height = (window.innerHeight - 48) - document.getElementById('chat-form').clientHeight;
            document.getElementById('chat-list').style.height = height + 'px';
        },
        sessionDialogClick() {
            if (this.dialogClicks === 3) {
                this.showProgressBar = false;
                this.progressInterval = null;
                return;
            }
            this.sessionDialog = true;
            this.dialogClicks += 1;

            if (this.dialogClicks === 1) {
                this.initializeProgressBar();
                this.resetClicks();
            }

            if (this.dialogClicks === 3) {
                this.sessionDialog = false;
                document.querySelector('body').dispatchEvent(new Event('shutdown-session'));
            }
        },
        resetClicks() {
            setTimeout(() => {
                if (this.dialogClicks < 3) {
                    this.dialogClicks = 0;
                    this.showProgressBar = false;
                    this.progressInterval = null;
                }
            }, 5000);
        },
        initializeProgressBar() {
            if (this.progressInterval) {
                clearInterval(this.progressInterval);
                this.progressInterval = null;
            }

            this.showProgressBar = true;
            let progress = 0;
            this.$refs.progressBar.style.width = `0%`;

            this.progressInterval = setInterval(() => {
                progress += 2;
                this.$refs.progressBar.style.width = `${progress}%`;
                if (progress >= 100) {
                    clearInterval(this.progressInterval);
                    this.showProgressBar = false;
                    this.$refs.progressBar.style.width = `0%`;
                    this.progressInterval = null;
                }
            }, 100);
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
        toastContainerBody.scrollTop = toastContainerBody.scrollHeight;
    }
}

function onMakeToast(e) {
    const toast = new Toast(e.detail.level, e.detail.message);
    toast.show();
}

function handleHideToast() {
    const closeButton = document.getElementById('toast-container-close');
    closeButton.addEventListener('click', function() {
        const toastContainer = document.getElementById('toast-container');
        toastContainer.style.display = 'none';
    });
}

function init() {
    handleHideToast();

    document.body.addEventListener("showToast", onMakeToast);
}

init();
