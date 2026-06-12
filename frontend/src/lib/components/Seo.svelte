<script lang="ts">
	import { page } from '$app/state';

	let {
		title,
		description,
		image,
		type = 'website',
		jsonLd
	}: {
		title: string;
		description: string;
		image?: string;
		type?: 'website' | 'product' | 'article';
		jsonLd?: Record<string, unknown>;
	} = $props();

	const canonical = $derived(`${page.url.origin}${page.url.pathname}`);
	const ogImage = $derived(image ?? `${page.url.origin}/og-cover.png`);

	const jsonLdScript = $derived(
		jsonLd ? `<script type="application/ld+json">${JSON.stringify(jsonLd)}<\/script>` : ''
	);
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<link rel="canonical" href={canonical} />

	<meta property="og:type" content={type} />
	<meta property="og:site_name" content="Despacha Aí" />
	<meta property="og:locale" content="pt_AO" />
	<meta property="og:title" content={title} />
	<meta property="og:description" content={description} />
	<meta property="og:url" content={canonical} />
	<meta property="og:image" content={ogImage} />

	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content={title} />
	<meta name="twitter:description" content={description} />
	<meta name="twitter:image" content={ogImage} />

	<meta name="geo.region" content="AO-BGU" />
	<meta name="geo.placename" content="Benguela, Angola" />

	{#if jsonLdScript}
		<!-- eslint-disable-next-line svelte/no-at-html-tags -->
		{@html jsonLdScript}
	{/if}
</svelte:head>
