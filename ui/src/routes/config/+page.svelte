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
    async function uploadJSON(event: Event) {
        const input = event.target as HTMLInputElement;
        if (!input.files || input.files.length === 0) return;

        const file = input.files[0];
        const reader = new FileReader();
        reader.onload = async (e) => {
            const content = e.target?.result;
            try {
                const response = await fetch('/api/import/scope', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${getMasterToken()}`
                    },
                    body: content as string
                });

                if (response.ok) {
                    alert('Import successful');
                } else {
                    const err = await response.text();
                    alert(`Import failed: ${err}`);
                }
            } catch (err) {
                alert(`Error: ${err}`);
            }
        };
        reader.readAsText(file);
    }
</script>

<title>Lober Config</title>

<h1>Config</h1>
<p>Beep boop meep morp. I am the lober</p>

<div class="config-section">
    <h2>Data Management</h2>
    <div class="button-group">
        <button onclick={downloadCSV}>
            <i class="fa-solid fa-file-csv"></i> Download Events CSV
        </button>
        
        <label class="file-upload">
            <input type="file" accept=".json" onchange={uploadJSON} style="display: none;" />
            <i class="fa-solid fa-file-import"></i> Import Scope JSON
        </label>
    </div>
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

    .button-group {
        display: flex;
        gap: 1rem;
        flex-wrap: wrap;
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

    button, .file-upload {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-weight: 600;
        background: #0066cc;
        color: white;
        border: none;
        padding: 0.75rem 1.25rem;
        cursor: pointer;
        border-radius: 4px;
    }

    button:hover, .file-upload:hover {
        background: #0052a3;
    }
</style>
