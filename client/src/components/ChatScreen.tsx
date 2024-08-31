import React, { useState, useEffect, useRef } from 'react';

interface Message {
  text: string;
  isUser: boolean;
}

interface ChatScreenProps {
  setError: (error: string | null) => void;
}

const ChatScreen: React.FC<ChatScreenProps> = ({ setError }) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isRecording, setIsRecording] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const chatContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const initialText = localStorage.getItem('initialText');
    if (initialText) {
      setMessages([{ text: initialText, isUser: false }]);
      playAudio(localStorage.getItem('initialAudio'));
    }
  }, []);

  useEffect(() => {
    if (chatContainerRef.current) {
      chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
    }
  }, [messages]);

  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      mediaRecorderRef.current = new MediaRecorder(stream);
      
      const audioChunks: BlobPart[] = [];
      mediaRecorderRef.current.ondataavailable = (event) => {
        audioChunks.push(event.data);
      };

      mediaRecorderRef.current.onstop = () => {
        const audioBlob = new Blob(audioChunks, { type: 'audio/webm' });
        sendAudioToServer(audioBlob);
      };

      mediaRecorderRef.current.start();
      setIsRecording(true);
    } catch (error) {
      console.error('Error accessing microphone:', error);
      setError('Failed to access microphone. Please check your permissions and try again.');
    }
  };

  const stopRecording = () => {
    if (mediaRecorderRef.current && isRecording) {
      mediaRecorderRef.current.stop();
      setIsRecording(false);
    }
  };

  const sendAudioToServer = async (audioBlob: Blob) => {
    const formData = new FormData();
    formData.append('file', audioBlob, 'audio.webm');

    const id = localStorage.getItem('interviewId');
    const secret = localStorage.getItem('interviewSecret');
    const authString = btoa(`${id}:${secret}`);

    setIsProcessing(true);

    try {
      const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/chat/answer`, {
        method: 'POST',
        headers: {
          'Authorization': `Basic ${authString}`,
        },
        body: formData,
      });

      const data = await response.json();

      if (response.ok && data.data) {
        setMessages(prev => [
          ...prev,
          { text: data.data.prompt.text, isUser: true },
          { text: data.data.answer.text, isUser: false }
        ]);
        playAudio(data.data.answer.audio);
      } else {
        const errorMessage = data.message || 'Failed to process your response. Please try again.';
        setError(errorMessage);
      }
    } catch (error) {
      console.error('Error sending audio:', error);
      setError('Failed to send your response. Please check your connection and try again.');
    } finally {
      setIsProcessing(false);
    }
  };

  const playAudio = (base64Audio: string | null) => {
    if (base64Audio) {
      const audio = new Audio(`data:audio/mp3;base64,${base64Audio}`);
      audio.play();
    }
  };

  return (
    <div className="container mx-auto p-4 h-screen flex flex-col">
      <div ref={chatContainerRef} className="flex-grow overflow-y-auto mb-4 bg-dark-surface rounded-xl p-4 shadow-inner">
        {messages.map((message, index) => (
          <div key={index} className={`mb-4 ${message.isUser ? 'text-right' : 'text-left'}`}>
            <span className={`inline-block p-3 rounded-lg ${message.isUser ? 'bg-dark-primary text-dark-on-surface' : 'bg-dark-secondary text-dark-on-surface'}`}>
              {message.text}
            </span>
          </div>
        ))}
      </div>
      <button
        onClick={isRecording ? stopRecording : startRecording}
        disabled={isProcessing}
        className={`w-full p-4 rounded-xl font-bold text-lg transition-all duration-300 ${
          isProcessing
            ? 'bg-dark-secondary text-dark-on-surface cursor-not-allowed'
            : isRecording 
              ? 'bg-dark-error text-dark-on-surface animate-pulse' 
              : 'bg-dark-primary text-dark-on-surface hover:bg-opacity-90'
        }`}
      >
        {isProcessing 
          ? 'Processing...' 
          : isRecording 
            ? 'Stop Recording' 
            : 'Start Recording'
        }
      </button>
    </div>
  );
};

export default ChatScreen;