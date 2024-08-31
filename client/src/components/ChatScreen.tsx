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
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);

  useEffect(() => {
    const initialText = localStorage.getItem('initialText');
    if (initialText) {
      setMessages([{ text: initialText, isUser: false }]);
      playAudio(localStorage.getItem('initialAudio'));
    }
  }, []);

  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      mediaRecorderRef.current = new MediaRecorder(stream);
      
      const audioChunks: BlobPart[] = [];
      mediaRecorderRef.current.ondataavailable = (event) => {
        audioChunks.push(event.data);
      };

      mediaRecorderRef.current.onstop = () => {
        const audioBlob = new Blob(audioChunks, { type: 'audio/wav' });
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
    formData.append('file', audioBlob, 'audio.wav');

    const id = localStorage.getItem('interviewId');
    const secret = localStorage.getItem('interviewSecret');
    const authString = btoa(`${id}:${secret}`);

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
    }
  };

  const playAudio = (base64Audio: string | null) => {
    if (base64Audio) {
      const audio = new Audio(`data:audio/mp3;base64,${base64Audio}`);
      audio.play();
    }
  };

  return (
    <div className="container mx-auto p-4">
      <div className="mb-4 h-[calc(100vh-8rem)] overflow-y-auto border rounded p-4">
        {messages.map((message, index) => (
          <div key={index} className={`mb-2 ${message.isUser ? 'text-right' : 'text-left'}`}>
            <span className={`inline-block p-2 rounded ${message.isUser ? 'bg-blue-200' : 'bg-gray-200'}`}>
              {message.text}
            </span>
          </div>
        ))}
      </div>
      <button
        onClick={isRecording ? stopRecording : startRecording}
        className={`w-full p-2 rounded ${isRecording ? 'bg-red-500' : 'bg-green-500'} text-white`}
      >
        {isRecording ? 'Stop Recording' : 'Start Recording'}
      </button>
    </div>
  );
};

export default ChatScreen;