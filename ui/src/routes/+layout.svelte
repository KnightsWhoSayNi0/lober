<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<nav>
	<h1>Lober</h1>
        <ul>
            <li><a class="active" href="/">Home</a></li>
            <li><a href="users">Users</a></li>
            <li><a href="teams">Teams</a></li>
            <li><a href="c2s">C2s</a></li>
            <li><a href="scope">Scope</a></li>
            <li><a href="config">Config</a></li>
        </ul>
        <img src="/imgs/lober.gif" alt="lober">
</nav>

<main>
	{@render children()}
</main>

<style>
	:global(body) {
		display: flex;
		font-family: "IBM Plex Sans", sans-serif;
	}

	nav ul {
        list-style-type: none;
    }

    nav li a {
        display: block;
        color: black;
        padding: 8px 16px;
        text-decoration: none;
    }

    nav li a:hover {
        background-color: #555555;
        color: white;
    }
</style>