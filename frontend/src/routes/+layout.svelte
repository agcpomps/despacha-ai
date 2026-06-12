<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import BottomNav from '$lib/components/BottomNav.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Despacha Aí — Compra e venda em Angola</title>
</svelte:head>

<div class="flex min-h-screen flex-col">
	<Header />
	<main class="mx-auto w-full max-w-7xl flex-1 px-4 pt-6 pb-28 sm:px-6 sm:pb-6">
		{@render children()}
	</main>
	<Footer />
	<BottomNav />
</div>
