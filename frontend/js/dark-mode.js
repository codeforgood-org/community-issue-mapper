// Dark Mode Toggle

class DarkModeManager {
    constructor() {
        this.isDark = this.loadPreference();
        this.init();
    }

    init() {
        // Apply saved preference
        if (this.isDark) {
            this.enable();
        }

        // Create toggle button
        this.createToggleButton();

        // Listen for system preference changes
        if (window.matchMedia) {
            window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', e => {
                if (!localStorage.getItem('dark-mode-preference')) {
                    this.toggle();
                }
            });
        }
    }

    createToggleButton() {
        const button = document.createElement('button');
        button.className = 'dark-mode-toggle';
        button.setAttribute('aria-label', 'Toggle dark mode');
        button.innerHTML = this.isDark ? '<i class="fas fa-sun"></i>' : '<i class="fas fa-moon"></i>';

        button.addEventListener('click', () => this.toggle());

        document.body.appendChild(button);
        this.toggleButton = button;
    }

    enable() {
        document.body.classList.add('dark-mode');
        this.isDark = true;
        this.updateToggleButton();
        this.savePreference();
    }

    disable() {
        document.body.classList.remove('dark-mode');
        this.isDark = false;
        this.updateToggleButton();
        this.savePreference();
    }

    toggle() {
        if (this.isDark) {
            this.disable();
        } else {
            this.enable();
        }
    }

    updateToggleButton() {
        if (this.toggleButton) {
            this.toggleButton.innerHTML = this.isDark
                ? '<i class="fas fa-sun"></i>'
                : '<i class="fas fa-moon"></i>';
        }
    }

    savePreference() {
        localStorage.setItem('dark-mode-preference', this.isDark ? 'dark' : 'light');
    }

    loadPreference() {
        const saved = localStorage.getItem('dark-mode-preference');
        if (saved) {
            return saved === 'dark';
        }

        // Check system preference
        if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
            return true;
        }

        return false;
    }
}

// Initialize dark mode on page load
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        new DarkModeManager();
    });
} else {
    new DarkModeManager();
}
