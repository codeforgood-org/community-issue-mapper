// Main Application Logic

let issueMap;
let currentFilters = {};

// Initialize app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    initializeApp();
});

async function initializeApp() {
    // Initialize map
    issueMap = new IssueMap('map');
    issueMap.initialize();

    // Load initial data
    await loadIssues();
    await loadStats();

    // Setup event listeners
    setupEventListeners();

    // Try to get user's location
    tryAutoLocate();
}

function setupEventListeners() {
    // Form submission
    document.getElementById('issue-form').addEventListener('submit', handleIssueSubmit);

    // Location selected
    document.addEventListener('locationSelected', handleLocationSelected);

    // Filter controls
    document.getElementById('apply-filters').addEventListener('click', applyFilters);
    document.getElementById('reset-filters').addEventListener('click', resetFilters);

    // Map controls
    document.getElementById('refresh-btn').addEventListener('click', () => {
        loadIssues();
        loadStats();
        showToast('Data refreshed', 'success');
    });

    document.getElementById('locate-btn').addEventListener('click', () => {
        issueMap.locateUser();
    });
}

function handleLocationSelected(e) {
    const { latitude, longitude } = e.detail;

    // Show selected location
    const locationElement = document.getElementById('selected-location');
    const locationText = document.getElementById('location-text');
    locationElement.style.display = 'block';
    locationText.textContent = `Location: ${latitude.toFixed(6)}, ${longitude.toFixed(6)}`;

    // Scroll to form
    document.getElementById('issue-form').scrollIntoView({ behavior: 'smooth', block: 'start' });
}

async function handleIssueSubmit(e) {
    e.preventDefault();

    if (!issueMap.selectedLocation) {
        showToast('Please select a location on the map', 'error');
        return;
    }

    const formData = new FormData(e.target);
    const imageFile = formData.get('image');

    // Prepare issue data
    const issueData = {
        title: formData.get('title'),
        description: formData.get('description'),
        category: formData.get('category'),
        latitude: issueMap.selectedLocation.latitude,
        longitude: issueMap.selectedLocation.longitude,
        reporter_name: formData.get('reporter_name') || undefined,
        reporter_email: formData.get('reporter_email') || undefined,
    };

    try {
        showLoading(true);

        // Create issue
        const issue = await api.createIssue(issueData);

        // Upload image if provided
        if (imageFile && imageFile.size > 0) {
            await api.uploadImage(issue.id, imageFile);
        }

        // Reset form
        e.target.reset();
        issueMap.clearTempMarker();
        document.getElementById('selected-location').style.display = 'none';

        // Reload data
        await loadIssues();
        await loadStats();

        showToast('Issue reported successfully!', 'success');

        // Scroll to the new marker
        issueMap.map.setView([issue.latitude, issue.longitude], 16);

    } catch (error) {
        showToast(`Failed to report issue: ${error.message}`, 'error');
    } finally {
        showLoading(false);
    }
}

async function loadIssues(filters = {}) {
    try {
        showLoading(true);
        const issues = await api.getIssues(filters);
        issueMap.displayIssues(issues || []);
    } catch (error) {
        showToast(`Failed to load issues: ${error.message}`, 'error');
    } finally {
        showLoading(false);
    }
}

async function loadStats() {
    try {
        const stats = await api.getStats();

        document.getElementById('total-issues').textContent = stats.total_issues || 0;
        document.getElementById('new-issues').textContent = stats.issues_by_status.new || 0;
        document.getElementById('resolved-issues').textContent =
            (stats.issues_by_status.resolved || 0) + (stats.issues_by_status.closed || 0);
    } catch (error) {
        console.error('Failed to load stats:', error);
    }
}

async function applyFilters() {
    const category = document.getElementById('filter-category').value;
    const status = document.getElementById('filter-status').value;

    currentFilters = {
        category: category || undefined,
        status: status || undefined,
    };

    await loadIssues(currentFilters);
    showToast('Filters applied', 'info');
}

async function resetFilters() {
    document.getElementById('filter-category').value = '';
    document.getElementById('filter-status').value = '';
    currentFilters = {};
    await loadIssues();
    showToast('Filters reset', 'info');
}

async function handleVote(issueId) {
    try {
        showLoading(true);
        const updatedIssue = await api.voteIssue(issueId);
        issueMap.updateIssueMarker(updatedIssue);
        await loadStats();
        showToast('Vote recorded!', 'success');
    } catch (error) {
        showToast(`Failed to vote: ${error.message}`, 'error');
    } finally {
        showLoading(false);
    }
}

function tryAutoLocate() {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(
            (position) => {
                const { latitude, longitude } = position.coords;
                issueMap.map.setView([latitude, longitude], 13);
            },
            (error) => {
                console.log('Auto-location failed, using default location');
            },
            { timeout: 5000 }
        );
    }
}

// Utility Functions

function showLoading(show) {
    const overlay = document.getElementById('loading-overlay');
    overlay.style.display = show ? 'flex' : 'none';
}

function showToast(message, type = 'info') {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;

    const icons = {
        success: 'fa-check-circle',
        error: 'fa-exclamation-circle',
        info: 'fa-info-circle',
    };

    toast.innerHTML = `
        <i class="fas ${icons[type]}"></i>
        <span class="toast-message">${message}</span>
    `;

    container.appendChild(toast);

    // Auto remove after 3 seconds
    setTimeout(() => {
        toast.style.animation = 'slideIn 0.3s ease-out reverse';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

// Make handleVote globally accessible
window.handleVote = handleVote;
