<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
    import { onMount } from 'svelte';

	let { children } = $props();
    let isDark = $state(false);
    let useGif = $state(true);

    const LOBSTER_IMG = "/imgs/lobster.png";
    const LOBSTER_GIF = "/imgs/lober.gif";

    onMount(() => {
        isDark = localStorage.getItem('theme') === 'dark';
        if (isDark) {
            document.documentElement.classList.add('dark');
        }
        useGif = localStorage.getItem('useGif') !== 'false';
    });

    function toggleTheme() {
        isDark = !isDark;
        if (isDark) {
            document.documentElement.classList.add('dark');
            localStorage.setItem('theme', 'dark');
        } else {
            document.documentElement.classList.remove('dark');
            localStorage.setItem('theme', 'light');
        }
    }

    function toggleGif() {
        useGif = !useGif;
        localStorage.setItem('useGif', useGif.toString());
    }
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
    <link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:wght@400;600&family=IBM+Plex+Mono:wght@400;600&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">
</svelte:head>

<div class="app-container">
    <nav>
        <div class="nav-content">
            <div class="brand">
                <h1>Lober</h1>
                <img src={LOBSTER_IMG} alt="lobster" class="lobster-img">
            </div>
            <ul>
                <li><a href="/"><i class="fa-solid fa-house"></i> Home</a></li>
                <li><a href="/metrics"><i class="fa-solid fa-chart-line"></i> Metrics</a></li>
                <li><a href="/users"><i class="fa-solid fa-users"></i> Users</a></li>
                <li><a href="/teams"><i class="fa-solid fa-people-group"></i> Teams</a></li>
                <li><a href="/c2s"><i class="fa-solid fa-server"></i> C2s</a></li>
                <li><a href="/scope"><i class="fa-solid fa-bullseye"></i> Scope</a></li>
                <li><a href="/tokens"><i class="fa-solid fa-key"></i> Tokens</a></li>
                <li><a href="/config"><i class="fa-solid fa-gear"></i> Config</a></li>
            </ul>
        </div>
        <div class="nav-footer">
            <button onclick={toggleTheme} class="theme-toggle">
                <i class="fa-solid {isDark ? 'fa-sun' : 'fa-moon'}"></i>
                {isDark ? 'Light' : 'Dark'} Mode
            </button>
            <button class="img-toggle" onclick={toggleGif} title=":3">
                <img src={useGif ? LOBSTER_GIF : LOBSTER_IMG} alt="lober" class="lober-footer-img">
            </button>
        </div>
    </nav>

    <main>
        {@render children()}
    </main>
</div>

<style>
	:global(body) {
		font-family: "IBM Plex Sans", sans-serif;
        margin: 0;
        padding: 0;
        height: 100vh;
        overflow: hidden;
	}

    :global(html.dark body) {
        background-color: #1a1a1a;
        color: #e0e0e0;
    }

    .app-container {
        display: flex;
        height: 100vh;
        width: 100vw;
    }

	nav {
        width: 180px;
        flex-shrink: 0;
        border-right: 1px solid #ddd;
        padding: 1.25rem;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        background: #fcfcfc;
    }

    :global(html.dark) nav {
        border-right-color: #333;
        background: #1e1e1e;
    }

    .brand {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 2rem;
    }

    .lobster-img {
        height: 28px;
        width: auto;
    }

    nav h1 {
        margin: 0;
        font-weight: 600;
        letter-spacing: -0.5px;
        font-size: 1.5rem;
    }

	nav ul {
        list-style-type: none;
        padding: 0;
        margin: 0;
    }

    nav li a {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 0.6rem 0.75rem;
        text-decoration: none;
        color: inherit;
        border-radius: 6px;
        margin-bottom: 0.25rem;
        font-size: 0.9rem;
    }

    nav li a i {
        width: 16px;
        text-align: center;
        color: #666;
    }

    :global(html.dark) nav li a i {
        color: #aaa;
    }

    nav li a:hover {
        background-color: #f0f0f0;
    }

    :global(html.dark nav li a:hover) {
        background-color: #2d2d2d;
    }

    main {
        flex-grow: 1;
        display: flex;
        flex-direction: column;
        padding: 2rem;
        overflow: hidden;
        width: 0;
    }

    .theme-toggle {
        padding: 0.5rem;
        cursor: pointer;
        width: 100%;
        border: 1px solid #ddd;
        background: transparent;
        color: inherit;
        border-radius: 4px;
        margin-bottom: 1rem;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        font-size: 0.85rem;
    }

    :global(html.dark) .theme-toggle {
        border-color: #444;
    }

    .img-toggle {
        background: none;
        border: none;
        padding: 0;
        cursor: pointer;
        width: 100%;
        display: flex;
        justify-content: center;
    }

    .lober-footer-img {
        max-width: 100%;
        border-radius: 4px;
        object-fit: contain;
    }

    /* Minimal Global Components */
    :global(h1) {
        margin: 0;
        font-size: 1.75rem;
    }

    :global(.header-container) {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.5rem;
        gap: 2rem;
    }

    :global(.header-actions) {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.5rem;
        gap: 2rem;
    }

    :global(.table-container) {
        flex-grow: 1;
        overflow: auto;
        border: 1px solid #ddd;
        border-radius: 4px;
        background: #fff;
        width: 100%;
    }

    :global(html.dark .table-container) {
        border-color: #333;
        background: #252525;
    }

    :global(table) {
        font-family: "IBM Plex Mono", monospace;
        width: 100%;
        border-collapse: collapse;
        font-size: 0.9rem;
        table-layout: fixed;
    }

    :global(th) {
        position: sticky;
        top: 0;
        background-color: #f8f9fa;
        z-index: 10;
        border-bottom: 2px solid #eee;
        text-align: left;
        padding: 12px;
        font-weight: 600;
        color: #666;
    }

    :global(html.dark th) {
        background-color: #2d2d2d;
        border-bottom-color: #444;
        color: #aaa;
    }

    :global(td) {
        padding: 10px 12px;
        border-bottom: 1px solid #eee;
        overflow: hidden;
    }

    :global(html.dark td) {
        border-bottom-color: #333;
    }

    :global(.cell-scroll) {
        overflow-x: auto;
        white-space: nowrap;
        scrollbar-width: thin;
    }

    :global(.cell-scroll::-webkit-scrollbar) {
        height: 4px;
    }

    :global(.cell-scroll::-webkit-scrollbar-thumb) {
        background: #ddd;
        border-radius: 2px;
    }

    :global(html.dark .cell-scroll::-webkit-scrollbar-thumb) {
        background: #444;
    }

    :global(a) {
        color: inherit;
        text-decoration: none;
    }

    :global(a:hover) {
        text-decoration: underline;
    }

    :global(button) {
        padding: 0.5rem 1rem;
        border-radius: 4px;
        border: 1px solid #ddd;
        background: #fff;
        cursor: pointer;
        font-family: inherit;
    }

    :global(html.dark button) {
        background: #333;
        border-color: #444;
        color: #fff;
    }

    :global(button:hover) {
        background: #f0f0f0;
    }

    :global(html.dark button:hover) {
        background: #444;
    }

    :global(button.danger) {
        color: #dc3545;
        border-color: #f5c2c7;
    }

    :global(html.dark button.danger) {
        color: #ea868f;
        border-color: #842029;
    }

    :global(button.danger:hover) {
        background: #f8d7da;
    }

    :global(html.dark button.danger:hover) {
        background: #41141a;
    }
</style>
