<script lang="ts">
    import { onMount } from 'svelte';
    import { Chart, registerables } from 'chart.js';
    Chart.register(...registerables);

    let metrics = $state({
        timeline: [] as Array<{ time: string; count: number }>,
        by_team: [] as Array<{ label: string; count: number }>,
        by_scope: [] as Array<{ label: string; count: number }>,
        by_c2: [] as Array<{ label: string; count: number }>,
        by_user: [] as Array<{ label: string; count: number }>
    });

    let timeRange = $state("24 hours");

    let charts: { [key: string]: Chart | null } = {
        timeline: null,
        team: null,
        scope: null,
        c2: null,
        user: null
    };

    let canvasElements: { [key: string]: HTMLCanvasElement } = $state({});

    async function fetchMetrics() {
        const res = await fetch(`/api/metrics?range=${encodeURIComponent(timeRange)}`);
        if (res.ok) {
            metrics = await res.json();
            updateCharts();
        }
    }

    function createBarChart(canvas: HTMLCanvasElement, data: Array<{label: string, count: number}>, label: string, color: string) {
        return new Chart(canvas, {
            type: 'bar',
            data: {
                labels: data.map(d => d.label),
                datasets: [{
                    label: label,
                    data: data.map(d => d.count),
                    backgroundColor: color,
                    borderRadius: 4
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                indexAxis: 'y',
                scales: {
                    x: { beginAtZero: true, grid: { display: false } },
                    y: { grid: { display: false } }
                },
                plugins: {
                    legend: { display: false }
                }
            }
        });
    }

    function updateCharts() {
        // Destroy existing charts
        Object.values(charts).forEach(c => c?.destroy());

        if (canvasElements.timeline) {
            charts.timeline = new Chart(canvasElements.timeline, {
                type: 'line',
                data: {
                    labels: metrics.timeline.map(t => {
                        const date = new Date(t.time);
                        return timeRange.includes('day') 
                            ? date.toLocaleDateString([], { month: 'short', day: 'numeric' })
                            : date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
                    }),
                    datasets: [{
                        label: 'Events',
                        data: metrics.timeline.map(t => t.count),
                        borderColor: '#0066cc',
                        backgroundColor: 'rgba(0, 102, 204, 0.1)',
                        fill: true,
                        tension: 0.4
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        y: { beginAtZero: true },
                        x: { grid: { display: false } }
                    },
                    plugins: {
                        legend: { display: false }
                    }
                }
            });
        }

        if (canvasElements.team) charts.team = createBarChart(canvasElements.team, metrics.by_team, 'Teams', '#0066cc');
        if (canvasElements.scope) charts.scope = createBarChart(canvasElements.scope, metrics.by_scope, 'Scopes', '#fd7e14');
        if (canvasElements.c2) charts.c2 = createBarChart(canvasElements.c2, metrics.by_c2, 'C2s', '#28a745');
        if (canvasElements.user) charts.user = createBarChart(canvasElements.user, metrics.by_user, 'Users', '#6f42c1');
    }

    onMount(() => {
        fetchMetrics();
    });

    $effect(() => {
        if (timeRange) fetchMetrics();
    });
</script>

<div class="header-container">
    <h1>Metrics</h1>
    <div class="config-group">
        <label>
            Range:
            <select bind:value={timeRange}>
                <option value="1 hour">Last Hour</option>
                <option value="6 hours">Last 6 Hours</option>
                <option value="24 hours">Last 24 Hours</option>
                <option value="7 days">Last 7 Days</option>
                <option value="30 days">Last 30 Days</option>
            </select>
        </label>
    </div>
</div>

<div class="metrics-grid">
    <section class="card full-width">
        <h2>Events Timeline</h2>
        <div class="chart-container timeline">
            <canvas bind:this={canvasElements.timeline}></canvas>
        </div>
    </section>

    <section class="card">
        <h2>Events by Team</h2>
        <div class="chart-container">
            <canvas bind:this={canvasElements.team}></canvas>
        </div>
    </section>

    <section class="card">
        <h2>Events by Scope</h2>
        <div class="chart-container">
            <canvas bind:this={canvasElements.scope}></canvas>
        </div>
    </section>

    <section class="card">
        <h2>Events by C2</h2>
        <div class="chart-container">
            <canvas bind:this={canvasElements.c2}></canvas>
        </div>
    </section>

    <section class="card">
        <h2>Events by User</h2>
        <div class="chart-container">
            <canvas bind:this={canvasElements.user}></canvas>
        </div>
    </section>
</div>

<style>
    .config-group {
        display: flex;
        align-items: center;
        gap: 1rem;
        background: #f0f0f0;
        padding: 0.5rem 1rem;
        border-radius: 6px;
    }

    :global(html.dark) .config-group {
        background: #2d2d2d;
    }

    .config-group label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.9rem;
        font-weight: 600;
    }

    select {
        padding: 0.4rem;
        border-radius: 4px;
        border: 1px solid #ddd;
        background: white;
        color: inherit;
    }

    :global(html.dark) select {
        background: #1a1a1a;
        border-color: #444;
    }

    .metrics-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1.5rem;
        overflow-y: auto;
        padding-bottom: 2rem;
        margin-top: 1rem;
    }

    .full-width {
        grid-column: 1 / -1;
    }

    .card {
        background: #fff;
        border: 1px solid #ddd;
        border-radius: 8px;
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
    }

    :global(html.dark) .card {
        background: #252525;
        border-color: #333;
    }

    h2 {
        margin: 0 0 1rem 0;
        font-size: 0.9rem;
        font-weight: 600;
        color: #666;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    :global(html.dark) h2 {
        color: #aaa;
    }

    .chart-container {
        height: 250px;
        position: relative;
        flex-grow: 1;
    }

    .timeline {
        height: 300px;
    }
</style>
