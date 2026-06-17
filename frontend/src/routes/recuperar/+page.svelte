<script lang="ts">
	import { env } from '$env/dynamic/public';
	import Seo from '$lib/components/Seo.svelte';

	const adminWhatsApp = (env.PUBLIC_ADMIN_WHATSAPP ?? '').replace(/\D/g, '');

	let phone = $state('');

	const link = $derived.by(() => {
		const text = encodeURIComponent(
			`Olá! Esqueci-me da palavra-passe da minha conta no Despacha Aí.\nO meu telemóvel de registo é: ${phone || '(escreve aqui)'}`
		);
		return `https://wa.me/${adminWhatsApp}?text=${text}`;
	});
</script>

<Seo
	title="Recuperar palavra-passe | Despacha Aí"
	description="Recupera o acesso à tua conta no Despacha Aí."
/>

<div class="mx-auto mt-6 max-w-md">
	<div class="rounded-2xl border border-neutral-200 bg-white p-8">
		<h1 class="text-2xl font-bold text-neutral-900">Recuperar palavra-passe</h1>
		<p class="mt-2 text-sm leading-relaxed text-neutral-600">
			Indica o telemóvel com que te registaste e fala connosco pelo WhatsApp. Vamos verificar a tua
			conta e enviar-te uma nova palavra-passe.
		</p>

		<label class="mt-6 block">
			<span class="text-sm font-semibold text-neutral-700">Telemóvel de registo</span>
			<input
				type="tel"
				bind:value={phone}
				placeholder="923 456 789"
				autocomplete="tel"
				class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
			/>
		</label>

		{#if adminWhatsApp}
			<a
				href={link}
				target="_blank"
				rel="noopener noreferrer"
				class="mt-5 flex h-12 w-full items-center justify-center gap-2 rounded-full bg-whatsapp text-sm font-bold text-white transition hover:brightness-95"
			>
				<svg class="h-5 w-5" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
					<path
						d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347Z"
					/>
				</svg>
				Contactar pelo WhatsApp
			</a>
		{:else}
			<p class="mt-5 rounded-xl bg-neutral-50 p-3 text-sm text-neutral-500">
				Contacta o suporte para repor a tua palavra-passe.
			</p>
		{/if}

		<p class="mt-6 text-center text-sm text-neutral-500">
			<a href="/login" class="font-semibold text-brand-700 hover:text-brand-800">Voltar ao início de sessão</a>
		</p>
	</div>
</div>
