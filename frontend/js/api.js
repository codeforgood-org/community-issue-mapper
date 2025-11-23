// API Client for Community Issue Mapper

const API_BASE_URL = window.location.origin + '/api/v1';

class APIClient {
    async request(endpoint, options = {}) {
        const url = `${API_BASE_URL}${endpoint}`;
        const config = {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            ...options,
        };

        try {
            const response = await fetch(url, config);

            if (!response.ok) {
                const error = await response.json().catch(() => ({ error: 'An error occurred' }));
                throw new Error(error.error || `HTTP ${response.status}`);
            }

            // Handle 204 No Content
            if (response.status === 204) {
                return null;
            }

            return await response.json();
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    // Issues
    async getIssues(filters = {}) {
        const params = new URLSearchParams();
        Object.keys(filters).forEach(key => {
            if (filters[key] !== '' && filters[key] !== null && filters[key] !== undefined) {
                params.append(key, filters[key]);
            }
        });

        const query = params.toString();
        const endpoint = query ? `/issues?${query}` : '/issues';
        return this.request(endpoint);
    }

    async getIssue(id) {
        return this.request(`/issues/${id}`);
    }

    async createIssue(data) {
        return this.request('/issues', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateIssue(id, data) {
        return this.request(`/issues/${id}`, {
            method: 'PATCH',
            body: JSON.stringify(data),
        });
    }

    async deleteIssue(id) {
        return this.request(`/issues/${id}`, {
            method: 'DELETE',
        });
    }

    async uploadImage(issueId, file) {
        const formData = new FormData();
        formData.append('image', file);

        return this.request(`/issues/${issueId}/upload`, {
            method: 'POST',
            headers: {}, // Let browser set Content-Type with boundary
            body: formData,
        });
    }

    async voteIssue(id) {
        return this.request(`/issues/${id}/vote`, {
            method: 'POST',
        });
    }

    async getStats() {
        return this.request('/stats');
    }
}

const api = new APIClient();
