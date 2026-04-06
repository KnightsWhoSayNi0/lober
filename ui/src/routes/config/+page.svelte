<script lang="ts">
    import { getContext } from 'svelte';
    const getMasterToken = getContext<() => string>('masterToken');

    async function downloadCSV() {
        const response = await fetch('/api/export/events', {
            headers: {
                'Authorization': `Bearer ${getMasterToken()}`
            }
        });
        
        if (response.ok) {
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'lober_events.csv';
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
        } else {
            alert('Failed to download CSV');
        }
    }
</script>

<title>Lober Config</title>

<h1>Config</h1>
<p>Beep boop meep morp. I am the lober</p>

<div class="config-section">
    <h2>Data Export</h2>
    <button onclick={downloadCSV}>
        <i class="fa-solid fa-file-csv"></i> Download Events CSV
    </button>
</div>

<style>
    .config-section {
        background: #fff;
        border: 1px solid #ddd;
        border-radius: 8px;
        padding: 1.5rem;
        margin-top: 2rem;
        max-width: 600px;
    }

    :global(html.dark) .config-section {
        background: #252525;
        border-color: #333;
    }

    h2 {
        margin-top: 0;
        font-size: 1.25rem;
        margin-bottom: 1rem;
    }

    p {
        color: #666;
        margin-bottom: 1.5rem;
    }

    :global(html.dark) p {
        color: #aaa;
    }

    button {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-weight: 600;
        background: #0066cc;
        color: white;
        border: none;
        padding: 0.75rem 1.25rem;
    }

    button:hover {
        background: #0052a3;
    }
</style>
