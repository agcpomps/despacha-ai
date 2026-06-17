<script lang="ts">
	import { env } from '$env/dynamic/public';
	import Seo from '$lib/components/Seo.svelte';
	import Logo from '$lib/components/Logo.svelte';

	const adminWhatsApp = (env.PUBLIC_ADMIN_WHATSAPP ?? '').replace(/\D/g, '');

	const steps = [
		{
			title: 'Publica grátis',
			text: 'Cria o teu anúncio em menos de 2 minutos, com fotos tiradas na hora pelo telemóvel.'
		},
		{
			title: 'Negoceia pelo WhatsApp',
			text: 'Os compradores falam contigo directamente, sem intermediários nem comissões.'
		},
		{
			title: 'Vende perto de ti',
			text: 'Chega a milhares de pessoas em Benguela, Lobito, Catumbela e em toda a Angola.'
		}
	];
</script>

<Seo
	title="Sobre o Despacha Aí — o mercado online de Benguela"
	description="O Despacha Aí é o mercado online de Angola, nascido em Benguela. Compra e vende artigos novos e usados, de graça, e negoceia directamente pelo WhatsApp."
/>

<!-- Hero -->
<section
	class="overflow-hidden rounded-3xl bg-linear-to-br from-brand-900 via-brand-700 to-brand-500 px-6 py-12 text-white sm:px-10 sm:py-16"
>
	<div class="flex items-center gap-3">
		<span class="rounded-2xl bg-white/95 px-3 py-2"><Logo class="text-xl" /></span>
	</div>
	<h1 class="mt-6 max-w-2xl text-3xl font-extrabold tracking-tight sm:text-4xl">
		O mercado online de Angola, nascido em Benguela
	</h1>
	<p class="mt-3 max-w-2xl text-brand-100">
		Criámos o Despacha Aí para facilitar a compra e venda entre angolanos — de forma simples,
		gratuita e segura, com o WhatsApp no centro de tudo.
	</p>
</section>

<!-- Missão -->
<section class="mt-10 grid gap-6 sm:grid-cols-2">
	<div class="rounded-2xl border border-neutral-200 bg-white p-6">
		<h2 class="text-lg font-bold text-neutral-900">A nossa missão</h2>
		<p class="mt-2 text-sm leading-relaxed text-neutral-600">
			Dar a cada angolano uma forma fácil de transformar o que já não usa em dinheiro, e de
			encontrar o que precisa a um bom preço — perto de casa. Começámos por Benguela porque é a nossa
			terra, e estamos a crescer para todo o país.
		</p>
	</div>
	<div class="rounded-2xl border border-neutral-200 bg-white p-6">
		<h2 class="text-lg font-bold text-neutral-900">Porquê grátis?</h2>
		<p class="mt-2 text-sm leading-relaxed text-neutral-600">
			Publicar um anúncio não custa nada e nunca cobramos comissão sobre as tuas vendas. Negoceias
			directamente com o comprador, pelo WhatsApp, sem intermediários.
		</p>
	</div>
</section>

<!-- Como funciona -->
<section class="mt-10">
	<h2 class="text-xl font-bold text-neutral-900">Como funciona</h2>
	<div class="mt-4 grid gap-4 sm:grid-cols-3">
		{#each steps as step, index (step.title)}
			<div class="rounded-2xl border border-neutral-200 bg-white p-6">
				<span
					class="flex h-9 w-9 items-center justify-center rounded-full bg-brand-100 text-sm font-bold text-brand-800"
					>{index + 1}</span
				>
				<h3 class="mt-3 font-bold text-neutral-900">{step.title}</h3>
				<p class="mt-1 text-sm leading-relaxed text-neutral-600">{step.text}</p>
			</div>
		{/each}
	</div>
</section>

<!-- Segurança -->
<section class="mt-10 rounded-2xl bg-amber-50 p-6">
	<h2 class="text-lg font-bold text-amber-900">Dicas de segurança</h2>
	<ul class="mt-3 space-y-2 text-sm text-amber-800">
		<li>• Encontra-te com o vendedor ou comprador num local público e movimentado.</li>
		<li>• Verifica bem o artigo antes de pagar.</li>
		<li>• Desconfia de preços bons demais para ser verdade e de pedidos de pagamento adiantado.</li>
	</ul>
</section>

<!-- CTA -->
<section
	class="mt-10 flex flex-col items-center gap-4 rounded-3xl border border-brand-200 bg-brand-50 p-8 text-center sm:flex-row sm:justify-between sm:text-left"
>
	<div>
		<h2 class="text-lg font-bold text-brand-900">Pronto para começar?</h2>
		<p class="mt-1 text-sm text-brand-800/80">Publica o teu primeiro anúncio hoje, de graça.</p>
	</div>
	<div class="flex flex-wrap items-center justify-center gap-3">
		<a
			href="/anunciar"
			class="rounded-full bg-brand-700 px-6 py-3 text-sm font-semibold text-white transition hover:bg-brand-800"
			>Anunciar agora</a
		>
		{#if adminWhatsApp}
			<a
				href={`https://wa.me/${adminWhatsApp}`}
				target="_blank"
				rel="noopener noreferrer"
				class="rounded-full border border-brand-300 px-6 py-3 text-sm font-semibold text-brand-800 transition hover:bg-white"
				>Falar connosco</a
			>
		{/if}
	</div>
</section>
