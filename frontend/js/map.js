// Map Management for Community Issue Mapper

class IssueMap {
    constructor(containerId) {
        this.map = null;
        this.markers = new Map();
        this.markerCluster = null;
        this.selectedLocation = null;
        this.tempMarker = null;
        this.containerId = containerId;
        this.clusteringEnabled = true;
    }

    initialize(lat = 37.7749, lng = -122.4194, zoom = 13) {
        // Initialize map
        this.map = L.map(this.containerId).setView([lat, lng], zoom);

        // Add tile layer
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
            maxZoom: 19,
        }).addTo(this.map);

        // Initialize marker cluster group
        this.markerCluster = L.markerClusterGroup({
            maxClusterRadius: 50,
            spiderfyOnMaxZoom: true,
            showCoverageOnHover: true,
            zoomToBoundsOnClick: true
        });
        this.map.addLayer(this.markerCluster);

        // Add click handler for selecting location
        this.map.on('click', (e) => this.onMapClick(e));

        return this;
    }

    toggleClustering(enabled) {
        this.clusteringEnabled = enabled;
        if (enabled && !this.map.hasLayer(this.markerCluster)) {
            this.map.addLayer(this.markerCluster);
        } else if (!enabled && this.map.hasLayer(this.markerCluster)) {
            this.map.removeLayer(this.markerCluster);
        }
    }

    onMapClick(e) {
        this.selectedLocation = {
            latitude: e.latlng.lat,
            longitude: e.latlng.lng,
        };

        // Remove previous temp marker
        if (this.tempMarker) {
            this.map.removeLayer(this.tempMarker);
        }

        // Add temp marker
        this.tempMarker = L.marker([e.latlng.lat, e.latlng.lng], {
            icon: this.createIcon('blue'),
        }).addTo(this.map);

        this.tempMarker.bindPopup('📍 Selected location for new issue').openPopup();

        // Trigger custom event
        const event = new CustomEvent('locationSelected', {
            detail: this.selectedLocation,
        });
        document.dispatchEvent(event);
    }

    clearTempMarker() {
        if (this.tempMarker) {
            this.map.removeLayer(this.tempMarker);
            this.tempMarker = null;
        }
        this.selectedLocation = null;
    }

    createIcon(color = 'red') {
        const colors = {
            red: '#ef4444',
            orange: '#f59e0b',
            green: '#10b981',
            blue: '#3b82f6',
        };

        return L.divIcon({
            className: 'custom-marker',
            html: `
                <div style="
                    width: 24px;
                    height: 24px;
                    background: ${colors[color] || colors.red};
                    border: 2px solid white;
                    border-radius: 50%;
                    box-shadow: 0 2px 4px rgba(0,0,0,0.3);
                "></div>
            `,
            iconSize: [24, 24],
            iconAnchor: [12, 12],
        });
    }

    getMarkerColor(status) {
        const colorMap = {
            'new': 'red',
            'in_progress': 'orange',
            'resolved': 'green',
            'closed': 'green',
        };
        return colorMap[status] || 'red';
    }

    getCategoryIcon(category) {
        const icons = {
            'pothole': '🕳️',
            'accessibility': '♿',
            'streetlight': '💡',
            'graffiti': '🎨',
            'trash': '🗑️',
            'other': '📋',
        };
        return icons[category] || '📋';
    }

    createPopupContent(issue) {
        const statusClass = `status-${issue.status}`;
        const categoryIcon = this.getCategoryIcon(issue.category);
        const date = new Date(issue.created_at).toLocaleDateString();

        let imageHtml = '';
        if (issue.image_url) {
            imageHtml = `<img src="${issue.image_url}" alt="${issue.title}" class="issue-image">`;
        }

        return `
            <div class="issue-popup">
                <h3>${issue.title}</h3>
                <span class="category-badge">${categoryIcon} ${this.formatCategory(issue.category)}</span>
                <span class="status-badge ${statusClass}">${this.formatStatus(issue.status)}</span>
                <p>${issue.description}</p>
                ${imageHtml}
                ${issue.address ? `<p><i class="fas fa-map-marker-alt"></i> ${issue.address}</p>` : ''}
                <div class="issue-meta">
                    <span><i class="fas fa-calendar"></i> ${date}</span>
                    <span><i class="fas fa-thumbs-up"></i> ${issue.votes} votes</span>
                </div>
                <button class="vote-btn" onclick="handleVote(${issue.id})">
                    <i class="fas fa-thumbs-up"></i> Vote
                </button>
            </div>
        `;
    }

    formatCategory(category) {
        return category.split('_').map(word =>
            word.charAt(0).toUpperCase() + word.slice(1)
        ).join(' ');
    }

    formatStatus(status) {
        return status.split('_').map(word =>
            word.charAt(0).toUpperCase() + word.slice(1)
        ).join(' ');
    }

    addIssueMarker(issue) {
        const marker = L.marker([issue.latitude, issue.longitude], {
            icon: this.createIcon(this.getMarkerColor(issue.status)),
        });

        marker.bindPopup(this.createPopupContent(issue));

        if (this.clusteringEnabled) {
            this.markerCluster.addLayer(marker);
        } else {
            marker.addTo(this.map);
        }

        this.markers.set(issue.id, marker);
        return marker;
    }

    updateIssueMarker(issue) {
        if (this.markers.has(issue.id)) {
            const oldMarker = this.markers.get(issue.id);
            if (this.clusteringEnabled) {
                this.markerCluster.removeLayer(oldMarker);
            } else {
                this.map.removeLayer(oldMarker);
            }
        }
        this.addIssueMarker(issue);
    }

    removeIssueMarker(issueId) {
        if (this.markers.has(issueId)) {
            const marker = this.markers.get(issueId);
            if (this.clusteringEnabled) {
                this.markerCluster.removeLayer(marker);
            } else {
                this.map.removeLayer(marker);
            }
            this.markers.delete(issueId);
        }
    }

    clearMarkers() {
        if (this.clusteringEnabled) {
            this.markerCluster.clearLayers();
        } else {
            this.markers.forEach(marker => this.map.removeLayer(marker));
        }
        this.markers.clear();
    }

    displayIssues(issues) {
        this.clearMarkers();
        issues.forEach(issue => this.addIssueMarker(issue));

        // Fit bounds if there are issues
        if (issues.length > 0) {
            const bounds = L.latLngBounds(
                issues.map(issue => [issue.latitude, issue.longitude])
            );
            this.map.fitBounds(bounds, { padding: [50, 50] });
        }
    }

    exportToGeoJSON() {
        const features = [];
        this.markers.forEach((marker, issueId) => {
            const latlng = marker.getLatLng();
            features.push({
                type: 'Feature',
                geometry: {
                    type: 'Point',
                    coordinates: [latlng.lng, latlng.lat]
                },
                properties: {
                    issueId: issueId
                }
            });
        });

        return {
            type: 'FeatureCollection',
            features: features
        };
    }

    locateUser() {
        if (!navigator.geolocation) {
            showToast('Geolocation is not supported by your browser', 'error');
            return;
        }

        showLoading(true);
        navigator.geolocation.getCurrentPosition(
            (position) => {
                showLoading(false);
                const { latitude, longitude } = position.coords;
                this.map.setView([latitude, longitude], 15);

                L.marker([latitude, longitude], {
                    icon: this.createIcon('blue'),
                })
                    .addTo(this.map)
                    .bindPopup('📍 You are here')
                    .openPopup();

                showToast('Location found!', 'success');
            },
            (error) => {
                showLoading(false);
                showToast('Unable to retrieve your location', 'error');
                console.error('Geolocation error:', error);
            }
        );
    }

    getBounds() {
        const bounds = this.map.getBounds();
        return {
            min_lat: bounds.getSouth(),
            max_lat: bounds.getNorth(),
            min_lng: bounds.getWest(),
            max_lng: bounds.getEast(),
        };
    }
}
