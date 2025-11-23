// Admin Dashboard JavaScript

let charts = {};
let currentSection = 'dashboard';

document.addEventListener('DOMContentLoaded', () => {
    initializeAdmin();
    setupEventListeners();
    loadDashboardData();
});

function initializeAdmin() {
    // Check authentication
    const token = localStorage.getItem('auth_token');
    if (!token) {
        window.location.href = 'login.html';
        return;
    }
}

function setupEventListeners() {
    // Navigation
    document.querySelectorAll('.nav-item[data-section]').forEach(item => {
        item.addEventListener('click', (e) => {
            e.preventDefault();
            const section = item.dataset.section;
            switchSection(section);
        });
    });

    // Refresh data
    document.getElementById('refresh-data')?.addEventListener('click', loadDashboardData);

    // Issue filters
    document.getElementById('filter-status-admin')?.addEventListener('change', loadIssuesTable);
    document.getElementById('search-issues')?.addEventListener('input', debounce(loadIssuesTable, 300));

    // Modal close
    document.querySelector('.modal-close')?.addEventListener('click', () => {
        document.getElementById('issue-modal').style.display = 'none';
    });
}

function switchSection(section) {
    // Update nav
    document.querySelectorAll('.nav-item').forEach(item => item.classList.remove('active'));
    document.querySelector(`[data-section="${section}"]`)?.classList.add('active');

    // Update sections
    document.querySelectorAll('.admin-section').forEach(sec => sec.classList.remove('active'));
    document.getElementById(section)?.classList.add('active');

    currentSection = section;

    // Load section data
    switch (section) {
        case 'dashboard':
            loadDashboardData();
            break;
        case 'issues':
            loadIssuesTable();
            break;
        case 'analytics':
            loadAnalytics();
            break;
    }
}

async function loadDashboardData() {
    try {
        showLoading(true);

        // Load analytics
        const analytics = await api.request('/analytics');

        // Update stats
        document.getElementById('total-issues-stat').textContent = analytics.total_issues || 0;
        document.getElementById('resolved-stat').textContent =
            (analytics.issues_by_status.resolved || 0) + (analytics.issues_by_status.closed || 0);
        document.getElementById('pending-stat').textContent =
            (analytics.issues_by_status.new || 0) + (analytics.issues_by_status.in_progress || 0);
        document.getElementById('users-stat').textContent = analytics.total_users || 0;

        // Render charts
        renderTrendChart(analytics.issues_trend);
        renderCategoryChart(analytics.issues_by_category);

        // Render activity
        renderRecentActivity(analytics.recent_activity);

    } catch (error) {
        showToast('Failed to load dashboard data', 'error');
        console.error(error);
    } finally {
        showLoading(false);
    }
}

function renderTrendChart(trendData) {
    const ctx = document.getElementById('issues-trend-chart');
    if (!ctx) return;

    if (charts.trend) {
        charts.trend.destroy();
    }

    const labels = trendData.map(d => d.date);
    const data = trendData.map(d => d.count);

    charts.trend = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Issues',
                data: data,
                borderColor: '#3b82f6',
                backgroundColor: 'rgba(59, 130, 246, 0.1)',
                tension: 0.4,
                fill: true
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    display: false
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        precision: 0
                    }
                }
            }
        }
    });
}

function renderCategoryChart(categoryData) {
    const ctx = document.getElementById('category-chart');
    if (!ctx) return;

    if (charts.category) {
        charts.category.destroy();
    }

    const labels = Object.keys(categoryData);
    const data = Object.values(categoryData);
    const colors = [
        '#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#6366f1'
    ];

    charts.category = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels.map(l => l.charAt(0).toUpperCase() + l.slice(1)),
            datasets: [{
                data: data,
                backgroundColor: colors,
                borderWidth: 0
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    position: 'bottom'
                }
            }
        }
    });
}

function renderRecentActivity(activities) {
    const container = document.getElementById('recent-activity');
    if (!container) return;

    if (!activities || activities.length === 0) {
        container.innerHTML = '<p style="text-align: center; padding: 2rem; color: var(--text-secondary);">No recent activity</p>';
        return;
    }

    container.innerHTML = activities.slice(0, 10).map(activity => {
        const icon = getActivityIcon(activity.action);
        const time = formatRelativeTime(new Date(activity.created_at));

        return `
            <div class="activity-item">
                <div class="activity-icon">
                    <i class="${icon}"></i>
                </div>
                <div class="activity-content">
                    <div class="activity-text">
                        Issue #${activity.issue_id}: ${activity.action}
                    </div>
                    <div class="activity-time">${time}</div>
                </div>
            </div>
        `;
    }).join('');
}

function getActivityIcon(action) {
    const icons = {
        'created': 'fas fa-plus-circle',
        'updated': 'fas fa-edit',
        'commented': 'fas fa-comment',
        'voted': 'fas fa-thumbs-up',
        'resolved': 'fas fa-check-circle'
    };
    return icons[action] || 'fas fa-circle';
}

async function loadIssuesTable() {
    try {
        showLoading(true);

        const status = document.getElementById('filter-status-admin')?.value || '';
        const search = document.getElementById('search-issues')?.value || '';

        const filters = {};
        if (status) filters.status = status;
        if (search) filters.search = search;

        const issues = await api.getIssues(filters);

        const tbody = document.getElementById('issues-table-body');
        if (!tbody) return;

        if (!issues || issues.length === 0) {
            tbody.innerHTML = '<tr><td colspan="8" style="text-align: center; padding: 2rem;">No issues found</td></tr>';
            return;
        }

        tbody.innerHTML = issues.map(issue => {
            const statusClass = `status-${issue.status}`;
            const date = new Date(issue.created_at).toLocaleDateString();

            return `
                <tr>
                    <td>#${issue.id}</td>
                    <td>${escapeHtml(issue.title)}</td>
                    <td>${formatCategory(issue.category)}</td>
                    <td><span class="status-badge ${statusClass}">${formatStatus(issue.status)}</span></td>
                    <td>${escapeHtml(issue.reporter_name || 'Anonymous')}</td>
                    <td>${issue.votes}</td>
                    <td>${date}</td>
                    <td class="table-actions">
                        <button class="btn btn-sm btn-primary" onclick="viewIssue(${issue.id})">
                            <i class="fas fa-eye"></i> View
                        </button>
                        <button class="btn btn-sm btn-secondary" onclick="updateIssueStatus(${issue.id})">
                            <i class="fas fa-edit"></i> Update
                        </button>
                    </td>
                </tr>
            `;
        }).join('');

    } catch (error) {
        showToast('Failed to load issues', 'error');
        console.error(error);
    } finally {
        showLoading(false);
    }
}

async function viewIssue(issueId) {
    try {
        const issue = await api.getIssue(issueId);
        const modal = document.getElementById('issue-modal');
        const detail = document.getElementById('issue-detail');

        detail.innerHTML = `
            <h2>${escapeHtml(issue.title)}</h2>
            <div style="margin: 1rem 0;">
                <span class="category-badge">${formatCategory(issue.category)}</span>
                <span class="status-badge status-${issue.status}">${formatStatus(issue.status)}</span>
            </div>
            <p style="margin: 1rem 0;">${escapeHtml(issue.description)}</p>
            ${issue.image_url ? `<img src="${issue.image_url}" style="max-width: 100%; border-radius: 0.5rem; margin: 1rem 0;" />` : ''}
            <div style="margin-top: 1rem;">
                <p><strong>Reporter:</strong> ${escapeHtml(issue.reporter_name || 'Anonymous')}</p>
                <p><strong>Location:</strong> ${issue.latitude.toFixed(6)}, ${issue.longitude.toFixed(6)}</p>
                <p><strong>Votes:</strong> ${issue.votes}</p>
                <p><strong>Created:</strong> ${new Date(issue.created_at).toLocaleString()}</p>
            </div>
        `;

        modal.style.display = 'flex';
    } catch (error) {
        showToast('Failed to load issue details', 'error');
    }
}

async function updateIssueStatus(issueId) {
    const newStatus = prompt('Enter new status (new, in_progress, resolved, closed):');
    if (!newStatus) return;

    const validStatuses = ['new', 'in_progress', 'resolved', 'closed'];
    if (!validStatuses.includes(newStatus)) {
        showToast('Invalid status', 'error');
        return;
    }

    try {
        await api.updateIssue(issueId, { status: newStatus });
        showToast('Status updated successfully', 'success');
        loadIssuesTable();
    } catch (error) {
        showToast('Failed to update status', 'error');
    }
}

async function loadAnalytics() {
    try {
        showLoading(true);

        const analytics = await api.request('/analytics');

        // Update analytics stats
        document.getElementById('avg-resolution-time').textContent =
            Math.round(analytics.avg_resolution_time_hours || 0);
        document.getElementById('issues-this-week').textContent = analytics.issues_this_week || 0;

        const resolutionRate = analytics.total_issues > 0
            ? Math.round(((analytics.issues_by_status.resolved || 0) + (analytics.issues_by_status.closed || 0)) / analytics.total_issues * 100)
            : 0;
        document.getElementById('resolution-rate').textContent = resolutionRate;

        // Top contributors
        renderTopContributors(analytics.top_contributors);

    } catch (error) {
        showToast('Failed to load analytics', 'error');
        console.error(error);
    } finally {
        showLoading(false);
    }
}

function renderTopContributors(contributors) {
    const container = document.getElementById('top-contributors');
    if (!container) return;

    if (!contributors || contributors.length === 0) {
        container.innerHTML = '<p style="text-align: center; padding: 2rem;">No contributors yet</p>';
        return;
    }

    container.innerHTML = contributors.map((contributor, index) => {
        const totalActivity = contributor.issue_count + contributor.comment_count + contributor.vote_count;

        return `
            <div class="contributor-item">
                <div class="contributor-rank">${index + 1}</div>
                <div class="contributor-info">
                    <div class="contributor-name">${escapeHtml(contributor.user.name)}</div>
                    <div class="contributor-stats">
                        ${contributor.issue_count} issues · ${contributor.comment_count} comments · ${contributor.vote_count} votes
                    </div>
                </div>
                <div style="font-weight: 700; color: var(--primary-color);">
                    ${totalActivity} total
                </div>
            </div>
        `;
    }).join('');
}

// Utility Functions

function formatCategory(category) {
    return category.split('_').map(word =>
        word.charAt(0).toUpperCase() + word.slice(1)
    ).join(' ');
}

function formatStatus(status) {
    return status.split('_').map(word =>
        word.charAt(0).toUpperCase() + word.slice(1)
    ).join(' ');
}

function formatRelativeTime(date) {
    const now = new Date();
    const diff = now - date;
    const seconds = Math.floor(diff / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (days > 0) return `${days} day${days > 1 ? 's' : ''} ago`;
    if (hours > 0) return `${hours} hour${hours > 1 ? 's' : ''} ago`;
    if (minutes > 0) return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
    return 'Just now';
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

function showLoading(show) {
    // Implement loading indicator
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

    setTimeout(() => {
        toast.style.animation = 'slideIn 0.3s ease-out reverse';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

// Make functions globally accessible
window.viewIssue = viewIssue;
window.updateIssueStatus = updateIssueStatus;
