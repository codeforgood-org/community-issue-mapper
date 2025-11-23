// Heatmap functionality for Community Issue Mapper

class HeatmapLayer {
    constructor(map) {
        this.map = map;
        this.heatmapLayer = null;
        this.enabled = false;
    }

    enable(issues) {
        if (!issues || issues.length === 0) {
            return;
        }

        // Convert issues to heatmap data points
        const points = issues.map(issue => {
            // Weight by votes and status
            let intensity = 0.5;
            if (issue.votes > 5) intensity = 0.7;
            if (issue.votes > 10) intensity = 0.9;
            if (issue.status === 'new') intensity += 0.2;

            return {
                lat: issue.latitude,
                lng: issue.longitude,
                intensity: Math.min(intensity, 1.0)
            };
        });

        // Create heatmap layer using canvas overlay
        this.createHeatmap(points);
        this.enabled = true;
    }

    createHeatmap(points) {
        // Remove existing heatmap
        if (this.heatmapLayer) {
            this.map.map.removeLayer(this.heatmapLayer);
        }

        // Create canvas overlay for heatmap
        const HeatmapOverlay = L.Layer.extend({
            onAdd: function(map) {
                const canvas = L.DomUtil.create('canvas', 'leaflet-heatmap-layer');
                const size = map.getSize();
                canvas.width = size.x;
                canvas.height = size.y;
                canvas.style.position = 'absolute';
                canvas.style.pointerEvents = 'none';

                this._canvas = canvas;
                this._ctx = canvas.getContext('2d');
                this._map = map;

                map.getPanes().overlayPane.appendChild(canvas);

                map.on('moveend', this._reset, this);
                map.on('resize', this._resize, this);

                this._reset();
            },

            onRemove: function(map) {
                L.DomUtil.remove(this._canvas);
                map.off('moveend', this._reset, this);
                map.off('resize', this._resize, this);
            },

            _resize: function() {
                const size = this._map.getSize();
                this._canvas.width = size.x;
                this._canvas.height = size.y;
                this._reset();
            },

            _reset: function() {
                const topLeft = this._map.containerPointToLayerPoint([0, 0]);
                L.DomUtil.setPosition(this._canvas, topLeft);
                this._draw();
            },

            _draw: function() {
                const ctx = this._ctx;
                const canvas = this._canvas;

                // Clear canvas
                ctx.clearRect(0, 0, canvas.width, canvas.height);

                // Draw heatmap points
                points.forEach(point => {
                    const pos = this._map.latLngToContainerPoint([point.lat, point.lng]);

                    const gradient = ctx.createRadialGradient(
                        pos.x, pos.y, 0,
                        pos.x, pos.y, 40
                    );

                    const alpha = point.intensity * 0.8;
                    gradient.addColorStop(0, `rgba(255, 0, 0, ${alpha})`);
                    gradient.addColorStop(0.5, `rgba(255, 165, 0, ${alpha * 0.6})`);
                    gradient.addColorStop(1, `rgba(255, 255, 0, 0)`);

                    ctx.fillStyle = gradient;
                    ctx.fillRect(pos.x - 40, pos.y - 40, 80, 80);
                });

                // Apply blur
                ctx.globalCompositeOperation = 'source-over';
            }
        });

        this.heatmapLayer = new HeatmapOverlay();
        this.map.map.addLayer(this.heatmapLayer);
    }

    disable() {
        if (this.heatmapLayer) {
            this.map.map.removeLayer(this.heatmapLayer);
            this.heatmapLayer = null;
        }
        this.enabled = false;
    }

    toggle(issues) {
        if (this.enabled) {
            this.disable();
        } else {
            this.enable(issues);
        }
        return this.enabled;
    }

    isEnabled() {
        return this.enabled;
    }
}
