import type { SpeechRecognition, SpeechRecognitionEvent } from "@/@types/speech";
import { ref, onMounted, onUnmounted } from "vue";

export function useSpeechToText() {
	const transcript = ref("");
	const isListening = ref(false);
	const volume = ref(0);
	let recognition: SpeechRecognition | null = null;
	let audioContext: AudioContext | null = null;
	let analyser: AnalyserNode | null = null;
	let micStream: MediaStreamAudioSourceNode | null = null;
	let mediaStream: MediaStream | null = null;
	let animationFrameId: number;
	let onStopCallback: (() => void) | null = null;

	function updateVolume() {
		if (analyser) {
			const data = new Uint8Array(analyser.frequencyBinCount);
			analyser.getByteFrequencyData(data);
			const avg = data.reduce((sum, val) => sum + val, 0) / data.length;
			volume.value = avg / 255;
		}
		animationFrameId = requestAnimationFrame(updateVolume);
	}

	async function startMic() {
		console.log("Starting microphone...");
		audioContext = new AudioContext();
		mediaStream = await navigator.mediaDevices.getUserMedia({ audio: true });
		micStream = audioContext.createMediaStreamSource(mediaStream);
		analyser = audioContext.createAnalyser();
		analyser.fftSize = 64;
		micStream.connect(analyser);
		updateVolume();
		console.log("Microphone started.");
	}

	function stopMic() {
		if (animationFrameId) cancelAnimationFrame(animationFrameId);
		mediaStream?.getTracks().forEach(track => track.stop());
		audioContext?.close();
		audioContext = null;
		micStream = null;
		analyser = null;
		mediaStream = null;
		volume.value = 0;
	}

	onMounted(() => {
		const SpeechRecognition =
			(window as any).SpeechRecognition || (window as any).webkitSpeechRecognition;
		if (!SpeechRecognition) {
			console.log("Web Speech API is not supported in this browser.");
			return;
		}

		recognition = new SpeechRecognition() as SpeechRecognition;
		recognition.lang = "en-US";
		recognition.continuous = true;
		recognition.interimResults = true;

		recognition.onresult = (event: SpeechRecognitionEvent) => {
			let interim = "";
			for (let i = event.resultIndex; i < event.results.length; ++i) {
				const result = event.results[i];
				if (result.isFinal) {
					transcript.value += result[0].transcript + " ";
				} else {
					interim += result[0].transcript;
				}
			}
		};

		recognition.onstart = () => isListening.value = true;
		recognition.onend = () => {
			isListening.value = false;
			stopMic();
			if (onStopCallback) {
				onStopCallback();
				onStopCallback = null;
			}
		};
	});

	onUnmounted(() => {
		recognition?.stop();
		stopMic();
	});

	const start = async () => {
		transcript.value = "";
		await startMic();
		recognition?.start();
	};

	const stop = (callback?: () => void) => {
		onStopCallback = callback || null;
		recognition?.stop();
	};

	return { transcript, isListening, start, stop, volume };
}
