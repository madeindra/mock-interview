import React, { useState, useEffect, useRef } from 'react';

interface Message {
  text: string;
  isUser: boolean;
}

const ChatScreen: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isRecording, setIsRecording] = useState(false);
  const [audioBlob, setAudioBlob] = useState<Blob | null>(null);
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
      mediaRecorderRef.current.ondataavailable = (event) => {
        setAudioBlob(event.data);
      };
      mediaRecorderRef.current.start();
      setIsRecording(true);
    } catch (error) {
      console.error('Error accessing microphone:', error);
    }
  };

  const stopRecording = () => {
    if (mediaRecorderRef.current) {
      mediaRecorderRef.current.stop();
      setIsRecording(false);
    }
  };

  const sendAudioToServer = async () => {
    if (!audioBlob) return;

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

      if (response.ok) {
        const data = await response.json();
        setMessages(prev => [
          ...prev,
          { text: data.prompt.text, isUser: true },
          { text: data.answer.text, isUser: false }
        ]);
        playAudio(data.answer.audio);
      } else {
        console.error('Failed to send audio');
      }
    } catch (error) {
      console.error('Error sending audio:', error);
    }
  };

  const playAudio = (base64Audio: string | null) => {
    if (base64Audio) {
      const audio = new Audio(`data:audio/mp3;base64,${base64Audio}`);
      audio.play();
    }
  };

  return (
    <div className="container mx-auto mt-10 p-4">
      <div className="mb-4 h-96 overflow-y-auto border rounded p-4">
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
      {audioBlob && !isRecording && (
        <button
          onClick={sendAudioToServer}
          className="w-full mt-2 p-2 rounded bg-blue-500 text-white"
        >
          Send Response
        </button>
      )}
    </div>
  );
};

export default ChatScreen;