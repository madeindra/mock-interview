import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import Navbar from './Navbar';

interface Message {
  text: string;
  isUser: boolean;
  displayedText?: string;
}

interface ChatScreenProps {
  backendHost: string;
  setError: (error: string | null) => void;
}

const languageEN: string = "en";
const languageID: string = "id";

const languageMap: Record<string, string> = {
  [languageEN]: "en-US",
  [languageID]: "id-ID"
};

const ChatScreen: React.FC<ChatScreenProps> = ({ backendHost, setError }) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isRecording, setIsRecording] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const [isTyping, setIsTyping] = useState(false);
  const [hasStarted, setHasStarted] = useState(false);
  const [hasEnded, setHasEnded] = useState(false);

  const navigate = useNavigate();

  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const chatContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const initialText = sessionStorage.getItem('initialText');
    const storedMessages = sessionStorage.getItem('messages');
    const chatLanguage = sessionStorage.getItem('chatLanguage') || languageEN;

    if (storedMessages) {
      setMessages(JSON.parse(storedMessages));
    } else if (initialText) {
      const initialMessage = { text: initialText, isUser: false, displayedText: '' };
      setMessages([initialMessage]);
      typeMessage(0, initialText);
      if (chatLanguage == languageEN) { 
        playAudio(sessionStorage.getItem('initialAudio'));
      } else {
        synthesizeText(initialText, languageMap[chatLanguage]);
      }
      sessionStorage.setItem('messages', JSON.stringify([initialMessage]));
    } else {
      navigate('/');
    }
  }, [navigate]);

  useEffect(() => {
    if (messages.length > 0) {
      sessionStorage.setItem('messages', JSON.stringify(messages));
    }
  }, [messages]);

  useEffect(() => {
    if (chatContainerRef.current) {
      chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
    }
  }, [messages]);

  const typeMessage = (messageIndex: number, text: string) => {
    setIsTyping(true);

    let i = 0;
    const typingInterval = setInterval(() => {
      setMessages(prevMessages => {
        const newMessages = [...prevMessages];
        if (newMessages[messageIndex]) {
          newMessages[messageIndex] = {
            ...newMessages[messageIndex],
            displayedText: text.slice(0, i)
          };
        }
        return newMessages;
      });

      i++;

      if (i > text.length) {
        clearInterval(typingInterval);
        setIsTyping(false);
      }
    }, 25);
  };

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

    const id = sessionStorage.getItem('interviewId');
    const secret = sessionStorage.getItem('interviewSecret');
    const authString = btoa(`${id}:${secret}`);

    setIsProcessing(true);

    try {
      const response = await fetch(`${backendHost}/chat/answer`, {
        method: 'POST',
        headers: {
          'Authorization': `Basic ${authString}`,
        },
        body: formData,
      });

      const data = await response.json();

      if (response.ok && data.data) {
        const userMessage: Message = { text: data.data.prompt.text, isUser: true };
        const botMessage: Message = { text: data.data.answer.text, isUser: false, displayedText: '' };

        setMessages(prev => [...prev, userMessage, botMessage]);

        // Start typing effect for the bot message
        requestAnimationFrame(() => {
          typeMessage(messages.length + 1, data.data.answer.text);
        });

        if (data.data.language == languageEN) { 
          playAudio(data.data.answer.audio);
        } else {
          synthesizeText(data.data.answer.text, languageMap[data?.data?.language]);
        }

        setHasStarted(true);
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

  const synthesizeText = async (text: string, language: string) => {
    const audio = new SpeechSynthesisUtterance(text);
    audio.lang = language;
    window.speechSynthesis.speak(audio);
  }

  const endInterview = async () => {
    const id = sessionStorage.getItem('interviewId');
    const secret = sessionStorage.getItem('interviewSecret');
    const authString = btoa(`${id}:${secret}`);

    setIsProcessing(true);

    try {
      const response = await fetch(`${backendHost}/chat/end`, {
        method: 'GET',
        headers: {
          'Authorization': `Basic ${authString}`,
        },
      });

      const data = await response.json();

      if (response.ok && data.data) {
        const botMessage: Message = { text: data.data.answer.text, isUser: false, displayedText: '' };

        setMessages(prevMessages => {
          const newMessages = [...prevMessages, botMessage];
          // Start typing effect for the bot message after state update
          setTimeout(() => {
            typeMessage(newMessages.length - 1, data.data.answer.text);
          }, 0);
          return newMessages;
        });

        if (data.data.language == languageEN) { 
          playAudio(data.data.answer.audio);
        } else {
          synthesizeText(data.data.answer.text, languageMap[data?.data?.language]);
        }
      } else {
        setError(data.message || 'Failed to end the interview. Please try again.');
      }
    } catch (error) {
      console.error('Error ending interview:', error);
      setError('Failed to end the interview. Please check your connection and try again.');
    } finally {
      setHasEnded(true);
      setIsProcessing(false);
    }
  };

  const handleStartOver = () => {
    sessionStorage.removeItem('messages');
    sessionStorage.removeItem('initialText');
    sessionStorage.removeItem('initialAudio');
    sessionStorage.removeItem('chatLanguage');
    navigate('/');
  };

  const handleBack = () => {
    navigate('/');
  };

  return (
    <div className="flex flex-col h-screen bg-[#1E1E2E] text-white">
      <Navbar
        backendHost={backendHost}
        showBackIcon
        showForwardIcon
        showStartOver
        onBack={handleBack}
        onStartOver={handleStartOver}
        disableForward={true}
      />
      <div ref={chatContainerRef} className="flex-grow overflow-y-auto px-4 py-2">
        {messages.map((message, index) => (
          <div key={index} className={`mb-4 ${message.isUser ? 'text-right' : 'text-left'}`}>
            <span className={`inline-block p-3 rounded-2xl ${message.isUser
                ? 'bg-[#3E64FF] text-white'
                : 'bg-[#2B2B3B] text-white'
              }`}>
              {message.isUser ? message.text : (message.displayedText || '')}
            </span>
          </div>
        ))}
      </div>

      <div className="flex justify-between items-center space-x-4 p-4 bg-[#1E1E2E]">
        <button
          onClick={isRecording ? stopRecording : startRecording}
          disabled={isProcessing || isTyping || hasEnded}
          className={`w-full p-4 rounded-xl font-bold text-lg transition-all duration-300 ${isProcessing || isTyping || hasEnded
              ? 'bg-[#2B2B3B] text-gray-400 cursor-not-allowed'
              : isRecording
                ? 'bg-[#FF3E3E] text-white animate-pulse'
                : 'bg-[#3E64FF] text-white hover:bg-opacity-90'
            }`}
        >
          {isProcessing
            ? 'Processing...'
            : isTyping
              ? 'Responding...'
              : isRecording
                ? 'Stop Recording'
                : hasEnded
                  ? 'This interview has ended'
                  : 'Start Recording'
          }
        </button>
        {hasStarted && !hasEnded && (
          <button
            onClick={endInterview}
            disabled={isProcessing || isTyping || isRecording || hasEnded}
            className={`w-3/12 p-4 rounded-xl font-bold text-lg hover:bg-opacity-90 transition-all duration-300 ${isProcessing || isTyping || isRecording || hasEnded
              ? 'bg-[#2B2B3B] text-gray-400 cursor-not-allowed'
                : 'bg-[#FF3E3E] text-white hover:bg-opacity-90'
            }`}
          >
            End
          </button>
        )}
      </div>
    </div>
  );
};

export default ChatScreen;