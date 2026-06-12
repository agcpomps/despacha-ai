<script lang="ts">
	let { class: className = 'text-2xl' }: { class?: string } = $props();
</script>

<!-- O tamanho do lockup acompanha o font-size: usa class="text-3xl" etc. -->
<span class="inline-flex items-center gap-[0.3em] leading-none {className}">
	<svg
		class="h-[1.5em] w-auto shrink-0"
		viewBox="0 0 72 72"
		fill="none"
		role="img"
		aria-label="Despacha Aí"
	>
		<!-- saco de trás (camada clara) -->
		<path
			d="M28 30 L55 30 Q58 30 58.4 33 L61 60 Q61.6 66 55 66 L34 66 Q27.4 66 28 60 Z"
			fill="#8dd6a1"
		/>
		<!-- alça -->
		<path
			d="M26 32 V21 a10 10 0 0 1 20 0 v6"
			stroke="#157033"
			stroke-width="5.5"
			stroke-linecap="round"
		/>
		<!-- saco da frente -->
		<path
			d="M17 28 L48 28 Q51 28 51.4 31 L54.5 59.5 Q55.2 66 48.5 66 L16.5 66 Q9.8 66 10.5 59.5 L13.6 31 Q14 28 17 28 Z"
			fill="#22a04a"
		/>
		<!-- seta (swoosh) -->
		<path d="M13 60 C30 56 42 42 49 24" stroke="#f0f7f0" stroke-width="6" stroke-linecap="round" />
		<!-- ponta da seta -->
		<path d="M42 16 L58 10 L53 26 Z" fill="#8dd6a1" />
	</svg>
	<span class="font-extrabold tracking-tight whitespace-nowrap text-brand-700">
		Despacha <span class="text-brand-500">Aí</span>
	</span>
</span>
