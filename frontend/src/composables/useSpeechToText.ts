import type { SpeechRecognition, SpeechRecognitionEvent } from "@/@types/speech";
import { useStore } from "@/store/store";
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
	let processor: any = null
	let onStopCallback: (() => void) | null = null;
	let outputBuffer = ref<AudioBuffer | null>(null);
	let bufferQueue: Float32Array[] = [];

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
		audioContext = new AudioContext();
		mediaStream = await navigator.mediaDevices.getUserMedia({ audio: true });
		micStream = audioContext.createMediaStreamSource(mediaStream);

		processor = audioContext.createScriptProcessor(2048, 1, 1);
		processor.onaudioprocess = (event: AudioProcessingEvent) => {
			const input = event.inputBuffer.getChannelData(0);
			const inputCopy = new Float32Array(input.length);
			inputCopy.set(input);
			bufferQueue.push(inputCopy);
		};


		micStream.connect(processor);
		processor.connect(audioContext.destination);

		analyser = audioContext.createAnalyser();
		analyser.fftSize = 64;
		micStream.connect(analyser);
		updateVolume();
	}

	function replayAudio(outputBuffer: AudioBuffer) {
		const replayAudioContext = new AudioContext();
		const source = replayAudioContext.createBufferSource();
		source.buffer = outputBuffer;
		source.connect(replayAudioContext.destination);
		source.start();
		source.onended = () => {
			replayAudioContext.close();
		};
	}

	function stopMic() {
		if (animationFrameId) cancelAnimationFrame(animationFrameId);
		mediaStream?.getTracks().forEach(track => track.stop());
		if (audioContext && bufferQueue.length > 0) {
			const totalLength = bufferQueue.reduce((acc, buf) => acc + buf.length, 0);
			outputBuffer.value = audioContext.createBuffer(1, totalLength, audioContext.sampleRate);
			const combined = outputBuffer.value.getChannelData(0);

			let offset = 0;
			for (const buf of bufferQueue) {
				combined.set(buf, offset);
				offset += buf.length;
			}
		}

		// replayAudio();

		audioContext?.close();
		audioContext = null;
		micStream = null;
		analyser = null;
		bufferQueue = [];
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

	return { transcript, isListening, start, stop, replayAudio, outputBuffer, volume };
}
